package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/product/domain"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	FindMany(ctx context.Context, pageID int32, pageSize int32) ([]domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
	Count(ctx context.Context) (int64, error)
}

func NewProductRepository(store db.Store) ProductRepository {
	return &productRepositoryImpl{
		store: store,
	}
}
