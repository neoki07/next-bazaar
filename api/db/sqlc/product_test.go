package db

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/ot07/next-bazaar/test_util"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func createRandomProduct(t *testing.T, testQueries *Queries) Product {
	price := util.RandomPrice()

	category := createRandomCategory(t, testQueries)
	user := createRandomUser(t, testQueries)

	arg := CreateProductParams{
		Name:          util.RandomName(),
		Description:   sql.NullString{String: util.RandomName(), Valid: true},
		Price:         price.String(),
		StockQuantity: rand.Int31n(100),
		CategoryID:    category.ID,
		SellerID:      user.ID,
		ImageUrl:      sql.NullString{String: util.RandomImageUrl(), Valid: true},
	}

	product, err := testQueries.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.Equal(t, arg.Name, product.Name)
	require.Equal(t, arg.Description, product.Description)
	require.Equal(t, arg.Price, product.Price)
	require.Equal(t, arg.StockQuantity, product.StockQuantity)
	require.Equal(t, arg.CategoryID, product.CategoryID)
	require.Equal(t, arg.SellerID, product.SellerID)
	require.Equal(t, arg.ImageUrl, product.ImageUrl)

	require.NotEmpty(t, product.ID)
	require.NotZero(t, product.CreatedAt)

	return product
}

func TestCreateProduct(t *testing.T) {
	t.Parallel()

	db := test_util.OpenTestDB(t)
	defer db.Close()

	testQueries := New(db)

	createRandomProduct(t, testQueries)
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	db := test_util.OpenTestDB(t)
	defer db.Close()

	testQueries := New(db)

	product1 := createRandomProduct(t, testQueries)
	product2, err := testQueries.GetProduct(context.Background(), product1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ID, product2.ID)
	require.Equal(t, product1.Name, product2.Name)
	require.Equal(t, product1.Description, product2.Description)
	require.Equal(t, product1.Price, product2.Price)
	require.Equal(t, product1.StockQuantity, product2.StockQuantity)
	require.Equal(t, product1.CategoryID, product2.CategoryID)
	require.Equal(t, product1.SellerID, product2.SellerID)
	require.Equal(t, product1.ImageUrl, product2.ImageUrl)
	require.WithinDuration(t, product1.CreatedAt, product2.CreatedAt, time.Second)
}
