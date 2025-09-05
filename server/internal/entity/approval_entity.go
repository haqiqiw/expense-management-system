package entity

import "time"

type ApprovalStatus string

const (
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
)

type Approval struct {
	ID         uint64         `db:"id"`
	ExpenseID  uint64         `db:"expense_id"`
	ApproverID uint64         `db:"approver_id"`
	Status     ApprovalStatus `db:"status"`
	Notes      *string        `db:"notes"`
	CreatedAt  time.Time      `db:"created_at"`
}

type ApprovalDetail struct {
	ID            uint64         `db:"id"`
	ApproverID    uint64         `db:"approver_id"`
	ApproverEmail string         `db:"approver_email"`
	ApproverName  string         `db:"approver_name"`
	Status        ApprovalStatus `db:"status"`
	Notes         *string        `db:"notes"`
	CreatedAt     time.Time      `db:"created_at"`
}
