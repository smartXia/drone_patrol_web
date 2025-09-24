#!/bin/bash

# 大疆无人机监控系统 - 后端启动脚本
# 使用方法: ./run.sh [参数]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认参数
HOST="0.0.0.0"
PORT="18080"
RELOAD=""
LOG_LEVEL="info"

# 显示帮助信息
show_help() {
    echo "🚀 大疆无人机监控系统 - 后端启动脚本"
    echo ""
    echo "使用方法:"
    echo "  ./run.sh [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --host HOST     指定主机地址 (默认: 0.0.0.0)"
    echo "  -p, --port PORT     指定端口 (默认: 18080)"
    echo "  -r, --reload        启用热重载 (开发模式)"
    echo "  -l, --log-level     日志级别 (默认: info)"
    echo "  --help              显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  ./run.sh                    # 使用默认参数启动"
    echo "  ./run.sh -r                 # 开发模式启动"
    echo "  ./run.sh -h 127.0.0.1 -p 8080  # 指定主机和端口"
    echo ""
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--host)
            HOST="$2"
            shift 2
            ;;
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        -r|--reload)
            RELOAD="--reload"
            shift
            ;;
        -l|--log-level)
            LOG_LEVEL="$2"
            shift 2
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}❌ 未知参数: $1${NC}"
            echo "使用 --help 查看帮助信息"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🚀 启动大疆无人机监控系统后端...${NC}"
echo ""

# 检查是否在项目根目录
if [ ! -f "requirements.txt" ] || [ ! -d "backend" ]; then
    echo -e "${RED}❌ 错误: 请在项目根目录下运行此脚本${NC}"
    exit 1
fi

# 检查 Python 是否安装
if ! command -v python3 &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到 Python3，请先安装 Python 3.6+${NC}"
    exit 1
fi

# 检查虚拟环境
if [ ! -d ".venv" ]; then
    echo -e "${YELLOW}⚠️  虚拟环境不存在，正在创建...${NC}"
    python3 -m venv .venv
    echo -e "${GREEN}✅ 虚拟环境创建成功${NC}"
fi

# 激活虚拟环境
echo -e "${BLUE}📦 激活虚拟环境...${NC}"
source .venv/bin/activate

# 检查依赖是否已安装
echo -e "${BLUE}📦 检查 Python 依赖...${NC}"
if ! python -c "import fastapi, uvicorn, redis, paho.mqtt.client, websockets" 2>/dev/null; then
    echo -e "${YELLOW}⚠️  依赖未安装，正在安装...${NC}"
    pip install fastapi==0.65.0 uvicorn==0.13.4 redis==3.5.3 paho-mqtt==1.5.1 websockets==8.1
    echo -e "${GREEN}✅ 依赖安装完成${NC}"
else
    echo -e "${GREEN}✅ 依赖检查通过${NC}"
fi

# 检查 Redis 是否运行
echo -e "${BLUE}🔍 检查 Redis 服务...${NC}"
if ! pgrep -x "redis-server" > /dev/null; then
    echo -e "${YELLOW}⚠️  Redis 服务未运行，请先启动 Redis${NC}"
    echo "启动命令: sudo systemctl start redis"
    echo "或者: redis-server"
    echo ""
    read -p "是否继续启动后端服务？(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "已取消启动"
        exit 0
    fi
fi

# 显示启动信息
echo ""
echo -e "${GREEN}✅ 环境检查完成${NC}"
echo -e "${BLUE}📋 启动参数:${NC}"
echo "  主机: $HOST"
echo "  端口: $PORT"
echo "  热重载: ${RELOAD:-否}"
echo "  日志级别: $LOG_LEVEL"
echo ""

# 启动后端服务
echo -e "${GREEN}🚀 启动后端服务...${NC}"
echo -e "${BLUE}📱 访问地址: http://$HOST:$PORT${NC}"
echo -e "${BLUE}📚 API 文档: http://$HOST:$PORT/docs${NC}"
echo ""
echo -e "${YELLOW}按 Ctrl+C 停止服务${NC}"
echo ""

# 启动 uvicorn
python -m uvicorn backend.app.main:app \
    --host "$HOST" \
    --port "$PORT" \
    --log-level "$LOG_LEVEL" \
    $RELOAD
