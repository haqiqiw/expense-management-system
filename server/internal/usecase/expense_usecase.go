package usecase

import (
	"context"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/messaging"
	"expense-management-system/internal/model"
	"expense-management-system/internal/model/serializer"
	"fmt"

	"go.uber.org/zap"
)

type expenseUsecase struct {
	log                     *zap.Logger
	expenseRepository       ExpenseRepository
	expenseApprovedProducer *messaging.ExpenseApprovedProducer
}

func NewExpenseUsecase(log *zap.Logger, expenseRepository ExpenseRepository,
	expenseApprovedProducer *messaging.ExpenseApprovedProducer) ExpenseUsecase {
	return &expenseUsecase{
		log:                     log,
		expenseRepository:       expenseRepository,
		expenseApprovedProducer: expenseApprovedProducer,
	}
}

func (c *expenseUsecase) Create(ctx context.Context, req *model.CreateExpenseRequest) (*model.ExpenseCreateResponse, error) {
	if req.AmountIDR < entity.MinExpenseAmount {
		return nil, model.ErrExpenseMinAmount
	} else if req.AmountIDR > entity.MaxExpenseAmount {
		return nil, model.ErrExpenseMaxAmount
	}

	var status entity.ExpenseStatus
	if req.AmountIDR >= entity.ApprovalThresholdAmount {
		status = entity.ExpenseStatusAwaitingApproval
	} else {
		status = entity.ExpenseStatusApproved
	}

	expense := &entity.Expense{
		UserID:      req.UserID,
		Amount:      req.AmountIDR,
		Description: req.Description,
		ReceiptURL:  req.ReceiptURL,
		Status:      status,
	}

	err := c.expenseRepository.Create(ctx, expense)
	if err != nil {
		return nil, fmt.Errorf("failed to create expense = %w", err)
	}

	// publish after commit to avoid event emitted but underlying DB change never committed
	// if publish fails, just log, since the DB is already committed
	// in a real system, we can handle this with outbox pattern or retry mechanism
	if expense.Status == entity.ExpenseStatusApproved {
		event := model.ExpenseApprovedEvent{
			ID:             expense.ID,
			UserID:         expense.UserID,
			Amount:         expense.Amount,
			IdempotencyKey: expense.GetKey(),
		}
		err = c.expenseApprovedProducer.Send(&event)
		if err != nil {
			c.log.Error(
				fmt.Sprintf("failed to send expense-approved event for id (%d) = %s", event.ID, err.Error()),
				zap.Any("event", event),
				zap.Strings("tags", []string{"expense", "create", "send-event", "expense-approved"}),
			)
		}
	}

	return serializer.ExpenseToCreateResponse(expense), nil
}

func (c *expenseUsecase) List(ctx context.Context, req *model.ListExpenseRequest) ([]model.ExpenseWithUserResponse, int, error) {
	if req.View == model.ExpenseViewApprovalQueue && req.UserRole != string(entity.UserRoleManager) {
		return []model.ExpenseWithUserResponse{}, 0, model.ErrForbidden
	}

	expenses, total, err := c.expenseRepository.List(ctx, req)
	if err != nil {
		return []model.ExpenseWithUserResponse{}, 0, fmt.Errorf("failed to get expenses = %w", err)
	}

	if len(expenses) == 0 {
		return []model.ExpenseWithUserResponse{}, 0, nil
	}

	return serializer.ListExpenseWithUserToResponse(expenses), total, nil
}

func (c *expenseUsecase) FindByID(ctx context.Context, req *model.GetExpenseRequest) (*model.ExpenseDetailResponse, error) {
	expense, err := c.expenseRepository.FindDetailByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find expense by id (%d) = %w", req.ID, err)
	}

	if expense == nil {
		return nil, model.ErrExpenseNotFound
	}

	if req.UserID != expense.UserID && req.UserRole != string(entity.UserRoleManager) {
		return nil, model.ErrForbidden
	}

	return serializer.ExpenseDetailToResponse(expense), nil
}
