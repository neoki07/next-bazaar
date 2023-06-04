package api

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	mockdb "github.com/ot07/next-bazaar/db/mock"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/stretchr/testify/require"
)

func addSessionTokenInCookie(
	request *http.Request,
	sessionToken string,
) {
	cookie := &http.Cookie{
		Name:     cookieSessionTokenKey,
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	request.AddCookie(cookie)
}

func buildValidSessionStubs(store *mockdb.MockStore, session db.Session) {
	store.EXPECT().
		GetSession(gomock.Any(), gomock.Any()).
		Return(session, nil)
}

func TestAuthMiddleware(t *testing.T) {
	validName := "testuser"
	validEmail := "test@example.com"
	validHashedPassword := "test-hashed-password"

	validSessionToken := token.NewToken(time.Minute)
	createValidSessionSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		createdUser, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validName,
			Email:          validEmail,
			HashedPassword: validHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       createdUser.ID,
			SessionToken: validSessionToken.ID,
			ExpiredAt:    validSessionToken.ExpiredAt,
		})
		require.NoError(t, err)
	}

	expiredSessionToken := token.NewToken(-time.Minute)
	createExpiredSessionSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		createdUser, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validName,
			Email:          validEmail,
			HashedPassword: validHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       createdUser.ID,
			SessionToken: expiredSessionToken.ID,
			ExpiredAt:    expiredSessionToken.ExpiredAt,
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
			createSeed: createValidSessionSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(request *http.Request) {
			},
			buildStore: buildTestDBStore,
			createSeed: createValidSessionSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InvalidTokenFormat",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, "invalid")
			},
			buildStore: buildTestDBStore,
			createSeed: createValidSessionSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, expiredSessionToken.ID.String())
			},
			buildStore: buildTestDBStore,
			createSeed: createExpiredSessionSeed,
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

			server := newTestServer(t, store)

			authPath := "/auth"
			server.app.Get(
				authPath,
				authMiddleware(server),
				func(c *fiber.Ctx) error {
					return c.SendStatus(fiber.StatusOK)
				},
			)

			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(request)
			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}
