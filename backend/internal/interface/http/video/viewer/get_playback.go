package viewer

import (
	"net/http"

	"example.com/m/internal/usecase/video/view"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetPlaybackInfoResponse struct {
	PlaybackURL string `json:"playback_url"`
	MIMEType    string `json:"mime_type"`
}

func (h *VideoViewingHandler) GetPlaybackInfo(c echo.Context) error {
	ctx := c.Request().Context()

	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid video ID: "+err.Error())
	}

	info, err := h.usecase.GetPlaybackInfo(ctx, videoID)
	if err != nil {
		switch err {
		case view.ErrVideoNotReady:
			return echo.NewHTTPError(http.StatusConflict, "video not ready: "+err.Error())
		case view.ErrVideoForbidden:
			return echo.NewHTTPError(http.StatusForbidden, "video is not public: "+err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get playback info: "+err.Error())
		}
	}

	resp := GetPlaybackInfoResponse{
		PlaybackURL: info.PlaybackURL,
		MIMEType:    info.MIMEType,
	}

	return c.JSON(http.StatusOK, resp)
}
