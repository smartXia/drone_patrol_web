 <template>
   <div class="dashboard">
     <div class="dashboard-layout">
       <!-- 左侧主题列表 -->
       <div class="topic-sidebar">
         <el-card class="topic-list-card" shadow="hover">
           <template #header>
             <div class="card-header">
               <span>DJI主题列表</span>
               <el-button 
                 type="text" 
                 size="small" 
                 @click="resetTopicOrder"
                 style="margin-left: auto;"
               >
                 重置排序
               </el-button>
             </div>
           </template>
           
                     <div class="topic-categories" ref="topicCategoriesRef">
            <!-- 动态渲染主题分类 -->
            <div 
              v-for="category in orderedTopicCategories" 
              :key="category.key"
              class="topic-category" 
              draggable="true"
              @dragstart="handleDragStart($event, category.key)"
              @dragover.prevent
              @drop="handleDrop($event, category.key)"
              @dragenter.prevent
              @dragleave="handleDragLeave"
              @dragend="handleDragEnd"
            >
              <div class="category-header">
                <h4>{{ category.name }}</h4>
                <el-icon class="drag-handle"><Rank /></el-icon>
              </div>
              <div class="topic-item" v-for="(topic, index) in category.data" :key="`${category.key}-${index}`">
                <div class="topic-info">
                  <div class="topic-name">{{ topic.name }}</div>
                  <div class="topic-path">{{ formatTopicPath(topic.topic) }}</div>
                  <el-tooltip 
                    :content="topic.description" 
                    placement="top" 
                    :show-after="500"
                    popper-class="topic-tooltip"
                  >
                    <div class="topic-desc">{{ topic.description }}</div>
                  </el-tooltip>
                </div>
                <el-button 
                  type="primary" 
                  size="small"
                  :disabled="!mqttStore.isConnected || subscribedTopics.includes(topic.topic)"
                  @click="subscribeTopic(topic.topic)"
                >
                  {{ subscribedTopics.includes(topic.topic) ? '已订阅' : '订阅' }}
                </el-button>
              </div>
            </div>
          </div>
         </el-card>
       </div>
       
       <!-- 右侧主要内容区域 -->
       <div class="main-content">
         <!-- 连接控制区域 -->
         <el-card class="connection-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>MQTT连接控制</span>
          <el-button type="text" size="small" style="margin-left:auto;" @click="openErrorDoc">错误码文档</el-button>
        </div>
      </template>
      
             <div class="connection-controls">
         <el-button 
           type="primary" 
           :loading="mqttStore.isConnecting"
           :disabled="mqttStore.isConnected"
           @click="handleConnect"
         >
           连接MQTT
         </el-button>
         
                   <el-button 
            type="success" 
            :disabled="!mqttStore.isConnected"
            @click="handleSubscribe"
          >
            订阅主题
          </el-button>
         
         <el-button 
           type="danger" 
           :disabled="!mqttStore.isConnected"
           @click="handleDisconnect"
         >
           断开连接
         </el-button>
         
         <el-button 
           type="warning" 
           @click="clearHistory"
         >
           清空历史
         </el-button>
         <el-button @click="openConnMgr">连接管理</el-button>
         <el-button @click="gotoRedis">Redis 连接</el-button>
         <el-button type="info" @click="runDiagnostics">网络诊断</el-button>
       </div>
       
                <div class="topic-input" style="margin-top: 15px;">
           <el-input
             v-model="topicInput"
             placeholder="输入订阅主题"
             style="width: 400px; margin-right: 10px;"
             @keyup.enter="handleSubscribe"
           />
           <span style="color: #909399; font-size: 12px;">示例: thing/product/8UUXN2B00A00ST/events</span>
         </div>
      
                      <div class="status-display" style="margin-top: 15px;">
           <el-tag 
             :type="mqttStore.isConnected ? 'success' : 'info'"
             style="margin-right: 10px;"
           >
             MQTT状态: {{ mqttStore.isConnected ? '已连接' : '未连接' }}
           </el-tag>
           
                       <el-tag 
              :type="subscribedTopics.length > 0 ? 'success' : 'info'"
            >
              订阅状态: {{ subscribedTopics.length > 0 ? `已订阅 ${subscribedTopics.length} 个主题` : '未订阅' }}
            </el-tag>
         </div>
         
         <!-- 已订阅主题列表 -->
         <div v-if="subscribedTopics.length > 0" class="subscribed-topics" style="margin-top: 10px;">
           <h4 style="margin-bottom: 8px; color: #606266;">已订阅主题:</h4>
           <div class="topic-tags">
             <el-tag
               v-for="(topic, index) in subscribedTopics"
               :key="topic"
               type="success"
               size="small"
               closable
               @close="removeTopic(index)"
               style="margin-right: 8px; margin-bottom: 5px;"
             >
               {{ topic }}
             </el-tag>
           </div>
           <!-- 每个订阅主题的独立卡片 -->
           <div class="subscribed-topic-cards" style="margin-top: 12px; display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 12px;">
             <el-card
               v-for="topic in subscribedTopics"
               :key="`card-${topic}`"
               class="subscribed-topic-card clickable-card"
               shadow="hover"
               @click="handleTopicCardClick(topic)"
             >
               <template #header>
                 <div class="card-header" style="display:flex; align-items:center; gap:8px; width:100%;">
                   <span style="font-weight:600;">{{ topic }}</span>
                   <el-tag v-if="isRecording(topic)" :type="getFilteredTopicMessages(topic).length > 0 ? 'success' : 'info'" size="small">
                     {{ getFilteredTopicMessages(topic).length }} 条
                   </el-tag>
                   <div style="margin-left:auto; display:flex; gap:8px; align-items:center;">
                     <el-button :type="isRecording(topic) ? 'success' : 'info'" size="small" @click.stop="toggleRecording(topic)">
                       {{ isRecording(topic) ? '暂停记录' : '开启记录' }}
                     </el-button>
                     <el-button type="primary" size="small" @click.stop="openTopicDialog(topic)" :disabled="!isRecording(topic)">查看对话</el-button>
                     <el-button type="warning" size="small" @click.stop="clearTopicMessages(topic)">清空</el-button>
                     <el-button type="info" size="small" @click.stop="openTopicCardFull(topic)">放大</el-button>
                     <el-button type="danger" size="small" @click.stop="removeTopic(subscribedTopics.indexOf(topic))">退订</el-button>
                   </div>
                 </div>
               </template>
               <div v-if="isRecording(topic) && getFilteredTopicMessages(topic).length > 0" class="latest-message-preview">
                 <div class="preview-header" style="display:flex; align-items:center; gap:8px; margin-bottom:6px;">
                   <el-tag :type="getMessageDirectionType(getFilteredTopicMessages(topic)[0])" size="small">
                     {{ getMessageDirectionText(getFilteredTopicMessages(topic)[0]) }}
                   </el-tag>
                   <el-tag v-if="getMessageMethod(getFilteredTopicMessages(topic)[0])" :type="getMethodType(getFilteredTopicMessages(topic)[0])" size="small">
                     {{ getMessageMethod(getFilteredTopicMessages(topic)[0]) }}
                   </el-tag>
                   <span style="margin-left:auto; color:#909399; font-size:12px;">{{ formatTime(getFilteredTopicMessages(topic)[0].timestamp) }}</span>
                 </div>
                 <pre class="preview-payload" style="max-height:120px; overflow:auto; white-space:pre-wrap; word-break:break-word; background:#f8f9fa; border:1px solid #ebeef5; border-radius:6px; padding:8px;">
{{ formatMessagePayload(getFilteredTopicMessages(topic)[0].payload) }}
                 </pre>
               </div>
               <div v-else style="color:#909399; font-size:13px;">{{ isRecording(topic) ? '暂无匹配消息' : '未开启记录' }}</div>
             </el-card>
           </div>
         </div>
       
       <div v-if="mqttStore.connectionError" class="error-message">
         <el-alert 
           :title="mqttStore.connectionError" 
           type="error" 
           show-icon 
           :closable="false"
         />
       </div>
    </el-card>

    <!-- 设备状态展示区域 -->
    <div class="device-grid">
      <el-card 
        v-for="(device, deviceId) in mqttStore.deviceStatus" 
        :key="deviceId"
        class="device-card"
        shadow="hover"
        v-if="shouldShowDeviceCard(device)"
      >
        <template #header>
          <div class="device-header">
            <span class="device-id">设备ID: {{ deviceId }}</span>
            <div style="margin-left:auto; display:flex; gap:8px; align-items:center;">
              <el-tag :type="getDeviceStatusType(device)" size="small">{{ getDeviceStatusText(device) }}</el-tag>
              <el-button circle size="small" @click="openDeviceCardFull(deviceId)"><el-icon><FullScreen /></el-icon></el-button>
            </div>
          </div>
        </template>
        
        <div class="device-content">
          <div class="device-info">
            <p><strong>最后更新:</strong> {{ formatTime(device.lastUpdate) }}</p>
            <p><strong>主题:</strong> {{ device.topic }}</p>
          </div>
          
          <div class="device-data">
            <h4>设备数据:</h4>
            <pre>{{ JSON.stringify(device, null, 2) }}</pre>
          </div>
        </div>
      </el-card>
      
      <!-- 无设备状态时的提示 -->
      <el-empty 
        v-if="getVisibleDeviceCount() === 0"
        description="暂无设备状态数据"
      >
        <el-button type="primary" @click="handleConnect">
          连接MQTT获取数据
        </el-button>
      </el-empty>
    </div>

                   <!-- 主题消息对话框 -->
          <div class="topic-messages-container">
            <!-- 删除下方重复的已订阅主题列表，保留为空容器以便未来扩展 -->
          </div>

         <!-- 消息详情对话框 -->
           <el-dialog 
        v-model="messageDetailVisible" 
        title="消息详情" 
        width="60%"
      >
        <div v-if="selectedMessage" class="message-detail">
          <div class="detail-section">
            <h4>基本信息</h4>
            <div class="detail-row">
              <span class="detail-label">时间:</span>
              <span class="detail-value">{{ formatTime(selectedMessage.timestamp) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">设备ID:</span>
              <span class="detail-value">{{ selectedMessage.deviceId || '未知' }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">主题:</span>
              <span class="detail-value">{{ selectedMessage.topic }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">方向:</span>
              <el-tag 
                :type="getMessageDirectionType(selectedMessage)" 
                size="small"
              >
                {{ getMessageDirectionText(selectedMessage) }}
              </el-tag>
            </div>
          </div>
          
          <div class="detail-section" v-if="getMessageMethod(selectedMessage)">
            <h4>方法信息</h4>
            <div class="detail-row">
              <span class="detail-label">方法:</span>
              <el-tag 
                :type="getMethodType(selectedMessage)" 
                size="small"
              >
                {{ getMessageMethod(selectedMessage) }}
              </el-tag>
            </div>
          </div>
          
          <div class="detail-section" v-if="getMessageError(selectedMessage)">
            <h4>错误信息</h4>
            <div class="detail-row">
              <span class="detail-label">错误码:</span>
              <el-tag 
                type="danger" 
                size="small"
              >
                {{ getMessageError(selectedMessage) }}
              </el-tag>
            </div>
          </div>
          
          <div class="detail-section">
            <h4>原始数据</h4>
            <pre class="detail-payload">{{ JSON.stringify(selectedMessage.payload, null, 2) }}</pre>
          </div>
        </div>
      </el-dialog>

     <!-- 聊天式主题消息对话框 -->
     <el-dialog 
       v-model="topicDialogVisible" 
       :title="activeTopicDialog + ' - 消息对话'" 
       width="90%"
       :before-close="closeTopicDialog"
       class="chat-dialog"
     >
       <div v-if="activeTopicDialog" class="chat-container">
         <!-- 聊天头部 -->
         <div class="chat-header">
           <div class="chat-topic-info">
             <el-tag type="primary" size="large">{{ activeTopicDialog }}</el-tag>
             <el-tag 
               :type="getFilteredTopicMessages(activeTopicDialog).length > 0 ? 'success' : 'info'"
               size="large"
               style="margin-left: 10px;"
             >
               共 {{ getFilteredTopicMessages(activeTopicDialog).length }} 条消息
             </el-tag>
           </div>
           <div class="chat-actions">
             <el-input
               :model-value="getTopicSearch(activeTopicDialog)"
               placeholder="搜索关键词"
               size="small"
               clearable
               style="width: 220px; margin-right: 8px;"
               @update:model-value="onTopicDialogSearchUpdate"
               @clear="onTopicDialogSearchClear"
             />
             <el-button 
               type="danger" 
               size="small"
               @click="removeTopic(subscribedTopics.indexOf(activeTopicDialog))"
             >
               取消订阅
             </el-button>
             <el-button 
               type="warning" 
               size="small"
               @click="clearTopicMessages(activeTopicDialog)"
             >
               清空消息
             </el-button>
           </div>
         </div>
         
         <!-- 方法统计区域 -->
        <div v-if="activeTopicDialog" class="method-stats" style="margin: 10px 0; display:flex; align-items:center; flex-wrap: wrap; gap:8px;">
          <span style="color:#606266;">方法统计:</span>
          <el-tag
            v-for="item in getMethodStats(activeTopicDialog)"
            :key="item.method"
            :type="getTopicSearch(activeTopicDialog) === item.method ? 'primary' : 'info'"
            size="small"
            style="cursor: pointer;"
            @click="applyMethodFilter(item.method)"
          >
            {{ item.method }} ({{ item.count }})
          </el-tag>
          <el-button size="small" type="text" @click="clearMethodFilter" v-if="getTopicSearch(activeTopicDialog)">清除</el-button>
        </div>
        
        <!-- 聊天消息区域 -->
         <div class="chat-messages" ref="chatMessagesRef">
           <div 
             v-if="getTopicMessages(activeTopicDialog).length === 0"
             class="no-messages"
           >
             <el-empty description="暂无消息，等待接收..." :image-size="80" />
           </div>
           
           <div 
             v-else
             class="message-list"
           >
                           <div 
                v-for="(message, index) in getFilteredTopicMessages(activeTopicDialog).slice(-50).reverse()"
                :key="`${activeTopicDialog}-${index}`"
                class="message-bubble"
                :class="getMessageDirectionClass(message)"
                @click="showMessageDetail(message)"
              >
                <div class="message-avatar">
                  <el-avatar 
                    :size="32" 
                    :src="getMessageAvatar(message)"
                    :style="{ backgroundColor: getMessageDirectionColor(message) }"
                  >
                    {{ getMessageInitial(message) }}
                  </el-avatar>
                  <div class="message-direction-icon">
                    <el-icon v-if="getMessageDirection(message) === 'up'">
                      <ArrowUp />
                    </el-icon>
                    <el-icon v-else-if="getMessageDirection(message) === 'down'">
                      <ArrowDown />
                    </el-icon>
                  </div>
                </div>
                <div class="message-content">
                  <div class="message-header">
                    <span class="message-sender">{{ message.deviceId || '设备' }}</span>
                    <div class="message-meta">
                      <el-tag 
                        :type="getMessageDirectionType(message)" 
                        size="small"
                        style="margin-right: 8px;"
                      >
                        {{ getMessageDirectionText(message) }}
                      </el-tag>
                      <span class="message-time">{{ formatTime(message.timestamp) }}</span>
                    </div>
                  </div>
                                     <div class="message-body">
                     <div class="message-method" v-if="getMessageMethod(message)">
                       <el-tag 
                         :type="getMethodType(message)" 
                         size="small"
                         class="method-tag"
                       >
                         {{ getMessageMethod(message) }}
                       </el-tag>
                     </div>
                     <div class="message-error" v-if="getMessageError(message)">
                       <el-tag 
                         type="danger" 
                         size="small"
                         class="error-tag"
                       >
                         错误码: {{ getMessageError(message) }}
                       </el-tag>
                     </div>
                     <div class="message-payload">
                       {{ formatMessagePayload(message.payload) }}
                     </div>
                   </div>
                </div>
              </div>
           </div>
         </div>
         
         <!-- 聊天底部 -->
         <div class="chat-footer">
           <div class="chat-status">
             <el-tag 
               :type="mqttStore.isConnected ? 'success' : 'danger'"
               size="small"
             >
               {{ mqttStore.isConnected ? '已连接' : '未连接' }}
             </el-tag>
             <span class="status-text">
               {{ mqttStore.isConnected ? '正在监听消息...' : '连接已断开' }}
             </span>
           </div>
         </div>
       </div>
     </el-dialog>
       </div>
     </div>
   </div>
   
   <!-- 错误码文档抽屉（移入模板内部） -->
   <el-drawer 
     v-model="errorDocVisible" 
     title="错误码文档"
     direction="ltr"
     size="480px"
   >
     <el-tabs v-model="errorDocTab" type="card">
       <el-tab-pane label="DJI 错误码" name="dji">
         <div style="height: calc(100vh - 160px); overflow: hidden; display:flex; flex-direction: column;">
           <iframe 
             src="/docs/dji-error-codes.md" 
             style="flex:1; border:none; width:100%;" 
           ></iframe>
         </div>
       </el-tab-pane>
       <el-tab-pane label="HMS 错误码" name="hms">
         <div style="height: calc(100vh - 160px); overflow: hidden; display:flex; flex-direction: column;">
           <iframe 
             src="/hms.json" 
             style="flex:1; border:none; width:100%;" 
           ></iframe>
         </div>
       </el-tab-pane>
     </el-tabs>
   </el-drawer>
   <ConnectionManager ref="connMgrRef" />
 </template>

 <script setup>
 import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useMqttStore } from '@/stores/mqtt'
import dayjs from 'dayjs'
import { Rank, ArrowUp, ArrowDown } from '@element-plus/icons-vue'
import { FullScreen } from '@element-plus/icons-vue'
import '@/styles/dashboard.css'
import ConnectionManager from '../components/mqtt/ConnectionManager.vue'
import { useRouter, useRoute } from 'vue-router'

 const mqttStore = useMqttStore()
 const messageDetailVisible = ref(false)
 const selectedMessage = ref(null)
 const topicInput = ref('thing/product/8UUXN2B00A00ST/events')
 const subscribedTopics = ref([]) // 已订阅的主题列表
 const topicDialogVisible = ref(false) // 主题消息对话框显示状态
 const activeTopicDialog = ref('') // 当前激活的主题对话框
 const chatMessagesRef = ref(null) // 聊天消息区域引用
 const topicCategoriesRef = ref(null) // 主题分类容器引用
 const connMgrRef = ref(null)
 const router = useRouter()
 const route = useRoute()
 const gotoRedis = () => { router.push({ name: 'RedisExplorer' }) }

 onMounted(() => {
   if (route.query && (route.query.openMqtt === '1' || route.query.openMqtt === 1)) {
     openConnMgr()
   }
 })
 
 // 拖拽相关状态
 const draggedCategory = ref(null)
 const draggedOverCategory = ref(null)
 const originalTopicOrder = ref([]) // 保存原始主题顺序
 
   // 主题分类顺序 - 先使用默认值，后面再初始化
  const topicCategoryOrder = ref([
    'deviceOsdTopics',
    'deviceServiceTopics', 
    'deviceEventTopics',
    'deviceRequestTopics',
    'systemStatusTopics',
    'devicePropertyTopics',
    'deviceDrcTopics'
  ])
  
  // 计算属性：根据保存的顺序返回主题分类
  const orderedTopicCategories = computed(() => {
    const order = topicCategoryOrder.value
    const categories = {
      // deviceOsdTopics: { name: '设备OSD', data: deviceOsdTopics.value },
      deviceServiceTopics: { name: '设备服务', data: deviceServiceTopics.value },
      deviceEventTopics: { name: '设备事件', data: deviceEventTopics.value },
      deviceRequestTopics: { name: '设备请求', data: deviceRequestTopics.value },
      systemStatusTopics: { name: '系统状态', data: systemStatusTopics.value },
      devicePropertyTopics: { name: '设备属性', data: devicePropertyTopics.value },
      deviceDrcTopics: { name: 'DRC协议', data: deviceDrcTopics.value }
    }
    
    return order.map(categoryKey => ({
      key: categoryKey,
      name: categories[categoryKey]?.name || categoryKey,
      data: categories[categoryKey]?.data || []
    }))
  })

             // DJI主题定义 - 核心主题
    const deviceOsdTopics = ref([
      {
        name: '设备OSD',
        topic: 'thing/product/{device_sn}/osd',
        description: '设备端定频向云平台推送的设备属性（properties）'
      },
      {
        name: '设备状态',
        topic: 'thing/product/{device_sn}/state',
        description: '设备端按需上报向云平台推送的设备属性（properties）'
      }
    ])

    const deviceServiceTopics = ref([
      {
        name: '服务调用',
        topic: 'thing/product/{gateway_sn}/services',
        description: '云平台向设备发送的服务'
      },
      {
        name: '服务响应',
        topic: 'thing/product/{gateway_sn}/services_reply',
        description: '设备对service的回复、处理结果'
      }
    ])

    const deviceEventTopics = ref([
      {
        name: '事件通知',
        topic: 'thing/product/{gateway_sn}/events',
        description: '设备端向云平台发送的，需要关注和处理的事件'
      },
      {
        name: '事件响应',
        topic: 'thing/product/{gateway_sn}/events_reply',
        description: '云平台对设备事件的回复、处理结果'
      }
    ])

    const deviceRequestTopics = ref([
      {
        name: '设备请求',
        topic: 'thing/product/{gateway_sn}/requests',
        description: '设备端向云平台发送请求，为了获取一些信息'
      },
      {
        name: '请求响应',
        topic: 'thing/product/{gateway_sn}/requests_reply',
        description: '云平台对设备请求的回复'
      }
    ])

    const systemStatusTopics = ref([
      {
        name: '系统状态',
        topic: 'sys/product/{gateway_sn}/status',
        description: '设备上下线、更新拓扑'
      },
      {
        name: '状态响应',
        topic: 'sys/product/{gateway_sn}/status_reply',
        description: '平台响应'
      }
    ])

    const devicePropertyTopics = ref([
      {
        name: '属性设置',
        topic: 'thing/product/{gateway_sn}/property/set',
        description: '设备属性设置'
      },
      {
        name: '属性设置响应',
        topic: 'thing/product/{gateway_sn}/property/set_reply',
        description: '设备属性设置的响应'
      }
    ])

    const deviceDrcTopics = ref([
      {
        name: 'DRC上行',
        topic: 'thing/product/{gateway_sn}/drc/up',
        description: 'DRC协议上行'
      },
      {
        name: 'DRC下行',
        topic: 'thing/product/{gateway_sn}/drc/down',
        description: 'DRC协议下行'
      }
    ])

// 记录开关：哪些主题开启了消息记录
const recordingTopics = ref(new Set())

// 每个主题的搜索关键字
const topicSearch = ref({})

const getTopicSearch = (topic) => {
  return (topicSearch.value && topic && topicSearch.value[topic]) ? topicSearch.value[topic] : ''
}

const setTopicSearch = (topic, value) => {
  const next = { ...(topicSearch.value || {}) }
  next[topic] = value
  topicSearch.value = next
}

// 获取过滤后的主题消息列表
const getFilteredTopicMessages = (topic) => {
  const list = getTopicMessages(topic)
  const kw = getTopicSearch(topic).trim().toLowerCase()
  if (!kw) return list
  return list.filter(msg => {
    try {
      const payloadStr = typeof msg.payload === 'string' ? msg.payload : JSON.stringify(msg.payload)
      const haystack = [
        msg.deviceId || '',
        msg.topic || '',
        getMessageMethod(msg) || '',
        getMessageDirectionText(msg) || '',
        payloadStr || ''
      ].join(' ').toLowerCase()
      return haystack.includes(kw)
    } catch (e) {
      return false
    }
  })
}

const isRecording = (topic) => {
  return recordingTopics.value.has(topic)
}

const toggleRecording = (topic) => {
  if (recordingTopics.value.has(topic)) {
    recordingTopics.value.delete(topic)
  } else {
    recordingTopics.value.add(topic)
  }
}

// 连接MQTT
const handleConnect = async () => {
  try {
    await mqttStore.connect()
    ElMessage.success('MQTT连接成功')
  } catch (error) {
    ElMessage.error(`连接失败: ${error.message}`)
  }
}

               // 订阅主题
        const handleSubscribe = async () => {
          try {
            const topic = topicInput.value.trim()
            if (!topic) {
              ElMessage.warning('请输入主题')
              return
            }
            if (subscribedTopics.value.includes(topic)) {
              ElMessage.info('该主题已订阅')
              topicInput.value = ''
              return
            }
            
            await mqttStore.subscribeToTopics(topic)
            subscribedTopics.value.push(topic)
            recordingTopics.value.add(topic)
            topicInput.value = '' // 清空输入框
            ElMessage.success(`成功订阅主题: ${topic}`)
          } catch (error) {
            ElMessage.error(`订阅失败: ${error.message}`)
          }
        }

               // 断开连接
        const handleDisconnect = () => {
          mqttStore.disconnect()
          subscribedTopics.value = []
          ElMessage.info('MQTT连接已断开')
        }
        
                 // 删除主题
         const removeTopic = async (index) => {
           const topic = subscribedTopics.value[index]
           try {
             await mqttStore.unsubscribeTopics(topic)
           } catch (e) {
             console.error('取消订阅失败:', e)
           }
           subscribedTopics.value.splice(index, 1)
           ElMessage.info(`已移除主题: ${topic}`)
         }
         
                   // 订阅指定主题
          const subscribeTopic = async (topicTemplate) => {
            try {
              // 将模板中的 {device_sn} 和 {gateway_sn} 替换为实际的设备ID
              const topic = topicTemplate
                .replace(/\{device_sn\}/g, '8UUXN2B00A00ST')
                .replace(/\{gateway_sn\}/g, '8UUXN2B00A00ST')
              
              if (subscribedTopics.value.includes(topic)) {
                ElMessage.warning('该主题已订阅')
                return
              }
              
              await mqttStore.subscribeToTopics(topic)
              subscribedTopics.value.push(topic)
              recordingTopics.value.add(topic)
              ElMessage.success(`成功订阅主题: ${topic}`)
            } catch (error) {
              ElMessage.error(`订阅失败: ${error.message}`)
            }
          }

// 清空历史记录
const clearHistory = () => {
  mqttStore.clearMessageHistory()
  ElMessage.success('历史记录已清空')
}



// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '未知'
  return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
}

// 获取设备状态类型
const getDeviceStatusType = (device) => {
  if (!device) return 'info'
  
  // 根据设备数据判断状态
  if (device.status === 'online' || device.connected) return 'success'
  if (device.status === 'offline' || device.error) return 'danger'
  if (device.status === 'warning') return 'warning'
  
  return 'info'
}

// 获取设备状态文本
const getDeviceStatusText = (device) => {
  if (!device) return '未知'
  
  if (device.status) return device.status
  if (device.connected) return '在线'
  if (device.error) return '错误'
  
  return '正常'
}

 // 显示消息详情
 const showMessageDetail = (message) => {
   selectedMessage.value = message
   messageDetailVisible.value = true
 }
 
 // 获取指定主题的消息
 const getTopicMessages = (topic) => {
   return mqttStore.messageHistory.filter(message => message.topic === topic)
 }
 
 // 获取指定主题的消息数量
 const getTopicMessageCount = (topic) => {
   return getTopicMessages(topic).length
 }

 // 获取主题的最新一条消息
 const getLatestMessage = (topic) => {
   const list = getTopicMessages(topic)
   if (!list || list.length === 0) return null
   return list[list.length - 1]
 }
 
 // 格式化消息内容
 const formatMessagePayload = (payload) => {
   try {
     if (typeof payload === 'string') {
       const parsed = JSON.parse(payload)
       return JSON.stringify(parsed, null, 2)
     }
     return JSON.stringify(payload, null, 2)
   } catch (error) {
     return payload
   }
 }

 // 格式化主题路径显示
 const formatTopicPath = (topicTemplate) => {
   // 将模板中的 {device_sn} 和 {gateway_sn} 替换为通配符 *
   return topicTemplate.replace(/\{device_sn\}/g, '*').replace(/\{gateway_sn\}/g, '*')
 }

 // 打开主题消息对话框
 const openTopicDialog = (topic) => {
   activeTopicDialog.value = topic
   topicDialogVisible.value = true
 }

 // 关闭主题消息对话框
 const closeTopicDialog = () => {
   topicDialogVisible.value = false
   activeTopicDialog.value = ''
 }

 // 清空指定主题的消息
 const clearTopicMessages = (topic) => {
   // 从消息历史中移除指定主题的消息
   mqttStore.messageHistory = mqttStore.messageHistory.filter(
     message => message.topic !== topic
   )
   ElMessage.success(`已清空主题 ${topic} 的消息`)
 }

 // 获取消息头像
 const getMessageAvatar = (message) => {
   // 可以根据设备ID或消息类型返回不同的头像
   return ''
 }

 // 获取消息头像初始字母
 const getMessageInitial = (message) => {
   const deviceId = message.deviceId || '设备'
   return deviceId.charAt(0).toUpperCase()
 }

 // 获取消息方向
 const getMessageDirection = (message) => {
   const topic = message.topic || ''
   if (topic.includes('/up')) return 'up'
   if (topic.includes('/down')) return 'down'
   return 'unknown'
 }

 // 获取消息方向文本
 const getMessageDirectionText = (message) => {
   const direction = getMessageDirection(message)
   switch (direction) {
     case 'up': return '上行'
     case 'down': return '下行'
     default: return '未知'
   }
 }

 // 获取消息方向类型（用于标签颜色）
 const getMessageDirectionType = (message) => {
   const direction = getMessageDirection(message)
   switch (direction) {
     case 'up': return 'success'
     case 'down': return 'warning'
     default: return 'info'
   }
 }

 // 获取消息方向颜色
 const getMessageDirectionColor = (message) => {
   const direction = getMessageDirection(message)
   switch (direction) {
     case 'up': return '#67c23a'
     case 'down': return '#e6a23c'
     default: return '#909399'
   }
 }

 // 获取消息气泡样式类
 const getMessageDirectionClass = (message) => {
   const direction = getMessageDirection(message)
   switch (direction) {
     case 'up': return 'message-up'
     case 'down': return 'message-down'
     default: return 'message-unknown'
   }
 }

 // 获取消息方法
 const getMessageMethod = (message) => {
   try {
     const payload = typeof message.payload === 'string' ? JSON.parse(message.payload) : message.payload
     
     // 检查不同的方法字段
     if (payload.method) return payload.method
     if (payload.service) return payload.service
     if (payload.event) return payload.event
     if (payload.request) return payload.request
     if (payload.action) return payload.action
     if (payload.type) return payload.type
     
     return null
   } catch (error) {
     return null
   }
 }

 // 获取方法类型（用于标签颜色）
 const getMethodType = (message) => {
   const method = getMessageMethod(message)
   if (!method) return 'info'
   
   // 根据方法名称判断类型
   const methodLower = method.toLowerCase()
   
   // 服务相关方法
   if (methodLower.includes('service') || methodLower.includes('call')) return 'primary'
   
   // 事件相关方法
   if (methodLower.includes('event') || methodLower.includes('notify')) return 'success'
   
   // 请求相关方法
   if (methodLower.includes('request') || methodLower.includes('get')) return 'warning'
   
   // 设置相关方法
   if (methodLower.includes('set') || methodLower.includes('update')) return 'info'
   
   // 控制相关方法
   if (methodLower.includes('control') || methodLower.includes('command')) return 'danger'
   
   return 'info'
 }

 // 获取消息错误码
 const getMessageError = (message) => {
   try {
     const payload = typeof message.payload === 'string' ? JSON.parse(message.payload) : message.payload
     
     // 检查不同的错误字段
     if (payload.error) return payload.error
     if (payload.errorCode) return payload.errorCode
     if (payload.code && payload.code !== 200) return payload.code
     if (payload.status && payload.status !== 'success') return payload.status
     if (payload.result && payload.result.error) return payload.result.error
     if (payload.data && payload.data.error) return payload.data.error
     
     // 检查是否有错误信息
     if (payload.message && payload.message.toLowerCase().includes('error')) return payload.message
     if (payload.msg && payload.msg.toLowerCase().includes('error')) return payload.msg
     
     return null
   } catch (error) {
     return null
   }
 }

 // 自动滚动到最新消息
 const scrollToBottom = () => {
   nextTick(() => {
     if (chatMessagesRef.value) {
       chatMessagesRef.value.scrollTop = chatMessagesRef.value.scrollHeight
     }
   })
 }

 // 拖拽开始
 const handleDragStart = (event, categoryName) => {
   draggedCategory.value = categoryName
   event.dataTransfer.effectAllowed = 'move'
   event.dataTransfer.setData('text/plain', categoryName)
   event.target.style.opacity = '0.5'
 }

 // 拖拽结束
 const handleDragEnd = (event) => {
   event.target.style.opacity = '1'
   draggedCategory.value = null
   draggedOverCategory.value = null
 }

 // 拖拽离开
 const handleDragLeave = (event) => {
   if (event.target.classList.contains('topic-category')) {
     event.target.style.borderColor = ''
     event.target.style.backgroundColor = ''
   }
 }

 // 拖拽放置
 const handleDrop = (event, targetCategoryName) => {
   event.preventDefault()
   
   if (draggedCategory.value && draggedCategory.value !== targetCategoryName) {
     // 重新排序主题分类
     reorderTopicCategories(draggedCategory.value, targetCategoryName)
   }
   
   // 重置样式
   event.target.style.borderColor = ''
   event.target.style.backgroundColor = ''
 }

 // 重新排序主题分类
 const reorderTopicCategories = (fromCategory, toCategory) => {
   // 获取当前的主题分类顺序
   const categoryOrder = getCurrentCategoryOrder()
   
   const fromIndex = categoryOrder.indexOf(fromCategory)
   const toIndex = categoryOrder.indexOf(toCategory)
   
   if (fromIndex !== -1 && toIndex !== -1 && fromIndex !== toIndex) {
     // 移除拖拽的分类
     categoryOrder.splice(fromIndex, 1)
     // 在目标位置插入
     categoryOrder.splice(toIndex, 0, fromCategory)
     
     // 更新主题分类顺序
     updateTopicCategoryOrder(categoryOrder)
     
     ElMessage.success('主题分类顺序已更新')
   }
 }

 // 获取当前主题分类顺序
 const getCurrentCategoryOrder = () => {
   const savedOrder = localStorage.getItem('topicCategoryOrder')
   if (savedOrder) {
     try {
       return JSON.parse(savedOrder)
     } catch (error) {
       console.error('解析保存的主题顺序失败:', error)
     }
   }
   
   // 默认顺序
   return [
     'deviceOsdTopics',
     'deviceServiceTopics', 
     'deviceEventTopics',
     'deviceRequestTopics',
     'systemStatusTopics',
     'devicePropertyTopics',
     'deviceDrcTopics'
   ]
 }

 // 更新主题分类顺序
 const updateTopicCategoryOrder = (newOrder) => {
   // 更新响应式变量
   topicCategoryOrder.value = [...newOrder]
   
   // 保存到localStorage
   localStorage.setItem('topicCategoryOrder', JSON.stringify(newOrder))
   
   console.log('新的主题分类顺序:', newOrder)
 }

 // 重置主题顺序
 const resetTopicOrder = () => {
   const defaultOrder = [
     'deviceOsdTopics',
     'deviceServiceTopics', 
     'deviceEventTopics',
     'deviceRequestTopics',
     'systemStatusTopics',
     'devicePropertyTopics',
     'deviceDrcTopics'
   ]
   
   updateTopicCategoryOrder(defaultOrder)
   localStorage.removeItem('topicCategoryOrder')
   ElMessage.success('主题分类顺序已重置')
 }

 // 加载保存的主题顺序
 const loadTopicCategoryOrder = () => {
   const savedOrder = localStorage.getItem('topicCategoryOrder')
   if (savedOrder) {
     try {
       const order = JSON.parse(savedOrder)
       topicCategoryOrder.value = [...order]
     } catch (error) {
       console.error('加载主题顺序失败:', error)
     }
   }
 }

 // 仅展示特定主题来源的设备卡片
const shouldShowDeviceCard = (device) => {
  if (!device || !device.topic) return false
  const t = device.topic
  // 仅显示 osd/state/sys 状态等，不显示 events/requests/drc/property 等
  if (t.includes('/osd')) return true
  if (t.includes('/state')) return true
  if (t.startsWith('sys/')) return true
  return false
}

// 获取可见设备卡片数量
const getVisibleDeviceCount = () => {
  const all = mqttStore.deviceStatus
  if (!all) return 0
  return Object.values(all).filter(d => shouldShowDeviceCard(d)).length
}

const enlargedTopicVisible = ref(false)
const enlargedTopic = ref('')
const openTopicCardFull = (topic) => {
  enlargedTopic.value = topic
  enlargedTopicVisible.value = true
}
const closeTopicCardFull = () => {
  enlargedTopicVisible.value = false
  enlargedTopic.value = ''
}

const enlargedDeviceVisible = ref(false)
const enlargedDeviceId = ref('')
const openDeviceCardFull = (deviceId) => {
  enlargedDeviceId.value = deviceId
  enlargedDeviceVisible.value = true
}
const closeDeviceCardFull = () => {
  enlargedDeviceVisible.value = false
  enlargedDeviceId.value = ''
}

// 处理已订阅主题卡片点击
const handleTopicCardClick = (topic) => {
  if (!isRecording(topic)) return
  openTopicDialog(topic)
}

// 方法统计：统计某主题内各method出现次数
const getMethodStats = (topic) => {
  if (!topic) return []
  const list = getTopicMessages(topic)
  const counter = {}
  list.forEach((msg) => {
    const m = getMessageMethod(msg) || '未知'
    counter[m] = (counter[m] || 0) + 1
  })
  return Object.entries(counter)
    .map(([method, count]) => ({ method, count }))
    .sort((a, b) => b.count - a.count)
}

// 快速按method筛选/清除
const applyMethodFilter = (method) => {
  if (!activeTopicDialog.value) return
  setTopicSearch(activeTopicDialog.value, method || '')
}
const clearMethodFilter = () => {
  if (!activeTopicDialog.value) return
  setTopicSearch(activeTopicDialog.value, '')
}

const onTopicDialogSearchUpdate = (val) => {
  if (!activeTopicDialog.value) return
  setTopicSearch(activeTopicDialog.value, val)
}
const onTopicDialogSearchClear = () => {
  if (!activeTopicDialog.value) return
  setTopicSearch(activeTopicDialog.value, '')
}

const errorDocVisible = ref(false)
const errorDocTab = ref('dji')
const openErrorDoc = () => {
  errorDocTab.value = 'dji'
  errorDocVisible.value = true
}
const closeErrorDoc = () => {
  errorDocVisible.value = false
}

const openConnMgr = () => { connMgrRef.value && connMgrRef.value.open() }

// 网络诊断功能
const runDiagnostics = async () => {
  if (!mqttStore.currentConfig) {
    ElMessage.warning('请先配置MQTT连接信息')
    return
  }

  const config = mqttStore.currentConfig
  ElMessage.info('开始网络诊断...')
  
  try {
    // 检查网络连通性
    console.log('开始网络连通性检查...', { host: config.host, port: config.port })
    const networkResult = await mqttStore.checkNetworkConnectivity(config.host, config.port)
    console.log('网络连通性检查结果:', networkResult)
    
    if (networkResult.success) {
      ElMessage.success(`网络连通性检查通过 - ${config.host}:${config.port} (耗时: ${networkResult.duration}ms)`)
    } else {
      ElMessage.error(`网络连通性检查失败 - ${networkResult.error || '未知错误'}`)
    }
    
    // 测试MQTT连接
    try {
      console.log('开始MQTT连接测试...', config)
      await mqttStore.testConnection({ config })
      ElMessage.success('MQTT连接测试通过')
    } catch (mqttError) {
      console.error('MQTT连接测试失败:', mqttError)
      ElMessage.error(`MQTT连接测试失败 - ${mqttError.message}`)
    }
    
  } catch (error) {
    console.error('诊断过程中出现错误:', error)
    ElMessage.error(`诊断过程中出现错误: ${error.message}`)
  }
}

// 组件挂载时自动连接
onMounted(() => {
  // 加载保存的主题顺序
  loadTopicCategoryOrder()
})

// 组件卸载时断开连接
onUnmounted(() => {
  // 可以选择是否在组件卸载时断开连接
  // mqttStore.disconnect()
})
 </script>
