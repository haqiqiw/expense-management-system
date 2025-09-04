package usecase_test

import (
	"context"
	"errors"
	"expense-management-system/internal/auth"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseSuite struct {
	suite.Suite
	log *zap.Logger
	ctx context.Context
}

type MockFunc func(
	c context.Context,
	rc *mocks.RedisClient,
	jwt *mocks.JWTToken,
	ur *mocks.UserRepository,
)

func (s *AuthUsecaseSuite) SetupTest() {
	s.log, _ = zap.NewDevelopment()
	s.ctx = context.Background()
}

func (s *AuthUsecaseSuite) TestAuthUsecase_Login() {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	now := time.Now()

	tests := []struct {
		name       string
		request    *model.LoginRequest
		mockFunc   MockFunc
		wantRes    *model.LoginResponse
		wantErrMsg string
	}{
		{
			name: "error on find by email",
			request: &model.LoginRequest{
				Email:    "john@mail.com",
				Password: "password",
			},
			mockFunc: func(
				c context.Context,
				rc *mocks.RedisClient,
				jwt *mocks.JWTToken,
				ur *mocks.UserRepository,
			) {
				ur.On("FindByEmail", mock.Anything, "john@mail.com").
					Return(nil, errors.New("something error"))
			},
			wantRes:    nil,
			wantErrMsg: "failed to find user by email: something error",
		},
		{
			name: "error user not found",
			request: &model.LoginRequest{
				Email:    "john@mail.com",
				Password: "password",
			},
			mockFunc: func(
				c context.Context,
				rc *mocks.RedisClient,
				jwt *mocks.JWTToken,
				ur *mocks.UserRepository,
			) {
				ur.On("FindByEmail", mock.Anything, "john@mail.com").
					Return(nil, nil)
			},
			wantRes:    nil,
			wantErrMsg: "user not found",
		},
		{
			name: "error invalid password",
			request: &model.LoginRequest{
				Email:    "john@mail.com",
				Password: "invalid_password",
			},
			mockFunc: func(
				c context.Context,
				rc *mocks.RedisClient,
				jwt *mocks.JWTToken,
				ur *mocks.UserRepository,
			) {
				ur.On("FindByEmail", mock.Anything, "john@mail.com").Return(&entity.User{
					ID:           uint64(1),
					Email:        "john@mail.com",
					Name:         "John Doe",
					PasswordHash: string(passwordHash),
					Role:         "manager",
					CreatedAt:    now,
				}, nil)
			},
			wantRes:    nil,
			wantErrMsg: "invalid password",
		},
		{
			name: "error on create jwt token",
			request: &model.LoginRequest{
				Email:    "john@mail.com",
				Password: "password",
			},
			mockFunc: func(
				c context.Context,
				rc *mocks.RedisClient,
				jwt *mocks.JWTToken,
				ur *mocks.UserRepository,
			) {
				ur.On("FindByEmail", mock.Anything, "john@mail.com").Return(&entity.User{
					ID:           uint64(1),
					Email:        "john@mail.com",
					Name:         "John Doe",
					PasswordHash: string(passwordHash),
					Role:         "manager",
					CreatedAt:    now,
				}, nil)
				jwt.On("Create", "1", "manager").Return("", errors.New("something error"))
			},
			wantRes:    nil,
			wantErrMsg: "failed to create access token: something error",
		},
		{
			name: "success",
			request: &model.LoginRequest{
				Email:    "john@mail.com",
				Password: "password",
			},
			mockFunc: func(
				c context.Context,
				rc *mocks.RedisClient,
				jwt *mocks.JWTToken,
				ur *mocks.UserRepository,
			) {
				ur.On("FindByEmail", mock.Anything, "john@mail.com").Return(&entity.User{
					ID:           uint64(1),
					Email:        "john@mail.com",
					Name:         "John Doe",
					PasswordHash: string(passwordHash),
					Role:         "manager",
					CreatedAt:    now,
				}, nil)
				jwt.On("Create", "1", "manager").Return("qwerty-12345", nil)
			},
			wantRes: &model.LoginResponse{
				AccessToken: "qwerty-12345",
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			rc := mocks.NewRedisClient(s.T())
			jwt := mocks.NewJWTToken(s.T())
			ur := mocks.NewUserRepository(s.T())
			usecase := usecase.NewAuthUsecase(s.log, rc, jwt, ur)
			tt.mockFunc(s.ctx, rc, jwt, ur)

			res, err := usecase.Login(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Nil(res)
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				s.Equal(*tt.wantRes, *res)
				s.Nil(err)
			}
		})
	}
}

func (s *AuthUsecaseSuite) TestAuthUsecase_Logout() {
	now := time.Now()

	tests := []struct {
		name       string
		request    *model.LogoutRequest
		mockFunc   func(ctx context.Context, rc *mocks.RedisClient)
		wantErrMsg string
	}{
		{
			name: "error on set revoke token cache",
			request: &model.LogoutRequest{
				Claims: &auth.JWTClaims{
					UserID: "1",
					Role:   "manager",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Minute)),
						ID:        "asd-789",
					},
				},
			},
			mockFunc: func(ctx context.Context, rc *mocks.RedisClient) {
				setCmd := redis.NewStatusCmd(s.ctx)
				setCmd.SetErr(errors.New("something error"))
				rc.On("SetEx", mock.Anything, "revoke-jwt-token:asd-789", "true", mock.Anything).
					Return(setCmd)
			},
			wantErrMsg: "failed to set revoke token: something error",
		},
		{
			name: "success",
			request: &model.LogoutRequest{
				Claims: &auth.JWTClaims{
					UserID: "1",
					Role:   "manager",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Minute)),
						ID:        "asd-789",
					},
				},
			},
			mockFunc: func(ctx context.Context, rc *mocks.RedisClient) {
				setCmd := redis.NewStatusCmd(s.ctx)
				rc.On("SetEx", mock.Anything, "revoke-jwt-token:asd-789", "true", mock.Anything).
					Return(setCmd)
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			rc := mocks.NewRedisClient(s.T())
			jwt := mocks.NewJWTToken(s.T())
			ur := mocks.NewUserRepository(s.T())
			usecase := usecase.NewAuthUsecase(s.log, rc, jwt, ur)
			tt.mockFunc(s.ctx, rc)

			err := usecase.Logout(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				s.Nil(err)
			}
		})
	}
}

func TestAuthUsecaseSuite(t *testing.T) {
	suite.Run(t, new(AuthUsecaseSuite))
}
