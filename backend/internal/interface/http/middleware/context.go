package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type contextKey string

const userIDKey contextKey = "user_id"

func SetUserID(ctx echo.Context, id uuid.UUID) {
	ctx.Set(string(userIDKey), id)
}

func GetUserID(ctx echo.Context) (uuid.UUID, error) {
	v := ctx.Get(string(userIDKey))
	if v == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	id, ok := v.(uuid.UUID)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError, "invalid user id in context")
	}
	return id, nil
}
