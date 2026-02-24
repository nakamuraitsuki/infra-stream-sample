package user

import "github.com/labstack/echo/v4"

type RegisterRequest struct {
	Name string `json:"name"`
}

type RegisterResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Bio     string  `json:"bio"`
	IconKey *string `json:"icon_key,omitempty"`
	Role    string  `json:"role"`
}

func (h *Handler) Register(c echo.Context) error {

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	ctx := c.Request().Context()
	result, err := h.usecase.Register(ctx, req.Name)
	if err != nil {
		return err
	}

	response := RegisterResponse{
		ID:      result.ID.String(),
		Name:    result.Name,
		Bio:     result.Bio,
		IconKey: result.IconKey,
		Role:    result.Role,
	}

	return c.JSON(201, response)
}
