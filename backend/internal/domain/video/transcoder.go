package video

import "context"

type Transcoder interface {
	Transcode(
		ctx context.Context,
		sourceKey string,
		streamKey string,
	) error
}
