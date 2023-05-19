package product_service

import (
	"context"

	"github.com/google/uuid"
	product_domain "github.com/ot07/next-bazaar/api/domain/product"
	product_repository "github.com/ot07/next-bazaar/api/repository/product"
)

type ProductService struct {
	repository product_repository.ProductRepository
}

func NewProductService(repository product_repository.ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*product_domain.Product, error) {
	product, err := s.repository.FindByID(ctx, id)
	return product, err
}

func (s *ProductService) GetProducts(ctx context.Context, pageID int32, pageSize int32) ([]product_domain.Product, error) {
	products, err := s.repository.FindMany(ctx, pageID, pageSize)
	return products, err
}

func (s *ProductService) CountAllProducts(ctx context.Context) (int64, error) {
	count, err := s.repository.Count(ctx)
	return count, err
}
