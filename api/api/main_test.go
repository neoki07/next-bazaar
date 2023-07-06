package api

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/ot07/next-bazaar/db/mock"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/test_util"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

const testDBDriverName = "txdb-api"

func newTestDBStore(t *testing.T) (store *db.SQLStore, cleanup func()) {
	conn, err := sql.Open(testDBDriverName, uuid.New().String())
	require.NoError(t, err)
	return db.NewStore(conn), func() { conn.Close() }
}

func buildTestDBStore(t *testing.T) (store db.Store, cleanup func()) {
	return newTestDBStore(t)
}

func newMockStore(t *testing.T) (store *mockdb.MockStore, cleanup func()) {
	ctrl := gomock.NewController(t)
	return mockdb.NewMockStore(ctrl), func() { ctrl.Finish() }
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		SessionTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func createTestUserSeed(
	t *testing.T,
	ctx context.Context,
	store db.Store,
	name string,
	email string,
	hashedPassword string,
	sessionToken *token.Token,
) db.User {
	user, err := store.CreateUser(ctx, db.CreateUserParams{
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
	})
	require.NoError(t, err)

	_, err = store.CreateSession(ctx, db.CreateSessionParams{
		UserID:       user.ID,
		SessionToken: sessionToken.ID,
		ExpiredAt:    sessionToken.ExpiredAt,
	})
	require.NoError(t, err)

	return user
}

func TestMain(m *testing.M) {
	var err error

	dbConfig := test_util.DatabaseConfig{
		Image:            "postgres:15-alpine",
		Port:             5432,
		User:             "postgres",
		Password:         "secret",
		DBName:           "next-bazaar",
		DriverName:       "postgres",
		MigrateSourceURL: "file://../db/migration",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sourceName, err := test_util.NewTestDB(ctx, dbConfig)
	if err != nil {
		log.Fatal("cannot create test db:", err)
	}

	txdb.Register(testDBDriverName, dbConfig.DriverName, sourceName)

	os.Exit(m.Run())
}
