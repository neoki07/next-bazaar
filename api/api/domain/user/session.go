package user_domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID
	SessionToken uuid.UUID
	UserID       uuid.UUID
	ExpiredAt    time.Time
	CreatedAt    time.Time
}
