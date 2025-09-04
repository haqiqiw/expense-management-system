package http

import (
	"expense-management-system/internal/delivery/http/middleware"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type UserController struct {
	Log         *zap.Logger
	Validate    *validator.Validate
	UserUsecase usecase.UserUsecase
}

func NewUserController(log *zap.Logger, validate *validator.Validate, userUsecase usecase.UserUsecase) *UserController {
	return &UserController{
		Log:         log,
		Validate:    validate,
		UserUsecase: userUsecase,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	request := new(model.CreateUserRequest)
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		LogWarn(ctx, c.Log, "failed to parse request body", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	err = c.Validate.Struct(request)
	if err != nil {
		LogWarn(ctx, c.Log, "failed to validate request body", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	res, err := c.UserUsecase.Create(ctx.Request.Context(), request)
	if err != nil {
		LogWarn(ctx, c.Log, "failed to register user", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		model.NewSuccessResponse(res, http.StatusCreated),
	)
}

func (c *UserController) Me(ctx *gin.Context) {
	claims, err := middleware.GetJWTClaims(ctx)
	if err != nil {
		LogWarn(ctx, c.Log, "failed to get jwt claims", err)
		ctx.Error(model.ErrUnauthorized)
		return
	}

	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		LogWarn(ctx, c.Log, "failed to convert id", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	res, err := c.UserUsecase.FindByID(ctx.Request.Context(), &model.GetUserRequest{
		ID: userID,
	})
	if err != nil {
		LogWarn(ctx, c.Log, "failed to get user", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		model.NewSuccessResponse(res, http.StatusOK),
	)
}
