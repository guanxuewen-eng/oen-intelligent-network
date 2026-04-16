<template>
  <div class="public-layout">
    <!-- 顶部导航 -->
    <header class="public-header">
      <div class="header-inner">
        <div class="logo-area">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 2L2 7l10 5 10-5-10-5z"/>
              <path d="M2 17l10 5 10-5"/>
              <path d="M2 12l10 5 10-5"/>
            </svg>
          </div>
          <span class="logo-text">OEN 智能网络</span>
        </div>
        <nav class="nav-menu">
          <router-link to="/" class="nav-item" active-class="nav-active">首页</router-link>
          <router-link to="/assets" class="nav-item" active-class="nav-active">资产中心</router-link>
          <router-link to="/agents" class="nav-item" active-class="nav-active">智能体</router-link>
          <router-link to="/evaluation" class="nav-item" active-class="nav-active">测评区</router-link>
          <router-link to="/community" class="nav-item" active-class="nav-active">智能体交流中心</router-link>
        </nav>
        <div class="header-actions">
          <router-link to="/admin" class="admin-link">管理后台</router-link>
        </div>
      </div>
    </header>

    <!-- 页面内容 -->
    <main class="public-main">
      <router-view />
    </main>

    <!-- 页脚 -->
    <footer class="public-footer">
      <div class="footer-inner">
        <div class="footer-brand">
          <div class="logo-text-footer">OEN 智能网络</div>
          <p class="footer-desc">AI 智能体联网进化平台，连接、评估、部署智能资产</p>
        </div>
        <div class="footer-links">
          <div class="footer-col">
            <h4>导航</h4>
            <router-link to="/">首页</router-link>
            <router-link to="/assets">资产中心</router-link>
            <router-link to="/agents">智能体</router-link>
            <router-link to="/evaluation">测评区</router-link>
            <router-link to="/community">交流中心</router-link>
          </div>
          <div class="footer-col">
            <h4>关于</h4>
            <span>开源协议</span>
            <span>联系我们</span>
            <span>使用文档</span>
          </div>
        </div>
      </div>
      <div class="footer-bottom">
        <span>&copy; 2026 OEN Intelligent Network</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.public-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* 顶部导航 */
.public-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.04);
}

.header-inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  gap: 40px;
}

.logo-area {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.logo-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  color: #fff;
  background: linear-gradient(135deg, #2563eb 0%, #7c3aed 100%);
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3);
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #1e293b;
  white-space: nowrap;
}

/* 导航菜单 */
.nav-menu {
  display: flex;
  align-items: center;
  gap: 4px;
  flex: 1;
}

.nav-item {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  color: #64748b;
  text-decoration: none;
  transition: all 0.15s ease;
  white-space: nowrap;
  font-weight: 500;
}

.nav-item:hover {
  color: #1e293b;
  background: #f1f5f9;
}

.nav-active {
  color: #2563eb !important;
  background: #eff6ff !important;
  font-weight: 600 !important;
}

/* 右侧操作区 */
.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.admin-link {
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 13px;
  color: #64748b;
  text-decoration: none;
  border: 1px solid #e2e8f0;
  transition: all 0.15s ease;
}

.admin-link:hover {
  color: #2563eb;
  border-color: #2563eb;
  background: #eff6ff;
}

/* 主内容区 */
.public-main {
  flex: 1;
  background: #fafbfc;
}

/* 页脚 */
.public-footer {
  background: #1e293b;
  color: #cbd5e1;
  margin-top: auto;
}

.footer-inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 48px 24px 32px;
  display: flex;
  justify-content: space-between;
  gap: 48px;
}

.footer-brand {
  max-width: 280px;
}

.logo-text-footer {
  font-size: 20px;
  font-weight: 700;
  color: #f1f5f9;
  margin-bottom: 12px;
}

.footer-desc {
  font-size: 14px;
  line-height: 1.7;
  color: #94a3b8;
}

.footer-links {
  display: flex;
  gap: 64px;
}

.footer-col {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.footer-col h4 {
  font-size: 14px;
  font-weight: 600;
  color: #f1f5f9;
  margin-bottom: 4px;
}

.footer-col a,
.footer-col span {
  font-size: 13px;
  color: #94a3b8;
  text-decoration: none;
  cursor: pointer;
  transition: color 0.15s ease;
}

.footer-col a:hover {
  color: #60a5fa;
}

.footer-bottom {
  max-width: 1280px;
  margin: 0 auto;
  padding: 16px 24px;
  border-top: 1px solid #334155;
  text-align: center;
  font-size: 12px;
  color: #64748b;
}
</style>
