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

export default router
