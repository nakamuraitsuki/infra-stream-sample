package shared

import (
	"context"
	"time"

	"example.com/m/internal/domain/event"
	"github.com/google/uuid"
)

type OutboxEntry struct {
	ID         uuid.UUID
	EventType  string
	Payload    []byte
	OccurredAt time.Time
}

type OutboxRepository interface {
	Save(ctx context.Context, events []event.Event) error

	ListUnpublished(ctx context.Context, limit int) ([]OutboxEntry, error)

	MarkAsPublished(ctx context.Context, id uuid.UUID) error
}
