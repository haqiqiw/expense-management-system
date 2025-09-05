package http

import (
	"expense-management-system/internal/delivery/http/middleware"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApprovalController struct {
	log             *zap.Logger
	approvalUsecase usecase.ApprovalUsecase
}

func NewApprovalController(log *zap.Logger, approvalUsecase usecase.ApprovalUsecase) *ApprovalController {
	return &ApprovalController{
		log:             log,
		approvalUsecase: approvalUsecase,
	}
}

func (c *ApprovalController) Approve(ctx *gin.Context) {
	claims, err := middleware.GetJWTClaims(ctx)
	if err != nil {
		LogWarn(ctx, c.log, "failed to get jwt claims", err)
		ctx.Error(model.ErrUnauthorized)
		return
	}

	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		LogWarn(ctx, c.log, "failed to convert user id", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		LogWarn(ctx, c.log, "failed to convert id", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	request := new(model.ApprovalExpenseRequest)

	if ctx.Request.Body != nil && ctx.Request.ContentLength > 0 {
		err = ctx.ShouldBindJSON(request)
		if err != nil {
			LogWarn(ctx, c.log, "failed to parse request body", err)
			ctx.Error(model.ErrBadRequest)
			return
		}
	}

	request.ID = id
	request.UserID = userID
	request.UserRole = claims.Role
	err = c.approvalUsecase.Approve(ctx.Request.Context(), request)
	if err != nil {
		LogWarn(ctx, c.log, "failed to approve expense", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		model.NewSuccessMessageResponse("Expense approved", http.StatusOK),
	)
}

func (c *ApprovalController) Reject(ctx *gin.Context) {
	claims, err := middleware.GetJWTClaims(ctx)
	if err != nil {
		LogWarn(ctx, c.log, "failed to get jwt claims", err)
		ctx.Error(model.ErrUnauthorized)
		return
	}

	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		LogWarn(ctx, c.log, "failed to convert user id", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		LogWarn(ctx, c.log, "failed to convert id", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	request := new(model.ApprovalExpenseRequest)

	if ctx.Request.Body != nil && ctx.Request.ContentLength > 0 {
		err = ctx.ShouldBindJSON(request)
		if err != nil {
			LogWarn(ctx, c.log, "failed to parse request body", err)
			ctx.Error(model.ErrBadRequest)
			return
		}
	}

	request.ID = id
	request.UserID = userID
	request.UserRole = claims.Role
	err = c.approvalUsecase.Reject(ctx.Request.Context(), request)
	if err != nil {
		LogWarn(ctx, c.log, "failed to reject expense", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		model.NewSuccessMessageResponse("Expense rejected", http.StatusOK),
	)
}
