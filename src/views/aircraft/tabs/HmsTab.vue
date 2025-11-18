<template>
  <div style="padding: 8px 0; color: var(--el-text-color-regular); font-size: 13px; display: flex; flex-direction: column; gap: 10px;">
    
    <!-- HMS健康告警监控 -->
    <el-card shadow="never">
      <template #header>
        <div style="display: flex; align-items: center; justify-content: space-between; cursor: pointer;" @click="toggleHmsCollapse">
          <div style="display: flex; align-items: center; gap: 8px;">
            <el-icon :style="{ transform: hmsCollapsed ? 'rotate(0deg)' : 'rotate(90deg)', transition: 'transform 0.3s' }">
              <ArrowRight />
            </el-icon>
            <span style="font-weight: 600; font-size: 14px;">HMS健康告警监控</span>
            <div v-if="!hmsCollapsed" style="display: flex; align-items: center; gap: 8px;">
              <el-tag size="small" type="success">实时监控</el-tag>
              <el-tag size="small" type="info">Topic: thing/product/{gateway_sn}/events</el-tag>
              <el-tag size="small" type="warning">Direction: up</el-tag>
              <el-tag size="small" type="primary">Method: hms</el-tag>
            </div>
          </div>
        </div>
      </template>
      
      <el-collapse-transition>
        <div v-show="!hmsCollapsed">
          <!-- 方法说明 -->
          <div style="margin-bottom: 15px; padding: 10px; background: var(--el-bg-color-page); border-radius: 4px; border: 1px solid var(--el-border-color-light);">
            <div style="color: var(--el-text-color-regular); font-size: 13px; line-height: 1.6;">
              <div style="margin-bottom: 8px;"><strong>方法说明：</strong></div>
              <div><strong>Topic:</strong> thing/product/{gateway_sn}/events</div>
              <div><strong>Direction:</strong> up (设备端向服务端发送)</div>
              <div><strong>Method:</strong> hms</div>
              <div style="margin-top: 8px; color: var(--el-text-color-placeholder); font-size: 12px;">
                请手动订阅此Topic来接收实时的健康告警通知
              </div>
              
              <!-- 手动订阅按钮 -->
              <div style="margin-top: 10px;">
                <el-button 
                  type="success" 
                  @click="subscribeHmsEvents"
                  :loading="hmsSubscribing"
                  size="small"
                >
                  <el-icon><Connection /></el-icon>
                  手动订阅HMS告警
                </el-button>
              </div>
            </div>
          </div>
          
          <!-- 检测状态信息 -->
          <div style="margin-bottom: 15px; padding: 10px; background: var(--el-bg-color-page); border-radius: 4px; border: 1px solid var(--el-border-color-light);">
            <div style="margin-bottom: 8px;">
              <div style="color: var(--el-text-color-regular); font-size: 13px; font-weight: 600;">
                健康检测状态 - 最新检测: {{ getLastCheckTime() }}
              </div>
            </div>
            <div style="display: flex; gap: 15px; font-size: 12px;">
              <div>
                <span style="color: var(--el-text-color-regular);">告警数量:</span>
                <span style="color: var(--el-color-danger); font-weight: 600;">{{ hmsData.length }}</span>
              </div>
              <div>
                <span style="color: var(--el-text-color-regular);">检测状态:</span>
                <el-tag :type="hmsData.length > 0 ? 'danger' : 'success'" size="small">
                  {{ hmsData.length > 0 ? '有告警' : '正常' }}
                </el-tag>
              </div>
            </div>
          </div>
          
          <!-- 告警数据表格 -->
          <div v-if="hmsData.length > 0" style="margin-bottom: 15px;">
            <div style="margin-bottom: 8px; color: var(--el-text-color-regular); font-size: 13px;">
              健康告警列表 ({{ hmsData.length }}条)：
            </div>
            <el-table 
              :data="hmsData" 
              size="small" 
              height="300" 
              border 
              stripe
              style="margin-bottom: 15px;"
            >
              <el-table-column prop="level" label="告警等级" width="100">
                <template #default="scope">
                  <el-tag 
                    :type="getLevelType(scope.row.level)" 
                    size="small"
                  >
                    {{ getLevelText(scope.row.level) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="module" label="事件模块" width="100">
                <template #default="scope">
                  <el-tag 
                    :type="getModuleType(scope.row.module)" 
                    size="small"
                  >
                    {{ getModuleText(scope.row.module) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="in_the_sky" label="飞行状态" width="100">
                <template #default="scope">
                  <el-tag 
                    :type="scope.row.in_the_sky ? 'warning' : 'info'" 
                    size="small"
                  >
                    {{ scope.row.in_the_sky ? '在天上' : '在地上' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="code" label="告警码" width="200">
                <template #default="scope">
                  <div style="font-size: 12px; cursor: pointer;" @click="showAlarmDetail(scope.row)">
                    <div style="font-family: monospace; color: var(--el-text-color-regular); margin-bottom: 2px;">
                      {{ scope.row.code }}
                    </div>
                    <div style="color: var(--el-text-color-primary); font-weight: 500;">
                      {{ getAlarmDescription(scope.row.code) }}
                    </div>
                    <div style="color: var(--el-text-color-placeholder); font-size: 10px; margin-top: 2px;">
                      点击查看详情
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="device_type" label="设备类型" width="120">
                <template #default="scope">
                  <span style="font-size: 12px;">{{ scope.row.device_type }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="imminent" label="及时性" width="80">
                <template #default="scope">
                  <el-tag 
                    :type="scope.row.imminent ? 'danger' : 'success'" 
                    size="small"
                  >
                    {{ scope.row.imminent ? '是' : '否' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="args" label="参数" width="120">
                <template #default="scope">
                  <div style="font-size: 11px;">
                    <div>组件: {{ scope.row.args?.component_index || '--' }}</div>
                    <div>传感器: {{ scope.row.args?.sensor_index || '--' }}</div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="timestamp" label="告警时间" width="150">
                <template #default="scope">
                  <div style="font-size: 11px; color: var(--el-text-color-regular);">
                    {{ formatAlarmTime(scope.row.timestamp) }}
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
          
          <!-- 调试信息 -->
          <div style="margin-top: 15px; padding: 10px; background: var(--el-bg-color-page); border-radius: 4px; border: 1px solid var(--el-border-color-light);">
            <div style="font-size: 12px; color: var(--el-text-color-regular); margin-bottom: 8px;">
              <strong>调试信息：</strong>
            </div>
            <div style="font-size: 11px; color: var(--el-text-color-placeholder); line-height: 1.4;">
              <div>hmsPayload: {{ hmsPayload ? '有数据' : '无数据' }}</div>
              <div>告警数据: {{ hmsData.length }}条</div>
              <div v-if="hmsPayload">hmsPayload类型: {{ typeof hmsPayload }}</div>
            </div>
          </div>
          
          <!-- 完整HMS消息展示 -->
          <div v-if="hmsPayload && hmsData.length > 0" style="margin-top: 15px;">
            <div style="margin-bottom: 8px; color: var(--el-text-color-regular); font-size: 13px;">
              完整HMS消息：
            </div>
            <pre style="background: var(--el-bg-color-page); padding: 10px; border-radius: 4px; border: 1px solid var(--el-border-color-light); font-size: 11px; max-height: 300px; overflow-y: auto;">{{ formatHmsPayload() }}</pre>
          </div>
          
          <!-- 无告警提示 -->
          <div v-if="hmsPayload && hmsData.length === 0" style="margin-top: 15px; text-align: center; padding: 20px; color: var(--el-text-color-regular); font-size: 12px;">
            <div style="margin-bottom: 10px;">
              <el-icon size="48" color="var(--el-color-success)">
                <Check />
              </el-icon>
            </div>
            <div style="font-size: 14px; margin-bottom: 8px; color: var(--el-color-success);">系统健康状态正常</div>
            <div style="color: var(--el-text-color-placeholder);">暂无健康告警</div>
          </div>
          
          <!-- 等待数据提示 -->
          <div v-if="!hmsPayload" style="margin-top: 15px; text-align: center; padding: 20px; color: var(--el-text-color-placeholder); font-size: 12px;">
            等待接收HMS健康检测数据...
          </div>
        </div>
      </el-collapse-transition>
    </el-card>
    
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { ElButton, ElTag, ElCard, ElTable, ElTableColumn, ElIcon, ElCollapseTransition } from 'element-plus'
import { ElMessage } from 'element-plus'
import { ArrowRight, Connection, Check } from '@element-plus/icons-vue'
import { useDeviceStore } from '@/stores/device'
import { useMqttProxyStore } from '@/stores/mqtt-proxy'

const props = defineProps({
  airportSn: String,
  hmsPayload: [Object, String],
  hmsData: Array
})

const emit = defineEmits(['subscribeHms'])

// 折叠状态
const hmsCollapsed = ref(false)
const hmsSubscribing = ref(false)

// HMS告警码映射数据
const hmsMapping = ref({})

// 切换折叠状态
const toggleHmsCollapse = () => {
  hmsCollapsed.value = !hmsCollapsed.value
}

// 加载HMS告警码映射数据
const loadHmsMapping = async () => {
  try {
    const response = await fetch('/docs/hms.json')
    if (response.ok) {
      hmsMapping.value = await response.json()
      console.log('✅ HMS告警码映射数据加载成功')
    } else {
      console.error('❌ 加载HMS映射数据失败:', response.status)
    }
  } catch (error) {
    console.error('❌ 加载HMS映射数据失败:', error)
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadHmsMapping()
})

// 手动订阅HMS告警
const subscribeHmsEvents = async () => {
  try {
    hmsSubscribing.value = true
    
    // 使用传入的机场序列号
    const gatewaySn = props.airportSn
    
    if (!gatewaySn || gatewaySn === '--') {
      ElMessage.error('缺少网关SN（机场SN）')
      return
    }
    
    // 获取MQTT代理store
    const mqttProxyStore = useMqttProxyStore()
    
    if (!mqttProxyStore.isConnected) {
      ElMessage.error('MQTT未连接，请先连接MQTT服务')
      return
    }
    
    // 订阅HMS告警Topic
    const hmsTopic = `thing/product/${gatewaySn}/events`
    console.log('=== 手动订阅HMS告警Topic ===')
    console.log('Topic:', hmsTopic)
    console.log('Gateway SN:', gatewaySn)
    
    // 通过emit通知父组件
    emit('subscribeHms', { topic: hmsTopic, gatewaySn })
    
    ElMessage.success('已发送订阅请求，请查看控制台确认订阅结果')
    console.log('📡 订阅请求已发送，等待父组件处理...')
    
  } catch (error) {
    console.error('订阅HMS告警失败:', error)
    ElMessage.error('订阅失败: ' + (error?.message || '未知错误'))
  } finally {
    hmsSubscribing.value = false
  }
}

// 获取告警等级类型
const getLevelType = (level) => {
  const typeMap = {
    0: 'info',    // 通知
    1: 'warning', // 提醒
    2: 'danger'   // 警告
  }
  return typeMap[level] || 'info'
}

// 获取告警等级文本
const getLevelText = (level) => {
  const textMap = {
    0: '通知',
    1: '提醒', 
    2: '警告'
  }
  return textMap[level] || '未知'
}

// 获取事件模块类型
const getModuleType = (module) => {
  const typeMap = {
    0: 'primary', // 飞行任务
    1: 'success', // 设备管理
    2: 'warning', // 媒体
    3: 'info'      // hms
  }
  return typeMap[module] || 'info'
}

// 获取事件模块文本
const getModuleText = (module) => {
  const textMap = {
    0: '飞行任务',
    1: '设备管理',
    2: '媒体',
    3: 'HMS'
  }
  return textMap[module] || '未知'
}

// 获取最新检测时间
const getLastCheckTime = () => {
  if (props.hmsPayload) {
    try {
      const data = typeof props.hmsPayload === 'string' 
        ? JSON.parse(props.hmsPayload) 
        : props.hmsPayload
      
      if (data.timestamp) {
        return new Date(data.timestamp).toLocaleString()
      }
    } catch (e) {
      console.error('解析HMS数据失败:', e)
    }
  }
  return '未检测'
}

// 获取告警码详细信息
const getAlarmDetails = (code) => {
  if (!code || !hmsMapping.value) return null
  
  // 尝试不同的匹配方式
  const possibleKeys = [
    code, // 直接匹配
    `dock_tip_${code}`, // dock_tip_前缀
    `aircraft_tip_${code}`, // aircraft_tip_前缀
    `payload_tip_${code}`, // payload_tip_前缀
    `remote_controller_tip_${code}` // remote_controller_tip_前缀
  ]
  
  for (const key of possibleKeys) {
    if (hmsMapping.value[key]) {
      return hmsMapping.value[key]
    }
  }
  
  return null
}

// 获取告警码描述
const getAlarmDescription = (code) => {
  const details = getAlarmDetails(code)
  if (details) {
    return details.zh || details.en || '未知告警'
  }
  return '未知告警'
}

// 格式化HMS消息
const formatHmsPayload = () => {
  if (!props.hmsPayload) return ''
  
  try {
    const data = typeof props.hmsPayload === 'string' 
      ? JSON.parse(props.hmsPayload) 
      : props.hmsPayload
    
    return JSON.stringify(data, null, 2)
  } catch (e) {
    console.error('格式化HMS消息失败:', e)
    return String(props.hmsPayload)
  }
}

// 格式化告警时间
const formatAlarmTime = (timestamp) => {
  if (!timestamp) return '--'
  
  try {
    // 如果是数字时间戳
    if (typeof timestamp === 'number') {
      return new Date(timestamp).toLocaleString()
    }
    // 如果是字符串时间戳
    if (typeof timestamp === 'string') {
      const numTimestamp = parseInt(timestamp)
      if (!isNaN(numTimestamp)) {
        return new Date(numTimestamp).toLocaleString()
      }
    }
    return '--'
  } catch (e) {
    console.error('格式化告警时间失败:', e)
    return '--'
  }
}

// 显示告警详情
const showAlarmDetail = (alarm) => {
  const details = getAlarmDetails(alarm.code)
  if (details) {
    ElMessage({
      message: `告警码: ${alarm.code}\n中文描述: ${details.zh || '无'}\n英文描述: ${details.en || '无'}`,
      type: 'info',
      duration: 5000,
      showClose: true
    })
  } else {
    ElMessage.warning('未找到该告警码的详细信息')
  }
}
</script>

<style scoped>
</style>


