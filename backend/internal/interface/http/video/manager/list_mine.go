package manager

import (
	"example.com/m/internal/interface/http/middleware"
	"example.com/m/internal/usecase/video/query"
	"github.com/labstack/echo/v4"
)

type ListMineRequest struct {
	Limit int `query:"limit"`
}

type ListMineResponse struct {
	Items []ListMineItem `json:"items"`
}

type ListMineItem struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
	Visibility  string   `json:"visibility"`
	CreatedAt   string   `json:"created_at"`
}

func (h *VideoManagementHandler) ListMine(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}

	var req ListMineRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "invalid query parameters: "+err.Error())
	}

	result, err := h.manageUsecase.ListMine(ctx, userID, query.VideoSearchQuery{
		Limit: req.Limit,
	})
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	items := make([]ListMineItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = ListMineItem{
			ID:          item.ID.String(),
			Title:       item.Title,
			Description: item.Description,
			Status:      item.Status,
			Tags:        item.Tags,
			Visibility:  item.Visibility,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	resp := ListMineResponse{
		Items: items,
	}

	return c.JSON(200, resp)
}
