package view

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"time"

	video_value "example.com/m/internal/domain/video/value"
	"github.com/google/uuid"
)

const (
	secretKey = "secret" // 本番環境では安全な方法で管理すること
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

	expires := time.Now().Add(2 * time.Hour).Unix()
	vid := videoID.String()

	hash := uc.generateHash(vid, expires, secretKey)

	// signedPath の組み立て順序を HASH / EXPIRES / VIDEO_ID に変更
	signedPath := fmt.Sprintf("/api/videos/%s/%d/%s/stream/index.m3u8", hash, expires, vid)

	return &PlaybackInfo{
		PlaybackURL: signedPath,
		MIMEType:    "application/x-mpegURL",
	}, nil
}

func (uc *VideoViewingUseCase) generateHash(videoID string, expires int64, secret string) string {
	input := fmt.Sprintf("%d%s%s", expires, videoID, secret)
	h := md5.Sum([]byte(input))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
