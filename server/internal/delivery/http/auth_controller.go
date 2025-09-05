package http

import (
	"expense-management-system/internal/delivery/http/middleware"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type AuthController struct {
	log         *zap.Logger
	validate    *validator.Validate
	authUsecase usecase.AuthUsecase
}

func NewAuthController(log *zap.Logger, validate *validator.Validate,
	authUsecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		log:         log,
		validate:    validate,
		authUsecase: authUsecase,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	request := new(model.LoginRequest)
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		LogWarn(ctx, c.log, "failed to parse request body", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	err = c.validate.Struct(request)
	if err != nil {
		LogWarn(ctx, c.log, "failed to validate request body", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	res, err := c.authUsecase.Login(ctx.Request.Context(), request)
	if err != nil {
		LogWarn(ctx, c.log, "failed to login", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		model.NewSuccessResponse(res, http.StatusOK),
	)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	claims, err := middleware.GetJWTClaims(ctx)
	if err != nil {
		LogWarn(ctx, c.log, "failed to get jwt claims", err)
		ctx.Error(model.ErrUnauthorized)
		return
	}

	request := &model.LogoutRequest{
		Claims: claims,
	}
	err = c.authUsecase.Logout(ctx.Request.Context(), request)
	if err != nil {
		LogWarn(ctx, c.log, "failed to logout", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		model.NewSuccessMessageResponse("Logged out", http.StatusOK),
	)
}
