#!/bin/bash
# OEN 错误检测采集脚本（Hook: PostToolUse）
# 参考 self-improving-agent 的 error-detector.sh

set -e

OUTPUT="${AGENT_TOOL_OUTPUT:-}"

ERROR_PATTERNS=(
    "error:" "Error:" "ERROR:" "failed" "FAILED"
    "command not found" "No such file" "Permission denied"
    "fatal:" "Exception" "Traceback"
    "npm ERR!" "ModuleNotFoundError" "SyntaxError" "TypeError"
)

contains_error=false
for pattern in "${ERROR_PATTERNS[@]}"; do
    if [[ "$OUTPUT" == *"$pattern"* ]]; then
        contains_error=true
        break
    fi
done

if [ "$contains_error" = true ]; then
    LEARNINGS_DIR="${OEN_LEARNINGS_DIR:-./.learnings}"
    mkdir -p "$LEARNINGS_DIR"
    echo "## [ERR-$(date +%Y%m%d)-$(date +%H%M%S)] auto_detected" >> "$LEARNINGS_DIR/ERRORS.md"
    echo "" >> "$LEARNINGS_DIR/ERRORS.md"
    echo "**Logged**: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$LEARNINGS_DIR/ERRORS.md"
    echo "**Status**: pending" >> "$LEARNINGS_DIR/ERRORS.md"
    echo "" >> "$LEARNINGS_DIR/ERRORS.md"
    echo "### Error Output" >> "$LEARNINGS_DIR/ERRORS.md"
    echo '```' >> "$LEARNINGS_DIR/ERRORS.md"
    echo "$OUTPUT" | head -100 >> "$LEARNINGS_DIR/ERRORS.md"
    echo '```' >> "$LEARNINGS_DIR/ERRORS.md"
    echo "" >> "$LEARNINGS_DIR/ERRORS.md"
    echo "---" >> "$LEARNINGS_DIR/ERRORS.md"
    echo "" >> "$LEARNINGS_DIR/ERRORS.md"
    echo "[OEN Collector] Error logged to $LEARNINGS_DIR/ERRORS.md"
fi
