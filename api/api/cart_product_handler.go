package api

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ot07/next-bazaar/api/domain"
	"github.com/ot07/next-bazaar/api/service"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type cartProductResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description db.NullString `json:"description" swaggertype:"string"`
	Price       string        `json:"price"`
	Quantity    int32         `json:"quantity"`
	Subtotal    string        `json:"subtotal"`
}

func newCartProductResponse(cartProduct domain.CartProduct) cartProductResponse {
	return cartProductResponse{
		ID:          cartProduct.ID,
		Name:        cartProduct.Name,
		Description: db.NullString{NullString: cartProduct.Description},
		Price:       cartProduct.Price,
		Quantity:    cartProduct.Quantity,
		Subtotal:    cartProduct.Subtotal,
	}
}

type cartProductsResponse []cartProductResponse

func newCartProductsResponse(products []domain.CartProduct) cartProductsResponse {
	rsp := make(cartProductsResponse, 0, len(products))
	for _, product := range products {
		rsp = append(rsp, newCartProductResponse(product))
	}
	return rsp
}

type getCartProductsRequest struct {
	ID uuid.UUID `params:"user_id"`
}

type cartProductHandler struct {
	service *service.CartProductService
}

func newCartProductHandler(s *service.CartProductService) *cartProductHandler {
	return &cartProductHandler{
		service: s,
	}
}

// @Summary      Get cart products
// @Tags         cartProducts
// @Param        userId path string true "User ID"
// @Success      200 {object} productResponse
// @Failure      400 {object} errorResponse
// @Failure      404 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart-products/{userId} [get]
func (h *cartProductHandler) getCartProducts(c *fiber.Ctx) error {
	req := new(getCartProductsRequest)
	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	cartProducts, err := h.service.GetCartProductsByUserID(c.Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newCartProductsResponse(cartProducts)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

type addProductToCartRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int32     `json:"quantity" validate:"required,min=1"`
}

// @Summary      Add product to cart
// @Tags         cartProducts
// @Param        body body addProductToCartRequest true "Cart product object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      403 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /cart-products [post]
func (h *cartProductHandler) addProductToCart(c *fiber.Ctx) error {
	session, ok := c.Locals(ctxLocalSessionKey).(db.Session)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(
			fmt.Errorf("session token not found"),
		))
	}

	req := new(addProductToCartRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := newValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err := h.service.AddProductToCart(c.Context(), service.NewAddProductToCartParams(
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

	rsp := messageResponse{
		Message: "Cart product created successfully",
	}

	return c.Status(fiber.StatusOK).JSON(rsp)
}
