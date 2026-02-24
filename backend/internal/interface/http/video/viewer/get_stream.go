package viewer

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"example.com/m/internal/usecase/video/query"
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

	rangeHeader := c.Request().Header.Get("Range")
	byteRangeQuery, err := h.parseRangeHeader(rangeHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid range header: "+err.Error())
	}

	stream, meta, mimeType, err := h.usecase.GetVideoStream(ctx, videoID, objectPath, byteRangeQuery)
	if err != nil {
		switch err {
		case view.ErrVideoNotReady:
			return echo.NewHTTPError(http.StatusConflict, "video not ready: "+err.Error())
		case view.ErrVideoForbidden:
			return echo.NewHTTPError(http.StatusForbidden, "video is not public: "+err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get video stream: "+err.Error())
		}
	}
	defer stream.Close()

	res := c.Response()

	// cf. https://datatracker.ietf.org/doc/html/rfc7233#section-4

	res.Header().Set(echo.HeaderContentType, mimeType)
	res.Header().Set("Accept-Ranges", "bytes")

	// Use a stricter cache policy for playlists and non-segment files to avoid
	// unintentionally caching private or rapidly changing content. Allow short
	// public caching only for segment files to improve streaming performance.
	cacheControl := "no-store"
	res.Header().Set(echo.HeaderCacheControl, cacheControl)
	if meta.ETag != "" {
		res.Header().Set("ETag", meta.ETag)
	}

	if !meta.LastModified.IsZero() {
		res.Header().Set(
			echo.HeaderLastModified,
			meta.LastModified.UTC().Format(http.TimeFormat),
		)
	}

	if h.etagMatches(
		c.Request().Header.Get("If-None-Match"),
		meta.ETag,
	) {
		return c.NoContent(http.StatusNotModified)
	}

	if byteRangeQuery != nil {
		res.Header().Set(
			echo.HeaderContentLength,
			strconv.FormatInt(meta.ContentLength, 10),
		)

		contentRange := "bytes " +
			strconv.FormatInt(meta.RangeStart, 10) + "-" +
			strconv.FormatInt(meta.RangeEnd, 10) + "/" +
			strconv.FormatInt(meta.TotalSize, 10)
		res.Header().Set("Content-Range", contentRange)

		res.WriteHeader(http.StatusPartialContent)
	} else {
		res.Header().Set(
			echo.HeaderContentLength,
			strconv.FormatInt(meta.TotalSize, 10),
		)

		res.WriteHeader(http.StatusOK)
	}

	_, err = io.Copy(res.Writer, stream)
	return err
}

func (v *VideoViewingHandler) parseRangeHeader(r string) (*query.VideoRangeQuery, error) {
	if r == "" {
		return nil, nil
	}

	if !strings.HasPrefix(r, "bytes=") {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid range header format")
	}

	rangeSpec := strings.TrimPrefix(r, "bytes=")
	// Explicitly reject multipart byte ranges such as "bytes=0-499,1000-1499".
	if strings.Contains(rangeSpec, ",") {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "multipart byte ranges are not supported")
	}

	parts := strings.Split(rangeSpec, "-")
	if len(parts) != 2 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid range header format")
	}

	// NOTE: suffix-byte-range (RFC 7233 §2.1) is not supported.
	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid range start value: "+err.Error())
	}

	if parts[1] == "" {
		return &query.VideoRangeQuery{
			Start: start,
			End:   nil,
		}, nil
	}

	end, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid range end value: "+err.Error())
	}

	return &query.VideoRangeQuery{
		Start: start,
		End:   &end,
	}, nil
}

// 弱ETagに対応
func (v *VideoViewingHandler) etagMatches(ifNoneMatch, currentETag string) bool {
	if ifNoneMatch == "" {
		return false
	}

	parts := strings.Split(ifNoneMatch, ",")
	for _, part := range parts {
		tag := strings.TrimSpace(part)
		if tag == "*" {
			return true
		}

		tag = strings.TrimPrefix(tag, "W/")
		if tag == currentETag {
			return true
		}
	}

	return false
}
