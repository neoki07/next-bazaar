package test_util

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"

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

func newTestDBContainer(ctx context.Context, config DatabaseConfig) (
	db *sql.DB,
	url string,
	purge func() error,
	err error,
) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", config.Password),
			fmt.Sprintf("POSTGRES_USER=%s", config.User),
			fmt.Sprintf("POSTGRES_DB=%s", config.DBName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort(fmt.Sprintf("%d/tcp", config.Port))
	url = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		config.User,
		config.Password,
		hostAndPort,
		config.DBName,
	)

	log.Println("Connecting to database on url:", url)

	resource.Expire(60) // Tell docker to hard kill the container in 60 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 60 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", url)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return nil, "", nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	purge = func() error {
		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("could not purge resource: %s", err)
		}
		return nil
	}

	return

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
func NewTestDB(ctx context.Context, config DatabaseConfig) (purge func() error, err error) {
	db, url, purge, err := newTestDBContainer(ctx, config)
	if err != nil {
		return nil, err
	}

	err = migrateUp(db, config)
	if err != nil {
		return nil, err
	}

	txdb.Register(testDBDriverName, config.DriverName, url)

	return
}

// OpenTestDB opens a new test database connection.
func OpenTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open(testDBDriverName, uuid.New().String())
	require.NoError(t, err)
	return db
}
