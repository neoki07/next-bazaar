package api

import (
	"database/sql"

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
// @Router       /cart [get]
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

// @Summary      Get cart products count
// @Tags         Cart
// @Success      200 {object} cart_domain.CartProductsCountResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart/count [get]
func (h *cartHandler) getCartProductsCount(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	cartProducts, err := h.service.GetProductsByUserID(c.Context(), session.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	var cartProductsCount int32
	for _, cartProduct := range cartProducts {
		cartProductsCount += cartProduct.Quantity
	}

	rsp := cart_domain.NewCartProductsCountResponse(cartProductsCount)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Add product to cart
// @Tags         Cart
// @Param        body body cart_domain.AddProductRequest true "Cart product object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart/add-product [post]
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

	err = h.service.AddProduct(c.Context(), cart_domain.AddProductServiceParams{
		UserID:    session.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Cart product added successfully")
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Update cart product quantity
// @Tags         Cart
// @Param        product_id path string true "Product ID"
// @Param        body body cart_domain.UpdateProductQuantityRequestBody true "Cart product object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart/{product_id} [put]
func (h *cartHandler) updateProductQuantity(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	reqParams := new(cart_domain.UpdateProductQuantityRequestParams)
	if err := c.ParamsParser(reqParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	reqBody := new(cart_domain.UpdateProductQuantityRequestBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = h.service.UpdateProductQuantity(c.Context(), cart_domain.UpdateProductQuantityServiceParams{
		UserID:    session.UserID,
		ProductID: reqParams.ProductID,
		Quantity:  reqBody.Quantity,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Cart product added successfully")
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Delete cart product
// @Tags         Cart
// @Param        product_id path string true "Product ID"
// @Success      204
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart/{product_id} [delete]
func (h *cartHandler) deleteProduct(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	req := new(cart_domain.DeleteProductRequest)
	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = h.service.DeleteProduct(c.Context(), cart_domain.DeleteProductServiceParams{
		UserID:    session.UserID,
		ProductID: req.ProductID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
