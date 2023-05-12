package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/product/domain"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
}

func NewProductRepository(store db.Store) ProductRepository {
	return &productRepositoryImpl{
		store: store,
	}
}
