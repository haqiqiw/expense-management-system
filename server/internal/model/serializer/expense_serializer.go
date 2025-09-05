package serializer

import (
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"time"
)

func ExpenseToCreateResponse(e *entity.Expense) *model.ExpenseCreateResponse {
	return &model.ExpenseCreateResponse{
		ID:               e.ID,
		AmountIDR:        e.Amount,
		Description:      e.Description,
		ReceiptURL:       e.ReceiptURL,
		Status:           string(e.Status),
		RequiresApproval: e.RequiresApproval(),
		AutoApproved:     e.AutoApproved(),
		CreatedAt:        e.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func ExpenseWithUserToResponse(e *entity.ExpenseWithUser) *model.ExpenseWithUserResponse {
	return &model.ExpenseWithUserResponse{
		ID:               e.ID,
		AmountIDR:        e.Amount,
		Description:      e.Description,
		ReceiptURL:       e.ReceiptURL,
		Status:           string(e.Status),
		RequiresApproval: e.RequiresApproval(),
		AutoApproved:     e.AutoApproved(),
		CreatedAt:        e.CreatedAt.UTC().Format(time.RFC3339),
		User:             *UserSimpleToResponse(&e.User),
	}
}

func ListExpenseWithUserToResponse(expenses []entity.ExpenseWithUser) []model.ExpenseWithUserResponse {
	res := make([]model.ExpenseWithUserResponse, len(expenses))

	for i, e := range expenses {
		res[i] = *ExpenseWithUserToResponse(&e)
	}

	return res
}

func ExpenseDetailToResponse(e *entity.ExpenseDetail) *model.ExpenseDetailResponse {
	var approval *model.ApprovalDetailResponse
	if e.Approval != nil {
		approval = ApprovalDetailResponse(e.Approval)
	}

	return &model.ExpenseDetailResponse{
		ID:               e.ID,
		AmountIDR:        e.Amount,
		Description:      e.Description,
		ReceiptURL:       e.ReceiptURL,
		Status:           string(e.Status),
		RequiresApproval: e.RequiresApproval(),
		AutoApproved:     e.AutoApproved(),
		CreatedAt:        e.CreatedAt.UTC().Format(time.RFC3339),
		User:             *UserSimpleToResponse(&e.User),
		Approval:         approval,
	}
}
