package usecase

import (
	"context"
	"expense-management-system/internal/db"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/messaging"
	"expense-management-system/internal/metrics"
	"expense-management-system/internal/model"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type approvalUsecase struct {
	log                     *zap.Logger
	tx                      db.Transactioner
	approvalRepository      ApprovalRepository
	expenseRepository       ExpenseRepository
	expenseApprovedProducer *messaging.ExpenseApprovedProducer
}

func NewApprovalUsecase(log *zap.Logger, tx db.Transactioner, approvalRepository ApprovalRepository,
	expenseRepository ExpenseRepository, expenseApprovedProducer *messaging.ExpenseApprovedProducer) ApprovalUsecase {
	return &approvalUsecase{
		log:                     log,
		tx:                      tx,
		approvalRepository:      approvalRepository,
		expenseRepository:       expenseRepository,
		expenseApprovedProducer: expenseApprovedProducer,
	}
}

func (c *approvalUsecase) Approve(ctx context.Context, req *model.ApprovalExpenseRequest) error {
	return c.updateApproval(ctx, req, entity.ApprovalStatusApproved)
}

func (c *approvalUsecase) Reject(ctx context.Context, req *model.ApprovalExpenseRequest) error {
	return c.updateApproval(ctx, req, entity.ApprovalStatusRejected)
}

func (c *approvalUsecase) updateApproval(ctx context.Context, req *model.ApprovalExpenseRequest, approvalStatus entity.ApprovalStatus) error {
	if req.UserRole != string(entity.UserRoleManager) {
		return model.ErrForbidden
	}

	var (
		idemKey   string
		expAmount uint64
	)

	err := c.tx.Do(ctx, func(exec db.Executor) error {
		expense, txErr := c.expenseRepository.FindByIDWithLock(ctx, exec, req.ID)
		if txErr != nil {
			return fmt.Errorf("failed to find expense by id (%d) with lock = %w", req.ID, txErr)
		}

		if expense == nil {
			return model.ErrExpenseNotFound
		}
		if expense.UserID == req.UserID {
			return model.ErrForbidden
		}
		if expense.Status != entity.ExpenseStatusAwaitingApproval {
			return model.ErrExpenseAlreadyProcessed
		}
		if !expense.RequiresApproval() {
			return model.ErrExpenseNotRequireApproval
		}

		idemKey = expense.GetKey()
		expAmount = expense.Amount

		var notes *string
		if req.Notes != nil {
			n := strings.TrimSpace(*req.Notes)
			notes = &n
		}

		approval := &entity.Approval{
			ExpenseID:  req.ID,
			ApproverID: req.UserID,
			Status:     approvalStatus,
			Notes:      notes,
		}
		txErr = c.approvalRepository.CreateTx(ctx, exec, approval)
		if txErr != nil {
			return fmt.Errorf("failed to to create approval for expense id (%d) = %w", req.ID, txErr)
		}

		var expenseStatus entity.ExpenseStatus
		switch approvalStatus {
		case entity.ApprovalStatusApproved:
			expenseStatus = entity.ExpenseStatusApproved
		case entity.ApprovalStatusRejected:
			expenseStatus = entity.ExpenseStatusRejected
		}

		txErr = c.expenseRepository.UpdateStatusByIDTx(ctx, exec, req.ID, expenseStatus)
		if txErr != nil {
			return fmt.Errorf("failed to to update expense for id (%d) = %w", req.ID, txErr)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if approvalStatus == entity.ApprovalStatusApproved {
		event := model.ExpenseApprovedEvent{
			ID:             req.ID,
			UserID:         req.UserID,
			Amount:         expAmount,
			IdempotencyKey: idemKey,
		}

		// publish after the expense is successfully updated
		// if publishing fails, log the error and send a metric to notify us
		// later, we can run a script to retry sending failed events
		// a more robust solution would be to use the outbox pattern
		eventStatus := "success"
		err = c.expenseApprovedProducer.Send(&event)
		if err != nil {
			eventStatus = "fail"
			c.log.Error(
				fmt.Sprintf("failed to send expense-approved event for id (%d) = %s", event.ID, err.Error()),
				zap.Any("event", event),
				zap.Strings("tags", []string{"approval", "update", "send-event", "expense-approved"}),
			)
		}
		metrics.IncrementEvent(metrics.EventPusblishExpenseApprove, eventStatus)
	}

	return nil
}
