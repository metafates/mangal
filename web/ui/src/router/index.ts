import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/views/AboutView.vue')
    },
    {
      path: '/search/:provider/:query?',
      name: 'search',
      component: () => import('@/views/SearchView.vue'),
    },
    {
      path: '/search/:provider/:query/:manga',
      name: 'manga',
      component: () => import('@/views/MangaView.vue')
    }
  ]
})

export default router
