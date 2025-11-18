# 高德地图API配置说明

## 📍 API密钥配置

### 1. 获取API密钥
- 访问 [高德开放平台](https://lbs.amap.com/)
- 注册并创建应用
- 获取API密钥和安全密钥

### 2. 当前配置
```javascript
// 在 src/config/amap.js 中配置
export const AMAP_CONFIG = {
  key: 'd593eb16f78b90ac373b5770434f4b89',           // 您的API密钥
  securityKey: '81c6f50547753df15a0b661c5302deb5',   // 您的安全密钥
  version: '2.0',                                     // API版本
}
```

### 3. 环境变量配置（可选）
```bash
# 在 .env.local 文件中配置
VITE_AMAP_KEY=your_api_key_here
VITE_AMAP_SECURITY_KEY=your_security_key_here
```

## 🗺️ 功能特性

### 1. 地图功能
- ✅ 实时飞机位置显示
- ✅ 飞行轨迹绘制
- ✅ 地图控件（比例尺、工具条）
- ✅ 多种地图样式
- ✅ 3D视图支持

### 2. 轨迹功能
- ✅ 实时轨迹记录
- ✅ 轨迹线显示/隐藏
- ✅ 轨迹清除
- ✅ 飞行距离计算
- ✅ 飞行时间统计

### 3. 地图样式
```javascript
// 支持的地图样式
mapStyles: {
  normal: 'amap://styles/normal',     // 标准样式
  dark: 'amap://styles/dark',         // 暗色样式
  light: 'amap://styles/light',       // 浅色样式
  fresh: 'amap://styles/fresh',       // 清新样式
  grey: 'amap://styles/grey',         // 灰色样式
  graffiti: 'amap://styles/graffiti', // 涂鸦样式
  macaron: 'amap://styles/macaron',   // 马卡龙样式
  blue: 'amap://styles/blue',         // 蓝色样式
  darkblue: 'amap://styles/darkblue', // 深蓝样式
  wine: 'amap://styles/wine'          // 酒红样式
}
```

## 🔧 使用方法

### 1. 基本使用
```vue
<template>
  <FlightTrajectoryMap
    v-model="showMap"
    :aircraft-data="aircraftData"
    :aircraft-sn="aircraftSn"
  />
</template>
```

### 2. 组件属性
- `modelValue`: 控制模态框显示/隐藏
- `aircraftData`: 飞机OSD数据对象
- `aircraftSn`: 飞机序列号

### 3. 数据格式
```javascript
// aircraftData 格式
{
  latitude: 31.971145,        // 纬度
  longitude: 118.654765,       // 经度
  height: 50.0,               // 高度
  attitude_head: 45.0,        // 航向角
  horizontal_speed: 10.5      // 水平速度
}
```

## 🚀 部署说明

### 1. 开发环境
- 确保网络连接正常
- 检查API密钥是否正确
- 查看浏览器控制台是否有错误

### 2. 生产环境
- 配置正确的API密钥
- 设置域名白名单
- 启用安全密钥验证

### 3. 故障排除
```javascript
// 检查API状态
import { getAMapStatus } from '@/config/amap'
console.log('API状态:', getAMapStatus())
```

## 📊 性能优化

### 1. 轨迹点限制
- 最多保留1000个轨迹点
- 自动清理旧轨迹点
- 避免内存泄漏

### 2. 地图控件
- 按需加载插件
- 延迟初始化
- 错误重试机制

### 3. 内存管理
- 组件卸载时清理地图实例
- 避免重复初始化
- 及时释放资源

## 🔒 安全说明

### 1. API密钥安全
- 不要在客户端代码中硬编码密钥
- 使用环境变量管理密钥
- 定期轮换密钥

### 2. 域名限制
- 在高德开放平台设置域名白名单
- 限制API调用频率
- 监控API使用情况

### 3. 安全密钥
- 启用安全密钥验证
- 配置IP白名单
- 设置调用限制

## 📞 技术支持

如果遇到问题，请检查：
1. 网络连接是否正常
2. API密钥是否正确
3. 域名是否在白名单中
4. 浏览器控制台错误信息

更多信息请参考 [高德地图API文档](https://lbs.amap.com/api/jsapi-v2/documentation)
