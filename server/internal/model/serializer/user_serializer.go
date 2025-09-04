package serializer

import (
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"time"
)

func UserToResponse(u *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
	}
}
