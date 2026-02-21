#!/bin/bash
# Phase 1 整合预检脚本
# 在整合前运行，检查环境是否就绪

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "======================================"
echo "Phase 1 Integration Pre-check"
echo "======================================"
echo ""

ERRORS=0
WARNINGS=0

# Check 1: Git状态
echo "🔍 Checking git status..."
if git diff-index --quiet HEAD --; then
    echo -e "${GREEN}✓${NC} Working directory is clean"
else
    echo -e "${YELLOW}⚠${NC} Uncommitted changes detected"
    git status --short
    WARNINGS=$((WARNINGS + 1))
fi

# Check 2: Remote updates
echo ""
echo "🔍 Checking for remote updates..."
git fetch origin
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse @{u} 2>/dev/null || echo "")

if [ -z "$REMOTE" ]; then
    echo -e "${YELLOW}⚠${NC} No upstream branch configured"
elif [ "$LOCAL" = "$REMOTE" ]; then
    echo -e "${GREEN}✓${NC} Up to date with remote"
else
    echo -e "${YELLOW}⚠${NC} Remote has updates"
    echo "  Behind by: $(git rev-list --count HEAD..@{u}) commits"
    WARNINGS=$((WARNINGS + 1))
fi

# Check 3: InfraAgent commits
echo ""
echo "🔍 Checking InfraAgent commits..."
COMMIT_COUNT=$(git rev-list --count origin/main..HEAD 2>/dev/null || echo "0")
if [ "$COMMIT_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓${NC} $COMMIT_COUNT commits ready to integrate"
    echo "  Latest: $(git log -1 --pretty=format:'%s' HEAD)"
else
    echo -e "${RED}✗${NC} No commits found"
    ERRORS=$((ERRORS + 1))
fi

# Check 4: Check for Agent branches
echo ""
echo "🔍 Checking for Phase 1 Agent branches..."
AGENT_BRANCHES=$(git branch -r | grep -E "(portal|user|project|security)" | wc -l)
if [ "$AGENT_BRANCHES" -gt 0 ]; then
    echo -e "${GREEN}✓${NC} Found $AGENT_BRANCHES Agent branches:"
    git branch -r | grep -E "(portal|user|project|security)" | sed 's/^/  - /'
else
    echo -e "${YELLOW}⚠${NC} No Agent branches found"
    echo "  Phase 1 Agents haven't submitted yet"
fi

# Check 5: Environment readiness
echo ""
echo "🔍 Checking environment..."

# Node.js
if command -v node &> /dev/null; then
    NODE_VERSION=$(node -v)
    echo -e "${GREEN}✓${NC} Node.js: $NODE_VERSION"
else
    echo -e "${RED}✗${NC} Node.js not found"
    ERRORS=$((ERRORS + 1))
fi

# Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓${NC} Go: $GO_VERSION"
else
    echo -e "${RED}✗${NC} Go not found"
    ERRORS=$((ERRORS + 1))
fi

# npm
if command -v npm &> /dev/null; then
    NPM_VERSION=$(npm -v)
    echo -e "${GREEN}✓${NC} npm: $NPM_VERSION"
else
    echo -e "${RED}✗${NC} npm not found"
    ERRORS=$((ERRORS + 1))
fi

# Check 6: Project structure
echo ""
echo "🔍 Checking project structure..."
REQUIRED_DIRS=("apps/web" "services/api" "database" "deploy" "agents/outputs")
for dir in "${REQUIRED_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        echo -e "${GREEN}✓${NC} $dir/"
    else
        echo -e "${RED}✗${NC} $dir/ missing"
        ERRORS=$((ERRORS + 1))
    fi
done

# Check 7: Key files
echo ""
echo "🔍 Checking key files..."
REQUIRED_FILES=(
    "apps/web/package.json"
    "services/api/go.mod"
    "Makefile"
    "agents/WORKSPACE_REGISTRY.md"
)
for file in "${REQUIRED_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓${NC} $file"
    else
        echo -e "${RED}✗${NC} $file missing"
        ERRORS=$((ERRORS + 1))
    fi
done

# Summary
echo ""
echo "======================================"
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✓ Pre-check passed${NC}"
    echo ""
    echo "Ready for Phase 1 integration!"
    echo ""
    echo "Next steps when Agent submits:"
    echo "1. git fetch origin"
    echo "2. git merge <agent-branch>"
    echo "3. Resolve conflicts if any"
    echo "4. Run: make lint && make test"
    echo "5. git push origin main"
    exit 0
else
    echo -e "${RED}✗ Pre-check failed${NC}"
    echo "  Errors: $ERRORS"
    echo "  Warnings: $WARNINGS"
    echo ""
    echo "Please fix errors before integration"
    exit 1
fi
