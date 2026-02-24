package user

import (
	"example.com/m/internal/interface/http/middleware"
	"github.com/labstack/echo/v4"
)

type GetMeResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Bio     string  `json:"bio"`
	IconKey *string `json:"icon_key,omitempty"`
	Role    string  `json:"role"`
}

func (h *Handler) GetMe(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}

	ctx := c.Request().Context()
	userInfo, err := h.usecase.GetMe(ctx, userID)
	if err != nil {
		c.Logger().Errorf("failed to get user info for userID=%s: %v", userID, err)
		return echo.NewHTTPError(500, "failed to get user info")
	}

	response := GetMeResponse{
		ID:      userInfo.ID.String(),
		Name:    userInfo.Name,
		Bio:     userInfo.Bio,
		IconKey: userInfo.IconKey,
		Role:    userInfo.Role,
	}

	return c.JSON(200, response)
}
