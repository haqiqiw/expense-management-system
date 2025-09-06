import { defineStore } from 'pinia'
import apiClient, { setAuthToken } from '@/services/api'
import type { User } from '@/types'

interface AuthState {
  user: User | null
  token: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    token: localStorage.getItem('accessToken') || null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !!state.user,
    currentUser: (state) => state.user,
    userRole: (state) => state.user?.role,
  },

  actions: {
    async login(credentials: { email: string; password: string }): Promise<void> {
      try {
        const response = await apiClient.post('/auth/login', credentials)
        const token = response.data.data.access_token

        if (!token) {
          throw new Error('Login failed: No token received.')
        }

        this.token = token
        setAuthToken(token)

        await this.fetchUser()
      } catch (error) {
        this.logout()
        console.error('Login failed:', error)
        throw error
      }
    },

    async fetchUser(): Promise<void> {
      if (!this.token) return

      try {
        const response = await apiClient.get('/users/me')
        const user = response.data.data
        this.user = user
        localStorage.setItem('user', JSON.stringify(user))
      } catch (error) {
        this.logout()
        console.error('Failed to fetch user:', error)
        throw error
      }
    },

    async logout(): Promise<void> {
      try {
        const response = await apiClient.post('auth/logout')
        if (response.status == 200) {
          this.user = null
          this.token = null
          setAuthToken(null)
          localStorage.removeItem('user')
        } else {
          throw new Error('Logout failed: unexpected error')
        }
      } catch (error) {
        console.error('Logout failed:', error)
        throw error
      }
    },
  },
})
