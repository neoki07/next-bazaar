package product_domain

import (
	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

func productsToCategoryIDs(products []db.Product) []uuid.UUID {
	categoryIDs := make([]uuid.UUID, len(products))
	for i, product := range products {
		categoryIDs[i] = product.CategoryID
	}

	return categoryIDs
}

func productsToSellersIDs(products []db.Product) []uuid.UUID {
	sellersIDs := make([]uuid.UUID, len(products))
	for i, product := range products {
		sellersIDs[i] = product.SellerID
	}

	return sellersIDs
}

func toProductDomain(product db.Product, category db.Category, seller db.User) Product {
	return Product{
		ID:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		CategoryID:    category.ID,
		Category:      category.Name,
		SellerID:      seller.ID,
		Seller:        seller.Name,
		ImageUrl:      product.ImageUrl,
	}
}

func toCategoryDomain(category db.Category) Category {
	return Category{
		ID:   category.ID,
		Name: category.Name,
	}
}
