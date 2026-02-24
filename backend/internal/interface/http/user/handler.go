package user

import (
	user_usecase "example.com/m/internal/usecase/user"
)

type Handler struct {
	usecase user_usecase.UserUseCaseInterface
}

func NewHandler(usecase user_usecase.UserUseCaseInterface) *Handler {
	return &Handler{
		usecase: usecase,
	}
}
