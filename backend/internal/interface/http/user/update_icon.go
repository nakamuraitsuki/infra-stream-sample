package user

import (
	"io"

	"example.com/m/internal/interface/http/middleware"
	"github.com/labstack/echo/v4"
)

const MAX_ICON_FILE_SIZE = 5 * 1024 * 1024 // 5MB

func (h *Handler) UpdateIcon(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	// multipart/form-data からファイルを取得
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(400, "failed to get uploaded file: "+err.Error())
	}

	if file.Size > MAX_ICON_FILE_SIZE {
		return echo.NewHTTPError(400, "file size exceeds the limit")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(500, "failed to open uploaded file: "+err.Error())
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return echo.NewHTTPError(500, "failed to read uploaded file: "+err.Error())
	}

	if err := h.usecase.UpdateIcon(ctx, userID, data); err != nil {
		if httpErr, ok := err.(*echo.HTTPError); ok {
			return httpErr
		}
		return echo.NewHTTPError(500, "failed to update icon: "+err.Error())
	}

	return c.NoContent(204)
}
