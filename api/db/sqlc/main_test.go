package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ot07/next-bazaar/test_util"
)

var testDB *sql.DB

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
	defer cancel()

	testDB, err = test_util.NewTestDB(ctx, dbConfig)
	if err != nil {
		log.Fatal("cannot create test db:", err)
	}

	os.Exit(m.Run())
}
