package repository

import (
	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/product/domain"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type ProductRepository interface {
	FindByID(id uuid.UUID) (*domain.Product, error)
	Create(product *domain.Product) error
}

func NewProductRepository(store db.Store) ProductRepository {
	return &productRepositoryImpl{
		store: store,
	}
}
