package product_domain

import (
	"context"

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

type productRepositoryImpl struct {
	store db.Store
}

func (r *productRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (Product, error) {
	product, err := r.store.GetProduct(ctx, id)
	if err != nil {
		return Product{}, err
	}

	category, err := r.store.GetCategory(ctx, product.CategoryID)
	if err != nil {
		return Product{}, err
	}

	seller, err := r.store.GetUser(ctx, product.SellerID)
	if err != nil {
		return Product{}, err
	}

	return toProductDomain(product, category, seller), nil
}

func (r *productRepositoryImpl) FindMany(
	ctx context.Context,
	pageID int32,
	pageSize int32,
) ([]Product, error) {
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

	categoriesMap := make(map[uuid.UUID]db.Category)
	for _, category := range categories {
		categoriesMap[category.ID] = category
	}

	sellersIDs := productsToSellersIDs(products)
	sellers, err := r.store.GetUsersByIDs(ctx, sellersIDs)
	if err != nil {
		return nil, err
	}

	sellersMap := make(map[uuid.UUID]db.User)
	for _, seller := range sellers {
		sellersMap[seller.ID] = seller
	}

	rsp := make([]Product, len(products))
	for i, product := range products {
		rsp[i] = toProductDomain(product, categoriesMap[product.CategoryID], sellersMap[product.SellerID])
	}

	return rsp, nil
}

func (r *productRepositoryImpl) Create(ctx context.Context, product Product) error {
	_, err := r.store.CreateProduct(ctx, db.CreateProductParams{
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		CategoryID:    product.CategoryID,
		SellerID:      product.SellerID,
		ImageUrl:      product.ImageUrl,
	})

	return err
}

func (r *productRepositoryImpl) Count(ctx context.Context) (int64, error) {
	return r.store.CountProducts(ctx)
}
