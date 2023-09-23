package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
)

func getSession(c *fiber.Ctx) (db.Session, error) {
	session, ok := c.Locals(ctxLocalSessionKey).(db.Session)
	if !ok {
		return db.Session{}, fmt.Errorf("session token not found")
	}
	return session, nil
}

func refreshSessionToken(c *fiber.Ctx, server *Server, expiredSession db.Session) (db.Session, error) {
	// TODO: setup transaction
	if token.IsExpired(expiredSession.RefreshTokenExpiredAt) {
		return db.Session{}, token.ErrExpiredToken
	}

	err := server.store.DeleteSession(c.Context(), expiredSession.ID)
	if err != nil {
		return db.Session{}, err
	}

	newSessionToken := token.NewToken(server.config.SessionTokenDuration)

	newSession, err := server.store.CreateSession(c.Context(), db.CreateSessionParams{
		UserID:                expiredSession.UserID,
		SessionToken:          newSessionToken.ID,
		SessionTokenExpiredAt: newSessionToken.ExpiredAt,
		RefreshToken:          expiredSession.RefreshToken,
		RefreshTokenExpiredAt: expiredSession.RefreshTokenExpiredAt,
	})
	if err != nil {
		return db.Session{}, err
	}

	c.Cookie(&fiber.Cookie{
		Name:     cookieSessionTokenKey,
		Value:    newSessionToken.ID.String(),
		HTTPOnly: true,
		SameSite: "none",
		Secure:   true,
		MaxAge:   int(server.config.RefreshTokenDuration.Seconds()),
	})

	return newSession, nil
}
