# 部署说明

## 环境配置

项目支持两种环境配置：

### 开发环境 (Development)
- 配置文件：`.env.development`
- 后端地址：`http://127.0.0.1:18080`
- WebSocket地址：`ws://127.0.0.1:18080`
- 启动命令：`npm run dev`

### 生产环境 (Production)
- 配置文件：`.env.production`
- 后端地址：`http://180.102.24.183:18080`
- WebSocket地址：`ws://180.102.24.183:18080`
- 构建命令：`npm run build:prod`

## 可用命令

### 开发环境
```bash
# 启动前端开发服务器（开发环境）
npm run dev

# 启动后端服务
npm run dev:py

# 同时启动前端和后端
npm run dev:full

# 构建开发环境版本
npm run build:dev

# 预览开发环境构建结果
npm run preview:dev
```

### 生产环境
```bash
# 构建生产环境版本
npm run build:prod

# 预览生产环境构建结果
npm run preview:prod

# 构建并启动生产环境
npm start
```

## 环境变量说明

| 变量名 | 说明 | 开发环境 | 生产环境 |
|--------|------|----------|----------|
| VITE_API_BASE_URL | 后端API地址 | http://127.0.0.1:18080 | http://180.102.24.183:18080 |
| VITE_WS_BASE_URL | WebSocket地址 | ws://127.0.0.1:18080 | ws://180.102.24.183:18080 |
| VITE_DEV_MODE | 开发模式 | true | false |
| VITE_ENV | 环境标识 | development | production |

## 部署步骤

### 1. 开发环境部署
```bash
# 1. 安装依赖
npm install

# 2. 创建Python虚拟环境
python -m venv .venv

# 3. 激活虚拟环境并安装Python依赖
.venv\Scripts\activate
pip install -r backend/requirements.txt

# 4. 启动开发环境
npm run dev:full
```

### 2. 生产环境部署
```bash
# 1. 安装依赖
npm install

# 2. 构建生产版本
npm run build:prod

# 3. 部署dist目录到Web服务器
# 将dist目录内容复制到Web服务器根目录
```

## 注意事项

1. 确保生产环境的服务器地址正确配置
2. 确保WebSocket连接地址与后端服务地址一致
3. 生产环境部署前请先测试连接是否正常
4. 建议使用HTTPS和WSS协议用于生产环境
