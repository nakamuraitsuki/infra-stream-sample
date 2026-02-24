package view

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GetByIDResult struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Title       string
	Description string
	Tags        []string
	Visibility  string
	CreatedAt   time.Time
}

func (uc *VideoViewingUseCase) GetByID(
	ctx context.Context,
	videoID uuid.UUID,
) (*GetByIDResult, error) {

	video, err := uc.VideoRepo.FindByID(ctx, videoID)
	if err != nil {
		return nil, err
	}

	if video.Status() != "ready" {
		return nil, ErrVideoNotReady
	}

	if video.Visibility() != "public" {
		return nil, ErrVideoForbidden
	}

	tagsStr := make([]string, len(video.Tags()))
	for i, tag := range video.Tags() {
		tagsStr[i] = string(tag)
	}

	return &GetByIDResult{
		ID:          video.ID(),
		OwnerID:     video.OwnerID(),
		Title:       video.Title(),
		Description: video.Description(),
		Tags:        tagsStr,
		Visibility:  string(video.Visibility()),
		CreatedAt:   video.CreatedAt(),
	}, nil
}
