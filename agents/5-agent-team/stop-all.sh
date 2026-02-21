#!/bin/bash
# 停止所有5个Agent

echo "🛑 停止 RDP 5-Agent Team"
echo "=========================="
echo ""

# 查找并停止opencode进程
sessions=("rdp-pm" "rdp-architect" "rdp-backend" "rdp-frontend" "rdp-devops")

for session in "${sessions[@]}"; do
    echo "🛑 停止 $session..."
    
    # 查找进程并停止
    pids=$(pgrep -f "opencode.*$session" || true)
    
    if [ -n "$pids" ]; then
        echo "   找到PID: $pids"
        kill $pids 2>/dev/null || true
        sleep 1
        
        # 强制停止如果还在运行
        pids=$(pgrep -f "opencode.*$session" || true)
        if [ -n "$pids" ]; then
            kill -9 $pids 2>/dev/null || true
        fi
        
        echo "   ✅ 已停止"
    else
        echo "   ⚠️  未找到进程"
    fi
done

echo ""
echo "=========================="
echo "✅ 所有Agent已停止"
