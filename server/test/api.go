package test

import (
	"expense-management-system/internal/config"
	"expense-management-system/internal/delivery/http/middleware"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewApi(logger *zap.Logger) *gin.Engine {
	app := config.NewGin(logger)

	app.Use(requestid.New())
	app.Use(middleware.NewRequestLoggerMiddleware(logger))
	app.Use(middleware.NewRecoverMiddleware(logger))
	app.Use(middleware.NewErrorMiddleware(logger))

	return app
}
