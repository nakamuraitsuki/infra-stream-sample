package outbox

import (
	"example.com/m/internal/domain/shared"
	"github.com/jmoiron/sqlx"
)

type outboxRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) shared.OutboxRepository {
	return &outboxRepository{db: db}
}
