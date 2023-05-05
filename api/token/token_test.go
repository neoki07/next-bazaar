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

	err := token.Valid()
	require.NoError(t, err)

	require.NotZero(t, token.ID)
	require.WithinDuration(t, expiredAt, token.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	token := NewToken(-time.Minute)
	require.NotEmpty(t, token)

	err := token.Valid()
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
}
