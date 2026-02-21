#!/bin/bash
# Agent提交监控脚本
# 自动检测Phase 1 Agent的提交状态

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}Phase 1 Agent 提交监控${NC}"
echo -e "${BLUE}$(date)${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# 获取最新信息
git fetch origin --prune 2> /dev/null || true

echo "🔍 检查远程分支..."
echo ""

# 定义要监控的Agent分支
AGENT_PATTERNS=(
    "portal"
    "user"
    "project"
    "security"
    "phase1"
)

FOUND_AGENTS=0

for pattern in "${AGENT_PATTERNS[@]}"; do
    BRANCHES=$(git branch -r | grep -i "$pattern" || true)
    if [ -n "$BRANCHES" ]; then
        echo -e "${GREEN}✓ 发现Agent分支:${NC}"
        echo "$BRANCHES" | sed 's/^/  /'
        FOUND_AGENTS=$((FOUND_AGENTS + 1))
    fi
done

if [ $FOUND_AGENTS -eq 0 ]; then
    echo -e "${YELLOW}⚠ 尚未发现Phase 1 Agent分支${NC}"
    echo ""
    echo "监控中的Agent:"
    echo "  - PortalAgent (预期分支: portal*, feature/portal*)"
    echo "  - UserAgent (预期分支: user*, feature/user*)"
    echo "  - ProjectAgent (预期分支: project*, feature/project*)"
    echo "  - SecurityAgent (预期分支: security*, feature/security*)"
fi

echo ""
echo "📊 当前状态:"
echo "  本地提交数: $(git rev-list --count origin/main..HEAD 2>/dev/null || echo '0')"
echo "  远程提交数: $(git rev-list --count HEAD..origin/main 2>/dev/null || echo '0')"
echo ""

# 检查本地是否有未推送的工作
LOCAL_COMMITS=$(git rev-list --count origin/main..HEAD 2>/dev/null || echo "0")
if [ "$LOCAL_COMMITS" -gt 0 ]; then
    echo -e "${GREEN}✓ 本地有 $LOCAL_COMMITS 个提交等待推送${NC}"
    echo "  最新提交:"
    git log --oneline -3
fi

echo ""
echo -e "${BLUE}======================================${NC}"
echo "下次检查: 每小时或当Agent通知时"
echo ""
echo "当发现Agent分支时，运行:"
echo "  ./scripts/auto-integrate.sh <branch-name>"
echo ""
