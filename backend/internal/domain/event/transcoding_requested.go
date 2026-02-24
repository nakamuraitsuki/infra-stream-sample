package event

import (
	"time"

	"github.com/google/uuid"
)

const TranscodingRequestType = "video.transcoding_requested"

// Domain Model 内で使用するイベントの定義
type TranscodingRequested struct {
	// 全て大文字開始にして JSON タグを付与する（Marshal 時に無視されないようにするため）
	EventID   uuid.UUID `json:"id"`
	VideoID   uuid.UUID `json:"video_id"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewTranscodingRequested(videoID uuid.UUID) Event {
	return &TranscodingRequested{
		EventID:   uuid.New(),
		VideoID:   videoID,
		Timestamp: time.Now(),
	}
}

// インターフェースを満たすためのメソッド
func (e *TranscodingRequested) ID() uuid.UUID         { return e.EventID }
func (e *TranscodingRequested) EventType() string     { return TranscodingRequestType }
func (e *TranscodingRequested) OccurredAt() time.Time { return e.Timestamp }
func (e *TranscodingRequested) Payload() any {
	return e
}
