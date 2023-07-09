package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/go-faker/faker/v4"
	"github.com/ot07/next-bazaar/util"
)

const publicDir = "../web/public"

var (
	testImageRootDir       = filepath.Join(publicDir, "testdata/product-images")
	categoryToTestImageDir = map[string]string{
		"Jeans":   filepath.Join(testImageRootDir, "jeans"),
		"Sofa":    filepath.Join(testImageRootDir, "sofa"),
		"T-Shirt": filepath.Join(testImageRootDir, "tshirt"),
		"TV":      filepath.Join(testImageRootDir, "tv"),
	}
)

func CreateProductTestData(ctx context.Context, store *SQLStore, config util.Config) error {
	user, err := store.GetUserByEmail(ctx, config.TestAccountEmail)
	if err != nil {
		return err
	}

	arg := ListCategoriesParams{
		Offset: 0,
		Limit:  int32(len(categoryToTestImageDir)),
	}

	categories, err := store.ListCategories(ctx, arg)
	if err != nil {
		return err
	}

	for categoryName, testImageDir := range categoryToTestImageDir {
		imagePaths, err := filepath.Glob(filepath.Join(testImageDir, "*.jpg"))
		if err != nil {
			return err
		}

		for _, imagePath := range imagePaths {
			price := util.RandomPrice()

			categoryID, err := findCategoryIDByName(categories, categoryName)
			if err != nil {
				return err
			}

			imageRelPath, err := filepath.Rel(publicDir, imagePath)
			if err != nil {
				return err
			}

			arg := CreateProductParams{
				Name:          faker.Name(),
				Description:   sql.NullString{String: faker.Paragraph(), Valid: true},
				Price:         price.String(),
				StockQuantity: rand.Int31n(100),
				CategoryID:    categoryID,
				SellerID:      user.ID,
				ImageUrl:      sql.NullString{String: fmt.Sprintf("/%s", imageRelPath), Valid: true},
			}

			_, err = store.CreateProduct(ctx, arg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func findCategoryIDByName(categories []Category, name string) (uuid.UUID, error) {
	for _, category := range categories {
		if category.Name == name {
			return category.ID, nil
		}
	}

	return uuid.UUID{}, fmt.Errorf("category not found: %s", name)
}
