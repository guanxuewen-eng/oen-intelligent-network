<template>
  <div class="app-container">
    <AppLayout />
  </div>
</template>

<script setup lang="ts">
import AppLayout from './components/AppLayout.vue'
</script>

<style>
:root {
  --bg-page: #f5f7fb;
  --bg-panel: #ffffff;
  --bg-soft: #f8fafc;
  --bg-hover: #eef4ff;
  --border: #dbe3ef;
  --border-strong: #c8d4e3;
  --text-primary: #162033;
  --text-secondary: #4b5b74;
  --text-muted: #7b8799;
  --brand: #2563eb;
  --brand-weak: #eff6ff;
  --success: #16a34a;
  --warning: #d97706;
  --danger: #dc2626;
  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --shadow-soft: 0 8px 24px rgba(15, 23, 42, 0.06);
  --shadow-card: 0 1px 3px rgba(15, 23, 42, 0.04), 0 1px 2px rgba(15, 23, 42, 0.02);
  --shadow-elevated: 0 12px 32px rgba(15, 23, 42, 0.08);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  width: 100%;
  height: 100%;
  background: var(--bg-page);
  color: var(--text-primary);
  font-family: "Inter", "PingFang SC", "Microsoft YaHei", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.app-container {
  width: 100%;
  height: 100%;
}

/* Element Plus 全局覆盖 */
.el-card {
  border-radius: var(--radius-lg) !important;
  border: 1px solid var(--border) !important;
  box-shadow: var(--shadow-card) !important;
  overflow: hidden;
}

.el-card__header {
  border-bottom: 1px solid var(--border) !important;
  padding: 16px 20px !important;
}

.el-card__body {
  padding: 20px !important;
}

.el-table th.el-table__cell {
  background-color: var(--bg-soft) !important;
  color: var(--text-secondary) !important;
  font-weight: 600 !important;
  font-size: 13px !important;
}

.el-table td.el-table__cell {
  font-size: 13px !important;
}

.el-table__row:hover > td.el-table__cell {
  background-color: var(--bg-hover) !important;
}

.el-pagination {
  font-size: 13px !important;
}

.el-pagination.is-background .el-pager li {
  border-radius: var(--radius-sm) !important;
}

.el-dialog {
  border-radius: var(--radius-lg) !important;
  overflow: hidden;
}

.el-button {
  border-radius: var(--radius-sm) !important;
}

.el-tag {
  border-radius: 4px !important;
  border: none !important;
  font-weight: 500 !important;
}

.el-empty {
  padding: 40px 0 !important;
}

/* 全局动画 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.el-card {
  animation: fadeInUp 0.35s ease both;
}

.el-card:nth-child(2) {
  animation-delay: 0.06s;
}

.el-card:nth-child(3) {
  animation-delay: 0.12s;
}

.el-row .el-col:nth-child(1) > .el-card { animation-delay: 0.06s; }
.el-row .el-col:nth-child(2) > .el-card { animation-delay: 0.12s; }
.el-row .el-col:nth-child(3) > .el-card { animation-delay: 0.18s; }
</style>
