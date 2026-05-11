import { defineStore } from 'pinia'
import { api, clearToken, setToken } from '../api/client'
import type { User } from '../api/types'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    loading: false,
    error: ''
  }),
  actions: {
    async login(username: string, password: string) {
      this.loading = true
      this.error = ''
      try {
        const result = await api.login(username, password)
        setToken(result.token)
        this.user = result.user
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Login failed'
        throw error
      } finally {
        this.loading = false
      }
    },
    logout() {
      clearToken()
      this.user = null
    }
  }
})
