package api

import (
	"database/sql"
	"math"

	"github.com/gofiber/fiber/v2"
	product_domain "github.com/ot07/next-bazaar/api/domain/product"
	"github.com/ot07/next-bazaar/api/validation"
	"github.com/shopspring/decimal"
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
		PageID:     req.PageID,
		PageSize:   req.PageSize,
		CategoryID: req.CategoryID,
	}

	products, err := h.service.GetProducts(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	totalCount, err := h.service.CountProducts(c.Context())
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

// @Summary      List products by seller
// @Tags         Users
// @Param        query query product_domain.ListProductsBySellerRequest true "query"
// @Success      200 {object} product_domain.ListProductsResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/products [get]
func (h *productHandler) listProductsBySeller(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	req := new(product_domain.ListProductsBySellerRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	arg := product_domain.GetProductsBySellerServiceParams{
		PageID:   req.PageID,
		PageSize: req.PageSize,
		SellerID: session.UserID,
	}

	products, err := h.service.GetProductsBySeller(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	totalCount, err := h.service.CountProductsBySeller(c.Context(), session.UserID)
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

// @Summary      List product categories
// @Tags         Products
// @Param        query query product_domain.ListProductCategoriesRequest true "query"
// @Success      200 {object} product_domain.ListProductCategoriesResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /products/categories [get]
func (h *productHandler) listProductCategories(c *fiber.Ctx) error {
	req := new(product_domain.ListProductCategoriesRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	arg := product_domain.GetProductCategoriesServiceParams{
		PageID:   req.PageID,
		PageSize: req.PageSize,
	}

	categories, err := h.service.GetProductCategories(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := product_domain.ListProductCategoriesResponse{
		Meta: product_domain.ListProductCategoriesResponseMeta{
			PageID:   req.PageID,
			PageSize: req.PageSize,
		},
		Data: product_domain.NewProductCategoriesResponse(categories),
	}
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Add product
// @Tags         Users
// @Param        body body product_domain.AddProductRequest true "Product object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/products [post]
func (h *productHandler) addProduct(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	req := new(product_domain.AddProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = h.service.AddProduct(c.Context(), product_domain.AddProductServiceParams{
		Name:          req.Name,
		Description:   sql.NullString{String: req.Description, Valid: len(req.Description) > 0},
		Price:         price,
		StockQuantity: req.StockQuantity,
		CategoryID:    req.CategoryID,
		SellerID:      session.UserID,
		ImageUrl:      sql.NullString{String: req.ImageUrl, Valid: len(req.ImageUrl) > 0},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Product added successfully")
	return c.Status(fiber.StatusOK).JSON(rsp)
}

// @Summary      Update product
// @Tags         Users
// @Param        body body product_domain.UpdateProductRequestBody true "Product object"
// @Success      200 {object} messageResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /users/products/{id} [post]
func (h *productHandler) updateProduct(c *fiber.Ctx) error {
	session, err := getSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
	}

	reqParams := new(product_domain.UpdateProductRequestParams)
	if err := c.ParamsParser(reqParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	reqBody := new(product_domain.UpdateProductRequestBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := validation.NewValidator()
	if err := validate.Struct(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	price, err := decimal.NewFromString(reqBody.Price)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = h.service.UpdateProduct(c.Context(), product_domain.UpdateProductServiceParams{
		ID:            reqParams.ProductID,
		Name:          reqBody.Name,
		Description:   sql.NullString{String: reqBody.Description, Valid: len(reqBody.Description) > 0},
		Price:         price,
		StockQuantity: reqBody.StockQuantity,
		CategoryID:    reqBody.CategoryID,
		SellerID:      session.UserID,
		ImageUrl:      sql.NullString{String: reqBody.ImageUrl, Valid: len(reqBody.ImageUrl) > 0},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMessageResponse("Product updated successfully")
	return c.Status(fiber.StatusOK).JSON(rsp)
}
