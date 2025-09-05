package serializer_test

import (
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"expense-management-system/internal/model/serializer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpenseSerializer_ExpenseToCreateResponse(t *testing.T) {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)
	description := "dummy description"
	receipt := "https://example.com/receipt.jpg"

	tests := []struct {
		name    string
		param   *entity.Expense
		wantRes *model.ExpenseCreateResponse
	}{
		{
			name: "success",
			param: &entity.Expense{
				ID:          1,
				UserID:      1,
				Amount:      10000,
				Description: description,
				ReceiptURL:  &receipt,
				Status:      entity.ExpenseStatusApproved,
				CreatedAt:   now,
			},
			wantRes: &model.ExpenseCreateResponse{
				ID:               1,
				AmountIDR:        10000,
				Description:      description,
				ReceiptURL:       &receipt,
				Status:           "approved",
				RequiresApproval: false,
				AutoApproved:     true,
				CreatedAt:        now.Format(time.RFC3339),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := serializer.ExpenseToCreateResponse(tt.param)

			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestExpenseSerializer_ExpenseWithUserToResponse(t *testing.T) {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)
	description := "dummy description"
	receipt := "https://example.com/receipt.jpg"

	tests := []struct {
		name    string
		param   *entity.ExpenseWithUser
		wantRes *model.ExpenseWithUserResponse
	}{
		{
			name: "success",
			param: &entity.ExpenseWithUser{
				Expense: entity.Expense{
					ID:          1,
					UserID:      1,
					Amount:      10000,
					Description: description,
					ReceiptURL:  &receipt,
					Status:      entity.ExpenseStatusApproved,
					CreatedAt:   now,
				},
				User: entity.UserSimple{
					ID:    1,
					Email: "john@mail.com",
					Name:  "John Doe",
				},
			},
			wantRes: &model.ExpenseWithUserResponse{
				ID:               1,
				AmountIDR:        10000,
				Description:      description,
				ReceiptURL:       &receipt,
				Status:           "approved",
				RequiresApproval: false,
				AutoApproved:     true,
				CreatedAt:        now.Format(time.RFC3339),
				User: model.UserSimpleResponse{
					ID:    1,
					Email: "john@mail.com",
					Name:  "John Doe",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := serializer.ExpenseWithUserToResponse(tt.param)

			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestExpenseSerializer_ListExpenseWithUserToResponse(t *testing.T) {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)
	description := "dummy description"
	receipt := "https://example.com/receipt.jpg"

	tests := []struct {
		name    string
		param   []entity.ExpenseWithUser
		wantRes []model.ExpenseWithUserResponse
	}{
		{
			name: "success",
			param: []entity.ExpenseWithUser{
				{
					Expense: entity.Expense{
						ID:          1,
						UserID:      1,
						Amount:      10000,
						Description: description,
						ReceiptURL:  &receipt,
						Status:      entity.ExpenseStatusApproved,
						CreatedAt:   now,
					},
					User: entity.UserSimple{
						ID:    1,
						Email: "john@mail.com",
						Name:  "John Doe",
					},
				},
			},
			wantRes: []model.ExpenseWithUserResponse{
				{
					ID:               1,
					AmountIDR:        10000,
					Description:      description,
					ReceiptURL:       &receipt,
					Status:           "approved",
					RequiresApproval: false,
					AutoApproved:     true,
					CreatedAt:        now.Format(time.RFC3339),
					User: model.UserSimpleResponse{
						ID:    1,
						Email: "john@mail.com",
						Name:  "John Doe",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := serializer.ListExpenseWithUserToResponse(tt.param)

			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestExpenseSerializer_ExpenseDetailToResponse(t *testing.T) {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)
	description := "dummy description"
	receipt := "https://example.com/receipt.jpg"
	notes := "dummy notes"

	tests := []struct {
		name    string
		param   *entity.ExpenseDetail
		wantRes *model.ExpenseDetailResponse
	}{
		{
			name: "success with approval",
			param: &entity.ExpenseDetail{
				Expense: entity.Expense{
					ID:          1,
					UserID:      1,
					Amount:      10000,
					Description: description,
					ReceiptURL:  &receipt,
					Status:      entity.ExpenseStatusApproved,
					CreatedAt:   now,
				},
				User: entity.UserSimple{
					ID:    1,
					Email: "john@mail.com",
					Name:  "John Doe",
				},
				Approval: &entity.ApprovalDetail{
					ID:            1,
					ApproverID:    1,
					ApproverEmail: "john@mail.com",
					ApproverName:  "John Doe",
					Status:        entity.ApprovalStatusApproved,
					Notes:         &notes,
					CreatedAt:     now,
				},
			},
			wantRes: &model.ExpenseDetailResponse{
				ID:               1,
				AmountIDR:        10000,
				Description:      description,
				ReceiptURL:       &receipt,
				Status:           "approved",
				RequiresApproval: false,
				AutoApproved:     true,
				CreatedAt:        now.Format(time.RFC3339),
				User: model.UserSimpleResponse{
					ID:    1,
					Email: "john@mail.com",
					Name:  "John Doe",
				},
				Approval: &model.ApprovalDetailResponse{
					ID:            1,
					ApproverID:    1,
					ApproverEmail: "john@mail.com",
					ApproverName:  "John Doe",
					Status:        "approved",
					Notes:         &notes,
					CreatedAt:     now.Format(time.RFC3339),
				},
			},
		},
		{
			name: "success without approval",
			param: &entity.ExpenseDetail{
				Expense: entity.Expense{
					ID:          1,
					UserID:      1,
					Amount:      10000,
					Description: description,
					ReceiptURL:  &receipt,
					Status:      entity.ExpenseStatusApproved,
					CreatedAt:   now,
				},
				User: entity.UserSimple{
					ID:    1,
					Email: "john@mail.com",
					Name:  "John Doe",
				},
				Approval: nil,
			},
			wantRes: &model.ExpenseDetailResponse{
				ID:               1,
				AmountIDR:        10000,
				Description:      description,
				ReceiptURL:       &receipt,
				Status:           "approved",
				RequiresApproval: false,
				AutoApproved:     true,
				CreatedAt:        now.Format(time.RFC3339),
				User: model.UserSimpleResponse{
					ID:    1,
					Email: "john@mail.com",
					Name:  "John Doe",
				},
				Approval: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := serializer.ExpenseDetailToResponse(tt.param)

			assert.Equal(t, tt.wantRes, res)
		})
	}
}
