package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ot07/next-bazaar/test_util"
)

const testDBDriverName = "txdb-db"

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

	sourceName, err := test_util.InitTestDB(ctx, dbConfig)
	if err != nil {
		log.Fatal("cannot create test db:", err)
	}

	txdb.Register(testDBDriverName, dbConfig.DriverName, sourceName)

	os.Exit(m.Run())
}
