#!/bin/bash

# Traefik演示项目启动脚本
# 这个脚本会自动配置hosts文件并启动所有服务

echo "🚀 启动Traefik演示项目..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 检查Docker Compose是否可用
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 检查端口是否被占用
echo "🔍 检查端口占用情况..."
if lsof -Pi :80 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  警告: 端口80已被占用"
fi

if lsof -Pi :443 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  警告: 端口443已被占用"
fi

if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  警告: 端口8080已被占用"
fi

# 配置hosts文件
echo "📝 配置hosts文件..."
HOSTS_ENTRIES=(
    "127.0.0.1 traefik.localhost"
    "127.0.0.1 app.localhost"
    "127.0.0.1 app1.localhost"
    "127.0.0.1 app2.localhost"
)

# 检查是否已经配置过hosts
HOSTS_FILE="/etc/hosts"
for entry in "${HOSTS_ENTRIES[@]}"; do
    if ! grep -q "$entry" "$HOSTS_FILE" 2>/dev/null; then
        echo "添加hosts条目: $entry"
        echo "$entry" | sudo tee -a "$HOSTS_FILE" > /dev/null
    else
        echo "hosts条目已存在: $entry"
    fi
done

# 停止现有服务（如果存在）
echo "🛑 停止现有服务..."
docker-compose down

# 构建并启动服务
echo "🔨 构建并启动服务..."
docker-compose up -d --build

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "📊 检查服务状态..."
docker-compose ps

# 显示访问信息
echo ""
echo "🎉 Traefik演示项目启动成功！"
echo ""
echo "📱 访问地址："
echo "   Traefik Dashboard: http://traefik.localhost:8080"
echo "   自定义应用: http://app.localhost"
echo "   Nginx应用: http://app1.localhost"
echo "   Apache应用: http://app2.localhost"
echo ""
echo "🧪 测试API："
echo "   curl http://app.localhost/api/users"
echo "   curl http://app.localhost/health"
echo ""
echo "📋 查看日志："
echo "   docker-compose logs -f"
echo ""
echo "🛑 停止服务："
echo "   docker-compose down"
echo "" 