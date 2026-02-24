package s3

import "example.com/m/internal/infrastructure/env"

type Config struct {
	Endpoint       string
	PublicEndpoint string // クライアントがアクセスする際のエンドポイント（例: MinIOの外部公開URL）
	Region         string
	BucketName     string
	AccessKey      string
	SecretKey      string
	UsePathStyle   bool
	UseSSL         bool
}

func NewS3Config() Config {
	return Config{
		// NOTE: Default は MinIO のローカルを意識
		Endpoint:       env.GetString("S3_ENDPOINT", "http://minio:9000"),
		PublicEndpoint: env.GetString("S3_PUBLIC_URL", "http://localhost:9000"), // クライアントがアクセスする際のエンドポイント
		Region:         env.GetString("S3_REGION", "us-east-1"),                 // MinIOでは任意だが指定が必要
		BucketName:     env.GetString("S3_BUCKET", "my-video-app"),
		AccessKey:      env.GetString("S3_ACCESS_KEY", "minioadmin"),
		SecretKey:      env.GetString("S3_SECRET_KEY", "minioadmin"),

		// MinIOを使う場合、以下の2つは重要
		UsePathStyle: env.GetBool("S3_USE_PATH_STYLE", true), // MinIO は true 必須
		UseSSL:       env.GetBool("S3_USE_SSL", false),
	}
}
