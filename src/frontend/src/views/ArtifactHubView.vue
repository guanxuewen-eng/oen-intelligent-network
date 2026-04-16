<template>
  <div class="artifact-hub-view">
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">资产中心</span>
          <el-button type="primary" :icon="Refresh" @click="fetchArtifacts">刷新</el-button>
        </div>
      </template>

      <!-- 筛选 -->
      <div class="filter-row">
        <span class="filter-label">资产类型</span>
        <div class="filter-pills">
          <span
            v-for="t in typeOptions"
            :key="t.value"
            class="filter-pill"
            :class="{ 'filter-pill--active': filterType === t.value }"
            @click="filterType = t.value; currentPage = 1; fetchArtifacts()"
          >
            {{ t.label }}
          </span>
        </div>
      </div>
      <div class="filter-row">
        <span class="filter-label">风险等级</span>
        <div class="filter-pills">
          <span
            v-for="r in riskOptions"
            :key="r.value"
            class="filter-pill"
            :class="{ 'filter-pill--active': filterRisk === r.value }"
            @click="filterRisk = r.value; currentPage = 1; fetchArtifacts()"
          >
            {{ r.label }}
          </span>
        </div>
      </div>

      <el-table :data="artifacts" stripe style="width: 100%" v-loading="loading" class="artifact-table">
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <span class="artifact-title">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="140">
          <template #default="{ row }">
            <el-tag size="small" :type="typeTag(row.artifact_type)" effect="plain">{{ typeLabel(row.artifact_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="风险" width="90" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="getRiskType(row.risk_level)" effect="dark">{{ riskLabel(row.risk_level) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.verification_status)" size="small" effect="plain">
              {{ statusLabel(row.verification_status) || '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="目标系统" width="130">
          <template #default="{ row }">
            <span class="mono-text">{{ row.target_system || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">
            <span class="time-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="100" align="center">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="viewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && artifacts.length === 0" description="暂无资产" />

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

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" :title="selectedArtifact?.title" width="800px" class="artifact-detail">
      <el-tabs v-if="detailData?.versions && detailData.versions.length > 0" v-model="activeVersion" class="version-tabs">
        <el-tab-pane v-for="ver in detailData.versions" :key="ver.id" :name="String(ver.id)">
          <template #label>
            <span class="tab-label">v{{ ver.version_number }}</span>
          </template>
          <div class="version-meta">
            <span class="meta-item"><strong>变更说明：</strong>{{ ver.change_summary || '-' }}</span>
            <span class="meta-item">{{ formatDate(ver.created_at) }} 由 {{ ver.created_by }} 创建</span>
          </div>
          <el-divider />
          <div v-if="ver.views && ver.views.length > 0">
            <el-radio-group v-model="selectedView" size="small">
              <el-radio-button v-for="view in ver.views" :key="view.view_type" :label="view.view_type">
                {{ viewLabel(view.view_type) }}
              </el-radio-button>
            </el-radio-group>
            <pre class="view-content">{{ formatViewContent(currentViewContent) }}</pre>
          </div>
          <el-empty v-else description="此版本暂无视图" />
        </el-tab-pane>
      </el-tabs>
      <el-empty v-else description="暂无可用版本" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { Artifact } from '@/types'
import { getArtifacts, getArtifactDetail } from '@/api/artifacts'

const loading = ref(false)
const filterType = ref('')
const filterRisk = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const artifacts = ref<Artifact[]>([])

const detailVisible = ref(false)
const selectedArtifact = ref<Artifact | null>(null)
const detailData = ref<any>(null)
const activeVersion = ref('')
const selectedView = ref('')

const typeOptions = [
  { label: '全部', value: '' },
  { label: '知识条目', value: 'knowledge_entry' },
  { label: '优化指令集', value: 'optimization_set' },
  { label: '升级方案', value: 'upgrade_plan' },
  { label: '子智能体模板', value: 'subagent_template' },
]

const riskOptions = [
  { label: '全部', value: '' },
  { label: '低', value: 'low' },
  { label: '中', value: 'medium' },
  { label: '高', value: 'high' },
]

const currentViewContent = computed(() => {
  if (!detailData.value?.versions) return ''
  const ver = detailData.value.versions.find((v: any) => String(v.id) === activeVersion.value)
  if (!ver?.views) return ''
  const view = ver.views.find((v: any) => v.view_type === selectedView.value)
  return view?.content || ''
})

onMounted(() => {
  fetchArtifacts()
})

function typeLabel(type: string): string {
  const map: Record<string, string> = {
    knowledge_entry: '知识条目',
    optimization_set: '优化指令',
    upgrade_plan: '升级方案',
    subagent_template: '智能体模板',
  }
  return map[type] || type
}

function typeTag(type: string): 'primary' | 'success' | 'warning' | 'info' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info'> = {
    knowledge_entry: 'primary',
    optimization_set: 'success',
    upgrade_plan: 'warning',
    subagent_template: 'info',
  }
  return map[type] || 'info'
}

function riskLabel(risk: string): string {
  const map: Record<string, string> = { low: '低', medium: '中', high: '高' }
  return map[risk] || risk
}

function statusLabel(status: string): string {
  const map: Record<string, string> = { verified: '已验证', pending: '审核中' }
  return map[status] || status
}

function getRiskType(risk: string): 'success' | 'warning' | 'danger' {
  const map: Record<string, 'success' | 'warning' | 'danger'> = { low: 'success', medium: 'warning', high: 'danger' }
  return map[risk] || 'info'
}

function getStatusType(status: string): 'success' | 'warning' | 'info' {
  const map: Record<string, 'success' | 'warning' | 'info'> = { verified: 'success', pending: 'warning' }
  return map[status] || 'info'
}

function viewLabel(type: string): string {
  const map: Record<string, string> = { human: '人类视图', machine: '机器视图', feed: '推送视图' }
  return map[type] || type
}

function formatDate(date: string): string {
  return new Date(date).toLocaleString()
}

function formatViewContent(content: string): string {
  try {
    const parsed = JSON.parse(content)
    return typeof parsed === 'string' ? parsed : JSON.stringify(parsed, null, 2)
  } catch {
    return content
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchArtifacts()
}

async function fetchArtifacts() {
  loading.value = true
  try {
    const res = await getArtifacts({
      page: currentPage.value,
      page_size: pageSize.value,
      artifact_type: filterType.value || undefined,
      risk_level: filterRisk.value || undefined,
    })
    if (res.data.code === 0 && res.data.data) {
      artifacts.value = res.data.data.items || []
      total.value = res.data.data.total || 0
    }
  } catch (e: any) {
    ElMessage.error('加载资产失败: ' + (e.message || e))
  } finally {
    loading.value = false
  }
}

async function viewDetail(artifact: Artifact) {
  selectedArtifact.value = artifact
  detailVisible.value = true
  try {
    const res = await getArtifactDetail(artifact.id)
    if (res.data.code === 0) {
      detailData.value = res.data.data
      const versions = res.data.data?.versions || []
      if (versions.length > 0) {
        activeVersion.value = String(versions[0].id)
        if (versions[0].views?.length > 0) {
          selectedView.value = versions[0].views[0].view_type
        }
      }
    }
  } catch {}
}
</script>

<style scoped>
.artifact-hub-view {
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

/* 筛选 */
.filter-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.filter-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-muted);
  flex-shrink: 0;
  min-width: 64px;
}

.filter-pills {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-pill {
  display: inline-flex;
  align-items: center;
  padding: 5px 14px;
  border-radius: 20px;
  font-size: 13px;
  color: var(--text-secondary);
  background: var(--bg-soft);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
}

.filter-pill:hover {
  background: var(--bg-hover);
  border-color: var(--brand);
  color: var(--brand);
}

.filter-pill--active {
  background: var(--brand-weak);
  border-color: var(--brand);
  color: var(--brand);
  font-weight: 500;
}

/* 表格 */
.artifact-title {
  color: var(--text-primary);
  font-weight: 500;
}

.mono-text {
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 12px;
  color: var(--text-secondary);
}

.time-text {
  font-size: 12px;
  color: var(--text-muted);
}

/* 详情弹窗 */
.version-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;
  color: var(--text-secondary);
}

.meta-item {
  line-height: 1.5;
}

.view-content {
  background: var(--bg-soft);
  padding: 16px;
  border-radius: var(--radius-md);
  font-size: 13px;
  max-height: 400px;
  overflow-y: auto;
  white-space: pre-wrap;
  color: var(--text-primary);
  border: 1px solid var(--border);
  line-height: 1.6;
}

:deep(.version-tabs .el-tabs__header) {
  margin-bottom: 16px;
}

:deep(.tab-label) {
  font-weight: 500;
}
</style>
