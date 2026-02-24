package job

import (
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Attempt   int       `json:"attempt"`
	MaxRetry  int       `json:"max_retry"`
	CreatedAt time.Time `json:"created_at"`
}
