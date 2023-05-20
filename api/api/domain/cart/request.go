package cart_domain

import "github.com/google/uuid"

type GetProductsRequest struct {
	ID uuid.UUID `params:"user_id"`
}

type AddProductRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int32     `json:"quantity" validate:"required,min=1"`
}
