# 无人机设备监控系统

基于Vue 3 + Element Plus + MQTT.js的大疆无人机设备实时状态监控系统。

## 功能特性

- 🔌 **MQTT实时连接**: 支持连接MQTT broker，实时接收设备状态数据
- 📊 **设备状态监控**: 实时显示无人机设备状态、属性、遥测数据
- 📈 **数据可视化**: 支持原始数据和格式化数据两种显示模式
- 📝 **消息历史**: 记录和查看历史消息数据
- 🔄 **自动重连**: 支持网络断开自动重连机制
- 🎨 **现代化UI**: 基于Element Plus的现代化用户界面
- 📱 **响应式设计**: 支持多种屏幕尺寸

## 技术栈

- **前端框架**: Vue 3 (Composition API)
- **状态管理**: Pinia
- **UI组件库**: Element Plus
- **MQTT客户端**: MQTT.js
- **构建工具**: Vite
- **路由**: Vue Router 4
- **时间处理**: Day.js

## 快速开始

### 环境要求

- Node.js >= 16.0.0
- npm >= 8.0.0

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 项目结构

```
drone-patrol-web/
├── src/
│   ├── components/          # 组件目录
│   │   └── DeviceStatusCard.vue
│   ├── stores/             # 状态管理
│   │   └── mqtt.js
│   ├── utils/              # 工具函数
│   │   ├── mqtt.js
│   │   └── time.js
│   ├── views/              # 页面组件
│   │   └── Dashboard.vue
│   ├── router/             # 路由配置
│   │   └── index.js
│   ├── App.vue             # 根组件
│   └── main.js             # 入口文件
├── index.html              # HTML模板
├── package.json            # 项目配置
├── vite.config.js          # Vite配置
└── README.md               # 项目说明
```

## MQTT配置

系统默认配置连接到以下MQTT broker：

- **Broker**: `tcp://180.102.24.183:1883`
- **用户名**: `admin`
- **密码**: `njbc@2216446`

### 订阅主题

系统会自动订阅以下主题：

- `device/+/properties` - 设备属性
- `device/+/status` - 设备状态
- `device/+/telemetry` - 设备遥测数据
- `device/+/events` - 设备事件

## 使用说明

### 1. 连接MQTT

点击"连接MQTT"按钮建立与MQTT broker的连接。

### 2. 查看设备状态

连接成功后，系统会自动接收并显示设备状态数据。

### 3. 查看消息历史

在页面底部可以查看最近的消息历史记录。

### 4. 设备操作

- **刷新**: 刷新设备状态数据
- **历史记录**: 查看设备的历史消息
- **格式化/原始**: 切换数据显示模式

## 开发指南

### 添加新的设备类型

1. 在 `src/utils/mqtt.js` 中的 `DEVICE_TYPES` 枚举中添加新类型
2. 在 `getDeviceTypeIcon` 函数中添加对应的图标

### 自定义MQTT主题

修改 `src/stores/mqtt.js` 中的 `subscribeToTopics` 函数来订阅自定义主题。

### 扩展设备状态显示

在 `src/utils/mqtt.js` 中添加新的状态处理函数。

## 部署说明

### 构建生产版本

```bash
npm run build
```

构建完成后，`dist` 目录包含可部署的静态文件。

### 部署到Web服务器

将 `dist` 目录中的文件部署到任何支持静态文件的Web服务器。

### 环境变量配置

可以通过环境变量配置MQTT连接参数：

```bash
VITE_MQTT_HOST=your-mqtt-host
VITE_MQTT_PORT=1883
VITE_MQTT_USERNAME=your-username
VITE_MQTT_PASSWORD=your-password
```

## 故障排除

### 连接问题

1. 检查MQTT broker是否可访问
2. 验证用户名和密码是否正确
3. 检查网络连接

### 数据不显示

1. 确认设备正在发送数据
2. 检查订阅的主题是否正确
3. 查看浏览器控制台是否有错误信息

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue 或联系开发团队。
