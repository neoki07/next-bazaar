package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
)

const (
	cookieSessionTokenKey = "session_token"
	ctxLocalSessionKey    = "session"
)

func authMiddleware(server *Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionToken := c.Cookies(cookieSessionTokenKey)
		if len(sessionToken) == 0 {
			err := fmt.Errorf("session token not found")
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}

		parsedSessionToken, err := uuid.Parse(sessionToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}

		session, err := server.store.GetSession(c.Context(), parsedSessionToken)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
		}

		if token.IsExpired(session.SessionTokenExpiredAt) {
			if token.IsExpired(session.RefreshTokenExpiredAt) {
				return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(token.ErrExpiredToken))
			}

			newSession, err := refreshSessionToken(c, server, session)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
			}

			c.Locals(ctxLocalSessionKey, newSession)
		} else {
			c.Locals(ctxLocalSessionKey, session)
		}

		return c.Next()
	}
}

func getSession(c *fiber.Ctx) (db.Session, error) {
	session, ok := c.Locals(ctxLocalSessionKey).(db.Session)
	if !ok {
		return db.Session{}, fmt.Errorf("session token not found")
	}
	return session, nil
}

func refreshSessionToken(c *fiber.Ctx, server *Server, expiredSession db.Session) (db.Session, error) {
	err := server.store.DeleteSession(c.Context(), expiredSession.ID)
	if err != nil {
		return db.Session{}, err
	}

	newSessionToken := token.NewToken(time.Hour * 24 * 7)

	newSession, err := server.store.CreateSession(c.Context(), db.CreateSessionParams{
		UserID:                expiredSession.UserID,
		SessionToken:          newSessionToken.ID,
		SessionTokenExpiredAt: newSessionToken.ExpiredAt,
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
		MaxAge:   int(server.config.SessionTokenDuration.Seconds()),
	})

	return newSession, nil
}
