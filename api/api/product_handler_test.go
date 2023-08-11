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
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store) test_util.SeedData
		checkResponse  func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			buildStore: test_util.BuildTestDBStore,
			createSeedData: func(t *testing.T, store db.Store) test_util.SeedData {
				ctx := context.Background()

				user := test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
					Name:         "testuser",
					Email:        "test@example.com",
					Password:     "test-password",
					SessionToken: token.NewToken(time.Minute),
				})

				category, err := store.CreateCategory(ctx, "test-category")
				require.NoError(t, err)

				product, err := store.CreateProduct(ctx, db.CreateProductParams{
					Name:          "test-product",
					Description:   sql.NullString{String: "test-description", Valid: true},
					Price:         "100.00",
					StockQuantity: 10,
					CategoryID:    category.ID,
					SellerID:      user.ID,
					ImageUrl:      sql.NullString{String: "test-image-url", Valid: true},
				})
				require.NoError(t, err)

				return test_util.SeedData{
					"product_id": product.ID.String(),
				}
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
			createSeedData: func(t *testing.T, store db.Store) test_util.SeedData {
				return test_util.SeedData{
					"product_id": util.RandomUUID().String(),
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name:       "InvalidID",
			buildStore: test_util.BuildTestDBStore,
			createSeedData: func(t *testing.T, store db.Store) test_util.SeedData {
				return test_util.SeedData{
					"product_id": "InvalidID",
				}
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
			createSeedData: func(t *testing.T, store db.Store) test_util.SeedData {
				return test_util.SeedData{
					"product_id": util.RandomUUID().String(),
				}
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

			seedData := tc.createSeedData(t, store)

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/products/%s", seedData["product_id"].(string)),
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestListProducts(t *testing.T) {
	pageSize := 5

	defaultCreateSeedData := func(t *testing.T, store db.Store) test_util.SeedData {
		var err error

		ctx := context.Background()

		users := make([]db.User, 2)
		for i := range users {
			users[i] = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
				Name:         fmt.Sprintf("testuser-%d", i),
				Email:        fmt.Sprintf("test-%d@example.com", i),
				Password:     "test-password",
				SessionToken: token.NewToken(time.Minute),
			})
		}

		categories := make([]db.Category, 3)
		for i := range categories {
			categories[i], err = store.CreateCategory(ctx, fmt.Sprintf("test-category-%d", i))
			require.NoError(t, err)
		}

		for i := 0; i < 6; i++ {
			_, err = store.CreateProduct(ctx, db.CreateProductParams{
				Name:          fmt.Sprintf("test-product-%d", i),
				Description:   sql.NullString{String: fmt.Sprintf("test-description-%d", i), Valid: true},
				Price:         fmt.Sprintf("%d.00", (i+1)*10),
				StockQuantity: int32(i + 1),
				CategoryID:    categories[i%3].ID,
				SellerID:      users[i%2].ID,
				ImageUrl:      sql.NullString{String: fmt.Sprintf("test-image-url-%d", i), Valid: true},
			})
			require.NoError(t, err)
		}

		return test_util.SeedData{
			"users":      users,
			"categories": categories,
		}
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store) test_util.SeedData
		createQuery    func(t *testing.T, seedData test_util.SeedData) test_util.Query
		checkResponse  func(t *testing.T, response *http.Response)
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "1",
					"page_size": fmt.Sprintf("%d", pageSize),
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
			name:           "FilterByCategory",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				categories, ok := seedData["categories"].([]db.Category)
				require.True(t, ok)

				return test_util.Query{
					"page_id":     "1",
					"page_size":   fmt.Sprintf("%d", pageSize),
					"category_id": categories[0].ID.String(),
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotResponse := unmarshalListProductsResponse(t, response.Body)

				require.Len(t, gotResponse.Data, 2)

				for i := 0; i < 2; i++ {
					require.Equal(t, "test-category-0", gotResponse.Data[i].Category)
				}
			},
		},
		{
			name:           "PageIDNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateAndReturnSeed,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageIDLessThanLowerLimit",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateAndReturnSeed,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "0",
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageSizeNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateAndReturnSeed,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id": "1",
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageSizeLessThanLowerLimit",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateAndReturnSeed,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "1",
					"page_size": "0",
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageSizeMoreThanUpperLimit",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateAndReturnSeed,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "1",
					"page_size": "101",
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalServerError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Return([]db.Product{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeedData: test_util.NoopCreateAndReturnSeed,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "1",
					"page_size": "1",
				}
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

			seedData := tc.createSeedData(t, store)

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodGet,
				URL:    "/api/v1/products",
				Query:  tc.createQuery(t, seedData),
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})

	}
}

func TestListProductsBySeller(t *testing.T) {
	pageSize := 5

	sessionTokens := test_util.NewSessionTokens(2, time.Minute)

	defaultCreateSeedData := func(t *testing.T, store db.Store) test_util.SeedData {
		var err error

		ctx := context.Background()

		users := make([]db.User, 2)
		for i := range users {
			users[i] = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
				Name:         fmt.Sprintf("testuser-%d", i),
				Email:        fmt.Sprintf("test-%d@example.com", i),
				Password:     "test-password",
				SessionToken: sessionTokens[i],
			})
		}

		categories := make([]db.Category, 3)
		for i := range categories {
			categories[i], err = store.CreateCategory(ctx, fmt.Sprintf("test-category-%d", i))
			require.NoError(t, err)
		}

		for i := 0; i < 6; i++ {
			_, err = store.CreateProduct(ctx, db.CreateProductParams{
				Name:          fmt.Sprintf("test-product-%d", i),
				Description:   sql.NullString{String: fmt.Sprintf("test-description-%d", i), Valid: true},
				Price:         fmt.Sprintf("%d.00", (i+1)*10),
				StockQuantity: int32(i + 1),
				CategoryID:    categories[i%3].ID,
				SellerID:      users[i%2].ID,
				ImageUrl:      sql.NullString{String: fmt.Sprintf("test-image-url-%d", i), Valid: true},
			})
			require.NoError(t, err)
		}

		return test_util.SeedData{
			"users":      users,
			"categories": categories,
		}
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store) test_util.SeedData
		createQuery    func(t *testing.T, seedData test_util.SeedData) test_util.Query
		setupAuth      func(request *http.Request, sessionToken string)
		checkResponse  func(t *testing.T, response *http.Response)
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "1",
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotResponse := unmarshalListProductsResponse(t, response.Body)

				require.Equal(t, int32(1), gotResponse.Meta.PageID)
				require.Equal(t, int32(pageSize), gotResponse.Meta.PageSize)
				require.Equal(t, int64(1), gotResponse.Meta.PageCount)
				require.Equal(t, int64(3), gotResponse.Meta.TotalCount)

				require.Len(t, gotResponse.Data, 3)

				for i := 0; i < 3; i++ {
					productIndex := i * 2
					categoryIndex := productIndex % 3

					require.NotEmpty(t, gotResponse.Data[i].ID)
					require.Equal(t, fmt.Sprintf("test-product-%d", productIndex), gotResponse.Data[i].Name)
					require.Equal(t, fmt.Sprintf("test-description-%d", productIndex), gotResponse.Data[i].Description.String)
					require.True(t, decimal.NewFromInt(int64((productIndex+1)*10)).Equal(gotResponse.Data[i].Price.Decimal))
					require.Equal(t, int32(productIndex+1), gotResponse.Data[i].StockQuantity)
					require.Equal(t, fmt.Sprintf("test-category-%d", categoryIndex), gotResponse.Data[i].Category)
					require.Equal(t, "testuser-0", gotResponse.Data[i].Seller)
					require.Equal(t, fmt.Sprintf("test-image-url-%d", productIndex), gotResponse.Data[i].ImageUrl.String)
				}
			},
		},
		{
			name:           "NoAuthorization",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "1",
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			setupAuth: test_util.NoopSetupAuth,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:           "PageIDNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageIDLessThanLowerLimit",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "0",
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageSizeNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id": "1",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageSizeLessThanLowerLimit",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "1",
					"page_size": "0",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageSizeMoreThanUpperLimit",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "1",
					"page_size": "101",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalServerError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				test_util.BuildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: sessionTokens[0].ID,
					ExpiredAt:    sessionTokens[0].ExpiredAt,
					CreatedAt:    time.Now(),
				})

				mockStore.EXPECT().
					ListProductsBySeller(gomock.Any(), gomock.Any()).
					Return([]db.Product{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeedData: test_util.NoopCreateAndReturnSeed,
			createQuery: func(t *testing.T, seedData test_util.SeedData) test_util.Query {
				return test_util.Query{
					"page_id":   "1",
					"page_size": "1",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
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

			seedData := tc.createSeedData(t, store)

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodGet,
				URL:    "/api/v1/users/products",
				Query:  tc.createQuery(t, seedData),
			})

			tc.setupAuth(request, sessionTokens[0].ID.String())

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})

	}
}

func TestListProductCategories(t *testing.T) {
	pageSize := 5

	defaultCreateSeedData := func(t *testing.T, store db.Store) {
		var err error

		ctx := context.Background()

		categories := make([]db.Category, 6)
		for i := range categories {
			categories[i], err = store.CreateCategory(ctx, fmt.Sprintf("test-category-%d", i))
			require.NoError(t, err)
		}
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store)
		query          test_util.Query
		checkResponse  func(t *testing.T, response *http.Response)
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			query: test_util.Query{
				"page_id":   "1",
				"page_size": fmt.Sprintf("%d", pageSize),
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotResponse := unmarshalListProductCategoriesResponse(t, response.Body)

				require.Equal(t, int32(1), gotResponse.Meta.PageID)
				require.Equal(t, int32(pageSize), gotResponse.Meta.PageSize)

				require.Len(t, gotResponse.Data, pageSize)

				for i := 0; i < pageSize; i++ {
					require.NotEmpty(t, gotResponse.Data[i].ID)
					require.Equal(t, fmt.Sprintf("test-category-%d", i), gotResponse.Data[i].Name)
				}
			},
		},
		{
			name:           "PageIDNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateSeedData,
			query: test_util.Query{
				"page_size": fmt.Sprintf("%d", pageSize),
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageIDLessThanLowerLimit",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateSeedData,
			query: test_util.Query{
				"page_id":   "0",
				"page_size": fmt.Sprintf("%d", pageSize),
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageSizeNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateSeedData,
			query: test_util.Query{
				"page_id": "1",
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageSizeLessThanLowerLimit",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateSeedData,
			query: test_util.Query{
				"page_id":   "1",
				"page_size": "0",
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PageSizeMoreThanUpperLimit",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateSeedData,
			query: test_util.Query{
				"page_id":   "1",
				"page_size": "101",
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalServerError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					ListCategories(gomock.Any(), gomock.Any()).
					Return([]db.Category{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeedData: test_util.NoopCreateSeedData,
			query: test_util.Query{
				"page_id":   "1",
				"page_size": "1",
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

			tc.createSeedData(t, store)

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodGet,
				URL:    "/api/v1/products/categories",
				Query:  tc.query,
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})

	}
}

func TestProductHandlerAddProduct(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	defaultCreateSeedData := func(t *testing.T, store db.Store) test_util.SeedData {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
		})

		category, err := store.CreateCategory(ctx, "test-category")
		require.NoError(t, err)

		return test_util.SeedData{
			"category": category,
		}
	}

	defaultCreateBody := func(seedData test_util.SeedData) test_util.Body {
		return test_util.Body{
			"name":           "test-product",
			"description":    "test-description",
			"price":          "10.00",
			"stock_quantity": 10,
			"category_id":    seedData["category"].(db.Category).ID,
			"image_url":      "https://example.com/image.png",
		}
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store) test_util.SeedData
		createBody     func(seedData test_util.SeedData) test_util.Body
		setupAuth      func(request *http.Request, sessionToken string)
		checkResponse  func(t *testing.T, response *http.Response)
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody:     defaultCreateBody,
			setupAuth:      test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:           "OptionalFieldsNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody: func(seedData test_util.SeedData) test_util.Body {
				return test_util.Body{
					"name":           "test-product",
					"price":          "10.00",
					"stock_quantity": 10,
					"category_id":    seedData["category"].(db.Category).ID,
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:           "NoAuthorization",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody:     defaultCreateBody,
			setupAuth:      test_util.NoopSetupAuth,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:           "NameNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody: func(seedData test_util.SeedData) test_util.Body {
				return test_util.Body{
					"description":    "test-description",
					"price":          "10.00",
					"stock_quantity": 10,
					"category_id":    seedData["category"].(db.Category).ID,
					"image_url":      "https://example.com/image.png",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PriceNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody: func(seedData test_util.SeedData) test_util.Body {
				return test_util.Body{
					"name":           "test-product",
					"description":    "test-description",
					"stock_quantity": 10,
					"category_id":    seedData["category"].(db.Category).ID,
					"image_url":      "https://example.com/image.png",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PriceInvalidFormat",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody: func(seedData test_util.SeedData) test_util.Body {
				return test_util.Body{
					"name":           "test-product",
					"description":    "test-description",
					"price":          "Invalid Price",
					"stock_quantity": 10,
					"category_id":    seedData["category"].(db.Category).ID,
					"image_url":      "https://example.com/image.png",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "PriceIsZero",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody: func(seedData test_util.SeedData) test_util.Body {
				return test_util.Body{
					"name":           "test-product",
					"description":    "test-description",
					"price":          "0.00",
					"stock_quantity": 10,
					"category_id":    seedData["category"].(db.Category).ID,
					"image_url":      "https://example.com/image.png",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "ImageUrlInvalidFormat",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody: func(seedData test_util.SeedData) test_util.Body {
				return test_util.Body{
					"name":           "test-product",
					"description":    "test-description",
					"price":          "10.00",
					"stock_quantity": 10,
					"category_id":    seedData["category"].(db.Category).ID,
					"image_url":      "ttp://example.com/image.png",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "StockQuantityNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody: func(seedData test_util.SeedData) test_util.Body {
				return test_util.Body{
					"name":        "test-product",
					"description": "test-description",
					"price":       "10.00",
					"category_id": seedData["category"].(db.Category).ID,
					"image_url":   "https://example.com/image.png",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "CategoryIDNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			createBody: func(seedData test_util.SeedData) test_util.Body {
				return test_util.Body{
					"name":           "test-product",
					"description":    "test-description",
					"price":          "10.00",
					"stock_quantity": 10,
					"image_url":      "https://example.com/image.png",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				test_util.BuildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: sessionToken.ID,
					ExpiredAt:    sessionToken.ExpiredAt,
					CreatedAt:    time.Now(),
				})

				mockStore.EXPECT().
					AddProduct(gomock.Any(), gomock.Any()).
					Return(db.Product{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeedData: test_util.NoopCreateAndReturnSeed,
			createBody: func(seedData test_util.SeedData) test_util.Body {
				return test_util.Body{
					"name":           "test-product",
					"description":    "test-description",
					"price":          "10.00",
					"stock_quantity": 10,
					"category_id":    util.RandomUUID().String(),
					"seller_id":      util.RandomUUID().String(),
					"image_url":      "https://example.com/image.png",
				}
			},
			setupAuth: test_util.AddSessionTokenInCookie,
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

			seedData := tc.createSeedData(t, store)

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodPost,
				URL:    "/api/v1/users/products",
				Body:   tc.createBody(seedData),
			})

			tc.setupAuth(request, sessionToken.ID.String())

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

func unmarshalListProductCategoriesResponse(t *testing.T, body io.ReadCloser) product_domain.ListProductCategoriesResponse {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var parsed product_domain.ListProductCategoriesResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	return parsed
}
