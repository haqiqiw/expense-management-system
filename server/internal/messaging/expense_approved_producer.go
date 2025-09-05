package messaging

import (
	"expense-management-system/internal/model"

	"go.uber.org/zap"
)

type ExpenseApprovedProducer struct {
	Producer[*model.ExpenseApprovedEvent]
}

func NewExpenseApprovedProducer(logger *zap.Logger, kProducer KafkaProducer, topic string) *ExpenseApprovedProducer {
	return &ExpenseApprovedProducer{
		Producer: &producer[*model.ExpenseApprovedEvent]{
			Producer: kProducer,
			Topic:    topic,
			Log:      logger,
		},
	}
}
