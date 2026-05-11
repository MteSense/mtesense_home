import type { SearchEngineId } from './types'

export const searchEngines: Record<SearchEngineId, { label: string; url: (query: string) => string }> = {
  google: {
    label: 'Google',
    url: query => `https://www.google.com/search?q=${encodeURIComponent(query)}`
  },
  bing: {
    label: 'Bing',
    url: query => `https://www.bing.com/search?q=${encodeURIComponent(query)}`
  },
  baidu: {
    label: 'Baidu',
    url: query => `https://www.baidu.com/s?wd=${encodeURIComponent(query)}`
  }
}
