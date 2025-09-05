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

//go:generate mockery --name=ExpenseUsecase --structname ExpenseUsecase --outpkg=mocks --output=./../mocks
type ExpenseUsecase interface {
	Create(ctx context.Context, req *model.CreateExpenseRequest) (*model.ExpenseCreateResponse, error)
	List(ctx context.Context, req *model.ListExpenseRequest) ([]model.ExpenseWithUserResponse, int, error)
	FindByID(ctx context.Context, req *model.GetExpenseRequest) (*model.ExpenseDetailResponse, error)
}

//go:generate mockery --name=ApprovalUsecase --structname ApprovalUsecase --outpkg=mocks --output=./../mocks
type ApprovalUsecase interface {
	Approve(ctx context.Context, req *model.ApprovalExpenseRequest) error
	Reject(ctx context.Context, req *model.ApprovalExpenseRequest) error
}

//go:generate mockery --name=PaymentProcessorUsecase --structname PaymentProcessorUsecase --outpkg=mocks --output=./../mocks
type PaymentProcessorUsecase interface {
	Execute(ctx context.Context, req *model.PaymentProcessorRequest) error
}
