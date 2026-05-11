<script setup lang="ts">
import { Languages } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useUiStore } from '../stores/ui'

const { locale } = useI18n()
const ui = useUiStore()

function toggleLocale() {
  const next = ui.locale === 'zh' ? 'en' : 'zh'
  ui.setLocale(next)
  locale.value = next
}
</script>

<template>
  <button class="icon-button" type="button" title="Language" @click="toggleLocale">
    <Languages :size="18" />
    <span>{{ ui.locale === 'zh' ? 'EN' : '中' }}</span>
  </button>
</template>
