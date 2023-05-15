package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/domain"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/shopspring/decimal"
)

type cartProductRepositoryImpl struct {
	store db.Store
}

func (r *cartProductRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CartProduct, error) {
	cartProducts, err := r.store.GetCartProductsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	rsp := make([]domain.CartProduct, len(cartProducts))
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

		rsp[i] = *domain.NewCartProduct(
			cartProduct.ID,
			product.Name,
			product.Description,
			product.Price,
			cartProduct.Quantity,
			price.Mul(quantity).String(),
		)
	}

	return rsp, nil
}

func (r *cartProductRepositoryImpl) Create(ctx context.Context, cartProduct *domain.CartProduct) error {
	return nil
}
