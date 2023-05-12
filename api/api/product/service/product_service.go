package service

import (
	"github.com/ot07/next-bazaar/api/product/repository"
	"github.com/ot07/next-bazaar/util"
)

type ProductService struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) GetProduct() string {
	p, _ := s.repository.FindByID(util.RandomUUID())
	return p.Name
}
