#!/bin/bash
set -e

OEN_API_URL="http://183.223.249.216:58000/api/v1"
OEN_SITE_URL="http://183.223.249.216:58000"

echo "======================================"
echo "  OEN-SKILL 安装程序"
echo "======================================"
echo ""

# Detect platform
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$PLATFORM" in
  darwin)  PLATFORM="darwin" ;;
  linux)   PLATFORM="linux" ;;
  *)       echo "不支持的操作系统: $PLATFORM"; exit 1 ;;
esac

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  arm64)   ARCH="arm64" ;;
  aarch64) ARCH="arm64" ;;
  *)       echo "不支持的架构: $ARCH"; exit 1 ;;
esac

echo "平台: $PLATFORM/$ARCH"
echo "安装目录: /usr/local/bin"
echo ""

# Create skill directory
mkdir -p "$HOME/.oen-skill"

# Download OEN-SKILL CLI
CLI_NAME="oen-skill-${PLATFORM}-${ARCH}"
if [ "$PLATFORM" = "linux" ]; then
  CLI_URL="${OEN_SITE_URL}/install/${CLI_NAME}"
else
  CLI_URL="${OEN_SITE_URL}/install/${CLI_NAME}"
fi

echo "正在下载 OEN-SKILL CLI..."
if command -v curl &> /dev/null; then
  curl -sSL -o "$HOME/.oen-skill/oen-skill" "$CLI_URL" || {
    echo "错误: CLI 下载失败"
    exit 1
  }
elif command -v wget &> /dev/null; then
  wget -qO "$HOME/.oen-skill/oen-skill" "$CLI_URL" || {
    echo "错误: CLI 下载失败"
    exit 1
  }
else
  echo "错误: 需要 curl 或 wget"
  exit 1
fi

chmod +x "$HOME/.oen-skill/oen-skill"

# Create symlink
sudo ln -sf "$HOME/.oen-skill/oen-skill" /usr/local/bin/oen-skill 2>/dev/null || \
  ln -sf "$HOME/.oen-skill/oen-skill" "$HOME/.local/bin/oen-skill" 2>/dev/null || \
  echo "(请将 $HOME/.oen-skill/oen-skill 加入 PATH)"

# Save config
cat > "$HOME/.oen-skill/config.json" << EOF
{
  "api_url": "${OEN_API_URL}",
  "site_url": "${OEN_SITE_URL}"
}
EOF

echo ""
echo "======================================"
echo "  安装完成"
echo "======================================"
echo ""
echo "配置信息:"
echo "  API 地址: $OEN_API_URL"
echo "  网站地址: $OEN_SITE_URL"
echo "  技能目录: $HOME/.oen-skill"
echo ""
echo "使用方法:"
echo "  oen-skill list     # 列出可用技能"
echo "  oen-skill install  # 安装最新技能"
echo "  oen-skill status   # 查看服务器状态"
echo ""
