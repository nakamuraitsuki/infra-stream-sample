package user

import (
	"context"
	"time"

	"example.com/m/internal/domain/user"
	"example.com/m/internal/domain/user/value"
	"github.com/google/uuid"
)

type RegisterResult struct {
	ID      uuid.UUID
	Name    string
	Bio     string
	IconKey *string
	Role    string
}

func (uc *UserUseCase) Register(
	ctx context.Context,
	name string,
) (*RegisterResult, error) {

	user := user.NewUser(
		uuid.New(),
		name,
		"",
		nil,
		value.RoleUser,
		time.Now(),
	)

	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &RegisterResult{
		ID:      user.ID(),
		Name:    user.Name(),
		Bio:     user.Bio(),
		IconKey: user.IconKey(),
		Role:    string(user.Role()),
	}, nil
}
