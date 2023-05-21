package db

import (
	"context"
	"database/sql"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/ot07/next-bazaar/util"
)

func CreateProductTestData(ctx context.Context, store *SQLStore, config util.Config) error {
	user, err := store.GetUserByEmail(ctx, config.TestAccountEmail)
	if err != nil {
		return err
	}

	numProducts := 100

	arg := ListCategoriesParams{
		Offset: 0,
		Limit:  int32(numProducts),
	}

	categories, err := store.ListCategories(ctx, arg)
	if err != nil {
		return err
	}

	for i := 0; i < numProducts; i++ {
		price, err := util.RandomPrice()
		if err != nil {
			return err
		}

		arg := CreateProductParams{
			Name:          faker.Name(),
			Description:   sql.NullString{String: faker.Paragraph(), Valid: true},
			Price:         price,
			StockQuantity: rand.Int31n(100),
			CategoryID:    categories[rand.Intn(len(categories))].ID,
			SellerID:      user.ID,
			ImageUrl:      sql.NullString{String: util.RandomImageUrl(), Valid: true},
		}

		_, err = store.CreateProduct(ctx, arg)
		if err != nil {
			return err
		}
	}

	return nil
}
