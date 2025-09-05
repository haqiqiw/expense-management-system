package model

type ApprovalExpenseRequest struct {
	ID       uint64  `json:"id"`
	Notes    *string `json:"notes"`
	UserID   uint64  `json:"user_id"`   // current user id
	UserRole string  `json:"user_role"` // current user role
}

type ApprovalDetailResponse struct {
	ID            uint64  `json:"id"`
	ApproverID    uint64  `json:"approver_id"`
	ApproverEmail string  `json:"approver_email"`
	ApproverName  string  `json:"approver_name"`
	Status        string  `json:"status"`
	Notes         *string `json:"notes"`
	CreatedAt     string  `json:"created_at"`
}
