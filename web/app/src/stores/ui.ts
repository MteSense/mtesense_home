import { defineStore } from 'pinia'
import type { SearchEngineId, ThemeName } from '../api/types'

export const useUiStore = defineStore('ui', {
  state: () => ({
    theme: (localStorage.getItem('mtesense_theme') as ThemeName) || 'dark',
    locale: localStorage.getItem('mtesense_locale') || 'zh',
    searchEngine: (localStorage.getItem('mtesense_search_engine') as SearchEngineId) || 'google'
  }),
  actions: {
    setTheme(theme: ThemeName) {
      this.theme = theme
      localStorage.setItem('mtesense_theme', theme)
      document.documentElement.dataset.theme = theme
    },
    toggleTheme() {
      this.setTheme(this.theme === 'dark' ? 'light' : 'dark')
    },
    setLocale(locale: string) {
      this.locale = locale
      localStorage.setItem('mtesense_locale', locale)
    },
    setSearchEngine(engine: SearchEngineId) {
      this.searchEngine = engine
      localStorage.setItem('mtesense_search_engine', engine)
    }
  }
})
