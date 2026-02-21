#!/bin/bash
# RDP Agent Team 快速启动脚本
# 一键启动完整的Agent团队

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔═══════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     RDP Agent Team 快速启动器                    ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════╝${NC}"
echo ""

# 检查OpenCode安装
if ! command -v opencode &> /dev/null; then
    echo -e "${RED}错误: 未找到 opencode 命令${NC}"
    echo "请先安装OpenCode: https://docs.opencode.ai"
    exit 1
fi

# 显示菜单
echo "请选择启动模式:"
echo ""
echo "  1) 🚀 完整模式 - 启动所有Agent (推荐)"
echo "  2) ⚡ 快速模式 - 仅启动协调层 + 1个Feature Agent"
echo "  3) 🎯 单Agent模式 - 调试特定Agent"
echo "  4) 📊 监控模式 - 仅启动监控Dashboard"
echo "  5) ❌ 退出"
echo ""

read -p "请输入选项 (1-5): " choice

case $choice in
    1)
        echo -e "\n${GREEN}启动完整Agent团队...${NC}\n"
        
        # 协调层
        echo "🎯 启动协调层..."
        opencode --session rdp-architect --model claude-opus &
        sleep 2
        opencode --session rdp-pm --model claude-sonnet &
        sleep 2
        opencode --session rdp-reviewer --model claude-sonnet &
        sleep 2
        
        # Phase 1 Feature Agents
        echo "👥 启动Phase 1 Feature Agents..."
        opencode --session rdp-portal --model claude-sonnet &
        sleep 1
        opencode --session rdp-user --model claude-sonnet &
        sleep 1
        opencode --session rdp-project --model claude-sonnet &
        sleep 1
        opencode --session rdp-security --model claude-sonnet &
        sleep 1
        opencode --session rdp-infra --model claude-sonnet &
        
        echo -e "\n${GREEN}✅ 完整Agent团队已启动!${NC}"
        echo ""
        echo "活跃Sessions:"
        opencode --list | grep "rdp-"
        ;;
        
    2)
        echo -e "\n${YELLOW}启动快速模式...${NC}\n"
        
        echo "🎯 启动协调层..."
        opencode --session rdp-architect --model claude-opus &
        sleep 2
        opencode --session rdp-pm --model claude-sonnet &
        sleep 2
        
        echo "👤 启动PortalAgent (示例Feature Agent)..."
        opencode --session rdp-portal --model claude-sonnet &
        
        echo -e "\n${GREEN}✅ 快速模式已启动!${NC}"
        echo ""
        echo "提示: 在PortalAgent Session中输入:"
        echo "  '开始任务 P1-T1 部门门户首页'"
        ;;
        
    3)
        echo ""
        echo "选择要调试的Agent:"
        echo "  1) PortalAgent (门户界面)"
        echo "  2) UserAgent (用户管理)"
        echo "  3) ProjectAgent (项目管理)"
        echo "  4) SecurityAgent (安全合规)"
        echo "  5) InfraAgent (基础设施)"
        echo ""
        read -p "请输入选项 (1-5): " agent_choice
        
        case $agent_choice in
            1) AGENT="portal"; TASK="P1-T1~T4"; DESC="门户界面" ;;
            2) AGENT="user"; TASK="P1-T5~T8"; DESC="用户管理" ;;
            3) AGENT="project"; TASK="P1-T9~T12"; DESC="项目管理" ;;
            4) AGENT="security"; TASK="P1-T13~P1-T16"; DESC="安全合规" ;;
            5) AGENT="infra"; TASK="P0-T0"; DESC="基础设施" ;;
            *) echo "无效选项"; exit 1 ;;
        esac
        
        echo -e "\n${GREEN}启动 ${AGENT}agent...${NC}\n"
        opencode --session rdp-${AGENT} --model claude-sonnet
        ;;
        
    4)
        echo -e "\n${BLUE}启动监控Dashboard...${NC}\n"
        python3 agents/mcp/dashboard.py &
        echo "Dashboard地址: http://localhost:5000"
        echo "API端点: http://localhost:5000/api/agents/status"
        ;;
        
    5)
        echo "退出"
        exit 0
        ;;
        
    *)
        echo "无效选项"
        exit 1
        ;;
esac

echo ""
echo "📋 后续操作:"
echo "  1. 查看所有sessions: opencode --list"
echo "  2. 切换到某个session: opencode --session rdp-portal"
echo "  3. 发送消息给其他Agent: @rdp-pm 任务进度如何?"
echo "  4. 查看任务状态: read agents/WORKSPACE_REGISTRY.md"
echo ""
echo "🎯 Leader Session提示:"
echo "  在当前Session运行: ./scripts/integration-controller.sh"
echo "  自动监控并整合所有Agent工作"
