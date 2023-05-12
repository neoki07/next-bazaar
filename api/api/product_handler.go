package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/product/service"
)

type productResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type getMemberRequest struct {
	ID uuid.UUID `params:"id"`
}

type productHandler struct {
	service *service.ProductService
}

func newProductHandler(s *service.ProductService) *productHandler {
	return &productHandler{
		service: s,
	}
}

// @Summary      Get product
// @Tags         products
// @Param        id path string true "Product ID"
// @Success      200 {object} memberResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /members/{id} [get]
func (h *productHandler) getProduct(c *fiber.Ctx) error {
	p := h.service.GetProduct()
	return c.Status(fiber.StatusOK).JSON(p)
}
