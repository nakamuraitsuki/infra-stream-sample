package outbox

import (
	"context"

	"example.com/m/internal/domain/shared"
	"example.com/m/internal/infrastructure/persistence/postgres"
)

func (r *outboxRepository) ListUnpublished(ctx context.Context, limit int) ([]shared.OutboxEntry, error) {
	db := postgres.GetExt(ctx, r.db)

	const query = `
SELECT id, event_type, payload, occurred_at
FROM outbox
ORDER BY occurred_at ASC
LIMIT $1	
FOR UPDATE SKIP LOCKED
`
	var models []outboxModel
	if err := db.SelectContext(ctx, &models, query, limit); err != nil {
		return nil, err
	}

	entries := make([]shared.OutboxEntry, len(models))
	for i, model := range models {
		entry, err := model.toEntry()
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}

	return entries, nil
}
