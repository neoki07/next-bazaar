package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	user_domain "github.com/ot07/next-bazaar/api/domain/user"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	t.Parallel()

	user, password := randomUser(t)

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
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
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
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
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
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				store, cleanup = newTestDBStore(t)

				_, err := store.CreateUser(context.Background(), db.CreateUserParams{
					Name:           "user",
					Email:          user.Email,
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
				"name":     "user ",
				"email":    user.Email,
				"password": password,
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
				"name":     "user!",
				"email":    user.Email,
				"password": password,
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
				"name":     "user|",
				"email":    user.Email,
				"password": password,
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
				"name":     user.Name,
				"email":    "invalid-email",
				"password": password,
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
				"name":     user.Name,
				"email":    user.Email,
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
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"email":    user.Email,
				"password": password,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
				mockStore.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, nil)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name: "InvalidEmail",
			body: fiber.Map{
				"email":    "invalid-email",
				"password": password,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("invalid-email")).
					Times(0)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "TooShortPassword",
			body: fiber.Map{
				"email":    user.Email,
				"password": "1234567",
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(0)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "NoExistsUser",
			body: fiber.Map{
				"email":    user.Email,
				"password": password,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
				mockStore.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "GetUserInternalError",
			body: fiber.Map{
				"email":    user.Email,
				"password": password,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "CreateSessionInternalError",
			body: fiber.Map{
				"email":    user.Email,
				"password": password,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
				mockStore.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "MistakePassword",
			body: fiber.Map{
				"email":    user.Email,
				"password": "12345678",
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)

				return mockStore, cleanup
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
	session := randomSession()

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, session.SessionToken.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				buildValidSessionStubs(mockStore, session)
				mockStore.EXPECT().
					DeleteSession(gomock.Any(), gomock.Eq(session.SessionToken)).
					Times(1).
					Return(nil)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(request *http.Request) {},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newMockStore(t)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "NoExistsSessionToken",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, session.SessionToken.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				buildValidSessionStubs(mockStore, session)
				mockStore.EXPECT().
					DeleteSession(gomock.Any(), gomock.Eq(session.SessionToken)).
					Times(1).
					Return(sql.ErrNoRows)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, session.SessionToken.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				buildValidSessionStubs(mockStore, session)
				mockStore.EXPECT().
					DeleteSession(gomock.Any(), gomock.Eq(session.SessionToken)).
					Times(1).
					Return(sql.ErrConnDone)

				return mockStore, cleanup
			},
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
	user, _ := randomUser(t)
	session := randomExistsUserSession(user.ID)

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request)
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, session.SessionToken.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				buildValidSessionStubs(mockStore, session)
				mockStore.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(session.UserID)).
					Times(1).
					Return(user, nil)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchUser(t, response.Body, user)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(request *http.Request) {},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newMockStore(t)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "NoExistsUser",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, session.SessionToken.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				buildValidSessionStubs(mockStore, session)
				mockStore.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(session.UserID)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)

				return mockStore, cleanup
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, session.SessionToken.String())
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				buildValidSessionStubs(mockStore, session)
				mockStore.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(session.UserID)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
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

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:             util.RandomUUID(),
		Name:           util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}
	return
}

func randomExistsUserSession(userID uuid.UUID) db.Session {
	return db.Session{
		ID:           util.RandomUUID(),
		UserID:       userID,
		SessionToken: util.RandomUUID(),
		ExpiredAt:    time.Now().Add(time.Minute),
	}
}

func requireBodyMatchUser(t *testing.T, body io.ReadCloser, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser user_domain.UserResponse
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	requireUserResponseMatchUser(t, gotUser, user)

	err = body.Close()
	require.NoError(t, err)
}

func requireUserResponseMatchUser(t *testing.T, gotUser user_domain.UserResponse, user db.User) {
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.Email, gotUser.Email)
}
