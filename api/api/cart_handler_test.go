package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func TestGetCart(t *testing.T) {
	t.Parallel()

	validUserName := "testuser"
	validUserEmail := "test@example.com"
	validUserPassword := "test-password"

	validUserHashedPassword, err := util.HashPassword(validUserPassword)
	require.NoError(t, err)

	validSessionToken := token.NewToken(time.Minute)

	createSeed := func(t *testing.T, store db.Store) (userID string) {
		ctx := context.Background()

		createdUser, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validUserName,
			Email:          validUserEmail,
			HashedPassword: validUserHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       createdUser.ID,
			SessionToken: validSessionToken.ID,
			ExpiredAt:    validSessionToken.ExpiredAt,
		})
		require.NoError(t, err)

		createdCategory, err := store.CreateCategory(ctx, "test-category")
		require.NoError(t, err)

		createdProduct, err := store.CreateProduct(ctx, db.CreateProductParams{
			Name:          "test-product",
			Description:   sql.NullString{String: "test-description", Valid: true},
			Price:         "100.00",
			StockQuantity: 10,
			CategoryID:    createdCategory.ID,
			SellerID:      createdUser.ID,
			ImageUrl:      sql.NullString{String: "test-image-url", Valid: true},
		})
		require.NoError(t, err)

		createdCartProduct, err := store.CreateCartProduct(ctx, db.CreateCartProductParams{
			UserID:    createdUser.ID,
			ProductID: createdProduct.ID,
			Quantity:  5,
		})
		require.NoError(t, err)

		return createdCartProduct.UserID.String()
	}

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store) (userID string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(request *http.Request) {},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InvalidUserID",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: func(t *testing.T, store db.Store) (userID string) {
				_ = createSeed(t, store)
				return "InvalidUserID"
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
					GetCartProductsByUserId(gomock.Any(), gomock.Any()).
					Return([]db.CartProduct{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: func(t *testing.T, store db.Store) (userID string) {
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

			userID := tc.createSeed(t, store)

			server := newTestServer(t, store)

			url := fmt.Sprintf("/api/v1/cart-products/%s", userID)
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
