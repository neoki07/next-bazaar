package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/domain"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/shopspring/decimal"
)

type cartRepositoryImpl struct {
	store db.Store
}

func (r *cartRepositoryImpl) FindOneByUserIDAndProductID(
	ctx context.Context,
	params FindOneByUserIDAndProductIDParams,
) (*domain.CartProduct, error) {
	arg := db.GetCartProductByUserIdAndProductIdParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
	}

	cartProduct, err := r.store.GetCartProductByUserIdAndProductId(ctx, arg)
	if err != nil {
		return nil, err
	}

	product, err := r.store.GetProduct(ctx, cartProduct.ProductID)
	if err != nil {
		return nil, err
	}

	price, err := decimal.NewFromString(product.Price)
	if err != nil {
		return nil, err
	}

	quantity := decimal.NewFromInt32(cartProduct.Quantity)

	return domain.NewCartProduct(
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		cartProduct.Quantity,
		price.Mul(quantity).String(),
	), nil
}

func (r *cartRepositoryImpl) FindManyByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.CartProduct, error) {
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
			product.ID,
			product.Name,
			product.Description,
			product.Price,
			cartProduct.Quantity,
			price.Mul(quantity).String(),
		)
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
