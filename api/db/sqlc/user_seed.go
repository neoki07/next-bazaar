package db

import (
	"context"

	"github.com/ot07/next-bazaar/util"
)

const (
	testUsername  = "testuser"
	testUserEmail = "testuser@email.com"
)

func CreateUserTestData(ctx context.Context, store *SQLStore) error {
	hashedPassword, err := util.HashPassword("password")
	if err != nil {
		return err
	}

	arg := CreateUserParams{
		Name:           testUsername,
		Email:          testUserEmail,
		HashedPassword: hashedPassword,
	}

	_, err = store.CreateUser(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}
