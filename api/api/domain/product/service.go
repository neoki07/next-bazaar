package product_domain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/shopspring/decimal"
)

type ProductService struct {
	store db.Store
}

func NewProductService(store db.Store) *ProductService {
	return &ProductService{
		store: store,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (Product, error) {
	product, err := s.store.GetProduct(ctx, id)
	if err != nil {
		return Product{}, err
	}

	category, err := s.store.GetCategory(ctx, product.CategoryID)
	if err != nil {
		return Product{}, err
	}

	seller, err := s.store.GetUser(ctx, product.SellerID)
	if err != nil {
		return Product{}, err
	}

	return toProductDomain(product, category, seller), nil
}

type GetProductsServiceParams struct {
	PageID     int32
	PageSize   int32
	CategoryID uuid.NullUUID
}

func (s *ProductService) GetProducts(ctx context.Context, params GetProductsServiceParams) ([]Product, error) {
	arg := db.ListProductsParams{
		Limit:      params.PageSize,
		Offset:     (params.PageID - 1) * params.PageSize,
		CategoryID: params.CategoryID,
	}

	products, err := s.store.ListProducts(ctx, arg)
	if err != nil {
		return nil, err
	}

	categoryIDs := productsToCategoryIDs(products)
	categories, err := s.store.GetCategoriesByIDs(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	categoriesMap := make(map[uuid.UUID]db.Category)
	for _, category := range categories {
		categoriesMap[category.ID] = category
	}

	sellersIDs := productsToSellersIDs(products)
	sellers, err := s.store.GetUsersByIDs(ctx, sellersIDs)
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

func (s *ProductService) CountProducts(ctx context.Context) (int64, error) {
	return s.store.CountProducts(ctx)
}

type GetProductsBySellerServiceParams struct {
	PageID   int32
	PageSize int32
	SellerID uuid.UUID
}

func (s *ProductService) GetProductsBySeller(ctx context.Context, params GetProductsBySellerServiceParams) ([]Product, error) {
	arg := db.ListProductsBySellerParams{
		Limit:    params.PageSize,
		Offset:   (params.PageID - 1) * params.PageSize,
		SellerID: params.SellerID,
	}

	products, err := s.store.ListProductsBySeller(ctx, arg)
	if err != nil {
		return nil, err
	}

	categoryIDs := productsToCategoryIDs(products)
	categories, err := s.store.GetCategoriesByIDs(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	categoriesMap := make(map[uuid.UUID]db.Category)
	for _, category := range categories {
		categoriesMap[category.ID] = category
	}

	sellersIDs := productsToSellersIDs(products)
	sellers, err := s.store.GetUsersByIDs(ctx, sellersIDs)
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

func (s *ProductService) CountProductsBySeller(ctx context.Context, sellerID uuid.UUID) (int64, error) {
	return s.store.CountProductsBySeller(ctx, sellerID)
}

type GetProductCategoriesServiceParams struct {
	PageID   int32
	PageSize int32
}

func (s *ProductService) GetProductCategories(ctx context.Context, params GetProductCategoriesServiceParams) ([]Category, error) {
	arg := db.ListCategoriesParams{
		Limit:  params.PageSize,
		Offset: (params.PageID - 1) * params.PageSize,
	}

	categories, err := s.store.ListCategories(ctx, arg)
	if err != nil {
		return nil, err
	}

	rsp := make([]Category, len(categories))
	for i, category := range categories {
		rsp[i] = toCategoryDomain(category)
	}

	return rsp, nil
}

type AddProductServiceParams struct {
	Name          string
	Description   sql.NullString
	Price         decimal.Decimal
	StockQuantity int32
	CategoryID    uuid.UUID
	SellerID      uuid.UUID
	ImageUrl      sql.NullString
}

func (s *ProductService) AddProduct(ctx context.Context, params AddProductServiceParams) error {
	_, err := s.store.AddProduct(ctx, db.AddProductParams{
		Name:          params.Name,
		Description:   params.Description,
		Price:         params.Price.String(),
		StockQuantity: params.StockQuantity,
		CategoryID:    params.CategoryID,
		SellerID:      params.SellerID,
		ImageUrl:      params.ImageUrl,
	})

	return err
}

type UpdateProductServiceParams struct {
	ID            uuid.UUID
	Name          string
	Description   sql.NullString
	Price         decimal.Decimal
	StockQuantity int32
	CategoryID    uuid.UUID
	SellerID      uuid.UUID
	ImageUrl      sql.NullString
}

func (s *ProductService) UpdateProduct(ctx context.Context, params UpdateProductServiceParams) error {
	_, err := s.store.UpdateProduct(ctx, db.UpdateProductParams{
		ID:            params.ID,
		Name:          params.Name,
		Description:   params.Description,
		Price:         params.Price.String(),
		StockQuantity: params.StockQuantity,
		CategoryID:    params.CategoryID,
		SellerID:      params.SellerID,
		ImageUrl:      params.ImageUrl,
	})

	return err
}
