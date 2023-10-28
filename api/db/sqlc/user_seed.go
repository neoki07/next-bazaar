package db

import (
	"context"

	"github.com/ot07/next-bazaar/util"
)

func CreateUserTestData(ctx context.Context, store *SQLStore, config util.Config) error {
	for _, testAccount := range config.TestAccounts {
		hashedPassword, err := util.HashPassword(testAccount.Password)
		if err != nil {
			return err
		}
		arg := CreateUserParams{
			Name:           testAccount.Username,
			Email:          testAccount.Email,
			HashedPassword: hashedPassword,
		}

		_, err = store.CreateUser(ctx, arg)
		if err != nil {
			return err
		}
	}

	return nil
}
