package config

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewGin(logger *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	return engine
}
