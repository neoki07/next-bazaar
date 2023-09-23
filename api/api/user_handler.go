package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	user_domain "github.com/ot07/next-bazaar/api/domain/user"
	"github.com/ot07/next-bazaar/api/validation"
	"github.com/ot07/next-bazaar/util"
	"golang.org/x/crypto/bcrypt"
)

type userHandler struct {
	service *user_domain.UserService
	config  util.Config
}

func newUserHandler(s *user_domain.UserService, config util.Config) *userHandler {
	return &userHandler{
		service: s,
		config:  config,
	}
}

// @Summary      Register user
// @Tags         Users
// @Param        body body user_domain.RegisterRequest true "User object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      403 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/register [post]
func (h *userHandler) register(c *fiber.Ctx) error {
	req := new(user_domain.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err := h.service.Register(c.Context(), user_domain.RegisterServiceParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.Status(fiber.StatusForbidden).JSON(newErrorResponse(err))
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Congratulations! You are now a member of our online bazaar. Start exploring!")
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Login
// @Tags         Users
// @Param        body body user_domain.LoginRequest true "User object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/login [post]
func (h *userHandler) login(c *fiber.Ctx) error {
	req := new(user_domain.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	sessionToken, err := h.service.Login(c.Context(), user_domain.LoginServiceParams{
		Email:                req.Email,
		Password:             req.Password,
		SessionTokenDuration: h.config.SessionTokenDuration,
		RefreshTokenDuration: h.config.RefreshTokenDuration,
	})
	if err != nil {
		if err == sql.ErrNoRows || err == bcrypt.ErrMismatchedHashAndPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Welcome to our online bazaar! Get ready to discover unique treasures and amazing deals.")

	c.Cookie(&fiber.Cookie{
		Name:     cookieSessionTokenKey,
		Value:    sessionToken.ID.String(),
		HTTPOnly: true,
		SameSite: "none",
		Secure:   true,
		MaxAge:   int(h.config.RefreshTokenDuration.Seconds()),
	})

	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Logout
// @Tags         Users
// @Success      200 {object} messageResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/logout [post]
func (h *userHandler) logout(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	err = h.service.Logout(c.Context(), session.SessionToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Thank you for visiting us, we look forward to your next visit!")

	c.ClearCookie(cookieSessionTokenKey)

	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Get current user
// @Tags         Users
// @Success      200 {object} user_domain.UserResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/me [get]
func (h *userHandler) getCurrentUser(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	user, err := h.service.GetUser(c.Context(), session.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := user_domain.NewUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Update user information
// @Tags         Users
// @Param        body body user_domain.UpdateRequest true "User object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      403 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/me [patch]
func (h *userHandler) updateCurrentUser(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	req := new(user_domain.UpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = h.service.UpdateUser(c.Context(), user_domain.UpdateUserServiceParams{
		ID:    session.UserID,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.Status(fiber.StatusForbidden).JSON(newErrorResponse(err))
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Your information has been updated successfully!")
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Update user password
// @Tags         Users
// @Param        body body user_domain.UpdatePasswordRequest true "User object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/me/password [patch]
func (h *userHandler) updateCurrentUserPassword(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	req := new(user_domain.UpdatePasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = h.service.UpdateUserPassword(c.Context(), user_domain.UpdateUserPasswordServiceParams{
		ID:          session.UserID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Your password has been updated successfully!")
	return c.Status(fiber.StatusOK).JSON(rsp)
}
