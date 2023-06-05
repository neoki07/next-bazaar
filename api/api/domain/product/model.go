package product_domain

import (
	"database/sql"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
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
	PageID   int32 `query:"page_id" json:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" json:"page_size" validate:"required,min=1,max=100"`
}

type ProductResponse struct {
	ID            uuid.UUID     `json:"id"`
	Name          string        `json:"name"`
	Description   db.NullString `json:"description" swaggertype:"string"`
	Price         string        `json:"price"`
	StockQuantity int32         `json:"stock_quantity"`
	Category      string        `json:"category"`
	Seller        string        `json:"seller"`
	ImageUrl      db.NullString `json:"image_url" swaggertype:"string"`
}

func NewProductResponse(product Product) ProductResponse {
	return ProductResponse{
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

type ProductsResponse []ProductResponse

func NewProductsResponse(products []Product) ProductsResponse {
	rsp := make(ProductsResponse, 0, len(products))
	for _, product := range products {
		rsp = append(rsp, NewProductResponse(product))
	}
	return rsp
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
