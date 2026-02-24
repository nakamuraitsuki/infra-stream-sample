package video

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	video_domain "example.com/m/internal/domain/video"
	"example.com/m/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
)

func (r *videoRepository) FindByID(ctx context.Context, id uuid.UUID) (*video_domain.Video, error) {
	db := postgres.GetExt(ctx, r.db)

	// NOTE: postgresの集約関数でタグを配列として取得
	const query = `
SELECT 
	v.id, v.owner_id, v.source_key, v.stream_key, v.status,
	v.title, v.description, v.retry_count, v.failure_reason,
	v.visibility, v.created_at,
	array_remove(array_agg(t.name), NULL) as tags
FROM videos v
LEFT JOIN video_tags vt ON v.id = vt.video_id
LEFT JOIN tags t ON vt.tag_id = t.id
WHERE v.id = $1
GROUP BY v.id
`
	var m videoModel
	if err := db.GetContext(ctx, &m, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("video not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find video by id: %w", err)
	}

	entity, err := m.toEntity()
	if err != nil {
		return nil, fmt.Errorf("failed to convert video model to entity: %w", err)
	}

	return entity, nil
}
