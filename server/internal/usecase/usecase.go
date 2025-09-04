package usecase

import (
	"context"
	"expense-management-system/internal/model"
)

//go:generate mockery --name=AuthUsecase --structname AuthUsecase --outpkg=mocks --output=./../mocks
type AuthUsecase interface {
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	Logout(ctx context.Context, req *model.LogoutRequest) error
}

//go:generate mockery --name=UserUsecase --structname UserUsecase --outpkg=mocks --output=./../mocks
type UserUsecase interface {
	Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
	FindByID(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error)
}
