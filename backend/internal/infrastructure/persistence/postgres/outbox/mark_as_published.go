package outbox

import (
	"context"
	"fmt"
	"time"

	"example.com/m/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
)

func (r *outboxRepository) MarkAsPublished(ctx context.Context, id uuid.UUID) error {
	db := postgres.GetExt(ctx, r.db)

	// NOTE: 今回は物理削除を採用
	//       履歴はfmtログ等で管理する想定
	const query = `DELETE FROM outbox WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	fmt.Printf("Outbox event %s marked as published and deleted at %s\n", id.String(), time.Now().Format(time.RFC3339))
	return nil
}
