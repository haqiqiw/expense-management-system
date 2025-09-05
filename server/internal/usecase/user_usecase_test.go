package usecase_test

import (
	"context"
	"errors"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type UserUsecaseSuite struct {
	suite.Suite
	log *zap.Logger
	ctx context.Context
}

func (s *UserUsecaseSuite) SetupTest() {
	s.log, _ = zap.NewDevelopment()
	s.ctx = context.Background()
}

func (s *UserUsecaseSuite) TestUserUsecase_Create() {
	now := time.Now()

	tests := []struct {
		name       string
		request    *model.CreateUserRequest
		mockFunc   func(r *mocks.UserRepository)
		wantUser   *model.UserResponse
		wantErrMsg string
	}{
		{
			name: "error on count",
			request: &model.CreateUserRequest{
				Email:    "john@mail.com",
				Name:     "John Doe",
				Password: "password",
				Role:     "manager",
			},
			mockFunc: func(r *mocks.UserRepository) {
				r.On("CountByEmail", mock.Anything, "john@mail.com").
					Return(0, errors.New("something error"))
			},
			wantUser:   nil,
			wantErrMsg: "failed to count by email (john@mail.com) = something error",
		},
		{
			name: "error on duplicate email",
			request: &model.CreateUserRequest{
				Email:    "john@mail.com",
				Name:     "John Doe",
				Password: "password",
				Role:     "manager",
			},
			mockFunc: func(r *mocks.UserRepository) {
				r.On("CountByEmail", mock.Anything, "john@mail.com").
					Return(1, nil)
			},
			wantUser:   nil,
			wantErrMsg: "Email already exist",
		},
		{
			name: "error on parse role",
			request: &model.CreateUserRequest{
				Email:    "john@mail.com",
				Name:     "John Doe",
				Password: "password",
				Role:     "unknown",
			},
			mockFunc: func(r *mocks.UserRepository) {
				r.On("CountByEmail", mock.Anything, "john@mail.com").
					Return(0, nil)
			},
			wantUser:   nil,
			wantErrMsg: "failed to parse user role for email (john@mail.com) = invalid user role = unknown",
		},
		{
			name: "error on create",
			request: &model.CreateUserRequest{
				Email:    "john@mail.com",
				Name:     "John Doe",
				Password: "password",
				Role:     "manager",
			},
			mockFunc: func(r *mocks.UserRepository) {
				r.On("CountByEmail", mock.Anything, "john@mail.com").
					Return(0, nil)
				r.On("Create", mock.Anything, mock.Anything).
					Return(errors.New("something error"))
			},
			wantUser:   nil,
			wantErrMsg: "failed to create user for email (john@mail.com) = something error",
		},
		{
			name: "success",
			request: &model.CreateUserRequest{
				Email:    "john@mail.com",
				Name:     "John Doe",
				Password: "password",
				Role:     "manager",
			},
			mockFunc: func(r *mocks.UserRepository) {
				r.On("CountByEmail", mock.Anything, "john@mail.com").
					Return(0, nil)
				r.On("Create", mock.Anything, mock.Anything).
					Return(nil)
			},
			wantUser: &model.UserResponse{
				ID:        1,
				Email:     "john@mail.com",
				Name:      "John Doe",
				Role:      "manager",
				CreatedAt: now.Format(time.RFC3339),
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			userRepository := mocks.NewUserRepository(s.T())
			usecase := usecase.NewUserUsecase(s.log, userRepository)
			tt.mockFunc(userRepository)

			_, err := usecase.Create(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				s.Nil(err)
			}
		})
	}
}

func (s *UserUsecaseSuite) TestUserUsecase_FindByID() {
	now := time.Now()

	tests := []struct {
		name       string
		request    *model.GetUserRequest
		mockFunc   func(r *mocks.UserRepository)
		wantUser   *model.UserResponse
		wantErrMsg string
	}{
		{
			name: "error on find",
			request: &model.GetUserRequest{
				ID: 1,
			},
			mockFunc: func(r *mocks.UserRepository) {
				r.On("FindByID", mock.Anything, uint64(1)).
					Return(nil, errors.New("something error"))
			},
			wantUser:   nil,
			wantErrMsg: "failed to find user by id (1) = something error",
		},
		{
			name: "not found",
			request: &model.GetUserRequest{
				ID: 1,
			},
			mockFunc: func(r *mocks.UserRepository) {
				r.On("FindByID", mock.Anything, uint64(1)).
					Return(nil, nil)
			},
			wantUser:   nil,
			wantErrMsg: "User not found",
		},
		{
			name: "success",
			request: &model.GetUserRequest{
				ID: 1,
			},
			mockFunc: func(r *mocks.UserRepository) {
				r.On("FindByID", mock.Anything, uint64(1)).Return(&entity.User{
					ID:           1,
					Email:        "john@mail.com",
					Name:         "John Doe",
					PasswordHash: "password",
					Role:         "manager",
					CreatedAt:    now,
				}, nil)
			},
			wantUser: &model.UserResponse{
				ID:        1,
				Email:     "john@mail.com",
				Name:      "John Doe",
				Role:      "manager",
				CreatedAt: now.UTC().Format(time.RFC3339),
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			userRepository := mocks.NewUserRepository(s.T())
			usecase := usecase.NewUserUsecase(s.log, userRepository)
			tt.mockFunc(userRepository)

			res, err := usecase.FindByID(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Nil(res)
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				s.Equal(*tt.wantUser, *res)
				s.Nil(err)
			}
		})
	}
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}
