# OEN Backend - 智能体联网进化网格

中心化治理的智能优化与升级资产平台后端 API 服务。

## 技术栈

- FastAPI + Uvicorn
- SQLAlchemy (异步) + SQLite
- Pydantic V2

## 快速开始

```bash
# 安装依赖
pip install -r requirements.txt

# 启动服务
python -m app.main
```

服务默认运行在 `http://localhost:8000`

## API 文档

启动后访问：
- Swagger UI: `http://localhost:8000/docs`
- ReDoc: `http://localhost:8000/redoc`

## API 接口

### 智能体管理
- `POST /api/v1/agents/register` — 智能体注册
- `GET /api/v1/agents/:id/status` — 查询智能体状态
- `POST /api/v1/agents/:id/heartbeat` — 心跳上报

### 技能资产
- `POST /api/v1/artifacts/upload` — 上传技能/方案
- `GET /api/v1/artifacts/recommend` — 获取推荐技能列表
- `GET /api/v1/artifacts/:id` — 获取技能详情

### 记忆存储
- `POST /api/v1/memories/store` — 存储子代理记忆
- `GET /api/v1/memories/:agent_id` — 查询智能体记忆

## 统一返回格式

```json
{
  "code": 0,
  "data": {},
  "message": "ok"
}
```

## 配置

复制 `.env.example` 为 `.env` 并修改配置：

```
APP_NAME=oen-intelligent-network
APP_ENV=development
DATABASE_URL=sqlite+aiosqlite:///./oen.db
LOG_LEVEL=INFO
```

## 目录结构

```
src/backend/
├── app/
│   ├── __init__.py
│   ├── main.py          # FastAPI 入口
│   ├── config.py        # 配置
│   ├── database.py      # 数据库连接
│   ├── models/          # SQLAlchemy 模型
│   │   ├── agent.py
│   │   ├── artifact.py
│   │   └── memory.py
│   ├── schemas/         # Pydantic 请求/响应模型
│   │   ├── agent.py
│   │   ├── artifact.py
│   │   └── memory.py
│   └── routers/         # API 路由
│       ├── agents.py
│       ├── artifacts.py
│       └── memories.py
├── requirements.txt
├── .env.example
└── README.md
```
