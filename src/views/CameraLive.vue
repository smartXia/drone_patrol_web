<template>
  <div class="camera-live-page">
    <el-container class="page-container">
      <!-- 头部 -->
      <el-header class="page-header">
        <div class="header-content">
          <div class="header-left">
            <el-icon class="header-icon"><VideoCamera /></el-icon>
            <h1>机场直播摄像头</h1>
            <div v-if="getAirportSn() !== '--' || getAircraftSn() !== '--'" class="device-info">
              <el-tag v-if="getAirportSn() !== '--'" size="small" type="success">
                机场: {{ getAirportSn() }}
              </el-tag>
              <el-tag v-if="getAircraftSn() !== '--'" size="small" type="info" style="margin-left: 8px;">
                飞机: {{ getAircraftSn() }}
              </el-tag>
            </div>
          </div>
          <div class="header-right">
            <el-button 
              type="primary" 
              @click="goBack"
              :icon="ArrowLeft"
            >
              返回面板
            </el-button>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main class="page-main">
        <el-row :gutter="16">
          <!-- 摄像头直播区域 -->
          <el-col :span="24">
            <CameraLiveStream 
              ref="cameraStreamRef" 
              :airport-sn="getAirportSn()"
              :aircraft-sn="getAircraftSn()"
            />
          </el-col>
          
        </el-row>
      </el-main>
    </el-container>

  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  VideoCamera, 
  ArrowLeft
} from '@element-plus/icons-vue'
import CameraLiveStream from '@/components/CameraLiveStream.vue'
import { useDeviceStore } from '@/stores/device'

const router = useRouter()
const route = useRoute()
const deviceStore = useDeviceStore()

// 响应式数据
const cameraStreamRef = ref(null)
const streamDuration = ref(0)
const isStreaming = ref(false)
const currentResolution = ref('')
const currentFps = ref('')

// 获取设备SN
const getAirportSn = () => {
  return route.params.airportSn || deviceStore.currentDevice?.airport_sn || deviceStore.currentDeviceSn || '--'
}

const getAircraftSn = () => {
  return route.params.aircraftSn || route.params.deviceSn || '--'
}

// 计算属性 - 获取设备列表
const deviceList = computed(() => deviceStore.deviceList)
const airportDevices = computed(() => 
  deviceList.value.filter(device => device.type === 'airport' || device.airport_sn)
)
const aircraftDevices = computed(() => 
  deviceList.value.filter(device => device.type === 'aircraft' && device.airport_sn)
)


// 方法
const goBack = () => {
  router.push('/')
}



const formatDuration = (seconds) => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  
  if (hours > 0) {
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  } else {
    return `${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
}

// 生命周期
onMounted(() => {
  // 加载保存的摄像头列表
  const savedCameras = localStorage.getItem('camera-list')
  if (savedCameras) {
    try {
      cameraList.value = JSON.parse(savedCameras)
    } catch (error) {
      console.error('加载摄像头列表失败:', error)
    }
  }
  
  // 设置默认摄像头
  if (cameraList.value.length > 0) {
    currentCameraId.value = cameraList.value[0].id
    switchCamera(cameraList.value[0])
  }
  
  // 保存摄像头列表
  localStorage.setItem('camera-list', JSON.stringify(cameraList.value))
})
</script>

<style scoped>
.camera-live-page {
  height: 100vh;
  background-color: #f5f7fa;
}

.page-container {
  height: 100%;
}

.page-header {
  background-color: white;
  border-bottom: 1px solid #e4e7ed;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon {
  font-size: 24px;
  color: #409eff;
}

.header-left h1 {
  margin: 0;
  font-size: 20px;
  font-weight: 500;
  color: #303133;
}

.page-main {
  padding: 20px;
  height: calc(100vh - 60px);
  overflow-y: auto;
}

.camera-list-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.camera-list {
  max-height: 300px;
  overflow-y: auto;
}

.camera-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
}

.camera-item:hover {
  border-color: #409eff;
  background-color: #f0f9ff;
}

.camera-item.active {
  border-color: #409eff;
  background-color: #e6f7ff;
}

.camera-info {
  flex: 1;
}

.camera-name {
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.camera-url {
  color: #909399;
  font-size: 12px;
  font-family: monospace;
  margin-bottom: 4px;
}

.camera-status {
  margin-top: 4px;
}

.camera-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s;
}

.camera-item:hover .camera-actions {
  opacity: 1;
}

.delete-btn {
  color: #f56c6c;
}

.empty-camera-list {
  padding: 20px;
  text-align: center;
}

.stats-card {
  margin-bottom: 16px;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.stat-item:last-child {
  border-bottom: none;
}

.stat-label {
  color: #606266;
  font-size: 14px;
}

.stat-value {
  color: #303133;
  font-weight: 500;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .page-main .el-col {
    margin-bottom: 16px;
  }
}

@media (max-width: 768px) {
  .page-header {
    padding: 0 16px;
  }
  
  .page-main {
    padding: 16px;
  }
  
  .header-left h1 {
    font-size: 18px;
  }
}
</style>
