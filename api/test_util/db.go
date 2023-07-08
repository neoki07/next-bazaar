package test_util

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const testDBDriverName = "txdb-api"

type DatabaseConfig struct {
	Image            string
	Port             int
	User             string
	Password         string
	DBName           string
	DriverName       string
	MigrateSourceURL string
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
		config.MigrateSourceURL,
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

// NewTestDB creates a new test database.
func NewTestDB(ctx context.Context, config DatabaseConfig) error {
	_, mappedPort, err := newTestContainer(ctx, config)
	if err != nil {
		return err
	}

	sourceName := fmt.Sprintf("postgresql://%s:%s@127.0.0.1:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		mappedPort.Int(),
		config.DBName,
	)

	testDB, err := sql.Open(config.DriverName, sourceName)
	if err != nil {
		return err
	}

	err = migrateUp(testDB, config)
	if err != nil {
		return err
	}

	txdb.Register(testDBDriverName, config.DriverName, sourceName)

	return nil
}

// OpenTestDB opens a new test database connection.
func OpenTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open(testDBDriverName, uuid.New().String())
	require.NoError(t, err)
	return db
}
