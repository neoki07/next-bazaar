package test_util

import (
	"context"
	"testing"
	"time"

	"github.com/ot07/next-bazaar/util"

	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/stretchr/testify/require"
)

func CreateUserTestData(
	t *testing.T,
	ctx context.Context,
	store db.Store,
	name string,
	email string,
	password string,
) (db.User, *token.Token) {
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user, err := store.CreateUser(ctx, db.CreateUserParams{
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
	})
	require.NoError(t, err)

	sessionToken := token.NewToken(time.Minute)

	_, err = store.CreateSession(ctx, db.CreateSessionParams{
		UserID:       user.ID,
		SessionToken: sessionToken.ID,
		ExpiredAt:    sessionToken.ExpiredAt,
	})
	require.NoError(t, err)

	return user, sessionToken
}
