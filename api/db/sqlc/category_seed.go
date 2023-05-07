package db

import (
	"context"
)

func CreateCategoryTestData(ctx context.Context, store *SQLStore) error {
	categoryNames := []string{
		"Jeans",
		"Sofa",
		"T-Shirt",
		"TV",
	}

	for _, name := range categoryNames {
		_, err := store.CreateCategory(ctx, name)
		if err != nil {
			return err
		}
	}

	return nil
}
