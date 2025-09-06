import { defineStore } from 'pinia'
import apiClient from '@/services/api'
import type { Expense, ExpenseFilters } from '@/types'

interface ApprovalState {
  isLoading: boolean
  expenses: Expense[]
  total: number
  filters: ExpenseFilters
}

export const useApprovalStore = defineStore('approval', {
  state: (): ApprovalState => ({
    isLoading: false,
    expenses: [],
    total: 0,
    filters: {
      status: null,
      auto_approved: false,
      limit: 5,
      offset: 0,
    },
  }),

  actions: {
    async fetchExpenses() {
      this.expenses = []
      this.isLoading = true

      try {
        const params = new URLSearchParams()
        params.append('view', 'approval_queue')
        params.append('limit', this.filters.limit.toString())
        params.append('offset', this.filters.offset.toString())

        const response = await apiClient.get(`/expenses?${params.toString()}`)
        this.expenses = response.data.data
        this.total = response.data.meta.total
      } catch (error) {
        console.error('Failed to fetch approval expenses:', error)
      } finally {
        this.isLoading = false
      }
    },

    setPage(page: number) {
      this.filters.offset = (page - 1) * this.filters.limit
      this.fetchExpenses()
    },
  },
})
