package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T, testQueries *Queries) Session {
	user := createRandomUser(t, testQueries)

	arg := CreateSessionParams{
		UserID:       user.ID,
		SessionToken: util.RandomUUID(),
		ExpiredAt:    time.Now().Add(time.Minute),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.UserID, session.UserID)
	require.Equal(t, arg.SessionToken, session.SessionToken)
	require.WithinDuration(t, arg.ExpiredAt, session.ExpiredAt, time.Second)

	require.NotEmpty(t, session.ID)
	require.NotZero(t, session.CreatedAt)

	return session
}

func TestCreateSession(t *testing.T) {
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	createRandomSession(t, testQueries)
}

func TestGetSession(t *testing.T) {
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	session1 := createRandomSession(t, testQueries)
	session2, err := testQueries.GetSession(context.Background(), session1.SessionToken)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session1.ID, session2.ID)
	require.Equal(t, session1.UserID, session2.UserID)
	require.Equal(t, session1.SessionToken, session2.SessionToken)
	require.Equal(t, session1.ExpiredAt, session2.ExpiredAt)
	require.WithinDuration(t, session1.CreatedAt, session2.CreatedAt, time.Second)
}

func TestDeleteSession(t *testing.T) {
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	session1 := createRandomSession(t, testQueries)
	err := testQueries.DeleteSession(context.Background(), session1.ID)
	require.NoError(t, err)

	session2, err := testQueries.GetSession(context.Background(), session1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, session2)
}
