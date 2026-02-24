package value

import "errors"

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

var (
	ErrInvalidVisibility = errors.New("invalid visibility")
)

func NewVisibility(s string) (Visibility, error) {
	switch s {
	case string(VisibilityPublic),
		string(VisibilityPrivate):
		return Visibility(s), nil
	default:
		return "", ErrInvalidVisibility
	}
}
