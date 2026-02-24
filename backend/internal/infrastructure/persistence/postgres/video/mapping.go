package video

import (
	"time"

	"example.com/m/internal/domain/event"
	video_domain "example.com/m/internal/domain/video"
	"example.com/m/internal/domain/video/value"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type videoModel struct {
	ID            uuid.UUID `db:"id"`
	OwnerID       uuid.UUID `db:"owner_id"`
	SourceKey     string    `db:"source_key"`
	StreamKey     string    `db:"stream_key"`
	Status        string    `db:"status"`
	Title         string    `db:"title"`
	Description   string    `db:"description"`
	RetryCount    int       `db:"retry_count"`
	FailureReason *string   `db:"failure_reason"`
	Visibility    string    `db:"visibility"`
	CreatedAt     time.Time `db:"created_at"`
	// Tags はJOINで取得するため，テーブルには含めない
	Tags pq.StringArray `db:"tags"`
	// event は Outbox パターンで使用するため除外
}

// 構造として明示的に定義しておく（クエリで暗黙的に使われていても）
type videoTagModel struct {
	VideoID int64     `db:"video_id"`
	TagID   uuid.UUID `db:"tag_id"`
}

func fromEntity(v *video_domain.Video) *videoModel {
	var failureReason *string
	if v.FailureReason() != nil {
		fr := string(*v.FailureReason())
		failureReason = &fr
	}

	tagsStr := make([]string, len(v.Tags()))
	for i, tag := range v.Tags() {
		tagsStr[i] = string(tag)
	}

	return &videoModel{
		ID:            v.ID(),
		OwnerID:       v.OwnerID(),
		SourceKey:     v.SourceKey(),
		StreamKey:     v.StreamKey(),
		Status:        string(v.Status()),
		Title:         v.Title(),
		Description:   v.Description(),
		RetryCount:    v.RetryCount(),
		FailureReason: failureReason,
		Visibility:    string(v.Visibility()),
		CreatedAt:     v.CreatedAt(),
		Tags:          tagsStr,
	}
}

// toEntity はDomainコンストラクタのラッパーとして運用する
// DB から変換の際のバリデーションを担保するため
func (m *videoModel) toEntity() (*video_domain.Video, error) {
	tags := make([]value.Tag, len(m.Tags))
	for i, tagStr := range m.Tags {
		tag, err := value.NewTag(tagStr)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}

	failureReason, err := value.NewFailureReason(m.FailureReason)
	if err != nil {
		return nil, err
	}

	status, err := value.NewStatus(m.Status)
	if err != nil {
		return nil, err
	}

	visibility, err := value.NewVisibility(m.Visibility)
	if err != nil {
		return nil, err
	}

	return video_domain.NewVideo(
		m.ID,
		m.OwnerID,
		m.SourceKey,
		m.StreamKey,
		status,
		m.Title,
		m.Description,
		tags,
		m.RetryCount,
		failureReason,
		visibility,
		m.CreatedAt,
		[]event.Event{},
	), nil

}
