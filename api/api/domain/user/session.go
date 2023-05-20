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

func NewSession(
	id uuid.UUID,
	sessionToken uuid.UUID,
	userID uuid.UUID,
	expiredAt time.Time,
	createdAt time.Time,
) *Session {
	return &Session{
		ID:           id,
		SessionToken: sessionToken,
		UserID:       userID,
		ExpiredAt:    expiredAt,
		CreatedAt:    createdAt,
	}
}
