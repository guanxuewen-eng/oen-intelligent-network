# OEN 智能体联网进化网格 — 技术栈选型

> 版本：V1.0
> 日期：2026-04-13
> 编制：关小飞

---

## 一、项目定位

面向智能体的复杂技能中心平台，核心是**服务器中心页面**（类似 Claw Hub），提供：
- Agent 状态监控与管理
- 技能安装、授权、运行
- Artifact 资产沉淀与推荐
- 每日推荐闭环（学习→草稿→推荐→决策）

**不是**传统业务系统，不是论坛/社区，不是高并发电商。

**核心特征：**
- API 驱动，前后端分离
- 重状态管理（Agent 状态机）
- 重数据模型（Artifact 三视图、JSON 内容）
- 推荐引擎（规则优先，后期 AI）
- 审计与安全（授权、撤销、降级）

---

## 二、技术栈选型

### 2.1 后端：Go (Gin + GORM)

| 维度 | Go (Gin) | Spring Boot | Python (FastAPI) |
|------|----------|-------------|------------------|
| 开发效率 | ★★★★ | ★★★ | ★★★★★ |
| 部署复杂度 | 极低（单二进制） | 中（JVM） | 低 |
| 内存占用 | ~30MB | ~300MB | ~80MB |
| 并发性能 | 极好（goroutine） | 好（线程池） | 一般（异步） |
| 类型安全 | 强类型 | 强类型 | 弱类型 |
| Claude Code 支持 | 优秀 | 优秀 | 优秀 |
| 生态成熟度 | 成熟 | 极成熟 | 成熟 |

**选择 Go 的理由：**
1. **部署极简**：编译成单二进制文件，直接跑，无需 JVM/Python 环境
2. **内存友好**：服务器 31G 内存，Go 服务只需 ~30MB，给 AI 模型留足空间
3. **并发天然优势**：Agent 心跳、推荐生成、审计日志写入天然适合 goroutine
4. **结构清晰**：Go 的目录约定优于配置，Claude Code 生成的代码更规范
5. **与公安项目隔离**：Java 跑公安项目，Go 跑 OEN，互不干扰

### 2.2 前端：Vue 3 + TypeScript + Element Plus

- 关总熟悉，组件库成熟
- Element Plus 适合管理后台/数据中心类页面
- TypeScript 保证类型安全

### 2.3 数据库：PostgreSQL

| 维度 | PostgreSQL | MySQL |
|------|-----------|-------|
| JSON 支持 | ★★★★★（jsonb + 索引） | ★★★（JSON 类型） |
| 全文搜索 | ★★★★（内置 tsvector） | ★★★（全文索引） |
| 并发性能 | 极好 | 好 |
| 复杂查询 | 极强 | 强 |
| Artifact 三视图 | 原生 jsonb，查询高效 | 需 JSON 函数 |

**选择 PostgreSQL 的理由：**
1. **Artifact 内容存储**：三视图（Human/Machine/Feed）用 jsonb 存储，原生支持高效查询
2. **推荐引擎**：复杂规则筛选，PostgreSQL 的窗口函数和 CTE 更优雅
3. **审计日志**：分区表支持，适合时间序列数据

### 2.4 缓存：Redis
- 复用服务器现有 Redis
- 用途：会话管理、推荐缓存、限流

### 2.5 部署：Docker Compose
- 后端 + 前端 + PostgreSQL + Redis 一键启动
- 开发/生产环境统一

---

## 三、项目结构

```
oen-intelligent-network/
├── docs/                          # 方案文档
│   ├── 01-tech-stack-selection.md # 技术栈选型（本文档）
│   ├── 02-architecture-design.md  # 架构设计
│   ├── 03-api-contract.md         # API 契约
│   └── tasks.md                   # 任务清单
│
├── src/
│   ├── backend/                   # Go 后端
│   │   ├── cmd/server/            # 入口
│   │   ├── internal/
│   │   │   ├── handler/           # HTTP 处理器
│   │   │   ├── service/           # 业务逻辑
│   │   │   ├── model/             # 数据模型
│   │   │   ├── repository/        # 数据访问
│   │   │   ├── scheduler/         # 定时任务
│   │   │   └── middleware/        # 中间件
│   │   ├── pkg/                   # 公共包
│   │   ├── migrations/            # 数据库迁移
│   │   ├── go.mod
│   │   └── config.yaml
│   │
│   ├── frontend/                  # Vue 3 前端
│   │   ├── src/
│   │   │   ├── api/               # API 请求
│   │   │   ├── views/             # 页面
│   │   │   ├── components/        # 组件
│   │   │   ├── stores/            # Pinia 状态
│   │   │   ├── router/            # 路由
│   │   │   └── types/             # TypeScript 类型
│   │   ├── package.json
│   │   └── vite.config.ts
│   │
│   └── deploy/                    # 部署配置
│       ├── docker-compose.yml
│       ├── Dockerfile.backend
│       └── Dockerfile.frontend
│
├── infra/
│   ├── docker/
│   ├── nginx/
│   └── scripts/
│
├── assets/
│   ├── images/
│   └── docs/
│
├── tests/
└── README.md
```

---

## 四、与 V7/V8 方案的映射

| V7/V8 模块 | 实现位置 | 说明 |
|-----------|---------|------|
| 中心 1.0 | backend/internal/handler + service | 登记、承接、筛选、推荐、留痕 |
| 技能 1.0 | backend/internal/handler + scheduler | 安装、授权、创建、学习、上报 |
| Agent 状态机 | backend/internal/service/agent_service.go | 状态流转 Guard |
| Artifact 管理 | backend/internal/service/artifact_service.go | 四类资产 + 三视图 |
| 推荐引擎 | backend/internal/service/recommendation_service.go | 规则优先，每日推荐 |
| 审计日志 | backend/internal/middleware/audit.go | 全局审计中间件 |
| 前端页面 | frontend/src/views/ | Home、Status Center、Artifact Hub、Ops Console |

---

## 五、服务器规划（118.119.60.205）

| 组件 | 端口 | 内存 | 说明 |
|------|------|------|------|
| OEN 后端 | 8090 | ~50MB | Go 二进制 |
| OEN 前端 | 5180 | ~20MB | Nginx 静态文件 |
| PostgreSQL | 5432 | ~200MB | Docker 容器 |
| Redis | 6379 | ~50MB | 复用现有 |

总内存占用：~320MB，服务器 31G 内存完全够用。
