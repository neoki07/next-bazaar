package test_util

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

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

type DatabaseConfig struct {
	Image      string
	Port       int
	User       string
	Password   string
	DBName     string
	DriverName string
}

func newTestContainer(ctx context.Context, config DatabaseConfig) (testcontainers.Container, nat.Port, error) {
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", config.Port)},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		Env: map[string]string{
			"POSTGRES_USER":     config.User,
			"POSTGRES_PASSWORD": config.Password,
			"POSTGRES_DB":       config.DBName,
		},
		WaitingFor: wait.ForListeningPort(
			nat.Port(
				fmt.Sprintf("%d/tcp", config.Port),
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

	mappedPort, err := testContainer.MappedPort(ctx, nat.Port(fmt.Sprintf("%d", config.Port)))
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
		config.DBName, driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}

func BeginTransaction(t *testing.T, testDB *sql.DB) *sql.Tx {
	tx, err := testDB.Begin()
	require.NoError(t, err)
	return tx
}

func RollbackTransaction(t *testing.T, tx *sql.Tx) {
	err := tx.Rollback()
	require.NoError(t, err)
}

func NewTestDB(ctx context.Context, config DatabaseConfig) (testDB *sql.DB, err error) {
	_, mappedPort, err := newTestContainer(ctx, config)
	if err != nil {
		return
	}

	sourceName := fmt.Sprintf("postgresql://%s:%s@127.0.0.1:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		mappedPort.Int(),
		config.DBName,
	)

	testDB, err = sql.Open(config.DriverName, sourceName)
	if err != nil {
		return
	}

	err = migrateUp(testDB, config)
	if err != nil {
		return
	}

	return
}
