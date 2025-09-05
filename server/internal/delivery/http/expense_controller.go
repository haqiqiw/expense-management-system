package http

import (
	"expense-management-system/internal/delivery/http/middleware"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ExpenseController struct {
	log            *zap.Logger
	validate       *validator.Validate
	expenseUsecase usecase.ExpenseUsecase
}

func NewExpenseController(log *zap.Logger, validate *validator.Validate, expenseUsecase usecase.ExpenseUsecase) *ExpenseController {
	return &ExpenseController{
		log:            log,
		validate:       validate,
		expenseUsecase: expenseUsecase,
	}
}

func (c *ExpenseController) Create(ctx *gin.Context) {
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

	request := new(model.CreateExpenseRequest)
	err = ctx.ShouldBindJSON(request)
	if err != nil {
		LogWarn(ctx, c.log, "failed to parse request body", err)
		ctx.Error(model.ErrBadRequest)
		return
	}

	err = c.validate.Struct(request)
	if err != nil {
		ctx.Error(err)
		return
	}

	request.UserID = userID
	res, err := c.expenseUsecase.Create(ctx.Request.Context(), request)
	if err != nil {
		LogWarn(ctx, c.log, "failed to create expense", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		model.NewSuccessResponse(res, http.StatusCreated),
	)
}

func (c *ExpenseController) List(ctx *gin.Context) {
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

	var (
		view   model.ExpenseView
		status *string
	)

	switch ctx.Query("view") {
	case "personal":
		view = model.ExpenseViewPersonal
	case "approval_queue":
		view = model.ExpenseViewApprovalQueue
	default:
		view = model.ExpenseViewPersonal
	}

	statusQuery := ctx.Query("status")
	_, err = entity.ParseExpenseStatus(statusQuery)
	if err == nil {
		status = &statusQuery
	}

	autoApproved := ctx.Query("auto_approved") == "true"

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil || limit < 0 {
		offset = 0
	}

	request := &model.ListExpenseRequest{
		UserID:       userID,
		UserRole:     claims.Role,
		View:         view,
		Status:       status,
		AutoApproved: autoApproved,
		Limit:        limit,
		Offset:       offset,
	}
	res, total, err := c.expenseUsecase.List(ctx.Request.Context(), request)
	if err != nil {
		LogWarn(ctx, c.log, "failed to get expenses", err)
		ctx.Error(err)
		return
	}

	meta := model.MetaWithPage{
		Limit:      limit,
		Offset:     offset,
		Total:      total,
		HTTPStatus: http.StatusOK,
	}
	ctx.JSON(
		http.StatusOK,
		model.NewSuccessListResponse(res, meta),
	)
}

func (c *ExpenseController) Get(ctx *gin.Context) {
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

	res, err := c.expenseUsecase.FindByID(ctx.Request.Context(), &model.GetExpenseRequest{
		ID:       id,
		UserID:   userID,
		UserRole: claims.Role,
	})
	if err != nil {
		LogWarn(ctx, c.log, "failed to get expense", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		model.NewSuccessResponse(res, http.StatusOK),
	)
}
