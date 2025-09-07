package middleware

import (
	"expense-management-system/internal/metrics"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRequestLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.FullPath(), "/swagger") || strings.HasPrefix(ctx.FullPath(), "/metrics") {
			ctx.Next()
			return
		}

		start := time.Now()
		ctx.Next()

		duration := time.Since(start)
		statusCode := ctx.Writer.Status()

		var status string
		if statusCode >= 500 {
			status = "fail"
		} else {
			status = "ok"
		}

		logger.Info("request finished",
			zap.Any("request_id", requestid.Get(ctx)),
			zap.Any("path", ctx.Request.RequestURI),
			zap.Any("method", ctx.Request.Method),
			zap.Any("status", statusCode),
			zap.Any("path", ctx.Request.RequestURI),
			zap.Duration("duration", duration),
		)

		metrics.RequestsTotal.WithLabelValues(ctx.Request.Method, ctx.FullPath(), status).Inc()
		metrics.RequestDuration.WithLabelValues(ctx.Request.Method, ctx.FullPath(), status).Observe(duration.Seconds())
	}
}
