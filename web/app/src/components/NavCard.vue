<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { NavLink } from '../api/types'

const props = defineProps<{ link: NavLink }>()
const iconIndex = ref(0)

function isImageIcon(icon: string) {
  const value = icon.trim()
  if (!value) return false
  if (value.startsWith('data:image/') || value.startsWith('blob:')) return true
  if (value.startsWith('/uploads/')) return true
  try {
    const url = new URL(value, window.location.origin)
    return /\.(png|jpe?g|webp|gif|svg|ico)$/i.test(url.pathname)
  } catch {
    return /\.(png|jpe?g|webp|gif|svg|ico)$/i.test(value)
  }
}

const customImageIcon = computed(() => {
  const icon = props.link.icon.trim()
  if (!icon) return ''
  return props.link.iconType === 'image' || isImageIcon(icon) ? icon : ''
})

const iconSources = computed(() => {
  const sources: string[] = []
  try {
    const url = new URL(props.link.url)
    sources.push(`${url.origin}/favicon.ico`, `${url.origin}/favicon.png`)
  } catch {
    // Links can still fall back to an admin configured image or text icon.
  }
  if (customImageIcon.value) sources.push(customImageIcon.value)
  return sources
})
const favicon = computed(() => iconSources.value[iconIndex.value] || '')
const fallbackIcon = computed(() => props.link.icon || props.link.title.slice(0, 1).toUpperCase())

watch(
  () => [props.link.icon, props.link.iconType, props.link.url],
  () => {
    iconIndex.value = 0
  }
)

function showNextIcon() {
  iconIndex.value += 1
}
</script>

<template>
  <a class="nav-card" :href="link.url" :target="link.openInNewTab ? '_blank' : '_self'" rel="noreferrer">
    <span class="nav-icon" :class="{ image: favicon }">
      <img v-if="favicon" :src="favicon" :alt="link.title" @error="showNextIcon" />
      <span v-else>{{ fallbackIcon }}</span>
    </span>
    <span class="nav-copy">
      <strong>{{ link.title }}</strong>
      <small v-if="link.description">{{ link.description }}</small>
    </span>
  </a>
</template>
