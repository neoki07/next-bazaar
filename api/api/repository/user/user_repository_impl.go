package user_repository

import (
	"context"

	"github.com/google/uuid"
	user_domain "github.com/ot07/next-bazaar/api/domain/user"
	db "github.com/ot07/next-bazaar/db/sqlc"
)

type userRepositoryImpl struct {
	store db.Store
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (user_domain.User, error) {
	user, err := r.store.GetUser(ctx, id)
	if err != nil {
		return user_domain.User{}, err
	}

	rsp := user_domain.User{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		HashedPassword:    user.HashedPassword,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	return rsp, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (user_domain.User, error) {
	user, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return user_domain.User{}, err
	}

	rsp := user_domain.User{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		HashedPassword:    user.HashedPassword,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	return rsp, nil
}

func (r *userRepositoryImpl) Create(ctx context.Context, params CreateParams) error {
	_, err := r.store.CreateUser(ctx, db.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.HashedPassword,
	})

	return err
}

func (r *userRepositoryImpl) CreateSession(ctx context.Context, params CreateSessionParams) error {
	_, err := r.store.CreateSession(ctx, db.CreateSessionParams{
		UserID:       params.UserID,
		SessionToken: params.SessionToken.ID,
		ExpiredAt:    params.SessionToken.ExpiredAt,
	})

	return err
}

func (r *userRepositoryImpl) DeleteSession(ctx context.Context, sessionToken uuid.UUID) error {
	return r.store.DeleteSession(ctx, sessionToken)
}
