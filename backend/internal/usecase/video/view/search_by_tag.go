package view

import (
	"context"
	"time"

	video_domain "example.com/m/internal/domain/video"
	video_value "example.com/m/internal/domain/video/value"
	"example.com/m/internal/usecase/video/query"
	"github.com/google/uuid"
)

type GetByTagsResults struct {
	Results []*GetByTagsResult
}

type GetByTagsResult struct {
	ID          uuid.UUID
	Title       string
	Description string
	Tags        []string
	OwnerID     uuid.UUID
	CreatedAt   time.Time
}

func (uc *VideoViewingUseCase) SearchByTag(
	ctx context.Context,
	tagStr string,
	query query.VideoSearchQuery,
) (*GetByTagsResults, error) {

	visibility := video_value.VisibilityPublic
	status := video_value.StatusReady

	tag, err := video_value.NewTag(tagStr)
	if err != nil {
		return nil, err
	}

	cond := video_domain.ListCondition{
		Tag:        &tag,
		Visibility: &visibility,
		Status:     &status,
		Limit:      query.Limit,
	}

	videos, err := uc.VideoRepo.FindByCondition(ctx, cond)
	if err != nil {
		return nil, err
	}

	results := make([]*GetByTagsResult, len(videos))
	for i, video := range videos {
		tagsStr := make([]string, len(video.Tags()))
		for j, t := range video.Tags() {
			tagsStr[j] = string(t)
		}
		results[i] = &GetByTagsResult{
			ID:          video.ID(),
			Title:       video.Title(),
			Description: video.Description(),
			Tags:        tagsStr,
			OwnerID:     video.OwnerID(),
			CreatedAt:   video.CreatedAt(),
		}
	}

	return &GetByTagsResults{
		Results: results,
	}, nil
}
