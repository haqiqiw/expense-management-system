package messaging_test

import (
	"context"
	"encoding/json"
	"errors"
	"expense-management-system/internal/delivery/messaging"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestExpenseApprovedHandler_Consume(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	topic := "expense-approved"

	validMsg := func() *kafka.Message {
		event := &model.ExpenseApprovedEvent{
			ID:             8,
			UserID:         2,
			Amount:         17000,
			IdempotencyKey: "EXP-00123ABC",
		}
		data, _ := json.Marshal(event)
		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: data,
			Key:   []byte(event.GetID()),
		}

		return msg
	}

	tests := []struct {
		name       string
		message    *kafka.Message
		mockFunc   func(t *mocks.PaymentProcessorUsecase)
		wantErrMsg string
	}{
		{
			name: "error on unrmarshal",
			message: func() *kafka.Message {
				data, _ := json.Marshal("dummy")
				msg := &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &topic,
						Partition: kafka.PartitionAny,
					},
					Value: data,
					Key:   []byte("1"),
				}

				return msg
			}(),
			mockFunc:   func(t *mocks.PaymentProcessorUsecase) {},
			wantErrMsg: "failed to unmarshal event for expense-approved",
		},
		{
			name:    "error on execute",
			message: validMsg(),
			mockFunc: func(t *mocks.PaymentProcessorUsecase) {
				t.On("Execute", mock.Anything, mock.Anything).Return(errors.New("something error"))
			},
			wantErrMsg: "something error",
		},
		{
			name:    "success",
			message: validMsg(),
			mockFunc: func(t *mocks.PaymentProcessorUsecase) {
				t.On("Execute", mock.Anything, mock.Anything).Return(nil)
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := mocks.NewPaymentProcessorUsecase(t)
			handler := messaging.NewExpenseApprovedHandler(logger, uc)
			tt.mockFunc(uc)

			err := handler.Consume(ctx, tt.message)

			if tt.wantErrMsg != "" {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
