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

	defaultCreateSeedData := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         validName,
			Email:        validEmail,
			Password:     validPassword,
			SessionToken: token.NewToken(time.Minute),
		})
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store)
		body           test_util.Body
		checkResponse  func(t *testing.T, response *http.Response)
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body:           defaultBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:           "InvalidEmail",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"email":    "invalid-email",
				"password": validPassword,
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "TooShortPassword",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"email":    validEmail,
				"password": "1234567",
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:           "NoExistsUser",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: test_util.NoopCreateSeedData,
			body:           defaultBody,
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
			createSeedData: test_util.NoopCreateSeedData,
			body:           defaultBody,
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
			createSeedData: test_util.NoopCreateSeedData,
			body:           defaultBody,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name:           "MistakePassword",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
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

			tc.createSeedData(t, store)

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

	defaultCreateSeedData := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
		})
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store)
		setupAuth      func(request *http.Request, sessionToken string)
		checkResponse  func(t *testing.T, response *http.Response)
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			setupAuth:      test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:           "NoAuthorization",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			setupAuth:      test_util.NoopSetupAuth,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:           "NoExistsSessionToken",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
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
			createSeedData: test_util.NoopCreateSeedData,
			setupAuth:      test_util.AddSessionTokenInCookie,
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

			tc.createSeedData(t, store)

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

	defaultCreateSeedData := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     "test-password",
			SessionToken: sessionToken,
		})
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store)
		setupAuth      func(request *http.Request, sessionToken string)
		checkResponse  func(t *testing.T, response *http.Response)
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			setupAuth:      test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)

				gotUser := unmarshalUserResponse(t, response.Body)

				require.Equal(t, "testuser", gotUser.Name)
				require.Equal(t, "test@example.com", gotUser.Email)
			},
		},
		{
			name:           "NoAuthorization",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			setupAuth:      test_util.NoopSetupAuth,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:           "NoExistsUser",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
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
			createSeedData: test_util.NoopCreateSeedData,
			setupAuth:      test_util.AddSessionTokenInCookie,
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

			tc.createSeedData(t, store)

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

func TestUpdateCurrentUserAPI(t *testing.T) {
	validName := "testuser"
	validName2 := "testuser2"

	validEmail := "test@example.com"
	validEmail2 := "test2@example.com"

	validPassword := "test-password"
	validPassword2 := "test-password2"

	validSessionToken := token.NewToken(time.Minute)
	validSessionToken2 := token.NewToken(time.Minute)

	defaultCreateSeedData := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         validName,
			Email:        validEmail,
			Password:     validPassword,
			SessionToken: validSessionToken,
		})

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         validName2,
			Email:        validEmail2,
			Password:     validPassword2,
			SessionToken: validSessionToken2,
		})
	}

	defaultBody := test_util.Body{
		"name":  validName,
		"email": validEmail,
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store)
		body           test_util.Body
		setupAuth      func(request *http.Request, sessionToken string)
		checkResponse  func(t *testing.T, response *http.Response)
		allowParallel  bool
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body:           defaultBody,
			setupAuth:      test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name: "InternalError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				test_util.BuildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: validSessionToken.ID,
					ExpiredAt:    validSessionToken.ExpiredAt,
					CreatedAt:    time.Now(),
				})

				mockStore.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Return(db.User{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeedData: test_util.NoopCreateSeedData,
			body:           defaultBody,
			setupAuth:      test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "NameAlreadyExists",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"name":  validName2,
				"email": validEmail,
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusForbidden, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "EmailAlreadyExists",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"name":  validName,
				"email": validEmail2,
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusForbidden, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "NameWithSpace",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"name":  "testuser ",
				"email": validEmail,
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "NameWithPunct",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"name":  "testuser!",
				"email": validEmail,
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "NameWithSymbol",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"name":  "testuser|",
				"email": validEmail,
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "InvalidEmail",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"name":  validName,
				"email": "invalid-email",
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: false,
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

			tc.createSeedData(t, store)

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodPatch,
				URL:    "/api/v1/users/me",
				Body:   tc.body,
			})

			tc.setupAuth(request, validSessionToken.ID.String())

			server := newTestServer(t, store)
			response := test_util.SendRequest(t, server.app, request)
			tc.checkResponse(t, response)
		})
	}
}

func TestUpdateCurrentUserPasswordAPI(t *testing.T) {
	validOldPassword := "test-old-password"
	validNewPassword := "test-new-password"

	validSessionToken := token.NewToken(time.Minute)

	defaultCreateSeedData := func(t *testing.T, store db.Store) {
		ctx := context.Background()

		_ = test_util.CreateWithSessionUser(t, ctx, store, test_util.WithSessionUserParams{
			Name:         "testuser",
			Email:        "test@example.com",
			Password:     validOldPassword,
			SessionToken: validSessionToken,
		})
	}

	defaultBody := test_util.Body{
		"old_password": validOldPassword,
		"new_password": validNewPassword,
	}

	testCases := []struct {
		name           string
		buildStore     func(t *testing.T) (store db.Store, cleanup func())
		createSeedData func(t *testing.T, store db.Store)
		body           test_util.Body
		setupAuth      func(request *http.Request, sessionToken string)
		checkResponse  func(t *testing.T, response *http.Response)
		allowParallel  bool
	}{
		{
			name:           "OK",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body:           defaultBody,
			setupAuth:      test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name: "InternalError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := test_util.NewMockStore(t)

				test_util.BuildValidSessionStubs(mockStore, db.Session{
					ID:           util.RandomUUID(),
					UserID:       util.RandomUUID(),
					SessionToken: validSessionToken.ID,
					ExpiredAt:    validSessionToken.ExpiredAt,
					CreatedAt:    time.Now(),
				})

				mockStore.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Return(db.User{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeedData: test_util.NoopCreateSeedData,
			body:           defaultBody,
			setupAuth:      test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "NoAuthorization",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body:           defaultBody,
			setupAuth:      test_util.NoopSetupAuth,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "OldPasswordNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"new_password": validNewPassword,
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "NewPasswordNotFound",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"old_password": validOldPassword,
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "MismatchOldPassword",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"old_password": "12345678",
				"new_password": validNewPassword,
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "TooShortOldPassword",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"old_password": "1234567",
				"new_password": validNewPassword,
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: false,
		},
		{
			name:           "TooShortNewPassword",
			buildStore:     test_util.BuildTestDBStore,
			createSeedData: defaultCreateSeedData,
			body: test_util.Body{
				"old_password": validOldPassword,
				"new_password": "1234567",
			},
			setupAuth: test_util.AddSessionTokenInCookie,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			allowParallel: false,
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

			tc.createSeedData(t, store)

			request := test_util.NewRequest(t, test_util.RequestParams{
				Method: http.MethodPatch,
				URL:    "/api/v1/users/me/password",
				Body:   tc.body,
			})

			tc.setupAuth(request, validSessionToken.ID.String())

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
