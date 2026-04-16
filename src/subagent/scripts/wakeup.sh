#!/bin/bash
# OEN 子代理唤醒脚本
# 触发方式：安装唤醒 / 定时唤醒 / 事件唤醒

set -e

CENTER_URL="${OEN_CENTER_URL:-http://localhost:8000}"
AGENT_ID="${OEN_AGENT_ID:-default}"

echo "[OEN SubAgent] Waking up for $AGENT_ID"

# 采集资源信息
CPU_USAGE=$(top -l 1 | head -5 | grep "CPU usage" | awk '{print $3}' | sed 's/%//' 2>/dev/null || echo "0")
MEM_USAGE=$(vm_stat | awk '/Pages active/ {print $3}' 2>/dev/null || echo "0")

# 上报心跳
curl -s -X POST "$CENTER_URL/api/v1/agents/$AGENT_ID/heartbeat"   -H "Content-Type: application/json"   -d "{"cpu_usage": $CPU_USAGE, "memory_usage": $MEM_USAGE, "status": "active"}"

echo "[OEN SubAgent] Heartbeat sent"
