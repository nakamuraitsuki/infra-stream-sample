package video

import (
	"context"
	"io"
	"time"
)

type ByteRange struct {
	Start int64
	End   *int64 // nilの場合はファイル末尾まで
}

type ObjectMeta struct {
	TotalSize     int64
	ContentLength int64
	RangeStart    int64
	RangeEnd      int64
	ETag          string
	LastModified  time.Time
}

type PartInfo struct {
	PartNumber int32
	ID         string
}

type Storage interface {
	// 元動画の保存（外部アップロード）
	// アップロードセッションを開始。セッションのIDを返す。
	StartUploadSession(ctx context.Context, key string) (string, error)

	// アップロードURLの生成。セッションIDとパート番号を指定して、アップロードURLを取得する。
	GenerateUploadPartURL(ctx context.Context, key string, sessionId string, partNum int32) (string, error)

	// アップロードセッションの完了
	CommitUploadSession(ctx context.Context, key string, sessionId string, parts []PartInfo) error

	// 変換後動画の保存（Transcodeプロセスからのアップロード）
	SaveStream(ctx context.Context, streamKey string, data io.Reader) error
	// 一時アクセスURLの取得
	GenerateTemporaryAccessURL(ctx context.Context, streamKey string, expires time.Duration) (string, error)

	GetStream(ctx context.Context, streamKey string, byteRange *ByteRange) (io.ReadCloser, *ObjectMeta, error)

	GetSource(ctx context.Context, sourceKey string) (io.ReadCloser, error)
	// 元動画の削除
	DeleteSource(ctx context.Context, sourceKey string) error
	// 変換後動画の削除
	DeleteStream(ctx context.Context, streamKey string) error
}
