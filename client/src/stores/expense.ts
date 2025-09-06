import { defineStore } from 'pinia'
import apiClient from '@/services/api'
import type { Expense, ExpenseFilters, ExpenseDetail } from '@/types'

interface ExpenseState {
  isLoading: boolean
  isCreateLoading: boolean
  expenses: Expense[]
  total: number
  filters: ExpenseFilters
  currentExpense: ExpenseDetail | null
}

export const useExpenseStore = defineStore('expense', {
  state: (): ExpenseState => ({
    isLoading: false,
    isCreateLoading: false,
    expenses: [],
    total: 0,
    filters: {
      status: null,
      auto_approved: false,
      limit: 5,
      offset: 0,
    },
    currentExpense: null,
  }),

  actions: {
    async fetchExpenses() {
      this.expenses = []
      this.isLoading = true

      try {
        const params = new URLSearchParams()
        params.append('view', 'personal')
        params.append('limit', this.filters.limit.toString())
        params.append('offset', this.filters.offset.toString())

        if (this.filters.status) {
          params.append('status', this.filters.status)
        } else if (this.filters.auto_approved) {
          params.append('auto_approved', 'true')
        }

        const response = await apiClient.get(`/expenses?${params.toString()}`)
        this.expenses = response.data.data
        this.total = response.data.meta.total
      } catch (error) {
        console.error('Failed to fetch current expenses:', error)
      } finally {
        this.isLoading = false
      }
    },

    setPage(page: number) {
      this.filters.offset = (page - 1) * this.filters.limit
      this.fetchExpenses()
    },

    setFilters(newFilters: Partial<ExpenseFilters>) {
      this.filters.offset = 0
      if (newFilters.auto_approved) {
        this.filters.status = null
        this.filters.auto_approved = true
      } else {
        this.filters.auto_approved = false
        this.filters.status = newFilters.status ?? null
      }
      this.fetchExpenses()
    },

    async createExpense(payload: {
      amount_idr: number
      description: string
      receipt_url: string | null
    }) {
      this.isCreateLoading = true
      try {
        await apiClient.post('/expenses', payload)

        await this.fetchExpenses()
      } catch (error) {
        console.error('Failed to create expense:', error)
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async fetchExpenseById(id: string) {
      this.isLoading = true
      this.currentExpense = null
      try {
        const response = await apiClient.get(`/expenses/${id}`)
        this.currentExpense = response.data.data
      } catch (error) {
        console.error('Failed to fetch expense detail:', error)
        throw error
      } finally {
        this.isLoading = false
      }
    },
  },
})
