package manager

import (
	"net/http"

	"example.com/m/internal/interface/http/middleware"
	"github.com/labstack/echo/v4"
)

type CreateRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type CreateResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
	Visibility  string   `json:"visibility"`
	CreatedAt   string   `json:"created_at"`
}

func (h *VideoManagementHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	res, err := h.manageUsecase.Create(
		ctx,
		userID,
		req.Title,
		req.Description,
		req.Tags,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := CreateResponse{
		ID:          res.ID.String(),
		Title:       res.Title,
		Description: res.Description,
		Status:      res.Status,
		Tags:        res.Tags,
		Visibility:  res.Visibility,
		CreatedAt:   res.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusCreated, resp)
}
