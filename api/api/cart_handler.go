package api

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	cart_domain "github.com/ot07/next-bazaar/api/domain/cart"
	cart_service "github.com/ot07/next-bazaar/api/service/cart"
	"github.com/ot07/next-bazaar/api/validation"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type cartHandler struct {
	service *cart_service.CartService
}

func newCartHandler(s *cart_service.CartService) *cartHandler {
	return &cartHandler{
		service: s,
	}
}

// @Summary      Get cart
// @Tags         cart
// @Param        userId path string true "User ID"
// @Success      200 {object} productResponse
// @Failure      400 {object} errorResponse
// @Failure      404 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart-products/{userId} [get]
func (h *cartHandler) getCart(c *fiber.Ctx) error {
	req := new(cart_domain.GetProductsRequest)
	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	cartProducts, err := h.service.GetProductsByUserID(c.Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := cart_domain.NewCartResponse(cartProducts)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Add product to cart
// @Tags         cart
// @Param        body body addProductRequest true "Cart product object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      403 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart-products [post]
func (h *cartHandler) addProduct(c *fiber.Ctx) error {
	session, ok := c.Locals(ctxLocalSessionKey).(db.Session)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(
			fmt.Errorf("session token not found"),
		))
	}

	req := new(cart_domain.AddProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err := h.service.AddProduct(c.Context(), cart_service.NewAddProductParams(
		session.UserID,
		req.ProductID,
		req.Quantity,
	))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.Status(fiber.StatusForbidden).JSON(newErrorResponse(err))
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Cart product added successfully")
	return c.Status(fiber.StatusOK).JSON(rsp)
}
