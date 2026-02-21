#!/bin/bash
# 冲突检测工具
# 检测潜在的代码冲突

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}Phase 1 冲突检测${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# 检测当前工作目录的冲突风险
echo -e "${BLUE}[1/5] 检测高风险文件...${NC}"

HIGH_RISK_FILES=(
    "apps/web/package.json"
    "apps/web/src/types/index.ts"
    "apps/web/src/App.tsx"
    "services/api/go.mod"
    "services/api/main.go"
)

echo "  高风险文件 (容易被多个Agent修改):"
for file in "${HIGH_RISK_FILES[@]}"; do
    if [ -f "$file" ]; then
        # 检查是否有未提交的修改
        if git diff --quiet HEAD -- "$file" 2>/dev/null; then
            echo -e "    ${GREEN}✓${NC} $file (未修改)"
        else
            echo -e "    ${YELLOW}⚠${NC} $file (有修改)"
        fi
    else
        echo -e "    ${RED}✗${NC} $file (不存在)"
    fi
done

echo ""

# 检测文件权限变化
echo -e "${BLUE}[2/5] 检测文件权限...${NC}"
PERMISSION_ISSUES=$(git diff --cached --diff-filter=M --summary 2>/dev/null | grep "mode change" || true)
if [ -n "$PERMISSION_ISSUES" ]; then
    echo -e "  ${YELLOW}⚠ 发现权限变更:${NC}"
    echo "$PERMISSION_ISSUES" | sed 's/^/    /'
else
    echo -e "  ${GREEN}✓${NC} 无权限变更"
fi

echo ""

# 检测大文件
echo -e "${BLUE}[3/5] 检测大文件...${NC}"
LARGE_FILES=$(find . -type f -size +1M ! -path "./.git/*" ! -path "./node_modules/*" 2>/dev/null | head -10)
if [ -n "$LARGE_FILES" ]; then
    echo -e "  ${YELLOW}⚠ 发现大文件 (>1MB):${NC}"
    echo "$LARGE_FILES" | sed 's/^/    /'
else
    echo -e "  ${GREEN}✓${NC} 无大文件"
fi

echo ""

# 检测未跟踪文件
echo -e "${BLUE}[4/5] 检测未跟踪文件...${NC}"
UNTRACKED=$(git ls-files --others --exclude-standard | wc -l)
if [ "$UNTRACKED" -gt 0 ]; then
    echo -e "  ${YELLOW}⚠ 发现 $UNTRACKED 个未跟踪文件${NC}"
    git ls-files --others --exclude-standard | head -5 | sed 's/^/    /'
    if [ "$UNTRACKED" -gt 5 ]; then
        echo "    ... 还有 $((UNTRACKED - 5)) 个"
    fi
else
    echo -e "  ${GREEN}✓${NC} 无未跟踪文件"
fi

echo ""

# 检测潜在冲突模式
echo -e "${BLUE}[5/5] 检测冲突模式...${NC}"
CONFLICT_MARKERS=$(git grep -l "<<<<<<< HEAD" 2>/dev/null || true)
if [ -n "$CONFLICT_MARKERS" ]; then
    echo -e "  ${RED}✗ 发现未解决的冲突标记:${NC}"
    echo "$CONFLICT_MARKERS" | sed 's/^/    /'
else
    echo -e "  ${GREEN}✓${NC} 无冲突标记"
fi

echo ""
echo -e "${BLUE}======================================${NC}"
echo -e "${GREEN}✓ 冲突检测完成${NC}"
echo ""
echo "如果发现问题，请先解决再提交。"
