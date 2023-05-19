package product_domain

import (
	"database/sql"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID
	Name          string
	Description   sql.NullString
	Price         string
	StockQuantity int32
	Category      string
	Seller        string
	ImageUrl      sql.NullString
}

func NewProduct(
	id uuid.UUID,
	name string,
	description sql.NullString,
	price string,
	stockQuantity int32,
	category string,
	seller string,
	imageUrl sql.NullString,
) *Product {
	return &Product{
		ID:            id,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      category,
		Seller:        seller,
		ImageUrl:      imageUrl,
	}
}
