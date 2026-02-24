//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"example.com/m/internal/infrastructure/persistence/postgres"
	"example.com/m/internal/infrastructure/persistence/postgres/outbox"
	user_repo "example.com/m/internal/infrastructure/persistence/postgres/user"
	video_repo "example.com/m/internal/infrastructure/persistence/postgres/video"
	"example.com/m/internal/infrastructure/storage/s3"
	"example.com/m/internal/infrastructure/transcoder/ffmpeg"
	"example.com/m/internal/interface/http"
	user_h "example.com/m/internal/interface/http/user"
	"example.com/m/internal/interface/http/video/manager"
	"example.com/m/internal/interface/http/video/viewer"
	user_uc "example.com/m/internal/usecase/user"
	"example.com/m/internal/usecase/video/manage"
	"example.com/m/internal/usecase/video/view"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type HTTPServerApp struct {
	Echo *echo.Echo
	DB   *sqlx.DB
}

func provideContext() context.Context {
	return context.Background()
}

func InitializeHTTPServer() (*HTTPServerApp, error) {
	wire.Build(
		provideContext,
		// Configs
		postgres.NewPostgresConfig,
		s3.NewS3Config,

		// Infra Clients
		postgres.NewClient,
		postgres.NewTransactor,
		s3.NewClient,

		// Domain Services
		s3.NewIconStorage,
		s3.NewVideoStorage,
		ffmpeg.NewFFmpegTranscoder,

		// Repositories
		video_repo.NewRepository,
		user_repo.NewRepository,
		outbox.NewRepository,

		// UseCases
		user_uc.NewUserUseCase,
		manage.NewVideoManagementUseCase,
		view.NewVideoViewingUseCase,

		// HTTP Handlers
		user_h.NewHandler,
		manager.NewVideoManagementHandler,
		viewer.NewVideoViewingHandler,

		// Router
		http.NewRouter,

		wire.Struct(new(HTTPServerApp), "*"),
	)

	return nil, nil
}
