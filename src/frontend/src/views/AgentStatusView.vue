<template>
  <div class="agent-status-view">
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">智能体管理中心</span>
          <div class="header-actions">
            <el-select v-model="filterState" placeholder="按状态筛选" clearable class="filter-select" @change="applyFilters">
              <el-option label="全部状态" value="" />
              <el-option label="未授权" value="unauthorized" />
              <el-option label="已授权" value="authorized" />
              <el-option label="运行中" value="running" />
              <el-option label="创建中" value="creating" />
              <el-option label="已降级" value="degraded" />
              <el-option label="已暂停" value="paused" />
              <el-option label="已撤销" value="revoked" />
            </el-select>
            <el-button type="primary" :icon="Refresh" @click="fetchAgents">刷新</el-button>
          </div>
        </div>
      </template>

      <!-- 智能体类型筛选 -->
      <div class="type-filter-row">
        <span class="type-filter-label">平台类型</span>
        <div class="type-pills">
          <span
            v-for="t in typeOptions"
            :key="t.value"
            class="type-pill"
            :class="{ 'type-pill--active': filterType === t.value }"
            @click="toggleType(t.value)"
          >
            {{ t.label }}
            <span class="type-count">{{ typeCounts[t.value] ?? 0 }}</span>
          </span>
        </div>
      </div>

      <el-table :data="agents" stripe style="width: 100%" v-loading="loading" class="agent-table">
        <el-table-column prop="agent_key" label="智能体标识" min-width="160">
          <template #default="{ row }">
            <span class="mono-key">{{ row.agent_key }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="140" />
        <el-table-column label="平台" width="130">
          <template #default="{ row }">
            <el-tag size="small" :type="agentTypeTag(row.agent_type)" effect="plain">{{ typeLabel(row.agent_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="role" label="角色" width="140">
          <template #default="{ row }">
            <span class="role-text">{{ row.role || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110" align="center">
          <template #default="{ row }">
            <span class="state-badge" :class="`state--${row.state}`">
              <span class="state-dot"></span>
              {{ stateLabel(row.state) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="route_mode" label="路由模式" width="110" />
        <el-table-column label="最后心跳" width="160">
          <template #default="{ row }">
            <span :class="{ 'heartbeat-old': !isRecent(row.last_heartbeat) }">{{ formatTime(row.last_heartbeat) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <div class="action-group">
              <el-button size="small" type="primary" link @click="viewHeartbeats(row)">心跳</el-button>
              <el-button v-if="row.state === 'unauthorized'" size="small" type="success" link @click="handleConsent(row)">授权</el-button>
              <el-button v-if="row.state === 'degraded' || row.state === 'paused'" size="small" type="warning" link @click="handleRebuild(row)">重建</el-button>
              <el-button v-if="row.state === 'running'" size="small" type="info" link @click="handlePause(row)">暂停</el-button>
              <el-button v-if="row.state === 'paused'" size="small" type="success" link @click="handleResume(row)">恢复</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > pageSize"
        style="margin-top: 20px; justify-content: flex-end"
        background
        layout="total, prev, pager, next"
        :total="total"
        :page-size="pageSize"
        :current-page="currentPage"
        @current-change="handlePageChange"
      />
    </el-card>

    <!-- 心跳对话框 -->
    <el-dialog v-model="hbDialogVisible" :title="`心跳记录 — ${selectedAgent?.name}`" width="720px" class="hb-dialog">
      <el-timeline v-if="heartbeats.length > 0" class="hb-timeline">
        <el-timeline-item
          v-for="hb in heartbeats"
          :key="hb.id"
          :timestamp="formatTime(hb.heartbeat_at)"
          placement="top"
          :color="hb.error_message ? '#dc2626' : '#16a34a'"
        >
          <div class="hb-item">
            <el-tag :type="hb.error_message ? 'danger' : 'success'" size="small" effect="plain">{{ hb.status === 'running' ? '正常' : hb.status || '未知' }}</el-tag>
            <span v-if="hb.route_mode" class="hb-meta"> 路由: {{ hb.route_mode }}</span>
            <span v-if="hb.cpu_usage > 0" class="hb-meta"> CPU: {{ hb.cpu_usage.toFixed(1) }}%</span>
            <span v-if="hb.memory_usage > 0" class="hb-meta"> 内存: {{ hb.memory_usage.toFixed(1) }}%</span>
            <div v-if="hb.error_message" class="hb-error">{{ hb.error_message }}</div>
          </div>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else description="暂无心跳记录" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Agent, AgentHeartbeat } from '@/types'
import { getAgents, consentAgent, rebuildAgent, pauseAgent, resumeAgent, getHeartbeats } from '@/api/agents'

const loading = ref(false)
const filterState = ref('')
const filterType = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const agents = ref<Agent[]>([])

const hbDialogVisible = ref(false)
const selectedAgent = ref<Agent | null>(null)
const heartbeats = ref<AgentHeartbeat[]>([])

const typeCounts = ref<Record<string, number>>({})

const typeOptions = [
  { label: '全部', value: '' },
  { label: 'OpenClaw', value: 'openclaw' },
  { label: 'Hermes', value: 'hermes' },
  { label: 'Claude Code', value: 'claude_code' },
  { label: 'Codex', value: 'codex' },
  { label: '其他', value: 'other' },
]

onMounted(() => {
  fetchAgents()
  fetchTypeCounts()
})

function toggleType(val: string) {
  filterType.value = val
  currentPage.value = 1
  fetchAgents()
}

async function fetchTypeCounts() {
  try {
    const res = await getAgents({ page: 1, page_size: 10000 })
    if (res.data.code === 0 && res.data.data) {
      const items = res.data.data.items || []
      const counts: Record<string, number> = {}
      items.forEach((a: Agent) => {
        const key = ['openclaw', 'hermes', 'claude_code', 'codex'].includes(a.agent_type) ? a.agent_type : 'other'
        counts[key] = (counts[key] || 0) + 1
      })
      counts[''] = items.length
      typeCounts.value = counts
    }
  } catch {}
}

function typeLabel(type: string): string {
  const map: Record<string, string> = {
    openclaw: 'OpenClaw',
    hermes: 'Hermes',
    claude_code: 'Claude Code',
    codex: 'Codex',
    generic: '通用',
  }
  return map[type] || type
}

function agentTypeTag(type: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    openclaw: 'primary',
    hermes: 'success',
    claude_code: 'warning',
    codex: 'danger',
  }
  return map[type] || 'info'
}

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

function isRecent(time: string | null): boolean {
  if (!time) return false
  return Date.now() - new Date(time).getTime() < 300000
}

function formatTime(time: string | null): string {
  if (!time) return '从未'
  const d = new Date(time)
  const now = Date.now()
  const diff = now - d.getTime()
  if (diff < 60000) return `${Math.floor(diff / 1000)} 秒前`
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  return d.toLocaleString()
}

function applyFilters() {
  currentPage.value = 1
  fetchAgents()
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchAgents()
}

async function fetchAgents() {
  loading.value = true
  try {
    const res = await getAgents({
      page: currentPage.value,
      page_size: pageSize.value,
      state: filterState.value || undefined,
      agent_type: (filterType.value && filterType.value !== 'other') ? filterType.value : undefined,
    })
    if (res.data.code === 0 && res.data.data) {
      let items = res.data.data.items || []
      if (filterType.value === 'other') {
        const knownTypes = new Set(['openclaw', 'hermes', 'claude_code', 'codex'])
        items = items.filter((a: Agent) => !knownTypes.has(a.agent_type))
      }
      agents.value = items
      total.value = res.data.data.total || 0
    }
  } catch (e: any) {
    ElMessage.error('加载智能体失败: ' + (e.message || e))
  } finally {
    loading.value = false
  }
}

async function handleConsent(agent: Agent) {
  try {
    await ElMessageBox.confirm(`确定要授权智能体「${agent.name}」吗？`, '授权确认', { type: 'warning' })
    await consentAgent(agent.id, { consent_type: 'full', granted_by: 'user' })
    ElMessage.success('授权成功')
    fetchAgents()
  } catch {}
}

async function handleRebuild(agent: Agent) {
  try {
    await ElMessageBox.confirm(`确定要重建智能体「${agent.name}」吗？（当前状态: ${stateLabel(agent.state)}）`, '重建确认', { type: 'warning' })
    await rebuildAgent(agent.id, 'user')
    ElMessage.success('已提交重建')
    fetchAgents()
  } catch {}
}

async function handlePause(agent: Agent) {
  try {
    await ElMessageBox.confirm(`确定要暂停智能体「${agent.name}」吗？`, '暂停确认', { type: 'warning' })
    await pauseAgent(agent.id, 'user')
    ElMessage.success('已暂停')
    fetchAgents()
  } catch {}
}

async function handleResume(agent: Agent) {
  try {
    await ElMessageBox.confirm(`确定要恢复智能体「${agent.name}」吗？`, '恢复确认', { type: 'info' })
    await resumeAgent(agent.id, 'user')
    ElMessage.success('已恢复')
    fetchAgents()
  } catch {}
}

async function viewHeartbeats(agent: Agent) {
  selectedAgent.value = agent
  hbDialogVisible.value = true
  try {
    const res = await getHeartbeats(agent.id, 20)
    if (res.data.code === 0) {
      heartbeats.value = res.data.data || []
    }
  } catch {
    heartbeats.value = []
  }
}
</script>

<style scoped>
.agent-status-view {
  max-width: 1400px;
  margin: 0 auto;
}

.section-card {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-select {
  width: 150px;
}

/* 类型筛选 */
.type-filter-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.type-filter-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-muted);
  flex-shrink: 0;
}

.type-pills {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.type-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  color: var(--text-secondary);
  background: var(--bg-soft);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
}

.type-pill:hover {
  background: var(--bg-hover);
  border-color: var(--brand);
  color: var(--brand);
}

.type-pill--active {
  background: var(--brand-weak);
  border-color: var(--brand);
  color: var(--brand);
  font-weight: 500;
}

.type-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 18px;
  padding: 0 5px;
  border-radius: 9px;
  font-size: 11px;
  font-weight: 600;
  background: rgba(0, 0, 0, 0.06);
  color: var(--text-muted);
}

.type-pill--active .type-count {
  background: rgba(37, 99, 235, 0.15);
  color: var(--brand);
}

/* 表格 */
.mono-key {
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 12px;
  color: var(--text-secondary);
}

.role-text {
  color: var(--text-secondary);
}

.state-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
}

.state-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.state--running .state-dot { background: #16a34a; box-shadow: 0 0 4px rgba(22,163,74,0.3); }
.state--running { color: #16a34a; }

.state--degraded .state-dot, .state--revoked .state-dot { background: #dc2626; }
.state--degraded, .state--revoked { color: #dc2626; }

.state--paused .state-dot, .state--creating .state-dot { background: #d97706; }
.state--paused, .state--creating { color: #d97706; }

.state--unauthorized .state-dot { background: #7b8799; }
.state--unauthorized { color: #7b8799; }

.state--authorized .state-dot { background: #16a34a; }
.state--authorized { color: #16a34a; }

.heartbeat-old {
  color: var(--text-muted);
}

.action-group {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
}

/* 心跳对话框 */
.hb-item {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
}

.hb-meta {
  font-size: 12px;
  color: var(--text-muted);
}

.hb-error {
  color: #dc2626;
  font-size: 12px;
  margin-top: 4px;
}
</style>
