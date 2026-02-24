package ffmpeg

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func (t *ffmpegTranscoder) workerPool(
	ctx context.Context,
	numWorkers int,
	pathCh <-chan string,
	streamKey string,
) error {

	eg, ctx := errgroup.WithContext(ctx)

	for range numWorkers {
		eg.Go(func() error {
			for {
				select {
				// -- contextキャンセルの検知 --
				case <-ctx.Done():
					log.Println("Worker received shutdown signal, exiting...")
					return nil // contextがキャンセルされたら終了

				// -- アップロード --
				case path, ok := <-pathCh:
					if !ok {
						log.Println("Path channel closed, worker exiting...")
						return nil // channelが閉じられた場合
					}
					log.Printf("Worker received path: %s", path)
					if err := t.uploadFile(ctx, path, streamKey); err != nil {
						log.Printf("Failed to upload file: %v", err)
						return err
					}
				}
			}
		})
	}

	return eg.Wait()
}

func (t *ffmpegTranscoder) uploadFile(
	ctx context.Context,
	filePath, streamKeyPrefix string,
) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	s3Key := fmt.Sprintf("%s/%s", strings.TrimRight(streamKeyPrefix, "/"), fileName)

	return t.storage.SaveStream(ctx, s3Key, file)
}
