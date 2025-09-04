package model

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,min=4,max=100,email"`
	Name     string `json:"name" validate:"required,min=4,max=100"`
	Password string `json:"password" validate:"required,min=4,max=100"`
	Role     string `json:"role" validate:"required,oneof=employee manager"`
}

type GetUserRequest struct {
	ID uint64 `json:"id"`
}

type UserResponse struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
