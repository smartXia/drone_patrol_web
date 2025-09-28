@echo off
REM 无人机监控系统 Go后端启动脚本

echo 启动无人机监控系统 Go后端...

REM 检查Go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: Go未安装，请先安装Go 1.21+
    pause
    exit /b 1
)

REM 设置环境变量
if not defined PORT set PORT=18080
if not defined DATABASE_PATH set DATABASE_PATH=./data/backend.db
if not defined ENV set ENV=development

echo 环境配置:
echo   端口: %PORT%
echo   数据库: %DATABASE_PATH%
echo   环境: %ENV%

REM 创建数据目录
if not exist data mkdir data

REM 下载依赖
echo 下载依赖...
go mod download

REM 构建应用
echo 构建应用...
go build -o drone-patrol-backend.exe main.go

REM 启动应用
echo 启动应用...
drone-patrol-backend.exe
