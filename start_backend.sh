#!/bin/bash

# 启动后端服务脚本
# 功能：先杀死8002端口进程，然后启动后端服务

echo "🔄 正在启动后端服务..."

# 1. 杀死8002端口上的进程
echo "🛑 正在检查并杀死8002端口进程..."
PORT_PID=$(lsof -ti:8002)
if [ -n "$PORT_PID" ]; then
    echo "   发现8002端口进程: $PORT_PID"
    kill -9 $PORT_PID
    echo "   ✅ 已杀死8002端口进程"
    sleep 1
else
    echo "   ℹ️  8002端口没有运行的进程"
fi

# 2. 杀死所有go run main.go进程（防止有残留进程）
echo "🛑 正在杀死残留的go run进程..."
pkill -f "go run main.go" 2>/dev/null || echo "   ℹ️  没有发现go run残留进程"

# 3. 切换到后端目录
BACKEND_DIR="/Users/mac/go/src/rentPro/houduan/rentpro-admin-main"
echo "📁 切换到后端目录: $BACKEND_DIR"
cd "$BACKEND_DIR" || {
    echo "❌ 错误: 无法切换到后端目录 $BACKEND_DIR"
    exit 1
}

# 4. 验证main.go文件存在
if [ ! -f "main.go" ]; then
    echo "❌ 错误: 在 $BACKEND_DIR 中没有找到 main.go 文件"
    exit 1
fi

# 5. 启动后端服务
echo "🚀 正在启动后端服务..."
echo "   命令: go run main.go api --port 8002"
echo "   目录: $(pwd)"
echo "----------------------------------------"

# 启动服务（前台运行，可以看到日志）
go run main.go api --port 8002

# 如果需要后台运行，可以使用下面的命令替换上面的命令：
# nohup go run main.go api --port 8002 > backend.log 2>&1 &
# echo "✅ 后端服务已在后台启动，日志保存到 backend.log"
# echo "📋 查看日志: tail -f backend.log"
