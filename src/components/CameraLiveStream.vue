<template>
  <div class="camera-live-stream">
    <el-card shadow="hover" class="camera-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="camera-icon"><VideoCamera /></el-icon>
            <span>机场直播摄像头</span>
            <el-tag 
              :type="isStreaming ? 'success' : 'info'" 
              size="small"
              class="status-tag"
            >
              {{ isStreaming ? '直播中' : '未连接' }}
            </el-tag>
          </div>
          <div class="header-right">
            <el-button 
              v-if="!isStreaming"
              type="primary" 
              size="small" 
              @click="startStream"
              :loading="isConnecting"
              :icon="VideoPlay"
            >
              开始直播
            </el-button>
            <el-button 
              v-else
              type="danger" 
              size="small" 
              @click="stopStream"
              :icon="VideoPause"
            >
              停止直播
            </el-button>
            <el-button 
              size="small" 
              @click="toggleFullscreen"
              :icon="FullScreen"
            >
              全屏
            </el-button>
          </div>
        </div>
      </template>

      <!-- 摄像头配置 -->
      <div class="camera-config" v-if="!isStreaming">
        <el-form :model="cameraConfig" label-width="100px" size="small">
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="摄像头地址">
                <el-input 
                  v-model="cameraConfig.url" 
                  placeholder="rtsp://ip:port/stream"
                  clearable
                />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="用户名">
                <el-input 
                  v-model="cameraConfig.username" 
                  placeholder="用户名"
                  clearable
                />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="密码">
                <el-input 
                  v-model="cameraConfig.password" 
                  type="password" 
                  placeholder="密码"
                  show-password
                  clearable
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="8">
              <el-form-item label="分辨率">
                <el-select v-model="cameraConfig.resolution" placeholder="选择分辨率">
                  <el-option label="1920x1080" value="1920x1080" />
                  <el-option label="1280x720" value="1280x720" />
                  <el-option label="854x480" value="854x480" />
                  <el-option label="640x360" value="640x360" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="帧率">
                <el-select v-model="cameraConfig.fps" placeholder="选择帧率">
                  <el-option label="30 FPS" :value="30" />
                  <el-option label="25 FPS" :value="25" />
                  <el-option label="15 FPS" :value="15" />
                  <el-option label="10 FPS" :value="10" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="质量">
                <el-slider 
                  v-model="cameraConfig.quality" 
                  :min="1" 
                  :max="10" 
                  :step="1"
                  show-stops
                  show-tooltip
                />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </div>

      <!-- 视频播放区域 -->
      <div class="video-container" :class="{ 'fullscreen': isFullscreen }">
        <div v-if="!isStreaming" class="video-placeholder">
          <el-icon class="placeholder-icon"><VideoCamera /></el-icon>
          <p>点击"开始直播"连接摄像头</p>
          <p class="placeholder-desc">支持 RTSP、HTTP、WebRTC 等协议</p>
        </div>
        
        <div v-else class="video-wrapper">
          <video 
            ref="videoElement"
            :src="streamUrl"
            autoplay
            muted
            playsinline
            class="video-stream"
            @loadstart="onLoadStart"
            @loadeddata="onLoadedData"
            @error="onError"
            @play="onPlay"
            @pause="onPause"
          />
          
          <!-- 视频控制覆盖层 -->
          <div class="video-overlay" v-if="isStreaming">
            <div class="video-info">
              <div class="stream-info">
                <span class="info-item">
                  <el-icon><Timer /></el-icon>
                  {{ formatDuration(streamDuration) }}
                </span>
                <span class="info-item">
                  <el-icon><Connection /></el-icon>
                  {{ connectionStatus }}
                </span>
                <span class="info-item" v-if="cameraConfig.resolution">
                  <el-icon><Monitor /></el-icon>
                  {{ cameraConfig.resolution }}
                </span>
              </div>
            </div>
            
            <div class="video-controls">
              <el-button 
                circle 
                size="small" 
                @click="toggleMute"
                :icon="isMuted ? Mute : Microphone"
              />
              <el-button 
                circle 
                size="small" 
                @click="captureScreenshot"
                :icon="Camera"
              />
              <el-button 
                circle 
                size="small" 
                @click="toggleRecording"
                :type="isRecording ? 'danger' : 'default'"
                :icon="isRecording ? VideoPause : VideoPlay"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 截图预览对话框 -->
      <el-dialog
        v-model="screenshotDialogVisible"
        title="截图预览"
        width="800px"
        @close="screenshotDialogVisible = false"
      >
        <div class="screenshot-preview">
          <img :src="screenshotDataUrl" alt="截图预览" class="screenshot-image" />
        </div>
        <template #footer>
          <el-button @click="screenshotDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="downloadScreenshot">下载截图</el-button>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  VideoCamera, 
  VideoPlay, 
  VideoPause, 
  FullScreen, 
  Timer, 
  Connection, 
  Monitor,
  Mute,
  Microphone,
  Camera
} from '@element-plus/icons-vue'

// 响应式数据
const videoElement = ref(null)
const isStreaming = ref(false)
const isConnecting = ref(false)
const isFullscreen = ref(false)
const isMuted = ref(false)
const isRecording = ref(false)
const streamDuration = ref(0)
const connectionStatus = ref('未连接')
const screenshotDialogVisible = ref(false)
const screenshotDataUrl = ref('')

// 摄像头配置
const cameraConfig = reactive({
  url: 'rtsp://192.168.1.100:554/stream1',
  username: '',
  password: '',
  resolution: '1280x720',
  fps: 25,
  quality: 7
})

// 计算属性
const streamUrl = computed(() => {
  if (!cameraConfig.url) return ''
  
  let url = cameraConfig.url
  if (cameraConfig.username && cameraConfig.password) {
    const protocol = url.split('://')[0]
    const rest = url.split('://')[1]
    url = `${protocol}://${cameraConfig.username}:${cameraConfig.password}@${rest}`
  }
  
  return url
})

// 定时器
let durationTimer = null
let streamTimer = null

// 方法
const startStream = async () => {
  if (!cameraConfig.url) {
    ElMessage.warning('请输入摄像头地址')
    return
  }
  
  isConnecting.value = true
  connectionStatus.value = '连接中...'
  
  try {
    // 模拟连接延迟
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    if (videoElement.value) {
      videoElement.value.src = streamUrl.value
      videoElement.value.load()
      
      // 等待视频加载
      await new Promise((resolve, reject) => {
        const timeout = setTimeout(() => {
          reject(new Error('连接超时'))
        }, 10000)
        
        videoElement.value.onloadeddata = () => {
          clearTimeout(timeout)
          resolve()
        }
        
        videoElement.value.onerror = () => {
          clearTimeout(timeout)
          reject(new Error('视频加载失败'))
        }
      })
    }
    
    isStreaming.value = true
    isConnecting.value = false
    connectionStatus.value = '已连接'
    streamDuration.value = 0
    
    // 开始计时
    startDurationTimer()
    
    ElMessage.success('摄像头连接成功')
  } catch (error) {
    isConnecting.value = false
    connectionStatus.value = '连接失败'
    ElMessage.error(`连接失败: ${error.message}`)
  }
}

const stopStream = () => {
  if (videoElement.value) {
    videoElement.value.pause()
    videoElement.value.src = ''
  }
  
  isStreaming.value = false
  connectionStatus.value = '未连接'
  streamDuration.value = 0
  
  // 停止计时
  stopDurationTimer()
  
  ElMessage.info('已停止直播')
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}

const toggleMute = () => {
  if (videoElement.value) {
    videoElement.value.muted = !videoElement.value.muted
    isMuted.value = videoElement.value.muted
  }
}

const captureScreenshot = () => {
  if (!videoElement.value || !isStreaming.value) {
    ElMessage.warning('请先开始直播')
    return
  }
  
  try {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    
    canvas.width = videoElement.value.videoWidth
    canvas.height = videoElement.value.videoHeight
    
    ctx.drawImage(videoElement.value, 0, 0)
    
    screenshotDataUrl.value = canvas.toDataURL('image/png')
    screenshotDialogVisible.value = true
    
    ElMessage.success('截图成功')
  } catch (error) {
    ElMessage.error(`截图失败: ${error.message}`)
  }
}

const downloadScreenshot = () => {
  if (screenshotDataUrl.value) {
    const link = document.createElement('a')
    link.download = `camera-screenshot-${new Date().getTime()}.png`
    link.href = screenshotDataUrl.value
    link.click()
  }
}

const toggleRecording = () => {
  isRecording.value = !isRecording.value
  ElMessage.info(isRecording.value ? '开始录制' : '停止录制')
}

const startDurationTimer = () => {
  durationTimer = setInterval(() => {
    streamDuration.value++
  }, 1000)
}

const stopDurationTimer = () => {
  if (durationTimer) {
    clearInterval(durationTimer)
    durationTimer = null
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

// 视频事件处理
const onLoadStart = () => {
  connectionStatus.value = '加载中...'
}

const onLoadedData = () => {
  connectionStatus.value = '已连接'
}

const onError = (event) => {
  console.error('视频播放错误:', event)
  connectionStatus.value = '播放错误'
  ElMessage.error('视频播放失败，请检查摄像头配置')
}

const onPlay = () => {
  connectionStatus.value = '播放中'
}

const onPause = () => {
  connectionStatus.value = '已暂停'
}

// 生命周期
onMounted(() => {
  // 加载保存的配置
  const savedConfig = localStorage.getItem('camera-config')
  if (savedConfig) {
    try {
      Object.assign(cameraConfig, JSON.parse(savedConfig))
    } catch (error) {
      console.error('加载摄像头配置失败:', error)
    }
  }
})

onUnmounted(() => {
  stopStream()
  stopDurationTimer()
})

// 保存配置
const saveConfig = () => {
  localStorage.setItem('camera-config', JSON.stringify(cameraConfig))
}
</script>

<style scoped>
.camera-live-stream {
  width: 100%;
}

.camera-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.camera-icon {
  font-size: 18px;
  color: #409eff;
}

.status-tag {
  margin-left: 8px;
}

.header-right {
  display: flex;
  gap: 8px;
}

.camera-config {
  margin-bottom: 16px;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
}

.video-container {
  position: relative;
  width: 100%;
  height: 400px;
  background-color: #000;
  border-radius: 6px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.video-container.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 9999;
  border-radius: 0;
}

.video-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
}

.placeholder-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.placeholder-desc {
  font-size: 12px;
  margin-top: 8px;
  opacity: 0.7;
}

.video-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}

.video-stream {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    to bottom,
    rgba(0, 0, 0, 0.3) 0%,
    transparent 20%,
    transparent 80%,
    rgba(0, 0, 0, 0.3) 100%
  );
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.video-container:hover .video-overlay {
  opacity: 1;
}

.video-info {
  padding: 12px;
}

.stream-info {
  display: flex;
  gap: 16px;
  color: white;
  font-size: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 4px;
  background-color: rgba(0, 0, 0, 0.5);
  padding: 4px 8px;
  border-radius: 4px;
}

.video-controls {
  display: flex;
  justify-content: center;
  gap: 8px;
  padding: 12px;
}

.video-controls .el-button {
  background-color: rgba(0, 0, 0, 0.5);
  border-color: rgba(255, 255, 255, 0.3);
  color: white;
}

.video-controls .el-button:hover {
  background-color: rgba(0, 0, 0, 0.7);
  border-color: rgba(255, 255, 255, 0.5);
}

.screenshot-preview {
  text-align: center;
}

.screenshot-image {
  max-width: 100%;
  max-height: 500px;
  border-radius: 6px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .video-container {
    height: 300px;
  }
  
  .camera-config .el-row {
    margin-bottom: 8px;
  }
  
  .stream-info {
    flex-wrap: wrap;
    gap: 8px;
  }
}
</style>
