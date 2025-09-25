#!/bin/bash

# 大疆无人机监控系统 - 后台运行脚本

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
LOG_FILE="backend.log"
PID_FILE="backend.pid"

# 显示帮助信息
show_help() {
    echo "🚀 大疆无人机监控系统 - 后台运行脚本"
    echo ""
    echo "使用方法:"
    echo "  ./run-background.sh [选项] [动作]"
    echo ""
    echo "选项:"
    echo "  -h, --host HOST     指定主机地址 (默认: 0.0.0.0)"
    echo "  -p, --port PORT     指定端口 (默认: 18080)"
    echo "  -l, --log-file      指定日志文件 (默认: backend.log)"
    echo "  --help              显示此帮助信息"
    echo ""
    echo "动作:"
    echo "  start               启动服务"
    echo "  stop                停止服务"
    echo "  restart             重启服务"
    echo "  status              查看状态"
    echo "  logs                查看日志"
    echo ""
    echo "示例:"
    echo "  ./run-background.sh start              # 启动服务"
    echo "  ./run-background.sh stop               # 停止服务"
    echo "  ./run-background.sh restart            # 重启服务"
    echo "  ./run-background.sh status             # 查看状态"
    echo "  ./run-background.sh logs               # 查看日志"
    echo ""
}

# 解析命令行参数
ACTION=""
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
        -l|--log-file)
            LOG_FILE="$2"
            shift 2
            ;;
        --help)
            show_help
            exit 0
            ;;
        start|stop|restart|status|logs)
            ACTION="$1"
            shift
            ;;
        *)
            echo -e "${RED}❌ 未知参数: $1${NC}"
            echo "使用 --help 查看帮助信息"
            exit 1
            ;;
    esac
done

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
source .venv/bin/activate

# 检查依赖是否已安装
if ! python -c "import fastapi, uvicorn, redis, paho.mqtt.client, websockets" 2>/dev/null; then
    echo -e "${YELLOW}⚠️  依赖未安装，正在安装...${NC}"
    pip install fastapi==0.65.0 uvicorn==0.13.4 redis==3.5.3 paho-mqtt==1.5.1 websockets==8.1
    echo -e "${GREEN}✅ 依赖安装完成${NC}"
fi

# 启动服务
start_service() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo -e "${YELLOW}⚠️  服务已在运行 (PID: $(cat $PID_FILE))${NC}"
        return 1
    fi
    
    echo -e "${BLUE}🚀 启动后端服务...${NC}"
    echo -e "${BLUE}📱 访问地址: http://$HOST:$PORT${NC}"
    echo -e "${BLUE}📚 API 文档: http://$HOST:$PORT/docs${NC}"
    echo -e "${BLUE}📝 日志文件: $LOG_FILE${NC}"
    
    # 后台启动服务
    nohup python -m uvicorn backend.app.main:app \
        --host "$HOST" \
        --port "$PORT" \
        --log-level info \
        > "$LOG_FILE" 2>&1 &
    
    # 保存 PID
    echo $! > "$PID_FILE"
    
    # 等待服务启动
    sleep 3
    
    if kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo -e "${GREEN}✅ 服务启动成功 (PID: $(cat $PID_FILE))${NC}"
    else
        echo -e "${RED}❌ 服务启动失败${NC}"
        rm -f "$PID_FILE"
        return 1
    fi
}

# 停止服务
stop_service() {
    if [ ! -f "$PID_FILE" ]; then
        echo -e "${YELLOW}⚠️  服务未运行${NC}"
        return 1
    fi
    
    PID=$(cat "$PID_FILE")
    if ! kill -0 "$PID" 2>/dev/null; then
        echo -e "${YELLOW}⚠️  服务未运行${NC}"
        rm -f "$PID_FILE"
        return 1
    fi
    
    echo -e "${BLUE}⏹️  停止服务 (PID: $PID)...${NC}"
    kill "$PID"
    
    # 等待进程结束
    for i in {1..10}; do
        if ! kill -0 "$PID" 2>/dev/null; then
            break
        fi
        sleep 1
    done
    
    # 强制杀死进程
    if kill -0 "$PID" 2>/dev/null; then
        echo -e "${YELLOW}⚠️  强制停止服务...${NC}"
        kill -9 "$PID"
    fi
    
    rm -f "$PID_FILE"
    echo -e "${GREEN}✅ 服务已停止${NC}"
}

# 查看状态
show_status() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        PID=$(cat "$PID_FILE")
        echo -e "${GREEN}✅ 服务正在运行${NC}"
        echo "  PID: $PID"
        echo "  主机: $HOST"
        echo "  端口: $PORT"
        echo "  日志: $LOG_FILE"
        echo "  访问: http://$HOST:$PORT"
    else
        echo -e "${RED}❌ 服务未运行${NC}"
        rm -f "$PID_FILE"
    fi
}

# 查看日志
show_logs() {
    if [ -f "$LOG_FILE" ]; then
        echo -e "${BLUE}📋 显示日志 (按 Ctrl+C 退出):${NC}"
        tail -f "$LOG_FILE"
    else
        echo -e "${YELLOW}⚠️  日志文件不存在: $LOG_FILE${NC}"
    fi
}

# 执行动作
case "$ACTION" in
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        stop_service
        sleep 2
        start_service
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    "")
        echo -e "${RED}❌ 请指定动作${NC}"
        echo "使用 --help 查看帮助信息"
        exit 1
        ;;
    *)
        echo -e "${RED}❌ 未知动作: $ACTION${NC}"
        echo "使用 --help 查看帮助信息"
        exit 1
        ;;
esac
