package cart_domain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type CartService struct {
	repository CartRepository
}

func NewCartService(repository CartRepository) *CartService {
	return &CartService{
		repository: repository,
	}
}

func (s *CartService) GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]CartProduct, error) {
	cartProducts, err := s.repository.FindManyByUserID(ctx, userID)
	return cartProducts, err
}

type AddProductParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func NewAddProductParams(
	userID uuid.UUID,
	productID uuid.UUID,
	quantity int32,
) AddProductParams {
	return AddProductParams{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
}

func (s *CartService) AddProduct(ctx context.Context, params AddProductParams) error {
	arg := NewFindOneByUserIDAndProductIDParams(
		params.UserID,
		params.ProductID,
	)

	cartProduct, err := s.repository.FindOneByUserIDAndProductID(ctx, arg)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err != nil && err == sql.ErrNoRows {
		return s.repository.Create(ctx, NewCreateParams(
			params.UserID,
			params.ProductID,
			params.Quantity,
		))
	}

	return s.repository.Update(ctx, NewUpdateParams(
		params.UserID,
		params.ProductID,
		params.Quantity+cartProduct.Quantity,
	))
}

type UpdateProductQuantityParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func NewUpdateProductQuantityParams(
	userID uuid.UUID,
	productID uuid.UUID,
	quantity int32,
) UpdateProductQuantityParams {
	return UpdateProductQuantityParams{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
}

func (s *CartService) UpdateProductQuantity(ctx context.Context, params UpdateProductQuantityParams) error {
	return s.repository.Update(ctx, NewUpdateParams(
		params.UserID,
		params.ProductID,
		params.Quantity,
	))
}
