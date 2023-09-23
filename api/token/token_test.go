package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidToken(t *testing.T) {
	duration := time.Minute
	expiredAt := time.Now().Add(duration)

	token := NewToken(duration)

	require.NotZero(t, token.ID)
	require.False(t, IsExpired(token.ExpiredAt))
	require.WithinDuration(t, expiredAt, token.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	token := NewToken(-time.Minute)
	require.NotEmpty(t, token)

	require.NotZero(t, token.ID)
	require.True(t, IsExpired(token.ExpiredAt))
}
