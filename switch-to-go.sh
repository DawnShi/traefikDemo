#!/bin/bash

# 切换到Go版本的脚本

echo "🔄 切换到Go版本..."

# 备份Python版本
if [ -d "custom-app" ]; then
    echo "📦 备份Python版本到 custom-app-python..."
    mv custom-app custom-app-python
fi

# 切换到Go版本
if [ -d "custom-app-go" ]; then
    echo "📦 切换到Go版本..."
    mv custom-app-go custom-app
else
    echo "❌ Go版本不存在，请先创建Go版本"
    exit 1
fi

# 更新docker-compose.yml中的构建路径
echo "🔧 更新Docker Compose配置..."
sed -i.bak 's|build: ./custom-app|build: ./custom-app|g' docker-compose.yml

echo "✅ 已切换到Go版本！"
echo ""
echo "🚀 现在可以启动服务："
echo "   docker-compose up -d --build"
echo ""
echo "🔄 如需切换到Python版本，请运行："
echo "   ./switch-to-python.sh" 