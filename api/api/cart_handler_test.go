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
	cart_domain "github.com/ot07/next-bazaar/api/domain/cart"
	"github.com/ot07/next-bazaar/api/test_util"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestGetCart(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		user := test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
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

		_, err = store.CreateCartProduct(ctx, db.CreateCartProductParams{
			UserID:    user.ID,
			ProductID: product.ID,
			Quantity:  5,
		})
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request, sessionToken string)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotResponse := unmarshalCartResponse(t, response.Body)

				require.Equal(t, 1, len(gotResponse.Products))

				require.Equal(t, "test-product", gotResponse.Products[0].Name)
				require.True(t, decimal.NewFromFloat(100.00).Equal(gotResponse.Products[0].Price.Decimal))
				require.Equal(t, int32(5), gotResponse.Products[0].Quantity)
				require.True(t, decimal.NewFromFloat(500.00).Equal(gotResponse.Products[0].Subtotal.Decimal))
				require.Equal(t, "test-image-url", gotResponse.Products[0].ImageUrl.NullString.String)

				require.True(t, decimal.NewFromFloat(500.00).Equal(gotResponse.Subtotal.Decimal))
				require.True(t, decimal.NewFromFloat(5.00).Equal(gotResponse.Shipping.Decimal))
				require.True(t, decimal.NewFromFloat(50.00).Equal(gotResponse.Tax.Decimal))
				require.True(t, decimal.NewFromFloat(555.00).Equal(gotResponse.Total.Decimal))
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  test_util.NoopSetupAuth,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:      "InternalError",
			setupAuth: test_util.AddSessionTokenInCookie,
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
					GetCartProductsByUserID(gomock.Any(), gomock.Any()).
					Return([]db.CartProduct{}, sql.ErrConnDone)

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
				URL:    "/api/v1/cart",
			})

			tc.setupAuth(request, sessionToken.ID.String())

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestGetCartProductsCount(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		user := test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
		})

		category, err := store.CreateCategory(ctx, "test-category")
		require.NoError(t, err)

		products := make([]db.Product, 2)
		for i := range products {
			products[i], err = store.CreateProduct(ctx, db.CreateProductParams{
				Name:          "test-product",
				Description:   sql.NullString{String: "test-description", Valid: true},
				Price:         "100.00",
				StockQuantity: 10,
				CategoryID:    category.ID,
				SellerID:      user.ID,
				ImageUrl:      sql.NullString{String: "test-image-url", Valid: true},
			})
			require.NoError(t, err)
		}

		_, err = store.CreateCartProduct(ctx, db.CreateCartProductParams{
			UserID:    user.ID,
			ProductID: products[0].ID,
			Quantity:  3,
		})
		require.NoError(t, err)

		_, err = store.CreateCartProduct(ctx, db.CreateCartProductParams{
			UserID:    user.ID,
			ProductID: products[1].ID,
			Quantity:  5,
		})
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request, sessionToken string)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotResponse := unmarshalCartProductsCountResponse(t, response.Body)

				require.Equal(t, int32(8), gotResponse.Count)
			},
		},
		{
			name:       "Empty",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) {
				ctx := context.Background()

				_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
					Name:         "testuser",
					Email:        "test@example.com",
					Password:     "test-password",
					SessionToken: sessionToken,
				})
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotResponse := unmarshalCartProductsCountResponse(t, response.Body)

				require.Equal(t, int32(0), gotResponse.Count)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  test_util.NoopSetupAuth,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:      "InternalError",
			setupAuth: test_util.AddSessionTokenInCookie,
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
					GetCartProductsByUserID(gomock.Any(), gomock.Any()).
					Return([]db.CartProduct{}, sql.ErrConnDone)

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
				URL:    "/api/v1/cart/count",
			})

			tc.setupAuth(request, sessionToken.ID.String())

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestAddProduct(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) (productID string) {
		ctx := context.Background()

		user := test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
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
	}

	defaultCreateBody := func(t *testing.T, productID string) fiber.Map {
		return fiber.Map{
			"product_id": productID,
			"quantity":   1,
		}
	}

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request, sessionToken string)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) (productID string)
		createBody    func(t *testing.T, productID string) fiber.Map
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			createBody: defaultCreateBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "AddExistingProduct",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) (productID string) {
				ctx := context.Background()

				user := test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
					Name:         "testuser",
					Email:        "test@example.com",
					Password:     "test-password",
					SessionToken: sessionToken,
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

				_, err = store.CreateCartProduct(ctx, db.CreateCartProductParams{
					UserID:    user.ID,
					ProductID: product.ID,
					Quantity:  5,
				})
				require.NoError(t, err)

				return product.ID.String()
			},
			createBody: defaultCreateBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  test_util.NoopSetupAuth,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			createBody: func(t *testing.T, productID string) fiber.Map {
				return fiber.Map{
					"product_id": productID,
					"quantity":   1,
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:       "ProductIDNotFound",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			createBody: func(t *testing.T, productID string) fiber.Map {
				return fiber.Map{
					"quantity": 1,
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:       "QuantityNotFound",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			createBody: func(t *testing.T, productID string) fiber.Map {
				return fiber.Map{
					"product_id": productID,
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:       "QuantityIsZero",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			createBody: func(t *testing.T, productID string) fiber.Map {
				return fiber.Map{
					"product_id": productID,
					"quantity":   0,
				}
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:      "InternalError",
			setupAuth: test_util.AddSessionTokenInCookie,
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
					GetCartProductByUserIDAndProductID(gomock.Any(), gomock.Any()).
					Return(db.CartProduct{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: func(t *testing.T, store db.Store) (productID string) {
				return util.RandomUUID().String()
			},
			createBody: defaultCreateBody,
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
				Method: http.MethodPost,
				URL:    "/api/v1/cart/add-product",
				Body:   tc.createBody(t, productID),
			})

			tc.setupAuth(request, sessionToken.ID.String())

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestUpdateProductQuantity(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) (productID string) {
		ctx := context.Background()

		user := test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
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

		_, err = store.CreateCartProduct(ctx, db.CreateCartProductParams{
			UserID:    user.ID,
			ProductID: product.ID,
			Quantity:  5,
		})
		require.NoError(t, err)

		return product.ID.String()
	}

	defaultBody := fiber.Map{
		"quantity": 1,
	}

	testCases := []struct {
		name          string
		body          fiber.Map
		setupAuth     func(request *http.Request, sessionToken string)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) (productID string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			body:       defaultBody,
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			body:       defaultBody,
			setupAuth:  test_util.NoopSetupAuth,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:       "ProductNotFound",
			body:       defaultBody,
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) (productID string) {
				_ = defaultCreateSeed(t, store)
				return util.RandomUUID().String()
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name:       "QuantityNotFound",
			body:       fiber.Map{},
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "QuantityIsZero",
			body: fiber.Map{
				"quantity": 0,
			},
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:      "InternalError",
			body:      defaultBody,
			setupAuth: test_util.AddSessionTokenInCookie,
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
					UpdateCartProduct(gomock.Any(), gomock.Any()).
					Return(db.CartProduct{}, sql.ErrConnDone)

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
				Method: http.MethodPut,
				URL:    fmt.Sprintf("/api/v1/cart/%s", productID),
				Body:   tc.body,
			})

			tc.setupAuth(request, sessionToken.ID.String())

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) (productID string) {
		ctx := context.Background()

		user := test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
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

		_, err = store.CreateCartProduct(ctx, db.CreateCartProductParams{
			UserID:    user.ID,
			ProductID: product.ID,
			Quantity:  5,
		})
		require.NoError(t, err)

		return product.ID.String()
	}

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request, sessionToken string)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) (productID string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNoContent, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  test_util.NoopSetupAuth,
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:       "ProductNotFound",
			setupAuth:  test_util.AddSessionTokenInCookie,
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) (productID string) {
				_ = defaultCreateSeed(t, store)
				return util.RandomUUID().String()
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNoContent, response.StatusCode)
			},
		},
		{
			name:      "InternalError",
			setupAuth: test_util.AddSessionTokenInCookie,
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
					DeleteCartProduct(gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone)

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
				Method: http.MethodDelete,
				URL:    fmt.Sprintf("/api/v1/cart/%s", productID),
			})

			tc.setupAuth(request, sessionToken.ID.String())

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestCartAPIScenario(t *testing.T) {
	ctx := context.Background()

	store, cleanupStore := test_util.BuildTestDBStore(t)
	defer cleanupStore()

	server := newTestServer(t, store)

	// Create user
	sessionToken := token.NewToken(time.Minute)
	user := test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
		Name:         "testuser",
		Email:        "test@example.com",
		Password:     "test-password",
		SessionToken: sessionToken,
	})

	setupAuth := func(request *http.Request) {
		test_util.AddSessionTokenInCookie(request, sessionToken.ID.String())
	}

	// Create get cart function
	getCart := func() *http.Response {
		request := test_util.NewRequest(t, test_util.RequestParams{
			Method: http.MethodGet,
			URL:    "/api/v1/cart",
		})
		setupAuth(request)
		response := test_util.SendRequest(t, server.app, request)
		require.Equal(t, http.StatusOK, response.StatusCode)

		return response
	}

	// Create product
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

	// Add product to cart
	request := test_util.NewRequest(t, test_util.RequestParams{
		Method: http.MethodPost,
		URL:    "/api/v1/cart/add-product",
		Body: fiber.Map{
			"product_id": product.ID,
			"quantity":   5,
		},
	})
	setupAuth(request)
	response := test_util.SendRequest(t, server.app, request)
	require.Equal(t, http.StatusOK, response.StatusCode)

	// Get cart
	response = getCart()
	require.Equal(t, http.StatusOK, response.StatusCode)

	gotResponse := unmarshalCartResponse(t, response.Body)
	require.Equal(t, int32(5), gotResponse.Products[0].Quantity)

	// Update product quantity
	request = test_util.NewRequest(t, test_util.RequestParams{
		Method: http.MethodPut,
		URL:    fmt.Sprintf("/api/v1/cart/%s", product.ID),
		Body: fiber.Map{
			"quantity": 10,
		},
	})
	setupAuth(request)
	response = test_util.SendRequest(t, server.app, request)
	require.Equal(t, http.StatusOK, response.StatusCode)

	// Get cart
	response = getCart()
	require.Equal(t, http.StatusOK, response.StatusCode)

	gotResponse = unmarshalCartResponse(t, response.Body)
	require.Equal(t, int32(10), gotResponse.Products[0].Quantity)

	// Delete product
	request = test_util.NewRequest(t, test_util.RequestParams{
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("/api/v1/cart/%s", product.ID),
	})
	setupAuth(request)
	response = test_util.SendRequest(t, server.app, request)
	require.Equal(t, http.StatusNoContent, response.StatusCode)

	// Get cart
	response = getCart()
	require.Equal(t, http.StatusOK, response.StatusCode)

	gotResponse = unmarshalCartResponse(t, response.Body)
	require.Equal(t, 0, len(gotResponse.Products))
}

func unmarshalCartResponse(t *testing.T, body io.ReadCloser) cart_domain.CartResponse {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var parsed cart_domain.CartResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	return parsed
}

func unmarshalCartProductsCountResponse(t *testing.T, body io.ReadCloser) cart_domain.CartProductsCountResponse {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var parsed cart_domain.CartProductsCountResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	return parsed
}
