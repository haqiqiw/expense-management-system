package usecase

import (
	"context"
	"expense-management-system/internal/db"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"expense-management-system/internal/storage"
	"fmt"
	"time"

	"go.uber.org/zap"
)

const (
	prefixLockKey = "expense-payment:lock:"
	lockValue     = "lock"
)

type paymentProcessorUsecase struct {
	log                      *zap.Logger
	redisClient              storage.RedisClient
	tx                       db.Transactioner
	expenseRepository        ExpenseRepository
	paymentPartnerRepository PaymentPartnerRepository
	paymentLockDuration      int
}

func NewPaymentProcessorUsecase(log *zap.Logger, redisClient storage.RedisClient, tx db.Transactioner,
	expenseRepository ExpenseRepository, paymentPartnerRepository PaymentPartnerRepository,
	paymentLockDuration int) PaymentProcessorUsecase {
	return &paymentProcessorUsecase{
		log:                      log,
		redisClient:              redisClient,
		tx:                       tx,
		expenseRepository:        expenseRepository,
		paymentPartnerRepository: paymentPartnerRepository,
		paymentLockDuration:      paymentLockDuration,
	}
}

func (c *paymentProcessorUsecase) Execute(ctx context.Context, req *model.PaymentProcessorRequest) error {
	expense, err := c.expenseRepository.FindByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to find expense by id (%d) = %w", req.ID, err)
	}

	if expense == nil {
		c.log.Info(
			fmt.Sprintf("expense with id (%d) is not found", req.ID),
			zap.Strings("tags", []string{"payment-processor", "execute", "find"}),
		)
		return nil
	}

	if expense.Status != entity.ExpenseStatusApproved {
		c.log.Info(
			fmt.Sprintf("invalid expense status for id (%d) = %s", req.ID, expense.Status),
			zap.Strings("tags", []string{"payment-processor", "execute", "invalid-status"}),
		)
		return nil
	}

	lockKey := fmt.Sprintf("%s%d", prefixLockKey, req.ID)
	locked, err := c.redisClient.SetNX(
		ctx,
		lockKey,
		lockValue,
		time.Second*time.Duration(c.paymentLockDuration),
	).Result()
	if err != nil {
		return fmt.Errorf("failed to set lock for expense id (%d) = %w", req.ID, err)
	}

	if !locked {
		c.log.Info(
			fmt.Sprintf("expense with id (%d) is still locked by another processs", req.ID),
			zap.Strings("tags", []string{"payment-processor", "execute", "lock"}),
		)
		return nil
	}
	defer c.redisClient.Del(ctx, lockKey)

	err = c.tx.Do(ctx, func(exec db.Executor) error {
		// update the expense to complete, then call the partner
		// if the partner call fails, we can still rollback the expense update
		// even though, there’s a chance both succeed but the transaction commit fails
		// in that case, the partner already processed the payment
		//
		// since the partner guarantees idempotency, next retries will succeed
		// and indicate that the payment was already executed, avoiding double payment
		// afterwards, we can safely retry updating the expense to completed
		err = c.expenseRepository.CompleteByIDTx(ctx, exec, req.ID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to update expense for id (%d) = %w", req.ID, err)
		}

		partnerReq := &model.PaymentPartnerRequest{
			Amount:     req.Amount,
			ExternalID: req.IdempotencyKey,
		}
		_, err = c.paymentPartnerRepository.Execute(ctx, partnerReq)
		if err != nil {
			return fmt.Errorf("failed to call partner for expense id (%d) = %w", req.ID, err)
		}

		return nil
	})

	return err
}
