package messaging

import (
	"context"
	"encoding/json"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type ExpenseApprovedHandler struct {
	log                     *zap.Logger
	paymentProcessorUsecase usecase.PaymentProcessorUsecase
}

func NewExpenseApprovedHandler(log *zap.Logger,
	paymentProcessorUsecase usecase.PaymentProcessorUsecase) *ExpenseApprovedHandler {
	return &ExpenseApprovedHandler{
		log:                     log,
		paymentProcessorUsecase: paymentProcessorUsecase,
	}
}

func (c *ExpenseApprovedHandler) Consume(ctx context.Context, message *kafka.Message) error {
	c.log.Info(
		fmt.Sprintf("processing event for %s with key %s", message.TopicPartition.String(), string(message.Key)),
		zap.Any("event", string(message.Value)),
	)

	event := new(model.ExpenseApprovedEvent)
	err := json.Unmarshal(message.Value, &event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event for %s with key %s = %w", message.TopicPartition.String(), string(message.Key), err)
	}

	req := &model.PaymentProcessorRequest{
		ID:             event.ID,
		UserID:         event.UserID,
		Amount:         event.Amount,
		IdempotencyKey: event.IdempotencyKey,
	}
	err = c.paymentProcessorUsecase.Execute(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to execute payment for %s with key %s = %w", message.TopicPartition.String(), string(message.Key), err)
	}

	c.log.Info(
		fmt.Sprintf("successfuly proceed event for %s with key %s", message.TopicPartition.String(), string(message.Key)),
		zap.Any("event", string(message.Value)),
	)

	return nil
}
