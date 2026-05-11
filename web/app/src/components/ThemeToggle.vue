<script setup lang="ts">
import { Moon, Sun } from 'lucide-vue-next'
import { useUiStore } from '../stores/ui'

const ui = useUiStore()
</script>

<template>
  <button class="icon-button" type="button" :title="ui.theme === 'dark' ? 'Light' : 'Dark'" @click="ui.toggleTheme">
    <Sun v-if="ui.theme === 'dark'" :size="18" />
    <Moon v-else :size="18" />
  </button>
</template>
