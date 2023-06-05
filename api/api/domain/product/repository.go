package product_domain

import (
	"context"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (Product, error)
	FindMany(ctx context.Context, pageID int32, pageSize int32) ([]Product, error)
	Create(ctx context.Context, product Product) error
	Count(ctx context.Context) (int64, error)
}

func NewProductRepository(store db.Store) ProductRepository {
	return &productRepositoryImpl{
		store: store,
	}
}
