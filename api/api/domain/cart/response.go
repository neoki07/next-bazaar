package cart_domain

import (
	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type CartProductResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description db.NullString `json:"description" swaggertype:"string"`
	Price       string        `json:"price"`
	Quantity    int32         `json:"quantity"`
	Subtotal    string        `json:"subtotal"`
}

func NewCartProductResponse(cartProduct CartProduct) CartProductResponse {
	return CartProductResponse{
		ID:          cartProduct.ID,
		Name:        cartProduct.Name,
		Description: db.NullString{NullString: cartProduct.Description},
		Price:       cartProduct.Price,
		Quantity:    cartProduct.Quantity,
		Subtotal:    cartProduct.Subtotal,
	}
}

type CartResponse []CartProductResponse

func NewCartResponse(products []CartProduct) CartResponse {
	rsp := make(CartResponse, 0, len(products))
	for _, product := range products {
		rsp = append(rsp, NewCartProductResponse(product))
	}
	return rsp
}
