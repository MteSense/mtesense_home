<script setup lang="ts">
import { CheckCircle2, CircleAlert } from 'lucide-vue-next'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
</script>

<template>
  <Transition name="toast">
    <div v-if="toast.message" class="toast-message" :class="toast.tone" role="status" aria-live="polite">
      <CheckCircle2 v-if="toast.tone === 'success'" :size="18" />
      <CircleAlert v-else :size="18" />
      <span>{{ toast.message }}</span>
    </div>
  </Transition>
</template>
