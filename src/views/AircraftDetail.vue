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
        <!-- 功能菜单栏 -->
        <el-card shadow="never" style="margin-bottom: 12px;">
          <el-tabs v-model="activeFeatureTab" type="border-card">
            <el-tab-pane label="机场信息" name="airport-detail">
              <AirportInfoTab 
                v-if="activeFeatureTab === 'airport-detail'"
                :device-sn="deviceSn"
                :airport-sn="getAirportSn()"
                @data-received="handleAirportDataReceived"
                @subscription-status="handleAirportSubscriptionStatus"
              />
            </el-tab-pane>
            <el-tab-pane label="飞机信息" name="aircraft-detail">
              <AircraftInfoTab 
                v-if="activeFeatureTab === 'aircraft-detail'"
                :device-sn="deviceSn"
                :aircraft-sn="getAircraftSn()"
                @data-received="handleAircraftDataReceived"
                @subscription-status="handleAircraftSubscriptionStatus"
              />
            </el-tab-pane>
            <el-tab-pane label="直播信息" name="live-stream">
              <CameraLiveStream 
                v-if="activeFeatureTab === 'live-stream'"
                :airport-sn="getAirportSn()"
                :aircraft-sn="getAircraftSn()"
              />
            </el-tab-pane>
            <el-tab-pane label="远程日志" name="remote-log">
              <RemoteLogTab 
                :sending="logSending"
                :airportSn="getAirportSn()"
                :requestTopic="logRequestTopic"
                :replyTopic="logReplyTopic"
                :requestPayload="logRequestPayload"
                :replyPayload="logReplyPayload"
                :uploadSending="uploadSending"
                :uploadRequestTopic="uploadRequestTopic"
                :uploadReplyTopic="uploadReplyTopic"
                :uploadRequestPayload="uploadRequestPayload"
                :uploadReplyPayload="uploadReplyPayload"
                :progressPayload="progressPayload"
                :allEventsMessages="allEventsMessages"
                @request="quickFetchFileList"
                @upload="handleFileUpload"
                @statusUpdate="handleStatusUpdate"
                @subscribeProgress="handleSubscribeProgress"
              />
            </el-tab-pane>
            <el-tab-pane label="计划库任务记录" name="task-records">
              <div style="padding: 8px 0; display:flex; justify-content: space-between; align-items:center;">
                <div style="color: var(--el-text-color-regular); font-size: 13px;">任务下发后的本地记录（后端接入后可替换为真实数据）</div>
                <div style="display:flex; gap:8px;">
                  <el-button size="small" @click="refreshTaskRecords">刷新记录</el-button>
                </div>
              </div>
              <el-table :data="taskRecords" size="small" height="360" border>
                <el-table-column prop="taskId" label="任务ID" width="160" />
                <el-table-column prop="method" label="方法" width="140" />
                <el-table-column prop="tid" label="TID" width="220" />
                <el-table-column prop="gatewaySn" label="机场SN" width="160" />
                <el-table-column label="请求参数">
                  <template #default="scope">
                    <span style="font-family: monospace;">{{ briefJson(scope.row.payload) }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="result" label="结果码" width="100" />
                <el-table-column label="时间" width="180">
                  <template #default="scope">{{ formatTime(scope.row.createdAt) }}</template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="远程升级功能" name="remote-upgrade">
              <RemoteUpgradeTab 
                :device-sn="deviceSn"
                :osd-data="osdData"
                :upgrade-status="upgradeStatus"
                :upgrade-status-payload="upgradeStatusPayload"
                @subscribe-upgrade="handleSubscribeUpgrade"
                @start-upgrade="handleStartUpgrade"
                @check-status="handleCheckUpgradeStatus"
              />
            </el-tab-pane>
            <el-tab-pane label="HMS 功能" name="hms">
              <HmsTab 
                :airportSn="getAirportSn()"
                :hmsPayload="hmsPayload"
                :hmsData="hmsData"
                @subscribeHms="handleSubscribeHms"
              />
            </el-tab-pane>
            <el-tab-pane label="远程调试" name="remote-debug">
              <div style="padding: 8px 0; color: var(--el-text-color-regular); font-size: 13px;">
                在此处集成远程调试操作（远程控制、参数查询/设置等）。
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useMqttProxyStore } from '@/stores/mqtt-proxy'
import { useDeviceStore } from '@/stores/device'

// 导入组件
import AirportInfoTab from '@/views/aircraft/tabs/AirportInfoTab.vue'
import AircraftInfoTab from '@/views/aircraft/tabs/AircraftInfoTab.vue'
import RemoteLogTab from '@/views/aircraft/tabs/RemoteLogTab.vue'
import RemoteUpgradeTab from '@/views/aircraft/tabs/RemoteUpgradeTab.vue'
import HmsTab from '@/views/aircraft/tabs/HmsTab.vue'
import CameraLiveStream from '@/components/CameraLiveStream.vue'

// 路由
const route = useRoute()
const router = useRouter()

// 状态管理
const mqttProxyStore = useMqttProxyStore()
const deviceStore = useDeviceStore()

// 响应式数据
const deviceSn = ref(route.params.deviceSn || route.params.aircraftSn || '--')
const activeFeatureTab = ref('airport-detail')

// 基本状态数据
const osdData = ref(null)
const upgradeStatus = ref(null)
const upgradeStatusPayload = ref(null)
const hmsPayload = ref(null)
const hmsData = ref([])

// 远程日志相关
const logSending = ref(false)
const logRequestTopic = ref('')
const logReplyTopic = ref('')
const logRequestPayload = ref(null)
const logReplyPayload = ref(null)

// 文件上传相关
const uploadSending = ref(false)
const uploadRequestTopic = ref('')
const uploadReplyTopic = ref('')
const uploadRequestPayload = ref(null)
const uploadReplyPayload = ref(null)
const progressPayload = ref(null)
const allEventsMessages = ref([])

// 任务记录
const taskRecords = ref([])

// 方法
const goBack = () => {
  router.go(-1)
}

// 获取机场序列号
const getAirportSn = () => {
  const routeAirportSn = route.params.airportSn
  return routeAirportSn || deviceStore.currentDevice?.airport_sn || deviceSn.value || '--'
}

// 获取飞机序列号
const getAircraftSn = () => {
  const routeAircraftSn = route.params.aircraftSn
  return routeAircraftSn || deviceSn.value || '--'
}

// 计算属性 - 获取设备列表
const deviceList = computed(() => deviceStore.deviceList)
const airportDevices = computed(() => 
  deviceList.value.filter(device => device.type === 'airport' || device.airport_sn)
)
const aircraftDevices = computed(() => 
  deviceList.value.filter(device => device.type === 'aircraft' && device.airport_sn)
)

// 事件处理器
const handleAirportDataReceived = (data) => {
  console.log('机场数据接收事件:', data)
  osdData.value = data
}

const handleAirportSubscriptionStatus = (status) => {
  console.log('机场订阅状态事件:', status)
}

const handleAircraftDataReceived = (data) => {
  console.log('飞机数据接收事件:', data)
}

const handleAircraftSubscriptionStatus = (status) => {
  console.log('飞机订阅状态事件:', status)
}

const handleSubscribeHms = () => {
  console.log('订阅HMS事件')
}

const handleSubscribeUpgrade = () => {
  console.log('订阅升级事件')
}

const handleStartUpgrade = () => {
  console.log('开始升级事件')
}

const handleCheckUpgradeStatus = () => {
  console.log('检查升级状态事件')
}

const quickFetchFileList = () => {
  console.log('快速获取文件列表')
}

const handleFileUpload = () => {
  console.log('文件上传事件')
}

const handleStatusUpdate = () => {
  console.log('状态更新事件')
}

const handleSubscribeProgress = () => {
  console.log('订阅进度事件')
}

const refreshTaskRecords = () => {
  console.log('刷新任务记录')
}

const formatTime = (timestamp) => {
  if (!timestamp) return '--'
  return new Date(timestamp).toLocaleString()
}

const briefJson = (obj) => {
  if (!obj) return '--'
  return JSON.stringify(obj).substring(0, 50) + '...'
}

// 生命周期
onMounted(async () => {
  console.log('AircraftDetail页面已挂载')
  
  try {
    // 加载设备列表
    await deviceStore.loadDevices()
    console.log('设备列表已加载:', deviceStore.deviceList.length, '个设备')
    
    // 如果MQTT未连接，尝试使用默认连接
    if (!mqttProxyStore.isConnected) {
      console.log('MQTT未连接，尝试使用默认连接')
      try {
        await mqttProxyStore.connect()
        console.log('MQTT默认连接成功')
      } catch (error) {
        console.error('MQTT默认连接失败:', error)
        ElMessage.warning('MQTT连接失败，请检查网络连接')
      }
    }
    
    // 等待MQTT连接稳定后自动订阅首个机场和飞机
    setTimeout(async () => {
      await autoSubscribeFirstDevices()
    }, 1000)
    
  } catch (error) {
    console.error('页面初始化失败:', error)
  }
})

// 自动订阅首个设备
const autoSubscribeFirstDevices = async () => {
  try {
    // 等待MQTT连接完成
    if (!mqttProxyStore.isConnected) {
      console.log('等待MQTT连接...')
      return
    }
    
    // 获取首个机场设备
    const firstAirport = airportDevices.value[0]
    if (firstAirport) {
      console.log('自动订阅首个机场:', firstAirport.sn)
      try {
        // 订阅机场OSD数据
        const airportOsdTopic = `thing/product/${firstAirport.sn}/osd`
        await mqttProxyStore.subscribeToTopics(airportOsdTopic, 1)
        console.log('机场OSD订阅成功:', airportOsdTopic)
  } catch (error) {
        console.error('机场OSD订阅失败:', error)
      }
    }
    
    // 获取首个飞机设备
    const firstAircraft = aircraftDevices.value[0]
    if (firstAircraft) {
      console.log('自动订阅首个飞机:', firstAircraft.sn)
      try {
        // 订阅飞机OSD数据
        const aircraftOsdTopic = `thing/product/${firstAircraft.sn}/osd`
        await mqttProxyStore.subscribeToTopics(aircraftOsdTopic, 1)
        console.log('飞机OSD订阅成功:', aircraftOsdTopic)
  } catch (error) {
        console.error('飞机OSD订阅失败:', error)
      }
    }
    
    // 如果当前没有设备SN，使用首个设备
    if (deviceSn.value === '--' && firstAircraft) {
      deviceSn.value = firstAircraft.sn
      console.log('自动设置当前设备SN:', firstAircraft.sn)
    }
    
    // 显示成功消息
    if (firstAirport || firstAircraft) {
      ElMessage.success('已自动订阅首个设备数据')
    }
    
  } catch (error) {
    console.error('自动订阅设备失败:', error)
    ElMessage.warning('自动订阅设备失败，请手动订阅')
  }
}
</script>

<style scoped>
.aircraft-detail {
  height: 100vh;
  background: var(--el-bg-color-page);
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.detail-main {
  padding: 20px;
}
</style>
