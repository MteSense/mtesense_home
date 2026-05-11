export type ThemeName = 'light' | 'dark'
export type SearchEngineId = 'google' | 'bing' | 'baidu'

export interface User {
  id: number
  username: string
  role: string
}

export interface NavGroup {
  id: number
  title: string
  sortOrder: number
  visible: boolean
  links: NavLink[]
}

export interface NavLink {
  id: number
  groupId: number
  title: string
  url: string
  icon: string
  iconType: 'emoji' | 'text' | 'image'
  description: string
  sortOrder: number
  visible: boolean
  openInNewTab: boolean
}

export interface AppearanceSettings {
  siteTitle: string
  subtitle: string
  backgroundImage: string
  defaultTheme: ThemeName
  cardOpacity: number
  blurStrength: number
}

export interface SearchSettings {
  defaultSearchEngine: SearchEngineId
  enabledSearchEngines: SearchEngineId[]
}

export interface PublicSettings {
  appearance: AppearanceSettings
  search: SearchSettings
}

export interface ApiEnvelope<T> {
  data?: T
  error?: string
}
