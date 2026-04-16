#!/bin/bash
set -e

echo "========================================"
echo "  OEN V8 云服务器部署脚本"
echo "========================================"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check prerequisites
check_prereqs() {
    info "检查前置条件..."

    if ! command -v docker &> /dev/null; then
        error "Docker 未安装，请先安装: curl -fsSL https://get.docker.com | sh"
        exit 1
    fi

    if ! command -v docker compose &> /dev/null; then
        error "Docker Compose 未安装"
        exit 1
    fi

    info "Docker $(docker --version)"
    info "Docker Compose $(docker compose version)"
}

# Setup environment
setup_env() {
    info "检查环境配置..."

    if [ ! -f ".env" ]; then
        warn ".env 文件不存在，使用默认配置"
        cp .env.example .env
        warn "请编辑 .env 文件修改数据库密码等配置"
    fi

    source .env
}

# Build and deploy
deploy() {
    info "开始构建和部署..."

    # Pull latest images
    info "拉取基础镜像..."
    docker compose -f docker-compose.prod.yml pull

    # Build application images
    info "构建应用镜像..."
    docker compose -f docker-compose.prod.yml build --no-cache

    # Start services
    info "启动服务..."
    docker compose -f docker-compose.prod.yml up -d

    # Wait for services to be ready
    info "等待服务就绪..."
    sleep 10
}

# Health check
health_check() {
    info "执行健康检查..."

    # Check containers
    echo ""
    docker compose -f docker-compose.prod.yml ps
    echo ""

    # Check backend
    if curl -sf http://localhost:8080/api/v1/ops/health > /dev/null 2>&1; then
        info "后端 API: 正常"
    else
        warn "后端 API: 未响应（可能还在启动中）"
    fi

    # Check frontend via nginx
    if curl -sf http://localhost:80/ > /dev/null 2>&1; then
        info "前端页面: 正常"
    else
        warn "前端页面: 未响应"
    fi

    # Check database
    if docker exec oen-postgres pg_isready -U "${DB_USER:-oen_user}" -d "${DB_NAME:-oen_db}" > /dev/null 2>&1; then
        info "PostgreSQL: 正常"
    else
        warn "PostgreSQL: 未响应"
    fi
}

# Show info
show_info() {
    echo ""
    echo "========================================"
    info "部署完成！"
    echo "========================================"
    echo ""
    echo "  访问地址:"
    echo "    前端:   http://<服务器IP>/"
    echo "    后端:   http://<服务器IP>:8080/api/v1/"
    echo "    健康:   http://<服务器IP>/health"
    echo ""
    echo "  常用指令:"
    echo "    查看日志:  docker compose -f docker-compose.prod.yml logs -f"
    echo "    重启服务:  docker compose -f docker-compose.prod.yml restart"
    echo "    停止服务:  docker compose -f docker-compose.prod.yml down"
    echo "    查看状态:  docker compose -f docker-compose.prod.yml ps"
    echo ""
}

# Main
check_prereqs
setup_env
deploy
health_check
show_info
