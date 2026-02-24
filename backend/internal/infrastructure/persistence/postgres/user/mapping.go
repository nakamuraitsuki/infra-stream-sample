package user

import (
	"time"

	user_domain "example.com/m/internal/domain/user"
	user_value "example.com/m/internal/domain/user/value"
	"github.com/google/uuid"
)

type userModel struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Bio       string    `db:"bio"`
	IconKey   *string   `db:"icon_key"`
	Role      string    `db:"role"` // value.Role を文字列として保存
	CreatedAt time.Time `db:"created_at"`
}

func fromEntity(u *user_domain.User) *userModel {
	return &userModel{
		ID:        u.ID(),
		Name:      u.Name(),
		Bio:       u.Bio(),
		IconKey:   u.IconKey(),
		Role:      string(u.Role()),
		CreatedAt: u.CreatedAt(),
	}
}

// toEntity は Userのコンストラクタのラッパー
// DB特有の不正値のチェックなどをしたい場合に追記する
func (m *userModel) toEntity() (*user_domain.User, error) {
	return user_domain.NewUser(
		m.ID,
		m.Name,
		m.Bio,
		m.IconKey,
		user_value.Role(m.Role),
		m.CreatedAt,
	), nil
}
