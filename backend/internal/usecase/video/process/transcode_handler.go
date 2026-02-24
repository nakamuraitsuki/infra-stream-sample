package process

import (
	"context"
	"encoding/json"

	"example.com/m/internal/usecase/job"
)

type TranscodeHandler struct {
	UseCase VideoProcessUseCaseInterface
}

// Handle は、Consumerが呼び出せるようにUseCaseをラップする。
func (h *TranscodeHandler) Handle(ctx context.Context, meta job.Metadata, payload []byte) error {
	var p TranscodePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}

	err := h.UseCase.Handle(
		ctx,
		p.VideoID,
		meta.Attempt+1 >= meta.MaxRetry,
	)
	if err != nil {
		if meta.Attempt >= meta.MaxRetry {
			return nil
		}
		return err
	}

	return nil
}
