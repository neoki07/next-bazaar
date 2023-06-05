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
	product, err := s.repository.FindByID(ctx, id)
	return product, err
}

func (s *ProductService) GetProducts(ctx context.Context, pageID int32, pageSize int32) ([]Product, error) {
	products, err := s.repository.FindMany(ctx, pageID, pageSize)
	return products, err
}

func (s *ProductService) CountAllProducts(ctx context.Context) (int64, error) {
	count, err := s.repository.Count(ctx)
	return count, err
}
