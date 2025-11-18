<template>
  <el-dialog
    v-model="visible"
    title="飞机实时3D飞行轨迹"
    width="90%"
    :before-close="handleClose"
    destroy-on-close
  >
    <div class="flight-map-container">
      <div class="map-header">
        <div class="flight-info">
          <span><strong>飞机序列号:</strong> {{ aircraftSn }}</span>
          <span><strong>当前位置:</strong> {{ formatCoordinates() }}</span>
          <span><strong>飞行高度:</strong> {{ aircraftData?.height?.toFixed(2) || '未知' }}m</span>
          <span><strong>飞行速度:</strong> {{ aircraftData?.horizontal_speed?.toFixed(2) || '未知' }}m/s</span>
        </div>
          <div class="map-controls">
            <el-button @click="clearTrajectory" size="small">清除轨迹</el-button>
            <el-button @click="toggleTrajectory" size="small" :type="showTrajectory ? 'primary' : 'default'">
              {{ showTrajectory ? '隐藏轨迹' : '显示轨迹' }}
            </el-button>
            <el-button @click="centerMap" size="small">居中地图</el-button>
            <el-button @click="set3DView" size="small" type="success">3D视角</el-button>
            <el-button @click="setTopView" size="small">俯视图</el-button>
            <el-button @click="reloadMap" size="small" type="warning">重新加载地图</el-button>
          </div>
      </div>
      <div id="flight-map" class="flight-map"></div>
      <div class="trajectory-info">
        <div class="info-item">
          <span>轨迹点数: {{ flightTrajectory.length }}</span>
          <span>飞行距离: {{ calculateTotalDistance() }}km</span>
          <span>飞行时间: {{ formatFlightDuration() }}</span>
        </div>
        <div class="height-info" v-if="flightTrajectory.length > 0">
          <span>高度范围: {{ getHeightRange() }}</span>
          <span>当前高度: {{ aircraftData?.height?.toFixed(1) || '未知' }}m</span>
          <span>平均高度: {{ getAverageHeight() }}m</span>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { AMAP_CONFIG, waitForAMap, getAMapStatus, COORDINATE_CONVERTER } from '@/config/amap'

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  aircraftData: {
    type: Object,
    default: () => ({})
  },
  aircraftSn: {
    type: String,
    default: ''
  }
})

// Emits
const emit = defineEmits(['update:modelValue'])

// 响应式数据
const visible = ref(false)
const showTrajectory = ref(true)
const flightTrajectory = ref([])
const mapInstance = ref(null)
const aircraftMarker = ref(null)
const trajectoryLine = ref(null)
const mapInitialized = ref(false)
const mapErrorCount = ref(0)
const maxErrorCount = 3
const isRecovering = ref(false)
const heightLabels = ref([])
const heightBars = ref([])
const buildingHighlights = ref([])

// 监听显示状态
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) {
    nextTick(() => {
      initMap()
    })
  }
})

// 监听visible变化
watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

// 监听飞机数据变化
watch(() => props.aircraftData, (newData) => {
  if (newData && mapInitialized.value && visible.value) {
    updateAircraftMarker()
  }
}, { deep: true })

// 全局错误监听器
const setupGlobalErrorHandler = () => {
  // 监听高德地图相关的错误
  const originalConsoleError = console.error
  let errorCount = 0
  let lastErrorTime = 0
  
  console.error = (...args) => {
    const message = args.join(' ')
    const now = Date.now()
    
    // 检查是否是高德地图相关的错误
    if (message.includes("Cannot read properties of null (reading 'split')") ||
        message.includes("Cannot read properties of undefined (reading 'Vg')") ||
        message.includes("Cannot read properties of null (reading 'Ud')")) {
      
      // 防止频繁触发恢复
      if (now - lastErrorTime < 5000) { // 5秒内不重复触发
        console.warn('🚨 高德地图错误已检测到，但恢复间隔太短，跳过')
        originalConsoleError.apply(console, args)
        return
      }
      
      errorCount++
      lastErrorTime = now
      
      console.warn(`🚨 检测到高德地图内部错误 (第${errorCount}次):`, message)
      
      // 只有在错误次数较多时才尝试恢复
      if (errorCount >= 3) {
        console.warn('🚨 错误次数过多，尝试恢复地图')
        
        // 延迟恢复，避免频繁触发
        setTimeout(() => {
          if (mapInstance.value && !isRecovering.value) {
            recoverMap()
          }
        }, 2000)
        
        // 重置错误计数
        errorCount = 0
      }
    }
    
    // 调用原始console.error
    originalConsoleError.apply(console, args)
  }
}

// 暂时禁用全局错误监听器，避免过度恢复
// setupGlobalErrorHandler()

// 关闭处理
const handleClose = (done) => {
  visible.value = false
  if (done) done()
}

// 初始化地图
const initMap = async () => {
  if (mapInitialized.value) return

  try {
    console.log('🗺️ 开始初始化飞行轨迹地图...')
    console.log('📊 高德地图API状态:', getAMapStatus())

    // 检查API是否已经可用
    if (typeof window.AMap !== 'undefined') {
      console.log('✅ 检测到AMap已可用，直接使用')
      window.AMapLoaded = true
    }

    // 等待高德地图API加载完成
    await waitForAMap()
    console.log('✅ 高德地图API加载完成')

    // 使用全局AMap对象
    const AMap = window.AMap
    
    // 检查AMap对象是否完整
    if (!AMap || !AMap.Map) {
      throw new Error('AMap对象不完整，无法创建地图')
    }
    
    const { latitude, longitude } = props.aircraftData
    let center
    if (latitude && longitude) {
      // WGS84转GCJ02坐标转换
      const [gcjLng, gcjLat] = COORDINATE_CONVERTER.wgs84ToGcj02(longitude, latitude)
      center = [gcjLng, gcjLat]
      console.log('🗺️ 地图中心点转换:', {
        wgs84: { lng: longitude, lat: latitude },
        gcj02: { lng: gcjLng, lat: gcjLat }
      })
    } else {
      center = AMAP_CONFIG.defaultConfig.center
    }
    
    console.log('📍 地图中心点:', center)
    console.log('🔧 地图配置:', AMAP_CONFIG.defaultConfig)
    
    mapInstance.value = new AMap.Map('flight-map', {
      zoom: AMAP_CONFIG.defaultConfig.zoom,
      center: center,
      mapStyle: AMAP_CONFIG.defaultConfig.mapStyle,
      viewMode: '3D', // 强制使用3D模式
      pitch: 45, // 俯仰角，3D视角
      features: ['bg', 'road', 'building', 'point'],
      showLabel: true,
      buildingAnimation: true, // 启用建筑动画
      skyColor: '#87CEEB' // 天空颜色
    })
    
    // 检查地图实例是否创建成功
    if (!mapInstance.value) {
      throw new Error('地图实例创建失败')
    }
    
    // 监听地图瓦片加载事件
    mapInstance.value.on('complete', () => {
      console.log('✅ 地图瓦片加载完成')
    })
    
    mapInstance.value.on('error', (e) => {
      console.error('❌ 地图瓦片加载错误:', e)
    })
    
    // 监听地图缩放事件，添加错误处理
    mapInstance.value.on('zoomchange', () => {
      try {
        console.log('🔍 地图缩放级别:', mapInstance.value.getZoom())
      } catch (e) {
        console.warn('⚠️ 获取地图缩放级别时出错:', e)
      }
    })
    
    // 监听地图移动事件，添加错误处理
    mapInstance.value.on('moveend', () => {
      try {
        const center = mapInstance.value.getCenter()
        console.log('📍 地图中心点:', center)
      } catch (e) {
        console.warn('⚠️ 获取地图中心点时出错:', e)
      }
    })
    
    // 添加地图状态检查
    setTimeout(() => {
      if (mapInstance.value && mapInstance.value.getContainer) {
        try {
          const container = mapInstance.value.getContainer()
          if (!container || container.offsetWidth === 0) {
            console.warn('⚠️ 地图容器可能未正确初始化')
          }
        } catch (e) {
          console.warn('⚠️ 检查地图容器状态时出错:', e)
        }
      }
    }, 1000)
    
    // 添加控件（简化版本）
    try {
      mapInstance.value.addControl(new AMap.Scale())
      console.log('✅ 比例尺控件添加成功')
    } catch (error) {
      console.warn('⚠️ 比例尺控件添加失败:', error)
    }
    
    try {
      mapInstance.value.addControl(new AMap.ToolBar())
      console.log('✅ 工具条控件添加成功')
    } catch (error) {
      console.warn('⚠️ 工具条控件添加失败:', error)
    }
    
    // 添加飞机标记
    updateAircraftMarker()
    mapInitialized.value = true
    console.log('✅ 飞行轨迹地图初始化成功')
    
    // 检查地图瓦片加载状态
    setTimeout(() => {
      checkMapTilesLoaded()
    }, 3000)
    
    ElMessage.success('飞行轨迹地图已加载')
  } catch (error) {
    console.error('❌ 地图初始化失败:', error)
    console.error('📊 当前API状态:', getAMapStatus())
    
    // 尝试备用方案
    if (typeof window.AMap !== 'undefined') {
      console.log('🔄 尝试备用初始化方案...')
      try {
        const { latitude, longitude } = props.aircraftData
        const center = [longitude || AMAP_CONFIG.defaultConfig.center[0], latitude || AMAP_CONFIG.defaultConfig.center[1]]
        
        mapInstance.value = new window.AMap.Map('flight-map', {
          zoom: 15,
          center: center
        })
        
        updateAircraftMarker()
        mapInitialized.value = true
        console.log('✅ 备用方案初始化成功')
        ElMessage.success('飞行轨迹地图已加载（备用方案）')
        return
      } catch (backupError) {
        console.error('❌ 备用方案也失败:', backupError)
      }
    }
    
    ElMessage.error(`地图服务加载失败: ${error.message}`)
  }
}

// 更新飞机标记
const updateAircraftMarker = async () => {
  if (!mapInstance.value || !props.aircraftData) return
  
  // 简化健康检查，只检查基本状态
  if (!mapInstance.value || !mapInitialized.value) {
    console.warn('⚠️ 地图未初始化，跳过标记更新')
    return
  }
  
  const { latitude, longitude, height, attitude_head } = props.aircraftData
  
  // 严格检查坐标有效性
  if (latitude && longitude && 
      typeof latitude === 'number' && typeof longitude === 'number' &&
      !isNaN(latitude) && !isNaN(longitude) &&
      latitude >= -90 && latitude <= 90 && 
      longitude >= -180 && longitude <= 180) {
    
    try {
      // 移除旧标记
      if (aircraftMarker.value) {
        mapInstance.value.remove(aircraftMarker.value)
      }
      
      // WGS84转GCJ02坐标转换
      const [gcjLng, gcjLat] = COORDINATE_CONVERTER.wgs84ToGcj02(longitude, latitude)
      
      // 验证转换结果
      if (gcjLng === null || gcjLat === null || isNaN(gcjLng) || isNaN(gcjLat)) {
        console.error('❌ 坐标转换失败，使用原始坐标', { longitude, latitude, gcjLng, gcjLat })
        // 使用原始坐标
        aircraftMarker.value = new AMap.Marker({
          position: [longitude, latitude],
          title: `飞机位置\n高度: ${height?.toFixed(2) || '未知'}m\n航向: ${attitude_head?.toFixed(1) || '未知'}°\nWGS84: ${latitude.toFixed(8)}, ${longitude.toFixed(8)}`,
          icon: new AMap.Icon({
            size: new AMap.Size(32, 32),
            image: 'data:image/svg+xml;base64,' + btoa(`
              <svg width="32" height="32" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                <path d="M16 2L20 12L30 12L22 18L26 28L16 22L6 28L10 18L2 12L12 12Z" fill="#ff4444" stroke="#ffffff" stroke-width="2"/>
              </svg>
            `)
          }),
          rotation: attitude_head || 0
        })
      } else {
        console.log('🔄 坐标转换:', {
          wgs84: { lng: longitude, lat: latitude },
          gcj02: { lng: gcjLng, lat: gcjLat }
        })
        
        // 创建新标记（使用转换后的GCJ02坐标）
        aircraftMarker.value = new AMap.Marker({
          position: [gcjLng, gcjLat],
          title: `飞机位置\n高度: ${height?.toFixed(2) || '未知'}m\n航向: ${attitude_head?.toFixed(1) || '未知'}°\nWGS84: ${latitude.toFixed(8)}, ${longitude.toFixed(8)}\nGCJ02: ${gcjLat.toFixed(8)}, ${gcjLng.toFixed(8)}`,
          icon: new AMap.Icon({
            size: new AMap.Size(32, 32),
            image: 'data:image/svg+xml;base64,' + btoa(`
              <svg width="32" height="32" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                <path d="M16 2L20 12L30 12L22 18L26 28L16 22L6 28L10 18L2 12L12 12Z" fill="#ff4444" stroke="#ffffff" stroke-width="2"/>
              </svg>
            `)
          }),
          rotation: attitude_head || 0
        })
      }
      
      mapInstance.value.add(aircraftMarker.value)
      
      // 添加到轨迹（使用转换后的坐标或原始坐标）
      if (gcjLng !== null && gcjLat !== null && !isNaN(gcjLng) && !isNaN(gcjLat)) {
        addToTrajectory(gcjLng, gcjLat, height)
      } else {
        addToTrajectory(longitude, latitude, height)
      }
      
      // 更新轨迹线
      updateTrajectoryLine()
      
      console.log('✅ 飞机标记更新成功（已转换坐标）')
    } catch (error) {
      console.error('❌ 更新飞机标记失败:', error)
      
      // 简化错误处理，只记录错误，不进行恢复
      console.warn('⚠️ 飞机标记更新失败，将在下次数据更新时重试')
    }
  } else {
    console.warn('⚠️ 无效的坐标数据:', { latitude, longitude })
  }
}

// 添加到轨迹
const addToTrajectory = (lng, lat, height) => {
  // 严格检查输入参数
  if (!lng || !lat || 
      typeof lng !== 'number' || typeof lat !== 'number' ||
      isNaN(lng) || isNaN(lat)) {
    console.warn('⚠️ 无效的轨迹点数据:', { lng, lat })
    return
  }
  
  try {
    const point = {
      lng: Number(lng),
      lat: Number(lat),
      height: height && typeof height === 'number' ? Number(height) : 0,
      timestamp: Date.now()
    }
    
    flightTrajectory.value.push(point)
    
    // 限制轨迹点数量（保留最近1000个点）
    if (flightTrajectory.value.length > 1000) {
      flightTrajectory.value.shift()
    }
    
    console.log('✅ 轨迹点添加成功:', point)
  } catch (error) {
    console.error('❌ 添加轨迹点失败:', error)
  }
}

// 更新轨迹线
const updateTrajectoryLine = () => {
  if (!mapInstance.value || !showTrajectory.value || flightTrajectory.value.length < 2) return
  
  try {
    // 移除旧轨迹线
    if (trajectoryLine.value) {
      mapInstance.value.remove(trajectoryLine.value)
    }
    
    // 过滤有效的轨迹点
    const validPoints = flightTrajectory.value.filter(point => 
      point && 
      typeof point.lng === 'number' && typeof point.lat === 'number' &&
      !isNaN(point.lng) && !isNaN(point.lat)
    )
    
    if (validPoints.length < 2) {
      console.warn('⚠️ 有效轨迹点不足，无法绘制轨迹线')
      return
    }
    
    // 创建带高度信息的轨迹线
    const path = validPoints.map(point => [point.lng, point.lat])
    
    // 根据高度变化设置轨迹线颜色
    const heights = validPoints.map(p => p.height || 0)
    const minHeight = Math.min(...heights)
    const maxHeight = Math.max(...heights)
    const heightRange = maxHeight - minHeight
    
    // 创建分段轨迹线，根据高度变化显示不同颜色
    const segments = []
    for (let i = 0; i < validPoints.length - 1; i++) {
      const currentPoint = validPoints[i]
      const nextPoint = validPoints[i + 1]
      
      // 根据高度计算颜色
      const currentHeight = currentPoint.height || 0
      const heightRatio = heightRange > 0 ? (currentHeight - minHeight) / heightRange : 0.5
      
      // 高度越高，颜色越红；高度越低，颜色越蓝
      const red = Math.floor(255 * heightRatio)
      const blue = Math.floor(255 * (1 - heightRatio))
      const color = `rgb(${red}, 0, ${blue})`
      
      segments.push({
        path: [[currentPoint.lng, currentPoint.lat], [nextPoint.lng, nextPoint.lat]],
        color: color,
        height: currentHeight
      })
    }
    
    // 创建3D分段轨迹线
    const polylines = segments.map(segment => {
      return new AMap.Polyline({
        path: segment.path,
        strokeColor: segment.color,
        strokeWeight: 6, // 增加线宽，3D效果更明显
        strokeOpacity: 0.9,
        strokeStyle: 'solid',
        zIndex: 1000 + Math.floor(segment.height / 10) // 根据高度设置层级
      })
    })
    
    // 将多个线段组合成一个轨迹线对象
    trajectoryLine.value = {
      polylines: polylines,
      add: function(map) {
        polylines.forEach(polyline => map.add(polyline))
      },
      remove: function(map) {
        polylines.forEach(polyline => map.remove(polyline))
      }
    }
    
    trajectoryLine.value.add(mapInstance.value)
    
    // 添加高度标签
    addHeightLabels(validPoints)
    
    // 添加3D高度柱状图
    addHeightBars(validPoints)
    
    // 添加3D建筑高亮效果
    addBuildingHighlight(validPoints)
    
    console.log('✅ 3D轨迹线更新成功，点数:', validPoints.length, '高度范围:', minHeight.toFixed(1), '-', maxHeight.toFixed(1), 'm')
  } catch (error) {
    console.error('❌ 更新轨迹线失败:', error)
  }
}

// 添加高度标签
const addHeightLabels = (points) => {
  // 清除旧的高度标签
  clearHeightLabels()
  
  if (!mapInstance.value || points.length === 0) return
  
  try {
    // 每隔一定间隔添加高度标签
    const labelInterval = Math.max(1, Math.floor(points.length / 10)) // 最多显示10个标签
    
    for (let i = 0; i < points.length; i += labelInterval) {
      const point = points[i]
      if (point && point.height !== undefined && point.height !== null) {
        const label = new AMap.Text({
          position: [point.lng, point.lat],
          text: `${point.height.toFixed(1)}m`,
          style: {
            'background-color': 'rgba(255, 255, 255, 0.8)',
            'border': '1px solid #ccc',
            'border-radius': '3px',
            'padding': '2px 4px',
            'font-size': '12px',
            'color': '#333'
          },
          offset: new AMap.Pixel(0, -20)
        })
        
        mapInstance.value.add(label)
        heightLabels.value.push(label)
      }
    }
    
    console.log('✅ 高度标签添加成功，数量:', heightLabels.value.length)
  } catch (error) {
    console.error('❌ 添加高度标签失败:', error)
  }
}

// 清除高度标签
const clearHeightLabels = () => {
  if (mapInstance.value && heightLabels.value.length > 0) {
    heightLabels.value.forEach(label => {
      try {
        mapInstance.value.remove(label)
      } catch (e) {
        console.warn('⚠️ 移除高度标签时出错:', e)
      }
    })
    heightLabels.value = []
  }
}

// 添加3D高度柱状图
const addHeightBars = (points) => {
  // 清除旧的高度柱状图
  clearHeightBars()
  
  if (!mapInstance.value || points.length === 0) return
  
  try {
    // 每隔一定间隔添加高度柱状图
    const barInterval = Math.max(1, Math.floor(points.length / 20)) // 最多显示20个柱状图
    
    for (let i = 0; i < points.length; i += barInterval) {
      const point = points[i]
      if (point && point.height !== undefined && point.height !== null && point.height > 0) {
        // 使用圆形标记模拟3D柱状图
        const heightBar = new AMap.Circle({
          center: [point.lng, point.lat],
          radius: 20, // 半径（米）
          fillColor: getHeightColor(point.height, points),
          fillOpacity: 0.6,
          strokeColor: '#ffffff',
          strokeWeight: 2,
          strokeOpacity: 0.8,
          zIndex: 1000 + Math.floor(point.height / 10) // 根据高度设置层级
        })
        
        // 添加高度标签
        const heightLabel = new AMap.Text({
          position: [point.lng, point.lat],
          text: `${point.height.toFixed(0)}m`,
          style: {
            'background-color': 'rgba(255, 255, 255, 0.9)',
            'border': '1px solid #ccc',
            'border-radius': '3px',
            'padding': '2px 4px',
            'font-size': '11px',
            'color': '#333',
            'font-weight': 'bold'
          },
          offset: new AMap.Pixel(0, -30)
        })
        
        mapInstance.value.add(heightBar)
        mapInstance.value.add(heightLabel)
        heightBars.value.push({ circle: heightBar, label: heightLabel })
      }
    }
    
    console.log('✅ 3D高度柱状图添加成功，数量:', heightBars.value.length)
  } catch (error) {
    console.error('❌ 添加3D高度柱状图失败:', error)
  }
}

// 根据高度获取颜色
const getHeightColor = (height, allPoints) => {
  const heights = allPoints.map(p => p.height).filter(h => h !== undefined && h !== null)
  const minHeight = Math.min(...heights)
  const maxHeight = Math.max(...heights)
  const heightRange = maxHeight - minHeight
  
  if (heightRange === 0) return '#00ff00'
  
  const ratio = (height - minHeight) / heightRange
  const red = Math.floor(255 * ratio)
  const blue = Math.floor(255 * (1 - ratio))
  
  return `rgb(${red}, 0, ${blue})`
}

// 清除3D高度柱状图
const clearHeightBars = () => {
  if (mapInstance.value && heightBars.value.length > 0) {
    heightBars.value.forEach(bar => {
      try {
        if (bar.circle) {
          mapInstance.value.remove(bar.circle)
        }
        if (bar.label) {
          mapInstance.value.remove(bar.label)
        }
      } catch (e) {
        console.warn('⚠️ 移除3D高度柱状图时出错:', e)
      }
    })
    heightBars.value = []
  }
}

// 添加3D建筑高亮效果
const addBuildingHighlight = (points) => {
  // 清除旧的建筑高亮
  clearBuildingHighlight()
  
  if (!mapInstance.value || points.length === 0) return
  
  try {
    // 在轨迹点附近添加建筑高亮
    const highlightInterval = Math.max(1, Math.floor(points.length / 10)) // 最多显示10个高亮
    
    for (let i = 0; i < points.length; i += highlightInterval) {
      const point = points[i]
      if (point && point.lng && point.lat) {
        // 创建圆形区域高亮
        const highlight = new AMap.Circle({
          center: [point.lng, point.lat],
          radius: 50, // 半径（米）
          fillColor: '#ff6b6b',
          fillOpacity: 0.2,
          strokeColor: '#ff6b6b',
          strokeWeight: 3,
          strokeOpacity: 0.8,
          zIndex: 500
        })
        
        mapInstance.value.add(highlight)
        buildingHighlights.value.push(highlight)
      }
    }
    
    console.log('✅ 3D建筑高亮添加成功，数量:', buildingHighlights.value.length)
  } catch (error) {
    console.error('❌ 添加3D建筑高亮失败:', error)
  }
}

// 清除3D建筑高亮
const clearBuildingHighlight = () => {
  if (mapInstance.value && buildingHighlights.value.length > 0) {
    buildingHighlights.value.forEach(highlight => {
      try {
        mapInstance.value.remove(highlight)
      } catch (e) {
        console.warn('⚠️ 移除3D建筑高亮时出错:', e)
      }
    })
    buildingHighlights.value = []
  }
}

// 清除轨迹
const clearTrajectory = () => {
  flightTrajectory.value = []
  clearHeightLabels()
  clearHeightBars()
  clearBuildingHighlight()
  if (trajectoryLine.value) {
    trajectoryLine.value.remove(mapInstance.value)
    trajectoryLine.value = null
  }
  ElMessage.success('3D轨迹已清除')
}

// 切换轨迹显示
const toggleTrajectory = () => {
  showTrajectory.value = !showTrajectory.value
  if (showTrajectory.value) {
    updateTrajectoryLine()
  } else {
    // 隐藏轨迹线、高度标签、3D柱状图和建筑高亮
    if (trajectoryLine.value) {
      trajectoryLine.value.remove(mapInstance.value)
      trajectoryLine.value = null
    }
    clearHeightLabels()
    clearHeightBars()
    clearBuildingHighlight()
  }
}

// 居中地图
const centerMap = () => {
  if (mapInstance.value && props.aircraftData) {
    const { latitude, longitude } = props.aircraftData
    if (latitude && longitude) {
      // WGS84转GCJ02坐标转换
      const [gcjLng, gcjLat] = COORDINATE_CONVERTER.wgs84ToGcj02(longitude, latitude)
      mapInstance.value.setCenter([gcjLng, gcjLat])
      mapInstance.value.setZoom(15)
      console.log('🎯 地图居中（已转换坐标）:', { wgs84: [longitude, latitude], gcj02: [gcjLng, gcjLat] })
    }
  }
}

// 设置3D视角
const set3DView = () => {
  if (mapInstance.value) {
    try {
      // 高德地图3D视角设置
      mapInstance.value.setPitch(45) // 设置俯仰角
      mapInstance.value.setRotation(0) // 设置旋转角度
      mapInstance.value.setZoom(mapInstance.value.getZoom()) // 刷新地图
      console.log('🎬 切换到3D视角')
      ElMessage.success('已切换到3D视角')
    } catch (error) {
      console.warn('⚠️ 3D视角设置失败:', error)
      ElMessage.warning('3D视角设置失败，地图已为3D模式')
    }
  }
}

// 设置俯视图
const setTopView = () => {
  if (mapInstance.value) {
    try {
      // 高德地图俯视图设置
      mapInstance.value.setPitch(0) // 俯视角度
      mapInstance.value.setRotation(0) // 无旋转
      mapInstance.value.setZoom(mapInstance.value.getZoom()) // 刷新地图
      console.log('📐 切换到俯视图')
      ElMessage.success('已切换到俯视图')
    } catch (error) {
      console.warn('⚠️ 俯视图设置失败:', error)
      ElMessage.warning('俯视图设置失败')
    }
  }
}

// 重新加载地图
const reloadMap = () => {
  if (!mapInstance.value) return
  
  try {
    console.log('🔄 重新加载地图...')
    ElMessage.info('正在重新加载地图...')
    
    // 安全销毁当前地图
    safeDestroyMap()
    
    // 延迟重新初始化
    setTimeout(() => {
      initMap()
    }, 500)
  } catch (error) {
    console.error('❌ 重新加载地图失败:', error)
    ElMessage.error('重新加载地图失败')
  }
}

// 格式化坐标
const formatCoordinates = () => {
  try {
    const { latitude, longitude } = props.aircraftData || {}
    
    if (latitude && longitude && 
        typeof latitude === 'number' && typeof longitude === 'number' &&
        !isNaN(latitude) && !isNaN(longitude)) {
      // 显示转换后的GCJ02坐标
      const [gcjLng, gcjLat] = COORDINATE_CONVERTER.wgs84ToGcj02(longitude, latitude)
      return `GCJ02: ${gcjLat.toFixed(8)}, ${gcjLng.toFixed(8)}`
    }
    
    return '未知'
  } catch (error) {
    console.error('❌ 格式化坐标失败:', error)
    return '未知'
  }
}

// 计算总飞行距离
const calculateTotalDistance = () => {
  if (flightTrajectory.value.length < 2) return '0.00'
  
  try {
    // 过滤有效的轨迹点
    const validPoints = flightTrajectory.value.filter(point => 
      point && 
      typeof point.lng === 'number' && typeof point.lat === 'number' &&
      !isNaN(point.lng) && !isNaN(point.lat)
    )
    
    if (validPoints.length < 2) return '0.00'
    
    let totalDistance = 0
    for (let i = 1; i < validPoints.length; i++) {
      const prev = validPoints[i - 1]
      const curr = validPoints[i]
      
      if (prev && curr && 
          typeof prev.lng === 'number' && typeof prev.lat === 'number' &&
          typeof curr.lng === 'number' && typeof curr.lat === 'number') {
        const distance = AMap.GeometryUtil.distance([prev.lng, prev.lat], [curr.lng, curr.lat])
        if (!isNaN(distance) && distance > 0) {
          totalDistance += distance
        }
      }
    }
    
    return (totalDistance / 1000).toFixed(2)
  } catch (error) {
    console.error('❌ 计算飞行距离失败:', error)
    return '0.00'
  }
}

// 格式化飞行时长
const formatFlightDuration = () => {
  if (flightTrajectory.value.length < 2) return '00:00:00'
  
  const startTime = flightTrajectory.value[0].timestamp
  const endTime = flightTrajectory.value[flightTrajectory.value.length - 1].timestamp
  const duration = Math.floor((endTime - startTime) / 1000)
  
  const hours = Math.floor(duration / 3600)
  const minutes = Math.floor((duration % 3600) / 60)
  const seconds = duration % 60
  
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
}

// 获取高度范围
const getHeightRange = () => {
  if (flightTrajectory.value.length === 0) return '未知'
  
  const heights = flightTrajectory.value
    .map(point => point.height)
    .filter(height => height !== undefined && height !== null && !isNaN(height))
  
  if (heights.length === 0) return '未知'
  
  const minHeight = Math.min(...heights)
  const maxHeight = Math.max(...heights)
  
  return `${minHeight.toFixed(1)}m - ${maxHeight.toFixed(1)}m`
}

// 获取平均高度
const getAverageHeight = () => {
  if (flightTrajectory.value.length === 0) return '未知'
  
  const heights = flightTrajectory.value
    .map(point => point.height)
    .filter(height => height !== undefined && height !== null && !isNaN(height))
  
  if (heights.length === 0) return '未知'
  
  const averageHeight = heights.reduce((sum, height) => sum + height, 0) / heights.length
  return averageHeight.toFixed(1)
}

// 检查地图瓦片加载状态
const checkMapTilesLoaded = () => {
  if (!mapInstance.value) return
  
  try {
    // 检查地图容器
    const mapContainer = document.getElementById('flight-map')
    if (!mapContainer) {
      console.error('❌ 地图容器未找到')
      return
    }
    
    // 检查地图瓦片
    const tiles = mapContainer.querySelectorAll('img[src*="amap.com"]')
    console.log('🔍 地图瓦片检查:', {
      tilesCount: tiles.length,
      containerSize: {
        width: mapContainer.offsetWidth,
        height: mapContainer.offsetHeight
      }
    })
    
    if (tiles.length === 0) {
      console.warn('⚠️ 未检测到地图瓦片，尝试重新加载')
      ElMessage.warning('地图瓦片未加载，尝试重新加载...')
      
      // 尝试重新设置地图中心
      const { latitude, longitude } = props.aircraftData
      if (latitude && longitude) {
        mapInstance.value.setCenter([longitude, latitude])
        mapInstance.value.setZoom(15)
      }
    } else {
      console.log('✅ 地图瓦片已加载')
    }
  } catch (error) {
    console.error('❌ 检查地图瓦片状态失败:', error)
  }
}

// 组件挂载时初始化
onMounted(() => {
  if (props.modelValue) {
    nextTick(() => {
      initMap()
    })
  }
})

// 地图健康检查
const checkMapHealth = () => {
  if (!mapInstance.value) return false
  
  try {
    // 检查地图实例是否仍然有效
    const center = mapInstance.value.getCenter()
    const zoom = mapInstance.value.getZoom()
    
    if (!center || !zoom || isNaN(zoom)) {
      console.warn('⚠️ 地图实例状态异常')
      return false
    }
    
    return true
  } catch (error) {
    console.warn('⚠️ 地图健康检查失败:', error)
    return false
  }
}

// 地图错误恢复
const recoverMap = async () => {
  if (isRecovering.value) return
  
  isRecovering.value = true
  mapErrorCount.value++
  
  console.log(`🔄 开始地图恢复 (第${mapErrorCount.value}次尝试)`)
  
  try {
    // 保存当前飞机数据
    const currentAircraftData = props.aircraftData
    
    // 完全销毁当前地图
    safeDestroyMap()
    
    // 等待一段时间让资源完全释放
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 重新初始化地图
    await initMap()
    
    // 恢复飞机标记
    if (currentAircraftData && mapInitialized.value) {
      console.log('🔄 恢复飞机标记')
      await updateAircraftMarker()
    }
    
    console.log('✅ 地图恢复成功')
    mapErrorCount.value = 0
  } catch (error) {
    console.error('❌ 地图恢复失败:', error)
    
    if (mapErrorCount.value >= maxErrorCount) {
      console.error('❌ 地图恢复次数超限，停止尝试')
      ElMessage.error('地图初始化失败，请刷新页面重试')
    }
  } finally {
    isRecovering.value = false
  }
}

// 安全销毁地图
const safeDestroyMap = () => {
  if (mapInstance.value) {
    try {
      // 清理所有覆盖物
      if (aircraftMarker.value) {
        try {
          mapInstance.value.remove(aircraftMarker.value)
        } catch (e) {
          console.warn('⚠️ 移除飞机标记时出错:', e)
        }
        aircraftMarker.value = null
      }

      if (trajectoryLine.value) {
        try {
          trajectoryLine.value.remove(mapInstance.value)
        } catch (e) {
          console.warn('⚠️ 移除轨迹线时出错:', e)
        }
        trajectoryLine.value = null
      }

      // 清理高度标签、3D柱状图和建筑高亮
      clearHeightLabels()
      clearHeightBars()
      clearBuildingHighlight()

      // 销毁地图实例
      mapInstance.value.destroy()
      console.log('🗑️ 3D地图实例已安全销毁')
    } catch (error) {
      console.error('❌ 销毁地图时出错:', error)
    } finally {
      mapInstance.value = null
      mapInitialized.value = false
      flightTrajectory.value = []
    }
  }
}

// 组件卸载时清理
onUnmounted(() => {
  safeDestroyMap()
})
</script>

<style scoped>
/* 飞行地图样式 */
.flight-map-container {
  height: 600px;
  display: flex;
  flex-direction: column;
}

.map-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  background: var(--el-color-info-light-9);
  border-radius: 4px;
  margin-bottom: 10px;
}

.flight-info {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.flight-info span {
  font-size: 14px;
  color: var(--el-text-color-primary);
}

.map-controls {
  display: flex;
  gap: 10px;
}

.flight-map {
  flex: 1;
  border-radius: 4px;
  border: 1px solid var(--el-border-color);
}

.trajectory-info {
  padding: 10px;
  background: var(--el-color-success-light-9);
  border-radius: 4px;
  margin-top: 10px;
}

.trajectory-info .info-item {
  display: flex;
  gap: 20px;
  justify-content: center;
}

.trajectory-info span {
  font-size: 14px;
  color: var(--el-text-color-primary);
  font-weight: 500;
}

.height-info {
  display: flex;
  gap: 20px;
  justify-content: center;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid var(--el-border-color-light);
}

.height-info span {
  font-size: 13px;
  color: var(--el-color-warning);
  font-weight: 600;
  background: var(--el-color-warning-light-8);
  padding: 4px 8px;
  border-radius: 4px;
}
</style>
