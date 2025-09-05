package model_test

import (
	"expense-management-system/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpenseApprovedEvent_GetID(t *testing.T) {
	tests := []struct {
		name   string
		event  *model.ExpenseApprovedEvent
		wantID string
	}{
		{
			name: "success",
			event: &model.ExpenseApprovedEvent{
				ID:             8,
				UserID:         2,
				Amount:         17000,
				IdempotencyKey: "EXP-00123ABC",
			},
			wantID: "expense-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.event.GetID()

			assert.Equal(t, tt.wantID, id)
		})
	}
}
