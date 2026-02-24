package viewer

import (
	"net/http"

	"example.com/m/internal/usecase/video/query"
	"github.com/labstack/echo/v4"
)

type ListPublicRequest struct {
	Limit int `query:"limit"`
}

type ListPublicResponse struct {
	Items []ListPublicResponseItem `json:"items"`
}

type ListPublicResponseItem struct {
	ID          string   `json:"id"`
	OwnerID     string   `json:"owner_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at"`
}

func (h *VideoViewingHandler) ListPublic(c echo.Context) error {
	ctx := c.Request().Context()

	var req ListPublicRequest
	if err := echo.QueryParamsBinder(c).
		Int("limit", &req.Limit).
		BindError(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters: "+err.Error())
	}

	result, err := h.usecase.ListPublic(ctx, query.VideoSearchQuery{
		Limit: req.Limit,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list public videos: "+err.Error())
	}

	items := make([]ListPublicResponseItem, len(result.Videos))
	for i, item := range result.Videos {
		items[i] = ListPublicResponseItem{
			ID:          item.ID.String(),
			OwnerID:     item.OwnerID.String(),
			Title:       item.Title,
			Description: item.Description,
			Tags:        item.Tags,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	resp := ListPublicResponse{
		Items: items,
	}

	return c.JSON(200, resp)
}
