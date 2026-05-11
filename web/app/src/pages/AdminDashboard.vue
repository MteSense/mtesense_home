<script setup lang="ts">
import AdminLayout from '../layouts/AdminLayout.vue'
import { useUiStore } from '../stores/ui'

const ui = useUiStore()
ui.setTheme(ui.theme)
</script>

<template>
  <AdminLayout />
</template>
