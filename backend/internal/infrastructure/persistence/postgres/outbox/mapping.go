package outbox

import (
	"encoding/json"
	"time"

	"example.com/m/internal/domain/event"
	"example.com/m/internal/domain/shared"
	"github.com/google/uuid"
)

type outboxModel struct {
	ID         uuid.UUID `db:"id"`
	EventType  string    `db:"event_type"`
	Payload    []byte    `db:"payload"`
	OccurredAt time.Time `db:"occurred_at"`
}

func fromEntity(e event.Event) (*outboxModel, error) {
	payload, err := json.Marshal(e.Payload())
	if err != nil {
		return nil, err
	}

	return &outboxModel{
		ID:         e.ID(),
		EventType:  e.EventType(),
		Payload:    payload,
		OccurredAt: e.OccurredAt(),
	}, nil
}

func (m *outboxModel) toEntry() (shared.OutboxEntry, error) {
	return shared.OutboxEntry{
		ID:         m.ID,
		EventType:  m.EventType,
		Payload:    m.Payload,
		OccurredAt: m.OccurredAt,
	}, nil
}
