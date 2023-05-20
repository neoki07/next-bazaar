package user_domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID
	Name              string
	Email             string
	HashedPassword    string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
}

func NewUser(
	id uuid.UUID,
	name string,
	email string,
	hashedPassword string,
	passwordChangedAt time.Time,
	createdAt time.Time,
) *User {
	return &User{
		ID:                id,
		Name:              name,
		Email:             email,
		HashedPassword:    hashedPassword,
		PasswordChangedAt: passwordChangedAt,
		CreatedAt:         createdAt,
	}

}
