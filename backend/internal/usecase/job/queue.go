package job

import "context"

type Queue interface {
	Enqueue(ctx context.Context, meta Metadata, payload []byte) error
}

type Consumer interface {
	Register(jobType string, handler Handler)
	Start(ctx context.Context) error
}
