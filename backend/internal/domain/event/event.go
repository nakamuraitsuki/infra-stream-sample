package event

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	ID() uuid.UUID
	EventType() string
	OccurredAt() time.Time
	Payload() any
}
