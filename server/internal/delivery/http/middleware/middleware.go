package middleware

import (
	"errors"
	"expense-management-system/internal/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func RecoverMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = fmt.Errorf("%+v", t)
				}

				logger.Error(fmt.Sprintf("panic recovered: %+v", err),
					zap.Any("request_id", requestid.Get(ctx)),
					zap.Any("path", ctx.Request.RequestURI),
					zap.Any("method", ctx.Request.Method),
					zap.Error(err),
				)

				resp := model.ErrorResponse{}
				resp.Errors = model.ErrInternalServerError.Errors
				resp.Meta.HTTPStatus = model.ErrInternalServerError.HTTPStatus

				ctx.JSON(http.StatusInternalServerError, resp)
			}
		}()

		ctx.Next()
	}
}

func ErrorMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			logger.Error(err.Error(),
				zap.Any("request_id", requestid.Get(ctx)),
				zap.Any("path", ctx.Request.RequestURI),
				zap.Any("method", ctx.Request.Method),
				zap.Error(err),
			)

			resp := model.ErrorResponse{}

			var customErr *model.CustomError
			if errors.As(err, &customErr) {
				resp.Errors = customErr.Errors
				resp.Meta.HTTPStatus = customErr.HTTPStatus

				ctx.JSON(customErr.HTTPStatus, resp)
				ctx.Status(customErr.HTTPStatus)
				return
			}

			if valErrs, ok := err.(validator.ValidationErrors); ok {
				errItems := make([]model.ErrorItem, len(valErrs))
				for i, fe := range valErrs {
					errItems[i] = model.ErrorItem{
						Code:    2000 + i,
						Message: fmt.Sprintf("%s failed on the '%s' rule", fe.Field(), fe.Tag()),
					}
				}

				resp.Errors = errItems
				resp.Meta.HTTPStatus = http.StatusBadRequest

				ctx.JSON(http.StatusBadRequest, resp)
				ctx.Status(http.StatusBadRequest)
				return
			}

			resp.Errors = model.ErrInternalServerError.Errors
			resp.Meta.HTTPStatus = model.ErrInternalServerError.HTTPStatus

			ctx.JSON(http.StatusInternalServerError, resp)
			ctx.Status(http.StatusInternalServerError)
		}
	}
}

func RequestLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		logger.Info("request finished",
			zap.Any("request_id", requestid.Get(ctx)),
			zap.Any("path", ctx.Request.RequestURI),
			zap.Any("method", ctx.Request.Method),
			zap.Any("status", ctx.Writer.Status()),
			zap.Any("path", ctx.Request.RequestURI),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

func CorsMiddleware(origins []string) gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(config)
}
