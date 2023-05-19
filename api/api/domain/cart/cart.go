package cart_domain

import (
	"database/sql"

	"github.com/google/uuid"
)

type CartProduct struct {
	ID          uuid.UUID
	Name        string
	Description sql.NullString
	Price       string
	Quantity    int32
	Subtotal    string
}

func NewCartProduct(
	id uuid.UUID,
	name string,
	description sql.NullString,
	price string,
	quantity int32,
	subtotal string,
) *CartProduct {
	return &CartProduct{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
	}
}
