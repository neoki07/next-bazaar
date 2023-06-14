package api

import (
	"github.com/gofiber/fiber/v2"
	cart_domain "github.com/ot07/next-bazaar/api/domain/cart"
	"github.com/ot07/next-bazaar/api/validation"
)

type cartHandler struct {
	service *cart_domain.CartService
}

func newCartHandler(s *cart_domain.CartService) *cartHandler {
	return &cartHandler{
		service: s,
	}
}

// @Summary      Get cart
// @Tags         Cart
// @Success      200 {object} cart_domain.CartResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart-products [get]
func (h *cartHandler) getCart(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	cartProducts, err := h.service.GetProductsByUserID(c.Context(), session.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := cart_domain.NewCartResponse(cartProducts)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Add product to cart
// @Tags         Cart
// @Param        body body cart_domain.AddProductRequest true "Cart product object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart-products [post]
func (h *cartHandler) addProduct(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	req := new(cart_domain.AddProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = h.service.AddProduct(c.Context(), cart_domain.NewAddProductParams(
		session.UserID,
		req.ProductID,
		req.Quantity,
	))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Cart product added successfully")
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Update cart product quantity
// @Tags         Cart
// @Param        body body cart_domain.UpdateProductQuantityRequest true "Cart product object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart-products [put]
func (h *cartHandler) updateProductQuantity(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	req := new(cart_domain.UpdateProductQuantityRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = h.service.UpdateProductQuantity(c.Context(), cart_domain.NewUpdateProductQuantityParams(
		session.UserID,
		req.ProductID,
		req.Quantity,
	))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Cart product added successfully")
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Delete cart product
// @Tags         Cart
// @Param        body body cart_domain.DeleteProductRequest true "Cart product object"
// @Success      204
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart-products [delete]
func (h *cartHandler) deleteProduct(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	req := new(cart_domain.DeleteProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = h.service.DeleteProduct(c.Context(), cart_domain.NewDeleteProductParams(
		session.UserID,
		req.ProductID,
	))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
