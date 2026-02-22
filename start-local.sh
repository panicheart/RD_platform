#!/bin/bash
#
# 本地一键启动脚本
# Usage: ./start-local.sh

set -e

echo "🚀 RDP 本地部署启动脚本"
echo "========================"
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装"
    echo "请先安装 Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装"
    echo "请先安装 Docker Compose"
    exit 1
fi

echo "✅ Docker 检查通过"
echo ""

# 进入 docker 目录
cd "$(dirname "$0")/deploy/docker"

echo "🐳 启动 Docker 容器..."
docker-compose -f docker-compose.dev.yml down 2>/dev/null || true
docker-compose -f docker-compose.dev.yml up --build -d

echo ""
echo "⏳ 等待服务启动..."
sleep 5

# 检查服务状态
echo ""
echo "🔍 检查服务状态..."

# 检查前端
if curl -s http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ 前端服务: http://localhost:3000"
else
    echo "⏳ 前端服务启动中..."
fi

# 检查后端
if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "✅ 后端 API: http://localhost:8080"
else
    echo "⏳ 后端服务启动中..."
fi

echo ""
echo "🎉 部署完成!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📱 前端界面: http://localhost:3000"
echo "🔌 后端 API: http://localhost:8080"
echo "🗄️  数据库:  localhost:5432"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📋 常用命令:"
echo "  查看日志: docker-compose -f deploy/docker/docker-compose.dev.yml logs -f"
echo "  停止服务: docker-compose -f deploy/docker/docker-compose.dev.yml down"
echo "  重启服务: docker-compose -f deploy/docker/docker-compose.dev.yml restart"
echo ""
echo "⚠️  首次启动需要等待数据库初始化 (约30秒)"
echo ""

# 尝试自动打开浏览器
if command -v open &> /dev/null; then
    sleep 3
    echo "🌐 正在打开浏览器..."
    open http://localhost:3000
elif command -v xdg-open &> /dev/null; then
    sleep 3
    echo "🌐 正在打开浏览器..."
    xdg-open http://localhost:3000
fi
