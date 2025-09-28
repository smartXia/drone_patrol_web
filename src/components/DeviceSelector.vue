<template>
  <div class="device-selector">
    <el-card shadow="hover" class="device-card">
      <template #header>
        <div class="card-header">
          <span>设备管理</span>
          <el-button 
            type="primary" 
            size="small" 
            @click="showAddDialog = true"
            :icon="Plus"
          >
            添加设备
          </el-button>
        </div>
      </template>

      <!-- 当前设备信息 -->
      <div class="current-device">
        <div class="device-info" v-if="deviceStore.currentDevice">
          <el-tag 
            :type="deviceStore.currentDevice.status === 'online' ? 'success' : 'info'"
            size="small"
            class="status-tag"
          >
            {{ deviceStore.currentDevice.status === 'online' ? '在线' : '离线' }}
          </el-tag>
          <span class="device-name">{{ deviceStore.currentDevice.name }}</span>
          <span class="device-sn">{{ deviceStore.currentDeviceSn }}</span>
          <span v-if="deviceStore.currentDevice.airport_sn" class="airport-sn">
            机场: {{ deviceStore.currentDevice.airport_sn }}
          </span>
        </div>
        <div class="device-info" v-else>
          <el-tag type="info" size="small" class="status-tag">
            未选择
          </el-tag>
          <span class="device-name">请选择设备</span>
        </div>
      </div>

      <!-- 设备列表 -->
      <div class="device-list">
        <div 
          v-for="device in deviceStore.deviceList" 
          :key="device.id"
          class="device-item"
          :class="{ 
            'active': device.sn === deviceStore.currentDeviceSn,
            'online': device.status === 'online',
            'offline': device.status === 'offline'
          }"
          @click="selectDevice(device)"
        >
          <div class="device-main">
            <el-tag 
              :type="device.status === 'online' ? 'success' : 'info'"
              size="small"
              class="status-indicator"
            >
              {{ device.status === 'online' ? '在线' : '离线' }}
            </el-tag>
            <div class="device-details">
              <div class="device-name">{{ device.name }}</div>
              <div class="device-sn">{{ device.sn }}</div>
              <div class="device-type">{{ getDeviceTypeName(device.type) }}</div>
              <div v-if="device.airport_sn" class="airport-sn">机场: {{ device.airport_sn }}</div>
            </div>
          </div>
          <div class="device-actions">
            <el-button 
              type="text" 
              size="small" 
              @click.stop="editDevice(device)"
              :icon="Edit"
            >
              编辑
            </el-button>
            <el-button 
              type="text" 
              size="small" 
              @click.stop="removeDevice(device)"
              :icon="Delete"
              class="delete-btn"
            >
              删除
            </el-button>
          </div>
        </div>
        <div v-if="deviceStore.deviceList.length === 0" class="empty-device-list">
          <el-empty description="暂无设备，请添加设备" />
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="device-actions-footer">
        <el-button 
          type="primary" 
          @click="showAddDialog = true"
          :icon="Plus"
        >
          添加设备
        </el-button>
        <el-button 
          @click="refreshDevices"
          :icon="Refresh"
        >
          刷新
        </el-button>
        <el-button 
          @click="removeDefaultDevices"
          :icon="Delete"
          type="warning"
        >
          删除默认设备
        </el-button>
        <el-button 
          @click="clearAllDevices"
          :icon="Delete"
          type="danger"
        >
          清空所有
        </el-button>
      </div>
    </el-card>

    <!-- 添加/编辑设备对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="editingDevice ? `编辑设备 - ${editingDevice.name}` : '添加设备'"
      width="500px"
      @close="resetForm"
    >
      <el-form :model="deviceForm" :rules="deviceRules" ref="deviceFormRef" label-width="100px">
        <el-form-item label="设备名称" prop="name">
          <el-input v-model="deviceForm.name" placeholder="请输入设备名称" />
        </el-form-item>
        <el-form-item label="设备序列号" prop="sn">
          <el-input v-model="deviceForm.sn" placeholder="请输入设备序列号" />
        </el-form-item>
        <el-form-item label="设备类型" prop="type">
          <el-select v-model="deviceForm.type" placeholder="请选择设备类型">
            <el-option label="无人机" value="drone" />
            <el-option label="网关" value="gateway" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="设备状态" prop="status">
          <el-select v-model="deviceForm.status" placeholder="请选择设备状态">
            <el-option label="在线" value="online" />
            <el-option label="离线" value="offline" />
          </el-select>
        </el-form-item>
        <el-form-item label="机场序列号" prop="airport_sn">
          <el-input v-model="deviceForm.airport_sn" placeholder="请输入机场序列号（可选）" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveDevice">
          {{ editingDevice ? '更新' : '添加' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Refresh } from '@element-plus/icons-vue'
import { useDeviceStore } from '../stores/device'

const deviceStore = useDeviceStore()

// 响应式数据
const showAddDialog = ref(false)
const editingDevice = ref(null)
const deviceFormRef = ref()

// 设备表单
const deviceForm = reactive({
  name: '',
  sn: '',
  type: 'drone',
  status: 'offline',
  airport_sn: ''
})

// 表单验证规则
const deviceRules = {
  name: [
    { required: true, message: '请输入设备名称', trigger: 'blur' }
  ],
  sn: [
    { required: true, message: '请输入设备序列号', trigger: 'blur' },
    { min: 3, message: '设备序列号至少3个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择设备类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择设备状态', trigger: 'change' }
  ]
}

// 方法
const selectDevice = async (device) => {
  try {
    await deviceStore.setCurrentDevice(device.sn)
    ElMessage.success(`已切换到设备: ${device.name}`)
  } catch (error) {
    ElMessage.error(`切换设备失败: ${error.message}`)
  }
}

const editDevice = (device) => {
  editingDevice.value = device
  Object.assign(deviceForm, {
    name: device.name,
    sn: device.sn,
    type: device.type,
    status: device.status,
    airport_sn: device.airport_sn || ''
  })
  showAddDialog.value = true
}

const removeDevice = async (device) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除设备 "${device.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deviceStore.removeDevice(device.sn)
    ElMessage.success('设备删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`删除设备失败: ${error.message}`)
    }
  }
}

const saveDevice = async () => {
  try {
    await deviceFormRef.value.validate()
    
    if (editingDevice.value) {
      // 编辑设备
      await deviceStore.updateDevice(editingDevice.value.id, deviceForm)
      ElMessage.success('设备更新成功')
    } else {
      // 添加设备
      await deviceStore.addDevice(deviceForm)
      ElMessage.success('设备添加成功')
    }
    
    showAddDialog.value = false
    resetForm()
  } catch (error) {
    if (error.message) {
      ElMessage.error(`操作失败: ${error.message}`)
    } else {
      console.error('表单验证失败:', error)
    }
  }
}

const resetForm = () => {
  editingDevice.value = null
  Object.assign(deviceForm, {
    name: '',
    sn: '',
    type: 'drone',
    status: 'offline',
    airport_sn: ''
  })
  if (deviceFormRef.value) {
    deviceFormRef.value.resetFields()
  }
}

const refreshDevices = async () => {
  try {
    await deviceStore.loadDevices()
    ElMessage.success('设备列表已刷新')
  } catch (error) {
    ElMessage.error('刷新设备列表失败')
  }
}


const removeDefaultDevices = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要删除所有默认设备吗？',
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deviceStore.removeDefaultDevices()
    ElMessage.success('默认设备已删除')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`删除默认设备失败: ${error.message}`)
    }
  }
}

const clearAllDevices = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清空所有设备吗？此操作不可恢复！',
      '确认清空',
      {
        confirmButtonText: '确定清空',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )
    
    await deviceStore.clearAllDevices()
    ElMessage.success('所有设备已清空')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`清空设备失败: ${error.message}`)
    }
  }
}

const getDeviceTypeName = (type) => {
  const typeMap = {
    drone: '无人机',
    gateway: '网关',
    other: '其他'
  }
  return typeMap[type] || '未知'
}

// 生命周期
onMounted(async () => {
  await deviceStore.loadDevices()
  // 自动清理默认设备
  await deviceStore.autoCleanDefaultDevices()
})
</script>

<style scoped>
.device-selector {
  width: 100%;
}

.device-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.current-device {
  margin-bottom: 16px;
  padding: 12px;
  background-color: #f5f7fa;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
}

.device-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-tag {
  margin-right: 8px;
}

.device-name {
  font-weight: 500;
  color: #303133;
}

.device-sn {
  color: #909399;
  font-size: 12px;
  font-family: monospace;
}

.airport-sn {
  color: #67c23a;
  font-size: 12px;
  font-weight: 500;
  margin-left: 8px;
}

.device-list {
  max-height: 300px;
  overflow-y: auto;
}

.device-item {
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

.device-item:hover {
  border-color: #409eff;
  background-color: #f0f9ff;
}

.device-item.active {
  border-color: #409eff;
  background-color: #e6f7ff;
}

.device-item.offline {
  opacity: 0.6;
}

.device-main {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.status-indicator {
  flex-shrink: 0;
}

.device-details {
  flex: 1;
}

.device-details .device-name {
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.device-details .device-sn {
  color: #909399;
  font-size: 12px;
  font-family: monospace;
  margin-bottom: 2px;
}

.device-details .device-type {
  color: #606266;
  font-size: 12px;
}

.device-details .airport-sn {
  color: #67c23a;
  font-size: 12px;
  font-weight: 500;
  margin-top: 2px;
}

.device-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s;
}

.device-item:hover .device-actions {
  opacity: 1;
}

.delete-btn {
  color: #f56c6c;
}

.empty-device-list {
  padding: 20px;
  text-align: center;
}

.device-actions-footer {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e4e7ed;
}
</style>
