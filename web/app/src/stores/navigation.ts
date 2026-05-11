import { defineStore } from 'pinia'
import { api } from '../api/client'
import type { NavGroup, NavLink } from '../api/types'

export const useNavigationStore = defineStore('navigation', {
  state: () => ({
    groups: [] as NavGroup[],
    loading: false,
    error: ''
  }),
  actions: {
    async loadPublic() {
      await this.load(false)
    },
    async loadAdmin() {
      await this.load(true)
    },
    async load(admin: boolean) {
      this.loading = true
      this.error = ''
      try {
        this.groups = admin ? await api.adminNavigation() : await api.navigation()
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to load navigation'
      } finally {
        this.loading = false
      }
    },
    async saveGroup(group: Partial<NavGroup>) {
      if (group.id) {
        await api.updateGroup(group.id, group)
      } else {
        await api.createGroup(group)
      }
      await this.loadAdmin()
    },
    async deleteGroup(id: number) {
      await api.deleteGroup(id)
      await this.loadAdmin()
    },
    async saveLink(link: Partial<NavLink>) {
      if (link.id) {
        await api.updateLink(link.id, link)
      } else {
        await api.createLink(link)
      }
      await this.loadAdmin()
    },
    async deleteLink(id: number) {
      await api.deleteLink(id)
      await this.loadAdmin()
    }
  }
})
