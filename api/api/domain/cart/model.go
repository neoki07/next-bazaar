package cart_domain

import (
	"database/sql"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/shopspring/decimal"
)

type CartProduct struct {
	ID          uuid.UUID
	Name        string
	Description sql.NullString
	Price       decimal.Decimal
	Quantity    int32
	Subtotal    decimal.Decimal
	ImageUrl    sql.NullString
}

type Cart struct {
	Products []CartProduct
	Subtotal decimal.Decimal
	Shipping decimal.Decimal
	Tax      decimal.Decimal
	Total    decimal.Decimal
}

type GetProductsRequest struct {
	ID uuid.UUID `params:"user_id"`
}

type AddProductRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int32     `json:"quantity" validate:"required,min=1"`
}

type UpdateProductQuantityRequestParams struct {
	ProductID uuid.UUID `params:"product_id"`
}

type UpdateProductQuantityRequestBody struct {
	Quantity int32 `json:"quantity" validate:"required,min=1"`
}

type DeleteProductRequest struct {
	ProductID uuid.UUID `params:"product_id"`
}

type CartProductResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description db.NullString `json:"description" swaggertype:"string"`
	Price       db.Decimal    `json:"price" swaggertype:"string"`
	Quantity    int32         `json:"quantity"`
	Subtotal    db.Decimal    `json:"subtotal" swaggertype:"string"`
	ImageUrl    db.NullString `json:"image_url" swaggertype:"string"`
}

func NewCartProductResponse(cartProduct CartProduct) CartProductResponse {
	return CartProductResponse{
		ID:          cartProduct.ID,
		Name:        cartProduct.Name,
		Description: db.NullString{NullString: cartProduct.Description},
		Price:       db.Decimal{Decimal: cartProduct.Price},
		Quantity:    cartProduct.Quantity,
		Subtotal:    db.Decimal{Decimal: cartProduct.Subtotal},
		ImageUrl:    db.NullString{NullString: cartProduct.ImageUrl},
	}
}

type CartResponse struct {
	Products []CartProductResponse `json:"products"`
	Subtotal db.Decimal            `json:"subtotal" swaggertype:"string"`
	Shipping db.Decimal            `json:"shipping" swaggertype:"string"`
	Tax      db.Decimal            `json:"tax" swaggertype:"string"`
	Total    db.Decimal            `json:"total" swaggertype:"string"`
}

func NewCartResponse(products []CartProduct) CartResponse {
	productsRsp := make([]CartProductResponse, 0, len(products))
	for _, product := range products {
		productsRsp = append(productsRsp, NewCartProductResponse(product))
	}

	subtotal := calculateSubtotal(products)
	shipping := calculateShipping()
	tax := calculateTax(subtotal)
	total := calculateTotal(subtotal, shipping, tax)

	return CartResponse{
		Products: productsRsp,
		Subtotal: db.Decimal{Decimal: subtotal},
		Shipping: db.Decimal{Decimal: shipping},
		Tax:      db.Decimal{Decimal: tax},
		Total:    db.Decimal{Decimal: total},
	}
}

type CartProductsCountResponse struct {
	Count int32 `json:"count"`
}

func NewCartProductsCountResponse(count int32) CartProductsCountResponse {
	return CartProductsCountResponse{
		Count: count,
	}
}
