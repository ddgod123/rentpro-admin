#!/bin/bash

# 启动前端服务脚本
# 功能：启动前端开发服务器

echo "🔄 正在启动前端服务..."

# 1. 切换到前端目录
FRONTEND_DIR="/Users/mac/go/src/rentPro/houduan/rent-foren"
echo "📁 切换到前端目录: $FRONTEND_DIR"
cd "$FRONTEND_DIR" || {
    echo "❌ 错误: 无法切换到前端目录 $FRONTEND_DIR"
    exit 1
}

# 2. 验证package.json文件存在
if [ ! -f "package.json" ]; then
    echo "❌ 错误: 在 $FRONTEND_DIR 中没有找到 package.json 文件"
    exit 1
fi

# 3. 检查node_modules是否存在，如果不存在则安装依赖
if [ ! -d "node_modules" ]; then
    echo "📦 正在安装前端依赖..."
    npm install
fi

# 4. 启动前端开发服务器
echo "🚀 正在启动前端开发服务器..."
echo "   命令: npm run dev"
echo "   目录: $(pwd)"
echo "----------------------------------------"

# 启动前端服务
npm run dev

# 如果需要后台运行，可以使用下面的命令替换上面的命令：
# nohup npm run dev > frontend.log 2>&1 &
# echo "✅ 前端服务已在后台启动，日志保存到 frontend.log"
# echo "📋 查看日志: tail -f frontend.log"
