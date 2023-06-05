package cart_domain

import (
	"context"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/shopspring/decimal"
)

type cartRepositoryImpl struct {
	store db.Store
}

func (r *cartRepositoryImpl) FindOneByUserIDAndProductID(
	ctx context.Context,
	params FindOneByUserIDAndProductIDParams,
) (CartProduct, error) {
	arg := db.GetCartProductByUserIDAndProductIDParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
	}

	cartProduct, err := r.store.GetCartProductByUserIDAndProductID(ctx, arg)
	if err != nil {
		return CartProduct{}, err
	}

	product, err := r.store.GetProduct(ctx, cartProduct.ProductID)
	if err != nil {
		return CartProduct{}, err
	}

	price, err := decimal.NewFromString(product.Price)
	if err != nil {
		return CartProduct{}, err
	}

	quantity := decimal.NewFromInt32(cartProduct.Quantity)

	return CartProduct{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    cartProduct.Quantity,
		Subtotal:    price.Mul(quantity).String(),
	}, nil
}

func (r *cartRepositoryImpl) FindManyByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]CartProduct, error) {
	cartProducts, err := r.store.GetCartProductsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	rsp := make([]CartProduct, len(cartProducts))
	for i, cartProduct := range cartProducts {
		product, err := r.store.GetProduct(ctx, cartProduct.ProductID)
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
			Price:       product.Price,
			Quantity:    cartProduct.Quantity,
			Subtotal:    price.Mul(quantity).String(),
		}
	}

	return rsp, nil
}

func (r *cartRepositoryImpl) Create(ctx context.Context, params CreateParams) error {
	_, err := r.store.CreateCartProduct(ctx, db.CreateCartProductParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	})

	return err
}

func (r *cartRepositoryImpl) Update(ctx context.Context, params UpdateParams) error {
	_, err := r.store.UpdateCartProduct(ctx, db.UpdateCartProductParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	})

	return err
}
