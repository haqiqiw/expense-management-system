package middleware

import (
	"expense-management-system/internal/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	redisrate "github.com/go-redis/redis_rate/v10"
)

const limiterRateRequest = "rate_request_%s"

func NewLimiterMiddleware(rate *redisrate.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, _ := rate.Allow(ctx, fmt.Sprintf(limiterRateRequest, ctx.ClientIP()), redisrate.Limit{
			Rate:   10,
			Burst:  10,
			Period: time.Second,
		})

		if res.Allowed <= 0 {
			resp := model.ErrorResponse{}
			resp.Errors = model.ErrTooManyRequest.Errors
			resp.Meta.HTTPStatus = model.ErrTooManyRequest.HTTPStatus

			ctx.JSON(http.StatusTooManyRequests, resp)
			ctx.Abort()
		}

		ctx.Next()
	}
}
