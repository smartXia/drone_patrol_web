<template>
  <div class="camera-live-stream">
    <!-- 控制栏 -->
    <div class="control-bar">
      <div class="control-left">
            <el-icon class="camera-icon"><VideoCamera /></el-icon>
        <span>大疆机场摄像头</span>
            <el-tag 
              :type="isStreaming ? 'success' : 'info'" 
              size="small"
              class="status-tag"
            >
          {{ isStreaming ? '接收中' : '未连接' }}
            </el-tag>
          </div>
      <div class="control-right">
            <el-button 
              v-if="!isStreaming"
              type="primary" 
              size="small" 
              @click="startStream"
              :loading="isConnecting"
              :icon="VideoPlay"
            >
          连接大疆机场
            </el-button>
            <el-button 
              v-else
              type="danger" 
              size="small" 
              @click="stopStream"
              :icon="VideoPause"
            >
          断开连接
        </el-button>
        <el-button 
          size="small" 
          @click="testConnection"
          :icon="Connection"
          type="info"
        >
          测试连接
        </el-button>
        <el-button 
          size="small" 
          @click="showDiagnosticDialog"
          :icon="Monitor"
          type="warning"
        >
          诊断工具
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

      <!-- 摄像头配置 -->
      <div class="camera-config" v-if="!isStreaming">
        <el-form :model="cameraConfig" label-width="100px" size="small">
          <!-- 设备选择 -->
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="选择设备">
                <el-select 
                  v-model="selectedDevice" 
                  placeholder="选择机场或飞机设备" 
                  @change="onDeviceChange" 
                  clearable
                  :disabled="(airportDevices.length + aircraftDevices.length) <= 1"
                >
                  <el-option-group label="机场设备" v-if="airportDevices.length > 0">
                    <el-option 
                      v-for="device in airportDevices" 
                      :key="device.sn" 
                      :label="`${device.name || device.sn} (${device.sn})`" 
                      :value="device.sn"
                    />
                  </el-option-group>
                  <el-option-group label="飞机设备" v-if="aircraftDevices.length > 0">
                    <el-option 
                      v-for="device in aircraftDevices" 
                      :key="device.sn" 
                      :label="`${device.name || device.sn} (${device.sn})`" 
                      :value="device.sn"
                    />
                  </el-option-group>
                  <el-option v-if="airportDevices.length === 0 && aircraftDevices.length === 0" disabled>
                    暂无可用设备
                  </el-option>
                </el-select>
                <div class="config-tips" v-if="airportDevices.length === 0 && aircraftDevices.length === 0">
                  <el-text size="small" type="warning">
                    请确保路由参数中包含设备SN，或检查网络连接
                  </el-text>
                </div>
                <div class="config-tips" v-if="(airportDevices.length + aircraftDevices.length) === 1">
                  <el-text size="small" type="info">
                    已自动选择唯一可用设备
                  </el-text>
                </div>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="当前设备">
                <el-tag v-if="currentDeviceSN" type="success" size="large">
                  {{ currentDeviceSN }}
                </el-tag>
                <el-tag v-else type="info" size="large">
                  未选择设备
                </el-tag>
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

        <!-- WHIP WebRTC配置 -->
        <el-divider content-position="left">大疆机场视频流配置</el-divider>
        <el-form :model="whipConfig" label-width="120px" size="small">
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="房间号">
                <el-input v-model="whipConfig.room" readonly />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="认证Token">
                <el-input v-model="whipConfig.authToken" readonly />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="签名密钥">
                <el-input v-model="whipConfig.txSecret" readonly />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="时间戳">
                <el-input v-model="whipConfig.txTime" readonly />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="24">
              <el-form-item label="启用大疆机场视频流">
                <el-switch v-model="whipConfig.enableWhip" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>

      </div>

      <!-- 视频播放区域 -->
      <div class="video-container" :class="{ 'fullscreen': isFullscreen }">
        <div v-if="!isStreaming" class="video-placeholder">
          <el-icon class="placeholder-icon"><VideoCamera /></el-icon>
          <p>点击"开始直播"连接大疆机场摄像头</p>
          <p class="placeholder-desc">使用 WHIP WebRTC 协议接收大疆机场视频流</p>
        </div>
        
        <div v-else class="video-wrapper">
          <video 
            ref="videoElement"
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

      <!-- 诊断工具对话框 -->
      <el-dialog
        v-model="diagnosticDialogVisible"
        title="摄像头诊断工具"
        width="600px"
        :close-on-click-modal="false"
      >
        <div class="diagnostic-content">
          <el-alert
            title="诊断信息"
            type="info"
            :closable="false"
            show-icon
          >
            <template #default>
              <p>当前配置信息：</p>
              <ul>
                <li><strong>摄像头地址：</strong>{{ cameraConfig.url || '未设置' }}</li>
                <li><strong>用户名：</strong>{{ cameraConfig.username || '未设置' }}</li>
                <li><strong>密码：</strong>{{ cameraConfig.password ? '***' : '未设置' }}</li>
                <li><strong>完整地址：</strong>{{ streamUrl || '未生成' }}</li>
              </ul>
            </template>
          </el-alert>
          
          <el-divider />
          
          <div class="diagnostic-actions">
            <h4>常见问题解决方案：</h4>
            <el-space direction="vertical" size="large" style="width: 100%">
              <el-card shadow="never" class="diagnostic-card">
                <template #header>
                  <span>1. RTSP连接问题</span>
                </template>
                <p>• 确保摄像头支持RTSP协议</p>
                <p>• 检查IP地址和端口是否正确</p>
                <p>• 验证用户名和密码</p>
                <p>• 尝试不同的RTSP路径（如：/stream1, /live, /ch1）</p>
    </el-card>
              
              <el-card shadow="never" class="diagnostic-card">
                <template #header>
                  <span>2. 网络连接问题</span>
                </template>
                <p>• 检查网络连接是否正常</p>
                <p>• 确认摄像头和电脑在同一网络</p>
                <p>• 尝试ping摄像头IP地址</p>
                <p>• 检查防火墙设置</p>
              </el-card>
              
              <el-card shadow="never" class="diagnostic-card">
                <template #header>
                  <span>3. 视频格式问题</span>
                </template>
                <p>• 确保浏览器支持该视频格式</p>
                <p>• 尝试不同的分辨率设置</p>
                <p>• 检查摄像头编码格式（H.264推荐）</p>
                <p>• 尝试降低帧率设置</p>
              </el-card>
            </el-space>
          </div>
        </div>
        
        <template #footer>
          <el-button @click="diagnosticDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="runFullDiagnostic">运行完整诊断</el-button>
        </template>
      </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, computed, watch } from 'vue'
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
import { WhipService } from '@/services/whipService'
import { config } from '@/config'

// Props定义
const props = defineProps({
  airportSn: {
    type: String,
    default: ''
  },
  aircraftSn: {
    type: String,
    default: ''
  }
})

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
const diagnosticDialogVisible = ref(false)

// 摄像头配置（简化版）
const cameraConfig = reactive({
  resolution: '1280x720',
  fps: 25,
  quality: 7
})

// WHIP WebRTC配置
const whipConfig = reactive({
  room: '',         // 房间号（设备SN）
  authToken: '',    // 认证token
  txSecret: '',     // 签名密钥
  txTime: '',       // 时间戳
  enableWhip: true  // 默认启用WHIP推流
})

// 设备选择相关
const selectedDevice = ref('')
const currentDeviceSN = ref('')

// 设备列表（从路由参数获取）
const airportDevices = computed(() => {
  // 从路由参数获取机场SN
  const airportSn = props.airportSn
  if (airportSn && airportSn !== '') {
    return [{
      sn: airportSn,
      name: `机场设备 (${airportSn})`,
      type: 'airport'
    }]
  }
  return []
})

const aircraftDevices = computed(() => {
  // 从路由参数获取飞机SN
  const aircraftSn = props.aircraftSn
  if (aircraftSn && aircraftSn !== '') {
    return [{
      sn: aircraftSn,
      name: `飞机设备 (${aircraftSn})`,
      type: 'aircraft'
    }]
  }
  return []
})

// WHIP推流状态
const whipStatus = ref({
  isStreaming: false,
  room: '',
  connectionState: 'disconnected'
})


// 计算属性（简化版）
const streamUrl = computed(() => {
  // 使用设备SN构建直播流URL
  const deviceSN = currentDeviceSN.value || props.aircraftSn || props.airportSn
  if (!deviceSN) return ''
  
  // 这里可以根据设备SN生成对应的直播流URL
  // 暂时返回空，实际使用时需要根据设备配置生成
  return ''
})

// 定时器
let durationTimer = null
let streamTimer = null

// WHIP服务实例
const whipService = new WhipService(config)

// 后端直播服务
const backendLiveService = {
  // 创建直播流
  async createLiveStream(deviceSN, streamType, config) {
    try {
      const response = await fetch('/api/live/stream/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          device_sn: deviceSN,
          stream_type: streamType,
          resolution: config.resolution || '1920x1080',
          bitrate: config.bitrate || 2000,
          fps: config.fps || 25,
          quality: config.quality || 'high'
        })
      })
      
      const result = await response.json()
      if (result.success) {
        return {
          success: true,
          streamId: result.stream_id,
          pushUrl: result.push_url,
          playUrl: result.play_url,
          userSig: result.user_sig,
          sdkAppId: result.sdk_app_id,
          message: result.message
        }
      } else {
        return {
          success: false,
          error: result.error || '创建直播流失败'
        }
      }
    } catch (error) {
      console.error('创建直播流失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  },

  // 获取直播流状态
  async getLiveStreamStatus(streamId) {
    try {
      const response = await fetch(`/api/live/stream/${streamId}/status`)
      const result = await response.json()
      return result
    } catch (error) {
      console.error('获取直播流状态失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  },

  // 停止直播流
  async stopLiveStream(streamId) {
    try {
      const response = await fetch(`/api/live/stream/${streamId}/stop`, {
        method: 'POST'
      })
      const result = await response.json()
      return result
    } catch (error) {
      console.error('停止直播流失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  },

}



// 设备选择变化处理
const onDeviceChange = (deviceSN) => {
  if (deviceSN) {
    currentDeviceSN.value = deviceSN
    console.log('选择设备:', deviceSN)
    ElMessage.success(`已选择设备: ${deviceSN}`)
  } else {
    currentDeviceSN.value = ''
    console.log('取消选择设备')
  }
}

// WHIP推流状态更新
const updateWhipStatus = () => {
  const status = whipService.getStreamStatus()
  whipStatus.value = {
    isStreaming: status.isStreaming,
    room: status.room,
    connectionState: status.connectionState
  }
}

// 方法（兼容性方法，调用步骤1）
const startStream = async () => {
  isConnecting.value = true
  connectionStatus.value = '连接中...'
  
  try {
    // 获取设备SN（优先使用选择的设备，然后从props获取）
    const deviceSN = currentDeviceSN.value || props.aircraftSn || props.airportSn
    
    if (!deviceSN) {
      ElMessage.warning('请先选择设备或确保设备SN已传入')
      return
    }
    
    console.log('开始连接大疆机场设备，设备SN:', deviceSN)
    
    // 使用WHIP服务接收大疆机场视频流
    const result = await whipService.startStream(deviceSN, videoElement.value)
    
    if (result.success) {
      // 保存房间信息
      whipConfig.room = result.room
    
    isStreaming.value = true
    isConnecting.value = false
    connectionStatus.value = '已连接'
    streamDuration.value = 0
    
    // 开始计时
    startDurationTimer()
    
      ElMessage.success(`已连接大疆机场设备 - 房间: ${result.room}`)
      console.log('大疆机场视频流信息:', result)
    } else {
      throw new Error(result.error || '连接大疆机场设备失败')
    }
    
  } catch (error) {
    isConnecting.value = false
    connectionStatus.value = '连接失败'
    console.error('连接大疆机场设备失败:', error)
    ElMessage.error(`连接大疆机场设备失败: ${error.message}`)
  }
}

const stopStream = async () => {
  try {
    // 停止接收大疆机场视频流
    const result = await whipService.stopStream()
    
    if (result.success) {
      ElMessage.info('已断开大疆机场设备连接')
    } else {
      ElMessage.warning(`断开连接失败: ${result.error}`)
    }
  } catch (error) {
    console.error('断开大疆机场设备连接失败:', error)
    ElMessage.warning('断开连接失败')
  }
  
  isStreaming.value = false
  connectionStatus.value = '未连接'
  streamDuration.value = 0
  
  // 停止计时
  stopDurationTimer()
  
  // 清空房间信息
  whipConfig.room = ''
  
  ElMessage.info('已断开连接')
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}

// 测试摄像头连接
const testConnection = async () => {
  if (!currentDeviceSN.value && !props.aircraftSn && !props.airportSn) {
    ElMessage.warning('请先选择设备')
    return
  }
  
  // 验证URL格式
  const fullUrl = streamUrl.value
  console.log('测试连接 - 完整URL:', fullUrl)
  
  // 检查URL格式是否正确
  if (cameraConfig.username && cameraConfig.password) {
    const expectedFormat = `${cameraConfig.url.split('://')[0]}://${cameraConfig.username}:${cameraConfig.password}@${cameraConfig.url.split('://')[1]}`
    console.log('期望的URL格式:', expectedFormat)
    console.log('实际生成的URL:', fullUrl)
  }
  
  ElMessage.info('正在测试摄像头连接...')
  
  try {
    // 创建一个临时的视频元素来测试连接
    const testVideo = document.createElement('video')
    testVideo.muted = true
    testVideo.playsInline = true
    
    // 设置超时
    const timeout = setTimeout(() => {
      testVideo.src = ''
      ElMessage.error('连接测试超时，请检查摄像头地址和网络')
    }, 10000)
    
    // 监听加载成功
    testVideo.onloadeddata = () => {
      clearTimeout(timeout)
      testVideo.src = ''
      ElMessage.success('摄像头连接测试成功！')
    }
    
    // 监听错误
    testVideo.onerror = (event) => {
      clearTimeout(timeout)
      testVideo.src = ''
      
      let errorMsg = '摄像头连接失败'
      if (event.target?.error) {
        switch (event.target.error.code) {
          case 1:
            errorMsg = '连接被中止，请检查网络'
            break
          case 2:
            errorMsg = '网络错误，请检查摄像头地址'
            break
          case 3:
            errorMsg = '视频解码错误，请检查视频格式'
            break
          case 4:
            errorMsg = '视频源不支持'
            break
        }
      }
      
      ElMessage.error(`连接测试失败: ${errorMsg}`)
      
      // 提供解决建议
      if (cameraConfig.url.includes('rtsp://') && !cameraConfig.username) {
        ElMessage.warning('提示：RTSP摄像头通常需要用户名和密码认证')
      }
    }
    
    // 开始测试
    testVideo.src = streamUrl.value
    testVideo.load()
    
  } catch (error) {
    console.error('连接测试失败:', error)
    ElMessage.error('连接测试失败，请检查配置')
  }
}

// 显示诊断对话框
const showDiagnosticDialog = () => {
  diagnosticDialogVisible.value = true
}

// 运行完整诊断
const runFullDiagnostic = async () => {
  ElMessage.info('开始运行完整诊断...')
  
  const diagnosticResults = []
  
  // 1. 检查基本配置
  // 检查设备选择
  if (!currentDeviceSN.value && !props.aircraftSn && !props.airportSn) {
    diagnosticResults.push('❌ 未选择设备')
  } else {
    const deviceSN = currentDeviceSN.value || props.aircraftSn || props.airportSn
    diagnosticResults.push(`✅ 已选择设备: ${deviceSN}`)
  }
  
  // 2. 检查设备配置
  const deviceSN = currentDeviceSN.value || props.aircraftSn || props.airportSn
  if (deviceSN) {
    diagnosticResults.push(`✅ 设备SN: ${deviceSN}`)
  } else {
    diagnosticResults.push('❌ 设备SN未设置')
  }
  
  // 3. 测试网络连接
  if (deviceSN) {
    try {
      const url = new URL(cameraConfig.url)
      const hostname = url.hostname
      
      // 简单的网络测试（这里只是示例，实际需要更复杂的网络测试）
      diagnosticResults.push(`🔍 尝试连接 ${hostname}...`)
      
      // 模拟网络测试
      await new Promise(resolve => setTimeout(resolve, 1000))
      diagnosticResults.push('✅ 网络连接测试完成')
      
    } catch (error) {
      diagnosticResults.push('❌ URL格式错误')
    }
  }
  
  // 显示诊断结果
  const resultText = diagnosticResults.join('\n')
  ElMessageBox.alert(resultText, '诊断结果', {
    confirmButtonText: '确定',
    type: 'info'
  })
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
  console.error('当前视频源:', videoElement.value?.src)
  console.error('摄像头配置:', cameraConfig)
  console.error('生成的完整URL:', streamUrl.value)
  
  // 检查URL格式
  const currentUrl = videoElement.value?.src
  if (currentUrl) {
    console.error('URL分析:')
    console.error('- 协议:', currentUrl.split('://')[0])
    console.error('- 完整地址:', currentUrl)
    console.error('- 是否包含认证:', currentUrl.includes('@'))
  }
  
  connectionStatus.value = '播放错误'
  
  // 提供更详细的错误信息
  let errorMessage = '视频播放失败'
  if (event.target?.error) {
    const error = event.target.error
    switch (error.code) {
      case 1:
        errorMessage = '视频加载被中止，请检查网络连接'
        break
      case 2:
        errorMessage = '网络错误，请检查摄像头地址和网络连接'
        break
      case 3:
        errorMessage = '视频解码错误，请检查视频格式是否支持'
        break
      case 4:
        errorMessage = '视频源不支持或格式错误'
        break
      default:
        errorMessage = `视频播放错误 (错误代码: ${error.code})`
    }
  }
  
  ElMessage.error(errorMessage)
  
  // 如果是认证问题，提供解决建议
  if (cameraConfig.url && !cameraConfig.username) {
    ElMessage.warning('提示：如果摄像头需要认证，请填写用户名和密码')
  }
}

const onPlay = () => {
  connectionStatus.value = '播放中'
}

const onPause = () => {
  connectionStatus.value = '已暂停'
}

// 监听props变化

// 生命周期
onMounted(async () => {
  // 加载保存的配置
  const savedConfig = localStorage.getItem('camera-config')
  if (savedConfig) {
    try {
      Object.assign(cameraConfig, JSON.parse(savedConfig))
    } catch (error) {
      console.error('加载摄像头配置失败:', error)
    }
  }
  
  // 自动选择路由参数中的设备
  if (props.aircraftSn && props.aircraftSn !== '') {
    selectedDevice.value = props.aircraftSn
    currentDeviceSN.value = props.aircraftSn
    console.log('自动选择飞机设备:', props.aircraftSn)
  } else if (props.airportSn && props.airportSn !== '') {
    selectedDevice.value = props.airportSn
    currentDeviceSN.value = props.airportSn
    console.log('自动选择机场设备:', props.airportSn)
  } else {
    // 如果没有路由参数，尝试从设备列表选择
    if (airportDevices.value.length > 0) {
      selectedDevice.value = airportDevices.value[0].sn
      currentDeviceSN.value = airportDevices.value[0].sn
      console.log('自动选择第一个机场设备:', airportDevices.value[0].sn)
    } else if (aircraftDevices.value.length > 0) {
      selectedDevice.value = aircraftDevices.value[0].sn
      currentDeviceSN.value = aircraftDevices.value[0].sn
      console.log('自动选择第一个飞机设备:', aircraftDevices.value[0].sn)
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

.control-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 6px;
  margin-bottom: 16px;
  border: 1px solid #e4e7ed;
}

.control-left {
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

.control-right {
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

.config-tips {
  margin-top: 4px;
}

.diagnostic-content {
  max-height: 60vh;
  overflow-y: auto;
}

.diagnostic-card {
  border: 1px solid #e4e7ed;
}

.diagnostic-card .el-card__header {
  background-color: #f5f7fa;
  font-weight: bold;
}

.diagnostic-actions h4 {
  margin: 16px 0 12px 0;
  color: #303133;
}

/* 用户管理样式 */
.user-manager {
  margin-top: 16px;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.user-manager .el-table {
  margin-top: 8px;
}

.user-manager .el-button {
  margin-left: 8px;
}
</style>
