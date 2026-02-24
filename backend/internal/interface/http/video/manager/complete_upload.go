package manager

import (
	"net/http"

	"example.com/m/internal/usecase/video/manage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PartInfo struct {
	PartNumber int32  `json:"part_number"`
	ETag       string `json:"etag"`
}

type CompleteUploadRequest struct {
	UploadID string     `json:"upload_id"`
	Parts    []PartInfo `json:"parts"`
}

func (h *VideoManagementHandler) CompleteUpload(c echo.Context) error {
	ctx := c.Request().Context()

	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid video ID: "+err.Error())
	}

	var req CompleteUploadRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body: "+err.Error())
	}

	parts := make([]manage.UploadPart, len(req.Parts))
	for i, part := range req.Parts {
		parts[i] = manage.UploadPart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
		}
	}

	err = h.manageUsecase.CompleteUploadSession(ctx, manage.CompleteUploadRequest{
		VideoID:  videoID,
		UploadID: req.UploadID,
		Parts:    parts,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to complete upload session: "+err.Error())
	}

	return c.NoContent(http.StatusOK)
}
