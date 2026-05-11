<script setup lang="ts">
import { Home, LayoutGrid, LogOut, Palette } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import ThemeToggle from '../components/ThemeToggle.vue'
import LanguageToggle from '../components/LanguageToggle.vue'

const router = useRouter()
const auth = useAuthStore()
const { t } = useI18n()

function logout() {
  auth.logout()
  router.push('/')
}
</script>

<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <RouterLink class="brand-mini" to="/">MteSense</RouterLink>
      <nav>
        <RouterLink to="/admin/links"><LayoutGrid :size="18" />{{ t('links') }}</RouterLink>
        <RouterLink to="/admin/appearance"><Palette :size="18" />{{ t('appearance') }}</RouterLink>
      </nav>
      <div class="admin-sidebar-actions">
        <RouterLink class="icon-button" to="/" title="Home"><Home :size="18" /></RouterLink>
        <LanguageToggle />
        <ThemeToggle />
        <button class="icon-button" type="button" :title="t('logout')" @click="logout"><LogOut :size="18" /></button>
      </div>
    </aside>
    <main class="admin-content">
      <RouterView />
    </main>
  </div>
</template>
