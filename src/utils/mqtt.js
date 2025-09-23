/**
 * MQTT工具函数
 */

// MQTT主题常量
export const MQTT_TOPICS = {
  // 设备属性主题
  DEVICE_PROPERTIES: 'device/+/properties',
  // 设备状态主题
  DEVICE_STATUS: 'device/+/status',
  // 设备遥测数据主题
  DEVICE_TELEMETRY: 'device/+/telemetry',
  // 设备事件主题
  DEVICE_EVENTS: 'device/+/events',
  // 设备命令主题
  DEVICE_COMMANDS: 'device/+/commands',
  // 设备响应主题
  DEVICE_RESPONSES: 'device/+/responses'
}

// 设备状态枚举
export const DEVICE_STATUS = {
  ONLINE: 'online',
  OFFLINE: 'offline',
  ERROR: 'error',
  WARNING: 'warning',
  MAINTENANCE: 'maintenance'
}

// 设备类型枚举
export const DEVICE_TYPES = {
  DRONE: 'drone',
  DOCK: 'dock',
  CAMERA: 'camera',
  SENSOR: 'sensor'
}

/**
 * 解析MQTT主题获取设备ID
 * @param {string} topic MQTT主题
 * @returns {string} 设备ID
 */
export function parseDeviceIdFromTopic(topic) {
  const parts = topic.split('/')
  return parts[1] || 'unknown'
}

/**
 * 解析MQTT主题获取消息类型
 * @param {string} topic MQTT主题
 * @returns {string} 消息类型
 */
export function parseMessageTypeFromTopic(topic) {
  const parts = topic.split('/')
  return parts[2] || 'unknown'
}

/**
 * 验证MQTT消息格式
 * @param {Object} message MQTT消息对象
 * @returns {boolean} 是否有效
 */
export function validateMqttMessage(message) {
  if (!message || typeof message !== 'object') {
    return false
  }
  
  // 检查必要字段
  const requiredFields = ['timestamp', 'deviceId', 'data']
  return requiredFields.every(field => message.hasOwnProperty(field))
}

/**
 * 格式化设备状态显示文本
 * @param {Object} device 设备数据
 * @returns {string} 状态文本
 */
export function formatDeviceStatus(device) {
  if (!device) return '未知'
  
  if (device.status) {
    const statusMap = {
      [DEVICE_STATUS.ONLINE]: '在线',
      [DEVICE_STATUS.OFFLINE]: '离线',
      [DEVICE_STATUS.ERROR]: '错误',
      [DEVICE_STATUS.WARNING]: '警告',
      [DEVICE_STATUS.MAINTENANCE]: '维护中'
    }
    return statusMap[device.status] || device.status
  }
  
  if (device.connected) return '在线'
  if (device.error) return '错误'
  
  return '正常'
}

/**
 * 获取设备状态类型（用于UI显示）
 * @param {Object} device 设备数据
 * @returns {string} 状态类型
 */
export function getDeviceStatusType(device) {
  if (!device) return 'info'
  
  if (device.status === DEVICE_STATUS.ONLINE || device.connected) return 'success'
  if (device.status === DEVICE_STATUS.OFFLINE) return 'info'
  if (device.status === DEVICE_STATUS.ERROR || device.error) return 'danger'
  if (device.status === DEVICE_STATUS.WARNING) return 'warning'
  if (device.status === DEVICE_STATUS.MAINTENANCE) return 'warning'
  
  return 'info'
}

/**
 * 生成MQTT客户端ID
 * @param {string} prefix 前缀
 * @returns {string} 客户端ID
 */
export function generateClientId(prefix = 'drone_monitor') {
  const timestamp = Date.now()
  const random = Math.random().toString(36).substring(2, 8)
  return `${prefix}_${timestamp}_${random}`
}

/**
 * 构建MQTT连接URL
 * @param {Object} config MQTT配置
 * @returns {string} 连接URL
 */
export function buildMqttUrl(config) {
  const { protocol, host, port, path } = config
  
  // 在浏览器环境中，TCP 协议需要转换为 WebSocket
  if (protocol === 'tcp') {
    // TCP 协议在浏览器中通过 WebSocket 实现，使用 1883 端口和 /mqtt 路径
    return `ws://${host}:1883/mqtt`
  }
  
  // WebSocket 协议需要路径
  const suffix = path && typeof path === 'string' && path.length > 0 ? (path.startsWith('/') ? path : `/${path}`) : '/mqtt'
  return `${protocol}://${host}:${port}${suffix}`
}

/**
 * 检查设备是否在线
 * @param {Object} device 设备数据
 * @returns {boolean} 是否在线
 */
export function isDeviceOnline(device) {
  if (!device) return false
  
  if (device.status === DEVICE_STATUS.ONLINE) return true
  if (device.connected === true) return true
  
  // 检查最后更新时间，如果超过5分钟则认为离线
  if (device.lastUpdate) {
    const lastUpdate = new Date(device.lastUpdate)
    const now = new Date()
    const diffMinutes = (now - lastUpdate) / (1000 * 60)
    return diffMinutes < 5
  }
  
  return false
}

/**
 * 获取设备类型图标
 * @param {string} deviceType 设备类型
 * @returns {string} 图标类名
 */
export function getDeviceTypeIcon(deviceType) {
  const iconMap = {
    [DEVICE_TYPES.DRONE]: 'el-icon-plane',
    [DEVICE_TYPES.DOCK]: 'el-icon-house',
    [DEVICE_TYPES.CAMERA]: 'el-icon-camera',
    [DEVICE_TYPES.SENSOR]: 'el-icon-monitor'
  }
  return iconMap[deviceType] || 'el-icon-cpu'
}
