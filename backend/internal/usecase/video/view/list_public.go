package view

import (
	"context"
	"time"

	video_domain "example.com/m/internal/domain/video"
	video_value "example.com/m/internal/domain/video/value"
	"example.com/m/internal/usecase/video/query"
	"github.com/google/uuid"
)

type ListPublicResult struct {
	Videos []*PublicVideo
}

type PublicVideo struct {
	ID          uuid.UUID
	Title       string
	Description string
	Tags        []string
	OwnerID     uuid.UUID
	CreatedAt   time.Time
}

func (uc *VideoViewingUseCase) ListPublic(
	ctx context.Context,
	query query.VideoSearchQuery,
) (*ListPublicResult, error) {

	visibility := video_value.VisibilityPublic
	status := video_value.StatusReady

	cond := video_domain.ListCondition{
		Visibility: &visibility,
		Status:     &status,
		Limit:      query.Limit,
	}

	videos, err := uc.VideoRepo.FindByCondition(ctx, cond)
	if err != nil {
		return nil, err
	}

	publicVideos := make([]*PublicVideo, len(videos))
	for i, video := range videos {
		tagsStr := make([]string, len(video.Tags()))
		for j, tag := range video.Tags() {
			tagsStr[j] = string(tag)
		}

		publicVideos[i] = &PublicVideo{
			ID:          video.ID(),
			Title:       video.Title(),
			Description: video.Description(),
			Tags:        tagsStr,
			OwnerID:     video.OwnerID(),
			CreatedAt:   video.CreatedAt(),
		}
	}

	return &ListPublicResult{
		Videos: publicVideos,
	}, nil
}
