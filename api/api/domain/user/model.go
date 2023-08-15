package user_domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID
	Name              string
	Email             string
	HashedPassword    string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
}

type Session struct {
	ID           uuid.UUID
	SessionToken uuid.UUID
	UserID       uuid.UUID
	ExpiredAt    time.Time
	CreatedAt    time.Time
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,without_space,without_punct,without_symbol"`
	Email    string `json:"email" validate:"required,email" swaggertype:"string"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" swaggertype:"string"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateRequest struct {
	Name  string `json:"name" validate:"required,without_space,without_punct,without_symbol"`
	Email string `json:"email" validate:"required,email" swaggertype:"string"`
}

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email" swaggertype:"string"`
}

func NewUserResponse(user User) UserResponse {
	return UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}
}
