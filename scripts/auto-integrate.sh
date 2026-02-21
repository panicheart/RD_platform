#!/bin/bash
# Phase 1 自动化整合脚本
# 用法: ./scripts/auto-integrate.sh <agent-branch-name>

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

AGENT_BRANCH=$1

if [ -z "$AGENT_BRANCH" ]; then
    echo -e "${RED}错误: 请提供Agent分支名${NC}"
    echo "用法: ./scripts/auto-integrate.sh feature/portal-phase1"
    exit 1
fi

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}Phase 1 自动整合${NC}"
echo -e "${BLUE}Agent分支: $AGENT_BRANCH${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Step 1: 预检
echo -e "${BLUE}[1/6] 运行预检...${NC}"
if ! ./scripts/integration-precheck.sh; then
    echo -e "${RED}✗ 预检失败，请修复问题后重试${NC}"
    exit 1
fi
echo ""

# Step 2: 获取最新代码
echo -e "${BLUE}[2/6] 获取最新代码...${NC}"
git fetch origin
echo -e "${GREEN}✓${NC} 已获取最新代码"
echo ""

# Step 3: 创建整合分支
echo -e "${BLUE}[3/6] 创建整合分支...${NC}"
INTEGRATION_BRANCH="integration/$(date +%Y%m%d)-$(echo $AGENT_BRANCH | tr '/' '-')"
git checkout -b $INTEGRATION_BRANCH
echo -e "${GREEN}✓${NC} 创建分支: $INTEGRATION_BRANCH"
echo ""

# Step 4: 合并Agent代码
echo -e "${BLUE}[4/6] 合并Agent代码...${NC}"
if git merge origin/$AGENT_BRANCH --no-edit; then
    echo -e "${GREEN}✓${NC} 合并成功，无冲突"
else
    echo -e "${YELLOW}⚠ 检测到冲突，需要手动解决${NC}"
    echo ""
    echo "冲突文件:"
    git diff --name-only --diff-filter=U | sed 's/^/  - /'
    echo ""
    echo "解决步骤:"
    echo "1. 编辑冲突文件，解决冲突"
    echo "2. git add <冲突文件>"
    echo "3. git commit -m 'resolve: merge conflicts from $AGENT_BRANCH'"
    echo ""
    exit 1
fi
echo ""

# Step 5: 验证
echo -e "${BLUE}[5/6] 运行验证...${NC}"

# 检查前端
echo "  检查前端..."
if [ -f "apps/web/package.json" ]; then
    cd apps/web
    if npm run typecheck > /dev/null 2>&1; then
        echo -e "    ${GREEN}✓${NC} TypeScript检查通过"
    else
        echo -e "    ${YELLOW}⚠${NC} TypeScript检查失败"
    fi
    cd ../..
fi

# 检查后端
echo "  检查后端..."
if [ -f "services/api/go.mod" ]; then
    cd services/api
    if go build ./... > /dev/null 2>&1; then
        echo -e "    ${GREEN}✓${NC} Go编译通过"
    else
        echo -e "    ${YELLOW}⚠${NC} Go编译失败"
    fi
    cd ../..
fi

echo ""

# Step 6: 完成
echo -e "${BLUE}[6/6] 整合完成${NC}"
echo ""
echo -e "${GREEN}✓ 整合成功！${NC}"
echo ""
echo "分支: $INTEGRATION_BRANCH"
echo ""
echo "下一步:"
echo "1. 运行测试: make test"
echo "2. 如无问题，合并到main: git checkout main && git merge $INTEGRATION_BRANCH"
echo "3. 推送: git push origin main"
echo ""
