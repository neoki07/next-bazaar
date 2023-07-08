package test_util

import (
	"net/http"

	"github.com/golang/mock/gomock"
	mockdb "github.com/ot07/next-bazaar/db/mock"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

func AddSessionTokenInCookie(
	key string,
	sessionToken string,
	request *http.Request,
) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	request.AddCookie(cookie)
}

func BuildValidSessionStubs(store *mockdb.MockStore, session db.Session) {
	store.EXPECT().
		GetSession(gomock.Any(), gomock.Any()).
		Return(session, nil)
}
