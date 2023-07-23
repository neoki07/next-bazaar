package api

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/test_util"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		SessionTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
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

	purge, err := test_util.NewTestDB(ctx, dbConfig)
	if err != nil {
		log.Fatal("cannot create test db:", err)
	}

	code := m.Run()

	if err := purge(); err != nil {
		log.Fatal("cannot purge test db:", err)
	}

	cancel()

	os.Exit(code)
}
