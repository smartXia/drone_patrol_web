import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'

// 设置中文语言包
dayjs.locale('zh-cn')

/**
 * 格式化时间戳
 * @param {string|Date} timestamp 时间戳
 * @param {string} format 格式字符串
 * @returns {string} 格式化后的时间
 */
export function formatTime(timestamp, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!timestamp) return '未知'
  
  try {
    return dayjs(timestamp).format(format)
  } catch (error) {
    console.error('时间格式化错误:', error)
    return '格式错误'
  }
}

/**
 * 格式化相对时间
 * @param {string|Date} timestamp 时间戳
 * @returns {string} 相对时间字符串
 */
export function formatRelativeTime(timestamp) {
  if (!timestamp) return '未知'
  
  try {
    const now = dayjs()
    const target = dayjs(timestamp)
    const diff = now.diff(target, 'minute')
    
    if (diff < 1) return '刚刚'
    if (diff < 60) return `${diff}分钟前`
    
    const diffHours = now.diff(target, 'hour')
    if (diffHours < 24) return `${diffHours}小时前`
    
    const diffDays = now.diff(target, 'day')
    if (diffDays < 30) return `${diffDays}天前`
    
    const diffMonths = now.diff(target, 'month')
    if (diffMonths < 12) return `${diffMonths}个月前`
    
    const diffYears = now.diff(target, 'year')
    return `${diffYears}年前`
  } catch (error) {
    console.error('相对时间格式化错误:', error)
    return '时间错误'
  }
}

/**
 * 获取当前时间戳
 * @returns {string} ISO格式时间戳
 */
export function getCurrentTimestamp() {
  return dayjs().toISOString()
}

/**
 * 检查时间是否过期
 * @param {string|Date} timestamp 时间戳
 * @param {number} minutes 过期分钟数
 * @returns {boolean} 是否过期
 */
export function isExpired(timestamp, minutes = 5) {
  if (!timestamp) return true
  
  try {
    const target = dayjs(timestamp)
    const now = dayjs()
    const diff = now.diff(target, 'minute')
    return diff > minutes
  } catch (error) {
    console.error('时间过期检查错误:', error)
    return true
  }
}

/**
 * 格式化持续时间
 * @param {number} seconds 秒数
 * @returns {string} 格式化的持续时间
 */
export function formatDuration(seconds) {
  if (!seconds || seconds < 0) return '0秒'
  
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  
  if (hours > 0) {
    return `${hours}小时${minutes}分钟${secs}秒`
  } else if (minutes > 0) {
    return `${minutes}分钟${secs}秒`
  } else {
    return `${secs}秒`
  }
}

/**
 * 获取时间范围
 * @param {string} range 时间范围类型
 * @returns {Object} 开始和结束时间
 */
export function getTimeRange(range = 'today') {
  const now = dayjs()
  
  switch (range) {
    case 'today':
      return {
        start: now.startOf('day').toISOString(),
        end: now.endOf('day').toISOString()
      }
    case 'yesterday':
      const yesterday = now.subtract(1, 'day')
      return {
        start: yesterday.startOf('day').toISOString(),
        end: yesterday.endOf('day').toISOString()
      }
    case 'week':
      return {
        start: now.startOf('week').toISOString(),
        end: now.endOf('week').toISOString()
      }
    case 'month':
      return {
        start: now.startOf('month').toISOString(),
        end: now.endOf('month').toISOString()
      }
    case 'last7days':
      return {
        start: now.subtract(7, 'day').toISOString(),
        end: now.toISOString()
      }
    case 'last30days':
      return {
        start: now.subtract(30, 'day').toISOString(),
        end: now.toISOString()
      }
    default:
      return {
        start: now.subtract(1, 'day').toISOString(),
        end: now.toISOString()
      }
  }
}
