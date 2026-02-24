package video

import (
	video_domain "example.com/m/internal/domain/video"
	"github.com/jmoiron/sqlx"
)

type videoRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) video_domain.Repository {
	return &videoRepository{db: db}
}
