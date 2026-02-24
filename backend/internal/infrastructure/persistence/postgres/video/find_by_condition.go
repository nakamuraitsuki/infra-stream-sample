package video

import (
	"context"
	"fmt"
	"strings"

	video_domain "example.com/m/internal/domain/video"
	"example.com/m/internal/infrastructure/persistence/postgres"
)

func (r *videoRepository) FindByCondition(
	ctx context.Context,
	cond video_domain.ListCondition,
) ([]*video_domain.Video, error) {
	db := postgres.GetExt(ctx, r.db)

	// Condition をもとに動的な WHERE 句を構築
	var whereClauses []string
	var args []interface{}
	argCount := 1

	if cond.OwnerID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("v.owner_id = $%d", argCount))
		args = append(args, *cond.OwnerID)
		argCount++
	}

	if cond.Visibility != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("v.visibility = $%d", argCount))
		args = append(args, *cond.Visibility)
		argCount++
	}

	if cond.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("v.status = $%d", argCount))
		args = append(args, *cond.Status)
		argCount++
	}

	// Tagに関しては、サブクエリで存在チェックを行う
	if cond.Tag != nil {
		whereClauses = append(whereClauses, fmt.Sprintf(`
EXISTS (
		SELECT 1 FROM video_tags vt2
		JOIN tags t2 ON vt2.tag_id = t2.id
		WHERE vt2.video_id = v.id AND t2.name = $%d
)
		`, argCount))
		args = append(args, *cond.Tag)
		argCount++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + fmt.Sprintf("%s", strings.Join(whereClauses, " AND "))
	}

	query := fmt.Sprintf(`
SELECT
	v.id, v.owner_id, v.source_key, v.stream_key, v.status,
	v.title, v.description, v.retry_count, v.failure_reason,
	v.visibility, v.created_at,
	array_remove(array_agg(t.name), NULL) AS tags
FROM
	videos v
	LEFT JOIN video_tags vt ON v.id = vt.video_id
	LEFT JOIN tags t ON vt.tag_id = t.id
%s
GROUP BY v.id
ORDER BY v.created_at DESC
LIMIT $%d OFFSET $%d
`, whereSQL, argCount, argCount+1)
	args = append(args, cond.Limit, cond.Offset)

	var models []videoModel
	if err := db.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	results := make([]*video_domain.Video, len(models))
	for i, m := range models {
		entity, err := m.toEntity()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to entity: %w", err)
		}
		results[i] = entity
	}
	return results, nil
}
