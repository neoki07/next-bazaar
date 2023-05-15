package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/domain"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type CartProductRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CartProduct, error)
	Create(ctx context.Context, cartProduct *domain.CartProduct) error
}

func NewCartProductRepository(store db.Store) CartProductRepository {
	return &cartProductRepositoryImpl{
		store: store,
	}
}
