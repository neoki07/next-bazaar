package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	user_domain "github.com/ot07/next-bazaar/api/domain/user"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "test-password"

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		checkResponse func(t *testing.T, response *http.Response)
		allowParallel bool
	}{
		{
			name: "OK",
			body: fiber.Map{
				"name":     validName,
				"email":    validEmail,
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name: "InternalError",
			body: fiber.Map{
				"name":     validName,
				"email":    validEmail,
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name: "DuplicateEmail",
			body: fiber.Map{
				"name":     validName,
				"email":    validEmail,
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				store, cleanup = newTestDBStore(t)

				_, err := store.CreateUser(context.Background(), db.CreateUserParams{
					Name:           "testuser0",
					Email:          validEmail,
					HashedPassword: "hashed_password",
				})
				require.NoError(t, err)

				return
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusForbidden, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name: "NameWithSpace",
			body: fiber.Map{
				"name":     "testuser ",
				"email":    validEmail,
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name: "NameWithPunct",
			body: fiber.Map{
				"name":     "testuser!",
				"email":    validEmail,
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name: "NameWithSymbol",
			body: fiber.Map{
				"name":     "testuser|",
				"email":    validEmail,
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name: "InvalidEmail",
			body: fiber.Map{
				"name":     validName,
				"email":    "invalid-email",
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name: "TooShortPassword",
			body: fiber.Map{
				"name":     validName,
				"email":    validEmail,
				"password": "1234567",
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: true,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			if tc.allowParallel {
				t.Parallel()
			}

			store, cleanupStore := tc.buildStore(t)
			defer cleanupStore()

			server := newTestServer(t, store)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "test-password"

	validHashedPassword, err := util.HashPassword(validPassword)
	require.NoError(t, err)

	createSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validName,
			Email:          validEmail,
			HashedPassword: validHashedPassword,
		})
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"email":    validEmail,
				"password": validPassword,
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
			name: "InvalidEmail",
			body: fiber.Map{
				"email":    "invalid-email",
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "TooShortPassword",
			body: fiber.Map{
				"email":    validEmail,
				"password": "1234567",
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "NoExistsUser",
			body: fiber.Map{
				"email":    validEmail,
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "GetUserInternalError",
			body: fiber.Map{
				"email":    validEmail,
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(validEmail)).
					Return(db.User{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "CreateSessionInternalError",
			body: fiber.Map{
				"email":    validEmail,
				"password": validPassword,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				user := db.User{
					ID:             util.RandomUUID(),
					Name:           validName,
					Email:          validEmail,
					HashedPassword: validHashedPassword,
				}
				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(validEmail)).
					Return(user, nil)

				mockStore.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Return(db.Session{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "MistakePassword",
			body: fiber.Map{
				"email":    validEmail,
				"password": "12345678",
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store, cleanupStore := tc.buildStore(t)
			defer cleanupStore()

			tc.createSeed(t, store)

			server := newTestServer(t, store)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/users/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestLogoutUserAPI(t *testing.T) {
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "test-password"

	validHashedPassword, err := util.HashPassword(validPassword)
	require.NoError(t, err)

	sessionToken := token.NewToken(time.Minute)

	createSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		createdUser, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validName,
			Email:          validEmail,
			HashedPassword: validHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       createdUser.ID,
			SessionToken: sessionToken.ID,
			ExpiredAt:    sessionToken.ExpiredAt,
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
				addSessionTokenInCookie(request, sessionToken.ID.String())
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
			name: "NoExistsSessionToken",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, util.RandomUUID().String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, sessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				buildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: sessionToken.ID,
					ExpiredAt:    sessionToken.ExpiredAt,
					CreatedAt:    time.Now(),
				})

				mockStore.EXPECT().
					DeleteSession(gomock.Any(), gomock.Eq(sessionToken.ID)).
					Return(sql.ErrConnDone)

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
			store, cleanupStore := tc.buildStore(t)
			defer cleanupStore()

			tc.createSeed(t, store)

			server := newTestServer(t, store)

			url := "/api/v1/users/logout"
			request, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(request)
			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestGetLoggedInUserAPI(t *testing.T) {
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "test-password"

	validHashedPassword, err := util.HashPassword(validPassword)
	require.NoError(t, err)

	sessionToken := token.NewToken(time.Minute)

	createSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		createdUser, err := store.CreateUser(ctx, db.CreateUserParams{
			Name:           validName,
			Email:          validEmail,
			HashedPassword: validHashedPassword,
		})
		require.NoError(t, err)

		_, err = store.CreateSession(ctx, db.CreateSessionParams{
			UserID:       createdUser.ID,
			SessionToken: sessionToken.ID,
			ExpiredAt:    sessionToken.ExpiredAt,
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
				addSessionTokenInCookie(request, sessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotUser := unmarshalUserResponse(t, response.Body)

				require.Equal(t, validName, gotUser.Name)
				require.Equal(t, validEmail, gotUser.Email)
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
			name: "NoExistsUser",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, util.RandomUUID().String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, sessionToken.ID.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				session := db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: sessionToken.ID,
					ExpiredAt:    sessionToken.ExpiredAt,
					CreatedAt:    time.Now(),
				}

				buildValidSessionStubs(mockStore, session)

				mockStore.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(session.UserID)).
					Return(db.User{}, sql.ErrConnDone)

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
			store, cleanupStore := tc.buildStore(t)
			defer cleanupStore()

			tc.createSeed(t, store)

			server := newTestServer(t, store)

			url := "/api/v1/users/me"
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

func unmarshalUserResponse(t *testing.T, body io.ReadCloser) user_domain.UserResponse {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var parsed user_domain.UserResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	return parsed
}
