package process

import (
	"context"

	video_value "example.com/m/internal/domain/video/value"
	"github.com/google/uuid"
)

func (uc *VideoProcessUseCase) Handle(ctx context.Context, videoID uuid.UUID, isFinalAttempt bool) error {
	video, err := uc.VideoRepo.FindByID(ctx, videoID)
	if err != nil {
		return err
	}

	if video.Status() != video_value.StatusProcessing {
		return nil
	}

	transcodeErr := uc.Transcoder.Transcode(ctx, video.SourceKey(), video.StreamKey())

	return uc.UoW.Do(ctx, func(ctx context.Context) error {
		v, err := uc.VideoRepo.FindByID(ctx, videoID)
		if err != nil {
			return err
		}

		if transcodeErr != nil {
			if isFinalAttempt {
				v.MarkTranscodeFailed(video_value.FailureTranscode)
			} else {
				return transcodeErr // 再試行のためにエラーを返す
			}
		} else {
			v.MarkTranscodeSucceeded()
		}

		return uc.VideoRepo.Save(ctx, v)
	})
}
