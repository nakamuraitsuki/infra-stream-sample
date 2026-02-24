package view

import (
	"context"
	"path"
	"strings"
	"time"

	video_value "example.com/m/internal/domain/video/value"
	"github.com/google/uuid"
)

func (uc *VideoViewingUseCase) GetVideoStream(
	ctx context.Context,
	videoID uuid.UUID,
	objectPath string,
) (string, error) {

	video, err := uc.VideoRepo.FindByID(ctx, videoID)
	if err != nil {
		return "", err
	}

	if video.Status() != video_value.StatusReady {
		return "", ErrVideoNotReady
	}

	// Check visibility
	if video.Visibility() != video_value.VisibilityPublic {
		return "", ErrVideoForbidden
	}

	if objectPath == "" {
		objectPath = "index.m3u8"
	}

	cleanPath := path.Clean(objectPath)
	if strings.Contains(cleanPath, "..") {
		return "", ErrVideoForbidden
	}
	if path.IsAbs(cleanPath) {
		return "", ErrVideoForbidden
	}

	fullKey := path.Join(video.StreamKey(), cleanPath)

	url, err := uc.Storage.GenerateTemporaryAccessURL(ctx, fullKey, 5*time.Minute)
	if err != nil {
		return "", err
	}

	return url, nil
}
