package user

import (
	"time"

	"example.com/m/internal/domain/user/value"
	"github.com/google/uuid"
)

type User struct {
	id        uuid.UUID
	name      string
	bio       string
	iconKey   *string
	role      value.Role
	createdAt time.Time
}

func NewUser(
	id uuid.UUID,
	name string,
	bio string,
	iconKey *string,
	role value.Role,
	createdAt time.Time,
) *User {
	return &User{
		id:        id,
		name:      name,
		bio:       bio,
		iconKey:   iconKey,
		role:      role,
		createdAt: createdAt,
	}
}

func (u *User) ID() uuid.UUID {
	return u.id
}
func (u *User) Name() string {
	return u.name
}
func (u *User) Bio() string {
	return u.bio
}
func (u *User) IconKey() *string {
	return u.iconKey
}
func (u *User) Role() value.Role {
	return u.role
}
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// TODO: 名前のバリデーション
func (u *User) UpdateProfile(name string, bio string) error {
	u.name = name
	u.bio = bio
	return nil
}

func (u *User) UpdateIcon(iconKey *string) {
	u.iconKey = iconKey
}
