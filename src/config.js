// 应用配置文件
export const config = {
  // 后端 API 配置
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:18080',
  wsBaseUrl: import.meta.env.VITE_WS_BASE_URL || 'ws://127.0.0.1:18080',
  
  // 开发环境配置
  devMode: import.meta.env.VITE_DEV_MODE === 'true' || import.meta.env.DEV,
  
  // 超时配置
  timeout: 15000,
  
  // 其他配置
  appName: '大疆无人机监控系统',
  version: '1.0.0'
}

// 导出常用的配置
export const API_BASE_URL = config.apiBaseUrl
export const WS_BASE_URL = config.wsBaseUrl
export const TIMEOUT = config.timeout
