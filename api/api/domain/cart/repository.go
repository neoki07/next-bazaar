package cart_domain

import (
	"context"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type FindOneByUserIDAndProductIDRepositoryParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

type CreateRepositoryParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

type UpdateRepositoryParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

type DeleteRepositoryParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

type CartRepository interface {
	FindManyByUserID(ctx context.Context, userID uuid.UUID) ([]CartProduct, error)
	FindOneByUserIDAndProductID(ctx context.Context, params FindOneByUserIDAndProductIDRepositoryParams) (CartProduct, error)
	Create(ctx context.Context, params CreateRepositoryParams) error
	Update(ctx context.Context, params UpdateRepositoryParams) error
	Delete(ctx context.Context, params DeleteRepositoryParams) error
}

func NewCartRepository(store db.Store) CartRepository {
	return &cartRepositoryImpl{
		store: store,
	}
}
