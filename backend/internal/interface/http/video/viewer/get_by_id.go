package viewer

import (
	"net/http"

	view_usecase "example.com/m/internal/usecase/video/view"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetByIDResponse struct {
	ID          string   `json:"id"`
	OwnerID     string   `json:"owner_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Visibility  string   `json:"visibility"`
	CreatedAt   string   `json:"created_at"`
}

func (h *VideoViewingHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid video ID: "+err.Error())
	}

	result, err := h.usecase.GetByID(ctx, videoID)
	if err != nil {
		switch err {
		case view_usecase.ErrVideoNotReady:
			return echo.NewHTTPError(http.StatusNotFound, "video is not ready for playback: "+err.Error())
		case view_usecase.ErrVideoForbidden:
			return echo.NewHTTPError(http.StatusForbidden, "video is not accessible: "+err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get video: "+err.Error())
		}
	}

	resp := GetByIDResponse{
		ID:          result.ID.String(),
		OwnerID:     result.OwnerID.String(),
		Title:       result.Title,
		Description: result.Description,
		Tags:        result.Tags,
		Visibility:  result.Visibility,
		CreatedAt:   result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, resp)
}
