package cart_repository

import (
	"context"

	"github.com/google/uuid"
	cart_domain "github.com/ot07/next-bazaar/api/domain/cart"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type FindOneByUserIDAndProductIDParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

func NewFindOneByUserIDAndProductIDParams(
	userID uuid.UUID,
	productID uuid.UUID,
) FindOneByUserIDAndProductIDParams {
	return FindOneByUserIDAndProductIDParams{
		UserID:    userID,
		ProductID: productID,
	}
}

type CreateParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func NewCreateParams(
	userID uuid.UUID,
	productID uuid.UUID,
	quantity int32,
) CreateParams {
	return CreateParams{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
}

type UpdateParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func NewUpdateParams(
	userID uuid.UUID,
	productID uuid.UUID,
	quantity int32,
) UpdateParams {
	return UpdateParams{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
}

type CartRepository interface {
	FindManyByUserID(ctx context.Context, userID uuid.UUID) ([]cart_domain.CartProduct, error)
	FindOneByUserIDAndProductID(ctx context.Context, params FindOneByUserIDAndProductIDParams) (*cart_domain.CartProduct, error)
	Create(ctx context.Context, params CreateParams) error
	Update(ctx context.Context, params UpdateParams) error
}

func NewCartRepository(store db.Store) CartRepository {
	return &cartRepositoryImpl{
		store: store,
	}
}
