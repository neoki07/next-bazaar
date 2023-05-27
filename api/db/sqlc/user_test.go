package db

import (
	"context"
	"testing"
	"time"

	"github.com/ot07/next-bazaar/test_util"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T, testQueries *Queries) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Name:           util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.NotEmpty(t, user.ID)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tx := test_util.BeginTransaction(t, testDB)
	defer test_util.RollbackTransaction(t, tx)

	testQueries := New(tx)

	createRandomUser(t, testQueries)
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	tx := test_util.BeginTransaction(t, testDB)
	defer test_util.RollbackTransaction(t, tx)

	testQueries := New(tx)

	user1 := createRandomUser(t, testQueries)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()

	tx := test_util.BeginTransaction(t, testDB)
	defer test_util.RollbackTransaction(t, tx)

	testQueries := New(tx)

	user1 := createRandomUser(t, testQueries)
	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
