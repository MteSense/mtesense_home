import { defineStore } from 'pinia'
import { api } from '../api/client'
import type { PublicSettings } from '../api/types'

export const defaultSettings: PublicSettings = {
  appearance: {
    siteTitle: 'MteSense',
    subtitle: 'Personal navigation',
    backgroundImage: '',
    defaultTheme: 'dark',
    cardOpacity: 0.34,
    blurStrength: 18
  },
  search: {
    defaultSearchEngine: 'google',
    enabledSearchEngines: ['google', 'bing', 'baidu']
  }
}

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: structuredClone(defaultSettings) as PublicSettings,
    loading: false,
    error: ''
  }),
  actions: {
    async load() {
      this.loading = true
      this.error = ''
      try {
        this.settings = await api.settings()
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to load settings'
      } finally {
        this.loading = false
      }
    },
    async save(payload: PublicSettings) {
      this.settings = await api.saveSettings(payload)
    }
  }
})
