#!/bin/bash
# Phase 1 整合主控脚本
# 主动监控并自动整合Agent提交

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

LOG_FILE="agents/outputs/infra-scaffold/integration.log"

echo "$(date '+%Y-%m-%d %H:%M:%S') - 整合主控启动" >> $LOG_FILE

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}Phase 1 整合主控${NC}"
echo -e "${BLUE}$(date)${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# 配置
CHECK_INTERVAL=30  # 秒
MAX_WAIT_MINUTES=60
AGENT_PATTERNS=("portal" "user" "project" "security" "feature")

INTEGRATED_AGENTS=()
ATTEMPTS=0
MAX_ATTEMPTS=$((MAX_WAIT_MINUTES * 60 / CHECK_INTERVAL))

echo "配置:"
echo "  检查间隔: ${CHECK_INTERVAL}秒"
echo "  最大等待: ${MAX_WAIT_MINUTES}分钟"
echo "  监控Agent: PortalAgent, UserAgent, ProjectAgent, SecurityAgent"
echo ""

while [ $ATTEMPTS -lt $MAX_ATTEMPTS ]; do
    ATTEMPTS=$((ATTEMPTS + 1))
    
    echo "[$ATTEMPTS/$MAX_ATTEMPTS] $(date '+%H:%M:%S') - 检查Agent提交..."
    echo "$(date '+%Y-%m-%d %H:%M:%S') - 检查尝试 $ATTEMPTS" >> $LOG_FILE
    
    # 获取最新分支
    git fetch origin --prune 2> /dev/null || true
    
    # 检测Agent分支
    FOUND_NEW=false
    
    for pattern in "${AGENT_PATTERNS[@]}"; do
        BRANCHES=$(git branch -r | grep -i "$pattern" | sed 's/remotes\///' || true)
        
        if [ -n "$BRANCHES" ]; then
            while IFS= read -r branch; do
                [ -z "$branch" ] && continue
                
                AGENT_NAME=$(echo "$branch" | sed 's/origin\///' | sed 's/feature\///')
                
                # 检查是否已整合
                if [[ " ${INTEGRATED_AGENTS[@]} " =~ " ${AGENT_NAME} " ]]; then
                    continue
                fi
                
                echo ""
                echo -e "${GREEN}🎉 发现Agent提交!${NC}"
                echo "  Agent: $AGENT_NAME"
                echo "  分支: $branch"
                echo ""
                
                # 执行整合
                echo -e "${BLUE}开始自动整合...${NC}"
                echo "$(date '+%Y-%m-%d %H:%M:%S') - 开始整合 $AGENT_NAME" >> $LOG_FILE
                
                if ./scripts/auto-integrate.sh "$branch"; then
                    echo -e "${GREEN}✓ $AGENT_NAME 整合成功${NC}"
                    echo "$(date '+%Y-%m-%d %H:%M:%S') - $AGENT_NAME 整合成功" >> $LOG_FILE
                    INTEGRATED_AGENTS+=("$AGENT_NAME")
                    FOUND_NEW=true
                else
                    echo -e "${RED}✗ $AGENT_NAME 整合失败${NC}"
                    echo "$(date '+%Y-%m-%d %H:%M:%S') - $AGENT_NAME 整合失败" >> $LOG_FILE
                fi
                
            done <<< "$BRANCHES"
        fi
    done
    
    # 检查是否所有Agent都已整合
    if [ ${#INTEGRATED_AGENTS[@]} -ge 4 ]; then
        echo ""
        echo -e "${GREEN}======================================${NC}"
        echo -e "${GREEN}✓ 所有Agent整合完成!${NC}"
        echo -e "${GREEN}======================================${NC}"
        echo ""
        echo "已整合Agent:"
        for agent in "${INTEGRATED_AGENTS[@]}"; do
            echo "  - $agent"
        done
        echo ""
        echo "下一步:"
        echo "  1. 最终验证: make lint && make test"
        echo "  2. 推送到远程: git push origin main"
        echo ""
        
        echo "$(date '+%Y-%m-%d %H:%M:%S') - 所有Agent整合完成" >> $LOG_FILE
        exit 0
    fi
    
    if [ "$FOUND_NEW" = false ]; then
        echo "  未发现新Agent提交 (${#INTEGRATED_AGENTS[@]}/4 已整合)"
        echo "  ${YELLOW}等待 ${CHECK_INTERVAL}秒后重试...${NC}"
        sleep $CHECK_INTERVAL
    fi
done

echo ""
echo -e "${YELLOW}======================================${NC}"
echo -e "${YELLOW}⚠ 等待超时${NC}"
echo -e "${YELLOW}======================================${NC}"
echo ""
echo "已等待 ${MAX_WAIT_MINUTES} 分钟，未检测到所有Agent提交"
echo ""
echo "当前状态:"
echo "  已整合: ${#INTEGRATED_AGENTS[@]}/4"
echo ""
if [ ${#INTEGRATED_AGENTS[@]} -gt 0 ]; then
    echo "已整合Agent:"
    for agent in "${INTEGRATED_AGENTS[@]}"; do
        echo "  - $agent"
    done
    echo ""
fi

echo "选项:"
echo "  1. 继续等待: 重新运行此脚本"
echo "  2. 手动整合: ./scripts/auto-integrate.sh <branch>"
echo "  3. 检查状态: ./scripts/monitor-agent-submission.sh"
echo ""

echo "$(date '+%Y-%m-%d %H:%M:%S') - 等待超时" >> $LOG_FILE
exit 1
