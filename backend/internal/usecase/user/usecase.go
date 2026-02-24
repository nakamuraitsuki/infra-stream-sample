package user

import (
	"context"

	"example.com/m/internal/domain/user"
	"github.com/google/uuid"
)

type UserUseCaseInterface interface {
	Register(ctx context.Context, name string) (*RegisterResult, error)
	GetMe(ctx context.Context, id uuid.UUID) (*GetMeResult, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, name string, bio string) error
	UpdateIcon(ctx context.Context, id uuid.UUID, data []byte) error
}

type UserUseCase struct {
	repo        user.Repository
	iconStorage user.IconStorage
}

func NewUserUseCase(repo user.Repository, iconStorage user.IconStorage) UserUseCaseInterface {
	return &UserUseCase{
		repo:        repo,
		iconStorage: iconStorage,
	}
}
