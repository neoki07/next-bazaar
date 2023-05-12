package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/product/domain"
	"github.com/ot07/next-bazaar/api/product/repository"
)

type ProductService struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := s.repository.FindByID(ctx, id)
	return product, err
}

func (s *ProductService) GetProducts(ctx context.Context, pageID int32, pageSize int32) ([]domain.Product, error) {
	products, err := s.repository.FindMany(ctx, pageID, pageSize)
	return products, err
}

func (s *ProductService) CountAllProducts(ctx context.Context) (int64, error) {
	count, err := s.repository.Count(ctx)
	return count, err
}
