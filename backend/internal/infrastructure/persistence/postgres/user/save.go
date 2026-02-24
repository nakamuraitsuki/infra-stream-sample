package user

import (
	"context"
	"fmt"

	user_domain "example.com/m/internal/domain/user"
	"example.com/m/internal/infrastructure/persistence/postgres"
	"github.com/jmoiron/sqlx"
)

func (r *userRepository) Save(ctx context.Context, user *user_domain.User) error {
	db := postgres.GetExt(ctx, r.db)

	model := fromEntity(user)

	const query = `
INSERT INTO users (
	id, name, bio, icon_key, role, created_at
) VALUES (
	:id, :name, :bio, :icon_key, :role, :created_at
)
ON CONFLICT (id) DO UPDATE SET
	name = EXCLUDED.name,
	bio = EXCLUDED.bio,
	icon_key = EXCLUDED.icon_key,
	role = EXCLUDED.role
`

	_, err := sqlx.NamedExecContext(ctx, db, query, model)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
