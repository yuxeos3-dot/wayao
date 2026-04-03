import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    component: () => import('@/views/MainLayout.vue'),
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/Dashboard.vue') },
      { path: 'templates', name: 'Templates', component: () => import('@/views/Templates.vue') },
      { path: 'domains', name: 'Domains', component: () => import('@/views/Domains.vue') },
      { path: 'domains/:id/content', name: 'ContentEditor', component: () => import('@/views/ContentEditor.vue') },
      { path: 'keywords', name: 'Keywords', component: () => import('@/views/Keywords.vue') },
      { path: 'builder', name: 'Builder', component: () => import('@/views/Builder.vue') },
      { path: 'rankings', name: 'Rankings', component: () => import('@/views/Rankings.vue') },
      { path: 'settings', name: 'Settings', component: () => import('@/views/Settings.vue') },
      { path: 'health-check', name: 'HealthCheck', component: () => import('@/views/HealthCheck.vue') },
      { path: 'index-status', name: 'IndexStatus', component: () => import('@/views/IndexStatus.vue') },
      { path: 'clusters', name: 'Clusters', component: () => import('@/views/Clusters.vue') },
      { path: 'title-pool', name: 'TitlePool', component: () => import('@/views/TitlePool.vue') },
      { path: 'city-matrix', name: 'CityMatrix', component: () => import('@/views/CityMatrix.vue') },
      { path: 'security', name: 'Security', component: () => import('@/views/Security.vue') },
    ]
  }
]

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (!to.meta.public && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
