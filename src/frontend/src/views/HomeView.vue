<template>
  <div class="home-view">
    <!-- 顶部欢迎卡片 -->
    <div class="welcome-card">
      <div class="welcome-left">
        <h1 class="welcome-title">OEN 智能网络</h1>
        <p class="welcome-desc">
          中心化治理的智能优化与升级资产平台，兼容多种主流智能体平台，通过统一仪表盘管理智能体、追踪资产并监控系统健康。
        </p>
      </div>
      <div class="welcome-badge">
        <el-tag type="success" effect="plain" size="large">v{{ statusOverview?.version || '0.1.0' }}</el-tag>
      </div>
    </div>

    <!-- 统计卡片行 -->
    <el-row :gutter="20" class="stat-row">
      <el-col :span="8">
        <div class="stat-card stat-card--blue">
          <div class="stat-card-icon"><el-icon :size="28"><Cpu /></el-icon></div>
          <div class="stat-card-value">{{ totalAgents }}</div>
          <div class="stat-card-label">已注册智能体</div>
          <div class="stat-card-detail">
            <template v-for="s in agentStates" :key="s.state">
              <span class="stat-detail-item" :class="`stat-${s.state}`">{{ s.count }} {{ stateLabel(s.state) }}</span>
            </template>
          </div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-card stat-card--green">
          <div class="stat-card-icon"><el-icon :size="28"><Collection /></el-icon></div>
          <div class="stat-card-value">{{ statusOverview?.total_artifacts || 0 }}</div>
          <div class="stat-card-label">已管理资产</div>
          <div class="stat-card-detail">四类资产：知识条目、优化指令、升级方案、智能体模板</div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-card stat-card--orange">
          <div class="stat-card-icon"><el-icon :size="28"><Connection /></el-icon></div>
          <div class="stat-card-value">{{ statusOverview?.pending_recommendations || 0 }}</div>
          <div class="stat-card-label">待处理推荐</div>
          <div class="stat-card-detail">
            <span v-if="statusOverview?.pending_candidates">{{ statusOverview.pending_candidates }} 个候选待审核</span>
            <span v-else>暂无候选</span>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 推荐列表 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-title-row">
          <span class="card-title">推荐列表</span>
          <el-button size="small" @click="fetchRecommendations" :icon="Refresh">刷新</el-button>
        </div>
      </template>
      <el-table :data="recommendations" stripe style="width: 100%" v-loading="recLoading" :row-class-name="'home-table-row'">
        <el-table-column prop="title" label="推荐标题" min-width="200">
          <template #default="{ row }">
            <span class="table-link">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column label="置信度" width="160" align="center">
          <template #default="{ row }">
            <el-progress
              :percentage="Math.round(row.confidence_score * 100)"
              :stroke-width="6"
              :color="confidenceColor(row.confidence_score)"
            />
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="recStateType(row.state)" size="small" effect="plain">{{ recStateLabel(row.state) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" align="center">
          <template #default="{ row }">
            <el-button v-if="row.state === 'pending'" size="small" type="primary" @click="handleDecision(row, 'accept')">
              <el-icon><Check /></el-icon> 采纳
            </el-button>
            <el-button v-if="row.state === 'pending'" size="small" type="info" plain @click="handleDecision(row, 'later')">稍后</el-button>
            <el-button v-if="row.state === 'pending'" size="small" type="danger" plain link @click="handleDecision(row, 'ignore')">忽略</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="recommendations.length === 0 && !recLoading" description="暂无待处理的推荐" />
    </el-card>

    <!-- 系统健康 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <span class="card-title">系统健康</span>
      </template>
      <el-descriptions :column="3" border class="health-descriptions">
        <el-descriptions-item label="运行状态">
          <span class="health-running">
            <span class="health-dot"></span>
            {{ statusOverview?.status === 'running' ? '正常运行' : statusOverview?.status || '未知' }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="版本">v{{ statusOverview?.version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="智能体总数">{{ totalAgents }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Cpu, Collection, Connection, Refresh, Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getOpsStatus, getRecommendations, makeDecision } from '@/api/recommendations'
import type { Recommendation } from '@/types'

const statusOverview = ref<any>(null)
const recommendations = ref<Recommendation[]>([])
const recLoading = ref(false)

const totalAgents = computed(() => {
  if (!statusOverview.value?.agent_states) return 0
  return statusOverview.value.agent_states.reduce((sum: number, s: any) => sum + s.count, 0)
})

const agentStates = computed(() => statusOverview.value?.agent_states || [])

function stateLabel(state: string): string {
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

function confidenceColor(score: number): string {
  if (score >= 0.8) return '#16a34a'
  if (score >= 0.5) return '#d97706'
  return '#dc2626'
}

onMounted(() => {
  fetchStatus()
  fetchRecommendations()
})

async function fetchStatus() {
  try {
    const res = await getOpsStatus()
    if (res.data.code === 0) {
      statusOverview.value = res.data.data
    }
  } catch {}
}

async function fetchRecommendations() {
  recLoading.value = true
  try {
    const res = await getRecommendations({ page: 1, page_size: 10 })
    if (res.data.code === 0 && res.data.data) {
      recommendations.value = res.data.data.items || []
    }
  } catch {
    recommendations.value = []
  } finally {
    recLoading.value = false
  }
}

function recStateLabel(state: string): string {
  const map: Record<string, string> = {
    pending: '待处理',
    accepted: '已采纳',
    ignored: '已忽略',
    deferred: '已推迟',
  }
  return map[state] || state
}

function recStateType(state: string): 'success' | 'warning' | 'danger' | 'info' {
  const map: Record<string, 'success' | 'warning' | 'danger' | 'info'> = {
    pending: 'warning',
    accepted: 'success',
    ignored: 'info',
    deferred: 'info',
  }
  return map[state] || 'info'
}

async function handleDecision(rec: Recommendation, decision: 'accept' | 'ignore' | 'later') {
  try {
    await makeDecision(rec.id, { decision, decided_by: 'user' })
    const decisionLabels: Record<string, string> = { accept: '采纳', ignore: '忽略', later: '推迟' }
    ElMessage.success(`已${decisionLabels[decision]}推荐`)
    fetchRecommendations()
    fetchStatus()
  } catch (e: any) {
    ElMessage.error('操作失败: ' + (e.message || e))
  }
}
</script>

<style scoped>
.home-view {
  max-width: 1400px;
  margin: 0 auto;
}

/* 欢迎卡片 */
.welcome-card {
  background: linear-gradient(135deg, #2563eb 0%, #7c3aed 100%);
  border-radius: var(--radius-lg);
  padding: 28px 32px;
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  color: #fff;
  box-shadow: 0 8px 32px rgba(37, 99, 235, 0.2);
  position: relative;
  overflow: hidden;
}

.welcome-card::after {
  content: '';
  position: absolute;
  right: -40px;
  top: -40px;
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
}

.welcome-left {
  position: relative;
  z-index: 1;
}

.welcome-title {
  font-size: 26px;
  font-weight: 700;
  margin: 0 0 8px;
  letter-spacing: 1px;
}

.welcome-desc {
  font-size: 14px;
  line-height: 1.7;
  color: rgba(255, 255, 255, 0.85);
  max-width: 600px;
  margin: 0;
}

.welcome-badge {
  position: relative;
  z-index: 1;
  flex-shrink: 0;
}

:deep(.welcome-badge .el-tag) {
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #fff;
  font-weight: 600;
  font-size: 14px;
  padding: 6px 16px;
}

/* 统计卡片行 */
.stat-row {
  margin-bottom: 24px;
}

.stat-card {
  background: var(--bg-panel);
  border-radius: var(--radius-lg);
  padding: 24px;
  border: 1px solid var(--border);
  box-shadow: var(--shadow-card);
  position: relative;
  overflow: hidden;
  min-height: 140px;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
}

.stat-card--blue::before { background: linear-gradient(90deg, #2563eb, #60a5fa); }
.stat-card--green::before { background: linear-gradient(90deg, #16a34a, #4ade80); }
.stat-card--orange::before { background: linear-gradient(90deg, #d97706, #fbbf24); }

.stat-card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  margin-bottom: 12px;
}

.stat-card--blue .stat-card-icon { background: #eff6ff; color: #2563eb; }
.stat-card--green .stat-card-icon { background: #f0fdf4; color: #16a34a; }
.stat-card--orange .stat-card-icon { background: #fffbeb; color: #d97706; }

.stat-card-value {
  font-size: 36px;
  font-weight: 700;
  line-height: 1.2;
  margin-bottom: 4px;
}

.stat-card--blue .stat-card-value { color: #2563eb; }
.stat-card--green .stat-card-value { color: #16a34a; }
.stat-card--orange .stat-card-value { color: #d97706; }

.stat-card-label {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.stat-card-detail {
  font-size: 12px;
  color: var(--text-muted);
}

.stat-detail-item {
  margin-right: 12px;
  white-space: nowrap;
}

.stat-running { color: #16a34a; }
.stat-degraded { color: #dc2626; }
.stat-paused { color: #d97706; }

/* 区域卡片 */
.section-card {
  margin-bottom: 24px;
}

.card-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.table-link {
  color: var(--text-primary);
  cursor: default;
}

:deep(.home-table-row td) {
  padding: 12px 0;
}

/* 系统健康 */
.health-running {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--success);
  font-weight: 500;
}

.health-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--success);
  box-shadow: 0 0 6px rgba(22, 163, 74, 0.4);
}

:deep(.health-descriptions .el-descriptions__label) {
  color: var(--text-muted);
  font-weight: 500;
  width: 120px;
}
</style>
