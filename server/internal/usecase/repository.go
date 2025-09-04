package usecase

import (
	"context"
	"expense-management-system/internal/entity"
)

//go:generate mockery --name=UserRepository --structname UserRepository --outpkg=mocks --output=./../mocks
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uint64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	CountByEmail(ctx context.Context, email string) (int, error)
}
