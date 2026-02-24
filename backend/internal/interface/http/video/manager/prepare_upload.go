package manager

import (
	"errors"
	"net/http"

	"example.com/m/internal/usecase/video/manage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UploadSourceRequest struct {
	FileSize int64 `json:"file_size"`
}

type UploadSourceResponse struct {
	UploadID string   `json:"upload_id"`
	URLs     []string `json:"urls"`
	PartSize int64    `json:"part_size"`
	Key      string   `json:"key"`
}

func (h *VideoManagementHandler) PrepareUpload(c echo.Context) error {
	ctx := c.Request().Context()

	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid video ID: "+err.Error())
	}

	var req UploadSourceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body: "+err.Error())
	}
	if req.FileSize <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "file_size must be greater than 0")
	}

	result, err := h.manageUsecase.PrepareUploadSession(ctx, videoID, req.FileSize)
	if err != nil {
		if errors.Is(err, manage.ErrFileSizeTooLarge) {
			return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "file too large")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to prepare upload session: "+err.Error())
	}

	res := UploadSourceResponse{
		UploadID: result.UploadID,
		URLs:     result.URLs,
		PartSize: result.PartSize,
		Key:      result.Key,
	}

	return c.JSON(http.StatusCreated, res)
}
