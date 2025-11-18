// 高德地图配置
export const AMAP_CONFIG = {
  // 高德地图API密钥（支持环境变量）
  key: import.meta.env.VITE_AMAP_KEY || 'd593eb16f78b90ac373b5770434f4b89',
  
  // 安全密钥（支持环境变量）
  securityKey: import.meta.env.VITE_AMAP_SECURITY_KEY || '81c6f50547753df15a0b661c5302deb5',
  
  // API版本
  version: '1.4.10',
  
  // 需要加载的插件
  plugins: [
    'AMap.GeometryUtil',  // 几何计算工具
    'AMap.Scale',         // 比例尺控件
    'AMap.ToolBar'        // 工具条控件
  ],
  
  // 默认地图配置
  defaultConfig: {
    zoom: 15,
    mapStyle: 'amap://styles/normal',
    center: [118.654765, 31.971145], // 默认中心点（南京）
    viewMode: '2D', // 2D视图（更稳定）
    pitch: 0, // 俯仰角
    rotation: 0 // 旋转角
  },
  
  // 地图样式选项
  mapStyles: {
    normal: 'amap://styles/normal',           // 标准样式
    dark: 'amap://styles/dark',               // 暗色样式
    light: 'amap://styles/light',             // 浅色样式
    fresh: 'amap://styles/fresh',             // 清新样式
    grey: 'amap://styles/grey',               // 灰色样式
    graffiti: 'amap://styles/graffiti',       // 涂鸦样式
    macaron: 'amap://styles/macaron',         // 马卡龙样式
    blue: 'amap://styles/blue',               // 蓝色样式
    darkblue: 'amap://styles/darkblue',       // 深蓝样式
    wine: 'amap://styles/wine'                // 酒红样式
  }
}

// 获取高德地图API URL
export const getAMapApiUrl = () => {
  const plugins = AMAP_CONFIG.plugins.join(',')
  return `https://webapi.amap.com/maps?v=${AMAP_CONFIG.version}&key=${AMAP_CONFIG.key}&plugin=${plugins}`
}

// 检查高德地图API是否已加载
export const isAMapLoaded = () => {
  const hasAMap = typeof window.AMap !== 'undefined'
  const hasWindowFlag = window.AMapLoaded === true
  console.log('🔍 API加载状态检查:', {
    hasAMap,
    hasWindowFlag,
    AMapType: typeof window.AMap,
    windowAMapLoaded: window.AMapLoaded
  })
  return hasAMap && hasWindowFlag
}

// 等待高德地图API加载完成
export const waitForAMap = (timeout = 10000) => {
  return new Promise((resolve, reject) => {
    // 立即检查
    if (typeof window.AMap !== 'undefined') {
      console.log('✅ AMap已可用，直接返回')
      resolve(window.AMap)
      return
    }
    
    console.log('⏳ 等待高德地图API加载...')
    const startTime = Date.now()
    const checkInterval = setInterval(() => {
      if (typeof window.AMap !== 'undefined') {
        clearInterval(checkInterval)
        console.log('✅ 高德地图API加载完成')
        resolve(window.AMap)
      } else if (Date.now() - startTime > timeout) {
        clearInterval(checkInterval)
        console.error('❌ 高德地图API加载超时')
        reject(new Error(`高德地图API加载超时 (${timeout}ms)，请检查网络连接或API密钥是否正确`))
      }
    }, 200)
  })
}

// 获取高德地图API状态信息
export const getAMapStatus = () => {
  return {
    loaded: isAMapLoaded(),
    hasAMap: typeof window.AMap !== 'undefined',
    hasWindowFlag: window.AMapLoaded === true,
    key: AMAP_CONFIG.key,
    version: AMAP_CONFIG.version
  }
}

// WGS84坐标系转换工具
export const COORDINATE_CONVERTER = {
  // WGS84转GCJ02（火星坐标系）
  wgs84ToGcj02: (lng, lat) => {
    // 严格检查输入参数
    if (lng === null || lng === undefined || lat === null || lat === undefined) {
      console.warn('⚠️ 坐标转换: 输入参数为null或undefined', { lng, lat })
      return [0, 0] // 返回默认坐标
    }
    
    if (typeof lng !== 'number' || typeof lat !== 'number') {
      console.warn('⚠️ 坐标转换: 输入参数不是数字类型', { lng, lat })
      return [0, 0] // 返回默认坐标
    }
    
    if (isNaN(lng) || isNaN(lat)) {
      console.warn('⚠️ 坐标转换: 输入参数为NaN', { lng, lat })
      return [0, 0] // 返回默认坐标
    }
    
    const a = 6378245.0 // 长半轴
    const ee = 0.00669342162296594323 // 偏心率平方
    
    if (isOutOfChina(lng, lat)) {
      return [lng, lat]
    }
    
    try {
      let dLat = transformLat(lng - 105.0, lat - 35.0)
      let dLng = transformLng(lng - 105.0, lat - 35.0)
      
      const radLat = lat / 180.0 * Math.PI
      let magic = Math.sin(radLat)
      magic = 1 - ee * magic * magic
      const sqrtMagic = Math.sqrt(magic)
      
      dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * Math.PI)
      dLng = (dLng * 180.0) / (a / sqrtMagic * Math.cos(radLat) * Math.PI)
      
      const result = [lng + dLng, lat + dLat]
      
      // 检查结果是否有效
      if (isNaN(result[0]) || isNaN(result[1])) {
        console.warn('⚠️ 坐标转换结果包含NaN', { input: [lng, lat], output: result })
        return [lng, lat] // 返回原始坐标
      }
      
      return result
    } catch (error) {
      console.error('❌ 坐标转换失败:', error, { lng, lat })
      return [lng, lat] // 返回原始坐标作为备用
    }
  },
  
  // GCJ02转WGS84
  gcj02ToWgs84: (lng, lat) => {
    // 严格检查输入参数
    if (lng === null || lng === undefined || lat === null || lat === undefined) {
      console.warn('⚠️ 坐标转换: 输入参数为null或undefined', { lng, lat })
      return [0, 0] // 返回默认坐标
    }
    
    if (typeof lng !== 'number' || typeof lat !== 'number') {
      console.warn('⚠️ 坐标转换: 输入参数不是数字类型', { lng, lat })
      return [0, 0] // 返回默认坐标
    }
    
    if (isNaN(lng) || isNaN(lat)) {
      console.warn('⚠️ 坐标转换: 输入参数为NaN', { lng, lat })
      return [0, 0] // 返回默认坐标
    }
    
    const a = 6378245.0
    const ee = 0.00669342162296594323
    
    if (isOutOfChina(lng, lat)) {
      return [lng, lat]
    }
    
    try {
      let dLat = transformLat(lng - 105.0, lat - 35.0)
      let dLng = transformLng(lng - 105.0, lat - 35.0)
      
      const radLat = lat / 180.0 * Math.PI
      let magic = Math.sin(radLat)
      magic = 1 - ee * magic * magic
      const sqrtMagic = Math.sqrt(magic)
      
      dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * Math.PI)
      dLng = (dLng * 180.0) / (a / sqrtMagic * Math.cos(radLat) * Math.PI)
      
      const result = [lng - dLng, lat - dLat]
      
      // 检查结果是否有效
      if (isNaN(result[0]) || isNaN(result[1])) {
        console.warn('⚠️ 坐标转换结果包含NaN', { input: [lng, lat], output: result })
        return [lng, lat] // 返回原始坐标
      }
      
      return result
    } catch (error) {
      console.error('❌ 坐标转换失败:', error, { lng, lat })
      return [lng, lat] // 返回原始坐标作为备用
    }
  },
  
  // 批量转换WGS84到GCJ02
  batchWgs84ToGcj02: (coordinates) => {
    if (!Array.isArray(coordinates)) {
      console.warn('⚠️ 批量坐标转换: 输入不是数组', coordinates)
      return []
    }
    
    return coordinates.map((coord, index) => {
      if (!Array.isArray(coord) || coord.length < 2) {
        console.warn(`⚠️ 批量坐标转换: 第${index}个坐标无效`, coord)
        return [0, 0]
      }
      
      const [lng, lat] = coord
      if (lng === null || lng === undefined || lat === null || lat === undefined) {
        console.warn(`⚠️ 批量坐标转换: 第${index}个坐标包含null值`, coord)
        return [0, 0]
      }
      
      return COORDINATE_CONVERTER.wgs84ToGcj02(lng, lat)
    })
  }
}

// 判断是否在中国境外
function isOutOfChina(lng, lat) {
  return (lng < 72.004 || lng > 137.8347) || (lat < 0.8293 || lat > 55.8271)
}

// 纬度转换
function transformLat(lng, lat) {
  let ret = -100.0 + 2.0 * lng + 3.0 * lat + 0.2 * lat * lat + 0.1 * lng * lat + 0.2 * Math.sqrt(Math.abs(lng))
  ret += (20.0 * Math.sin(6.0 * lng * Math.PI) + 20.0 * Math.sin(2.0 * lng * Math.PI)) * 2.0 / 3.0
  ret += (20.0 * Math.sin(lat * Math.PI) + 40.0 * Math.sin(lat / 3.0 * Math.PI)) * 2.0 / 3.0
  ret += (160.0 * Math.sin(lat / 12.0 * Math.PI) + 320 * Math.sin(lat * Math.PI / 30.0)) * 2.0 / 3.0
  return ret
}

// 经度转换
function transformLng(lng, lat) {
  let ret = 300.0 + lng + 2.0 * lat + 0.1 * lng * lng + 0.1 * lng * lat + 0.1 * Math.sqrt(Math.abs(lng))
  ret += (20.0 * Math.sin(6.0 * lng * Math.PI) + 20.0 * Math.sin(2.0 * lng * Math.PI)) * 2.0 / 3.0
  ret += (20.0 * Math.sin(lng * Math.PI) + 40.0 * Math.sin(lng / 3.0 * Math.PI)) * 2.0 / 3.0
  ret += (150.0 * Math.sin(lng / 12.0 * Math.PI) + 300.0 * Math.sin(lng / 30.0 * Math.PI)) * 2.0 / 3.0
  return ret
}
