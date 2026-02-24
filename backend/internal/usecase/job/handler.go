package job

import "context"

type Handler interface {
	Handle(ctx context.Context, meta Metadata, payload []byte) error
}
