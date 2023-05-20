package user_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	user_domain "github.com/ot07/next-bazaar/api/domain/user"
	user_repository "github.com/ot07/next-bazaar/api/repository/user"
	"github.com/ot07/next-bazaar/token"
)

type UserService struct {
	repository user_repository.UserRepository
}

func NewUserService(repository user_repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (user_domain.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return user_domain.User{}, err
	}

	return user, err
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (user_domain.User, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return user_domain.User{}, err
	}

	return user, err
}

type CreateUserParams struct {
	Name           string
	Email          string
	HashedPassword string
}

func (s *UserService) CreateUser(ctx context.Context, params CreateUserParams) error {
	arg := user_repository.CreateParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.HashedPassword,
	}

	return s.repository.Create(ctx, arg)
}

func (s *UserService) CreateSession(ctx context.Context, userID uuid.UUID, duration time.Duration) (*token.Token, error) {
	sessionToken := token.NewToken(duration)

	arg := user_repository.CreateSessionParams{
		UserID:       userID,
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
