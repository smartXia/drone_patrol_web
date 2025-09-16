<template>
  <el-card class="device-status-card" shadow="hover">
    <template #header>
      <div class="device-header">
        <div class="device-info">
          <i :class="deviceIcon" class="device-icon"></i>
          <span class="device-id">{{ deviceId }}</span>
        </div>
        <el-tag 
          :type="statusType" 
          size="small"
          :class="{ 'pulse': isOnline }"
        >
          {{ statusText }}
        </el-tag>
      </div>
    </template>
    
    <div class="device-content">
      <div class="device-meta">
        <div class="meta-item">
          <span class="label">最后更新:</span>
          <span class="value">{{ formatTime(device.lastUpdate) }}</span>
        </div>
        <div class="meta-item">
          <span class="label">主题:</span>
          <span class="value topic">{{ device.topic }}</span>
        </div>
        <div v-if="device.deviceType" class="meta-item">
          <span class="label">设备类型:</span>
          <span class="value">{{ device.deviceType }}</span>
        </div>
      </div>
      
      <div class="device-data">
        <div class="data-header">
          <h4>设备数据</h4>
          <el-button 
            type="text" 
            size="small"
            @click="toggleDataView"
          >
            {{ showRawData ? '格式化' : '原始' }}
          </el-button>
        </div>
        
        <div class="data-content">
          <div v-if="showRawData" class="raw-data">
            <pre>{{ JSON.stringify(device, null, 2) }}</pre>
          </div>
          <div v-else class="formatted-data">
            <div 
              v-for="(value, key) in filteredDeviceData" 
              :key="key"
              class="data-item"
            >
              <span class="data-key">{{ key }}:</span>
              <span class="data-value">{{ formatValue(value) }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="device-actions">
        <el-button 
          type="primary" 
          size="small"
          @click="refreshDevice"
        >
          刷新
        </el-button>
        <el-button 
          type="info" 
          size="small"
          @click="viewHistory"
        >
          历史记录
        </el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  formatDeviceStatus, 
  getDeviceStatusType, 
  getDeviceTypeIcon,
  isDeviceOnline 
} from '@/utils/mqtt'
import { formatTime } from '@/utils/time'

const props = defineProps({
  deviceId: {
    type: String,
    required: true
  },
  device: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['refresh', 'viewHistory'])

const showRawData = ref(false)

// 计算属性
const statusText = computed(() => {
  return formatDeviceStatus(props.device)
})

const statusType = computed(() => {
  return getDeviceStatusType(props.device)
})

const isOnline = computed(() => {
  return isDeviceOnline(props.device)
})

const deviceIcon = computed(() => {
  return getDeviceTypeIcon(props.device.deviceType)
})

const filteredDeviceData = computed(() => {
  const { lastUpdate, topic, deviceType, ...data } = props.device
  return data
})

// 方法
const toggleDataView = () => {
  showRawData.value = !showRawData.value
}

const formatValue = (value) => {
  if (value === null || value === undefined) return '无'
  if (typeof value === 'boolean') return value ? '是' : '否'
  if (typeof value === 'object') return JSON.stringify(value)
  return String(value)
}

const refreshDevice = () => {
  emit('refresh', props.deviceId)
  ElMessage.success('刷新请求已发送')
}

const viewHistory = () => {
  emit('viewHistory', props.deviceId)
}
</script>

<style scoped>
.device-status-card {
  height: 100%;
  transition: all 0.3s ease;
}

.device-status-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.device-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.device-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.device-icon {
  font-size: 18px;
  color: #409eff;
}

.device-id {
  font-weight: bold;
  color: #303133;
}

.pulse {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
  100% {
    opacity: 1;
  }
}

.device-content {
  height: calc(100% - 60px);
  display: flex;
  flex-direction: column;
}

.device-meta {
  margin-bottom: 15px;
}

.meta-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.meta-item .label {
  color: #606266;
  font-weight: 500;
}

.meta-item .value {
  color: #303133;
}

.meta-item .topic {
  color: #409eff;
  font-family: monospace;
  font-size: 12px;
}

.device-data {
  flex: 1;
  margin-bottom: 15px;
}

.data-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.data-header h4 {
  margin: 0;
  color: #606266;
  font-size: 14px;
}

.data-content {
  height: 200px;
  overflow-y: auto;
}

.raw-data pre {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  font-size: 12px;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}

.formatted-data {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.data-item {
  display: flex;
  justify-content: space-between;
  padding: 8px;
  background-color: #f5f7fa;
  border-radius: 4px;
  font-size: 13px;
}

.data-key {
  color: #606266;
  font-weight: 500;
}

.data-value {
  color: #303133;
  max-width: 60%;
  word-break: break-all;
  text-align: right;
}

.device-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}
</style>
