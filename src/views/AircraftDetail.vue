<template>
  <div class="aircraft-detail">
    <el-container style="height: 100vh;">
      <!-- 头部导航 -->
      <el-header class="detail-header">
        <div class="header-left">
          <el-button 
            type="primary" 
            :icon="ArrowLeft" 
            @click="goBack"
            style="margin-right: 16px;"
          >
            返回
          </el-button>
          <h2>飞机信息详情 - {{ deviceSn }}</h2>
        </div>
        <div class="header-right">
          <div v-if="hasReceivedData && osdData" class="data-update-info">
            <el-tag type="info" size="small" style="margin-right: 8px;">
              最后更新: {{ formatTime(osdData.timestamp) }}
            </el-tag>
          </div>
          <el-tag 
            :type="mqttProxyStore.isConnected ? 'success' : 'danger'"
            size="large"
          >
            {{ mqttProxyStore.isConnected ? '已连接' : '未连接' }}
          </el-tag>
        </div>
      </el-header>

      <!-- 主要内容区域 -->
      <el-main class="detail-main">
        <!-- 连接状态提示 -->
        <div v-if="!mqttProxyStore.isConnected" class="connection-warning">
          <el-alert
            title="MQTT未连接"
            description="请先在仪表板页面连接MQTT服务器以获取实时数据"
            type="warning"
            show-icon
            :closable="false"
          />
        </div>

        <!-- 数据加载状态 -->
        <div v-if="mqttProxyStore.isConnected && !hasReceivedData" class="loading-state">
          <el-empty description="等待OSD数据...">
            <el-button type="primary" @click="subscribeToOSD">
              手动订阅OSD主题
            </el-button>
          </el-empty>
        </div>

        <!-- 实时数据展示 -->
        <div v-if="hasReceivedData" class="data-container">
          <!-- 基本信息卡片 -->
          <el-row :gutter="20" class="data-grid">
            <el-col :span="8">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Location /></el-icon>
                    <span>位置信息</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">纬度:</span>
                    <span class="value">{{ osdData?.data?.latitude?.toFixed(6) || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">经度:</span>
                    <span class="value">{{ osdData?.data?.longitude?.toFixed(6) || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">高度:</span>
                    <span class="value">{{ osdData?.data?.height?.toFixed(2) || '--' }}m</span>
                  </div>
                  <div class="info-item">
                    <span class="label">航向:</span>
                    <span class="value">{{ osdData?.data?.heading?.toFixed(1) || '--' }}°</span>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="8">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Lightning /></el-icon>
                    <span>电池状态</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">电量:</span>
                    <el-progress 
                      :percentage="osdData?.data?.drone_charge_state?.capacity_percent || 0"
                      :color="getBatteryColor(osdData?.data?.drone_charge_state?.capacity_percent)"
                      :show-text="true"
                    />
                  </div>
                  <div class="info-item">
                    <span class="label">充电状态:</span>
                    <el-tag :type="getChargeStateType(osdData?.data?.drone_charge_state?.state)">
                      {{ getChargeStateText(osdData?.data?.drone_charge_state?.state) }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">存储模式:</span>
                    <el-tag :type="getBatteryStoreModeType(osdData?.data?.battery_store_mode)">
                      {{ getBatteryStoreModeText(osdData?.data?.battery_store_mode) }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="8">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Cloudy /></el-icon>
                    <span>环境信息</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">环境温度:</span>
                    <span class="value">{{ osdData?.data?.environment_temperature?.toFixed(1) || '--' }}°C</span>
                  </div>
                  <div class="info-item">
                    <span class="label">设备温度:</span>
                    <span class="value">{{ osdData?.data?.temperature?.toFixed(1) || '--' }}°C</span>
                  </div>
                  <div class="info-item">
                    <span class="label">湿度:</span>
                    <span class="value">{{ osdData?.data?.humidity || '--' }}%</span>
                  </div>
                  <div class="info-item">
                    <span class="label">风速:</span>
                    <span class="value">{{ osdData?.data?.wind_speed || '--' }}m/s</span>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>

          <!-- 网络和状态信息 -->
          <el-row :gutter="20" class="data-grid" style="margin-top: 20px;">
            <el-col :span="12">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Connection /></el-icon>
                    <span>网络状态</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">网络类型:</span>
                    <el-tag :type="getNetworkTypeColor(osdData?.data?.network_state?.type)">
                      {{ getNetworkTypeText(osdData?.data?.network_state?.type) }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">信号质量:</span>
                    <el-progress 
                      :percentage="(osdData?.data?.network_state?.quality || 0) * 20"
                      :color="getNetworkQualityColor(osdData?.data?.network_state?.quality)"
                      :show-text="true"
                    />
                  </div>
                  <div class="info-item">
                    <span class="label">传输速率:</span>
                    <span class="value">{{ getNetworkRateText(osdData?.data?.network_state?.rate) }}</span>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="12">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Setting /></el-icon>
                    <span>设备状态</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">模式代码:</span>
                    <el-tag :type="getModeCodeType(osdData?.data?.mode_code)">
                      {{ getModeCodeText(osdData?.data?.mode_code) }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">在机库:</span>
                    <el-tag :type="osdData?.data?.drone_in_dock ? 'success' : 'info'">
                      {{ osdData?.data?.drone_in_dock ? '是' : '否' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">静默模式:</span>
                    <el-tag :type="osdData?.data?.silent_mode ? 'warning' : 'success'">
                      {{ osdData?.data?.silent_mode ? '开启' : '关闭' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">紧急停止:</span>
                    <el-tag :type="osdData?.data?.emergency_stop_state ? 'danger' : 'success'">
                      {{ osdData?.data?.emergency_stop_state ? '激活' : '正常' }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>

          <!-- 存储和子设备信息 -->
          <el-row :gutter="20" class="data-grid" style="margin-top: 20px;">
            <el-col :span="12">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Folder /></el-icon>
                    <span>存储信息</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">总容量:</span>
                    <span class="value">{{ formatStorage(osdData?.data?.storage?.total) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">已使用:</span>
                    <span class="value">{{ formatStorage(osdData?.data?.storage?.used) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">使用率:</span>
                    <el-progress 
                      :percentage="getStorageUsage(osdData?.data?.storage)"
                      :color="getStorageUsageColor(osdData?.data?.storage)"
                      :show-text="true"
                    />
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="12">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Monitor /></el-icon>
                    <span>子设备信息</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">设备SN:</span>
                    <span class="value">{{ osdData?.data?.sub_device?.device_sn || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">在线状态:</span>
                    <el-tag :type="osdData?.data?.sub_device?.device_online_status ? 'success' : 'danger'">
                      {{ osdData?.data?.sub_device?.device_online_status ? '在线' : '离线' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">配对状态:</span>
                    <el-tag :type="osdData?.data?.sub_device?.device_paired ? 'success' : 'warning'">
                      {{ osdData?.data?.sub_device?.device_paired ? '已配对' : '未配对' }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>

          <!-- 原始数据展示 -->
          <el-card class="raw-data-card" shadow="hover" style="margin-top: 20px;">
            <template #header>
              <div class="card-header">
                <el-icon><Document /></el-icon>
                <span>原始OSD数据</span>
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="copyRawData"
                  :icon="copySuccess ? 'Check' : 'CopyDocument'"
                  :class="{ 'copy-success': copySuccess }"
                  style="margin-left: auto;"
                >
                  {{ copySuccess ? '已复制' : '复制' }}
                </el-button>
              </div>
            </template>
            <div class="json-viewer">
              <pre class="json-content" v-html="formatJsonPayload(osdData)"></pre>
            </div>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  ArrowLeft, 
  Location, 
  Lightning, 
  Cloudy, 
  Connection, 
  Setting, 
  Folder, 
  Monitor, 
  Document 
} from '@element-plus/icons-vue'
import { useMqttProxyStore } from '@/stores/mqtt-proxy'

const route = useRoute()
const router = useRouter()
const mqttProxyStore = useMqttProxyStore()

// 响应式数据
const deviceSn = ref('')
const osdData = ref(null)
const hasReceivedData = ref(false)
const copySuccess = ref(false)
const messageListener = ref(null)
const lastProcessedMessageId = ref(null)
const updateTimeout = ref(null)
const dataStableCount = ref(0) // 数据稳定计数器

// 计算属性
const osdTopic = computed(() => `thing/product/${deviceSn.value}/osd`)

// 页面初始化
onMounted(() => {
  deviceSn.value = route.params.deviceSn || ''
  if (deviceSn.value) {
    subscribeToOSD()
  } else {
    ElMessage.error('未找到设备SN参数')
    goBack()
  }
})

// 页面卸载时取消订阅
onUnmounted(() => {
  stopMessageListener()
  if (updateTimeout.value) {
    clearTimeout(updateTimeout.value)
  }
  if (osdTopic.value) {
    mqttProxyStore.unsubscribeTopics(osdTopic.value)
  }
})

// 订阅OSD主题
const subscribeToOSD = async () => {
  if (!mqttProxyStore.isConnected) {
    ElMessage.warning('MQTT未连接，无法订阅OSD数据')
    return
  }

  try {
    await mqttProxyStore.subscribeToTopics(osdTopic.value, 1)
    ElMessage.success(`已订阅OSD主题: ${osdTopic.value}`)
    
    // 开始监听消息
    startMessageListener()
  } catch (error) {
    console.error('订阅OSD主题失败:', error)
    ElMessage.error('订阅OSD主题失败')
  }
}

// 开始消息监听
const startMessageListener = () => {
  if (messageListener.value) {
    clearInterval(messageListener.value)
  }
  
  messageListener.value = setInterval(() => {
    handleOSDMessage()
  }, 2000) // 每2秒检查一次新消息，进一步减少频率
}

// 停止消息监听
const stopMessageListener = () => {
  if (messageListener.value) {
    clearInterval(messageListener.value)
    messageListener.value = null
  }
}

// 处理OSD消息
const handleOSDMessage = () => {
  // 从消息历史中查找最新的OSD消息
  const osdMessages = mqttProxyStore.messageHistory.filter(msg => msg.topic === osdTopic.value)
  if (osdMessages.length > 0) {
    const latestMessage = osdMessages[0] // 最新的消息在数组开头
    
    // 检查是否是新消息（通过消息ID比较）
    if (latestMessage.id !== lastProcessedMessageId.value) {
      try {
        const newData = JSON.parse(latestMessage.payload)
        
        // 验证数据完整性
        if (isValidOSDData(newData)) {
          // 进一步检查数据是否真的有变化（比较关键字段）
          if (!osdData.value || hasDataChanged(osdData.value, newData)) {
            // 使用防抖更新，避免频繁更新
            if (updateTimeout.value) {
              clearTimeout(updateTimeout.value)
            }
            
            updateTimeout.value = setTimeout(() => {
              osdData.value = newData
              hasReceivedData.value = true
              lastProcessedMessageId.value = latestMessage.id
              
              console.log('OSD数据已更新:', new Date().toLocaleTimeString(), '消息ID:', latestMessage.id)
              console.log('OSD数据内容:', newData)
            }, 200) // 增加防抖延迟到200ms
          } else {
            // 数据没有实质性变化，只更新消息ID
            lastProcessedMessageId.value = latestMessage.id
            console.log('OSD数据无变化，跳过更新')
          }
        } else {
          console.warn('收到无效的OSD数据，跳过更新:', newData)
          // 即使数据无效，也更新消息ID，避免重复处理
          lastProcessedMessageId.value = latestMessage.id
        }
      } catch (error) {
        console.error('解析OSD数据失败:', error, '原始数据:', latestMessage.payload)
        // 解析失败也更新消息ID
        lastProcessedMessageId.value = latestMessage.id
      }
    }
  } else {
    // 如果没有OSD消息，但已经有数据，保持现有数据不变
    if (osdData.value && hasReceivedData.value) {
      console.log('暂无新OSD消息，保持现有数据显示')
      // 确保数据状态标志保持为true
      if (!hasReceivedData.value) {
        hasReceivedData.value = true
      }
    }
  }
}

// 验证OSD数据是否有效
const isValidOSDData = (data) => {
  if (!data || typeof data !== 'object') return false
  
  // 检查是否有data字段（OSD数据通常在这个字段中）
  if (data.data && typeof data.data === 'object') {
    // 如果有data字段，检查是否至少有一些基本字段
    const hasBasicFields = Object.keys(data.data).length > 0
    return hasBasicFields
  }
  
  // 如果没有data字段，检查是否直接包含OSD相关字段
  const osdFields = ['latitude', 'longitude', 'height', 'drone_charge_state', 'battery', 'altitude', 'speed']
  const hasOSDFields = osdFields.some(field => field in data)
  
  return hasOSDFields
}

// 检查数据是否真的有变化
const hasDataChanged = (oldData, newData) => {
  if (!oldData || !newData) return true
  
  // 比较关键字段
  const oldTimestamp = oldData.timestamp || 0
  const newTimestamp = newData.timestamp || 0
  
  // 如果时间戳不同，说明是新数据
  if (oldTimestamp !== newTimestamp) {
    return true
  }
  
  // 比较一些关键数据字段
  const keyFields = ['latitude', 'longitude', 'height', 'drone_charge_state', 'environment_temperature']
  for (const field of keyFields) {
    if (oldData.data?.[field] !== newData.data?.[field]) {
      return true
    }
  }
  
  return false
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 格式化存储大小
const formatStorage = (bytes) => {
  if (!bytes) return '--'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  
  return `${size.toFixed(2)} ${units[unitIndex]}`
}

// 获取存储使用率
const getStorageUsage = (storage) => {
  if (!storage || !storage.total) return 0
  return Math.round((storage.used / storage.total) * 100)
}

// 获取存储使用率颜色
const getStorageUsageColor = (storage) => {
  const usage = getStorageUsage(storage)
  if (usage > 90) return '#f56c6c'
  if (usage > 70) return '#e6a23c'
  return '#67c23a'
}

// 获取电池颜色
const getBatteryColor = (percentage) => {
  if (percentage > 50) return '#67c23a'
  if (percentage > 20) return '#e6a23c'
  return '#f56c6c'
}

// 获取充电状态类型
const getChargeStateType = (state) => {
  const stateMap = {
    0: 'success', // 正常
    1: 'warning', // 充电中
    2: 'danger'   // 异常
  }
  return stateMap[state] || 'info'
}

// 获取充电状态文本
const getChargeStateText = (state) => {
  const stateMap = {
    0: '正常',
    1: '充电中',
    2: '异常'
  }
  return stateMap[state] || '未知'
}

// 获取电池存储模式类型
const getBatteryStoreModeType = (mode) => {
  const modeMap = {
    0: 'success', // 正常模式
    1: 'warning', // 存储模式
    2: 'info'     // 其他
  }
  return modeMap[mode] || 'info'
}

// 获取电池存储模式文本
const getBatteryStoreModeText = (mode) => {
  const modeMap = {
    0: '正常模式',
    1: '存储模式',
    2: '其他模式'
  }
  return modeMap[mode] || '未知'
}

// 获取网络类型颜色
const getNetworkTypeColor = (type) => {
  const typeMap = {
    0: 'info',    // 未知
    1: 'success', // WiFi
    2: 'primary', // 4G
    3: 'warning'  // 其他
  }
  return typeMap[type] || 'info'
}

// 获取网络类型文本
const getNetworkTypeText = (type) => {
  const typeMap = {
    0: '未知',
    1: 'WiFi',
    2: '4G',
    3: '其他'
  }
  return typeMap[type] || '未知'
}

// 获取网络质量颜色
const getNetworkQualityColor = (quality) => {
  if (quality >= 4) return '#67c23a'
  if (quality >= 2) return '#e6a23c'
  return '#f56c6c'
}

// 获取网络传输速率文本
const getNetworkRateText = (rate) => {
  const rateMap = {
    0: '未知',
    1: '低',
    2: '中',
    3: '高',
    4: '最高'
  }
  return rateMap[rate] || '未知'
}

// 获取模式代码类型
const getModeCodeType = (code) => {
  const codeMap = {
    0: 'success', // 正常
    1: 'warning', // 警告
    2: 'danger'   // 错误
  }
  return codeMap[code] || 'info'
}

// 获取模式代码文本
const getModeCodeText = (code) => {
  const codeMap = {
    0: '正常',
    1: '警告',
    2: '错误'
  }
  return codeMap[code] || '未知'
}

// 格式化JSON数据
const formatJsonPayload = (payload) => {
  try {
    let jsonString
    if (typeof payload === 'string') {
      const parsed = JSON.parse(payload)
      jsonString = JSON.stringify(parsed, null, 2)
    } else if (typeof payload === 'object' && payload !== null) {
      jsonString = JSON.stringify(payload, null, 2)
    } else {
      return String(payload)
    }
    
    return jsonString
      .replace(/(".*?")\s*:/g, '<span class="json-key">$1</span>:')
      .replace(/:\s*(".*?")/g, ': <span class="json-string">$1</span>')
      .replace(/:\s*(\d+)/g, ': <span class="json-number">$1</span>')
      .replace(/:\s*(true|false)/g, ': <span class="json-boolean">$1</span>')
      .replace(/:\s*(null)/g, ': <span class="json-null">$1</span>')
      .replace(/([{}[\]])/g, '<span class="json-bracket">$1</span>')
      .replace(/(,)/g, '<span class="json-comma">$1</span>')
  } catch (error) {
    return String(payload)
  }
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '--'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 复制原始数据
const copyRawData = async () => {
  if (!osdData.value) return
  
  try {
    const jsonString = JSON.stringify(osdData.value, null, 2)
    await navigator.clipboard.writeText(jsonString)
    copySuccess.value = true
    
    setTimeout(() => {
      copySuccess.value = false
    }, 2000)
    
    ElMessage.success('OSD数据已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败，请手动复制')
  }
}
</script>

<style scoped>
.aircraft-detail {
  height: 100vh;
  background: var(--el-bg-color-page);
}

.detail-header {
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-left h2 {
  margin: 0;
  color: var(--el-text-color-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.detail-main {
  padding: 20px;
  background: var(--el-bg-color-page);
}

.connection-warning {
  margin-bottom: 20px;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
}

.data-container {
  max-width: 1400px;
  margin: 0 auto;
}

.data-grid {
  margin-bottom: 0;
}

.info-card {
  height: 100%;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.info-content {
  padding: 8px 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.info-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.label {
  font-weight: 500;
  color: var(--el-text-color-regular);
  min-width: 80px;
}

.value {
  color: var(--el-text-color-primary);
  font-weight: 600;
}

.raw-data-card {
  margin-top: 20px;
}

.json-viewer {
  background: #1e1e1e;
  border: 1px solid #3c3c3c;
  border-radius: 8px;
  padding: 16px;
  max-height: 400px;
  overflow-y: auto;
}

.json-content {
  margin: 0;
  font-family: 'Fira Code', 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #d4d4d4;
  white-space: pre-wrap;
  word-break: break-word;
}

/* JSON语法高亮 */
.json-key {
  color: #9cdcfe;
  font-weight: 600;
}

.json-string {
  color: #ce9178;
}

.json-number {
  color: #b5cea8;
}

.json-boolean {
  color: #569cd6;
  font-weight: 600;
}

.json-null {
  color: #569cd6;
  font-style: italic;
}

.json-bracket {
  color: #d4d4d4;
  font-weight: bold;
}

.json-comma {
  color: #d4d4d4;
}

.copy-success {
  background-color: var(--el-color-success) !important;
  border-color: var(--el-color-success) !important;
  color: white !important;
}
</style>
