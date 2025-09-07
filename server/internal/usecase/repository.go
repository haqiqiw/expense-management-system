package usecase

import (
	"context"
	"expense-management-system/internal/db"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"time"
)

//go:generate mockery --name=UserRepository --structname UserRepository --outpkg=mocks --output=./../mocks
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uint64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	CountByEmail(ctx context.Context, email string) (int, error)
}

//go:generate mockery --name=ExpenseRepository --structname ExpenseRepository --outpkg=mocks --output=./../mocks
type ExpenseRepository interface {
	Create(ctx context.Context, expense *entity.Expense) error
	List(ctx context.Context, req *model.ListExpenseRequest) ([]entity.ExpenseWithUser, int, error)
	FindDetailByID(ctx context.Context, id uint64) (*entity.ExpenseDetail, error)
	FindByID(ctx context.Context, id uint64) (*entity.Expense, error)
	FindByIDWithLock(ctx context.Context, exec db.Executor, id uint64) (*entity.Expense, error)
	UpdateStatusByIDTx(ctx context.Context, exec db.Executor, id uint64, status entity.ExpenseStatus) error
	CompleteByIDTx(ctx context.Context, exec db.Executor, id uint64, processedAt time.Time) error
}

//go:generate mockery --name=ApprovalRepository --structname ApprovalRepository --outpkg=mocks --output=./../mocks
type ApprovalRepository interface {
	CreateTx(ctx context.Context, exec db.Executor, approval *entity.Approval) error
}

//go:generate mockery --name=PaymentPartnerRepository --structname PaymentPartnerRepository --outpkg=mocks --output=./../mocks
type PaymentPartnerRepository interface {
	Execute(ctx context.Context, req *model.PaymentPartnerRequest) (*model.PaymentPartnerResponse, error)
}
