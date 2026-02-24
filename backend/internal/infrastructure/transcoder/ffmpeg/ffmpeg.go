package ffmpeg

import (
	"example.com/m/internal/domain/video"
)

type ffmpegTranscoder struct {
	storage video.Storage
}

func NewFFmpegTranscoder(storage video.Storage) video.Transcoder {
	return &ffmpegTranscoder{
		storage: storage,
	}
}
