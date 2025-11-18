<template>
  <div class="remote-upgrade-container">
    <!-- 固件属性上报 -->
    <el-card class="upgrade-card" shadow="hover" style="margin-bottom: 20px;">
      <template #header>
        <div class="card-header">
          <el-icon><Setting /></el-icon>
          <span>固件属性上报</span>
          <el-tag type="info" size="small" style="margin-left: auto;">
            Topic: thing/product/{gateway_sn}/osd
          </el-tag>
        </div>
      </template>
      
      <div class="firmware-info">
        <!-- 提示信息 -->
        <div class="info-tip" style="margin-bottom: 15px; padding: 12px; background: var(--el-color-info-light-9); border-radius: 6px; border-left: 4px solid var(--el-color-info);">
          <div style="display: flex; align-items: flex-start; gap: 8px; color: var(--el-color-info); font-size: 13px;">
            <el-icon style="margin-top: 2px;"><InfoFilled /></el-icon>
            <div>
              <div style="font-weight: 600; margin-bottom: 4px;">固件属性上报说明：</div>
              <div style="line-height: 1.5;">
                • 固件信息只有在设备变动时才会显示<br/>
                • 主要关注<strong>固件一致性状态</strong>，用于判断是否需要升级<br/>
                • 飞机开机状态表示机场停机坪上的飞行器是否开机<br/>
                • 对频状态表示飞行器是否与机场完成对频连接<br/>
                • 固件版本信息作为参考，实际升级以一致性状态为准
              </div>
            </div>
          </div>
        </div>
        
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="info-section">
              <h4>飞机固件信息</h4>
              <div class="info-item">
                <span class="label">固件版本:</span>
                <span class="value">{{ aircraftFirmwareVersion || '--' }}</span>
              </div>
              <div class="info-item" style="background: var(--el-color-warning-light-9); padding: 8px; border-radius: 4px; border: 1px solid var(--el-color-warning-light-7);">
                <span class="label" style="font-weight: 600; color: var(--el-color-warning-dark-2);">固件一致性:</span>
                <el-tag :type="aircraftCompatibleStatus ? 'warning' : 'success'" size="small" style="font-weight: 600;">
                  {{ aircraftCompatibleStatus ? '需要升级' : '正常' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">设备型号:</span>
                <span class="value">{{ aircraftModel || '--' }}</span>
              </div>
              <div class="info-item">
                <span class="label">设备序列号:</span>
                <span class="value">{{ aircraftSn || '--' }}</span>
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="info-section">
              <h4>机场固件信息</h4>
              <div class="info-item">
                <span class="label">固件版本:</span>
                <span class="value">{{ dockFirmwareVersion || '--' }}</span>
              </div>
              <div class="info-item" style="background: var(--el-color-warning-light-9); padding: 8px; border-radius: 4px; border: 1px solid var(--el-color-warning-light-7);">
                <span class="label" style="font-weight: 600; color: var(--el-color-warning-dark-2);">固件一致性:</span>
                <el-tag :type="dockCompatibleStatus ? 'warning' : 'success'" size="small" style="font-weight: 600;">
                  {{ dockCompatibleStatus ? '需要升级' : '正常' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">设备型号:</span>
                <span class="value">{{ dockModel || '--' }}</span>
              </div>
              <div class="info-item">
                <span class="label">设备序列号:</span>
                <span class="value">{{ dockSn || '--' }}</span>
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="info-section">
              <h4>系统状态</h4>
              <div class="info-item">
                <span class="label">飞机开机状态:</span>
                <el-tag :type="aircraftOnline ? 'success' : 'danger'" size="small">
                  {{ aircraftOnline ? '开机' : '关机' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">机场在线:</span>
                <el-tag :type="dockOnline ? 'success' : 'danger'" size="small">
                  {{ dockOnline ? '在线' : '离线' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">对频状态:</span>
                <el-tag :type="devicePaired ? 'success' : 'warning'" size="small">
                  {{ devicePaired ? '已对频' : '未对频' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">最后更新:</span>
                <span class="value">{{ lastUpdateTime || '--' }}</span>
              </div>
              <div class="info-item" style="margin-top: 10px; padding: 8px; background: var(--el-color-primary-light-9); border-radius: 4px; border: 1px solid var(--el-color-primary-light-7);">
                <span class="label" style="font-weight: 600; color: var(--el-color-primary-dark-2);">升级建议:</span>
                <span class="value" style="color: var(--el-color-primary-dark-2);">
                  {{ (aircraftCompatibleStatus || dockCompatibleStatus) ? '检测到设备需要升级' : '所有设备固件正常' }}
                </span>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <!-- 固件升级操作 -->
    <el-card class="upgrade-card" shadow="hover" style="margin-bottom: 20px;">
      <template #header>
        <div class="card-header">
          <el-icon><Upload /></el-icon>
          <span>固件升级操作</span>
          <el-tag type="primary" size="small" style="margin-left: auto;">
            Topic: thing/product/{gateway_sn}/services | Method: ota_create
          </el-tag>
        </div>
      </template>
      
      <div class="upgrade-actions">
        <el-form :model="upgradeForm" label-width="120px" size="small">
          <el-form-item label="设备类型">
            <el-select v-model="upgradeForm.deviceType" placeholder="选择设备类型" style="width: 200px;">
              <el-option label="飞机固件" value="aircraft" />
              <el-option label="机场固件" value="dock" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="升级类型">
            <el-select v-model="upgradeForm.type" placeholder="选择升级类型" style="width: 200px;">
              <el-option label="固件升级" value="firmware" />
              <el-option label="应用升级" value="application" />
              <el-option label="系统升级" value="system" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="固件文件URL">
            <el-input 
              v-model="upgradeForm.fileUrl" 
              placeholder="固件文件下载地址"
              style="width: 400px;"
            />
          </el-form-item>
          
          <el-form-item label="固件文件名">
            <el-input 
              v-model="upgradeForm.fileName" 
              placeholder="固件文件名称"
              style="width: 300px;"
            />
          </el-form-item>
          
          <el-form-item label="固件版本">
            <el-input 
              v-model="upgradeForm.productVersion" 
              placeholder="固件版本号"
              style="width: 200px;"
            />
          </el-form-item>
          
          <el-form-item label="文件大小">
            <el-input 
              v-model="upgradeForm.fileSize" 
              placeholder="文件大小(字节)"
              style="width: 200px;"
            />
          </el-form-item>
          
          <el-form-item label="MD5校验">
            <el-input 
              v-model="upgradeForm.md5" 
              placeholder="文件MD5值"
              style="width: 300px;"
            />
          </el-form-item>
          
          <el-form-item label="升级类型">
            <el-select v-model="upgradeForm.firmwareUpgradeType" placeholder="选择升级类型" style="width: 200px;">
              <el-option label="一致性升级" :value="2" />
              <el-option label="普通升级" :value="3" />
            </el-select>
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              :loading="upgradeSending"
              @click="startUpgrade"
              :disabled="!upgradeForm.firmwareUrl"
            >
              开始升级
            </el-button>
            <el-button @click="resetUpgradeForm">重置</el-button>
            <el-button type="warning" @click="checkUpgradeStatus">检查状态</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <!-- 升级状态监控 -->
    <el-card class="upgrade-card" shadow="hover" style="margin-bottom: 20px;">
      <template #header>
        <div class="card-header">
          <el-icon><Monitor /></el-icon>
          <span>升级状态监控</span>
          <el-tag type="info" size="small" style="margin-left: auto; margin-right: 10px;">
            Topic: thing/product/{gateway_sn}/services_reply | Method: ota_create
          </el-tag>
          <el-button 
            type="primary" 
            size="small" 
            :loading="statusSubscribing"
            @click="subscribeUpgradeStatus"
          >
            手动订阅状态
          </el-button>
        </div>
      </template>
      
      <div class="status-monitor">
        <div v-if="props.upgradeStatus" class="status-info">
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="status-item">
                <span class="label">升级状态:</span>
                <el-tag :type="getStatusType(props.upgradeStatus.status)" size="small">
                  {{ getStatusText(props.upgradeStatus.status) }}
                </el-tag>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="status-item">
                <span class="label">进度:</span>
                <el-progress 
                  :percentage="props.upgradeStatus.progress || 0" 
                  :status="getProgressStatus(props.upgradeStatus.status)"
                  style="width: 150px;"
                />
              </div>
            </el-col>
            <el-col :span="8">
              <div class="status-item">
                <span class="label">剩余时间:</span>
                <span class="value">{{ props.upgradeStatus.remainingTime || '--' }}</span>
              </div>
            </el-col>
          </el-row>
          
          <div v-if="props.upgradeStatus.error" class="error-info" style="margin-top: 15px;">
            <el-alert 
              :title="props.upgradeStatus.error" 
              type="error" 
              :closable="false"
              show-icon
            />
          </div>
        </div>
        
        <div v-else class="no-status">
          <el-empty description="等待升级状态数据..." :image-size="80" />
        </div>
      </div>
    </el-card>

    <!-- 升级历史记录 -->
    <el-card class="upgrade-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon><Document /></el-icon>
          <span>升级历史记录</span>
          <el-button type="primary" size="small" @click="refreshHistory" style="margin-left: auto;">
            刷新记录
          </el-button>
        </div>
      </template>
      
      <div class="upgrade-history">
        <el-table :data="upgradeHistory" size="small" height="300" border>
          <el-table-column prop="timestamp" label="时间" width="180">
            <template #default="scope">
              {{ formatTime(scope.row.timestamp) }}
            </template>
          </el-table-column>
          <el-table-column prop="deviceType" label="设备类型" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.deviceType === 'aircraft' ? 'primary' : 'info'" size="small">
                {{ scope.row.deviceType === 'aircraft' ? '飞机' : '机场' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="升级类型" width="100" />
          <el-table-column prop="version" label="版本" width="120" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)" size="small">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="result" label="结果" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.result === 0 ? 'success' : 'danger'" size="small">
                {{ scope.row.result === 0 ? '成功' : '失败' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="error" label="错误信息" />
        </el-table>
      </div>
    </el-card>

    <!-- 调试信息 -->
    <div class="debug-info" style="margin-top: 20px;">
      <el-collapse>
        <el-collapse-item title="调试信息" name="debug">
          <div style="padding: 10px; background: var(--el-bg-color-page); border-radius: 4px;">
            <div style="margin-bottom: 10px;">
              <strong>固件数据状态:</strong> {{ osdData ? '已接收' : '未接收' }}
            </div>
            <div style="margin-bottom: 10px;">
              <strong>升级状态数据:</strong> {{ props.upgradeStatus ? '已接收' : '未接收' }}
            </div>
            <div style="margin-bottom: 10px;">
              <strong>历史记录数量:</strong> {{ upgradeHistory.length }}
            </div>
            <div v-if="props.upgradeStatusPayload" style="margin-top: 10px;">
              <div style="margin-bottom: 5px; color: var(--el-text-color-regular); font-size: 13px;">
                最新升级状态数据：
              </div>
              <pre style="background: var(--el-bg-color); padding: 10px; border-radius: 4px; font-size: 11px; max-height: 200px; overflow-y: auto;">{{ formatUpgradeStatusPayload() }}</pre>
            </div>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Setting, Upload, Monitor, Document, InfoFilled } from '@element-plus/icons-vue'

// Props
const props = defineProps({
  deviceSn: {
    type: String,
    required: true
  },
  osdData: {
    type: Object,
    default: null
  },
  upgradeStatus: {
    type: Object,
    default: null
  },
  upgradeStatusPayload: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['subscribeUpgrade', 'startUpgrade', 'checkStatus'])

// 响应式数据
const upgradeSending = ref(false)
const statusSubscribing = ref(false)
const upgradeHistory = ref([])

// 升级表单
const upgradeForm = ref({
  deviceType: 'aircraft',
  type: 'firmware',
  fileUrl: '',
  fileName: '',
  productVersion: '',
  fileSize: '',
  md5: '',
  firmwareUpgradeType: 2
})

// 计算属性 - 从OSD数据中提取固件信息
// 飞机固件信息
const aircraftFirmwareVersion = computed(() => {
  return props.osdData?.data?.aircraft?.firmware_version || props.osdData?.data?.firmware_version || '--'
})

const aircraftCompatibleStatus = computed(() => {
  return props.osdData?.data?.aircraft?.compatible_status || props.osdData?.data?.compatible_status || false
})

const aircraftModel = computed(() => {
  return props.osdData?.data?.sub_device?.device_model_key || '--'
})

const aircraftSn = computed(() => {
  return props.osdData?.data?.sub_device?.device_sn || '--'
})

const aircraftOnline = computed(() => {
  const status = props.osdData?.data?.sub_device?.device_online_status
  return status === 1 // 1表示开机，0表示关机
})

// 机场固件信息
const dockFirmwareVersion = computed(() => {
  return props.osdData?.data?.dock?.firmware_version || '--'
})

const dockCompatibleStatus = computed(() => {
  return props.osdData?.data?.dock?.compatible_status || false
})

const dockModel = computed(() => {
  return props.osdData?.data?.dock?.device_model_key || '--'
})

const dockSn = computed(() => {
  return props.osdData?.data?.dock?.device_sn || props.deviceSn || '--'
})

const dockOnline = computed(() => {
  return props.osdData?.data?.dock?.device_online_status || true // 机场通常是在线的
})

// 系统状态
const devicePaired = computed(() => {
  const paired = props.osdData?.data?.sub_device?.device_paired
  return paired === 1 // 1表示已对频，0表示未对频
})

const lastUpdateTime = computed(() => {
  if (props.osdData?.timestamp) {
    return new Date(props.osdData.timestamp).toLocaleString()
  }
  return '--'
})

// 方法
const selectFirmwareFile = () => {
  ElMessage.info('文件选择功能待实现')
}

const loadPresetParams = () => {
  // 根据设备类型加载预设参数
  if (upgradeForm.value.deviceType === 'aircraft') {
    upgradeForm.value.fileUrl = 'https://s3.com/wm245_1.00.223.zip'
    upgradeForm.value.fileName = 'wm245_1.00.223.zip'
    upgradeForm.value.productVersion = '1.00.223'
    upgradeForm.value.fileSize = '653467234'
    upgradeForm.value.md5 = 'abcdefabcdefabcdef'
    upgradeForm.value.firmwareUpgradeType = 2
  } else {
    upgradeForm.value.fileUrl = 'https://s3.com/dock_1.00.223.zip'
    upgradeForm.value.fileName = 'dock_1.00.223.zip'
    upgradeForm.value.productVersion = '1.00.223'
    upgradeForm.value.fileSize = '123456789'
    upgradeForm.value.md5 = 'fedcbafedcbafedcba'
    upgradeForm.value.firmwareUpgradeType = 3
  }
  ElMessage.success('预设参数已加载')
}

const startUpgrade = async () => {
  // 验证必填字段
  if (!upgradeForm.value.fileUrl) {
    ElMessage.warning('请输入固件文件URL')
    return
  }

  if (!upgradeForm.value.fileName) {
    ElMessage.warning('请输入固件文件名')
    return
  }

  if (!upgradeForm.value.productVersion) {
    ElMessage.warning('请输入固件版本')
    return
  }

  if (!upgradeForm.value.deviceType) {
    ElMessage.warning('请选择设备类型')
    return
  }

  upgradeSending.value = true
  
  try {
    // 构建设备信息
    const deviceInfo = {
      sn: upgradeForm.value.deviceType === 'aircraft' ? props.osdData?.data?.sub_device?.device_sn : props.deviceSn,
      product_version: upgradeForm.value.productVersion,
      firmware_upgrade_type: upgradeForm.value.firmwareUpgradeType
    }

    // 如果是飞机且有文件信息，添加文件相关字段
    if (upgradeForm.value.deviceType === 'aircraft' && upgradeForm.value.fileUrl) {
      deviceInfo.file_url = upgradeForm.value.fileUrl
      deviceInfo.file_name = upgradeForm.value.fileName
      if (upgradeForm.value.md5) deviceInfo.md5 = upgradeForm.value.md5
      if (upgradeForm.value.fileSize) deviceInfo.file_size = parseInt(upgradeForm.value.fileSize)
    }

    const upgradePayload = {
      method: 'ota_create',
      data: {
        devices: [deviceInfo]
      }
    }

    emit('startUpgrade', upgradePayload)
    ElMessage.success(`${upgradeForm.value.deviceType === 'aircraft' ? '飞机' : '机场'}OTA升级请求已发送`)
  } catch (e) {
    ElMessage.error('升级请求发送失败: ' + e.message)
  } finally {
    upgradeSending.value = false
  }
}

const resetUpgradeForm = () => {
  upgradeForm.value = {
    deviceType: 'aircraft',
    type: 'firmware',
    fileUrl: '',
    fileName: '',
    productVersion: '',
    fileSize: '',
    md5: '',
    firmwareUpgradeType: 2
  }
}

const checkUpgradeStatus = () => {
  emit('checkStatus')
  ElMessage.info('状态检查请求已发送')
}

const subscribeUpgradeStatus = () => {
  statusSubscribing.value = true
  emit('subscribeUpgrade')
  
  setTimeout(() => {
    statusSubscribing.value = false
  }, 2000)
}

const refreshHistory = () => {
  // 模拟历史数据
  upgradeHistory.value = [
    {
      timestamp: Date.now() - 3600000,
      deviceType: 'aircraft',
      type: '固件升级',
      version: '1.0.0',
      status: 'completed',
      result: 0,
      error: ''
    },
    {
      timestamp: Date.now() - 7200000,
      deviceType: 'dock',
      type: '应用升级',
      version: '2.1.0',
      status: 'failed',
      result: 1,
      error: '升级超时'
    },
    {
      timestamp: Date.now() - 10800000,
      deviceType: 'aircraft',
      type: '系统升级',
      version: '1.5.0',
      status: 'completed',
      result: 0,
      error: ''
    }
  ]
  ElMessage.success('历史记录已刷新')
}

const getStatusType = (status) => {
  const statusMap = {
    'sent': 'info',
    'in_progress': 'primary',
    'ok': 'success',
    'failed': 'danger',
    'canceled': 'warning',
    'paused': 'warning',
    'rejected': 'danger',
    'timeout': 'danger',
    'pending': 'info',
    'downloading': 'primary',
    'installing': 'warning',
    'completed': 'success',
    'cancelled': 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    'sent': '已下发',
    'in_progress': '执行中',
    'ok': '执行成功',
    'failed': '失败',
    'canceled': '取消或终止',
    'paused': '暂停',
    'rejected': '拒绝',
    'timeout': '超时',
    'pending': '等待中',
    'downloading': '下载中',
    'installing': '安装中',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return statusMap[status] || '未知'
}

const getProgressStatus = (status) => {
  if (status === 'failed' || status === 'rejected' || status === 'timeout') return 'exception'
  if (status === 'completed' || status === 'ok') return 'success'
  return ''
}

const formatTime = (timestamp) => {
  if (!timestamp) return '--'
  return new Date(timestamp).toLocaleString()
}

const formatUpgradeStatusPayload = () => {
  if (!props.upgradeStatusPayload) return ''
  
  try {
    const data = typeof props.upgradeStatusPayload === 'string'
      ? JSON.parse(props.upgradeStatusPayload)
      : props.upgradeStatusPayload
    
    return JSON.stringify(data, null, 2)
  } catch (e) {
    return String(props.upgradeStatusPayload)
  }
}

// 组件挂载时初始化
onMounted(() => {
  refreshHistory()
})
</script>

<style scoped>
.remote-upgrade-container {
  padding: 8px 0;
}

.upgrade-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.firmware-info {
  padding: 10px 0;
}

.info-section {
  padding: 10px;
  background: var(--el-bg-color-page);
  border-radius: 4px;
  border: 1px solid var(--el-border-color-light);
}

.info-section h4 {
  margin: 0 0 15px 0;
  color: var(--el-text-color-primary);
  font-size: 14px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 13px;
}

.info-item:last-child {
  margin-bottom: 0;
}

.label {
  color: var(--el-text-color-regular);
  font-weight: 500;
}

.value {
  color: var(--el-text-color-primary);
  font-family: monospace;
}

.status-monitor {
  padding: 10px 0;
}

.status-info {
  padding: 15px;
  background: var(--el-bg-color-page);
  border-radius: 4px;
  border: 1px solid var(--el-border-color-light);
}

.status-item {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.status-item:last-child {
  margin-bottom: 0;
}

.no-status {
  text-align: center;
  padding: 40px 0;
}

.upgrade-history {
  padding: 10px 0;
}

.debug-info {
  margin-top: 20px;
}
</style>


