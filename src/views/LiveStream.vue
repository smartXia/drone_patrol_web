<template>
  <div class="live-stream-page">
    <div class="page-header">
      <h2>直播流播放</h2>
      <div class="stream-info">
        <el-tag type="primary">WHEP协议</el-tag>
        <span class="stream-url">{{ streamUrl }}</span>
      </div>
    </div>

    <div class="video-container">
      <div class="video-wrapper">
        <video
          ref="videoElement"
          class="video-player"
          controls
          autoplay
          muted
          playsinline
          @loadstart="onLoadStart"
          @loadeddata="onLoadedData"
          @error="onError"
          @play="onPlay"
          @pause="onPause"
        >
          您的浏览器不支持视频播放
        </video>
        
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-overlay">
          <el-icon class="loading-icon"><Loading /></el-icon>
          <p>正在加载直播流...</p>
        </div>

        <!-- 错误状态 -->
        <div v-if="error" class="error-overlay">
          <el-icon class="error-icon"><Warning /></el-icon>
          <p>{{ error }}</p>
          <el-button type="primary" @click="retryConnection">重试连接</el-button>
        </div>

        <!-- 连接状态 -->
        <div class="connection-status">
          <el-tag :type="connectionStatus.type" size="small">
            {{ connectionStatus.text }}
          </el-tag>
          <span v-if="lastUpdateTime" class="update-time">
            最后更新: {{ lastUpdateTime }}
          </span>
        </div>
      </div>
    </div>

    <!-- 控制面板 -->
    <div class="control-panel">
      <el-row :gutter="16">
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><VideoPlay /></el-icon>
                <span>播放控制</span>
              </div>
            </template>
            <div class="control-buttons">
              <el-button 
                type="primary" 
                :icon="isPlaying ? VideoPause : VideoPlay"
                @click="togglePlay"
                :disabled="loading || error"
              >
                {{ isPlaying ? '暂停' : '播放' }}
              </el-button>
              <el-button 
                type="info" 
                :icon="Refresh"
                @click="retryConnection"
                :loading="loading"
              >
                重新连接
              </el-button>
            </div>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Setting /></el-icon>
                <span>流信息</span>
              </div>
            </template>
            <div class="stream-details">
              <div class="detail-item">
                <span class="label">协议:</span>
                <span class="value">WHEP (WebRTC HTTP Egress Protocol)</span>
              </div>
              <div class="detail-item">
                <span class="label">应用:</span>
                <span class="value">live</span>
              </div>
              <div class="detail-item">
                <span class="label">流ID:</span>
                <span class="value">8UUXN2B00A00ST-165-0-7</span>
              </div>
              <div class="detail-item">
                <span class="label">服务器:</span>
                <span class="value">117.62.203.197:8080</span>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Monitor /></el-icon>
                <span>播放状态</span>
              </div>
            </template>
            <div class="playback-status">
              <div class="status-item">
                <span class="label">播放状态:</span>
                <el-tag :type="isPlaying ? 'success' : 'info'" size="small">
                  {{ isPlaying ? '播放中' : '已暂停' }}
                </el-tag>
              </div>
              <div class="status-item">
                <span class="label">连接状态:</span>
                <el-tag :type="connectionStatus.type" size="small">
                  {{ connectionStatus.text }}
                </el-tag>
              </div>
              <div class="status-item">
                <span class="label">视频尺寸:</span>
                <span class="value">{{ videoDimensions }}</span>
              </div>
              <div class="status-item">
                <span class="label">播放时长:</span>
                <span class="value">{{ playTime }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Loading, 
  Warning, 
  VideoPlay, 
  VideoPause, 
  Refresh, 
  Setting, 
  Monitor 
} from '@element-plus/icons-vue'

// 响应式数据
const videoElement = ref(null)
const loading = ref(false)
const error = ref('')
const isPlaying = ref(false)
const lastUpdateTime = ref('')
const playStartTime = ref(null)
const playTime = ref('00:00:00')

// 流URL
const streamUrl = 'http://117.62.203.197:8080/rtc/v1/whep/?app=live&stream=8UUXN2B00A00ST-165-0-7'

// 计算属性
const connectionStatus = computed(() => {
  if (loading.value) {
    return { type: 'warning', text: '连接中...' }
  }
  if (error.value) {
    return { type: 'danger', text: '连接失败' }
  }
  if (isPlaying.value) {
    return { type: 'success', text: '已连接' }
  }
  return { type: 'info', text: '未连接' }
})

const videoDimensions = computed(() => {
  if (videoElement.value && videoElement.value.videoWidth && videoElement.value.videoHeight) {
    return `${videoElement.value.videoWidth} x ${videoElement.value.videoHeight}`
  }
  return '未知'
})

// 定时器
let playTimeInterval = null

// 方法
const startPlayTime = () => {
  playStartTime.value = Date.now()
  playTimeInterval = setInterval(() => {
    if (playStartTime.value) {
      const elapsed = Math.floor((Date.now() - playStartTime.value) / 1000)
      const hours = Math.floor(elapsed / 3600)
      const minutes = Math.floor((elapsed % 3600) / 60)
      const seconds = elapsed % 60
      playTime.value = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
    }
  }, 1000)
}

const stopPlayTime = () => {
  if (playTimeInterval) {
    clearInterval(playTimeInterval)
    playTimeInterval = null
  }
}

const connectToStream = async () => {
  if (!videoElement.value) return

  loading.value = true
  error.value = ''
  lastUpdateTime.value = new Date().toLocaleString()

  try {
    // 使用WebRTC连接WHEP流
    await connectWhepStream()
    
    loading.value = false
    ElMessage.success('直播流连接成功')
  } catch (err) {
    loading.value = false
    error.value = `连接失败: ${err.message}`
    ElMessage.error('直播流连接失败')
    console.error('直播流连接错误:', err)
  }
}

const connectWhepStream = async () => {
  try {
    // 创建RTCPeerConnection
    const pc = new RTCPeerConnection({
      iceServers: [
        { urls: 'stun:stun.l.google.com:19302' },
        { urls: 'stun:stun1.l.google.com:19302' }
      ]
    })

    // 处理接收到的流
    pc.ontrack = (event) => {
      console.log('接收到流:', event)
      if (event.streams && event.streams[0]) {
        videoElement.value.srcObject = event.streams[0]
        videoElement.value.play()
      }
    }

    // 处理ICE连接状态
    pc.oniceconnectionstatechange = () => {
      console.log('ICE连接状态:', pc.iceConnectionState)
      if (pc.iceConnectionState === 'connected' || pc.iceConnectionState === 'completed') {
        console.log('WebRTC连接成功')
      } else if (pc.iceConnectionState === 'failed' || pc.iceConnectionState === 'disconnected') {
        console.log('WebRTC连接失败')
        error.value = 'WebRTC连接失败'
      }
    }

    // 创建接收offer
    const offer = await pc.createOffer({
      offerToReceiveAudio: true,
      offerToReceiveVideo: true
    })
    
    await pc.setLocalDescription(offer)

    // 发送WHEP请求
    const response = await fetch(streamUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/sdp'
      },
      body: offer.sdp
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    const answerSdp = await response.text()
    console.log('收到SDP应答:', answerSdp)

    // 设置远程描述
    await pc.setRemoteDescription({
      type: 'answer',
      sdp: answerSdp
    })

    console.log('WHEP连接建立成功')

  } catch (err) {
    console.error('WHEP连接错误:', err)
    throw err
  }
}

const retryConnection = () => {
  stopPlayTime()
  playTime.value = '00:00:00'
  connectToStream()
}

const togglePlay = () => {
  if (!videoElement.value) return

  if (isPlaying.value) {
    videoElement.value.pause()
  } else {
    videoElement.value.play()
  }
}

// 事件处理
const onLoadStart = () => {
  console.log('开始加载视频流')
}

const onLoadedData = () => {
  console.log('视频数据加载完成')
  loading.value = false
}

const onError = (event) => {
  console.error('视频播放错误:', event)
  error.value = '视频播放出错，请检查网络连接或流地址'
  loading.value = false
  stopPlayTime()
}

const onPlay = () => {
  isPlaying.value = true
  startPlayTime()
  console.log('开始播放')
}

const onPause = () => {
  isPlaying.value = false
  stopPlayTime()
  console.log('暂停播放')
}

// 生命周期
onMounted(() => {
  console.log('直播流页面已挂载')
  // 延迟连接，确保DOM已渲染
  setTimeout(() => {
    connectToStream()
  }, 500)
})

onUnmounted(() => {
  stopPlayTime()
  if (videoElement.value) {
    videoElement.value.pause()
    videoElement.value.src = ''
  }
})
</script>

<style scoped>
.live-stream-page {
  padding: 20px;
  background: var(--el-bg-color-page);
  min-height: 100vh;
}

.page-header {
  margin-bottom: 20px;
  text-align: center;
}

.page-header h2 {
  color: var(--el-text-color-primary);
  margin-bottom: 10px;
}

.stream-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.stream-url {
  font-family: 'Courier New', monospace;
  background: var(--el-color-info-light-9);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: var(--el-color-info);
}

.video-container {
  margin-bottom: 20px;
}

.video-wrapper {
  position: relative;
  background: #000;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.video-player {
  width: 100%;
  height: 500px;
  object-fit: contain;
  background: #000;
}

.loading-overlay,
.error-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  z-index: 10;
}

.loading-icon,
.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.loading-icon {
  color: var(--el-color-primary);
  animation: spin 1s linear infinite;
}

.error-icon {
  color: var(--el-color-danger);
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.connection-status {
  position: absolute;
  top: 10px;
  right: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
  z-index: 5;
}

.update-time {
  color: white;
  font-size: 12px;
  background: rgba(0, 0, 0, 0.6);
  padding: 2px 6px;
  border-radius: 4px;
}

.control-panel {
  margin-top: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.control-buttons {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.stream-details,
.playback-status {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-item,
.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
}

.label {
  color: var(--el-text-color-regular);
  font-size: 13px;
}

.value {
  color: var(--el-text-color-primary);
  font-weight: 500;
  font-size: 13px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .live-stream-page {
    padding: 10px;
  }
  
  .video-player {
    height: 300px;
  }
  
  .control-buttons {
    flex-direction: column;
  }
  
  .stream-info {
    flex-direction: column;
    gap: 5px;
  }
}
</style>
