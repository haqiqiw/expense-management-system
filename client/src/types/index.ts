export interface User {
  id: number
  email: string
  name: string
  role: 'manager' | 'employee'
  created_at: string
}

export interface Expense {
  id: number
  amount_idr: number
  description: string
  receipt_url: string | null
  status: 'awaiting_approval' | 'approved' | 'rejected' | 'completed'
  requires_approval: boolean
  auto_approved: boolean
  created_at: string
  user: {
    id: number
    email: string
    name: string
  }
}

export interface ExpenseDetail extends Expense {
  processed_at: string | null
  approval: ApprovalDetail | null
}

export interface ApprovalDetail {
  id: number
  approver_name: string
  status: 'approved' | 'rejected'
  notes: string | null
  created_at: string
}
export type ExpenseFiltersStatus =
  | 'awaiting_approval'
  | 'approved'
  | 'rejected'
  | 'completed'
  | null

export interface ExpenseFilters {
  status: ExpenseFiltersStatus
  auto_approved: boolean
  limit: number
  offset: number
}
