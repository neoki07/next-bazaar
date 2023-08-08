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

				return product.ID.String()
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

	defaultCreateSeed := func(t *testing.T, store db.Store) fiber.Map {
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

		return fiber.Map{
			"users":      users,
			"categories": categories,
		}
	}

	noopCreateSeed := func(t *testing.T, store db.Store) fiber.Map {
		return fiber.Map{}
	}

	testCases := []struct {
		name          string
		createQuery   func(t *testing.T, seedData fiber.Map) fiber.Map
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) fiber.Map
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "1",
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
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
			name: "FilterByCategory",
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				categories, ok := seedData["categories"].([]db.Category)
				require.True(t, ok)

				return fiber.Map{
					"page_id":     "1",
					"page_size":   fmt.Sprintf("%d", pageSize),
					"category_id": categories[0].ID.String(),
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
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
			name: "PageIDNotFound",
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageIDLessThanLowerLimit",
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "0",
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeNotFound",
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id": "1",
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeLessThanLowerLimit",
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "1",
					"page_size": "0",
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeMoreThanUpperLimit",
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "1",
					"page_size": "101",
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalServerError",
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "1",
					"page_size": "1",
				}
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Return([]db.Product{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: noopCreateSeed,
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

			seedData := tc.createSeed(t, store)

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

	defaultCreateSeed := func(t *testing.T, store db.Store) fiber.Map {
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

		return fiber.Map{
			"users":      users,
			"categories": categories,
		}
	}

	noopCreateSeed := func(t *testing.T, store db.Store) fiber.Map {
		return fiber.Map{}
	}

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request, sessionToken string)
		createQuery   func(t *testing.T, seedData fiber.Map) fiber.Map
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) fiber.Map
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:      "OK",
			setupAuth: test_util.AddSessionTokenInCookie,
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "1",
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
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
			name:      "NoAuthorization",
			setupAuth: test_util.NoopSetupAuth,
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "1",
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:      "PageIDNotFound",
			setupAuth: test_util.AddSessionTokenInCookie,
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:      "PageIDLessThanLowerLimit",
			setupAuth: test_util.AddSessionTokenInCookie,
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "0",
					"page_size": fmt.Sprintf("%d", pageSize),
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:      "PageSizeNotFound",
			setupAuth: test_util.AddSessionTokenInCookie,
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id": "1",
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:      "PageSizeLessThanLowerLimit",
			setupAuth: test_util.AddSessionTokenInCookie,
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "1",
					"page_size": "0",
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:      "PageSizeMoreThanUpperLimit",
			setupAuth: test_util.AddSessionTokenInCookie,
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "1",
					"page_size": "101",
				}
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:      "InternalServerError",
			setupAuth: test_util.AddSessionTokenInCookie,
			createQuery: func(t *testing.T, seedData fiber.Map) fiber.Map {
				return fiber.Map{
					"page_id":   "1",
					"page_size": "1",
				}
			},
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
			createSeed: noopCreateSeed,
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

			seedData := tc.createSeed(t, store)

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

	defaultCreateSeed := func(t *testing.T, store db.Store) {
		var err error

		ctx := context.Background()

		categories := make([]db.Category, 6)
		for i := range categories {
			categories[i], err = store.CreateCategory(ctx, fmt.Sprintf("test-category-%d", i))
			require.NoError(t, err)
		}
	}

	noopCreateSeed := func(t *testing.T, store db.Store) {}

	testCases := []struct {
		name          string
		query         fiber.Map
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			query: fiber.Map{
				"page_id":   "1",
				"page_size": fmt.Sprintf("%d", pageSize),
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
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
			name: "PageIDNotFound",
			query: fiber.Map{
				"page_size": fmt.Sprintf("%d", pageSize),
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageIDLessThanLowerLimit",
			query: fiber.Map{
				"page_id":   "0",
				"page_size": fmt.Sprintf("%d", pageSize),
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeNotFound",
			query: fiber.Map{
				"page_id": "1",
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeLessThanLowerLimit",
			query: fiber.Map{
				"page_id":   "1",
				"page_size": "0",
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeMoreThanUpperLimit",
			query: fiber.Map{
				"page_id":   "1",
				"page_size": "101",
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: noopCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalServerError",
			query: fiber.Map{
				"page_id":   "1",
				"page_size": "1",
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					ListCategories(gomock.Any(), gomock.Any()).
					Return([]db.Category{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: noopCreateSeed,
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
				URL:    "/api/v1/products/categories",
				Query:  tc.query,
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

func unmarshalListProductCategoriesResponse(t *testing.T, body io.ReadCloser) product_domain.ListProductCategoriesResponse {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var parsed product_domain.ListProductCategoriesResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	return parsed
}
