<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { Home, LogIn } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import LanguageToggle from '../components/LanguageToggle.vue'
import ThemeToggle from '../components/ThemeToggle.vue'

const router = useRouter()
const auth = useAuthStore()
const settings = useSettingsStore()
const { t } = useI18n()
const form = reactive({ username: '', password: '' })

const backgroundStyle = computed(() => {
  const image = settings.settings.appearance.backgroundImage
  return image ? { backgroundImage: `url(${image})` } : {}
})

async function submit() {
  await auth.login(form.username, form.password)
  router.push('/admin/links')
}

onMounted(() => settings.load())
</script>

<template>
  <main class="login-page" :style="backgroundStyle">
    <div class="login-actions">
      <RouterLink class="icon-button" to="/" :title="settings.settings.appearance.siteTitle">
        <Home :size="18" />
      </RouterLink>
      <LanguageToggle />
      <ThemeToggle />
    </div>
    <form class="login-panel" @submit.prevent="submit">
      <h1>{{ t('adminLogin') }}</h1>
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
