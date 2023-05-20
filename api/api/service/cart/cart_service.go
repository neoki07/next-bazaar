package cart_service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	cart_domain "github.com/ot07/next-bazaar/api/domain/cart"
	cart_repository "github.com/ot07/next-bazaar/api/repository/cart"
)

type CartService struct {
	repository cart_repository.CartRepository
}

func NewCartService(repository cart_repository.CartRepository) *CartService {
	return &CartService{
		repository: repository,
	}
}

func (s *CartService) GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]cart_domain.CartProduct, error) {
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
	arg := cart_repository.NewFindOneByUserIDAndProductIDParams(
		params.UserID,
		params.ProductID,
	)

	cartProduct, err := s.repository.FindOneByUserIDAndProductID(ctx, arg)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err != nil && err == sql.ErrNoRows {
		return s.repository.Create(ctx, cart_repository.NewCreateParams(
			params.UserID,
			params.ProductID,
			params.Quantity,
		))
	}

	return s.repository.Update(ctx, cart_repository.NewUpdateParams(
		params.UserID,
		params.ProductID,
		params.Quantity+cartProduct.Quantity,
	))
}
