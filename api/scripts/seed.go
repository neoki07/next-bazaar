package main

import (
	"database/sql"
	"fmt"
	"log"

	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/util"
	"golang.org/x/net/context"

	_ "github.com/lib/pq"
	_ "github.com/ot07/next-bazaar/docs"
)

func setup(config util.Config) (context.Context, *db.SQLStore, error) {
	ctx := context.Background()

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	store := db.NewStore(conn)
	return ctx, store, nil
}

func truncateAllTables(ctx context.Context, store *db.SQLStore) error {
	log.Println("truncating all tables...")

	err := store.TruncateProductsTable(ctx)
	if err != nil {
		log.Fatal("cannot truncate products table:", err)
	}

	err = store.TruncateCategoriesTable(ctx)
	if err != nil {
		log.Fatal("cannot truncate categories table:", err)
	}

	err = store.TruncateSessionsTable(ctx)
	if err != nil {
		log.Fatal("cannot truncate sessions table:", err)
	}

	err = store.TruncateUsersTable(ctx)
	if err != nil {
		log.Fatal("cannot truncate users table:", err)
	}

	return nil
}

func runSeed(ctx context.Context, store *db.SQLStore, config util.Config) error {
	log.Println("creating user test data...")
	err := db.CreateUserTestData(ctx, store, config)
	if err != nil {
		log.Fatal("cannot create user test data:", err)
	}

	log.Println("creating category test data...")
	err = db.CreateCategoryTestData(ctx, store)
	if err != nil {
		log.Fatal("cannot create category test data:", err)
	}

	log.Println("creating product test data...")
	err = db.CreateProductTestData(ctx, store, config)
	if err != nil {
		log.Fatal("cannot create product test data:", err)
	}

	return nil
}

func main() {
	log.Println("starting seed...")

	config, err := util.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, store, err := setup(config)
	if err != nil {
		log.Fatalf("failed to set up: %v", err)
	}

	err = truncateAllTables(ctx, store)
	if err != nil {
		log.Fatalf("failed to truncate all tables: %v", err)
	}

	err = runSeed(ctx, store, config)
	if err != nil {
		log.Fatalf("failed to run seed: %v", err)
	}

	log.Println("seed completed successfully")
}
