package model

import "expense-management-system/internal/auth"

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LogoutRequest struct {
	Claims *auth.JWTClaims `json:"claims"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
