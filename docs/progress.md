# OEN V8 项目进度追踪

> 开始时间：2026-04-13 21:45
> 协调者：关小飞
> 目标：完成第一期（T1.1-T1.5）骨架建立

---

## 当前状态：T3 全链路联调闭环已完成

### T1.1 ~ T1.5 骨架（已完成）
（详见上文）

---

### T3 全链路联调闭环（已完成）
- [x] Agent 创建 → 状态: unauthorized
- [x] Agent 授权 (consent) → 状态: authorized
- [x] Agent 心跳上报 → last_heartbeat 更新
- [x] Artifact 创建 + 版本管理 + 三视图 (Human/Machine/Feed)
- [x] Candidate 创建 → 审核通过 (approve)
- [x] Recommendation 创建 → 决策 (accept → accepted)
- [x] 系统状态总览 (GET /ops/status) 返回 agent_states 聚合
- [x] 审计日志 (GET /ops/audit) 记录全部操作事件
- [x] 前端构建通过 (npm run build)
- [x] 后端构建通过 (go build)

#### T3 联调修复项
- [x] JSONB 字段空字符串问题：Agent/Artifact/AuditLog.Metadata 改为 `*string` 类型
- [x] ArtifactView.Content 改为 `text` 类型（Human View 是纯文本）
- [x] 新增 CreateArtifactView handler + 路由
- [x] 新增 CreateCandidate / CreateRecommendation handler + 路由
- [x] 修复 ArtifactHubView.vue TS 编译错误（row 作用域 + getArtifactDetail 导出）

### T2.1 Agent 完整 CRUD + 状态机（已完成）
- [x] 状态机 guard 实现（UNAUTHORIZED→AUTHORIZED→CREATING→RUNNING→DEGRADED/PAUSED/REVOKED）
- [x] 授权/撤销授权 API（POST/DELETE /agents/:id/consent）
- [x] 心跳上报（POST /agents/:id/heartbeat）+ 自动降级检测
- [x] 重建/暂停/恢复（POST /agents/:id/rebuild/pause/resume）
- [x] 心跳历史查询（GET /agents/:id/heartbeats）

### T2.2 Artifact 完整 CRUD + 三视图（已完成）
- [x] Artifact 版本管理（POST/GET /artifacts/:id/versions）
- [x] 三视图 API（GET /artifacts/:id/versions/:ver/view/:viewType）
- [x] 带版本+视图的完整详情（GET /artifacts/:id/detail）

### T2.3 推荐引擎基础（已完成）
- [x] 候选资源 CRUD + 审核（GET /candidates, POST /candidates/:id/review）
- [x] 推荐列表（GET /recommendations）
- [x] 用户决策（POST /recommendations/:id/decision → accept/ignore/later）
- [x] 候选→推荐生成链路

### T2.4 前端页面实现（已完成）
- [x] HomeView 对接真实 API（系统状态、推荐列表、决策操作）
- [x] AgentStatusView 对接真实 API（状态筛选、授权/重建/暂停/恢复、心跳历史）
- [x] ArtifactHubView 对接真实 API（类型/风险筛选、版本+三视图详情弹窗）
- [x] OpsConsoleView 对接真实 API（状态总览、审计日志）

---

### 修改文件清单（T2）

#### 新增文件（6 个）
- `src/backend/internal/service/agent_service.go` — 状态机 + 授权 + 心跳 + 重建
- `src/backend/internal/service/artifact_service.go` — 版本 + 三视图服务
- `src/backend/internal/service/recommend_service.go` — 候选 + 推荐 + 决策服务
- `src/backend/internal/service/ops_service.go` — 状态总览 + 审计
- `src/backend/internal/handler/agent_handler.go` — 新 API 端点 handler
- `src/frontend/src/api/recommendations.ts` — 推荐/候选/运维 API 客户端

#### 修改文件（6 个）
- `src/backend/cmd/server/main.go` — 新增全部路由
- `src/backend/internal/handler/handler.go` — SystemStatus 改用新 API，新增 ListAuditLogs
- `src/backend/internal/repository/agent.go` — 补回 Agent CRUD + 新增全部仓储方法
- `src/frontend/src/api/agents.ts` — 新增状态机/心跳 API
- `src/frontend/src/views/AgentStatusView.vue` — 真实 API 对接 + 心跳历史弹窗
- `src/frontend/src/views/HomeView.vue` — 真实 API 对接 + 推荐决策
- `src/frontend/src/views/ArtifactHubView.vue` — 真实 API 对接 + 版本/三视图弹窗
- `src/frontend/src/views/OpsConsoleView.vue` — 审计日志 API 对接

---

## 当前状态：第一期已完成

### T1.1 Go 工程骨架
- [x] Gin 路由框架搭建
- [x] GORM 数据库连接
- [x] 配置管理 (config.yaml)
- [x] CORS/日志中间件
- [x] 数据模型定义（10 张表）
- [x] 路由替换占位符为真实 handler

### T1.2 数据库建表+迁移
- [x] GORM 模型定义
- [x] AutoMigrate 实现（database.go 中已有）
- [x] 数据库连接验证（通过 GORM AutoMigrate）

### T1.3 基础 API 实现
- [x] Agent CRUD 完整实现
  - GET /api/v1/agents（分页列表）
  - GET /api/v1/agents/:id（详情）
  - POST /api/v1/agents（创建）
  - PUT /api/v1/agents/:id（更新）
  - DELETE /api/v1/agents/:id（删除）
- [x] Artifact CRUD 完整实现
  - GET /api/v1/artifacts（分页列表）
  - GET /api/v1/artifacts/:id（详情）
  - POST /api/v1/artifacts（创建）
  - PUT /api/v1/artifacts/:id（更新）
  - DELETE /api/v1/artifacts/:id（删除）
- [x] 健康检查/错误处理
  - GET /api/v1/ops/health（保持现有）
  - GET /api/v1/ops/status（系统状态总览）

### T1.4 前端页面填充
- [x] 路由配置
- [x] AppLayout 布局
- [x] HomeView 内容（平台定位卡片、推荐摘要、状态概览）
- [x] AgentStatusView 内容（Agent 列表表格、状态筛选、分页）
- [x] ArtifactHubView 内容（Artifact 列表、类型/风险筛选、分页）
- [x] OpsConsoleView 内容（系统状态面板、API 端点列表、审计日志表格）
- [x] API 客户端层预留（agents.ts、artifacts.ts）
- [x] Mock 数据填充

### T1.5 Docker Compose 开发环境
- [x] docker-compose.yml 配置
- [x] 后端 Dockerfile
- [x] 前端 Nginx 配置
- [x] 一键启动验证（配置正确）

---

## 修改文件清单

### 新增文件（7 个）
- `src/backend/internal/repository/artifact.go` — Artifact 仓储层 CRUD
- `src/backend/internal/handler/handler.go` — 重写 Handler（Agent/Artifact CRUD + Ops Status）
- `src/frontend/src/api/agents.ts` — Agent API 客户端
- `src/frontend/src/api/artifacts.ts` — Artifact API 客户端

### 修改文件（7 个）
- `src/backend/cmd/server/main.go` — 替换占位符路由为真实 handler，添加 /ops/status
- `src/backend/internal/service/service.go` — 实现 Agent/Artifact/Ops 服务层
- `src/frontend/src/types/index.ts` — Agent 类型添加 agent_type 字段
- `src/frontend/src/views/HomeView.vue` — 填充平台定位、推荐摘要、状态概览（mock）
- `src/frontend/src/views/AgentStatusView.vue` — 填充 Agent 列表表格（mock）
- `src/frontend/src/views/ArtifactHubView.vue` — 填充 Artifact 列表卡片（mock）
- `src/frontend/src/views/OpsConsoleView.vue` — 填充系统状态面板、审计日志（mock）

---

## 验收结果
1. [x] `go build ./cmd/server/` 编译通过
2. [x] `cd src/frontend && npm run build` 构建通过
3. [x] `docker-compose up -d` 配置正确，可一键启动
4. [x] 后端 /api/v1/ops/health 返回 200
5. [x] 前端页面正常访问并展示 mock 数据
6. [x] Agent CRUD API 可通过 curl 测试

---

## 启动和测试指令

### 启动完整环境（Docker）
```bash
cd ~/Projects/oen-intelligent-network/src/deploy
docker-compose up -d
```

### 单独启动后端（开发模式）
```bash
cd ~/Projects/oen-intelligent-network/src/backend
# 确保 PostgreSQL 已运行
go run ./cmd/server/
```

### 单独启动前端（开发模式）
```bash
cd ~/Projects/oen-intelligent-network/src/frontend
npm run dev
```

### 测试指令

# 健康检查
curl http://localhost:8080/api/v1/ops/health

# 系统状态
curl http://localhost:8080/api/v1/ops/status

# 创建 Agent
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Content-Type: application/json" \
  -d '{"agent_key":"test-agent-001","name":"Test Agent","agent_type":"hermes","role":"coder"}'

# 列出 Agent
curl http://localhost:8080/api/v1/agents

# 获取 Agent 详情
curl http://localhost:8080/api/v1/agents/1

# 更新 Agent
curl -X PUT http://localhost:8080/api/v1/agents/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Test Agent","state":"active"}'

# 删除 Agent
curl -X DELETE http://localhost:8080/api/v1/agents/1

# 创建 Artifact
curl -X POST http://localhost:8080/api/v1/artifacts \
  -H "Content-Type: application/json" \
  -d '{"artifact_key":"test-artifact-001","artifact_type":"script","title":"Test Script","risk_level":"low"}'

# 列出 Artifact
curl http://localhost:8080/api/v1/artifacts

### 停止环境
```bash
cd ~/Projects/oen-intelligent-network/src/deploy
docker-compose down
```

---

## 下一步
1. 实现 Agent 心跳上报和在线状态检测
2. 完善 Artifact 版本管理（artifact_version 表）
3. 实现 Recommendation 推荐系统
4. 前端对接真实 API（替换 mock 数据）
5. 添加用户认证和权限管理
