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

const lunarMonths: Record<string, number> = {
  正月: 1,
  二月: 2,
  三月: 3,
  四月: 4,
  五月: 5,
  六月: 6,
  七月: 7,
  八月: 8,
  九月: 9,
  十月: 10,
  十一月: 11,
  冬月: 11,
  十二月: 12,
  腊月: 12
}
const lunarDays = [
  '',
  '初一',
  '初二',
  '初三',
  '初四',
  '初五',
  '初六',
  '初七',
  '初八',
  '初九',
  '初十',
  '十一',
  '十二',
  '十三',
  '十四',
  '十五',
  '十六',
  '十七',
  '十八',
  '十九',
  '二十',
  '廿一',
  '廿二',
  '廿三',
  '廿四',
  '廿五',
  '廿六',
  '廿七',
  '廿八',
  '廿九',
  '三十'
]
const lunarFestivals: Record<string, string> = {
  '1-1': '春节',
  '1-15': '元宵节',
  '2-2': '龙抬头',
  '5-5': '端午节',
  '7-7': '七夕',
  '8-15': '中秋节',
  '9-9': '重阳节',
  '12-8': '腊八节',
  '12-23': '北方小年',
  '12-24': '南方小年'
}
const solarFestivals: Record<string, string> = {
  '1-1': '元旦',
  '5-1': '劳动节',
  '5-4': '青年节',
  '6-1': '儿童节',
  '8-1': '建军节',
  '9-10': '教师节',
  '10-1': '国庆节'
}

const enabledEngines = computed(() => settings.settings.search.enabledSearchEngines)
const dateTimeLabel = computed(() => {
  const year = now.value.getFullYear()
  const month = String(now.value.getMonth() + 1).padStart(2, '0')
  const day = String(now.value.getDate()).padStart(2, '0')
  const time = now.value.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  return `${year}.${month}.${day} ${time}`
})
const lunarLabel = computed(() => formatLunarDate(now.value))
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

function formatLunarDate(date: Date) {
  const parts = new Intl.DateTimeFormat('zh-CN-u-ca-chinese', {
    month: 'long',
    day: 'numeric'
  }).formatToParts(date)
  const monthName = parts.find(part => part.type === 'month')?.value || ''
  const day = Number(parts.find(part => part.type === 'day')?.value || 0)
  const lunarMonth = lunarMonths[monthName.replace('闰', '')]
  const lunarDay = lunarDays[day] || `${day}日`
  const solarKey = `${date.getMonth() + 1}-${date.getDate()}`
  const lunarKey = `${lunarMonth}-${day}`
  const festival = solarFestivals[solarKey] || lunarFestivals[lunarKey]
  return `农历${monthName}${lunarDay}${festival ? ` · ${festival}` : ''}`
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
    <section class="home-shell">
      <header class="home-header">
        <div class="title-block">
          <h1>{{ settings.settings.appearance.siteTitle }}</h1>
          <div class="clock">
            <strong>{{ dateTimeLabel }}</strong>
            <span class="lunar">{{ lunarLabel }}</span>
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
