package http

import (
	"example.com/m/internal/interface/http/user"
	"example.com/m/internal/interface/http/video"
	"example.com/m/internal/interface/http/video/manager"
	"example.com/m/internal/interface/http/video/viewer"
	"github.com/labstack/echo/v4"
)

func NewRouter(
	uh *user.Handler,
	vmh *manager.VideoManagementHandler,
	vvh *viewer.VideoViewingHandler,
) *echo.Echo {
	e := echo.New()
	g := e.Group("/api")
	// 共通ミドルウェア設定はここ

	// 各ドメインのルート登録
	user.RegisterRoutes(g, uh)
	video.RegisterRoutes(g, vmh, vvh)

	return e
}
