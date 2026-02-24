package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/m/internal/domain/event"
	"example.com/m/internal/infrastructure/persistence/postgres"
	outbox_repo "example.com/m/internal/infrastructure/persistence/postgres/outbox"
	"example.com/m/internal/infrastructure/persistence/postgres/video"
	"example.com/m/internal/infrastructure/queue/redis"
	"example.com/m/internal/infrastructure/relay/outbox"
	"example.com/m/internal/infrastructure/storage/s3"
	"example.com/m/internal/infrastructure/transcoder/ffmpeg"
	event_uc "example.com/m/internal/usecase/event"
	"example.com/m/internal/usecase/video/process"
	"golang.org/x/sync/errgroup"
)

func main() {
	// NOTE: Worker Queue は、Keyを切り替えることで複数用途に利用可能
	//       よって、同じ型を複数作りづらい Wire は導入しない

	// -- Configs --
	jobQueueCfg := redis.QueueConfig{Key: "job_queue"}
	redisCfg := redis.NewRedisConfig()
	postgresCfg := postgres.NewPostgresConfig()
	s3Cfg := s3.NewS3Config()

	// -- Clients --
	redisClient, err := redis.NewClient(redisCfg)
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()
	postgresClient, err := postgres.NewClient(postgresCfg)
	if err != nil {
		panic(err)
	}
	defer postgresClient.Close()
	s3ClientSet, err := s3.NewClient(context.Background(), s3Cfg)
	if err != nil {
		panic(err)
	}

	// -- Domain Services --
	videoStorage := s3.NewVideoStorage(s3ClientSet, s3Cfg)
	transcoder := ffmpeg.NewFFmpegTranscoder(videoStorage)

	// -- Repositories --
	outboxRepo := outbox_repo.NewRepository(postgresClient)
	videoRepo := video.NewRepository(postgresClient)

	// -- Use Cases --
	transactor := postgres.NewTransactor(postgresClient)
	jobQueue := redis.NewQueue(redisClient, jobQueueCfg.Key)
	eventRelayUseCase := event_uc.NewEventRelayUseCase(outboxRepo, jobQueue, transactor)
	processUseCase := process.NewVideoProcessUseCase(videoRepo, transcoder, transactor)

	// -- Workers --
	interval := 5 * time.Second
	jobConsumer := redis.NewConsumer(redisClient, jobQueue, jobQueueCfg.Key)
	relayWorker := outbox.NewRelayWorker(eventRelayUseCase, interval)

	// ConsumerにWorkerを登録
	jobConsumer.Register(event.TranscodingRequestType, &process.TranscodeHandler{UseCase: processUseCase})

	// シグナルをキャッチしてGraceful Shutdownを行うためのコンテキストを作成
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)

	go func() {
		<-ctx.Done() // シグナルを待つ
		log.Printf("ROOT CTX DONE: %v", ctx.Err())
		log.Printf("ROOT CASE: %v", context.Cause(ctx))
	}()

	// Relay Worker の起動
	eg.Go(func() error {
		log.Println("Relay Worker start")
		err := relayWorker.Start(ctx)
		log.Printf("Relay Worker exit: %v", err)
		return err
	})

	// Consumer の起動
	eg.Go(func() error {
		log.Println("Job Consumer start")
		err := jobConsumer.Start(ctx)
		log.Printf("Job Consumer exit: %v", err)
		return err
	})
	// いずれかのWorkerがエラーを出すか、中断信号が来たら終了
	if err := eg.Wait(); err != nil {
		log.Printf("Worker stopped with error: %v", err)
		return
	}
	log.Println("Successfully shutdown all workers.")
}
