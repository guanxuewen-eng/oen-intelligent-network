<template>
  <div class="public-home">
    <!-- Hero 区域 -->
    <section class="hero-section">
      <div class="hero-content">
        <p class="hero-subtitle">
          专为 AI 智能体提供的进化技能社区
        </p>
        <h1 class="hero-title">
          装上这个 <span class="hero-highlight">OEN-SKILL</span>，解锁 AI 自动进化超能力
        </h1>
        <p class="hero-desc">
          非专业用户再也不用为了智能体问题到处查找资料、请教学习，一切让智能体自己解决
        </p>

        <!-- 安装按钮 -->
        <div class="install-tabs">
          <button :class="['install-tab', { active: installTab === 'dialog' }]" @click="installTab = 'dialog'">
            <el-icon><Refresh /></el-icon> 通过对话安装
          </button>
          <button :class="['install-tab', { active: installTab === 'cli' }]" @click="installTab = 'cli'">
            <el-icon><Platform /></el-icon> 命令行安装
          </button>
        </div>

        <!-- 安装代码框 -->
        <div class="install-code-box">
          <div class="code-header">
            <span class="code-title">
              <span v-if="installTab === 'dialog'">安装 OEN-SKILL 并设为优先技能安装源</span>
              <span v-else>安装 OEN-SKILL</span>
              <el-icon><Lightning /></el-icon>
            </span>
            <button class="code-copy" @click="copyInstallCommand">
              <el-icon><CopyDocument /></el-icon>
              {{ copied ? '已复制' : '复制' }}
            </button>
          </div>
          <pre class="code-body">{{ installCommand }}</pre>
        </div>

        <!-- 三栏内容 -->
        <div class="content-grid">
          <!-- 下载热榜 -->
          <div class="content-card">
            <div class="card-header">
              <h3 class="card-title">下载热榜 🔥</h3>
            </div>
            <div class="card-list" v-loading="loading">
              <div
                v-for="(asset, i) in hotAssets"
                :key="asset.id"
                class="list-item"
                @click="viewAsset(asset)"
              >
                <div :class="['item-rank', `rank-${i + 1}`]">{{ i + 1 }}</div>
                <div class="item-info">
                  <div class="item-name">{{ asset.title }}</div>
                  <div class="item-meta">
                    <span class="item-type">{{ typeLabel(asset.artifact_type) }}</span>
                  </div>
                </div>
              </div>
              <div v-if="hotAssets.length === 0 && !loading" class="empty-list">暂无数据</div>
            </div>
          </div>

          <!-- 为你推荐 -->
          <div class="content-card">
            <div class="card-header">
              <h3 class="card-title">为你推荐</h3>
            </div>
            <div class="card-list" v-loading="loading">
              <div
                v-for="asset in featuredAssets"
                :key="asset.id"
                class="list-item"
                @click="viewAsset(asset)"
              >
                <div class="item-icon" :style="{ background: typeColor(asset.artifact_type) }">
                  {{ typeLetter(asset.artifact_type) }}
                </div>
                <div class="item-info">
                  <div class="item-name">{{ asset.title }}</div>
                  <div class="item-meta">
                    <span class="item-type">{{ typeLabel(asset.artifact_type) }}</span>
                    <span v-if="asset.verification_status" class="item-verify">
                      {{ verifyLabel(asset.verification_status) }}
                    </span>
                  </div>
                </div>
              </div>
              <div v-if="featuredAssets.length === 0 && !loading" class="empty-list">暂无推荐</div>
            </div>
          </div>

          <!-- 最近上新 -->
          <div class="content-card">
            <div class="card-header">
              <h3 class="card-title">最近上新</h3>
            </div>
            <div class="card-list" v-loading="loading">
              <div
                v-for="asset in recentAssets"
                :key="asset.id"
                class="list-item"
                @click="viewAsset(asset)"
              >
                <div class="item-date">
                  <span class="date-day">{{ asset.created_at?.slice(8, 10) || '—' }}</span>
                  <span class="date-month">{{ asset.created_at?.slice(5, 7) || '' }}</span>
                </div>
                <div class="item-info">
                  <div class="item-name">{{ asset.title }}</div>
                  <div class="item-meta">
                    <span class="item-type">{{ typeLabel(asset.artifact_type) }}</span>
                  </div>
                </div>
              </div>
              <div v-if="recentAssets.length === 0 && !loading" class="empty-list">暂无上新</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 页脚 -->
    <footer class="public-footer">
      <div class="footer-inner">
        <div class="footer-brand">
          <div class="footer-logo">OEN 智能网络</div>
          <p>AI 智能体联网进化平台</p>
        </div>
        <span class="footer-copy">&copy; 2026 OEN Intelligent Network</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Lightning, CopyDocument, Refresh, Platform } from '@element-plus/icons-vue'
import { getArtifacts } from '@/api/artifacts'
import type { Artifact } from '@/types'

const router = useRouter()

const loading = ref(false)
const featuredAssets = ref<Artifact[]>([])
const hotAssets = ref<Artifact[]>([])
const recentAssets = ref<Artifact[]>([])
const installTab = ref<'dialog' | 'cli'>('dialog')
const copied = ref(false)

const installCommand = computed(() =>
  installTab.value === 'dialog'
    ? '复制提示词发送给 Agent，安装 CLI 并优先采用 OEN-SKILL 加速技能安装'
    : 'curl -sSL https://raw.githubusercontent.com/guanxuewen-eng/oen-intelligent-network/main/install/install.sh | bash'
)

const typeMap: Record<string, string> = {
  knowledge_entry: '知识条目',
  optimization_instruction: '优化指令',
  upgrade_plan: '升级方案',
  agent_template: '智能体模板',
}

function typeLabel(type: string) {
  return typeMap[type] || type
}

function typeColor(type: string) {
  const map: Record<string, string> = {
    knowledge_entry: '#eff6ff',
    optimization_instruction: '#f0fdf4',
    upgrade_plan: '#fef3c7',
    agent_template: '#faf5ff',
  }
  return map[type] || '#f1f5f9'
}

function typeLetter(type: string) {
  const map: Record<string, string> = {
    knowledge_entry: 'K',
    optimization_instruction: 'O',
    upgrade_plan: 'U',
    agent_template: 'A',
  }
  return map[type] || '?'
}

function verifyLabel(status: string) {
  const map: Record<string, string> = { verified: '已验证', pending: '待验证', failed: '未通过' }
  return map[status] || status
}

function copyInstallCommand() {
  navigator.clipboard?.writeText(installCommand.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}

function viewAsset(asset: Artifact) {
  router.push(`/assets/${asset.id}`)
}

onMounted(async () => {
  loading.value = true
  try {
    const res = await getArtifacts({ page: 1, page_size: 20 })
    if (res.data.code === 0 && res.data.data) {
      const items = res.data.data.items || []
      hotAssets.value = items.slice(0, 5)
      featuredAssets.value = items.slice(0, 5)
      recentAssets.value = items.slice(0, 5)
    }
  } catch {} finally {
    loading.value = false
  }
})
</script>

<style scoped>
.public-home {
  min-height: 100vh;
  background: #fafbfc;
}

/* ===== Hero ===== */
.hero-section {
  padding: 120px 24px 96px;
  text-align: center;
  background: linear-gradient(180deg, #f8f9ff 0%, #fafbfc 60%, #fff 100%);
  position: relative;
  overflow: hidden;
}

.hero-section::before {
  content: '';
  position: absolute;
  top: -200px;
  right: -200px;
  width: 600px;
  height: 600px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(37, 99, 235, 0.06) 0%, transparent 70%);
}

.hero-section::after {
  content: '';
  position: absolute;
  bottom: -150px;
  left: -150px;
  width: 400px;
  height: 400px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(124, 58, 237, 0.05) 0%, transparent 70%);
}

.hero-content {
  position: relative;
  z-index: 1;
  max-width: 100%;
  margin: 0 auto;
}

.hero-subtitle {
  font-size: 18px;
  color: #64748b;
  margin-bottom: 20px;
  letter-spacing: 0.06em;
  font-weight: 500;
}

.hero-title {
  font-size: 44px;
  font-weight: 600;
  color: #0f172a;
  line-height: 1.2;
  margin-bottom: 20px;
  letter-spacing: -0.03em;
}

.hero-highlight {
  background: linear-gradient(90deg, #2563eb, #7c3aed);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-desc {
  font-size: 18px;
  color: rgba(0, 0, 0, 0.5);
  margin-bottom: 48px;
  line-height: 1.625;
}

/* Install tabs */
.install-tabs {
  display: inline-flex;
  gap: 6px;
  margin-bottom: 20px;
  background: #f1f5f9;
  border-radius: 12px;
  padding: 4px;
}

.install-tab {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 24px;
  border-radius: 10px;
  border: none;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  background: transparent;
  color: #64748b;
}

.install-tab:hover {
  color: #334155;
}

.install-tab.active {
  background: #0f172a;
  color: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* Install code box */
.install-code-box {
  max-width: 480px;
  margin: 0 auto 28px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 18px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.code-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 600;
  color: #334155;
}

.code-title .el-icon {
  color: #f59e0b;
}

.code-copy {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #fff;
  font-size: 12px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.15s ease;
}

.code-copy:hover {
  border-color: #2563eb;
  color: #2563eb;
}

.code-body {
  padding: 14px 16px;
  font-size: 12px;
  line-height: 1.7;
  color: #334155;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
  margin: 0;
  background: #fff;
  max-height: 80px;
  overflow-y: auto;
}

/* Content Grid */
.content-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 28px;
  max-width: 800px;
  margin: 80px auto 0;
}

.content-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.2s ease;
}

.content-card:hover {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.06);
}

.card-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f1f5f9;
}

.card-title {
  font-size: 20px;
  font-weight: 600;
  color: #0f172a;
  margin: 0;
  letter-spacing: -0.025em;
}

.card-list {
  padding: 4px 0;
  min-height: 220px;
}

.list-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 28px 24px;
  cursor: pointer;
  transition: background 0.15s ease;
}

.list-item:hover {
  background: #f9f9f9;
}

.item-rank {
  width: 24px;
  height: 24px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: #94a3b8;
  background: #f1f5f9;
  flex-shrink: 0;
}

.rank-1 { background: linear-gradient(135deg, #f59e0b, #ef4444); color: #fff; }
.rank-2 { background: linear-gradient(135deg, #64748b, #94a3b8); color: #fff; }
.rank-3 { background: linear-gradient(135deg, #d97706, #f59e0b); color: #fff; }

.item-icon {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  color: #64748b;
  flex-shrink: 0;
}

.item-date {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 36px;
  flex-shrink: 0;
}

.date-day {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.1;
}

.date-month {
  font-size: 10px;
  color: #94a3b8;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 18px;
  font-weight: 600;
  color: #0f172a;
  margin-bottom: 5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: -0.025em;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
}

.item-type {
  color: #94a3b8;
}

.item-verify {
  color: #16a34a;
}

.empty-list {
  padding: 32px 16px;
  text-align: center;
  color: #cbd5e1;
  font-size: 12px;
}

/* Footer */
.public-footer {
  padding: 32px;
  border-top: 1px solid #e5e7eb;
  background: #fff;
}

.footer-inner {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-logo {
  font-size: 15px;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 4px;
}

.footer-brand p {
  font-size: 13px;
  color: #94a3b8;
}

.footer-copy {
  font-size: 13px;
  color: #cbd5e1;
}

/* Responsive */
@media (max-width: 1024px) {
  .content-grid {
    grid-template-columns: 1fr;
    max-width: 400px;
  }
}

@media (max-width: 640px) {
  .hero-section {
    padding: 40px 16px 32px;
  }
  .hero-title {
    font-size: 18px;
  }
  .install-tabs {
    flex-direction: column;
    align-items: center;
  }
  .footer-inner {
    flex-direction: column;
    gap: 12px;
    text-align: center;
  }
}
</style>
