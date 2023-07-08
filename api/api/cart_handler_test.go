package api

import (
	"bytes"
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
	validUserName := "testuser"
	validUserEmail := "test@example.com"
	validUserPassword := "test-password"

	validUserHashedPassword, err := util.HashPassword(validUserPassword)
	require.NoError(t, err)

	validSessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		user, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validUserName,
			Email:          validUserEmail,
			HashedPassword: validUserHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       user.ID,
			SessionToken: validSessionToken.ID,
			ExpiredAt:    validSessionToken.ExpiredAt,
		})
		require.NoError(t, err)

		category, err := store.CreateCategory(ctx, "test-category")
		require.NoError(t, err)

		createdProduct, err := store.CreateProduct(ctx, db.CreateProductParams{
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
			ProductID: createdProduct.ID,
			Quantity:  5,
		})
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
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
			setupAuth:  func(request *http.Request) {},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				buildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: validSessionToken.ID,
					ExpiredAt:    validSessionToken.ExpiredAt,
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

			server := newTestServer(t, store)

			url := "/api/v1/cart"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(request)
			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestGetCartProductsCount(t *testing.T) {
	validUserName := "testuser"
	validUserEmail := "test@example.com"
	validUserPassword := "test-password"

	validUserHashedPassword, err := util.HashPassword(validUserPassword)
	require.NoError(t, err)

	validSessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		user, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validUserName,
			Email:          validUserEmail,
			HashedPassword: validUserHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       user.ID,
			SessionToken: validSessionToken.ID,
			ExpiredAt:    validSessionToken.ExpiredAt,
		})
		require.NoError(t, err)

		category, err := store.CreateCategory(ctx, "test-category")
		require.NoError(t, err)

		createdProducts := make([]db.Product, 2)
		for i := range createdProducts {
			createdProducts[i], err = store.CreateProduct(ctx, db.CreateProductParams{
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
			ProductID: createdProducts[0].ID,
			Quantity:  3,
		})
		require.NoError(t, err)

		_, err = store.CreateCartProduct(ctx, db.CreateCartProductParams{
			UserID:    user.ID,
			ProductID: createdProducts[1].ID,
			Quantity:  5,
		})
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotResponse := unmarshalCartProductsCountResponse(t, response.Body)

				require.Equal(t, int32(8), gotResponse.Count)
			},
		},
		{
			name: "Empty",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) {
				ctx := context.Background()

				user, err := store.CreateUser(ctx, db.CreateUserParams{
					Name:           validUserName,
					Email:          validUserEmail,
					HashedPassword: validUserHashedPassword,
				})
				require.NoError(t, err)

				_, err = store.CreateSession(ctx, db.CreateSessionParams{
					UserID:       user.ID,
					SessionToken: validSessionToken.ID,
					ExpiredAt:    validSessionToken.ExpiredAt,
				})
				require.NoError(t, err)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotResponse := unmarshalCartProductsCountResponse(t, response.Body)

				require.Equal(t, int32(0), gotResponse.Count)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  func(request *http.Request) {},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				buildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: validSessionToken.ID,
					ExpiredAt:    validSessionToken.ExpiredAt,
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

			server := newTestServer(t, store)

			url := "/api/v1/cart/count"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(request)
			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestAddProduct(t *testing.T) {
	validUserName := "testuser"
	validUserEmail := "test@example.com"
	validUserPassword := "test-password"

	validUserHashedPassword, err := util.HashPassword(validUserPassword)
	require.NoError(t, err)

	validSessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) (productID string) {
		ctx := context.Background()

		user, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validUserName,
			Email:          validUserEmail,
			HashedPassword: validUserHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       user.ID,
			SessionToken: validSessionToken.ID,
			ExpiredAt:    validSessionToken.ExpiredAt,
		})
		require.NoError(t, err)

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
		setupAuth     func(request *http.Request)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) (productID string)
		createBody    func(t *testing.T, productID string) fiber.Map
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			createBody: defaultCreateBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name: "AddExistingProduct",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) (productID string) {
				ctx := context.Background()

				user, err := store.CreateUser(ctx, db.CreateUserParams{
					Name:           validUserName,
					Email:          validUserEmail,
					HashedPassword: validUserHashedPassword,
				})
				require.NoError(t, err)

				_, err = store.CreateSession(ctx, db.CreateSessionParams{
					UserID:       user.ID,
					SessionToken: validSessionToken.ID,
					ExpiredAt:    validSessionToken.ExpiredAt,
				})
				require.NoError(t, err)

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
			setupAuth:  func(request *http.Request) {},
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
			name: "ProductIDNotFound",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
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
			name: "QuantityNotFound",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
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
			name: "QuantityIsZero",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
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
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				buildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: validSessionToken.ID,
					ExpiredAt:    validSessionToken.ExpiredAt,
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

			server := newTestServer(t, store)

			url := "/api/v1/cart/add-product"

			body, err := json.Marshal(tc.createBody(t, productID))
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(request)
			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestUpdateProductQuantity(t *testing.T) {
	validUserName := "testuser"
	validUserEmail := "test@example.com"
	validUserPassword := "test-password"

	validUserHashedPassword, err := util.HashPassword(validUserPassword)
	require.NoError(t, err)

	validSessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) (productID string) {
		ctx := context.Background()

		user, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validUserName,
			Email:          validUserEmail,
			HashedPassword: validUserHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       user.ID,
			SessionToken: validSessionToken.ID,
			ExpiredAt:    validSessionToken.ExpiredAt,
		})
		require.NoError(t, err)

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
		setupAuth     func(request *http.Request)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) (productID string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			body: defaultBody,
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			body:       defaultBody,
			setupAuth:  func(request *http.Request) {},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "ProductNotFound",
			body: defaultBody,
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
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
			name: "QuantityNotFound",
			body: fiber.Map{},
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
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
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			body: defaultBody,
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				buildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: validSessionToken.ID,
					ExpiredAt:    validSessionToken.ExpiredAt,
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

			request := test_util.NewRequest(
				t,
				test_util.RequestParams{
					Method:    http.MethodPut,
					URL:       fmt.Sprintf("/api/v1/cart/%s", productID),
					Body:      tc.body,
					SetupAuth: tc.setupAuth,
				},
			)

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	validUserName := "testuser"
	validUserEmail := "test@example.com"
	validUserPassword := "test-password"

	validUserHashedPassword, err := util.HashPassword(validUserPassword)
	require.NoError(t, err)

	validSessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) (productID string) {
		ctx := context.Background()

		user, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validUserName,
			Email:          validUserEmail,
			HashedPassword: validUserHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       user.ID,
			SessionToken: validSessionToken.ID,
			ExpiredAt:    validSessionToken.ExpiredAt,
		})
		require.NoError(t, err)

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
		setupAuth     func(request *http.Request)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) (productID string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNoContent, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  func(request *http.Request) {},
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "ProductNotFound",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
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
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				buildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: validSessionToken.ID,
					ExpiredAt:    validSessionToken.ExpiredAt,
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

			server := newTestServer(t, store)

			url := fmt.Sprintf("/api/v1/cart/%s", productID)

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(request)
			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

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
	user, sessionToken := test_util.CreateUserTestData(
		t,
		ctx,
		store,
		"testuser",
		"test@example.com",
		"test-password",
	)

	// Create get cart function
	getCart := func() *http.Response {
		url := "/api/v1/cart"
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		addSessionTokenInCookie(request, sessionToken.ID.String())
		response, err := server.app.Test(request, int(time.Second.Milliseconds()))
		require.NoError(t, err)

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
	url := "/api/v1/cart/add-product"

	body, err := json.Marshal(fiber.Map{
		"product_id": product.ID,
		"quantity":   5,
	})
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	addSessionTokenInCookie(request, sessionToken.ID.String())
	response, err := server.app.Test(request, int(time.Second.Milliseconds()))
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	// Get cart
	response = getCart()

	require.Equal(t, http.StatusOK, response.StatusCode)

	gotResponse := unmarshalCartResponse(t, response.Body)

	require.Equal(t, int32(5), gotResponse.Products[0].Quantity)

	// Update product quantity
	url = fmt.Sprintf("/api/v1/cart/%s", product.ID)

	body, err = json.Marshal(fiber.Map{
		"quantity": 10,
	})
	require.NoError(t, err)

	request, err = http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	addSessionTokenInCookie(request, sessionToken.ID.String())
	response, err = server.app.Test(request, int(time.Second.Milliseconds()))
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	// Get cart
	response = getCart()

	require.Equal(t, http.StatusOK, response.StatusCode)

	gotResponse = unmarshalCartResponse(t, response.Body)

	require.Equal(t, int32(10), gotResponse.Products[0].Quantity)

	// Delete product
	url = fmt.Sprintf("/api/v1/cart/%s", product.ID)

	request, err = http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	addSessionTokenInCookie(request, sessionToken.ID.String())
	response, err = server.app.Test(request, int(time.Second.Milliseconds()))
	require.NoError(t, err)

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
