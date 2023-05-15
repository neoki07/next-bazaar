package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
// @Router       /cart-products/{user-id} [get]
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
