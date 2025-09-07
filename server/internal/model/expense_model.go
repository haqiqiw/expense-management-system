package model

type ExpenseView string

const (
	ExpenseViewPersonal      ExpenseView = "personal"
	ExpenseViewApprovalQueue ExpenseView = "approval_queue" // manager only
)

type CreateExpenseRequest struct {
	UserID      uint64  `json:"user_id"` // current user id
	AmountIDR   uint64  `json:"amount_idr" validate:"required,number,gt=0"`
	Description string  `json:"description" validate:"required,max=255"`
	ReceiptURL  *string `json:"receipt_url" validate:"omitempty,url"`
}

type ListExpenseRequest struct {
	UserID       uint64      `json:"user_id"`   // current user id
	UserRole     string      `json:"user_role"` // current user role
	View         ExpenseView `json:"view"`
	Status       *string     `json:"status"`
	AutoApproved bool        `json:"auto_approved"` // flag to filter by amount
	Limit        int         `json:"limit"`
	Offset       int         `json:"offset"`
}

type GetExpenseRequest struct {
	ID       uint64 `json:"id"`
	UserID   uint64 `json:"user_id"`   // current user id
	UserRole string `json:"user_role"` // current user role
}

type ExpenseCreateResponse struct {
	ID               uint64  `json:"id"`
	AmountIDR        uint64  `json:"amount_idr"`
	Description      string  `json:"description"`
	ReceiptURL       *string `json:"receipt_url"`
	Status           string  `json:"status"`
	RequiresApproval bool    `json:"requires_approval"`
	AutoApproved     bool    `json:"auto_approved"`
	CreatedAt        string  `json:"created_at"`
}

type ExpenseWithUserResponse struct {
	ID               uint64             `json:"id"`
	AmountIDR        uint64             `json:"amount_idr"`
	Description      string             `json:"description"`
	ReceiptURL       *string            `json:"receipt_url"`
	Status           string             `json:"status"`
	RequiresApproval bool               `json:"requires_approval"`
	AutoApproved     bool               `json:"auto_approved"`
	CreatedAt        string             `json:"created_at"`
	User             UserSimpleResponse `json:"user"`
}

type ExpenseDetailResponse struct {
	ID               uint64                  `json:"id"`
	AmountIDR        uint64                  `json:"amount_idr"`
	Description      string                  `json:"description"`
	ReceiptURL       *string                 `json:"receipt_url"`
	Status           string                  `json:"status"`
	RequiresApproval bool                    `json:"requires_approval"`
	AutoApproved     bool                    `json:"auto_approved"`
	CreatedAt        string                  `json:"created_at"`
	ProcessedAt      *string                 `json:"processed_at"`
	User             UserSimpleResponse      `json:"user"`
	Approval         *ApprovalDetailResponse `json:"approval"`
}
