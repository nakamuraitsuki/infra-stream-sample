package event

import "context"

// UseCase を回すための Relay インターフェース
type Relay interface {
	Start(ctx context.Context) error
}
