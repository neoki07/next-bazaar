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
