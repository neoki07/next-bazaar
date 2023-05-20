package product_domain

import "github.com/google/uuid"

type GetProductRequest struct {
	ID uuid.UUID `params:"id"`
}

type ListProductsRequest struct {
	PageID   int32 `query:"page_id" json:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" json:"page_size" validate:"required,min=1,max=100"`
}
