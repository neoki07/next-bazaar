package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	user_domain "github.com/ot07/next-bazaar/api/domain/user"
	"github.com/ot07/next-bazaar/api/test_util"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func TestRegisterAPI(t *testing.T) {
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "test-password"

	defaultBody := test_util.Body{
		"name":     validName,
		"email":    validEmail,
		"password": validPassword,
	}

	testCases := []struct {
		name          string
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		body          test_util.Body
		checkResponse func(t *testing.T, response *http.Response)
		allowParallel bool
	}{
		{
			name:       "OK",
			buildStore: test_util.BuildTestDBStore,
			body:       defaultBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name: "InternalError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(db.User{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			body: defaultBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name: "DuplicateEmail",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				store, cleanup = test_util.NewTestDBStore(t)

				_ = test_util.CreateWithSessionUser(t, context.Background(), store, test_util.WithSessionUserParams{
					Name:         "testuser2",
					Email:        validEmail,
					Password:     validPassword,
					SessionToken: token.NewToken(time.Minute),
				})

				return
			},
			body: defaultBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusForbidden, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:       "NameWithSpace",
			buildStore: test_util.BuildTestDBStore,
			body: test_util.Body{
				"name":     "testuser ",
				"email":    validEmail,
				"password": validPassword,
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name:       "NameWithPunct",
			buildStore: test_util.BuildTestDBStore,
			body: test_util.Body{
				"name":     "testuser!",
				"email":    validEmail,
				"password": validPassword,
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name:       "NameWithSymbol",
			buildStore: test_util.BuildTestDBStore,
			body: test_util.Body{
				"name":     "testuser|",
				"email":    validEmail,
				"password": validPassword,
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name:       "InvalidEmail",
			buildStore: test_util.BuildTestDBStore,
			body: test_util.Body{
				"name":     validName,
				"email":    "invalid-email",
				"password": validPassword,
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: true,
		},
		{
			name:       "TooShortPassword",
			buildStore: test_util.BuildTestDBStore,
			body: test_util.Body{
				"name":     validName,
				"email":    validEmail,
				"password": "1234567",
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

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodPost,
				URL:    "/api/v1/users/register",
				Body:   tc.body,
			})

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestLoginAPI(t *testing.T) {
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "test-password"

	defaultBody := test_util.Body{
		"email":    validEmail,
		"password": validPassword,
	}

	defaultCreateSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         validName,
			Email:        validEmail,
			Password:     validPassword,
			SessionToken: token.NewToken(time.Minute),
		})
	}

	testCases := []struct {
		name          string
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		body          test_util.Body
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			body:       defaultBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "InvalidEmail",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			body: test_util.Body{
				"email":    "invalid-email",
				"password": validPassword,
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:       "TooShortPassword",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			body: test_util.Body{
				"email":    validEmail,
				"password": "1234567",
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:       "NoExistsUser",
			buildStore: test_util.BuildTestDBStore,
			createSeed: test_util.NoopCreateSeed,
			body:       defaultBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "GetUserInternalError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Return(db.User{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: test_util.NoopCreateSeed,
			body:       defaultBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "CreateSessionInternalError",
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
			createSeed: test_util.NoopCreateSeed,
			body:       defaultBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name:       "MistakePassword",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			body: test_util.Body{
				"email":    validEmail,
				"password": "12345678",
			},
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

func TestLogoutAPI(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
		})
	}

	testCases := []struct {
		name          string
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		setupAuth     func(request *http.Request, sessionToken string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			setupAuth:  test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:       "NoAuthorization",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			setupAuth:  test_util.NoopSetupAuth,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:       "NoExistsSessionToken",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			setupAuth: func(request *http.Request, sessionToken string) {
				test_util.AddSessionTokenInCookie(request, util.RandomUUID().String())
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
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
			createSeed: test_util.NoopCreateSeed,
			setupAuth:  test_util.AddSessionTokenInCookie,
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
				Method: http.MethodPost,
				URL:    "/api/v1/users/logout",
			})

			tc.setupAuth(request, sessionToken.ID.String())

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestGetCurrentUserAPI(t *testing.T) {
	sessionToken := token.NewToken(time.Minute)

	defaultCreateSeed := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
		})
	}

	testCases := []struct {
		name          string
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store)
		setupAuth     func(request *http.Request, sessionToken string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			setupAuth:  test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotUser := unmarshalUserResponse(t, response.Body)

				require.Equal(t, "testuser", gotUser.Name)
				require.Equal(t, "test@example.com", gotUser.Email)
			},
		},
		{
			name:       "NoAuthorization",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			setupAuth:  test_util.NoopSetupAuth,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:       "NoExistsUser",
			buildStore: test_util.BuildTestDBStore,
			createSeed: defaultCreateSeed,
			setupAuth: func(request *http.Request, sessionToken string) {
				test_util.AddSessionTokenInCookie(request, util.RandomUUID().String())
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
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
			createSeed: test_util.NoopCreateSeed,
			setupAuth:  test_util.AddSessionTokenInCookie,
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
				Method: http.MethodGet,
				URL:    "/api/v1/users/me",
			})

			tc.setupAuth(request, sessionToken.ID.String())

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
