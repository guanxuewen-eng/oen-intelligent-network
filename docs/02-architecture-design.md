# OEN V8 智能体联网进化网格 — 架构设计

> 版本：V1.0
> 日期：2026-04-13
> 编制：关小飞
> 基于：V7 一期边界设计 + V8 总体方案

## 零、智能体兼容性设计

OEN 需要兼容多种主流智能体平台，作为统一的"官方技能 + 进化子代理"中心。

### 0.1 支持的智能体类型

| 类型 | 标识 | 特点 |
|------|------|------|
| **Hermes** | `hermes` | CLI Agent，skill 安装机制，本地运行 |
| **OpenClaw** | `openclaw` | 自有生态，Agent Card 标准 |
| **Claude Code** | `claude_code` | Anthropic 编码 Agent |
| **Codex** | `codex` | OpenAI 编码 Agent |
| **通用** | `generic` | 通过 HTTP API 接入的其他智能体 |

### 0.2 兼容策略

1. **统一 Agent Card**：中心维护统一的智能体名片，对外兼容 A2A / Agent Card 标准
2. **技能安装适配层**：不同智能体的技能安装方式不同，中心提供统一安装入口，内部路由到对应适配
3. **心跳标准化**：各智能体按统一格式上报心跳，中心不关心底层实现
4. **推荐消费适配**：推荐结果按智能体类型转换为对应格式（Hermes skill config、OpenClaw feed 等）
5. **类型无关的核心逻辑**：状态机、Artifact、推荐引擎与智能体类型解耦

### 0.3 Agent 模型扩展

agent 表增加 `agent_type` 字段，标识智能体平台类型：
```
agent_type: hermes | openclaw | claude_code | codex | generic
```

不同智能体的连接方式：
- **Hermes**：通过 Hermes API Server（端口 8642）通信，Bearer Token 认证
- **OpenClaw**：通过 Agent Card + HTTP API
- **Claude Code**：通过 ACP 协议或文件轮询
- **Codex**：通过 OpenAI API 或本地 CLI

---

## 一、核心定位

**OEN 是中心化治理的智能优化与升级资产平台。**

不是论坛、不是社区、不是自动升级工具。

核心逻辑：
```
安装技能 → 授权 → 创建子代理 → 学习+检索 → 沉淀Artifact → 推荐给用户 → 用户确认
```

---

## 二、一期范围（MVP 闭环）

### 2.1 必做

| 能力 | 说明 |
|------|------|
| 官方技能安装与授权 | 用户安装技能、一次性授权 |
| 官方进化子代理创建与保活 | 自动创建子代理，维持运行，心跳检测 |
| 本地学习记录 | 记录本地运行错误、修正记录、优化经验 |
| 在线检索与匿名化摘要 | 子代理联网检索，生成匿名化摘要草稿 |
| Artifact 基础模型 | 四类资产：KnowledgeEntry、OptimizationInstructionSet、UpgradePlan、SubagentTemplate |
| 三视图模型 | Human View、Machine View、OpenClaw Feed View |
| Recommendation 生成与展示 | 从摘要生成推荐，用户可接受/忽略/稍后处理 |
| Agent Status Center | 显示授权状态、子代理状态、心跳、降级事件 |
| 安全骨架 | backup / degraded / local_safe 三种模式 |

### 2.2 不做

| 排除项 | 原因 |
|--------|------|
| 公开评论、帖子流、社交互动 | 不是社区产品 |
| 开放式技能市场交易 | 后期规划 |
| 自动执行关键升级 | 安全第一，推荐层不执行 |
| 多专家评测网络 | V0.1 不需要 |
| 区块链上链 | 仅预留 |

---

## 三、系统架构

### 3.1 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端 (Vue3)                          │
│  Home · Agent Status Center · Artifact Hub · Ops Console    │
├─────────────────────────────────────────────────────────────┤
│                    REST API (Go + Gin)                       │
├──────────┬──────────┬──────────┬──────────┬─────────────────┤
│ Agent    │ Artifact │ Recommend│ Audit    │ Consent         │
│ Service  │ Service  │ Service  │ Service  │ Service         │
├──────────┴──────────┴──────────┴──────────┴─────────────────┤
│                    数据访问层 (GORM)                          │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL (业务数据) · Redis (缓存/会话)                    │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 Agent 状态机

```
[UNAUTHORIZED] ──授权──→ [AUTHORIZED]
                              │
                         创建子代理
                              ▼
                         [CREATING]
                              │
                          创建成功
                              ▼
                         [RUNNING] ◄──── 重建完成 ────┐
                              │                      │
                    ┌─────────┼─────────┐             │
                    │         │         │             │
                 主通道异常  用户暂停   授权撤销       │
                    ▼         ▼         ▼             │
               [DEGRADED] [PAUSED]  [REVOKED]         │
                    │                                  │
              恢复/重建 ───────────────────────────────┘
```

状态机必须由后端统一执行 guard，不允许前端或脚本直接改状态。

---

## 四、核心数据模型

### 4.1 核心表

| 表名 | 说明 |
|------|------|
| agent | 子代理信息、状态、路由模式 |
| agent_heartbeat | 心跳记录 |
| artifact | 资产主表（四类） |
| artifact_version | 资产版本 |
| artifact_view | 三视图（Human/Machine/Feed） |
| candidate_resource | 候选资源池 |
| recommendation | 推荐记录 |
| recommendation_decision | 用户决策 |
| consent_record | 授权记录 |
| audit_log | 统一审计日志 |

### 4.2 Artifact 类型

| 类型 | 说明 |
|------|------|
| knowledge_entry | 知识条目（问题背景、现象、原因分析、处理方法） |
| optimization_set | 优化指令集（操作步骤、依赖、风险提示、验证方式） |
| upgrade_plan | 升级方案（版本差异、实施顺序、影响评估、验收标准） |
| subagent_template | 子代理模板（角色定位、输入输出、工具绑定、约束） |

---

## 五、API 契约（一期）

### 5.1 子代理管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/agents` | 创建子代理 |
| GET | `/api/v1/agents` | 查询列表 |
| GET | `/api/v1/agents/{id}` | 查询详情 |
| POST | `/api/v1/agents/{id}/heartbeat` | 上报心跳 |
| POST | `/api/v1/agents/{id}/consent` | 授权 |
| DELETE | `/api/v1/agents/{id}/consent` | 撤销授权 |
| POST | `/api/v1/agents/{id}/rebuild` | 触发重建 |

### 5.2 Artifact 资产

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/artifacts` | 查询列表 |
| GET | `/api/v1/artifacts/{id}` | 查询详情 |
| GET | `/api/v1/artifacts/{id}/versions` | 版本列表 |
| GET | `/api/v1/artifacts/{id}/versions/{ver}/view/{viewType}` | 指定视图 |

### 5.3 推荐

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/recommendations` | 推荐列表 |
| POST | `/api/v1/recommendations/{id}/decision` | 用户决策 |

### 5.4 候选资源

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/candidates` | 提交候选 |
| POST | `/api/v1/candidates/{id}/review` | 审核候选 |

### 5.5 运维控制台

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/ops/status` | 状态总览 |
| GET | `/api/v1/ops/audit` | 审计日志 |
| GET | `/api/v1/ops/health` | 健康检查 |

---

## 六、前端页面

### 6.1 页面清单

| 页面 | 说明 |
|------|------|
| Home | 平台定位、今日推荐卡片、Agent 运行摘要 |
| Agent Status Center | 状态卡 + 时间线视图（state、route_mode、degraded、recent_errors） |
| Artifact Hub | 统一卡片式列表（四类资产共用结构） |
| Artifact Detail | 三视图切换（Human/Machine/Feed） |
| Ops Console | 治理中心（授权、审计、状态机、降级管理） |

### 6.2 设计原则

- 围绕 Agent Hub + Artifact Hub 设计
- 不做成论坛/社区风格
- 状态清楚，闭环可走通
- 结构正确优先于视觉复杂

---

## 七、安全与治理

1. **授权撤销必须真实生效**：撤销后阻断技能高权限任务
2. **推荐不越权执行**：任何推荐型动作需主任确认
3. **状态机后端 Guard**：不允许前端直接改状态
4. **全量审计日志**：所有重要事件写入 audit_log
5. **降级机制**：primary → backup → local_safe

---

## 八、开发批次

### 第一批：骨架建立（1-2 天）

| 任务 | 内容 |
|------|------|
| T1.1 | 创建工程骨架（Go + Vue3） |
| T1.2 | 数据库建表 + 迁移脚本 |
| T1.3 | 基础 API 骨架（健康检查、基础 CRUD） |
| T1.4 | 前端路由和基础页面壳 |

### 第二批：核心功能（2-3 天）

| 任务 | 内容 |
|------|------|
| T2.1 | Agent 完整 CRUD + 状态机 |
| T2.2 | Artifact 完整 CRUD + 三视图 |
| T2.3 | 推荐引擎基础 |
| T2.4 | 前端页面实现 |

### 第三批：联调闭环（2-3 天）

| 任务 | 内容 |
|------|------|
| T3.1 | 授权→创建→运行→学习 链路 |
| T3.2 | 检索→摘要→Candidate→Artifact 链路 |
| T3.3 | 推荐流程联调 |
| T3.4 | 审计日志验证 |

---

## 九、部署

Docker Compose 一键启动：
- backend (Go)
- frontend (Vue3 + Nginx)
- PostgreSQL
- Redis（复用）
