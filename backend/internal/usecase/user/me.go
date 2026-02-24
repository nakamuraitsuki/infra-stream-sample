package user

import (
	"context"

	"github.com/google/uuid"
)

type GetMeResult struct {
	ID      uuid.UUID
	Name    string
	Bio     string
	IconKey *string
	Role    string
}

func (u *UserUseCase) GetMe(ctx context.Context, id uuid.UUID) (*GetMeResult, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetMeResult{
		ID:      user.ID(),
		Name:    user.Name(),
		Bio:     user.Bio(),
		IconKey: user.IconKey(),
		Role:    string(user.Role()),
	}, nil
}
