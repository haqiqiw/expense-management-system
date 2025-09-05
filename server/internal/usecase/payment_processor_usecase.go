package usecase

import (
	"context"
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
	expenseRepository        ExpenseRepository
	paymentPartnerRepository PaymentPartnerRepository
	paymentLockDuration      int
}

func NewPaymentProcessorUsecase(log *zap.Logger, redisClient storage.RedisClient,
	expenseRepository ExpenseRepository, paymentPartnerRepository PaymentPartnerRepository,
	paymentLockDuration int) PaymentProcessorUsecase {
	return &paymentProcessorUsecase{
		log:                      log,
		redisClient:              redisClient,
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

	partnerReq := &model.PaymentPartnerRequest{
		Amount:     req.Amount,
		ExternalID: req.IdempotencyKey,
	}
	_, err = c.paymentPartnerRepository.Execute(ctx, partnerReq)
	if err != nil {
		return fmt.Errorf("failed to call partner for expense id (%d) = %w", req.ID, err)
	}

	// at this point, the partner already processed the payment successfully
	// if updating the expense status in DB fails, we return an error so the consumer will retry
	// on next retry, the partner will return success again because their API is idempotent,
	// and we’ll retry the DB update to ensures the "completed" status is persisted
	err = c.expenseRepository.UpdateStatusByID(ctx, req.ID, entity.ExpenseStatusCompleted, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update expense for id (%d) = %w", req.ID, err)
	}

	return nil
}
