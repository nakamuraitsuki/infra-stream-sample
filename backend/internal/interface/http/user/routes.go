package user

import (
	"example.com/m/internal/interface/http/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.POST("/users", h.Register)
	g.POST("/users/login", h.DummyLogin) // 認証不要な仮ログインエンドポイント
	g.POST("/users/logout", h.Logout)

	auth := g.Group("/users")
	auth.Use(middleware.DummyAuthMiddleware) // Dummy認証ミドルウェアを適用

	auth.GET("/me", h.GetMe)
	auth.PATCH("/me/profile", h.UpdateProfile)
	auth.PUT("/me/icon", h.UpdateIcon)
}
