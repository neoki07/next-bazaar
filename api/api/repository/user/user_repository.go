package user_repository

import (
	"context"

	"github.com/google/uuid"
	user_domain "github.com/ot07/next-bazaar/api/domain/user"
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
	FindByEmail(ctx context.Context, email string) (*user_domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*user_domain.User, error)
	Create(ctx context.Context, params CreateParams) error
	CreateSession(ctx context.Context, params CreateSessionParams) error
	DeleteSession(ctx context.Context, sessionToken uuid.UUID) error
}

func NewUserRepository(store db.Store) UserRepository {
	return &userRepositoryImpl{
		store: store,
	}
}
