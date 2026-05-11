<script setup lang="ts">
import { ExternalLink } from 'lucide-vue-next'
import type { NavLink } from '../api/types'

defineProps<{ link: NavLink }>()
</script>

<template>
  <a class="nav-card" :href="link.url" :target="link.openInNewTab ? '_blank' : '_self'" rel="noreferrer">
    <span class="nav-icon" :class="{ image: link.iconType === 'image' }">
      <img v-if="link.iconType === 'image' && link.icon" :src="link.icon" :alt="link.title" />
      <span v-else>{{ link.icon || link.title.slice(0, 1).toUpperCase() }}</span>
    </span>
    <span class="nav-copy">
      <strong>{{ link.title }}</strong>
      <small v-if="link.description">{{ link.description }}</small>
    </span>
    <ExternalLink class="nav-open" :size="16" />
  </a>
</template>
