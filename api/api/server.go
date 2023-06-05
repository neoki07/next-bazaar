package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	cart_domain "github.com/ot07/next-bazaar/api/domain/cart"
	product_domain "github.com/ot07/next-bazaar/api/domain/product"
	user_domain "github.com/ot07/next-bazaar/api/domain/user"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/util"
)

type handlers struct {
	user    *userHandler
	product *productHandler
	cart    *cartHandler
}

func newHandlers(config util.Config, store db.Store) handlers {
	/* User */
	userRepository := user_domain.NewUserRepository(store)
	userService := user_domain.NewUserService(userRepository)
	userHandler := newUserHandler(userService, config)

	/* Product */
	productRepository := product_domain.NewProductRepository(store)
	productService := product_domain.NewProductService(productRepository)
	productHandler := newProductHandler(productService)

	/* Cart */
	cartRepository := cart_domain.NewCartRepository(store)
	cartService := cart_domain.NewCartService(cartRepository)
	cartHandler := newCartHandler(cartService)

	return handlers{
		user:    userHandler,
		product: productHandler,
		cart:    cartHandler,
	}
}

// Server serves HTTP requests for this app domain.
type Server struct {
	config   util.Config
	store    db.Store
	app      *fiber.App
	handlers handlers
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,https://next-bazaar.vercel.app",
		AllowCredentials: true,
	}))

	server := &Server{
		config:   config,
		store:    store,
		app:      app,
		handlers: newHandlers(config, store),
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	app := server.app

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/users", server.handlers.user.createUser)
	v1.Post("/users/login", server.handlers.user.loginUser)

	v1.Get("/products", server.handlers.product.listProducts)
	v1.Get("/products/:id", server.handlers.product.getProduct)

	v1.Use(authMiddleware(server))

	v1.Post("/users/logout", server.handlers.user.logoutUser)
	v1.Get("/users/me", server.handlers.user.getLoggedInUser)

	v1.Get("/cart-products/:user_id", server.handlers.cart.getCart)
	v1.Post("/cart-products", server.handlers.cart.addProduct)
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
