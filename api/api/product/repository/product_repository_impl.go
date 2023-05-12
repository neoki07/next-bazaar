package repository

import (
	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/product/domain"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type productRepositoryImpl struct {
	store db.Store
}

func (r *productRepositoryImpl) FindByID(id uuid.UUID) (*domain.Product, error) {
	product := domain.NewProduct(id, "product")
	return product, nil
}

func (r *productRepositoryImpl) Create(product *domain.Product) error {
	return nil
}
