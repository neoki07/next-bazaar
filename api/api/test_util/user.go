package test_util

import (
	"context"
	"testing"

	"github.com/ot07/next-bazaar/util"

	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/stretchr/testify/require"
)

type WithSessionUserParams struct {
	Name         string
	Email        string
	Password     string
	SessionToken *token.Token
}

func CreateWithSessionUser(
	t *testing.T,
	ctx context.Context,
	store db.Store,
	params WithSessionUserParams,
) db.User {
	hashedPassword, err := util.HashPassword(params.Password)
	require.NoError(t, err)

	user, err := store.CreateUser(ctx, db.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	require.NoError(t, err)

	_, err = store.CreateSession(ctx, db.CreateSessionParams{
		UserID:                user.ID,
		SessionToken:          params.SessionToken.ID,
		SessionTokenExpiredAt: params.SessionToken.ExpiredAt,
	})
	require.NoError(t, err)

	return user
}
