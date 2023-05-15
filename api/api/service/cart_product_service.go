package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/api/domain"
	"github.com/ot07/next-bazaar/api/repository"
)

type CartProductService struct {
	repository repository.CartProductRepository
}

func NewCartProductService(repository repository.CartProductRepository) *CartProductService {
	return &CartProductService{
		repository: repository,
	}
}

func (s *CartProductService) GetCartProductsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CartProduct, error) {
	cartProducts, err := s.repository.FindByUserID(ctx, userID)
	return cartProducts, err
}
