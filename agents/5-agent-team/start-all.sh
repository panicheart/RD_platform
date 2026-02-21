#!/bin/bash
# 一键启动5个Agent
# 用法: ./start-all.sh

set -e

PROJECT_ROOT="/Users/tancong/Code/RD_platform"
AGENT_DIR="$PROJECT_ROOT/agents/5-agent-team"
LOG_DIR="$PROJECT_ROOT/agents/outputs/logs"

echo "🚀 RDP 5-Agent Team 启动器"
echo "=========================="
echo ""

# 创建日志目录
mkdir -p "$LOG_DIR"

# 检查opencode是否安装
if ! command -v opencode &> /dev/null; then
    echo "❌ 错误: opencode 命令未找到"
    echo "请先安装 OpenCode: https://github.com/opencode-ai/opencode"
    exit 1
fi

echo "✅ OpenCode 已安装"
echo ""

# 初始化任务数据库
echo "📋 初始化任务数据库..."
cd "$PROJECT_ROOT"
python3 "$AGENT_DIR/coordinator.py" init 2>/dev/null || echo "数据库已存在，跳过初始化"
echo ""

# 启动函数
start_agent() {
    local name=$1
    local session=$2
    local model=$3
    local log_file="$LOG_DIR/${session}.log"
    
    echo "🚀 启动 $name (session: $session, model: $model)..."
    
    # 使用nohup在后台启动，输出到日志文件
    nohup opencode --session "$session" --model "$model" --working-dir "$PROJECT_ROOT" > "$log_file" 2>&1 &
    
    echo "   PID: $!"
    echo "   日志: $log_file"
    sleep 1
}

# 启动5个Agent
start_agent "PM-Agent" "rdp-pm" "claude-sonnet"
start_agent "Architect-Agent" "rdp-architect" "claude-sonnet"
start_agent "Backend-Agent" "rdp-backend" "claude-sonnet"
start_agent "Frontend-Agent" "rdp-frontend" "claude-sonnet"
start_agent "DevOps-Agent" "rdp-devops" "claude-sonnet"

echo ""
echo "=========================="
echo "✅ 5个Agent已启动"
echo ""
echo "📊 查看状态:"
echo "  任务看板: make agent-team-status"
echo "  查看日志: tail -f $LOG_DIR/*.log"
echo ""
echo "📝 各Agent启动指令已保存到:"
echo "  $AGENT_DIR/instructions/"
echo ""
echo "💡 提示:"
echo "  1. 使用 'opencode --session rdp-pm' 进入PM-Agent会话"
echo "  2. 在Agent会话中粘贴对应指令启动工作"
echo "  3. 使用 './agents/5-agent-team/stop-all.sh' 停止所有Agent"
