package middleware

import (
	"errors"
	"expense-management-system/internal/model"
	"fmt"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func NewRecoverMiddleware(logger *zap.Logger) gin.HandlerFunc {
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

func NewErrorMiddleware(logger *zap.Logger) gin.HandlerFunc {
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
				return
			}

			resp.Errors = model.ErrInternalServerError.Errors
			resp.Meta.HTTPStatus = model.ErrInternalServerError.HTTPStatus

			ctx.JSON(http.StatusInternalServerError, resp)
		}
	}
}
