package viewer

import "example.com/m/internal/usecase/video/view"

type VideoViewingHandler struct {
	usecase view.VideoViewingUseCaseInterface
}

func NewVideoViewingHandler(usecase view.VideoViewingUseCaseInterface) *VideoViewingHandler {
	return &VideoViewingHandler{
		usecase: usecase,
	}
}
