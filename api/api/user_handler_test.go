package api

import (
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
	"github.com/ot07/next-bazaar/api/test_util"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "test-password"

	defaultBody := fiber.Map{
		"name":     validName,
		"email":    validEmail,
		"password": validPassword,
	}

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		checkResponse func(t *testing.T, response *http.Response)
		allowParallel bool
	}{
		{
			name:       "OK",
			body:       defaultBody,
			buildStore: test_util.BuildTestDBStore,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name: "InternalError",
			body: defaultBody,
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
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
			body: defaultBody,
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				store, cleanup = test_util.NewTestDBStore(t)

				_ = test_util.CreateUserTestData(t, context.Background(), store,
					"testuser2",
					validEmail,
					validPassword,
					token.NewToken(time.Minute),
				)

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
			buildStore: test_util.BuildTestDBStore,
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
			buildStore: test_util.BuildTestDBStore,
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
			buildStore: test_util.BuildTestDBStore,
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
			buildStore: test_util.BuildTestDBStore,
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
			buildStore: test_util.BuildTestDBStore,
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

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodPost,
				URL:    "/api/v1/users",
				Body:   tc.body,
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "test-password"

	defaultBody := fiber.Map{
		"email":    validEmail,
		"password": validPassword,
	}

	createSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateUserTestData(t, ctx, store,
			validName,
			validEmail,
			validPassword,
			token.NewToken(time.Minute),
		)

	}

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			body:       defaultBody,
			buildStore: test_util.BuildTestDBStore,
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
			buildStore: test_util.BuildTestDBStore,
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
			buildStore: test_util.BuildTestDBStore,
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:       "NoExistsUser",
			body:       defaultBody,
			buildStore: test_util.BuildTestDBStore,
			createSeed: func(t *testing.T, store db.Store) {},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "GetUserInternalError",
			body: defaultBody,
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
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
			body: defaultBody,
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				validHashedPassword, err := util.HashPassword(validPassword)
				require.NoError(t, err)

				user := db.User{
					ID:             util.RandomUUID(),
					Name:           validName,
					Email:          validEmail,
					HashedPassword: validHashedPassword,
				}
				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
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
			buildStore: test_util.BuildTestDBStore,
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

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodPost,
				URL:    "/api/v1/users/login",
				Body:   tc.body,
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestLogoutUserAPI(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	createSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateUserTestData(t, ctx, store,
			"testuser",
			"test@example.com",
			"test-password",
			sessionToken,
		)
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
				test_util.AddSessionTokenInCookie(cookieSessionTokenKey, sessionToken.ID.String(), request)
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  func(request *http.Request) {},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "NoExistsSessionToken",
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(cookieSessionTokenKey, util.RandomUUID().String(), request)
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(cookieSessionTokenKey, sessionToken.ID.String(), request)
			},
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
					DeleteSession(gomock.Any(), gomock.Any()).
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

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method:    http.MethodPost,
				URL:       "/api/v1/users/logout",
				SetupAuth: tc.setupAuth,
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestGetLoggedInUserAPI(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	createSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateUserTestData(t, ctx, store,
			"testuser",
			"test@example.com",
			"test-password",
			sessionToken,
		)
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
				test_util.AddSessionTokenInCookie(cookieSessionTokenKey, sessionToken.ID.String(), request)
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotUser := unmarshalUserResponse(t, response.Body)

				require.Equal(t, "testuser", gotUser.Name)
				require.Equal(t, "test@example.com", gotUser.Email)
			},
		},
		{
			name:       "NoAuthorization",
			setupAuth:  func(request *http.Request) {},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "NoExistsUser",
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(cookieSessionTokenKey, util.RandomUUID().String(), request)
			},
			buildStore: test_util.BuildTestDBStore,
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				test_util.AddSessionTokenInCookie(cookieSessionTokenKey, sessionToken.ID.String(), request)
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				session := db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: sessionToken.ID,
					ExpiredAt:    sessionToken.ExpiredAt,
					CreatedAt:    time.Now(),
				}

				test_util.BuildValidSessionStubs(mockStore, session)

				mockStore.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
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

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method:    http.MethodGet,
				URL:       "/api/v1/users/me",
				SetupAuth: tc.setupAuth,
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
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
