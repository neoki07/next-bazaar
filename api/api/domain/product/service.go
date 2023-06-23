package product_domain

import (
	"context"

	"github.com/google/uuid"
)

type ProductService struct {
	repository ProductRepository
}

func NewProductService(repository ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (Product, error) {
	return s.repository.FindByID(ctx, id)
}

type GetProductsServiceParams struct {
	PageID   int32
	PageSize int32
}

func (s *ProductService) GetProducts(ctx context.Context, params GetProductsServiceParams) ([]Product, error) {
	return s.repository.FindMany(ctx, FindManyParams(params))
}

func (s *ProductService) CountAllProducts(ctx context.Context) (int64, error) {
	return s.repository.Count(ctx)
}
