package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ot07/next-bazaar/test_util"
)

func TestMain(m *testing.M) {
	var err error

	dbConfig := test_util.DatabaseConfig{
		Image:            "postgres:15-alpine",
		Port:             5432,
		User:             "postgres",
		Password:         "secret",
		DBName:           "next-bazaar",
		DriverName:       "postgres",
		MigrateSourceURL: "file://../migration",
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
