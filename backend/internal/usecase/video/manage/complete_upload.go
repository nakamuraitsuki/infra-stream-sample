package manage

import (
	"context"
	"fmt"

	"example.com/m/internal/domain/video"
	video_value "example.com/m/internal/domain/video/value"
	"github.com/google/uuid"
)

type UploadPart struct {
	PartNumber int32
	ETag       string
}

type CompleteUploadRequest struct {
	VideoID  uuid.UUID
	UploadID string
	Parts    []UploadPart
}

// TODO: 実行者が保有者であることの検証を追加する
func (uc *VideoManagementUseCase) CompleteUploadSession(
	ctx context.Context,
	req CompleteUploadRequest,
) error {
	// NOTE: sourceKey は ID から一意に決まるので、保存完了までDomainには反映させない。
	sourceKey := fmt.Sprintf("videos/%s/source", req.VideoID.String())

	parts := make([]video.PartInfo, len(req.Parts))
	for i, part := range req.Parts {
		parts[i] = video.PartInfo{
			PartNumber: part.PartNumber,
			ID:         part.ETag, // ETagをIDとして扱う
		}
	}

	err := uc.Storage.CommitUploadSession(
		ctx,
		sourceKey,
		req.UploadID,
		parts,
	)
	if err != nil {
		return fmt.Errorf("failed to commit upload session: %w", err)
	}

	err = uc.UoW.Do(ctx, func(ctx context.Context) error {
		video, err := uc.VideoRepo.FindByID(ctx, req.VideoID)
		if err != nil {
			return err
		}

		if video.Status() != video_value.StatusInitial {
			return nil
		}

		if err := video.MarkUploaded(sourceKey); err != nil {
			return err
		}

		streamKey := fmt.Sprintf(
			"videos/%s/stream/",
			req.VideoID.String(),
		)

		// アップロード完了したら、非同期Workerにジョブを投げていい
		if err = video.StartTranscoding(streamKey); err != nil {
			return err
		}

		if err = uc.VideoRepo.Save(ctx, video); err != nil {
			return err
		}

		ev := video.PullEvents()
		return uc.OutboxRepo.Save(ctx, ev)
	})

	if err != nil {
		// DBが失敗した時だけ、保存してしまったファイルを消す
		_ = uc.Storage.DeleteSource(ctx, sourceKey)
		return err
	}

	return nil
}
