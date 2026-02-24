package user

import (
	"context"

	"github.com/google/uuid"
)

func (uc *UserUseCase) UpdateProfile(
	ctx context.Context,
	id uuid.UUID,
	name string,
	bio string,
) error {

	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := user.UpdateProfile(name, bio); err != nil {
		return err
	}

	return uc.repo.Save(ctx, user)
}
