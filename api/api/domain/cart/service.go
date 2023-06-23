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
	return s.repository.FindManyByUserID(ctx, userID)
}

type AddProductServiceParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func (s *CartService) AddProduct(ctx context.Context, params AddProductServiceParams) error {
	arg := FindOneByUserIDAndProductIDRepositoryParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
	}

	cartProduct, err := s.repository.FindOneByUserIDAndProductID(ctx, arg)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err != nil && err == sql.ErrNoRows {
		return s.repository.Create(ctx, CreateRepositoryParams(params))
	}

	return s.repository.Update(ctx, UpdateRepositoryParams{
		UserID:    params.UserID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity + cartProduct.Quantity,
	})
}

type UpdateProductQuantityServiceParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func (s *CartService) UpdateProductQuantity(ctx context.Context, params UpdateProductQuantityServiceParams) error {
	return s.repository.Update(ctx, UpdateRepositoryParams(params))
}

type DeleteProductServiceParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

func (s *CartService) DeleteProduct(ctx context.Context, params DeleteProductServiceParams) error {
	return s.repository.Delete(ctx, DeleteRepositoryParams(params))
}
