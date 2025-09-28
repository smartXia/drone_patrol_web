#!/bin/bash

# 无人机监控系统 Go后端启动脚本

echo "启动无人机监控系统 Go后端..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: Go未安装，请先安装Go 1.21+"
    exit 1
fi

# 检查Go版本
GO_VERSION=$(go version | cut -d' ' -f3 | cut -d'o' -f2)
REQUIRED_VERSION="1.21"
if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "错误: Go版本过低，需要1.21+，当前版本: $GO_VERSION"
    exit 1
fi

# 设置环境变量
export PORT=${PORT:-18080}
export DATABASE_PATH=${DATABASE_PATH:-./data/backend.db}
export ENV=${ENV:-development}

echo "环境配置:"
echo "  端口: $PORT"
echo "  数据库: $DATABASE_PATH"
echo "  环境: $ENV"

# 创建数据目录
mkdir -p data

# 下载依赖
echo "下载依赖..."
go mod download

# 构建应用
echo "构建应用..."
go build -o drone-patrol-backend main.go

# 启动应用
echo "启动应用..."
./drone-patrol-backend
