#!/bin/bash
# Agent 启动验证脚本
# 每个Agent启动前必须执行，检查是否满足前提条件

AGENT_NAME=$1
TASK_ID=$2

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

ERRORS=0

echo "======================================"
echo "Agent Startup Validator"
echo "Agent: $AGENT_NAME"
echo "Task: $TASK_ID"
echo "======================================"
echo ""

# Check 1: 验证项目骨架是否存在
check_scaffold() {
    echo -n "Checking project scaffold... "
    if [ -f "apps/web/package.json" ] && [ -f "services/api/go.mod" ]; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗ Missing${NC}"
        echo "  ERROR: Project scaffold not found!"
        echo "  Action: Run 'git status' to check if infra-scaffold is committed"
        ERRORS=$((ERRORS + 1))
    fi
}

# Check 2: 验证Agent是否阅读了协作文档
check_collaboration_guide() {
    echo -n "Checking collaboration guide... "
    if [ -f "agents/outputs/COLLABORATION_GUIDE.md" ]; then
        echo -e "${GREEN}✓${NC}"
        echo "  ✓ Read: agents/outputs/COLLABORATION_GUIDE.md"
        echo "  ✓ Read: agents/outputs/infra-scaffold/CONFLICT_WARNING.md"
    else
        echo -e "${RED}✗ Missing${NC}"
        echo "  ERROR: Collaboration guide not found!"
        ERRORS=$((ERRORS + 1))
    fi
}

# Check 3: 验证依赖Agent是否已完成
check_dependencies() {
    echo "Checking dependencies for $AGENT_NAME..."
    
    # Phase 1 Agents (无依赖，只有infra-scaffold)
    if [[ "$AGENT_NAME" == "PortalAgent" ]] || \
       [[ "$AGENT_NAME" == "UserAgent" ]] || \
       [[ "$AGENT_NAME" == "ProjectAgent" ]] || \
       [[ "$AGENT_NAME" == "SecurityAgent" ]]; then
        echo "  Phase 1 Agent - depends on: infra-scaffold"
        if [ -f "agents/outputs/infra-scaffold/DELIVERY.md" ]; then
            echo -e "  ${GREEN}✓${NC} infra-scaffold: DELIVERY.md found"
        else
            echo -e "  ${RED}✗${NC} infra-scaffold: DELIVERY.md not found"
            ERRORS=$((ERRORS + 1))
        fi
    fi
    
    # Phase 2 Agents (依赖Phase 1)
    if [[ "$AGENT_NAME" == "WorkflowAgent" ]]; then
        echo "  Phase 2 Layer 1 - depends on: ProjectAgent"
        # Check if ProjectAgent has delivered
        if [ -f "agents/outputs/project-agent/DELIVERY.md" ]; then
            echo -e "  ${GREEN}✓${NC} ProjectAgent: Completed"
        else
            echo -e "  ${YELLOW}⚠${NC} ProjectAgent: Not completed yet"
            echo "  WARNING: Consider waiting or coordinate with ProjectAgent"
        fi
    fi
}

# Check 4: 验证工作目录是否清洁
check_git_status() {
    echo -n "Checking git status... "
    if git diff-index --quiet HEAD -- 2>/dev/null; then
        echo -e "${GREEN}✓ Working directory clean${NC}"
    else
        echo -e "${YELLOW}⚠ Uncommitted changes${NC}"
        echo "  Note: You have uncommitted changes. Consider committing first."
    fi
}

# Check 5: 冲突检查
check_conflicts() {
    echo "Checking for potential conflicts..."
    
    # 读取冲突预警表
    if [ -f "agents/outputs/infra-scaffold/CONFLICT_WARNING.md" ]; then
        echo "  ✓ Conflict warning table available"
        
        # 根据Agent名称检查特定路径
        case $AGENT_NAME in
            "PortalAgent")
                echo "  Your workspace: apps/web/src/{pages,components}/portal/"
                echo "  Your workspace: apps/web/src/components/workbench/"
                echo "  Your workspace: apps/web/src/components/notification/"
                echo "  Your workspace: apps/web/src/components/search/"
                ;;
            "UserAgent")
                echo "  Your workspace: services/api/{models,handlers,services}/user*"
                echo "  Your workspace: services/api/{models,handlers,services}/organization*"
                echo "  Your workspace: database/migrations/001_*.sql"
                ;;
            "ProjectAgent")
                echo "  Your workspace: services/api/{models,handlers,services}/project*"
                echo "  Your workspace: database/migrations/002_*.sql"
                ;;
            "SecurityAgent")
                echo "  Your workspace: services/api/{models,handlers,services}/{classification,audit,session}*"
                echo "  Your workspace: services/api/middleware/session.go"
                echo "  Your workspace: database/migrations/003_*.sql"
                echo "  Your workspace: database/migrations/004_*.sql"
                ;;
        esac
    fi
}

# Run all checks
check_scaffold
check_collaboration_guide
check_dependencies
check_git_status
check_conflicts

echo ""
echo "======================================"

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}All checks passed! You can start working.${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Read: agents/outputs/infra-scaffold/DELIVERY.md"
    echo "2. Check: agents/outputs/infra-scaffold/CONFLICT_WARNING.md"
    echo "3. Create: agents/outputs/$AGENT_NAME/DELIVERY.md (before you start)"
    echo ""
    exit 0
else
    echo -e "${RED}Found $ERRORS error(s). Please fix before proceeding.${NC}"
    exit 1
fi
