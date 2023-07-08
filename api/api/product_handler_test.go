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

	"github.com/gofiber/fiber/v2"

	"github.com/golang/mock/gomock"
	product_domain "github.com/ot07/next-bazaar/api/domain/product"
	"github.com/ot07/next-bazaar/api/test_util"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestGetProduct(t *testing.T) {
	testCases := []struct {
		name          string
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) (productID string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) (productID string) {
				ctx := context.Background()

				user := test_util.CreateUserTestData(t, ctx, store,
					"testuser",
					"test@example.com",
					"test-password",
					token.NewToken(time.Minute),
				)

				createdCategory, err := store.CreateCategory(ctx, "test-category")
				require.NoError(t, err)

				createdProduct, err := store.CreateProduct(ctx, db.CreateProductParams{
					Name:          "test-product",
					Description:   sql.NullString{String: "test-description", Valid: true},
					Price:         "100.00",
					StockQuantity: 10,
					CategoryID:    createdCategory.ID,
					SellerID:      user.ID,
					ImageUrl:      sql.NullString{String: "test-image-url", Valid: true},
				})
				require.NoError(t, err)

				return createdProduct.ID.String()
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotProduct := unmarshalProductResponse(t, response.Body)

				require.NotEmpty(t, gotProduct.ID)
				require.Equal(t, "test-product", gotProduct.Name)
				require.Equal(t, "test-description", gotProduct.Description.String)
				require.True(t, decimal.NewFromFloat(100.00).Equal(gotProduct.Price.Decimal))
				require.Equal(t, int32(10), gotProduct.StockQuantity)
				require.Equal(t, "test-category", gotProduct.Category)
				require.Equal(t, "testuser", gotProduct.Seller)
				require.Equal(t, "test-image-url", gotProduct.ImageUrl.String)
			},
		},
		{
			name:       "NotFound",
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) (productID string) {
				return util.RandomUUID().String()
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name:       "InvalidID",
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) (productID string) {
				return "InvalidID"
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Return(db.Product{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: func(t *testing.T, store db.Store) (productID string) {
				return util.RandomUUID().String()
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store, cleanupStore := tc.buildStore(t)
			defer cleanupStore()

			productID := tc.createSeed(t, store)

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/products/%s", productID),
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestListProducts(t *testing.T) {
	pageSize := 5

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: pageSize,
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) {
				var err error

				ctx := context.Background()

				users := make([]db.User, 2)
				for i := range users {
					users[i] = test_util.CreateUserTestData(t, ctx, store,
						fmt.Sprintf("testuser-%d", i),
						fmt.Sprintf("test-%d@example.com", i),
						"test-password",
						token.NewToken(time.Minute),
					)
				}

				createdCategories := make([]db.Category, 3)
				for i := range createdCategories {
					createdCategories[i], err = store.CreateCategory(ctx, fmt.Sprintf("test-category-%d", i))
					require.NoError(t, err)
				}

				for i := 0; i < 6; i++ {
					_, err = store.CreateProduct(ctx, db.CreateProductParams{
						Name:          fmt.Sprintf("test-product-%d", i),
						Description:   sql.NullString{String: fmt.Sprintf("test-description-%d", i), Valid: true},
						Price:         fmt.Sprintf("%d.00", (i+1)*10),
						StockQuantity: int32(i + 1),
						CategoryID:    createdCategories[i%3].ID,
						SellerID:      users[i%2].ID,
						ImageUrl:      sql.NullString{String: fmt.Sprintf("test-image-url-%d", i), Valid: true},
					})
					require.NoError(t, err)
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotResponse := unmarshalListProductsResponse(t, response.Body)

				require.Equal(t, int32(1), gotResponse.Meta.PageID)
				require.Equal(t, int32(pageSize), gotResponse.Meta.PageSize)
				require.Equal(t, int64(2), gotResponse.Meta.PageCount)
				require.Equal(t, int64(6), gotResponse.Meta.TotalCount)

				require.Len(t, gotResponse.Data, pageSize)

				for i := 0; i < pageSize; i++ {
					userIndex := i % 2
					categoryIndex := i % 3

					require.NotEmpty(t, gotResponse.Data[i].ID)
					require.Equal(t, fmt.Sprintf("test-product-%d", i), gotResponse.Data[i].Name)
					require.Equal(t, fmt.Sprintf("test-description-%d", i), gotResponse.Data[i].Description.String)
					require.True(t, decimal.NewFromInt(int64((i+1)*10)).Equal(gotResponse.Data[i].Price.Decimal))
					require.Equal(t, int32(i+1), gotResponse.Data[i].StockQuantity)
					require.Equal(t, fmt.Sprintf("test-category-%d", categoryIndex), gotResponse.Data[i].Category)
					require.Equal(t, fmt.Sprintf("testuser-%d", userIndex), gotResponse.Data[i].Seller)
					require.Equal(t, fmt.Sprintf("test-image-url-%d", i), gotResponse.Data[i].ImageUrl.String)
				}
			},
		},
		{
			name: "PageIDNotFound",
			query: Query{
				pageSize: pageSize,
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageIDLessThanLowerLimit",
			query: Query{
				pageID:   0,
				pageSize: pageSize,
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeNotFound",
			query: Query{
				pageID: 1,
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeLessThanLowerLimit",
			query: Query{
				pageID:   1,
				pageSize: 0,
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeMoreThanUpperLimit",
			query: Query{
				pageID:   1,
				pageSize: 101,
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				pageID:   1,
				pageSize: 1,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Return([]db.Product{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store, cleanupStore := tc.buildStore(t)
			defer cleanupStore()

			tc.createSeed(t, store)

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodGet,
				URL:    "/api/v1/products",
				Query: fiber.Map{
					"page_id":   fmt.Sprintf("%d", tc.query.pageID),
					"page_size": fmt.Sprintf("%d", tc.query.pageSize),
				},
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})

	}
}

func unmarshalProductResponse(t *testing.T, body io.ReadCloser) product_domain.ProductResponse {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var parsed product_domain.ProductResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	return parsed
}

func unmarshalListProductsResponse(t *testing.T, body io.ReadCloser) product_domain.ListProductsResponse {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var parsed product_domain.ListProductsResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	return parsed
}
