package watcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/grafov/m3u8"
)

type Video struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ListPublicResponse struct {
	Items []Video `json:"items"`
}

type PlaybackInfoResponse struct {
	PlaybackURL string `json:"playback_url"`
}

func SimulateWatcher(ctx context.Context, id int, baseURL string) error {
	client := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &dockerRewriteTransport{rt: http.DefaultTransport},
	}

	// get public videos
	getPublicParams := map[string]string{"limit": "10"}

	body, err := fetchURL(ctx, client, baseURL+"/api/videos", getPublicParams)
	if err != nil {
		return nil // NOTE: サーバに関係無いエラーは無視する
	}

	var listResp ListPublicResponse
	if err := json.Unmarshal(body, &listResp); err != nil || len(listResp.Items) == 0 {
		fmt.Println("Error parsing videos:", err)
		return nil
	}

	// choice a random video to watch
	target := listResp.Items[rand.Intn(len(listResp.Items))]

	// NOTE: simulate user thinking time before watching
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

	// get playback info
	infoBody, err := fetchURL(
		ctx,
		client,
		fmt.Sprintf("%s/api/videos/%s/playback-info", baseURL, target.ID),
		nil,
	)
	if err != nil {
		return nil // NOTE: サーバに関係無いエラーは無視する
	}

	var info PlaybackInfoResponse
	if err := json.Unmarshal(infoBody, &info); err != nil {
		fmt.Println("Error parsing playback info:", err)
		return nil // NOTE: サーバに関係無いエラーは無視する
	}

	return watchHLS(ctx, client, baseURL, info.PlaybackURL, id)
}

func watchHLS(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	playbackURL string,
	userID int,
) error {
	fullPlaylistURL := baseURL + playbackURL

	req, _ := http.NewRequestWithContext(ctx, "GET", fullPlaylistURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("User %d: Error fetching playlist: %v\n", userID, err)
		return err
	}
	defer resp.Body.Close() 

	finalURL := resp.Request.URL

	body, _ := io.ReadAll(resp.Body)
	playlist, listType, err := m3u8.Decode(*bytes.NewBuffer(body), true)
	if err != nil {
		return err
	}

	if listType == m3u8.MEDIA {
		p := playlist.(*m3u8.MediaPlaylist)
		for i, seg := range p.Segments {
			if seg == nil {
				continue
			}

			segURL, err := url.Parse(seg.URI)
			if err != nil {
				fmt.Printf("User %d: Error parsing segment URL: %v\n", userID, err)
				continue
			}
			finalSegURL := finalURL.ResolveReference(segURL).String()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(seg.Duration) * time.Second):
				start := time.Now()
				_, err := fetchURL(ctx, client, finalSegURL, nil)
				latency := time.Since(start)

				if err != nil {
					fmt.Printf("User %d: Error fetching segment %d: %v\n", userID, i, err)
				} else {
					fmt.Printf("User %d: Fetched segment %d in %v\n", userID, i, latency)
				}
			}

			// ある程度見たら離脱する
			if i >= 5 {
				break
			}
		}
	}
	return nil
}

func fetchURL(
	ctx context.Context,
	client *http.Client,
	url string,
	params map[string]string,
) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
