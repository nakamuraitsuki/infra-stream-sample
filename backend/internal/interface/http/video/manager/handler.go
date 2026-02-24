package manager

import (
	"example.com/m/internal/usecase/video/manage"
)

type VideoManagementHandler struct {
	manageUsecase manage.VideoManagementUseCaseInterface
}

func NewVideoManagementHandler(manageUsecase manage.VideoManagementUseCaseInterface) *VideoManagementHandler {
	return &VideoManagementHandler{
		manageUsecase: manageUsecase,
	}
}
