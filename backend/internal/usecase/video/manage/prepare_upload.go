package manage

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
)

const (
	// NOTE: 5MB程度を最小パーツサイズとする。
	// S3の仕様を基準にしているが、デファクトスタンダードなのでインフラ依存ではないと捉えている
	MinPartSize = 5 * 1024 * 1024 // 5MB

	// NOTE: 最大パーツ数は10,000。
	// 10,000パーツ × 5MB = 50GBまでの動画に対応可能。これ以上のサイズの動画はサポートしない方針。
	MaxPartCount = 10000
)

type PrepareUploadResponse struct {
	UploadID string
	URLs     []string
	PartSize int64
	Key      string
}

// TODO: 実行者が保有者であることの検証を追加する
func (uc *VideoManagementUseCase) PrepareUploadSession(
	ctx context.Context,
	videoID uuid.UUID,
	fileSize int64,
) (*PrepareUploadResponse, error) {
	// NOTE: sourceKey は ID から一意に決まるので、保存完了までDomainには反映させない。
	sourceKey := fmt.Sprintf("videos/%s/source", videoID.String())

	partCount := int32(math.Ceil(float64(fileSize) / float64(MinPartSize)))
	if partCount > MaxPartCount {
		return nil, ErrFileSizeTooLarge
	}

	sessionID, err := uc.Storage.StartUploadSession(ctx, sourceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to start upload session: %w", err)
	}

	urls := make([]string, partCount)
	for i := int32(1); i <= partCount; i++ {
		url, err := uc.Storage.GenerateUploadPartURL(
			ctx,
			sourceKey,
			sessionID,
			i,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate upload part URL for part %d: %w", i, err)
		}
		urls[i-1] = url
	}

	return &PrepareUploadResponse{
		UploadID: sessionID,
		URLs:     urls,
		PartSize: MinPartSize,
		Key:      sourceKey,
	}, nil
}
