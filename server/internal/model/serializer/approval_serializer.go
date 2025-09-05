package serializer

import (
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"time"
)

func ApprovalDetailResponse(a *entity.ApprovalDetail) *model.ApprovalDetailResponse {
	return &model.ApprovalDetailResponse{
		ID:            a.ID,
		ApproverID:    a.ApproverID,
		ApproverEmail: a.ApproverEmail,
		ApproverName:  a.ApproverName,
		Status:        string(a.Status),
		Notes:         a.Notes,
		CreatedAt:     a.CreatedAt.UTC().Format(time.RFC3339),
	}
}
