package user

import (
	"net/http"

	"example.com/m/internal/interface/http/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DummyLoginResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name,omitempty"`
	Bio     string    `json:"bio,omitempty"`
	IconKey *string   `json:"icon_key,omitempty"`
	Role    string    `json:"role,omitempty"`
}

// POST /users/login
func (h *Handler) DummyLogin(c echo.Context) error {
	// 仮ユーザーとして UUID ゼロ値を Cookie に書き込む
	c.SetCookie(&http.Cookie{
		Name:     middleware.DummyLoginCookieName,
		Value:    uuid.Nil.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // 開発環境では false、本番では true
	})

	return c.JSON(http.StatusOK, DummyLoginResponse{
		ID:   uuid.Nil,
		Name: "Dummy User",
		Bio:  "This is a dummy user for development purposes.",
		Role: "user",
	})
}
