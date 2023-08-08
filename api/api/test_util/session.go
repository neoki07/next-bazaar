package test_util

import (
	"net/http"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/ot07/next-bazaar/db/mock"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
)

const (
	cookieSessionTokenKey = "session_token"
)

func AddSessionTokenInCookie(
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

func BuildValidSessionStubs(store *mockdb.MockStore, session db.Session) {
	store.EXPECT().
		GetSession(gomock.Any(), gomock.Any()).
		Return(session, nil)
}

func NewSessionTokens(count int, duration time.Duration) []*token.Token {
	sessionTokens := make([]*token.Token, count)
	for i := range sessionTokens {
		sessionTokens[i] = token.NewToken(duration)
	}
	return sessionTokens
}
