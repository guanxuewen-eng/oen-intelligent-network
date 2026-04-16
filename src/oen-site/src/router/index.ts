import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import PublicLayout from '@/components/PublicLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: PublicLayout,
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/HomePage.vue'),
        meta: { title: '首页' },
      },
      {
        path: 'assets',
        name: 'AssetCenter',
        component: () => import('@/views/AssetCenter.vue'),
        meta: { title: '资产中心' },
      },
      {
        path: 'agents',
        name: 'AgentShowcase',
        component: () => import('@/views/AgentShowcase.vue'),
        meta: { title: '智能体' },
      },
      {
        path: 'evaluation',
        name: 'Evaluation',
        component: () => import('@/views/Evaluation.vue'),
        meta: { title: '测评区' },
      },
      {
        path: 'community',
        name: 'Community',
        component: () => import('@/views/Community.vue'),
        meta: { title: '智能体交流中心' },
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
