package serializer_test

import (
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"expense-management-system/internal/model/serializer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserSerializer_UserToResponse(t *testing.T) {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		param   *entity.User
		wantRes *model.UserResponse
	}{
		{
			name: "success",
			param: &entity.User{
				ID:           1,
				Email:        "john@mail.com",
				Name:         "John Doe",
				PasswordHash: "password",
				Role:         "manager",
				CreatedAt:    now,
			},
			wantRes: &model.UserResponse{
				ID:        1,
				Email:     "john@mail.com",
				Name:      "John Doe",
				Role:      "manager",
				CreatedAt: now.Format(time.RFC3339),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := serializer.UserToResponse(tt.param)

			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestUserSerializer_UserSimpleToResponse(t *testing.T) {
	tests := []struct {
		name    string
		param   *entity.UserSimple
		wantRes *model.UserSimpleResponse
	}{
		{
			name: "success",
			param: &entity.UserSimple{
				ID:    1,
				Email: "john@mail.com",
				Name:  "John Doe",
			},
			wantRes: &model.UserSimpleResponse{
				ID:    1,
				Email: "john@mail.com",
				Name:  "John Doe",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := serializer.UserSimpleToResponse(tt.param)

			assert.Equal(t, tt.wantRes, res)
		})
	}
}
