package api

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	user_domain "github.com/ot07/next-bazaar/api/domain/user"
	user_service "github.com/ot07/next-bazaar/api/service/user"
	"github.com/ot07/next-bazaar/api/validation"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/util"
)

type userHandler struct {
	service *user_service.UserService
	config  util.Config
}

func newUserHandler(s *user_service.UserService, config util.Config) *userHandler {
	return &userHandler{
		service: s,
		config:  config,
	}
}

// @Summary      Create user
// @Tags         Users
// @Param        body body user_domain.CreateUserRequest true "User object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      403 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users [post]
func (h *userHandler) createUser(c *fiber.Ctx) error {
	req := new(user_domain.CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	arg := user_service.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	err = h.service.CreateUser(c.Context(), arg)
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

// @Summary      Login user
// @Tags         Users
// @Param        body body user_domain.LoginUserRequest true "User object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      404 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/login [post]
func (h *userHandler) loginUser(c *fiber.Ctx) error {
	req := new(user_domain.LoginUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	user, err := h.service.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	sessionToken, err := h.service.CreateSession(c.Context(), user.ID, h.config.SessionTokenDuration)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(newErrorResponse(err))
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
		MaxAge:   int(h.config.SessionTokenDuration.Seconds()),
	})

	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Logout user
// @Tags         Users
// @Success      200 {object} messageResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/logout [post]
func (h *userHandler) logoutUser(c *fiber.Ctx) error {
	session, ok := c.Locals(ctxLocalSessionKey).(db.Session)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(
			fmt.Errorf("session token not found"),
		))
	}

	err := h.service.DeleteSession(c.Context(), session.SessionToken)
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

// @Summary      Get logged in user
// @Tags         Users
// @Success      200 {object} user_domain.UserResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/me [get]
func (h *userHandler) getLoggedInUser(c *fiber.Ctx) error {
	session, ok := c.Locals(ctxLocalSessionKey).(db.Session)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(
			fmt.Errorf("session token not found"),
		))
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
