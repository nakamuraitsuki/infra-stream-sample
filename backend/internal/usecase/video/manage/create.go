package manage

import (
	"context"
	"time"

	"example.com/m/internal/domain/event"
	video_domain "example.com/m/internal/domain/video"
	video_value "example.com/m/internal/domain/video/value"
	"github.com/google/uuid"
)

type CreateResponse struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      string
	Tags        []string
	Visibility  string
	CreatedAt   time.Time
}

func (uc *VideoManagementUseCase) Create(
	ctx context.Context,
	ownerID uuid.UUID,
	title string,
	description string,
	tagsStr []string,
) (*CreateResponse, error) {

	videoID := uuid.New()

	status := video_value.StatusInitial

	// NOTE: sourceKey and streamKey will be set after the video data is uploaded.
	sourceKey := ""
	streamKey := ""

	visibility := video_value.VisibilityPublic

	tags := make([]video_value.Tag, len(tagsStr))
	for i, tagStr := range tagsStr {
		tag, err := video_value.NewTag(tagStr)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}

	video := video_domain.NewVideo(
		videoID,
		ownerID,
		sourceKey,
		streamKey,
		status,
		title,
		description,
		tags,
		0,
		nil,
		visibility,
		time.Now(),
		[]event.Event{},
	)

	if err := uc.VideoRepo.Save(ctx, video); err != nil {
		return nil, err
	}

	return &CreateResponse{
		ID:          video.ID(),
		Title:       video.Title(),
		Description: video.Description(),
		Status:      string(video.Status()),
		Tags:        tagsStr,
		Visibility:  string(video.Visibility()),
		CreatedAt:   video.CreatedAt(),
	}, nil
}
