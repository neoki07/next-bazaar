package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T, testQueries *Queries) Category {
	name := util.RandomName()

	category, err := testQueries.CreateCategory(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.Equal(t, name, category.Name)

	require.NotEmpty(t, category.ID)
	require.NotZero(t, category.CreatedAt)

	return category
}

func TestCreateCategory(t *testing.T) {
	t.Parallel()

	db, err := sql.Open(testDBDriverName, uuid.New().String())
	require.NoError(t, err)
	defer db.Close()

	testQueries := New(db)

	createRandomCategory(t, testQueries)
}

func TestGetCategory(t *testing.T) {
	t.Parallel()

	db, err := sql.Open(testDBDriverName, uuid.New().String())
	require.NoError(t, err)
	defer db.Close()

	testQueries := New(db)

	category1 := createRandomCategory(t, testQueries)
	category2, err := testQueries.GetCategory(context.Background(), category1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.Equal(t, category1.ID, category2.ID)
	require.Equal(t, category1.Name, category2.Name)
	require.WithinDuration(t, category1.CreatedAt, category2.CreatedAt, time.Second)
}

func TestListCategories(t *testing.T) {
	t.Parallel()

	db, err := sql.Open(testDBDriverName, uuid.New().String())
	require.NoError(t, err)
	defer db.Close()

	testQueries := New(db)

	for i := 0; i < 10; i++ {
		createRandomCategory(t, testQueries)
	}

	arg := ListCategoriesParams{Limit: 5, Offset: 5}

	categories, err := testQueries.ListCategories(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, categories, 5)

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}

func TestDeleteCategory(t *testing.T) {
	t.Parallel()

	db, err := sql.Open(testDBDriverName, uuid.New().String())
	require.NoError(t, err)
	defer db.Close()

	testQueries := New(db)

	category1 := createRandomCategory(t, testQueries)
	err = testQueries.DeleteCategory(context.Background(), category1.ID)
	require.NoError(t, err)

	category2, err := testQueries.GetCategory(context.Background(), category1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, category2)
}
