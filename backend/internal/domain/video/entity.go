package video

import (
	"errors"
	"time"

	"example.com/m/internal/domain/event"
	"example.com/m/internal/domain/video/value"
	"github.com/google/uuid"
)

type Video struct {
	id            uuid.UUID
	ownerID       uuid.UUID
	sourceKey     string // 元動画の保存先キー
	streamKey     string // 配信・再生用動画の保存先キー
	status        value.Status
	title         string
	description   string
	tags          []value.Tag
	retryCount    int
	failureReason *value.FailureReason
	visibility    value.Visibility
	createdAt     time.Time
	events        []event.Event
}

func NewVideo(
	id uuid.UUID,
	ownerID uuid.UUID,
	sourceKey string,
	streamKey string,
	status value.Status,
	title string,
	description string,
	tags []value.Tag,
	retryCount int,
	failureReason *value.FailureReason,
	visibility value.Visibility,
	createdAt time.Time,
	events []event.Event,
) *Video {
	return &Video{
		id:            id,
		ownerID:       ownerID,
		sourceKey:     sourceKey,
		streamKey:     streamKey,
		status:        status,
		title:         title,
		description:   description,
		tags:          tags,
		retryCount:    retryCount,
		failureReason: failureReason,
		visibility:    visibility,
		createdAt:     createdAt,
		events:        events,
	}
}

func (v *Video) ID() uuid.UUID {
	return v.id
}

func (v *Video) OwnerID() uuid.UUID {
	return v.ownerID
}

func (v *Video) Title() string {
	return v.title
}

func (v *Video) Description() string {
	return v.description
}

func (v *Video) Tags() []value.Tag {
	return v.tags
}

func (v *Video) RetryCount() int {
	return v.retryCount
}

func (v *Video) FailureReason() *value.FailureReason {
	return v.failureReason
}

func (v *Video) PullEvents() []event.Event {
	events := v.events
	v.events = nil
	return events
}

func (v *Video) CreatedAt() time.Time {
	return v.createdAt
}

func (v *Video) Status() value.Status {
	return v.status
}

func (v *Video) SourceKey() string {
	return v.sourceKey
}

func (v *Video) StreamKey() string {
	return v.streamKey
}

func (v *Video) Visibility() value.Visibility {
	return v.visibility
}

func (v *Video) MarkUploaded(sourceKey string) error {
	if v.status != value.StatusInitial {
		return errors.New("video is not in initial state")
	}
	v.sourceKey = sourceKey
	v.status = value.StatusUploaded
	return nil
}

func (v *Video) StartTranscoding(streamKey string) error {
	if v.status != value.StatusUploaded {
		return errors.New("video is not transcodable")
	}
	v.streamKey = streamKey
	v.status = value.StatusProcessing
	v.events = append(v.events, event.NewTranscodingRequested(v.id))
	return nil
}

func (v *Video) RollbackToUploaded() error {
	if v.status != value.StatusProcessing {
		return errors.New("video is not in processing state")
	}
	v.status = value.StatusUploaded
	return nil
}

func (v *Video) MarkTranscodeSucceeded() {
	v.status = value.StatusReady
	v.failureReason = nil
	v.retryCount = 0
}

func (v *Video) MarkTranscodeFailed(reason value.FailureReason) {
	v.status = value.StatusFailed
	v.failureReason = &reason
	v.retryCount++
}
