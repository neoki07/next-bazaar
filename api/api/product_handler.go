package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/product/domain"
	"github.com/ot07/next-bazaar/api/product/service"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type productResponse struct {
	ID            uuid.UUID     `json:"id"`
	Name          string        `json:"name"`
	Description   db.NullString `json:"description" swaggertype:"string"`
	Price         string        `json:"price"`
	StockQuantity int32         `json:"stock_quantity"`
	Category      string        `json:"category"`
	Seller        string        `json:"seller"`
	ImageUrl      db.NullString `json:"image_url" swaggertype:"string"`
}

func newProductResponse(product domain.Product) productResponse {
	return productResponse{
		ID:            product.ID,
		Name:          product.Name,
		Description:   db.NullString{NullString: product.Description},
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		Category:      product.Category,
		Seller:        product.Seller,
		ImageUrl:      db.NullString{NullString: product.ImageUrl},
	}
}

type getProductRequest struct {
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
// @Success      200 {object} productResponse
// @Failure      400 {object} errorResponse
// @Failure      404 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /products/{id} [get]
func (h *productHandler) getProduct(c *fiber.Ctx) error {
	req := new(getProductRequest)
	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	product, err := h.service.GetProduct(c.Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newProductResponse(*product)
	return c.Status(fiber.StatusOK).JSON(rsp)
}
