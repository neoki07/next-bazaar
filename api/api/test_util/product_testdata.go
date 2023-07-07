package test_util

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func CreateCategoryTestData(t *testing.T, ctx context.Context, store db.Store, name string) db.Category {
	category, err := store.CreateCategory(ctx, name)
	require.NoError(t, err)

	return category
}

func CreateProductTestData(
	t *testing.T,
	ctx context.Context,
	store db.Store,
	name string,
	description string,
	price decimal.Decimal,
	stockQuantity int32,
	categoryID uuid.UUID,
	sellerID uuid.UUID,
	imageUrl string,
) db.Product {
	product, err := store.CreateProduct(ctx, db.CreateProductParams{
		Name:          name,
		Description:   sql.NullString{String: description, Valid: true},
		Price:         price.String(),
		StockQuantity: stockQuantity,
		CategoryID:    categoryID,
		SellerID:      sellerID,
		ImageUrl:      sql.NullString{String: imageUrl, Valid: true},
	})
	require.NoError(t, err)

	return product
}
