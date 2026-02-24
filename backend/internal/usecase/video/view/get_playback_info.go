package view

import (
	"context"

	video_value "example.com/m/internal/domain/video/value"
	"github.com/google/uuid"
)

type PlaybackInfo struct {
	PlaybackURL string // URL 組み立てはHandler側で行う（リダイレクト）
	MIMEType    string
}

func (uc *VideoViewingUseCase) GetPlaybackInfo(
	ctx context.Context,
	videoID uuid.UUID,
) (*PlaybackInfo, error) {

	video, err := uc.VideoRepo.FindByID(ctx, videoID)
	if err != nil {
		return nil, err
	}

	if video.Status() != video_value.StatusReady {
		return nil, ErrVideoNotReady
	}

	if video.Visibility() != video_value.VisibilityPublic {
		return nil, ErrVideoForbidden
	}

	url := "/api/videos/" + videoID.String() + "/stream/"

	// url, err := uc.Storage.GenerateTemporaryAccessURL(
	// 	ctx,
	// 	video.StreamKey(),
	// 	15*time.Minute,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return &PlaybackInfo{
		PlaybackURL: url,
		MIMEType:    "application/x-mpegURL",
	}, nil
}
