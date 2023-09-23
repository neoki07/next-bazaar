package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/util"
)

var (
	ErrExpiredToken = errors.New("token has expired")
)

// Token contains the token and expired at
type Token struct {
	ID        uuid.UUID `json:"id"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewToken creates a new token with a specific duration
func NewToken(duration time.Duration) *Token {
	tokenID := util.RandomUUID()

	token := &Token{
		ID:        tokenID,
		ExpiredAt: time.Now().Add(duration),
	}
	return token
}

// IsExpired checks if the `expiredAt` is expired or not
func IsExpired(expiredAt time.Time) bool {
	return time.Now().After(expiredAt)
}
