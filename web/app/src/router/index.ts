import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '../api/client'
import HomePage from '../pages/HomePage.vue'
import AdminLogin from '../pages/AdminLogin.vue'
import AdminDashboard from '../pages/AdminDashboard.vue'
import AdminLinks from '../pages/AdminLinks.vue'
import AdminAppearance from '../pages/AdminAppearance.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomePage },
    { path: '/admin/login', component: AdminLogin },
    {
      path: '/admin',
      component: AdminDashboard,
      meta: { auth: true },
      children: [
        { path: '', redirect: '/admin/links' },
        { path: 'links', component: AdminLinks },
        { path: 'appearance', component: AdminAppearance }
      ]
    }
  ]
})

router.beforeEach(to => {
  if (to.path === '/admin/login' && getToken()) {
    return '/admin/links'
  }
  if (to.meta.auth && !getToken()) {
    return '/admin/login'
  }
})

router.afterEach(to => {
  syncRobotsMeta(to.path.startsWith('/admin'))
  syncCanonicalLink(!to.path.startsWith('/admin'))
})

function syncRobotsMeta(noindex: boolean) {
  let tag = document.querySelector<HTMLMetaElement>('meta[name="robots"]')
  if (!noindex) {
    tag?.remove()
    return
  }
  if (!tag) {
    tag = document.createElement('meta')
    tag.name = 'robots'
    document.head.appendChild(tag)
  }
  tag.content = 'noindex,nofollow'
}

function syncCanonicalLink(enabled: boolean) {
  let tag = document.querySelector<HTMLLinkElement>('link[rel="canonical"]')
  if (!enabled) {
    tag?.remove()
    return
  }
  if (!tag) {
    tag = document.createElement('link')
    tag.rel = 'canonical'
    document.head.appendChild(tag)
  }
  tag.href = `${window.location.origin}/`
}

export default router
