package redis

import (
	"context"
	"log"

	"golang.org/x/sync/errgroup"
)

func (c *consumer) workerPool(
	ctx context.Context,
	numWorkers int,
	jobCh <-chan []byte,
) error {
	eg, ctx := errgroup.WithContext(ctx)

	for range numWorkers {
		eg.Go(func() error {
			for {
				select {
				// -- contextキャンセルの検知 --
				case <-ctx.Done():
					log.Println("Worker received shutdown signal, exiting...")
					return nil

				// -- ジョブの処理 --
				case job, ok := <-jobCh:
					if !ok {
						log.Println("Job channel closed, worker exiting...")
						return nil // channelが閉じられた場合
					}
					if err := c.dispatch(ctx, job); err != nil {
						log.Printf("Failed to dispatch job: %v", err)
						return err
					}
				}
			}
		})
	}

	return eg.Wait()
}
