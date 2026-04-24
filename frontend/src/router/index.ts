import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/views/HomePage.vue'
import HistoryPage from '@/views/HistoryPage.vue'
import SettingsPage from '@/views/SettingsPage.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomePage },
    { path: '/history', component: HistoryPage },
    { path: '/settings', component: SettingsPage }
  ]
})

export default router
