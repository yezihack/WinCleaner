import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('@/views/Dashboard.vue'),
    },
    {
      path: '/cleaner',
      name: 'cleaner',
      component: () => import('@/views/Cleaner.vue'),
    },
    {
      path: '/memory',
      name: 'memory',
      component: () => import('@/views/Memory.vue'),
    },
    {
      path: '/process',
      name: 'process',
      component: () => import('@/views/Process.vue'),
    },
    {
      path: '/network',
      name: 'network',
      component: () => import('@/views/Network.vue'),
    },
    {
      path: '/disk',
      name: 'disk',
      component: () => import('@/views/Disk.vue'),
    },
  ],
})

export default router
