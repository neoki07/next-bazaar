package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/product/domain"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type productRepositoryImpl struct {
	store db.Store
}

func (r *productRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := r.store.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	category, err := r.store.GetCategory(ctx, product.CategoryID)
	if err != nil {
		return nil, err
	}

	seller, err := r.store.GetUser(ctx, product.SellerID)
	if err != nil {
		return nil, err
	}

	rsp := domain.NewProduct(
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		category.Name,
		seller.Name,
		product.ImageUrl,
	)

	return rsp, nil
}

func (r *productRepositoryImpl) Create(ctx context.Context, product *domain.Product) error {
	return nil
}
