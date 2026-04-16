<template>
  <div class="public-home">
    <!-- Hero 区域 -->
    <section class="hero-section">
      <div class="hero-bg">
        <div class="hero-pattern"></div>
      </div>
      <div class="hero-content">
        <div class="hero-badge">
          <span class="badge-dot"></span>
          开源智能体资产平台
        </div>
        <h1 class="hero-title">
          发现、评估、部署<br/>
          <span class="hero-highlight">AI 智能资产</span>
        </h1>
        <p class="hero-desc">
          OEN 智能网络是一个面向 AI 智能体的资产化管理平台，提供智能体注册、能力评估、资产版本管理和一键部署能力。
        </p>
        <div class="hero-actions">
          <router-link to="/assets" class="btn btn-primary">浏览资产中心</router-link>
          <a href="#featured" class="btn btn-ghost">查看精选</a>
        </div>
        <div class="hero-stats">
          <div class="hero-stat">
            <div class="stat-value">{{ stats.totalAgents }}</div>
            <div class="stat-label">注册智能体</div>
          </div>
          <div class="hero-stat-divider"></div>
          <div class="hero-stat">
            <div class="stat-value">{{ stats.totalAssets }}</div>
            <div class="stat-label">管理资产</div>
          </div>
          <div class="hero-stat-divider"></div>
          <div class="hero-stat">
            <div class="stat-value">{{ stats.pendingRecs }}</div>
            <div class="stat-label">待处理推荐</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 分类浏览 -->
    <section class="categories-section">
      <div class="section-container">
        <div class="section-header">
          <h2 class="section-title">资产分类</h2>
          <p class="section-subtitle">按类型浏览，快速定位目标资产</p>
        </div>
        <div class="category-grid">
          <div
            v-for="cat in categories"
            :key="cat.name"
            class="category-card"
            @click="goToAssets(cat.type)"
          >
            <div class="category-icon" :style="{ background: cat.bgColor }">
              <el-icon :size="28" :color="cat.color"><component :is="cat.icon" /></el-icon>
            </div>
            <h3 class="category-name">{{ cat.name }}</h3>
            <p class="category-desc">{{ cat.desc }}</p>
            <div class="category-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 精选推荐 -->
    <section id="featured" class="featured-section">
      <div class="section-container">
        <div class="section-header">
          <h2 class="section-title">精选推荐</h2>
          <p class="section-subtitle">平台精选高置信度资产，快速发现可用能力</p>
          <router-link to="/assets" class="section-more">查看全部 <el-icon><ArrowRight /></el-icon></router-link>
        </div>
        <div class="asset-grid" v-loading="loading">
          <div
            v-for="asset in featuredAssets"
            :key="asset.id"
            class="asset-card"
            @click="viewAsset(asset)"
          >
            <div class="asset-card-header">
              <div class="asset-type-tag">{{ typeLabel(asset.artifact_type) }}</div>
              <div v-if="asset.risk_level" :class="['risk-badge', riskClass(asset.risk_level)]">
                {{ riskLabel(asset.risk_level) }}
              </div>
            </div>
            <h3 class="asset-title">{{ asset.title }}</h3>
            <p class="asset-desc">{{ asset.description || '暂无描述' }}</p>
            <div class="asset-meta">
              <span class="meta-item">
                <el-icon><Clock /></el-icon>
                {{ formatDate(asset.created_at) }}
              </span>
              <span v-if="asset.verification_status" class="meta-item meta-status">
                <el-icon><CircleCheck /></el-icon>
                {{ verifyLabel(asset.verification_status) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 智能体能力速览 -->
    <section class="agents-preview-section">
      <div class="section-container">
        <div class="section-header">
          <h2 class="section-title">智能体动态</h2>
          <p class="section-subtitle">了解平台中智能体的运行状态和最新变化</p>
          <router-link to="/agents" class="section-more">进入智能体频道 <el-icon><ArrowRight /></el-icon></router-link>
        </div>
        <div class="agent-cards-row">
          <div
            v-for="agent in recentAgents"
            :key="agent.id"
            class="agent-mini-card"
          >
            <div :class="['agent-state-dot', stateClass(agent.state)]"></div>
            <div class="agent-info">
              <div class="agent-name">{{ agent.name }}</div>
              <div class="agent-meta-text">{{ agent.agent_type }} · {{ stateLabel(agent.state) }}</div>
            </div>
          </div>
          <div v-if="recentAgents.length === 0" class="empty-hint">
            暂无已注册智能体
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Document, Setting, Connection, Files, ArrowRight,
  Clock, CircleCheck,
} from '@element-plus/icons-vue'
import { getOpsStatus, getRecommendations } from '@/api/recommendations'
import { getArtifacts } from '@/api/artifacts'
import { getAgents } from '@/api/agents'
import type { Artifact, Agent } from '@/types'

const router = useRouter()

const stats = ref({ totalAgents: 0, totalAssets: 0, pendingRecs: 0 })
const featuredAssets = ref<Artifact[]>([])
const recentAgents = ref<Agent[]>([])
const loading = ref(false)

const categories = [
  { name: '知识条目', type: 'knowledge_entry', desc: '结构化的领域知识与经验沉淀', icon: Document, bgColor: '#eff6ff', color: '#2563eb' },
  { name: '优化指令', type: 'optimization_instruction', desc: '针对特定场景的优化策略与指令集', icon: Setting, bgColor: '#f0fdf4', color: '#16a34a' },
  { name: '升级方案', type: 'upgrade_plan', desc: '系统升级、迁移、适配方案文档', icon: Connection, bgColor: '#fef3c7', color: '#d97706' },
  { name: '智能体模板', type: 'agent_template', desc: '可复用的智能体配置模板', icon: Files, bgColor: '#faf5ff', color: '#7c3aed' },
]

const typeMap: Record<string, string> = {
  knowledge_entry: '知识条目',
  optimization_instruction: '优化指令',
  upgrade_plan: '升级方案',
  agent_template: '智能体模板',
}

function typeLabel(type: string) {
  return typeMap[type] || type
}

function riskClass(level: string) {
  const map: Record<string, string> = { low: 'risk-low', medium: 'risk-medium', high: 'risk-high', critical: 'risk-critical' }
  return map[level] || ''
}

function riskLabel(level: string) {
  const map: Record<string, string> = { low: '低风险', medium: '中风险', high: '高风险', critical: '严重' }
  return map[level] || level
}

function verifyLabel(status: string) {
  const map: Record<string, string> = { verified: '已验证', pending: '待验证', failed: '未通过' }
  return map[status] || status
}

function stateClass(state: string) {
  const map: Record<string, string> = {
    running: 'dot-running',
    degraded: 'dot-degraded',
    paused: 'dot-paused',
    authorized: 'dot-authorized',
  }
  return map[state] || 'dot-default'
}

function stateLabel(state: string) {
  const map: Record<string, string> = {
    unauthorized: '未授权',
    authorized: '已授权',
    creating: '创建中',
    running: '运行中',
    degraded: '已降级',
    paused: '已暂停',
    revoked: '已撤销',
  }
  return map[state] || state
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  return dateStr.slice(0, 10)
}

function goToAssets(type?: string) {
  router.push(type ? `/assets?type=${type}` : '/assets')
}

function viewAsset(asset: Artifact) {
  router.push(`/assets/${asset.id}`)
}

onMounted(async () => {
  // 获取系统状态
  try {
    const res = await getOpsStatus()
    if (res.data.code === 0 && res.data.data) {
      const d = res.data.data
      stats.value.totalAgents = d.agent_states?.reduce((s: number, x: any) => s + x.count, 0) || 0
      stats.value.totalAssets = d.total_artifacts || 0
      stats.value.pendingRecs = d.pending_recommendations || 0
    }
  } catch {}

  // 获取资产列表
  loading.value = true
  try {
    const res = await getArtifacts({ page: 1, page_size: 8 })
    if (res.data.code === 0 && res.data.data) {
      featuredAssets.value = res.data.data.items || []
    }
  } catch {} finally {
    loading.value = false
  }

  // 获取智能体列表
  try {
    const res = await getAgents({ page: 1, page_size: 6 })
    if (res.data.code === 0 && res.data.data) {
      recentAgents.value = res.data.data.items || []
    }
  } catch {}
})
</script>

<style scoped>
.public-home {
  background: #fafbfc;
}

.section-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Hero */
.hero-section {
  position: relative;
  padding: 80px 24px 64px;
  text-align: center;
  overflow: hidden;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 40%, #312e81 100%);
  color: #fff;
}

.hero-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.hero-pattern {
  position: absolute;
  top: -50%;
  left: -20%;
  width: 140%;
  height: 200%;
  background:
    radial-gradient(circle at 20% 80%, rgba(37, 99, 235, 0.15) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(124, 58, 237, 0.12) 0%, transparent 50%),
    radial-gradient(circle at 50% 50%, rgba(6, 182, 212, 0.08) 0%, transparent 60%);
}

.hero-content {
  position: relative;
  z-index: 1;
  max-width: 720px;
  margin: 0 auto;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.15);
  font-size: 13px;
  color: #e2e8f0;
  margin-bottom: 24px;
}

.badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #4ade80;
  box-shadow: 0 0 6px rgba(74, 222, 128, 0.5);
}

.hero-title {
  font-size: 48px;
  font-weight: 800;
  line-height: 1.2;
  margin-bottom: 20px;
  letter-spacing: -0.5px;
}

.hero-highlight {
  background: linear-gradient(90deg, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-desc {
  font-size: 16px;
  line-height: 1.8;
  color: #94a3b8;
  margin-bottom: 32px;
  max-width: 560px;
  margin-left: auto;
  margin-right: auto;
}

.hero-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 48px;
}

.btn {
  padding: 12px 28px;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s ease;
  cursor: pointer;
}

.btn-primary {
  background: #2563eb;
  color: #fff;
  border: none;
  box-shadow: 0 4px 16px rgba(37, 99, 235, 0.3);
}

.btn-primary:hover {
  background: #1d4ed8;
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.4);
}

.btn-ghost {
  background: transparent;
  color: #cbd5e1;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-ghost:hover {
  border-color: rgba(255, 255, 255, 0.4);
  color: #fff;
}

.hero-stats {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 40px;
  padding-top: 32px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.hero-stat {
  text-align: center;
}

.stat-value {
  font-size: 32px;
  font-weight: 800;
  color: #f1f5f9;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #64748b;
  margin-top: 4px;
}

.hero-stat-divider {
  width: 1px;
  height: 40px;
  background: rgba(255, 255, 255, 0.1);
}

/* Section headers */
.section-header {
  text-align: center;
  margin-bottom: 40px;
  position: relative;
}

.section-title {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 8px;
}

.section-subtitle {
  font-size: 14px;
  color: #94a3b8;
}

.section-more {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  font-size: 14px;
  color: #2563eb;
  text-decoration: none;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 4px;
}

.section-more:hover {
  color: #1d4ed8;
}

/* Categories */
.categories-section {
  padding: 64px 0;
  background: #fff;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.category-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 28px 24px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;
}

.category-card:hover {
  border-color: #2563eb;
  box-shadow: 0 8px 24px rgba(37, 99, 235, 0.08);
  transform: translateY(-2px);
}

.category-icon {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  margin-bottom: 16px;
}

.category-name {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 8px;
}

.category-desc {
  font-size: 13px;
  color: #94a3b8;
  line-height: 1.6;
}

.category-arrow {
  position: absolute;
  right: 20px;
  bottom: 20px;
  color: #cbd5e1;
  transition: all 0.2s ease;
}

.category-card:hover .category-arrow {
  color: #2563eb;
  transform: translateX(4px);
}

/* Featured assets */
.featured-section {
  padding: 64px 0;
}

.asset-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.asset-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.asset-card:hover {
  border-color: #2563eb;
  box-shadow: 0 8px 24px rgba(37, 99, 235, 0.08);
  transform: translateY(-2px);
}

.asset-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.asset-type-tag {
  font-size: 12px;
  padding: 3px 10px;
  border-radius: 4px;
  background: #f1f5f9;
  color: #64748b;
  font-weight: 500;
}

.risk-badge {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.risk-low { background: #f0fdf4; color: #16a34a; }
.risk-medium { background: #fef3c7; color: #d97706; }
.risk-high { background: #fee2e2; color: #dc2626; }
.risk-critical { background: #fce7f3; color: #be185d; }

.asset-title {
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 8px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.asset-desc {
  font-size: 13px;
  color: #94a3b8;
  line-height: 1.6;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.asset-meta {
  display: flex;
  gap: 16px;
  align-items: center;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #94a3b8;
}

.meta-status {
  color: #16a34a;
}

/* Agents preview */
.agents-preview-section {
  padding: 64px 0;
  background: #fff;
}

.agent-cards-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.agent-mini-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: #fafbfc;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  transition: all 0.15s ease;
}

.agent-mini-card:hover {
  border-color: #2563eb;
  background: #eff6ff;
}

.agent-state-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.dot-running { background: #16a34a; box-shadow: 0 0 6px rgba(22, 163, 74, 0.4); }
.dot-degraded { background: #dc2626; }
.dot-paused { background: #d97706; }
.dot-authorized { background: #2563eb; }
.dot-default { background: #94a3b8; }

.agent-info {
  flex: 1;
  min-width: 0;
}

.agent-name {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.agent-meta-text {
  font-size: 12px;
  color: #94a3b8;
}

.empty-hint {
  grid-column: 1 / -1;
  text-align: center;
  padding: 40px;
  color: #94a3b8;
  font-size: 14px;
}

/* Responsive */
@media (max-width: 1024px) {
  .category-grid,
  .asset-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .agent-cards-row {
    grid-template-columns: repeat(2, 1fr);
  }
  .hero-title {
    font-size: 36px;
  }
}

@media (max-width: 640px) {
  .category-grid,
  .asset-grid {
    grid-template-columns: 1fr;
  }
  .agent-cards-row {
    grid-template-columns: 1fr;
  }
  .header-inner {
    gap: 16px;
  }
  .nav-menu {
    display: none;
  }
  .hero-title {
    font-size: 28px;
  }
  .hero-stats {
    gap: 24px;
  }
  .stat-value {
    font-size: 24px;
  }
  .footer-inner {
    flex-direction: column;
    gap: 32px;
  }
  .footer-links {
    gap: 32px;
  }
}
</style>
