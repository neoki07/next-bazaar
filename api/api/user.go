package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,without_space,without_punct,without_symbol"`
	Email    string `json:"email" validate:"required,email" swaggertype:"string"`
	Password string `json:"password" validate:"required,min=8"`
}

// @Summary      Create user
// @Tags         users
// @Param        body body createUserRequest true "User object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      403 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users [post]
func (server *Server) createUser(c *fiber.Ctx) error {
	req := new(createUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := newValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	_, err = server.store.CreateUser(c.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.Status(fiber.StatusForbidden).JSON(newErrorResponse(err))
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := messageResponse{
		Message: "Congratulations! You are now a member of our online bazaar. Start exploring!",
	}

	return c.Status(fiber.StatusOK).JSON(rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" validate:"required,email" swaggertype:"string"`
	Password string `json:"password" validate:"required,min=8"`
}

// @Summary      Login user
// @Tags         users
// @Param        body body loginUserRequest true "User object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      404 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/login [post]
func (server *Server) loginUser(c *fiber.Ctx) error {
	req := new(loginUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	user, err := server.store.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	sessionToken := token.NewToken(
		server.config.SessionTokenDuration,
	)

	arg := db.CreateSessionParams{
		UserID:       user.ID,
		SessionToken: sessionToken.ID,
		ExpiredAt:    sessionToken.ExpiredAt,
	}

	_, err = server.store.CreateSession(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := messageResponse{
		Message: "Welcome to our online bazaar! Get ready to discover unique treasures and amazing deals.",
	}

	c.Cookie(&fiber.Cookie{
		Name:     sessionTokenKey,
		Value:    sessionToken.ID.String(),
		HTTPOnly: true,
		SameSite: "none",
		Secure:   true,
		MaxAge:   int(server.config.SessionTokenDuration.Seconds()),
	})

	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Logout user
// @Tags         users
// @Success      200 {object} messageResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/logout [post]
func (server *Server) logoutUser(c *fiber.Ctx) error {
	sessionToken := c.Locals(sessionTokenKey).(uuid.UUID)

	err := server.store.DeleteSession(c.Context(), sessionToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := messageResponse{
		Message: "Thank you for visiting us, we look forward to your next visit!",
	}

	c.ClearCookie(sessionTokenKey)

	return c.Status(fiber.StatusOK).JSON(rsp)
}

type userResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email" swaggertype:"string"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

// @Summary      Get logged in user
// @Tags         users
// @Success      200 {object} userResponse
// @Failure      401 {object} errorResponse
// @Failure      404 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/me [get]
func (server *Server) getLoggedInUser(c *fiber.Ctx) error {
	sessionToken := c.Locals(sessionTokenKey).(uuid.UUID)

	session, err := server.store.GetSession(c.Context(), sessionToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	user, err := server.store.GetUser(c.Context(), session.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(rsp)
}
