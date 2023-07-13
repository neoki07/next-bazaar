package user_domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/token"
	"github.com/ot07/next-bazaar/util"
)

type UserService struct {
	store db.Store
}

func NewUserService(store db.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := s.store.GetUser(ctx, id)
	if err != nil {
		return User{}, err
	}

	rsp := User{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		HashedPassword:    user.HashedPassword,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	return rsp, err
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	rsp := User{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		HashedPassword:    user.HashedPassword,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	return rsp, err
}

type CreateUserServiceParams struct {
	Name           string
	Email          string
	HashedPassword string
}

func (s *UserService) CreateUser(ctx context.Context, params CreateUserServiceParams) error {
	_, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.HashedPassword,
	})

	return err
}

type CreateSessionServiceParams struct {
	UserID   uuid.UUID
	Duration time.Duration
}

func (s *UserService) CreateSession(ctx context.Context, params CreateSessionServiceParams) (*token.Token, error) {
	sessionToken := token.NewToken(params.Duration)

	_, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		UserID:       params.UserID,
		SessionToken: sessionToken.ID,
		ExpiredAt:    sessionToken.ExpiredAt,
	})
	if err != nil {
		return nil, err
	}

	return sessionToken, nil
}

type RegisterServiceParams struct {
	Name     string
	Email    string
	Password string
}

func (s *UserService) Register(ctx context.Context, params RegisterServiceParams) error {
	hashedPassword, err := util.HashPassword(params.Password)
	if err != nil {
		return err
	}

	return s.CreateUser(ctx, CreateUserServiceParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
}

type LoginServiceParams struct {
	Email    string
	Password string
}

func (s *UserService) Login(ctx context.Context, params LoginServiceParams) (*token.Token, error) {
	user, err := s.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return nil, err
	}

	err = util.CheckPassword(params.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}

	arg := CreateSessionServiceParams{
		UserID:   user.ID,
		Duration: time.Hour * 24 * 7,
	}

	sessionToken, err := s.CreateSession(ctx, arg)
	if err != nil {
		return nil, err
	}

	return sessionToken, nil
}

func (s *UserService) Logout(ctx context.Context, sessionTokenID uuid.UUID) error {
	return s.store.DeleteSession(ctx, sessionTokenID)
}
