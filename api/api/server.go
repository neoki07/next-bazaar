package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/util"
)

// Server serves HTTP requests for this app service.
type Server struct {
	config util.Config
	store  db.Store
	app    *fiber.App
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
	}))

	server := &Server{
		config: config,
		store:  store,
		app:    app,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	app := server.app

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/users", server.createUser)
	v1.Post("/users/login", server.loginUser)

	v1.Use(authMiddleware(server))

	v1.Post("/users/logout", server.logoutUser)
	v1.Get("/users/me", server.getLoggedInUser)
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.app.Listen(address)
}

type messageResponse struct {
	Message string `json:"message"`
}

func newMessageResponse(message string) messageResponse {
	return messageResponse{Message: message}
}

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(err error) errorResponse {
	return errorResponse{Error: err.Error()}
}
