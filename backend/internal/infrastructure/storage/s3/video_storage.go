package s3

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"example.com/m/internal/domain/video"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type videoStorage struct {
  client              *s3.Client
  internalPresignClient *s3.PresignClient // アップロード用（内部NW向け）
  publicPresignClient   *s3.PresignClient   // 視聴用（NGINX/外部公開向け）
  bucketName          string
}

func NewVideoStorage(clientSet *S3ClientSet, cfg Config) video.Storage {
	return &videoStorage{
		client:        clientSet.Client,
		internalPresignClient: s3.NewPresignClient(clientSet.Client),
		publicPresignClient: clientSet.PresignClient,
		bucketName:    cfg.BucketName,
	}
}

func (s *videoStorage) StartUploadSession(ctx context.Context, key string) (string, error) {
	out, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String("video/mp4"),
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(out.UploadId), nil
}

func (s *videoStorage) GenerateUploadPartURL(
	ctx context.Context,
	key string,
	sessionID string,
	partNumber int32,
) (string, error) {
	// NOTE: 有効期限15分の署名付きURLを生成
	// アップロードはフロントエンド側に任せる想定
	presignedReq, err := s.publicPresignClient.PresignUploadPart(ctx, &s3.UploadPartInput{
		Bucket:     aws.String(s.bucketName),
		Key:        aws.String(key),
		UploadId:   aws.String(sessionID),
		PartNumber: aws.Int32(partNumber),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return "", err
	}
	return presignedReq.URL, nil
}

func (s *videoStorage) CommitUploadSession(
	ctx context.Context,
	key string,
	sessionID string,
	parts []video.PartInfo,
) error {
	completedParts := make([]types.CompletedPart, len(parts))
	for i, p := range parts {
		completedParts[i] = types.CompletedPart{
			ETag:       aws.String(p.ID), // IDはETagとして扱う
			PartNumber: aws.Int32(p.PartNumber),
		}
	}

	_, err := s.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(s.bucketName),
		Key:      aws.String(key),
		UploadId: aws.String(sessionID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	return err
}

func (s *videoStorage) SaveStream(ctx context.Context, streamKey string, data io.Reader) error {
	contentType := s.detectContentType(streamKey)
	return s.upload(ctx, streamKey, data, contentType)
}

func (s *videoStorage) GenerateTemporaryAccessURL(ctx context.Context, streamKey string, expiresDuration time.Duration) (string, error) {
	// NOTE: Proxyからの X-Accell-Redirect を想定し、Nginxを含む内部NWむけの PresignClient を使用して署名付きURLを生成
	presignedReq, err := s.internalPresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(streamKey),
	}, s3.WithPresignExpires(expiresDuration))

	if err != nil {
		return "", err
	}
	return presignedReq.URL, nil
}

func (s *videoStorage) GetStream(ctx context.Context, streamKey string, byteRange *video.ByteRange) (io.ReadCloser, *video.ObjectMeta, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(streamKey),
	}

	if byteRange != nil {
		if byteRange.End != nil {
			input.Range = aws.String(
				fmt.Sprintf("bytes=%d-%d", byteRange.Start, *byteRange.End),
			)
		} else {
			input.Range = aws.String(
				fmt.Sprintf("bytes=%d-", byteRange.Start),
			)
		}
	}

	out, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	meta := &video.ObjectMeta{
		ContentLength: aws.ToInt64(out.ContentLength),
		ETag:          aws.ToString(out.ETag),
		LastModified:  aws.ToTime(out.LastModified),
	}
	// Range 系の返答Metaをoutから組み立て
	if out.ContentRange != nil {
		// 例: "bytes 0-1023/5000000"
		var start, end, total int64
		_, err := fmt.Sscanf(
			aws.ToString(out.ContentRange),
			"bytes %d-%d/%d",
			&start,
			&end,
			&total,
		)
		if err == nil {
			meta.RangeStart = start
			meta.RangeEnd = end
			meta.TotalSize = total
		} else {
			// フォーマット不正などで Content-Range のパースに失敗した場合は
			// Content-Range が無い場合と同様に全体サイズからレンジを推定する
			meta.TotalSize = aws.ToInt64(out.ContentLength)
			meta.RangeStart = 0
			if meta.TotalSize > 0 {
				meta.RangeEnd = meta.TotalSize - 1
			} else {
				meta.RangeEnd = 0
			}
		}
	} else {
		meta.TotalSize = aws.ToInt64(out.ContentLength)
		meta.RangeStart = 0
		meta.RangeEnd = meta.TotalSize - 1
	}

	return out.Body, meta, nil
}

func (s *videoStorage) GetSource(ctx context.Context, sourceKey string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(sourceKey),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (s *videoStorage) DeleteSource(ctx context.Context, sourceKey string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(sourceKey),
	})
	return err
}

// streamKey は Prefix を渡し、配下のオブジェクトを全て削除する
func (s *videoStorage) DeleteStream(ctx context.Context, streamKey string) error {
	// cf. https://qiita.com/hkak03key/items/410fda8f008a21a4ca43
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(streamKey),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		if len(page.Contents) == 0 {
			continue
		}

		var objects []types.ObjectIdentifier
		for _, obj := range page.Contents {
			objects = append(
				objects,
				types.ObjectIdentifier{Key: obj.Key},
			)
		}

		_, err = s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(s.bucketName),
			Delete: &types.Delete{Objects: objects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper method to upload data to S3
func (s *videoStorage) upload(ctx context.Context, key string, data io.Reader, contentType string) error {

	uploader := transfermanager.New(s.client)
	_, err := uploader.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        data,
		ContentType: aws.String(contentType),
	})

	return err
}

func (s *videoStorage) detectContentType(key string) string {
	switch {
	case strings.HasSuffix(key, ".m3u8"):
		return "application/vnd.apple.mpegurl"
	case strings.HasSuffix(key, ".ts"):
		return "video/mp2t"
	case strings.HasSuffix(key, ".mp4"):
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}
