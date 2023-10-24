package db

import (
	"context"

	"github.com/ot07/next-bazaar/util"
)

func CreateUserTestData(ctx context.Context, store *SQLStore, config util.Config) error {
	testAccounts := [3]map[string]string{
		{
			"username": config.TestAccountUsername1,
			"email":    config.TestAccountEmail1,
		},
		{
			"username": config.TestAccountUsername2,
			"email":    config.TestAccountEmail2,
		},
		{
			"username": config.TestAccountUsername3,
			"email":    config.TestAccountEmail3,
		},
	}

	hashedPassword, err := util.HashPassword(config.TestAccountPassword)
	if err != nil {
		return err
	}

	for _, testAccount := range testAccounts {
		arg := CreateUserParams{
			Name:           testAccount["username"],
			Email:          testAccount["email"],
			HashedPassword: hashedPassword,
		}

		_, err = store.CreateUser(ctx, arg)
		if err != nil {
			return err
		}
	}

	return nil
}
