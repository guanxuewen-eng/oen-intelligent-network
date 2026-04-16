<template>
  <div class="coming-soon">
    <div class="so-content">
      <div class="so-icon">
        <el-icon :size="48"><DataAnalysis /></el-icon>
      </div>
      <h1 class="so-title">测评区</h1>
      <p class="so-desc">
        对平台中的智能资产进行多维度评估，包括能力指标、兼容性、风险等级和适用场景的自动化测评。
      </p>
      <div class="so-features">
        <div class="feature-card">
          <el-icon :size="20"><TrendCharts /></el-icon>
          <span>能力指标评估</span>
        </div>
        <div class="feature-card">
          <el-icon :size="20"><Warning /></el-icon>
          <span>风险等级分析</span>
        </div>
        <div class="feature-card">
          <el-icon :size="20"><DocumentChecked /></el-icon>
          <span>自动化验证</span>
        </div>
      </div>
      <p class="so-hint">该模块即将上线，敬请期待</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { DataAnalysis, TrendCharts, Warning, DocumentChecked } from '@element-plus/icons-vue'
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
  background: linear-gradient(135deg, #d97706, #f59e0b);
  color: #fff;
  display: grid;
  place-items: center;
  margin: 0 auto 24px;
  box-shadow: 0 8px 24px rgba(217, 119, 6, 0.2);
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
