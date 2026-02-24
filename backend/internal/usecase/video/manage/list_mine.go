package manage

import (
	"context"
	"time"

	video_domain "example.com/m/internal/domain/video"
	"example.com/m/internal/usecase/video/query"
	"github.com/google/uuid"
)

type ListMineResults struct {
	Items []ListMineResultItem
}

type ListMineResultItem struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      string
	Tags        []string
	Visibility  string
	CreatedAt   time.Time
}

func (uc *VideoManagementUseCase) ListMine(
	ctx context.Context,
	ownerID uuid.UUID,
	query query.VideoSearchQuery,
) (*ListMineResults, error) {

	cond := video_domain.ListCondition{
		OwnerID: &ownerID,
		Limit:   query.Limit,
	}

	videos, err := uc.VideoRepo.FindByCondition(ctx, cond)
	if err != nil {
		return nil, err
	}

	items := make([]ListMineResultItem, len(videos))
	for i, video := range videos {
		tagsStr := make([]string, len(video.Tags()))
		for j, tag := range video.Tags() {
			tagsStr[j] = string(tag)
		}

		items[i] = ListMineResultItem{
			ID:          video.ID(),
			Title:       video.Title(),
			Description: video.Description(),
			Status:      string(video.Status()),
			Tags:        tagsStr,
			Visibility:  string(video.Visibility()),
			CreatedAt:   video.CreatedAt(),
		}
	}
	return &ListMineResults{Items: items}, nil
}
