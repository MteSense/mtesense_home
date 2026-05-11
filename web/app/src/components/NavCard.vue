<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { NavLink } from '../api/types'

const props = defineProps<{ link: NavLink }>()
const iconIndex = ref(0)

const iconSources = computed(() => {
  if (props.link.iconType === 'image' && props.link.icon) return [props.link.icon]
  try {
    const url = new URL(props.link.url)
    return [
      `https://www.google.com/s2/favicons?domain_url=${encodeURIComponent(url.origin)}&sz=64`,
      `https://icons.duckduckgo.com/ip3/${url.hostname}.ico`,
      `${url.origin}/favicon.ico`
    ]
  } catch {
    return []
  }
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
