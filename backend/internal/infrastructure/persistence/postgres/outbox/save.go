package outbox

import (
	"context"
	"fmt"

	"example.com/m/internal/domain/event"
	"example.com/m/internal/infrastructure/persistence/postgres"
	"github.com/jmoiron/sqlx"
)

func (r *outboxRepository) Save(ctx context.Context, events []event.Event) error {
	if len(events) == 0 {
		return nil
	}
	db := postgres.GetExt(ctx, r.db)

	models := make([]*outboxModel, len(events))
	for i, e := range events {
		m, err := fromEntity(e)
		if err != nil {
			return fmt.Errorf("failed to map event to outbox DTO: %w", err)
		}
		models[i] = m
	}

	const query = `
		INSERT INTO outbox (id, event_type, payload, occurred_at)
		VALUES (:id, :event_type, :payload, :occurred_at)
	`

	if _, err := sqlx.NamedExecContext(ctx, db, query, models); err != nil {
		return fmt.Errorf("failed to execute bulk insert query: %w", err)
	}

	return nil
}
