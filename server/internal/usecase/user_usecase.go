package usecase

import (
	"context"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"expense-management-system/internal/model/serializer"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	log            *zap.Logger
	userRepository UserRepository
}

func NewUserUsecase(log *zap.Logger, userRepository UserRepository) UserUsecase {
	return &userUsecase{
		log:            log,
		userRepository: userRepository,
	}
}

func (c *userUsecase) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	total, err := c.userRepository.CountByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to count by email (%s) = %w", req.Email, err)
	}

	if total > 0 {
		return nil, model.ErrEmailAlreadyExist
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password for email (%s) = %w", req.Email, err)
	}

	role, err := entity.ParseUserRole(req.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user role for email (%s) = %w", req.Email, err)
	}

	user := &entity.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(password),
		Role:         role,
	}

	err = c.userRepository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user for email (%s) = %w", req.Email, err)
	}

	return serializer.UserToResponse(user), nil
}

func (c *userUsecase) FindByID(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	user, err := c.userRepository.FindByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id (%d) = %w", req.ID, err)
	}

	if user == nil {
		return nil, model.ErrUserNotFound
	}

	return serializer.UserToResponse(user), nil
}
