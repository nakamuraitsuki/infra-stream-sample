package view

import (
	"context"
	"errors"
	"io"
	"log"
	"path"
	"strings"
	"time"

	video_domain "example.com/m/internal/domain/video"
	video_value "example.com/m/internal/domain/video/value"
	"example.com/m/internal/usecase/video/query"
	"github.com/google/uuid"
)

type GetVideoStreamMeta struct {
	TotalSize     int64
	ContentLength int64
	RangeStart    int64
	RangeEnd      int64
	ETag          string
	LastModified  time.Time
}

func (uc *VideoViewingUseCase) GetVideoStream(
	ctx context.Context,
	videoID uuid.UUID,
	objectPath string,
	byteRangeQuery *query.VideoRangeQuery,
) (io.ReadCloser, GetVideoStreamMeta, string, error) {

	video, err := uc.VideoRepo.FindByID(ctx, videoID)
	if err != nil {
		return nil, GetVideoStreamMeta{}, "", err
	}

	if video.Status() != video_value.StatusReady {
		return nil, GetVideoStreamMeta{}, "", ErrVideoNotReady
	}

	// Check visibility
	if video.Visibility() != video_value.VisibilityPublic {
		return nil, GetVideoStreamMeta{}, "", ErrVideoForbidden
	}

	if objectPath == "" {
		objectPath = "index.m3u8"
	}

	cleanPath := path.Clean(objectPath)
	if strings.Contains(cleanPath, "..") {
		return nil, GetVideoStreamMeta{}, "", ErrVideoForbidden
	}
	if path.IsAbs(cleanPath) {
		return nil, GetVideoStreamMeta{}, "", ErrVideoForbidden
	}

	fullKey := path.Join(video.StreamKey(), cleanPath)

	log.Println("Getting video stream with key:", fullKey)
	var byteRange *video_domain.ByteRange
	if byteRangeQuery != nil {
		byteRange = &video_domain.ByteRange{
			Start: byteRangeQuery.Start,
			End:   byteRangeQuery.End,
		}
	}

	rc, meta, err := uc.Storage.GetStream(ctx, fullKey, byteRange)
	if err != nil {
		return nil, GetVideoStreamMeta{}, "", errors.New("failed to get video stream: " + fullKey + " - " + err.Error())
	}

	metaResult := GetVideoStreamMeta{
		TotalSize:     meta.TotalSize,
		ContentLength: meta.ContentLength,
		RangeStart:    meta.RangeStart,
		RangeEnd:      meta.RangeEnd,
		ETag:          meta.ETag,
		LastModified:  meta.LastModified,
	}

	mimeType := uc.detectMimeTypeFromPath(cleanPath)

	return rc, metaResult, mimeType, nil
}

func (uc *VideoViewingUseCase) detectMimeTypeFromPath(path string) string {
	switch {
	case strings.HasSuffix(path, ".m3u8"):
		return "application/vnd.apple.mpegurl"
	case strings.HasSuffix(path, ".ts"):
		return "video/mp2t"
	case strings.HasSuffix(path, ".mp4"):
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}
