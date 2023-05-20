package product_repository

import (
	"context"

	"github.com/google/uuid"
	product_domain "github.com/ot07/next-bazaar/api/domain/product"
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

type productRepositoryImpl struct {
	store db.Store
}

func (r *productRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (product_domain.Product, error) {
	product, err := r.store.GetProduct(ctx, id)
	if err != nil {
		return product_domain.Product{}, err
	}

	category, err := r.store.GetCategory(ctx, product.CategoryID)
	if err != nil {
		return product_domain.Product{}, err
	}

	seller, err := r.store.GetUser(ctx, product.SellerID)
	if err != nil {
		return product_domain.Product{}, err
	}

	rsp := product_domain.Product{
		ID:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		Category:      category.Name,
		Seller:        seller.Name,
		ImageUrl:      product.ImageUrl,
	}

	return rsp, nil
}

func (r *productRepositoryImpl) FindMany(
	ctx context.Context,
	pageID int32,
	pageSize int32,
) ([]product_domain.Product, error) {
	arg := db.ListProductsParams{
		Limit:  pageSize,
		Offset: (pageID - 1) * pageSize,
	}

	products, err := r.store.ListProducts(ctx, arg)
	if err != nil {
		return nil, err
	}

	categoryIDs := productsToCategoryIDs(products)
	categories, err := r.store.GetCategoriesByIDs(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	categoriesMap := make(map[uuid.UUID]string)
	for _, category := range categories {
		categoriesMap[category.ID] = category.Name
	}

	sellersIDs := productsToSellersIDs(products)
	sellers, err := r.store.GetUsersByIDs(ctx, sellersIDs)
	if err != nil {
		return nil, err
	}

	sellersMap := make(map[uuid.UUID]string)
	for _, seller := range sellers {
		sellersMap[seller.ID] = seller.Name
	}

	rsp := make([]product_domain.Product, len(products))
	for i, product := range products {
		rsp[i] = product_domain.Product{
			ID:            product.ID,
			Name:          product.Name,
			Description:   product.Description,
			Price:         product.Price,
			StockQuantity: product.StockQuantity,
			Category:      categoriesMap[product.CategoryID],
			Seller:        sellersMap[product.SellerID],
			ImageUrl:      product.ImageUrl,
		}
	}

	return rsp, nil
}

func (r *productRepositoryImpl) Create(ctx context.Context, product product_domain.Product) error {
	return nil
}

func (r *productRepositoryImpl) Count(ctx context.Context) (int64, error) {
	return r.store.CountProducts(ctx)
}
