package user

import (
	"context"

	"github.com/google/uuid"
)

func (uc *UserUseCase) UpdateIcon(
	ctx context.Context,
	id uuid.UUID,
	data []byte,
) error {

	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	key := "icons/" + uuid.New().String() + ".png"

	if err := uc.iconStorage.Upload(ctx, key, data); err != nil {
		return err
	}

	user.UpdateIcon(&key)

	if err := uc.repo.Save(ctx, user); err != nil {
		_ = uc.iconStorage.Delete(ctx, key)
		return err
	}

	return nil
}
