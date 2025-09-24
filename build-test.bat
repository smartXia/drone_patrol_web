@echo off
chcp 65001 >nul
echo 🚀 测试版打包脚本...

REM 检查 Python
python --version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: 未找到 Python
    pause
    exit /b 1
)

REM 检查 Node.js
node --version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: 未找到 Node.js
    pause
    exit /b 1
)

REM 构建前端
echo 📦 构建前端...
npm run build
if errorlevel 1 (
    echo ❌ 前端构建失败
    pause
    exit /b 1
)

REM 安装 PyInstaller
echo 📦 安装 PyInstaller...
pip install pyinstaller

REM 测试运行
echo 🧪 测试运行...
python test-exe.py
if errorlevel 1 (
    echo ❌ 测试运行失败
    pause
    exit /b 1
)

REM 打包
echo 🏗️ 开始打包...
pyinstaller --onefile --console --name "大疆无人机监控系统" test-exe.py

echo.
echo ✅ 打包完成！
echo 📁 exe 文件位置: dist\大疆无人机监控系统.exe
echo 🚀 运行方法: 双击 exe 文件
echo.
pause
