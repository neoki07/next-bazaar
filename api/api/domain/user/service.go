package user_domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ot07/next-bazaar/token"
)

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, err
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	return user, err
}

type CreateUserServiceParams struct {
	Name           string
	Email          string
	HashedPassword string
}

func (s *UserService) CreateUser(ctx context.Context, params CreateUserServiceParams) error {
	arg := CreateRepositoryParams(params)

	return s.repository.Create(ctx, arg)
}

type CreateSessionServiceParams struct {
	UserID   uuid.UUID
	Duration time.Duration
}

func (s *UserService) CreateSession(ctx context.Context, params CreateSessionServiceParams) (*token.Token, error) {
	sessionToken := token.NewToken(params.Duration)

	arg := CreateSessionRepositoryParams{
		UserID:       params.UserID,
		SessionToken: sessionToken,
	}

	err := s.repository.CreateSession(ctx, arg)
	if err != nil {
		return nil, err
	}

	return sessionToken, nil
}

func (s *UserService) DeleteSession(ctx context.Context, sessionTokenID uuid.UUID) error {
	return s.repository.DeleteSession(ctx, sessionTokenID)
}
