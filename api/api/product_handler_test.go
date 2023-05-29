package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	product_domain "github.com/ot07/next-bazaar/api/domain/product"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func createSeed(t *testing.T, store db.Store, user db.User, category product_domain.Category, product product_domain.Product) (productId uuid.UUID) {
	ctx := context.Background()

	createdUser, err := store.CreateUser(ctx, db.CreateUserParams{
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	})
	require.NoError(t, err)

	createdCategory, err := store.CreateCategory(ctx, category.Name)
	require.NoError(t, err)

	createdProduct, err := store.CreateProduct(ctx, db.CreateProductParams{
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		CategoryID:    createdCategory.ID,
		SellerID:      createdUser.ID,
		ImageUrl:      product.ImageUrl,
	})
	require.NoError(t, err)

	return createdProduct.ID
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	u, _ := randomUser(t)
	c := randomCategory(t)
	p := randomProduct(t, u, c)

	testCases := []struct {
		name          string
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchProduct(t, response.Body, p)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store, cleanupStore := tc.buildStore(t)
			defer cleanupStore()

			productId := createSeed(t, store, u, c, p)

			server := newTestServer(t, store)

			url := fmt.Sprintf("/api/v1/products/%s", productId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func randomProduct(t *testing.T, user db.User, category product_domain.Category) product_domain.Product {
	price, err := util.RandomMoney()
	require.NoError(t, err)

	return product_domain.Product{
		Name:          util.RandomName(),
		Description:   sql.NullString{String: util.RandomString(30), Valid: true},
		Price:         price,
		StockQuantity: util.RandomInt32(10),
		CategoryID:    category.ID,
		Category:      category.Name,
		SellerID:      user.ID,
		Seller:        user.Name,
		ImageUrl:      sql.NullString{String: util.RandomImageUrl(), Valid: true},
	}
}

func randomCategory(t *testing.T) product_domain.Category {
	return product_domain.Category{
		Name: util.RandomName(),
	}
}

func requireBodyMatchProduct(t *testing.T, body io.ReadCloser, product product_domain.Product) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotProduct product_domain.ProductResponse
	err = json.Unmarshal(data, &gotProduct)
	require.NoError(t, err)

	require.NotEmpty(t, gotProduct.ID)
	requireProductResponseMatchProduct(t, gotProduct, product)

	err = body.Close()
	require.NoError(t, err)
}

func requireProductResponseMatchProduct(t *testing.T, gotProduct product_domain.ProductResponse, product product_domain.Product) {
	require.Equal(t, product.Name, gotProduct.Name)
	require.Equal(t, product.Description, gotProduct.Description.NullString)
	require.Equal(t, product.Price, gotProduct.Price)
	require.Equal(t, product.StockQuantity, gotProduct.StockQuantity)
	require.Equal(t, product.Category, gotProduct.Category)
	require.Equal(t, product.Seller, gotProduct.Seller)
	require.Equal(t, product.ImageUrl, gotProduct.ImageUrl.NullString)
}
