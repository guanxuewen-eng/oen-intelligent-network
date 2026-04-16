# OEN V8 云服务器部署指南

## 前置条件

- 云服务器（推荐配置：2C4G, 40G SSD, Ubuntu 22.04+）
- Docker + Docker Compose 已安装
- 开放端口：80 (HTTP), 443 (HTTPS, 可选)

## 快速部署

### 1. 上传代码到服务器

```bash
# 在本地执行
scp -r ~/Projects/oen-intelligent-network user@<SERVER_IP>:/opt/oen-intelligent-network
```

或使用 git：

```bash
# 在服务器上执行
cd /opt
git clone <your-repo-url> oen-intelligent-network
cd oen-intelligent-network/src/deploy
```

### 2. 配置环境变量

```bash
cp .env.example .env
nano .env  # 修改数据库密码等配置
```

### 3. 一键部署

```bash
./deploy.sh
```

### 4. 验证部署

```bash
# 健康检查
curl http://<SERVER_IP>/health

# 后端 API
curl http://<SERVER_IP>/api/v1/ops/health

# 前端页面（浏览器）
open http://<SERVER_IP>/
```

## 常用运维指令

```bash
cd /opt/oen-intelligent-network/src/deploy

# 查看日志
docker compose -f docker-compose.prod.yml logs -f backend

# 重启服务
docker compose -f docker-compose.prod.yml restart

# 查看状态
docker compose -f docker-compose.prod.yml ps

# 更新代码后重新部署
git pull && ./deploy.sh

# 停止服务
docker compose -f docker-compose.prod.yml down

# 清理数据（危险操作）
docker compose -f docker-compose.prod.yml down -v
```

## 多 Agent 测试配置

部署完成后，各用户 Agent 可通过以下方式接入：

1. 获取 Agent Key（在 OEN 前端 OpsConsole 创建，或通过 API）
2. Agent 配置指向 `http://<SERVER_IP>/api/v1/`
3. 每个 Agent 独立 `agent_key`，独立心跳和状态

### Agent 注册流程

```bash
# 注册新 Agent
curl -X POST http://<SERVER_IP>/api/v1/agents \
  -H "Content-Type: application/json" \
  -d '{"agent_key":"user-agent-001","name":"User Agent 1","agent_type":"hermes","role":"coder"}'

# 授权
curl -X POST http://<SERVER_IP>/api/v1/agents/1/consent \
  -H "Content-Type: application/json" \
  -d '{"consent_type":"full_access","granted_by":"admin"}'
```

## SSL 证书配置（可选）

```bash
# 放置证书文件
mkdir -p certs
cp your-cert.pem certs/cert.pem
cp your-key.pem certs/key.pem

# 切换为 SSL 配置
cp nginx/default-ssl.conf nginx/default.conf

# 重启
docker compose -f docker-compose.prod.yml up -d nginx
```
