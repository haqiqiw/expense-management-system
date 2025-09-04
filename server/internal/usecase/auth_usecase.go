package usecase

import (
	"context"
	"expense-management-system/internal/auth"
	"expense-management-system/internal/model"
	"expense-management-system/internal/storage"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	log            *zap.Logger
	redisClient    storage.RedisClient
	jwtToken       auth.JWTToken
	userRepository UserRepository
}

func NewAuthUsecase(log *zap.Logger, redisClient storage.RedisClient, jwtToken auth.JWTToken, userRepository UserRepository) AuthUsecase {
	return &authUsecase{
		log:            log,
		redisClient:    redisClient,
		jwtToken:       jwtToken,
		userRepository: userRepository,
	}
}

func (c *authUsecase) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := c.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	if user == nil {
		return nil, model.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, model.ErrInvalidPassword
	}

	accessToken, err := c.jwtToken.Create(fmt.Sprint(user.ID), user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	return &model.LoginResponse{
		AccessToken: accessToken,
	}, nil
}

func (c *authUsecase) Logout(ctx context.Context, req *model.LogoutRequest) error {
	revokeKey := fmt.Sprintf("%s:%s", auth.PrefixRevokeKey, req.Claims.ID)
	revokeTTL := time.Until(req.Claims.ExpiresAt.Time)

	err := c.redisClient.SetEx(ctx, revokeKey, "true", revokeTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set revoke token: %w", err)
	}

	return nil
}
