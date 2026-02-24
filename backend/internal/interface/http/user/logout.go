package user

import (
	"net/http"

	"example.com/m/internal/interface/http/middleware"
	"github.com/labstack/echo/v4"
)

type LogoutResponse struct {
	Message string `json:"message"`
}

// POST /users/logout
func (h *Handler) Logout(c echo.Context) error {
	// Cookie を削除してログアウト
	c.SetCookie(&http.Cookie{
		Name:     middleware.DummyLoginCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // 即時削除
		HttpOnly: true,
		Secure:   false, // 開発環境向け
	})

	return c.JSON(http.StatusOK, LogoutResponse{
		Message: "logout successful",
	})
}
