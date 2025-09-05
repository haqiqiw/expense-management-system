package serializer_test

import (
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"expense-management-system/internal/model/serializer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApprovalSerializer_ApprovalDetailResponse(t *testing.T) {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)
	notes := "dummy notes"

	tests := []struct {
		name    string
		param   *entity.ApprovalDetail
		wantRes *model.ApprovalDetailResponse
	}{
		{
			name: "success",
			param: &entity.ApprovalDetail{
				ID:            1,
				ApproverID:    1,
				ApproverEmail: "john@mail.com",
				ApproverName:  "John Doe",
				Status:        entity.ApprovalStatusApproved,
				Notes:         &notes,
				CreatedAt:     now,
			},
			wantRes: &model.ApprovalDetailResponse{
				ID:            1,
				ApproverID:    1,
				ApproverEmail: "john@mail.com",
				ApproverName:  "John Doe",
				Status:        "approved",
				Notes:         &notes,
				CreatedAt:     now.Format(time.RFC3339),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := serializer.ApprovalDetailResponse(tt.param)

			assert.Equal(t, tt.wantRes, res)
		})
	}
}
