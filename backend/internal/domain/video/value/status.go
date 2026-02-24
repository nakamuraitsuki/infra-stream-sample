package value

import "errors"

type Status string

const (
	StatusInitial    Status = "initial"    // 初期状態
	StatusUploaded   Status = "uploaded"   // アップロード完了
	StatusQueued     Status = "queued"     // 非同期処理待ち
	StatusProcessing Status = "processing" // 非同期処理中
	StatusReady      Status = "ready"      // 再生可能
	StatusFailed     Status = "failed"     // 処理失敗
)

var (
	ErrInvalidStatus = errors.New("invalid status")
)

func NewStatus(s string) (Status, error) {
	switch s {
	case string(StatusInitial),
		string(StatusUploaded),
		string(StatusQueued),
		string(StatusProcessing),
		string(StatusReady),
		string(StatusFailed):
		return Status(s), nil
	default:
		return "", ErrInvalidStatus
	}
}
