package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/domain"
	"github.com/ot07/next-bazaar/api/repository"
)

type CartProductService struct {
	repository repository.CartProductRepository
}

func NewCartProductService(repository repository.CartProductRepository) *CartProductService {
	return &CartProductService{
		repository: repository,
	}
}

func (s *CartProductService) GetCartProductsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CartProduct, error) {
	cartProducts, err := s.repository.FindManyByUserID(ctx, userID)
	return cartProducts, err
}

type AddProductToCartParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func NewAddProductToCartParams(
	userID uuid.UUID,
	productID uuid.UUID,
	quantity int32,
) AddProductToCartParams {
	return AddProductToCartParams{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
}

func (s *CartProductService) AddProductToCart(ctx context.Context, params AddProductToCartParams) error {
	arg := repository.NewFindOneByUserIDAndProductIDParams(
		params.UserID,
		params.ProductID,
	)

	cartProduct, err := s.repository.FindOneByUserIDAndProductID(ctx, arg)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err != nil && err == sql.ErrNoRows {
		return s.repository.Create(ctx, repository.NewCreateParams(
			params.UserID,
			params.ProductID,
			params.Quantity,
		))
	}

	return s.repository.Update(ctx, repository.NewUpdateParams(
		params.UserID,
		params.ProductID,
		params.Quantity+cartProduct.Quantity,
	))
}
