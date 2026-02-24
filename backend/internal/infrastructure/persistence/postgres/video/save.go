package video

import (
	"context"
	"fmt"

	video_domain "example.com/m/internal/domain/video"
	"example.com/m/internal/infrastructure/persistence/postgres"
	"github.com/jmoiron/sqlx"
)

func (r *videoRepository) Save(ctx context.Context, v *video_domain.Video) error {
	db := postgres.GetExt(ctx, r.db)

	model := fromEntity(v)

	const upsertVideoQuery = `
INSERT INTO videos (
	id, owner_id, source_key, stream_key, status, title,
	description, retry_count, failure_reason, visibility, created_at
) VALUES (
	:id, :owner_id, :source_key, :stream_key, :status, :title,
	:description, :retry_count, :failure_reason, :visibility, :created_at
)
ON CONFLICT (id) DO UPDATE SET
	source_key = EXCLUDED.source_key,
	stream_key = EXCLUDED.stream_key,
	status = EXCLUDED.status,
	title = EXCLUDED.title,
	description = EXCLUDED.description,
	retry_count = EXCLUDED.retry_count,
	failure_reason = EXCLUDED.failure_reason,
	visibility = EXCLUDED.visibility
`
	if _, err := sqlx.NamedExecContext(ctx, db, upsertVideoQuery, model); err != nil {
		return fmt.Errorf("failed to save video: %w", err)
	}

	const deleteTagsQuery = `DELETE FROM video_tags WHERE video_id = $1`
	if _, err := db.ExecContext(ctx, deleteTagsQuery, v.ID()); err != nil {
		return fmt.Errorf("failed to delete existing video tags: %w", err)
	}

	if len(v.Tags()) > 0 {
		// NOTE: 多対多を１クエリで解決する（postgres特有）
		//       Tagsテーブルに、もし新規がある場合は登録してから一覧を取得
		//       それをもとにvideoとの紐付けを更新する
		const syncTagsQuery = `
WITH inserted_tags AS (
	INSERT INTO tags (name)
	SELECT unnest($1::text[])
	ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
	RETURNING id
)
INSERT INTO video_tags (video_id, tag_id)
SELECT $2, id FROM inserted_tags
`
		tagNames := make([]string, len(v.Tags()))
		for i, tag := range v.Tags() {
			tagNames[i] = string(tag)
		}

		if _, err := db.ExecContext(ctx, syncTagsQuery, tagNames, v.ID()); err != nil {
			return fmt.Errorf("failed to sync video tags: %w", err)
		}
	}

	return nil
}
