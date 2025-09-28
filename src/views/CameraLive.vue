<template>
  <div class="camera-live-page">
    <el-container class="page-container">
      <!-- 头部 -->
      <el-header class="page-header">
        <div class="header-content">
          <div class="header-left">
            <el-icon class="header-icon"><VideoCamera /></el-icon>
            <h1>机场直播摄像头</h1>
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
          <el-col :span="24" :lg="16">
            <CameraLiveStream ref="cameraStreamRef" />
          </el-col>
          
          <!-- 侧边栏 -->
          <el-col :span="24" :lg="8">
            <!-- 摄像头列表 -->
            <el-card shadow="hover" class="camera-list-card">
              <template #header>
                <div class="card-header">
                  <span>摄像头列表</span>
                  <el-button 
                    type="primary" 
                    size="small" 
                    @click="showAddCameraDialog = true"
                    :icon="Plus"
                  >
                    添加摄像头
                  </el-button>
                </div>
              </template>
              
              <div class="camera-list">
                <div 
                  v-for="camera in cameraList" 
                  :key="camera.id"
                  class="camera-item"
                  :class="{ 'active': currentCameraId === camera.id }"
                  @click="switchCamera(camera)"
                >
                  <div class="camera-info">
                    <div class="camera-name">{{ camera.name }}</div>
                    <div class="camera-url">{{ camera.url }}</div>
                    <div class="camera-status">
                      <el-tag 
                        :type="camera.status === 'online' ? 'success' : 'info'" 
                        size="small"
                      >
                        {{ camera.status === 'online' ? '在线' : '离线' }}
                      </el-tag>
                    </div>
                  </div>
                  <div class="camera-actions">
                    <el-button 
                      type="text" 
                      size="small" 
                      @click.stop="editCamera(camera)"
                      :icon="Edit"
                    >
                      编辑
                    </el-button>
                    <el-button 
                      type="text" 
                      size="small" 
                      @click.stop="deleteCamera(camera)"
                      :icon="Delete"
                      class="delete-btn"
                    >
                      删除
                    </el-button>
                  </div>
                </div>
                
                <div v-if="cameraList.length === 0" class="empty-camera-list">
                  <el-empty description="暂无摄像头，请添加摄像头" />
                </div>
              </div>
            </el-card>

            <!-- 直播统计 -->
            <el-card shadow="hover" class="stats-card">
              <template #header>
                <span>直播统计</span>
              </template>
              
              <div class="stats-content">
                <div class="stat-item">
                  <div class="stat-label">直播时长</div>
                  <div class="stat-value">{{ formatDuration(streamDuration) }}</div>
                </div>
                <div class="stat-item">
                  <div class="stat-label">连接状态</div>
                  <div class="stat-value">
                    <el-tag :type="isStreaming ? 'success' : 'info'" size="small">
                      {{ isStreaming ? '已连接' : '未连接' }}
                    </el-tag>
                  </div>
                </div>
                <div class="stat-item">
                  <div class="stat-label">当前分辨率</div>
                  <div class="stat-value">{{ currentResolution || '未知' }}</div>
                </div>
                <div class="stat-item">
                  <div class="stat-label">帧率</div>
                  <div class="stat-value">{{ currentFps || '未知' }} FPS</div>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </el-main>
    </el-container>

    <!-- 添加/编辑摄像头对话框 -->
    <el-dialog
      v-model="showAddCameraDialog"
      :title="editingCamera ? `编辑摄像头 - ${editingCamera.name}` : '添加摄像头'"
      width="600px"
      @close="resetCameraForm"
    >
      <el-form :model="cameraForm" :rules="cameraRules" ref="cameraFormRef" label-width="100px">
        <el-form-item label="摄像头名称" prop="name">
          <el-input v-model="cameraForm.name" placeholder="请输入摄像头名称" />
        </el-form-item>
        <el-form-item label="摄像头地址" prop="url">
          <el-input v-model="cameraForm.url" placeholder="rtsp://ip:port/stream" />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="用户名">
              <el-input v-model="cameraForm.username" placeholder="用户名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码">
              <el-input v-model="cameraForm.password" type="password" placeholder="密码" show-password />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="分辨率">
              <el-select v-model="cameraForm.resolution" placeholder="选择分辨率">
                <el-option label="1920x1080" value="1920x1080" />
                <el-option label="1280x720" value="1280x720" />
                <el-option label="854x480" value="854x480" />
                <el-option label="640x360" value="640x360" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="帧率">
              <el-select v-model="cameraForm.fps" placeholder="选择帧率">
                <el-option label="30 FPS" :value="30" />
                <el-option label="25 FPS" :value="25" />
                <el-option label="15 FPS" :value="15" />
                <el-option label="10 FPS" :value="10" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述">
          <el-input 
            v-model="cameraForm.description" 
            type="textarea" 
            :rows="3" 
            placeholder="摄像头描述（可选）" 
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showAddCameraDialog = false">取消</el-button>
        <el-button type="primary" @click="saveCamera">
          {{ editingCamera ? '更新' : '添加' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  VideoCamera, 
  ArrowLeft, 
  Plus, 
  Edit, 
  Delete 
} from '@element-plus/icons-vue'
import CameraLiveStream from '@/components/CameraLiveStream.vue'

const router = useRouter()

// 响应式数据
const cameraStreamRef = ref(null)
const showAddCameraDialog = ref(false)
const editingCamera = ref(null)
const cameraFormRef = ref()
const currentCameraId = ref(null)
const streamDuration = ref(0)
const isStreaming = ref(false)
const currentResolution = ref('')
const currentFps = ref('')

// 摄像头列表
const cameraList = ref([
  {
    id: 1,
    name: '机场主摄像头',
    url: 'rtsp://192.168.1.100:554/stream1',
    username: 'admin',
    password: 'admin123',
    resolution: '1920x1080',
    fps: 25,
    status: 'offline',
    description: '机场主要监控摄像头'
  },
  {
    id: 2,
    name: '跑道摄像头',
    url: 'rtsp://192.168.1.101:554/stream1',
    username: 'admin',
    password: 'admin123',
    resolution: '1280x720',
    fps: 30,
    status: 'offline',
    description: '跑道监控摄像头'
  },
  {
    id: 3,
    name: '塔台摄像头',
    url: 'rtsp://192.168.1.102:554/stream1',
    username: 'admin',
    password: 'admin123',
    resolution: '1920x1080',
    fps: 25,
    status: 'offline',
    description: '塔台监控摄像头'
  }
])

// 摄像头表单
const cameraForm = reactive({
  name: '',
  url: '',
  username: '',
  password: '',
  resolution: '1280x720',
  fps: 25,
  description: ''
})

// 表单验证规则
const cameraRules = {
  name: [
    { required: true, message: '请输入摄像头名称', trigger: 'blur' }
  ],
  url: [
    { required: true, message: '请输入摄像头地址', trigger: 'blur' },
    { pattern: /^(rtsp|http|https):\/\/.+/, message: '请输入有效的摄像头地址', trigger: 'blur' }
  ],
  resolution: [
    { required: true, message: '请选择分辨率', trigger: 'change' }
  ],
  fps: [
    { required: true, message: '请选择帧率', trigger: 'change' }
  ]
}

// 方法
const goBack = () => {
  router.push('/')
}

const switchCamera = (camera) => {
  currentCameraId.value = camera.id
  // 更新摄像头配置
  if (cameraStreamRef.value) {
    Object.assign(cameraStreamRef.value.cameraConfig, {
      url: camera.url,
      username: camera.username,
      password: camera.password,
      resolution: camera.resolution,
      fps: camera.fps
    })
  }
  ElMessage.success(`已切换到摄像头: ${camera.name}`)
}

const editCamera = (camera) => {
  editingCamera.value = camera
  Object.assign(cameraForm, {
    name: camera.name,
    url: camera.url,
    username: camera.username,
    password: camera.password,
    resolution: camera.resolution,
    fps: camera.fps,
    description: camera.description
  })
  showAddCameraDialog.value = true
}

const deleteCamera = async (camera) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除摄像头 "${camera.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const index = cameraList.value.findIndex(c => c.id === camera.id)
    if (index > -1) {
      cameraList.value.splice(index, 1)
      ElMessage.success('摄像头删除成功')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除摄像头失败')
    }
  }
}

const saveCamera = async () => {
  try {
    await cameraFormRef.value.validate()
    
    if (editingCamera.value) {
      // 编辑摄像头
      const index = cameraList.value.findIndex(c => c.id === editingCamera.value.id)
      if (index > -1) {
        cameraList.value[index] = {
          ...cameraList.value[index],
          ...cameraForm
        }
      }
      ElMessage.success('摄像头更新成功')
    } else {
      // 添加摄像头
      const newCamera = {
        id: Date.now(),
        ...cameraForm,
        status: 'offline'
      }
      cameraList.value.push(newCamera)
      ElMessage.success('摄像头添加成功')
    }
    
    showAddCameraDialog.value = false
    resetCameraForm()
  } catch (error) {
    if (error.message) {
      ElMessage.error(`操作失败: ${error.message}`)
    } else {
      console.error('表单验证失败:', error)
    }
  }
}

const resetCameraForm = () => {
  editingCamera.value = null
  Object.assign(cameraForm, {
    name: '',
    url: '',
    username: '',
    password: '',
    resolution: '1280x720',
    fps: 25,
    description: ''
  })
  if (cameraFormRef.value) {
    cameraFormRef.value.resetFields()
  }
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
