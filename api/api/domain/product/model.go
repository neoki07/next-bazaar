package product_domain

import (
	"database/sql"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID            uuid.UUID
	Name          string
	Description   sql.NullString
	Price         string
	StockQuantity int32
	CategoryID    uuid.UUID
	Category      string
	SellerID      uuid.UUID
	Seller        string
	ImageUrl      sql.NullString
}

type Category struct {
	ID   uuid.UUID
	Name string
}

type GetProductRequest struct {
	ID uuid.UUID `params:"id"`
}

type ListProductsRequest struct {
	PageID     int32         `query:"page_id" json:"page_id" validate:"required,min=1"`
	PageSize   int32         `query:"page_size" json:"page_size" validate:"required,min=1,max=100"`
	CategoryID uuid.NullUUID `query:"category_id" json:"category_id" swaggertype:"string"`
}

type ListProductsBySellerRequest struct {
	PageID   int32 `query:"page_id" json:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" json:"page_size" validate:"required,min=1,max=100"`
}

type ProductResponse struct {
	ID            uuid.UUID     `json:"id"`
	Name          string        `json:"name"`
	Description   db.NullString `json:"description" swaggertype:"string"`
	Price         db.Decimal    `json:"price" swaggertype:"string"`
	StockQuantity int32         `json:"stock_quantity"`
	Category      string        `json:"category"`
	Seller        string        `json:"seller"`
	ImageUrl      db.NullString `json:"image_url" swaggertype:"string"`
}

func NewProductResponse(product Product) (ProductResponse, error) {
	dec, err := decimal.NewFromString(product.Price)
	if err != nil {
		return ProductResponse{}, err
	}

	return ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Description:   db.NullString{NullString: product.Description},
		Price:         db.Decimal{Decimal: dec},
		StockQuantity: product.StockQuantity,
		Category:      product.Category,
		Seller:        product.Seller,
		ImageUrl:      db.NullString{NullString: product.ImageUrl},
	}, nil
}

type ProductsResponse []ProductResponse

func NewProductsResponse(products []Product) (ProductsResponse, error) {
	rsp := make(ProductsResponse, 0, len(products))

	for _, product := range products {
		item, err := NewProductResponse(product)
		if err != nil {
			return nil, err
		}

		rsp = append(rsp, item)
	}

	return rsp, nil
}

type ListProductsResponseMeta struct {
	PageID     int32 `json:"page_id"`
	PageSize   int32 `json:"page_size"`
	PageCount  int64 `json:"page_count"`
	TotalCount int64 `json:"total_count"`
}

type ListProductsResponse struct {
	Meta ListProductsResponseMeta `json:"meta"`
	Data ProductsResponse         `json:"data"`
}

type ListProductCategoriesRequest struct {
	PageID   int32 `query:"page_id" json:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" json:"page_size" validate:"required,min=1,max=100"`
}

type ListProductCategoriesResponseMeta struct {
	PageID   int32 `json:"page_id"`
	PageSize int32 `json:"page_size"`
}

type ProductCategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProductCategoriesResponse []ProductCategoryResponse

func NewProductCategoriesResponse(categories []Category) ProductCategoriesResponse {
	rsp := make(ProductCategoriesResponse, 0, len(categories))

	for _, category := range categories {
		rsp = append(rsp, ProductCategoryResponse(category))
	}

	return rsp
}

type ListProductCategoriesResponse struct {
	Meta ListProductCategoriesResponseMeta `json:"meta"`
	Data ProductCategoriesResponse         `json:"data"`
}
