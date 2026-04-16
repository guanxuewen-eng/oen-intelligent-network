import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import PublicLayout from '@/components/PublicLayout.vue'
import AppLayout from '@/components/AppLayout.vue'

const routes: RouteRecordRaw[] = [
  // === 面向公众的网站 ===
  {
    path: '/',
    component: PublicLayout,
    children: [
      {
        path: '',
        name: 'PublicHome',
        component: () => import('@/views/public/HomePage.vue'),
        meta: { title: '首页', public: true },
      },
      {
        path: 'assets',
        name: 'AssetCenter',
        component: () => import('@/views/public/AssetCenter.vue'),
        meta: { title: '资产中心', public: true },
      },
      {
        path: 'agents',
        name: 'AgentShowcase',
        component: () => import('@/views/public/AgentShowcase.vue'),
        meta: { title: '智能体', public: true },
      },
      {
        path: 'evaluation',
        name: 'Evaluation',
        component: () => import('@/views/public/Evaluation.vue'),
        meta: { title: '测评区', public: true },
      },
      {
        path: 'community',
        name: 'Community',
        component: () => import('@/views/public/Community.vue'),
        meta: { title: '智能体交流中心', public: true },
      },
    ],
  },

  // === 管理后台 ===
  {
    path: '/admin',
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
    document.title = `${title} - OEN 智能网络`
  }
  next()
})

export default router
