package db

import (
	"context"

	"github.com/ot07/next-bazaar/util"
)

func CreateUserTestData(ctx context.Context, store *SQLStore, config util.Config) error {
	hashedPassword, err := util.HashPassword(config.TestAccountPassword)
	if err != nil {
		return err
	}

	arg := CreateUserParams{
		Name:           config.TestAccountUsername,
		Email:          config.TestAccountEmail,
		HashedPassword: hashedPassword,
	}

	_, err = store.CreateUser(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}
