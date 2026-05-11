<script setup lang="ts">
import { Search } from 'lucide-vue-next'
import type { SearchEngineId } from '../api/types'
import { searchEngines } from '../api/searchEngines'

defineProps<{
  query: string
  engine: SearchEngineId
  enabledEngines: SearchEngineId[]
  placeholder: string
}>()

const emit = defineEmits<{
  'update:query': [value: string]
  'update:engine': [value: SearchEngineId]
  submit: []
}>()
</script>

<template>
  <form class="search-bar" @submit.prevent="emit('submit')">
    <select :value="engine" class="search-engine" @change="emit('update:engine', ($event.target as HTMLSelectElement).value as SearchEngineId)">
      <option v-for="item in enabledEngines" :key="item" :value="item">
        {{ searchEngines[item].label }}
      </option>
    </select>
    <input :value="query" :placeholder="placeholder" @input="emit('update:query', ($event.target as HTMLInputElement).value)" />
    <button class="search-submit" type="submit" :title="placeholder">
      <Search :size="20" />
    </button>
  </form>
</template>
