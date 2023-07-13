package cart_domain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/shopspring/decimal"
)

type CartService struct {
	store db.Store
}

func NewCartService(store db.Store) *CartService {
	return &CartService{
		store: store,
	}
}

func (s *CartService) GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]CartProduct, error) {
	cartProducts, err := s.store.GetCartProductsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	rsp := make([]CartProduct, len(cartProducts))
	for i, cartProduct := range cartProducts {
		product, err := s.store.GetProduct(ctx, cartProduct.ProductID)
		if err != nil {
			return nil, err
		}

		price, err := decimal.NewFromString(product.Price)
		if err != nil {
			return nil, err
		}

		quantity := decimal.NewFromInt32(cartProduct.Quantity)

		rsp[i] = CartProduct{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       price,
			Quantity:    cartProduct.Quantity,
			Subtotal:    price.Mul(quantity),
			ImageUrl:    product.ImageUrl,
		}
	}

	return rsp, nil
}

type createServiceParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func (s *CartService) createProduct(ctx context.Context, params createServiceParams) error {
	_, err := s.store.CreateCartProduct(ctx, db.CreateCartProductParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	})

	return err
}

type updateServiceParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func (s *CartService) updateProduct(ctx context.Context, params updateServiceParams) error {
	_, err := s.store.UpdateCartProduct(ctx, db.UpdateCartProductParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	})

	return err
}

type AddProductServiceParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func (s *CartService) AddProduct(ctx context.Context, params AddProductServiceParams) error {
	arg := db.GetCartProductByUserIDAndProductIDParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
	}

	cartProduct, err := s.store.GetCartProductByUserIDAndProductID(ctx, arg)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err != nil && err == sql.ErrNoRows {
		return s.createProduct(ctx, createServiceParams(params))
	}

	return s.updateProduct(ctx, updateServiceParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity + cartProduct.Quantity,
	})
}

type UpdateProductQuantityServiceParams = updateServiceParams

func (s *CartService) UpdateProductQuantity(ctx context.Context, params UpdateProductQuantityServiceParams) error {
	return s.updateProduct(ctx, updateServiceParams(params))
}

type DeleteProductServiceParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

func (s *CartService) DeleteProduct(ctx context.Context, params DeleteProductServiceParams) error {
	return s.store.DeleteCartProduct(ctx, db.DeleteCartProductParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
	})
}
