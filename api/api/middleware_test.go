package api

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/golang/mock/gomock"
	"github.com/ot07/next-bazaar/api/test_util"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/stretchr/testify/require"
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
		name          string
		setupAuth     func(request *http.Request)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createValidSessionSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  func(request *http.Request) {},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createValidSessionSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InvalidTokenFormat",
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(request, "invalid")
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createValidSessionSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(request, expiredSessionToken.ID.String())
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createExpiredSessionSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(request, validSessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Return(db.Session{}, sql.ErrConnDone)

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
