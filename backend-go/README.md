# 无人机监控系统 - Go后端

这是无人机监控系统的Go语言后端实现，提供RESTful API和WebSocket支持。

## 功能特性

- 🔧 **设备管理** - 设备的增删改查、状态管理
- 📡 **MQTT配置** - MQTT连接配置管理
- 🔴 **Redis代理** - Redis数据库操作代理
- 📊 **错误码查询** - 大疆错误码查询服务
- 🌐 **WebSocket** - 实时数据推送
- 🐳 **Docker支持** - 容器化部署

## API端点

### 健康检查和工具
- `GET /api/health` - 健康检查
- `GET /api/error-codes` - 获取错误码列表
- `GET /api/network/ping` - 网络ping测试

### 设备管理
- `GET /api/devices` - 获取设备列表
- `GET /api/devices/current` - 获取当前设备信息
- `POST /api/devices` - 创建设备
- `PUT /api/devices/{device_id}` - 更新设备
- `DELETE /api/devices/{device_id}` - 删除设备
- `POST /api/devices/{device_id}/set-current` - 设置当前设备
- `POST /api/devices/{device_id}/set-gateway` - 设置网关设备
- `DELETE /api/devices/clear` - 清空所有设备
- `DELETE /api/devices/remove-defaults` - 删除默认设备

### MQTT配置管理
- `GET /api/mqtt/profiles` - 获取MQTT配置列表
- `POST /api/mqtt/profiles` - 创建MQTT配置
- `GET /api/mqtt/profiles/{pid}` - 获取单个MQTT配置
- `PUT /api/mqtt/profiles/{pid}` - 更新MQTT配置
- `DELETE /api/mqtt/profiles/{pid}` - 删除MQTT配置
- `POST /api/mqtt/profiles/{pid}/default` - 设置默认MQTT配置
- `POST /api/mqtt/test` - 测试MQTT连接

### Redis代理
- `POST /api/redis/connect/test` - 测试Redis连接
- `POST /scan` - 扫描Redis键
- `POST /type` - 获取键类型
- `POST /ttl` - 获取键TTL
- `POST /expire` - 设置键过期
- `POST /persist` - 持久化键
- `POST /rename` - 重命名键
- `POST /del` - 删除键
- `POST /get` - 获取字符串值
- `POST /set` - 设置字符串值
- `POST /hash/getall` - 获取哈希所有字段
- `POST /hash/set` - 设置哈希字段
- `POST /hash/del` - 删除哈希字段
- `POST /list/range` - 获取列表范围
- `POST /list/lpush` - 左推入列表
- `POST /list/rpush` - 右推入列表
- `POST /list/lpop` - 左弹出列表
- `POST /list/rpop` - 右弹出列表
- `POST /list/set` - 设置列表元素
- `POST /list/lrem` - 删除列表元素
- `POST /set/scan` - 扫描集合
- `POST /set/sadd` - 添加集合成员
- `POST /set/srem` - 删除集合成员
- `POST /zset/range` - 获取有序集合范围
- `POST /zset/zadd` - 添加有序集合成员
- `POST /zset/zrem` - 删除有序集合成员
- `POST /zset/zincrby` - 增加有序集合分数

## 快速开始

### 使用Docker

1. 构建并运行：
```bash
docker-compose up --build
```

2. 访问API：
```bash
curl http://localhost:18080/api/health
```

### 本地开发

1. 安装Go 1.21+
2. 安装依赖：
```bash
go mod download
```

3. 运行：
```bash
go run main.go
```

## 环境变量

- `PORT` - 服务端口 (默认: 18080)
- `DATABASE_PATH` - 数据库文件路径 (默认: ./data/backend.db)
- `ENV` - 环境 (development/production)

## 项目结构

```
backend-go/
├── main.go                 # 主入口文件
├── go.mod                  # Go模块文件
├── go.sum                  # 依赖校验文件
├── Dockerfile              # Docker构建文件
├── docker-compose.yml      # Docker Compose配置
├── README.md               # 项目说明
└── internal/               # 内部包
    ├── config/             # 配置管理
    ├── database/           # 数据库相关
    ├── handlers/           # HTTP处理器
    ├── middleware/         # 中间件
    ├── models/             # 数据模型
    └── services/           # 业务逻辑服务
```

## 技术栈

- **Go 1.21** - 编程语言
- **Gin** - Web框架
- **SQLite** - 数据库
- **Redis** - 缓存和消息队列
- **MQTT** - 物联网通信协议
- **Docker** - 容器化

## 开发说明

### 添加新的API端点

1. 在 `internal/models/` 中定义数据模型
2. 在 `internal/services/` 中实现业务逻辑
3. 在 `internal/handlers/` 中实现HTTP处理器
4. 在 `internal/handlers/routes.go` 中注册路由

### 数据库迁移

数据库表会在应用启动时自动创建，无需手动迁移。

## 许可证

MIT License
gox -os "linux" -arch amd64