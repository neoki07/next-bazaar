package user_domain

import (
	"context"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
)

type CreateParams struct {
	Name           string
	Email          string
	HashedPassword string
}

type CreateSessionParams struct {
	UserID       uuid.UUID
	SessionToken *token.Token
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id uuid.UUID) (User, error)
	Create(ctx context.Context, params CreateParams) error
	CreateSession(ctx context.Context, params CreateSessionParams) error
	DeleteSession(ctx context.Context, sessionToken uuid.UUID) error
}

func NewUserRepository(store db.Store) UserRepository {
	return &userRepositoryImpl{
		store: store,
	}
}
