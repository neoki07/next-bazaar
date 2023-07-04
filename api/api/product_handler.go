package api

import (
	"database/sql"
	"math"

	"github.com/gofiber/fiber/v2"
	product_domain "github.com/ot07/next-bazaar/api/domain/product"
	"github.com/ot07/next-bazaar/api/validation"
)

type productHandler struct {
	service *product_domain.ProductService
}

func newProductHandler(s *product_domain.ProductService) *productHandler {
	return &productHandler{
		service: s,
	}
}

// @Summary      Get product
// @Tags         Products
// @Param        id path string true "Product ID"
// @Success      200 {object} product_domain.ProductResponse
// @Failure      400 {object} errorResponse
// @Failure      404 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /products/{id} [get]
func (h *productHandler) getProduct(c *fiber.Ctx) error {
	req := new(product_domain.GetProductRequest)
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

	rsp, err := product_domain.NewProductResponse(product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      List products
// @Tags         Products
// @Param        query query product_domain.ListProductsRequest true "query"
// @Success      200 {object} product_domain.ListProductsResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /products [get]
func (h *productHandler) listProducts(c *fiber.Ctx) error {
	req := new(product_domain.ListProductsRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	arg := product_domain.GetProductsServiceParams{
		PageID:   req.PageID,
		PageSize: req.PageSize,
	}

	products, err := h.service.GetProducts(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	totalCount, err := h.service.CountAllProducts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	pageCount := int64(math.Ceil(float64(totalCount) / float64(req.PageSize)))

	rspData, err := product_domain.NewProductsResponse(products)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := product_domain.ListProductsResponse{
		Meta: product_domain.ListProductsResponseMeta{
			PageID:     req.PageID,
			PageSize:   req.PageSize,
			PageCount:  pageCount,
			TotalCount: totalCount,
		},
		Data: rspData,
	}
	return c.Status(fiber.StatusOK).JSON(rsp)
}
