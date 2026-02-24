package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	user_domain "example.com/m/internal/domain/user"
	"example.com/m/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
)

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*user_domain.User, error) {
	db := postgres.GetExt(ctx, r.db)

	const query = `
SELECT
	id, name, bio, icon_key, role, created_at
FROM
	users
WHERE id = $1
`

	var model userModel
	err := db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find user by id: %w", err)
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return model.toEntity()
}
