package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const DummyLoginCookieName = "dummy_login_user_id"

func DummyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(DummyLoginCookieName)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: no cookie")
		}

		userID, err := uuid.Parse(cookie.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: invalid user id")
		}

		SetUserID(c, userID)

		return next(c)
	}
}
