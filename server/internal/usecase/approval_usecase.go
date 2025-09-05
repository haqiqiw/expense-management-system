package usecase

import (
	"context"
	"expense-management-system/internal/db"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/messaging"
	"expense-management-system/internal/model"
	"fmt"

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

		idemKey = expense.GetKey()
		expAmount = expense.Amount

		approval := &entity.Approval{
			ExpenseID:  req.ID,
			ApproverID: req.UserID,
			Status:     approvalStatus,
			Notes:      req.Notes,
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
		// publish after commit to avoid event emitted but underlying DB change never committed
		// if publish fails, just log, since the DB is already committed
		// in a real system, we can handle this with outbox pattern or retry mechanism
		event := model.ExpenseApprovedEvent{
			ID:             req.ID,
			UserID:         req.UserID,
			Amount:         expAmount,
			IdempotencyKey: idemKey,
		}
		err = c.expenseApprovedProducer.Send(&event)
		if err != nil {
			c.log.Error(
				fmt.Sprintf("failed to send expense-approved event for id (%d) = %s", event.ID, err.Error()),
				zap.Any("event", event),
				zap.Strings("tags", []string{"approval", "update", "send-event", "expense-approved"}),
			)
		}
	}

	return nil
}
