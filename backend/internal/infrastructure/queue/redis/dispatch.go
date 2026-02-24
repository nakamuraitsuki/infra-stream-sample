package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

func (c *consumer) dispatch(ctx context.Context, jobBytes []byte) error {
	var jobMsg jobMessage
	if err := json.Unmarshal(jobBytes, &jobMsg); err != nil {
		log.Printf("Failed to unmarshal job message: %v", err)
		return err
	}

	defer func() {
		// ジョブの処理が終わったらprocessing Queueから削除する
		if err := c.client.Universal.LRem(ctx, c.processingKey, 1, jobBytes).Err(); err != nil {
			log.Printf("Failed to remove job from processing queue [type: %s, id: %s]: %v", jobMsg.Meta.Type, jobMsg.Meta.ID, err)
		}
	}()

	c.mu.RLock()
	handler, ok := c.handlers[jobMsg.Meta.Type]
	c.mu.RUnlock()

	if !ok {
		log.Printf("No handler registered for job type: %s", jobMsg.Meta.Type)
		// ハンドラ未登録の場合でもジョブを消失させないため、元のキューに再投入する
		if err := c.queue.Enqueue(ctx, jobMsg.Meta, jobMsg.Payload); err != nil {
			log.Printf("Failed to re-enqueue job without registered handler [type: %s, id: %s]: %v", jobMsg.Meta.Type, jobMsg.Meta.ID, err)
		}
		// NOTE: リトライ処理を図るので、Consumer群に伝播させない
		return nil
	}

	if err := handler.Handle(ctx, jobMsg.Meta, jobMsg.Payload); err != nil {
		log.Printf("Error handling job [type: %s, id: %s]: %v", jobMsg.Meta.Type, jobMsg.Meta.ID, err)
		// NOTE: retry
		if jobMsg.Meta.Attempt < jobMsg.Meta.MaxRetry {
			jobMsg.Meta.Attempt++

			// NOTE: 指数バックオフの計算: 2^attempt * 1秒 (例: 2s, 4s, 8s...)
			backoff := time.Duration(1<<jobMsg.Meta.Attempt) * time.Second
			log.Printf("Retrying job %s in %v (Attempt: %d)", jobMsg.Meta.ID, backoff, jobMsg.Meta.Attempt)

			// NOTE: 本来は「遅延キュー」に入れるのがベストだが......
			time.Sleep(backoff)

			if err := c.queue.Enqueue(ctx, jobMsg.Meta, jobMsg.Payload); err != nil {
				log.Printf("Failed to re-enqueue job [type: %s, id: %s]: %v", jobMsg.Meta.Type, jobMsg.Meta.ID, err)
			}
		} else {
			// 最大リトライ回数を超えた場合、DLQ にジョブを移動する
			if err := c.client.Universal.LPush(ctx, c.dlqKey, jobBytes).Err(); err != nil {
				log.Printf("Failed to move job to DLQ [type: %s, id: %s]: %v", jobMsg.Meta.Type, jobMsg.Meta.ID, err)
			}
		}
	}
	return nil
}
