package video

import (
	"context"

	"example.com/m/internal/domain/video/value"
	"github.com/google/uuid"
)

// ListCondition defines the conditions for listing videos.
type ListCondition struct {
	OwnerID    *uuid.UUID
	Tag        *value.Tag
	Visibility *value.Visibility
	Status     *value.Status

	Limit  int
	Offset int
	// add more fields as needed
}

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Video, error)
	Save(ctx context.Context, v *Video) error

	FindByCondition(
		ctx context.Context,
		cond ListCondition,
	) ([]*Video, error)
}
