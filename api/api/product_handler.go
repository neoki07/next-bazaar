package api

import (
	"database/sql"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/domain"
	"github.com/ot07/next-bazaar/api/service"
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

type productsResponse []productResponse

func newProductsResponse(products []domain.Product) productsResponse {
	rsp := make(productsResponse, 0, len(products))
	for _, product := range products {
		rsp = append(rsp, newProductResponse(product))
	}
	return rsp
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

type listProductsRequest struct {
	PageID   int32 `query:"page_id" json:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" json:"page_size" validate:"required,min=1,max=100"`
}

type listProductsResponseMeta struct {
	PageID     int32 `json:"page_id"`
	PageSize   int32 `json:"page_size"`
	PageCount  int64 `json:"page_count"`
	TotalCount int64 `json:"total_count"`
}

type listProductsResponse struct {
	Meta listProductsResponseMeta `json:"meta"`
	Data productsResponse         `json:"data"`
}

// @Summary      List products
// @Tags         products
// @Param        query query listProductsRequest true "query"
// @Success      200 {object} listProductsResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /products [get]
func (h *productHandler) listProducts(c *fiber.Ctx) error {
	req := new(listProductsRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := newValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	products, err := h.service.GetProducts(c.Context(), req.PageID, req.PageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	totalCount, err := h.service.CountAllProducts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	pageCount := int64(math.Ceil(float64(totalCount) / float64(req.PageSize)))

	rsp := listProductsResponse{
		Meta: listProductsResponseMeta{
			PageID:     req.PageID,
			PageSize:   req.PageSize,
			PageCount:  pageCount,
			TotalCount: totalCount,
		},
		Data: newProductsResponse(products),
	}
	return c.Status(fiber.StatusOK).JSON(rsp)
}
