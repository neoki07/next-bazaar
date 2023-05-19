package product_repository

import (
	"context"

	"github.com/google/uuid"
	product_domain "github.com/ot07/next-bazaar/api/domain/product"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*product_domain.Product, error)
	FindMany(ctx context.Context, pageID int32, pageSize int32) ([]product_domain.Product, error)
	Create(ctx context.Context, product *product_domain.Product) error
	Count(ctx context.Context) (int64, error)
}

func NewProductRepository(store db.Store) ProductRepository {
	return &productRepositoryImpl{
		store: store,
	}
}
