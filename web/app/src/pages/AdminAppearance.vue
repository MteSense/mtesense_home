<script setup lang="ts">
import { onMounted, reactive, watch } from 'vue'
import { Save, Upload } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { api } from '../api/client'
import type { PublicSettings, SearchEngineId } from '../api/types'
import { searchEngines } from '../api/searchEngines'
import { defaultSettings, useSettingsStore } from '../stores/settings'
import { useToastStore } from '../stores/toast'

const settingsStore = useSettingsStore()
const toast = useToastStore()
const { t } = useI18n()
const engines = Object.keys(searchEngines) as SearchEngineId[]

function cloneSettings(settings: PublicSettings) {
  return JSON.parse(JSON.stringify(settings)) as PublicSettings
}

const form = reactive<PublicSettings>(cloneSettings(defaultSettings))

onMounted(async () => {
  await settingsStore.load()
  Object.assign(form, cloneSettings(settingsStore.settings))
})

watch(
  () => form.appearance.browserTitle,
  title => {
    document.title = title || form.appearance.siteTitle || defaultSettings.appearance.browserTitle
  }
)

async function save() {
  await settingsStore.save(cloneSettings(form))
  Object.assign(form, cloneSettings(settingsStore.settings))
  toast.show(t('saveSuccess'))
}

async function upload(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  const result = await api.upload(file)
  form.appearance.backgroundImage = result.url
}
</script>

<template>
  <section class="admin-page">
    <header class="admin-heading">
      <h1>{{ t('appearance') }}</h1>
      <p>{{ form.appearance.siteTitle }}</p>
    </header>

    <section class="tool-panel">
      <div class="form-grid">
        <label>{{ t('siteTitle') }}<input v-model="form.appearance.siteTitle" /></label>
        <label>{{ t('browserTitle') }}<input v-model="form.appearance.browserTitle" /></label>
        <label>{{ t('backgroundImage') }}<input v-model="form.appearance.backgroundImage" /></label>
        <label>{{ t('defaultTheme') }}
          <select v-model="form.appearance.defaultTheme">
            <option value="dark">Dark</option>
            <option value="light">Light</option>
          </select>
        </label>
        <label>Card opacity<input v-model.number="form.appearance.cardOpacity" type="number" min="0.1" max="1" step="0.01" /></label>
        <label>Blur<input v-model.number="form.appearance.blurStrength" type="number" min="0" max="40" /></label>
        <label>{{ t('searchEngine') }}
          <select v-model="form.search.defaultSearchEngine">
            <option v-for="engine in form.search.enabledSearchEngines" :key="engine" :value="engine">
              {{ searchEngines[engine].label }}
            </option>
          </select>
        </label>
      </div>

      <div class="engine-checks">
        <label v-for="engine in engines" :key="engine" class="check-row">
          <input v-model="form.search.enabledSearchEngines" type="checkbox" :value="engine" />
          {{ searchEngines[engine].label }}
        </label>
      </div>

      <div class="form-actions">
        <label class="upload-button">
          <Upload :size="18" />{{ t('uploadBackground') }}
          <input type="file" accept="image/*,.svg" @change="upload" />
        </label>
        <button class="primary-button" type="button" @click="save"><Save :size="18" />{{ t('save') }}</button>
      </div>
    </section>
  </section>
</template>
