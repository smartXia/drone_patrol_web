import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '@/views/Dashboard.vue'
const RedisExplorer = () => import('@/views/RedisExplorer.vue')

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: {
      title: '设备监控面板'
    }
  },
  {
    path: '/redis',
    name: 'RedisExplorer',
    component: RedisExplorer,
    meta: {
      title: 'Redis 可视化'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title || '无人机设备监控系统'
  next()
})

export default router
