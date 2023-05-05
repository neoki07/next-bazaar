package api

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
			err := errors.New("session token not found")
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

		token := token.Token{
			ID:        session.ID,
			ExpiredAt: session.ExpiredAt,
		}

		err = token.Valid()
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}

		c.Locals(ctxLocalSessionKey, session)
		return c.Next()
	}
}
