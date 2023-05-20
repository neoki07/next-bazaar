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
