package process

import (
	"context"

	"example.com/m/internal/domain/video"
	"example.com/m/internal/usecase/tx"
	"github.com/google/uuid"
)

type VideoProcessUseCaseInterface interface {
	// Handle processes the video with the given ID.
	Handle(ctx context.Context, videoID uuid.UUID, isFinalAttempt bool) error
}

type VideoProcessUseCase struct {
	VideoRepo  video.Repository
	Transcoder video.Transcoder
	UoW        tx.UnitOfWork
}

func NewVideoProcessUseCase(
	videoRepo video.Repository,
	transcoder video.Transcoder,
	uow tx.UnitOfWork,
) VideoProcessUseCaseInterface {
	return &VideoProcessUseCase{
		VideoRepo:  videoRepo,
		Transcoder: transcoder,
		UoW:        uow,
	}
}
