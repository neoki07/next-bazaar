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
