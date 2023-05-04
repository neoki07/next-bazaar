package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

type DatabaseConfig struct {
	image      string
	port       int
	user       string
	password   string
	dbName     string
	driverName string
}

func newTestContainer(ctx context.Context, config DatabaseConfig) (testcontainers.Container, nat.Port, error) {
	req := testcontainers.ContainerRequest{
		Image:        config.image,
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", config.port)},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		Env: map[string]string{
			"POSTGRES_USER":     config.user,
			"POSTGRES_PASSWORD": config.password,
			"POSTGRES_DB":       config.dbName,
		},
		WaitingFor: wait.ForListeningPort(
			nat.Port(
				fmt.Sprintf("%d/tcp", config.port),
			),
		),
	}
	testContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	mappedPort, err := testContainer.MappedPort(ctx, nat.Port(fmt.Sprintf("%d", config.port)))
	if err != nil {
		return nil, "", err
	}

	return testContainer, mappedPort, nil
}

func migrateUp(db *sql.DB, config DatabaseConfig) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../migration",
		config.dbName, driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}

func beginTransaction(t *testing.T) *sql.Tx {
	tx, err := testDB.Begin()
	require.NoError(t, err)
	return tx
}

func rollbackTransaction(t *testing.T, tx *sql.Tx) {
	err := tx.Rollback()
	require.NoError(t, err)
}

func TestMain(m *testing.M) {
	dbConfig := DatabaseConfig{
		image:      "postgres:15-alpine",
		port:       5432,
		user:       "postgres",
		password:   "secret",
		dbName:     "coworker",
		driverName: "postgres",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, mappedPort, err := newTestContainer(ctx, dbConfig)
	if err != nil {
		log.Fatal("cannot create container for testing:", err)
	}

	sourceName := fmt.Sprintf("postgresql://%s:%s@127.0.0.1:%d/%s?sslmode=disable",
		dbConfig.user,
		dbConfig.password,
		mappedPort.Int(),
		dbConfig.dbName,
	)
	testDB, err = sql.Open(dbConfig.driverName, sourceName)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = migrateUp(testDB, dbConfig)
	if err != nil {
		log.Fatal("cannot migrate up:", err)
	}

	os.Exit(m.Run())
}
