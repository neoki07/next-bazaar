package test_util

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	mockdb "github.com/ot07/next-bazaar/db/mock"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

const testDBDriverName = "txdb-api"

func NewTestDBStore(t *testing.T) (store *db.SQLStore, cleanup func()) {
	conn, err := sql.Open(testDBDriverName, uuid.New().String())
	require.NoError(t, err)
	return db.NewStore(conn), func() { conn.Close() }
}

func BuildTestDBStore(t *testing.T) (store db.Store, cleanup func()) {
	return NewTestDBStore(t)
}

func NewMockStore(t *testing.T) (store *mockdb.MockStore, cleanup func()) {
	ctrl := gomock.NewController(t)
	return mockdb.NewMockStore(ctrl), func() { ctrl.Finish() }
}
