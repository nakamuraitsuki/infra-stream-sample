package viewer

import (
	"net/http"

	"example.com/m/internal/usecase/video/view"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *VideoViewingHandler) GetVideoStream(c echo.Context) error {
	ctx := c.Request().Context()

	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid video ID: "+err.Error())
	}

	objectPath := c.Param("*")

	signedURL, err := h.usecase.GetVideoStream(ctx, videoID, objectPath)
	if err != nil {
		switch err {
		case view.ErrVideoNotReady:
			return echo.NewHTTPError(http.StatusConflict, "video is not ready for playback")
		case view.ErrVideoForbidden:
			return echo.NewHTTPError(http.StatusForbidden, "video is not accessible")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get video stream: "+err.Error())
		}
	}

	// cf. https://moneyforward-dev.jp/entry/2021/01/13/s3-x-accel-redirect/
	c.Response().Header().Set("X-Accel-Redirect", "/internal-s3-proxy")
	c.Response().Header().Set("X-S3-Signed-URL", signedURL)

	return c.NoContent(http.StatusOK)
}
