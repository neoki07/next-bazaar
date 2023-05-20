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
