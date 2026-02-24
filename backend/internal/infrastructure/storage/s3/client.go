package s3

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3ClientSet struct {
	Client        *s3.Client        // 内部クライアント
	PresignClient *s3.PresignClient // 外部署名用
}

func NewClient(ctx context.Context, cfg Config) (*S3ClientSet, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKey,
				cfg.SecretKey,
				"", // session token（通常不要）
			),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
		o.UsePathStyle = cfg.UsePathStyle
	})

	presignS3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.PublicEndpoint != "" {
			o.BaseEndpoint = aws.String(cfg.PublicEndpoint)
		}
		o.UsePathStyle = cfg.UsePathStyle
	})
	presignClient := s3.NewPresignClient(presignS3Client)

	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(cfg.BucketName),
	})
	if err != nil {
		var ownedErr *types.BucketAlreadyOwnedByYou
		var existsErr *types.BucketAlreadyExists
		if !errors.As(err, &ownedErr) && !errors.As(err, &existsErr) {
			return nil, err
		}
	}

	return &S3ClientSet{
		Client:        client,
		PresignClient: presignClient,
	}, nil
}
