package api

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ot07/next-bazaar/api/test_util"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestAuthMiddleware(t *testing.T) {
	validSessionToken := token.NewToken(time.Minute)
	createValidSessionSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: validSessionToken,
		})
	}

	expiredSessionToken := token.NewToken(-time.Minute)
	createExpiredSessionSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: expiredSessionToken,
		})
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store)
		setupAuth      func(request *http.Request)
		checkResponse  func(t *testing.T, response *http.Response)
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: createValidSessionSeed,
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:           "NoAuthorization",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: createValidSessionSeed,
			setupAuth:      func(request *http.Request) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:           "InvalidTokenFormat",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: createValidSessionSeed,
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(request, "invalid")
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:           "ExpiredToken",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: createExpiredSessionSeed,
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(request, expiredSessionToken.ID.String())
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Return(db.Session{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeedData: func(t *testing.T, store db.Store) {},
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(request, validSessionToken.ID.String())
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

			authPath := "/auth"

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodGet,
				URL:    authPath,
			})

			tc.setupAuth(request)

			server := newTestServer(t, store)
			server.app.Get(
				authPath,
				authMiddleware(server),
				func(c *fiber.Ctx) error {
					return c.SendStatus(fiber.StatusOK)
				},
			)

			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}
