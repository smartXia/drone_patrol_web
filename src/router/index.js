import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '@/views/Dashboard.vue'
const RedisExplorer = () => import('@/views/RedisExplorer.vue')
const ErrorCodes = () => import('@/views/ErrorCodes.vue')
const AircraftDetail = () => import('@/views/AircraftDetail.vue')
const CameraLive = () => import('@/views/CameraLive.vue')

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
  },
  {
    path: '/error-codes',
    name: 'ErrorCodes',
    component: ErrorCodes,
    meta: {
      title: 'DJI 错误码查询'
    }
  },
  {
    path: '/aircraft/:deviceSn',
    name: 'AircraftDetail',
    component: AircraftDetail,
    meta: {
      title: '飞机信息详情'
    }
  },
  {
    path: '/camera',
    name: 'CameraLive',
    component: CameraLive,
    meta: {
      title: '机场直播摄像头'
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
