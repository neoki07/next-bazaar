package user_domain

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
