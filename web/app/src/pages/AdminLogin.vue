<script setup lang="ts">
import { reactive } from 'vue'
import { LogIn } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import LanguageToggle from '../components/LanguageToggle.vue'
import ThemeToggle from '../components/ThemeToggle.vue'

const router = useRouter()
const auth = useAuthStore()
const { t } = useI18n()
const form = reactive({ username: 'admin', password: '' })

async function submit() {
  await auth.login(form.username, form.password)
  router.push('/admin/links')
}
</script>

<template>
  <main class="login-page">
    <div class="login-actions">
      <LanguageToggle />
      <ThemeToggle />
    </div>
    <form class="login-panel" @submit.prevent="submit">
      <h1>{{ t('login') }}</h1>
      <label>
        {{ t('username') }}
        <input v-model="form.username" autocomplete="username" />
      </label>
      <label>
        {{ t('password') }}
        <input v-model="form.password" type="password" autocomplete="current-password" />
      </label>
      <p v-if="auth.error" class="error-text">{{ auth.error }}</p>
      <button class="primary-button" type="submit" :disabled="auth.loading">
        <LogIn :size="18" />{{ t('login') }}
      </button>
    </form>
  </main>
</template>
