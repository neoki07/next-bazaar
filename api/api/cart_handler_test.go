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
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
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
			buildStore: buildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  func(request *http.Request) {},
			buildStore: buildTestDBStore,
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
				mockStore, cleanup := newMockStore(t)

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

			url := "/api/v1/cart-products"
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
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
				mockStore, cleanup := newMockStore(t)

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

			url := "/api/v1/cart-products/count"
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
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
				mockStore, cleanup := newMockStore(t)

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

			url := "/api/v1/cart-products"

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
			buildStore: buildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			body:       defaultBody,
			setupAuth:  func(request *http.Request) {},
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
				mockStore, cleanup := newMockStore(t)

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

			server := newTestServer(t, store)

			url := fmt.Sprintf("/api/v1/cart-products/%s", productID)

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(request)
			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

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
			buildStore: buildTestDBStore,
			createSeed: defaultCreateSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNoContent, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  func(request *http.Request) {},
			buildStore: buildTestDBStore,
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
			buildStore: buildTestDBStore,
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
				mockStore, cleanup := newMockStore(t)

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

			url := fmt.Sprintf("/api/v1/cart-products/%s", productID)

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

func unmarshalCartProductsCountResponse(t *testing.T, body io.ReadCloser) cart_domain.CartProductsCountResponse {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var parsed cart_domain.CartProductsCountResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	return parsed
}
