#!/bin/bash
# 无人机监控系统后端构建和启动脚本

echo "=== 无人机监控系统后端构建和启动脚本 ==="

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# 进入后端目录
cd "$(dirname "$0")"

# 创建日志目录
mkdir -p ./logs

# 构建项目
echo "正在构建后端服务..."
go build -o drone-patrol-backend .

if [ $? -ne 0 ]; then
    echo "Error: 构建失败"
    exit 1
fi

echo "构建成功！"

# 设置可执行权限
chmod 777 drone-patrol-backend

# 启动服务
echo "正在启动后端服务..."
nohup ./drone-patrol-backend > ./logs/backend.log 2>&1 &

# 获取进程ID
PID=$!

# 等待一秒检查进程是否启动成功
sleep 1

if ps -p $PID > /dev/null; then
    echo "后端服务启动成功！"
    echo "进程ID: $PID"
    echo "日志文件: ./logs/backend.log"
    echo "访问地址: http://localhost:18080"
    echo ""
    echo "查看日志: tail -f ./logs/backend.log"
    echo "停止服务: kill $PID"
else
    echo "Error: 后端服务启动失败"
    echo "请查看日志文件: ./logs/backend.log"
    exit 1
fi
