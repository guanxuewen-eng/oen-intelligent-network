import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import AppLayout from '@/components/AppLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: AppLayout,
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/HomeView.vue'),
        meta: { title: '首页', admin: true },
      },
      {
        path: 'agents',
        name: 'AgentStatus',
        component: () => import('@/views/AgentStatusView.vue'),
        meta: { title: '智能体管理中心', admin: true },
      },
      {
        path: 'artifacts',
        name: 'ArtifactHub',
        component: () => import('@/views/ArtifactHubView.vue'),
        meta: { title: '资产中心', admin: true },
      },
      {
        path: 'ops',
        name: 'OpsConsole',
        component: () => import('@/views/OpsConsoleView.vue'),
        meta: { title: '运维控制台', admin: true },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const title = to.meta.title as string
  if (title) {
    document.title = `${title} - OEN 管理后台`
  }
  next()
})

export default router
