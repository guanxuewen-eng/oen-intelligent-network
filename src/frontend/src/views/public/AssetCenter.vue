<template>
  <div class="asset-center">
    <!-- 顶部横幅 -->
    <section class="asset-hero">
      <div class="asset-hero-content">
        <h1 class="asset-hero-title">资产中心</h1>
        <p class="asset-hero-desc">浏览平台中所有智能资产，包括知识条目、优化指令、升级方案和智能体模板</p>
      </div>
    </section>

    <!-- 筛选栏 -->
    <section class="filter-bar">
      <div class="filter-inner">
        <div class="filter-group">
          <span class="filter-label">类型</span>
          <div class="filter-chips">
            <span
              v-for="t in typeOptions"
              :key="t.value"
              :class="['filter-chip', { active: filters.type === t.value }]"
              @click="filters.type = t.value; applyFilters()"
            >{{ t.label }}</span>
          </div>
        </div>
        <div class="filter-group">
          <span class="filter-label">风险等级</span>
          <div class="filter-chips">
            <span
              v-for="r in riskOptions"
              :key="r.value"
              :class="['filter-chip', { active: filters.riskLevel === r.value }]"
              @click="filters.riskLevel = r.value; applyFilters()"
            >{{ r.label }}</span>
          </div>
        </div>
        <div class="filter-group">
          <span class="filter-label">验证状态</span>
          <div class="filter-chips">
            <span
              v-for="v in verifyOptions"
              :key="v.value"
              :class="['filter-chip', { active: filters.verifyStatus === v.value }]"
              @click="filters.verifyStatus = v.value; applyFilters()"
            >{{ v.label }}</span>
          </div>
        </div>
        <div class="search-box">
          <el-input
            v-model="filters.search"
            placeholder="搜索资产名称..."
            :prefix-icon="Search"
            clearable
            @clear="applyFilters()"
            @keyup.enter="applyFilters()"
          />
        </div>
      </div>
    </section>

    <!-- 资产列表 -->
    <section class="asset-list">
      <div class="list-header">
        <span class="list-count">共 {{ total }} 项资产</span>
        <div class="view-switch">
          <el-icon :class="['view-btn', { active: viewMode === 'grid' }]" @click="viewMode = 'grid'"><Grid /></el-icon>
          <el-icon :class="['view-btn', { active: viewMode === 'list' }]" @click="viewMode = 'list'"><List /></el-icon>
        </div>
      </div>

      <div v-loading="loading">
        <!-- Grid 模式 -->
        <div v-if="viewMode === 'grid'" class="asset-grid">
          <div
            v-for="item in assets"
            :key="item.id"
            class="asset-card"
            @click="openDetail(item)"
          >
            <div class="card-header">
              <div class="type-badge">{{ typeLabel(item.artifact_type) }}</div>
              <div v-if="item.risk_level" :class="['risk-tag', riskClass(item.risk_level)]">
                {{ riskLabel(item.risk_level) }}
              </div>
            </div>
            <h3 class="card-title">{{ item.title }}</h3>
            <p class="card-desc">{{ item.description || '暂无描述' }}</p>
            <div class="card-footer">
              <span class="footer-item">
                <el-icon><Calendar /></el-icon>
                {{ formatDate(item.created_at) }}
              </span>
              <span v-if="item.target_system" class="footer-item">
                <el-icon><Monitor /></el-icon>
                {{ item.target_system }}
              </span>
            </div>
          </div>
        </div>

        <!-- List 模式 -->
        <el-table
          v-if="viewMode === 'list'"
          :data="assets"
          class="asset-table"
          @row-click="openDetail"
          style="cursor: pointer"
        >
          <el-table-column label="类型" width="120">
            <template #default="{ row }">
              <span class="type-badge-inline">{{ typeLabel(row.artifact_type) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="资产名称" min-width="240">
            <template #default="{ row }">
              <span class="table-title">{{ row.title }}</span>
            </template>
          </el-table-column>
          <el-table-column label="描述" min-width="300">
            <template #default="{ row }">
              <span class="table-desc">{{ row.description || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="风险" width="100" align="center">
            <template #default="{ row }">
              <span v-if="row.risk_level" :class="['risk-tag', riskClass(row.risk_level)]">
                {{ riskLabel(row.risk_level) }}
              </span>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column label="验证状态" width="100" align="center">
            <template #default="{ row }">
              <span v-if="row.verification_status">{{ verifyLabel(row.verification_status) }}</span>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="120" align="center">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrap">
        <el-pagination
          v-if="total > 0"
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          :page-sizes="[8, 16, 32]"
          background
          layout="total, sizes, prev, pager, next"
          @current-change="handlePage"
          @size-change="handleSize"
        />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search, Grid, List, Calendar, Monitor } from '@element-plus/icons-vue'
import { getArtifacts } from '@/api/artifacts'
import type { Artifact } from '@/types'

const route = useRoute()
const router = useRouter()

const assets = ref<Artifact[]>([])
const total = ref(0)
const loading = ref(false)
const viewMode = ref<'grid' | 'list'>('grid')
const page = ref(1)
const pageSize = ref(16)

const filters = ref({
  type: '',
  riskLevel: '',
  verifyStatus: '',
  search: '',
})

const typeOptions = [
  { label: '全部', value: '' },
  { label: '知识条目', value: 'knowledge_entry' },
  { label: '优化指令', value: 'optimization_instruction' },
  { label: '升级方案', value: 'upgrade_plan' },
  { label: '智能体模板', value: 'agent_template' },
]

const riskOptions = [
  { label: '全部', value: '' },
  { label: '低风险', value: 'low' },
  { label: '中风险', value: 'medium' },
  { label: '高风险', value: 'high' },
]

const verifyOptions = [
  { label: '全部', value: '' },
  { label: '已验证', value: 'verified' },
  { label: '待验证', value: 'pending' },
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
  const map: Record<string, string> = { low: '低', medium: '中', high: '高', critical: '严重' }
  return map[level] || level
}

function verifyLabel(status: string) {
  const map: Record<string, string> = { verified: '已验证', pending: '待验证', failed: '未通过' }
  return map[status] || status
}

function formatDate(dateStr: string) {
  return dateStr ? dateStr.slice(0, 10) : '-'
}

function openDetail(item: Artifact) {
  router.push(`/assets/${item.id}`)
}

async function applyFilters() {
  page.value = 1
  await fetchAssets()
}

async function fetchAssets() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: page.value, page_size: pageSize.value }
    if (filters.value.type) params.artifact_type = filters.value.type
    if (filters.value.riskLevel) params.risk_level = filters.value.riskLevel
    if (filters.value.verifyStatus) params.verification_status = filters.value.verifyStatus
    const res = await getArtifacts(params)
    if (res.data.code === 0 && res.data.data) {
      assets.value = res.data.data.items || []
      total.value = res.data.data.total || 0
    }
  } catch {
    assets.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function handlePage(p: number) {
  page.value = p
  fetchAssets()
}

function handleSize(s: number) {
  pageSize.value = s
  page.value = 1
  fetchAssets()
}

onMounted(() => {
  // 从 URL query 读取初始筛选
  if (route.query.type) filters.value.type = route.query.type as string
  fetchAssets()
})
</script>

<style scoped>
/* Hero */
.asset-hero {
  padding: 48px 24px 40px;
  text-align: center;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 40%, #312e81 100%);
  color: #fff;
}

.asset-hero-title {
  font-size: 32px;
  font-weight: 800;
  margin-bottom: 12px;
}

.asset-hero-desc {
  font-size: 14px;
  color: #94a3b8;
  max-width: 480px;
  margin: 0 auto;
}

/* Filter bar */
.filter-bar {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  padding: 20px 24px;
}

.filter-inner {
  max-width: 1280px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-label {
  font-size: 13px;
  color: #64748b;
  font-weight: 500;
  min-width: 64px;
  flex-shrink: 0;
}

.filter-chips {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-chip {
  padding: 4px 14px;
  border-radius: 6px;
  font-size: 13px;
  background: #f1f5f9;
  color: #64748b;
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
}

.filter-chip:hover {
  background: #e2e8f0;
}

.filter-chip.active {
  background: #2563eb;
  color: #fff;
}

.search-box {
  max-width: 280px;
}

/* List area */
.asset-list {
  max-width: 1280px;
  margin: 0 auto;
  padding: 24px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.list-count {
  font-size: 14px;
  color: #64748b;
}

.view-switch {
  display: flex;
  gap: 4px;
}

.view-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: grid;
  place-items: center;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.15s ease;
}

.view-btn:hover {
  color: #1e293b;
  background: #f1f5f9;
}

.view-btn.active {
  color: #2563eb;
  background: #eff6ff;
}

/* Grid */
.asset-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 32px;
}

.asset-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
}

.asset-card:hover {
  border-color: #2563eb;
  box-shadow: 0 8px 24px rgba(37, 99, 235, 0.08);
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.type-badge {
  font-size: 12px;
  padding: 3px 10px;
  border-radius: 4px;
  background: #f1f5f9;
  color: #64748b;
  font-weight: 500;
}

.type-badge-inline {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  background: #f1f5f9;
  color: #64748b;
  font-weight: 500;
}

.risk-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.risk-low { background: #f0fdf4; color: #16a34a; }
.risk-medium { background: #fef3c7; color: #d97706; }
.risk-high { background: #fee2e2; color: #dc2626; }
.risk-critical { background: #fce7f3; color: #be185d; }

.card-title {
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

.card-desc {
  font-size: 13px;
  color: #94a3b8;
  line-height: 1.6;
  margin-bottom: 16px;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  gap: 12px;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f1f5f9;
}

.footer-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #94a3b8;
}

/* Table */
.asset-table {
  margin-bottom: 32px;
}

.table-title {
  font-weight: 500;
  color: #1e293b;
}

.table-desc {
  color: #94a3b8;
  font-size: 13px;
}

.text-muted {
  color: #cbd5e1;
}

/* Pagination */
.pagination-wrap {
  display: flex;
  justify-content: center;
  padding-top: 8px;
}

@media (max-width: 1024px) {
  .asset-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .asset-grid {
    grid-template-columns: 1fr;
  }
  .filter-group {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  .search-box {
    max-width: 100%;
  }
}
</style>
