<template>
  <div class="coming-soon">
    <div class="so-content">
      <div class="so-icon">
        <el-icon :size="48"><ChatDotRound /></el-icon>
      </div>
      <h1 class="so-title">智能体交流中心</h1>
      <p class="so-desc">
        智能体间的协作与通信平台，支持消息队列、事件广播、能力协商和协作任务编排。
      </p>
      <div class="so-features">
        <div class="feature-card">
          <el-icon :size="20"><Connection /></el-icon>
          <span>智能体通信</span>
        </div>
        <div class="feature-card">
          <el-icon :size="20"><Share /></el-icon>
          <span>能力协商</span>
        </div>
        <div class="feature-card">
          <el-icon :size="20"><Operation /></el-icon>
          <span>协作编排</span>
        </div>
      </div>
      <p class="so-hint">该模块即将上线，敬请期待</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ChatDotRound, Connection, Share, Operation } from '@element-plus/icons-vue'
</script>

<style scoped>
.coming-soon {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.so-content {
  text-align: center;
  max-width: 560px;
  padding: 48px;
}

.so-icon {
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: linear-gradient(135deg, #7c3aed, #a78bfa);
  color: #fff;
  display: grid;
  place-items: center;
  margin: 0 auto 24px;
  box-shadow: 0 8px 24px rgba(124, 58, 237, 0.2);
}

.so-title {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 12px;
}

.so-desc {
  font-size: 14px;
  color: #64748b;
  line-height: 1.8;
  margin-bottom: 32px;
}

.so-features {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-bottom: 32px;
}

.feature-card {
  padding: 12px 20px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #64748b;
}

.so-hint {
  font-size: 13px;
  color: #cbd5e1;
}
</style>
