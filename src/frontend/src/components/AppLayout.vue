<template>
  <el-container class="layout-container">
    <!-- Sidebar -->
    <el-aside :width="store.sidebarCollapsed ? '88px' : '260px'" class="layout-aside">
      <div class="sidebar-brand">
        <div class="brand-badge">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 2L2 7l10 5 10-5-10-5z"/>
            <path d="M2 17l10 5 10-5"/>
            <path d="M2 12l10 5 10-5"/>
          </svg>
        </div>
        <div v-if="!store.sidebarCollapsed" class="brand-copy">
          <div class="brand-main-row">
            <strong class="brand-title-main">智巡</strong>
            <span class="brand-title-en">OEN</span>
          </div>
          <span class="brand-title-sub">智能体联网进化网格</span>
        </div>
      </div>
      <el-menu
        :default-active="activeMenu"
        :collapse="store.sidebarCollapsed"
        router
        class="sidebar-menu"
        background-color="transparent"
        text-color="#35506f"
        active-text-color="#2563eb"
      >
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon>
          <template #title>首页</template>
        </el-menu-item>
        <el-menu-item index="/agents">
          <el-icon><Cpu /></el-icon>
          <template #title>智能体管理中心</template>
        </el-menu-item>
        <el-menu-item index="/artifacts">
          <el-icon><Collection /></el-icon>
          <template #title>资产中心</template>
        </el-menu-item>
        <el-menu-item index="/ops">
          <el-icon><Monitor /></el-icon>
          <template #title>运维控制台</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <!-- Header -->
      <el-header class="layout-header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="store.toggleSidebar()">
            <Fold v-if="!store.sidebarCollapsed" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentTitle">{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <span v-if="store.healthStatus" class="status-dot">
            <span class="status-dot-inner"></span>
            {{ store.healthStatus.service }} v{{ store.healthStatus.version }}
          </span>
        </div>
      </el-header>

      <!-- Main Content -->
      <el-main class="layout-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import {
  HomeFilled,
  Cpu,
  Collection,
  Monitor,
  Fold,
  Expand,
} from '@element-plus/icons-vue'

const route = useRoute()
const store = useAppStore()

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => route.meta.title as string)

onMounted(() => {
  store.fetchHealth()
})
</script>

<style scoped>
.layout-container {
  width: 100%;
  height: 100%;
  background: linear-gradient(180deg, #f8fbff 0%, #f5f7fb 100%);
}

.layout-aside {
  height: auto;
  align-self: stretch;
  margin-bottom: 16px;
  padding: 20px 14px;
  border-right: 1px solid var(--border);
  border-radius: 0 0 14px 14px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
  transition: width 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 0;
  overflow: hidden;
  box-shadow: 2px 0 12px rgba(15, 23, 42, 0.03);
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 8px 10px;
}

.brand-badge {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  color: #fff;
  background: linear-gradient(135deg, #2563eb 0%, #7c3aed 100%);
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

.brand-copy {
  display: flex;
  flex-direction: column;
  gap: 2px;
  white-space: nowrap;
}

.brand-main-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.brand-title-main {
  font-size: 20px;
  font-weight: 700;
  line-height: 1.2;
  color: var(--text-primary);
}

.brand-title-en {
  font-size: 11px;
  font-weight: 500;
  line-height: 1.2;
  color: var(--text-muted);
  letter-spacing: 0.5px;
}

.brand-title-sub {
  font-size: 11px;
  line-height: 1.2;
  color: var(--text-secondary);
}

.sidebar-menu {
  flex: 1;
  min-height: 0;
  border-right: 0;
  overflow: auto;
}

:deep(.el-menu--collapse) {
  width: 64px;
}

:deep(.el-menu--collapse .el-sub-menu__title),
:deep(.el-menu--collapse .el-menu-item) {
  justify-content: center;
}

:deep(.el-menu-item) {
  height: 42px;
  line-height: 42px;
  border-radius: 8px;
  margin: 3px 6px;
  font-size: 14px;
  color: var(--text-secondary);
  transition: all 0.15s ease;
}

:deep(.el-menu-item .el-icon) {
  font-size: 18px;
}

:deep(.el-menu-item:hover) {
  background: rgba(37, 99, 235, 0.06);
  color: var(--text-primary);
}

:deep(.el-menu-item.is-active) {
  background: var(--brand-weak);
  color: var(--brand);
  font-weight: 600;
}

.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border);
  padding: 0 24px;
  height: 56px;
  box-shadow: 0 1px 4px rgba(15, 23, 42, 0.02);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  cursor: pointer;
  font-size: 18px;
  color: var(--text-muted);
  transition: color 0.15s ease;
}

.collapse-btn:hover {
  color: var(--brand);
}

:deep(.el-breadcrumb) {
  font-size: 14px;
}

:deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
  color: var(--text-primary);
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-dot {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-secondary);
}

.status-dot-inner {
  display: inline-block;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--success);
  box-shadow: 0 0 6px rgba(22, 163, 74, 0.4);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.layout-main {
  background: transparent;
  padding: 24px;
}
</style>
