package messaging_test

import (
	"errors"
	"expense-management-system/internal/messaging"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type ExpenseApprovedProducerSuite struct {
	suite.Suite
	logger   *zap.Logger
	kafka    *mocks.KafkaProducer
	producer messaging.Producer[*model.ExpenseApprovedEvent]
	topic    string
}

func (s *ExpenseApprovedProducerSuite) SetupTest() {
	s.logger, _ = zap.NewDevelopment()
	s.kafka = mocks.NewKafkaProducer(s.T())
	s.topic = "expense-approved"
	s.producer = messaging.NewExpenseApprovedProducer(s.logger, s.kafka, s.topic)
}

func (s *ExpenseApprovedProducerSuite) TearDownTest() {
	s.kafka = mocks.NewKafkaProducer(s.T())
}

func (s *ExpenseApprovedProducerSuite) TestExpenseApprovedProducer_GetTopic() {
	t := s.producer.GetTopic()

	s.Equal("expense-approved", *t)
}

func (s *ExpenseApprovedProducerSuite) TestExpenseApprovedProducer_Send() {
	tests := []struct {
		name       string
		mockFunc   func(k *mocks.KafkaProducer)
		param      *model.ExpenseApprovedEvent
		wantErrMsg string
	}{
		{
			name: "error on produce",
			mockFunc: func(k *mocks.KafkaProducer) {
				k.On("Produce", mock.Anything, mock.Anything).
					Return(errors.New("something error"))
			},
			param: &model.ExpenseApprovedEvent{
				ID:             8,
				UserID:         2,
				Amount:         17000,
				IdempotencyKey: "EXP-00123ABC",
			},
			wantErrMsg: "failed to produce message for expense-approved = something error",
		},
		{
			name: "success",
			mockFunc: func(k *mocks.KafkaProducer) {
				k.On("Produce", mock.Anything, mock.Anything).Return(nil)
			},
			param: &model.ExpenseApprovedEvent{
				ID:             8,
				UserID:         2,
				Amount:         17000,
				IdempotencyKey: "EXP-00123ABC",
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.kafka = mocks.NewKafkaProducer(s.T())
			s.producer = messaging.NewExpenseApprovedProducer(s.logger, s.kafka, s.topic)
			tt.mockFunc(s.kafka)

			err := s.producer.Send(tt.param)

			if tt.wantErrMsg == "" {
				s.Nil(err)
			} else {
				s.Equal(tt.wantErrMsg, err.Error())
			}
		})
	}
}

func TestExpenseApprovedProducerSuite(t *testing.T) {
	suite.Run(t, new(ExpenseApprovedProducerSuite))
}
