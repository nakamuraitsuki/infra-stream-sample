package value

import "errors"

type FailureReason string

const (
	FailureTranscode FailureReason = "transcode_failure" // トランスコード失敗
	FailureUnknown   FailureReason = "unknown_failure"   // 不明な失敗
)

var (
	ErrInvalidFailureReason = errors.New("invalid failure reason")
)

func NewFailureReason(s *string) (*FailureReason, error) {
	if s == nil {
		return nil, nil
	}
	switch *s {
	case string(FailureTranscode):
		fr := FailureTranscode
		return &fr, nil
	case string(FailureUnknown):
		fr := FailureUnknown
		return &fr, nil
	default:
		return nil, ErrInvalidFailureReason
	}
}
