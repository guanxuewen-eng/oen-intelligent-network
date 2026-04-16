<template>
  <div class="ops-console-view">
    <el-row :gutter="20">
      <!-- 左侧：系统状态 + API -->
      <el-col :span="16">
        <!-- 系统状态 -->
        <el-card shadow="never" class="section-card">
          <template #header>
            <span class="card-title">系统状态</span>
          </template>
          <div class="status-grid">
            <div class="status-item">
              <span class="status-label">运行状态</span>
              <span class="status-value">
                <span class="status-indicator" :class="opsStatus?.status === 'running' ? 'status-indicator--running' : 'status-indicator--warning'"></span>
                {{ opsStatus?.status === 'running' ? '正常运行' : opsStatus?.status || '加载中' }}
              </span>
            </div>
            <div class="status-item">
              <span class="status-label">版本</span>
              <span class="status-value mono">v{{ opsStatus?.version || '0.1.0' }}</span>
            </div>
            <div class="status-item">
              <span class="status-label">资产总数</span>
              <span class="status-value highlight-blue">{{ opsStatus?.total_artifacts || 0 }}</span>
            </div>
            <div class="status-item">
              <span class="status-label">待处理推荐</span>
              <span class="status-value highlight-orange">{{ opsStatus?.pending_recommendations || 0 }}</span>
            </div>
          </div>
          <el-divider />
          <div class="agent-summary" v-if="agentStateSummary">
            <div
              v-for="s in agentStateSummary"
              :key="s.state"
              class="agent-summary-item"
              :class="`agent-summary--${s.state}`"
            >
              <span class="summary-dot"></span>
              <span class="summary-label">{{ stateLabel(s.state) }}</span>
              <span class="summary-count">{{ s.count }}</span>
            </div>
            <div v-if="agentStateSummary.length === 0" class="no-agents">暂无智能体</div>
          </div>
        </el-card>

        <!-- API 接口 -->
        <el-card shadow="never" class="section-card">
          <template #header>
            <span class="card-title">API 接口</span>
          </template>
          <el-table :data="endpoints" stripe size="small" class="api-table">
            <el-table-column prop="method" label="方法" width="90" align="center">
              <template #default="{ row }">
                <el-tag :type="getMethodType(row.method)" size="small" effect="dark" class="method-tag">{{ row.method }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="path" label="路径">
              <template #default="{ row }">
                <span class="mono path-text">{{ row.path }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="说明" width="140">
              <template #default="{ row }">
                <span class="desc-text">{{ row.description }}</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 右侧：快捷操作 + 最近活动 -->
      <el-col :span="8">
        <el-card shadow="never" class="section-card">
          <template #header>
            <span class="card-title">快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" style="width: 100%" @click="refreshAll" :icon="Refresh">
              全部刷新
            </el-button>
            <el-button type="default" style="width: 100%" @click="loadAuditLogs">
              加载审计日志
            </el-button>
          </div>
        </el-card>

        <el-card shadow="never" class="section-card activity-card">
          <template #header>
            <span class="card-title">最近活动</span>
          </template>
          <el-timeline v-if="recentLogs.length > 0" class="activity-timeline">
            <el-timeline-item
              v-for="log in recentLogs"
              :key="log.id"
              :timestamp="formatDate(log.created_at)"
              placement="top"
              :color="logColor(log.event_type)"
              class="timeline-item"
            >
              <div class="timeline-content">
                <span class="timeline-title">{{ log.description }}</span>
                <el-tag size="small" :type="logTagType(log.event_type)" effect="plain" class="timeline-type">{{ log.event_type }}</el-tag>
              </div>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else description="暂无活动" :image-size="80" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 审计日志 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">审计日志</span>
          <el-button size="small" :icon="Refresh" @click="loadAuditLogs">刷新</el-button>
        </div>
      </template>
      <el-table :data="auditLogs" stripe style="width: 100%" v-loading="auditLoading" class="audit-table">
        <el-table-column prop="event_type" label="事件" width="130">
          <template #default="{ row }">
            <el-tag size="small" :type="logTagType(row.event_type)" effect="plain">{{ row.event_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="event_subtype" label="子类型" width="120" />
        <el-table-column prop="actor_type" label="操作者" width="100" />
        <el-table-column prop="target_type" label="目标" width="100" />
        <el-table-column prop="target_id" label="ID" width="80">
          <template #default="{ row }">
            <span class="mono">{{ row.target_id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="250" show-overflow-tooltip />
        <el-table-column label="时间" width="170">
          <template #default="{ row }">
            <span class="time-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!auditLoading && auditLogs.length === 0" description="暂无审计日志" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { AuditLog } from '@/types'
import { getOpsStatus, getAuditLogs } from '@/api/recommendations'

const opsStatus = ref<any>(null)
const auditLogs = ref<AuditLog[]>([])
const auditLoading = ref(false)
const recentLogs = ref<AuditLog[]>([])
const agentStateSummary = ref<any[]>([])

const endpoints = ref([
  { method: 'GET', path: '/api/v1/agents', description: '列出智能体' },
  { method: 'POST', path: '/api/v1/agents/:id/consent', description: '授权智能体' },
  { method: 'POST', path: '/api/v1/agents/:id/heartbeat', description: '上报心跳' },
  { method: 'POST', path: '/api/v1/agents/:id/rebuild', description: '重建智能体' },
  { method: 'GET', path: '/api/v1/artifacts', description: '列出资产' },
  { method: 'GET', path: '/api/v1/artifacts/:id/versions', description: '列出版本' },
  { method: 'GET', path: '/api/v1/artifacts/:id/versions/:ver/view/:type', description: '获取资产视图' },
  { method: 'GET', path: '/api/v1/recommendations', description: '列出推荐' },
  { method: 'POST', path: '/api/v1/recommendations/:id/decision', description: '处理决策' },
  { method: 'GET', path: '/api/v1/candidates', description: '列出候选' },
  { method: 'GET', path: '/api/v1/ops/status', description: '系统状态' },
  { method: 'GET', path: '/api/v1/ops/audit', description: '审计日志' },
])

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

function getMethodType(method: string): 'success' | 'warning' | 'danger' {
  const map: Record<string, 'success' | 'warning' | 'danger'> = { GET: 'success', POST: 'warning', PUT: 'danger', DELETE: 'danger' }
  return map[method] || 'info'
}

function logTagType(event: string): 'success' | 'warning' | 'danger' | 'info' {
  const map: Record<string, 'success' | 'warning' | 'danger' | 'info'> = {
    consent: 'success', agent: 'info', artifact: 'success', recommendation: 'warning', candidate: 'info',
  }
  return map[event] || 'info'
}

function logColor(event: string): string {
  const map: Record<string, string> = {
    consent: '#16a34a', agent: '#7b8799', artifact: '#16a34a', recommendation: '#d97706', candidate: '#7b8799',
  }
  return map[event] || '#7b8799'
}

function formatDate(date: string): string {
  return new Date(date).toLocaleString()
}

onMounted(() => {
  refreshAll()
})

async function refreshAll() {
  try {
    const res = await getOpsStatus()
    if (res.data.code === 0) {
      opsStatus.value = res.data.data
      agentStateSummary.value = res.data.data?.agent_states || []
    }
  } catch {}
  loadAuditLogs()
}

async function loadAuditLogs() {
  auditLoading.value = true
  try {
    const res = await getAuditLogs({ page: 1, page_size: 50 })
    if (res.data.code === 0 && res.data.data) {
      auditLogs.value = res.data.data.items || []
      recentLogs.value = (res.data.data.items || []).slice(0, 5)
    }
  } catch (e: any) {
    ElMessage.error('加载审计日志失败: ' + (e.message || e))
  } finally {
    auditLoading.value = false
  }
}
</script>

<style scoped>
.ops-console-view {
  max-width: 1400px;
  margin: 0 auto;
}

.section-card {
  margin-bottom: 20px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* 系统状态网格 */
.status-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px 24px;
}

.status-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--bg-soft);
  border-radius: var(--radius-sm);
}

.status-label {
  font-size: 13px;
  color: var(--text-muted);
}

.status-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.mono {
  font-family: 'SF Mono', 'Fira Code', monospace;
}

.highlight-blue { color: #2563eb; }
.highlight-orange { color: #d97706; }

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator--running {
  background: #16a34a;
  box-shadow: 0 0 6px rgba(22, 163, 74, 0.4);
}

.status-indicator--warning {
  background: #d97706;
}

/* 智能体状态汇总 */
.agent-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.agent-summary-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  background: var(--bg-soft);
  color: var(--text-secondary);
}

.summary-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.agent-summary--running .summary-dot { background: #16a34a; }
.agent-summary--running { background: #f0fdf4; color: #16a34a; }

.agent-summary--degraded .summary-dot, .agent-summary--revoked .summary-dot { background: #dc2626; }
.agent-summary--degraded, .agent-summary--revoked { background: #fef2f2; color: #dc2626; }

.agent-summary--paused .summary-dot, .agent-summary--creating .summary-dot { background: #d97706; }
.agent-summary--paused, .agent-summary--creating { background: #fffbeb; color: #d97706; }

.agent-summary--unauthorized .summary-dot { background: #7b8799; }
.agent-summary--unauthorized { background: #f8fafc; color: #7b8799; }

.agent-summary--authorized .summary-dot { background: #16a34a; }
.agent-summary--authorized { background: #f0fdf4; color: #16a34a; }

.summary-count {
  font-weight: 700;
}

.no-agents {
  font-size: 13px;
  color: var(--text-muted);
}

/* API 表格 */
.api-table .path-text {
  font-size: 12px;
  color: var(--text-secondary);
}

.api-table .desc-text {
  font-size: 13px;
  color: var(--text-secondary);
}

.method-tag {
  min-width: 52px;
  text-align: center;
}

/* 快捷操作 */
.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.quick-actions .el-button {
  height: 40px;
}

/* 活动时间线 */
.activity-card {
  min-height: 200px;
}

.activity-timeline {
  margin-top: 0;
}

.timeline-item {
  padding-bottom: 12px;
}

.timeline-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.timeline-title {
  font-size: 13px;
  color: var(--text-primary);
  line-height: 1.4;
  flex: 1;
}

.timeline-type {
  flex-shrink: 0;
}

:deep(.timeline-item .el-timeline-item__timestamp) {
  font-size: 11px;
  color: var(--text-muted);
}

/* 审计日志 */
.audit-table .time-text {
  font-size: 12px;
  color: var(--text-muted);
}
</style>
