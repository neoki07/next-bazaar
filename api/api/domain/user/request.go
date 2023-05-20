package user_domain

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,without_space,without_punct,without_symbol"`
	Email    string `json:"email" validate:"required,email" swaggertype:"string"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email" swaggertype:"string"`
	Password string `json:"password" validate:"required,min=8"`
}
