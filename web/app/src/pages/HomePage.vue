<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { Settings } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import NavCard from '../components/NavCard.vue'
import SearchBar from '../components/SearchBar.vue'
import ThemeToggle from '../components/ThemeToggle.vue'
import LanguageToggle from '../components/LanguageToggle.vue'
import { searchEngines } from '../api/searchEngines'
import type { SearchEngineId } from '../api/types'
import { useNavigationStore } from '../stores/navigation'
import { useSettingsStore } from '../stores/settings'
import { useUiStore } from '../stores/ui'

const navigation = useNavigationStore()
const settings = useSettingsStore()
const ui = useUiStore()
const { t, locale } = useI18n()
const query = ref('')
const now = ref(new Date())
let timer = 0

const enabledEngines = computed(() => settings.settings.search.enabledSearchEngines)
const filteredGroups = computed(() => {
  const keyword = query.value.trim().toLowerCase()
  if (!keyword) return navigation.groups
  return navigation.groups
    .map(group => ({
      ...group,
      links: group.links.filter(link =>
        [link.title, link.url, link.description].some(value => value.toLowerCase().includes(keyword))
      )
    }))
    .filter(group => group.links.length > 0)
})

const backgroundStyle = computed(() => {
  const image = settings.settings.appearance.backgroundImage
  return image ? { backgroundImage: `url(${image})` } : {}
})

function submitSearch() {
  const keyword = query.value.trim()
  if (!keyword) return
  const engine = ui.searchEngine
  window.open(searchEngines[engine].url(keyword), '_blank', 'noreferrer')
}

onMounted(async () => {
  await Promise.all([navigation.loadPublic(), settings.load()])
  const preferred = ui.searchEngine
  const fallback = settings.settings.search.defaultSearchEngine
  ui.setSearchEngine(enabledEngines.value.includes(preferred) ? preferred : fallback)
  ui.setTheme(ui.theme || settings.settings.appearance.defaultTheme)
  locale.value = ui.locale
  timer = window.setInterval(() => (now.value = new Date()), 1000)
})

onUnmounted(() => window.clearInterval(timer))

watch(enabledEngines, engines => {
  if (!engines.includes(ui.searchEngine)) {
    ui.setSearchEngine((engines[0] || 'google') as SearchEngineId)
  }
})
</script>

<template>
  <main class="home-page" :class="ui.theme" :style="backgroundStyle">
    <div class="window-bar">
      <span></span><span></span><span></span>
    </div>

    <section class="home-shell">
      <header class="home-header">
        <div class="title-block">
          <h1>{{ settings.settings.appearance.siteTitle }}</h1>
          <div class="clock">
            <strong>{{ now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}</strong>
            <span>{{ now.toLocaleDateString() }}</span>
          </div>
        </div>
        <div class="home-actions">
          <RouterLink class="icon-button" to="/admin/login" :title="t('admin')"><Settings :size="18" /></RouterLink>
          <LanguageToggle />
          <ThemeToggle />
        </div>
      </header>

      <SearchBar
        v-model:query="query"
        :engine="ui.searchEngine"
        :enabled-engines="enabledEngines"
        :placeholder="t('searchPlaceholder')"
        @update:engine="ui.setSearchEngine"
        @submit="submitSearch"
      />

      <section v-for="group in filteredGroups" :key="group.id" class="nav-section">
        <h2>{{ group.title }}</h2>
        <div class="nav-grid">
          <NavCard v-for="link in group.links" :key="link.id" :link="link" />
        </div>
      </section>
    </section>
  </main>
</template>
