package entity_test

import (
	"expense-management-system/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpense_RequiresApproval(t *testing.T) {
	tests := []struct {
		name    string
		model   *entity.Expense
		wantRes bool
	}{
		{
			name:    "nil model",
			model:   nil,
			wantRes: false,
		},
		{
			name: "amount greater than threshold",
			model: &entity.Expense{
				Amount: 2500000,
			},
			wantRes: true,
		},
		{
			name: "amount equal to threshold",
			model: &entity.Expense{
				Amount: 1000000,
			},
			wantRes: true,
		},
		{
			name: "amount less than threshold",
			model: &entity.Expense{
				Amount: 15000,
			},
			wantRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.model.RequiresApproval()

			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestExpense_AutoApproved(t *testing.T) {
	tests := []struct {
		name    string
		model   *entity.Expense
		wantRes bool
	}{
		{
			name:    "nil model",
			model:   nil,
			wantRes: false,
		},
		{
			name: "amount greater than threshold",
			model: &entity.Expense{
				Amount: 2500000,
			},
			wantRes: false,
		},
		{
			name: "amount equal to threshold",
			model: &entity.Expense{
				Amount: 1000000,
			},
			wantRes: false,
		},
		{
			name: "amount less than threshold",
			model: &entity.Expense{
				Amount: 15000,
			},
			wantRes: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.model.AutoApproved()

			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestExpense_GetKey(t *testing.T) {
	tests := []struct {
		name    string
		model   *entity.Expense
		wantRes string
	}{
		{
			name:    "nil model",
			model:   nil,
			wantRes: "",
		},
		{
			name: "success",
			model: &entity.Expense{
				ID: 120,
			},
			wantRes: "EXP-00000003C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.model.GetKey()

			assert.Equal(t, tt.wantRes, res)
		})
	}
}
