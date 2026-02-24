package ffmpeg

import (
	"context"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func (t *ffmpegTranscoder) watchAndQueue(
	ctx context.Context,
	dir string,
	pathCh chan<- string,
) error {
	// cf. https://pkg.go.dev/github.com/fsnotify/fsnotify
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(dir); err != nil {
		return err
	}

	for {
		select {
		// -- contextキャンセルの検知 --
		case <-ctx.Done():
			log.Println("Watcher received shutdown signal, exiting...")
			return nil

		// -- ファイルイベントの検知 --
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("Watcher event channel closed, exiting...")
				return nil // watcherが閉じられた場合
			}

			// NOTE: ffmpegでtemp_fileフラグを使う前提
			//       Rename を検知してアップロードする
			if event.Op&(fsnotify.Create|fsnotify.Rename|fsnotify.Write) != 0 {
				if strings.HasSuffix(event.Name, ".ts") {
					select {
					case pathCh <- event.Name:
					case <-ctx.Done():
						return nil
					}
				}
			}

		// -- watcherのエラー検知 --
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Println("Watcher error channel closed, exiting...")
				return nil // watcherが閉じられた場合
			}
			log.Printf("Watcher error: %v", err)
			return nil
		}
	}
}
