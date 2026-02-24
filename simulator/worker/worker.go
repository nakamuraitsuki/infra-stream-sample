package worker

import (
	"context"
	"log"
	"math/rand"
	"time"

	"example.com/m/watcher"
	"golang.org/x/sync/errgroup"
)

func StartWorker(ctx context.Context, maxWorkses int, baseURL string) {
	eg, ctx := errgroup.WithContext(ctx)

	for i := 1; i <= maxWorkses; i++ {
		id := i
		eg.Go(func() error {
			initialWait := time.NewTimer(time.Duration(id) * time.Second)
			defer initialWait.Stop()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-initialWait.C:
			}

			// 無限ループで動画視聴を繰り返す
			for {
				log.Printf("Worker [%d]: Starting a new video session...", id)

				watcher.SimulateWatcher(ctx, id, baseURL)

				interval := time.Duration(1+rand.Intn(5)) * time.Second

				select {
				case <-ctx.Done():
					log.Printf("Worker [%d]: Stopping due to context cancellation.", id)
					return ctx.Err()
				case <-time.After(interval):
					// ループの先頭に戻って次の動画へ
				}
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatalf("\n Simulation stopped with error: %v", err)
	}
}
