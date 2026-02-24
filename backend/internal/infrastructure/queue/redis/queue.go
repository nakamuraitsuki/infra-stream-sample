package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"example.com/m/internal/usecase/job"
)

type jobMessage struct {
	Meta    job.Metadata `json:"meta"`
	Payload []byte       `json:"payload"`
}

type queue struct {
	client *Client
	key    string // Queueの識別子
}

func NewQueue(client *Client, key string) job.Queue {
	return &queue{
		client: client,
		key:    key,
	}
}

func (q *queue) Enqueue(ctx context.Context, meta job.Metadata, payload []byte) error {
	msg := jobMessage{
		Meta:    meta,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal job message: %w", err)
	}

	err = q.client.Universal.LPush(ctx, q.key, data).Err()
	if err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	return nil
}
