package video

import (
	"example.com/m/internal/interface/http/middleware"
	"example.com/m/internal/interface/http/video/manager"
	"example.com/m/internal/interface/http/video/viewer"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, mh *manager.VideoManagementHandler, vh *viewer.VideoViewingHandler) {
	g.GET("/videos/:id", vh.GetByID)
	g.GET("/videos", vh.ListPublic)
	g.GET("/videos/search", vh.SearchByTag)
	g.GET("/videos/:id/playback-info", vh.GetPlaybackInfo)
	g.GET("/videos/:id/stream/*", vh.GetVideoStream) // 比較用

	videos := g.Group("/videos")
	videos.Use(middleware.DummyAuthMiddleware) // Dummyの認証ミドルウェアを使用

	videos.POST("", mh.Create)
	videos.GET("/mine", mh.ListMine)
	videos.POST("/:id/upload/init", mh.PrepareUpload)
	videos.POST("/:id/upload/complete", mh.CompleteUpload)
}
