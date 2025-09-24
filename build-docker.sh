#!/bin/bash

# 大疆无人机监控系统 - Docker 构建脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🐳 大疆无人机监控系统 - Docker 构建脚本${NC}"
echo ""

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到 Docker，请先安装 Docker${NC}"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到 Docker Compose，请先安装 Docker Compose${NC}"
    exit 1
fi

# 检查必要文件
if [ ! -f "Dockerfile" ]; then
    echo -e "${RED}❌ 错误: Dockerfile 不存在${NC}"
    exit 1
fi

if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}❌ 错误: docker-compose.yml 不存在${NC}"
    exit 1
fi

# 显示选项
echo -e "${BLUE}请选择操作:${NC}"
echo "1. 构建镜像"
echo "2. 启动服务"
echo "3. 停止服务"
echo "4. 重启服务"
echo "5. 查看日志"
echo "6. 清理资源"
echo "7. 完整部署"
echo ""

read -p "请输入选项 (1-7): " choice

case $choice in
    1)
        echo -e "${BLUE}🏗️  构建 Docker 镜像...${NC}"
        docker build -t drone-patrol:latest .
        echo -e "${GREEN}✅ 镜像构建完成${NC}"
        ;;
    2)
        echo -e "${BLUE}🚀 启动服务...${NC}"
        docker-compose up -d
        echo -e "${GREEN}✅ 服务启动完成${NC}"
        echo -e "${BLUE}📱 访问地址: http://localhost:18080${NC}"
        ;;
    3)
        echo -e "${BLUE}⏹️  停止服务...${NC}"
        docker-compose down
        echo -e "${GREEN}✅ 服务已停止${NC}"
        ;;
    4)
        echo -e "${BLUE}🔄 重启服务...${NC}"
        docker-compose restart
        echo -e "${GREEN}✅ 服务已重启${NC}"
        ;;
    5)
        echo -e "${BLUE}📋 查看日志...${NC}"
        docker-compose logs -f
        ;;
    6)
        echo -e "${BLUE}🧹 清理资源...${NC}"
        docker-compose down -v
        docker system prune -f
        echo -e "${GREEN}✅ 清理完成${NC}"
        ;;
    7)
        echo -e "${BLUE}🚀 完整部署...${NC}"
        echo "1. 构建镜像..."
        docker build -t drone-patrol:latest .
        echo "2. 启动服务..."
        docker-compose up -d
        echo -e "${GREEN}✅ 部署完成！${NC}"
        echo -e "${BLUE}📱 访问地址: http://localhost:18080${NC}"
        echo -e "${BLUE}📚 API 文档: http://localhost:18080/docs${NC}"
        ;;
    *)
        echo -e "${RED}❌ 无效选项${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}🎉 操作完成！${NC}"
