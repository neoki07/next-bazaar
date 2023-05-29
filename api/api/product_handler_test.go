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

type product struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Price         string
	StockQuantity int32
	CategoryID    uuid.UUID
	SellerID      uuid.UUID
	ImageUrl      string
}

type category struct {
	ID   uuid.UUID
	Name string
}

func createSeed(t *testing.T, store db.Store, user db.User, category category, product product) (productId uuid.UUID) {
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
		Description:   sql.NullString{String: product.Description, Valid: true},
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		CategoryID:    createdCategory.ID,
		SellerID:      createdUser.ID,
		ImageUrl:      sql.NullString{String: product.ImageUrl, Valid: true},
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

func randomProduct(t *testing.T, user db.User, category category) product {
	price, err := util.RandomMoney()
	require.NoError(t, err)

	return product{
		Name:          util.RandomName(),
		Description:   util.RandomString(30),
		Price:         price,
		StockQuantity: util.RandomInt32(10),
		ImageUrl:      util.RandomImageUrl(),
	}
}

func randomCategory(t *testing.T) category {
	return category{
		Name: util.RandomName(),
	}
}

func requireBodyMatchProduct(t *testing.T, body io.ReadCloser, product product) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotProduct product_domain.ProductResponse
	err = json.Unmarshal(data, &gotProduct)
	require.NoError(t, err)

	require.NotEmpty(t, gotProduct.ID)
	requireProductResponseMatchUser(t, gotProduct, product)

	err = body.Close()
	require.NoError(t, err)
}

func requireProductResponseMatchUser(t *testing.T, gotProduct product_domain.ProductResponse, product product) {
	require.Equal(t, product.Name, gotProduct.Name)
	require.Equal(t, product.Description, gotProduct.Description.String)
}
