package test_util

import (
	"net/http"
	"time"

	mockdb "github.com/ot07/next-bazaar/db/mock"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	gomock "go.uber.org/mock/gomock"
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

func NewTokens(count int, duration time.Duration) []*token.Token {
	tokens := make([]*token.Token, count)
	for i := range tokens {
		tokens[i] = token.NewToken(duration)
	}
	return tokens
}
