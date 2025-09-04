package test

import (
	"expense-management-system/internal/auth"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewAuthMiddleware(userID uint64, role string) gin.HandlerFunc {
	now := time.Now()
	claims := &auth.JWTClaims{
		UserID: fmt.Sprint(userID),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	return func(ctx *gin.Context) {
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
