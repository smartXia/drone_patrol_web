<template>
  <div class="dashboard">
    <el-container style="height: 100vh;">
      <!-- 左侧面板 -->
      <el-aside :width="sidebarCollapsed ? '0px' : '350px'" class="dashboard-aside" :class="{ 'collapsed': sidebarCollapsed }">
        <!-- 收缩状态下的展开按钮 -->
        <div v-if="sidebarCollapsed" class="sidebar-collapsed-indicator">
          <el-button 
            type="primary" 
            size="default" 
            @click="toggleSidebar"
            class="expand-button"
            :title="'展开侧边栏'"
          >
            <el-icon><ArrowRight /></el-icon>
            展开
          </el-button>
        </div>
        
        <div class="sidebar-content" v-show="!sidebarCollapsed">
          <!-- 主题列表 -->
          <el-card class="topic-list-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>DJI主题列表</span>
                <div class="card-header-actions">
                  <el-button 
                    type="text" 
                    size="small" 
                    @click="resetTopicOrder"
                    style="font-size: 10px;"
                  >
                    重置排序
                  </el-button>
                  <el-button 
                    type="text" 
                    size="small" 
                    @click="toggleSidebar"
                    style="font-size: 10px; margin-left: 8px;"
                  >
                    {{ sidebarCollapsed ? '展开' : '收缩' }}
                  </el-button>
                </div>
              </div>
            </template>
           
                     <div class="topic-list" ref="topicCategoriesRef">
            <!-- 所有主题统一显示 -->
            <div class="topic-item" v-for="(topic, index) in allTopics" :key="`topic-${index}`">
              <el-button 
                :type="subscribedTopics.includes(topic.topic) ? 'success' : 'primary'"
                size="small"
                :disabled="!mqttProxyStore.isConnected"
                @click="subscribedTopics.includes(topic.topic) ? unsubscribeTopic(topic.topic) : subscribeTopic(topic.topic)"
                style="margin-right: 12px; flex-shrink: 0; font-size: 10px; padding: 4px 8px;"
              >
                {{ subscribedTopics.includes(topic.topic) ? '已订阅' : '订阅' }}
              </el-button>
              <div class="topic-info" style="flex: 1;">
                <div class="topic-path">{{ formatTopicPath(topic.topic) }}</div>
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
           :loading="mqttProxyStore.isConnecting"
           :disabled="mqttProxyStore.isConnected"
           @click="handleConnect"
         >
           连接MQTT
         </el-button>
         
                   <el-button 
            type="success" 
            :disabled="!mqttProxyStore.isConnected"
            @click="handleSubscribe"
          >
            订阅主题
          </el-button>
         
         <el-button 
           type="warning" 
           :disabled="!mqttProxyStore.isConnected || subscribedTopics.length === 0"
           @click="clearAllSubscriptions"
         >
           清空订阅
         </el-button>
         
         <el-button 
           type="danger" 
           :disabled="!mqttProxyStore.isConnected"
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
       </div>
       
                <div class="topic-input" style="margin-top: 15px;">
           <el-input
             v-model="topicInput"
             placeholder="输入订阅主题"
             style="width: 400px; margin-right: 10px;"
             @keyup.enter="handleSubscribe"
           />
           <span style="color: #909399; font-size: 12px;">示例: thing/product/{{ (deviceStore.currentDevice?.airport_sn || deviceStore.currentDeviceSn) || 'DEVICE_SN' }}/events</span>
         </div>
      
                      <div class="status-display" style="margin-top: 15px;">
           <el-tag 
             :type="mqttProxyStore.isConnected ? 'success' : 'info'"
             style="margin-right: 10px;"
           >
             MQTT状态: {{ mqttProxyStore.isConnected ? '已连接' : '未连接' }}
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
               {{ topic }} (QoS {{ getTopicQos(topic) }})
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
                   <el-tag type="primary" size="small">
                     QoS {{ getTopicQos(topic) }}
                   </el-tag>
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
         </el-card>
       </div>
      </el-aside>

      <!-- 主内容区域 -->
      <el-container>
        <el-header class="dashboard-header">
          <div class="header-left">
            <el-button 
              :icon="sidebarCollapsed ? 'Expand' : 'Fold'" 
              @click="toggleSidebar"
              class="collapse-btn"
              type="primary"
              size="default"
              :title="sidebarCollapsed ? '展开侧边栏' : '收缩侧边栏'"
            >
              {{ sidebarCollapsed ? '展开' : '收缩' }}
            </el-button>
            <h2>设备监控面板</h2>
            
            <!-- 当前设备选择器 -->
            <div class="current-device-selector">
              <el-select
                v-model="deviceStore.currentDeviceSn"
                :placeholder="deviceStore.deviceList.length > 0 ? '选择监控设备' : '请先添加设备'"
                @change="handleDeviceChange"
                clearable
                :disabled="deviceStore.deviceList.length === 0"
                style="width: 200px; margin-left: 20px;"
              >
                <el-option
                  v-for="device in deviceStore.deviceList"
                  :key="device.sn"
                  :label="`${device.name} (${device.airport_sn || device.sn})`"
                  :value="device.airport_sn || device.sn"
                >
                  <div class="device-option">
                    <el-tag 
                      :type="device.status === 'online' ? 'success' : 'info'"
                      size="small"
                      style="margin-right: 8px;"
                    >
                      {{ device.status === 'online' ? '在线' : '离线' }}
                    </el-tag>
                    <span>{{ device.name }}</span>
                    <span style="color: #909399; font-size: 12px; margin-left: 8px;">{{ device.airport_sn || device.sn }}</span>
                    <span v-if="device.airport_sn && device.sn !== device.airport_sn" style="color: #67c23a; font-size: 10px; margin-left: 4px;">(机场)</span>
                  </div>
                </el-option>
              </el-select>
              
              <!-- 设备列表为空时的提示 -->
              <el-tooltip 
                v-if="deviceStore.deviceList.length === 0" 
                content="点击设备管理按钮添加设备"
                placement="bottom"
              >
                <el-button 
                  type="warning" 
                  size="small" 
                  @click="toggleDeviceManager"
                  style="margin-left: 8px;"
                >
                  <el-icon><Plus /></el-icon>
                  添加设备
                </el-button>
              </el-tooltip>
              
              <!-- 当前设备信息显示 -->
              <div v-if="deviceStore.currentDevice" class="current-device-info">
                <el-tag 
                  :type="deviceStore.currentDevice.status === 'online' ? 'success' : 'info'"
                  size="small"
                >
                  {{ deviceStore.currentDevice.status === 'online' ? '在线' : '离线' }}
                </el-tag>
                <span class="device-name">{{ deviceStore.currentDevice.name }}</span>
                <span class="device-sn">{{ deviceStore.currentDevice.airport_sn || deviceStore.currentDevice.sn }}</span>
                <span v-if="deviceStore.currentDevice.airport_sn && deviceStore.currentDevice.sn !== deviceStore.currentDevice.airport_sn" style="color: #67c23a; font-size: 10px; margin-left: 4px;">(机场)</span>
                <div v-if="deviceStore.currentDevice.airport_sn && deviceStore.currentDevice.sn !== deviceStore.currentDevice.airport_sn" style="font-size: 10px; color: #909399; margin-top: 2px;">
                  飞机SN: {{ deviceStore.currentDevice.sn }}
                </div>
              </div>
              
              <!-- MQTT连接信息显示 -->
              <div class="mqtt-connection-info">
                <el-tag 
                  :type="mqttProxyStore.connectionStatus === 'connected' ? 'success' : 
                         mqttProxyStore.connectionStatus === 'connecting' ? 'warning' : 'info'"
                  size="small"
                >
                  {{ mqttProxyStore.connectionStatus === 'connected' ? '已连接' : 
                     mqttProxyStore.connectionStatus === 'connecting' ? '连接中' : '未连接' }}
                </el-tag>
                <span class="mqtt-profile-name">{{ mqttProxyStore.currentMQTTProfileName }}</span>
                <span v-if="mqttProxyStore.currentConfig" class="mqtt-server-info">
                  {{ mqttProxyStore.currentConfig.host }}:{{ mqttProxyStore.currentConfig.port }}
                </span>
              </div>
            </div>
          </div>
          <div class="header-right">
            <el-button type="success" @click="toggleConnectionManager">
              <el-icon><Setting /></el-icon>
              MQTT 连接管理
            </el-button>
            <el-button type="primary" @click="toggleDeviceManager">
              <el-icon><Monitor /></el-icon>
              设备管理
            </el-button>
            <el-button 
              type="warning" 
              @click="openAircraftDetail(deviceStore.currentDeviceSn)"
              :disabled="!deviceStore.currentDeviceSn"
            >
              <el-icon><Monitor /></el-icon>
              飞机详情
            </el-button>
            <el-button type="success" @click="openErrorCodeDialog">
              <el-icon><Document /></el-icon>
              DJI错误码查询
            </el-button>
            <el-button type="info" @click="openCameraLive">
              <el-icon><VideoCamera /></el-icon>
              机场直播
            </el-button>
          </div>
        </el-header>

        <el-main class="dashboard-main">
          <div v-if="mqttProxyStore.connectionError" class="error-message">
            <el-alert 
              :title="mqttProxyStore.connectionError" 
              type="error" 
              show-icon 
              :closable="false"
            />
          </div>

          <!-- 主要内容区域：设备监控 + 错误码查询 -->
          <div class="main-content">
            <!-- 设备状态展示区域 -->
            <div class="device-section">
              <div class="device-grid">
      <el-card 
        v-for="[topic, topicInfo] in allDisplayTopics" 
        :key="topic"
        class="device-card clickable-card"
        shadow="hover"
        @click="openTopicDialog(topic)"
      >
        <template #header>
          <div class="device-header">
            <span class="device-id">主题: {{ topic }}</span>
            <div style="margin-left:auto; display:flex; gap:6px; align-items:center;">
              <el-tag type="primary" size="small" style="font-size: 10px;">QoS {{ topicInfo.qos }}</el-tag>
              <el-tag :type="topicInfo.hasData ? 'success' : 'info'" size="small" style="font-size: 10px;">
                {{ topicInfo.hasData ? '有数据' : '等待数据' }}
              </el-tag>
              <el-button circle size="small" @click.stop="openDeviceCardFull(topic)" style="font-size: 10px;"><el-icon><FullScreen /></el-icon></el-button>
            </div>
          </div>
        </template>
        
        <div class="device-content">
          <div class="device-info">
            <p><strong>最后更新:</strong> {{ formatTime(topicInfo.timestamp) }}</p>
            <p><strong>状态:</strong> {{ topicInfo.hasData ? '已收到数据' : '等待数据' }}</p>
          </div>
          
          <div class="device-data" v-if="topicInfo.hasData && topicInfo.device">
            <h4>设备数据:</h4>
            <pre>{{ JSON.stringify(topicInfo.device, null, 2) }}</pre>
          </div>
          <div class="device-data" v-else>
            <h4>等待数据中...</h4>
            <p style="color: #909399;">该主题已订阅，等待接收消息</p>
          </div>
        </div>
      </el-card>
      
      <!-- 无设备状态时的提示 -->
      <el-empty 
        v-if="getVisibleDeviceCount() === 0"
        :description="deviceStore.deviceList.length === 0 ? '请先添加设备' : '暂无设备状态数据'"
      >
        <div v-if="deviceStore.deviceList.length === 0">
          <p style="color: #909399; margin-bottom: 16px;">要查看飞机详情，请先添加设备：</p>
          <el-button type="primary" @click="toggleDeviceManager">
            <el-icon><Plus /></el-icon>
            添加设备
          </el-button>
        </div>
        <div v-else>
          <el-button type="primary" @click="handleConnect">
            连接MQTT获取数据
          </el-button>
        </div>
      </el-empty>
              </div>
            </div>

          </div>

          <!-- 主题消息对话框 -->
          <div class="topic-messages-container">
            <!-- 删除下方重复的已订阅主题列表，保留为空容器以便未来扩展 -->
          </div>
        </el-main>
      </el-container>
    </el-container>

    <!-- MQTT 连接管理抽屉 -->
    <ConnectionManager ref="connMgrRef" :drawer-mode="true" />
    
    <!-- 设备管理抽屉 -->
    <el-drawer
      v-model="deviceManagerVisible"
      title="设备管理"
      direction="rtl"
      size="600px"
      :before-close="handleDeviceManagerClose"
    >
      <div style="padding: 20px;">
        <DeviceSelector />
      </div>
    </el-drawer>

    <!-- 错误码查询抽屉 -->
    <el-drawer
      v-model="errorCodePageVisible"
      title="DJI错误码查询"
      direction="rtl"
      size="60%"
      :before-close="handleCloseErrorCodePage"
      class="error-code-drawer"
    >
      <div class="error-code-page">
        <!-- 页面头部搜索 -->
        <div class="page-header">
          <div class="search-controls">
            <el-input
              v-model="pageSearchQuery"
              placeholder="搜索所有错误码..."
              clearable
              @input="handlePageSearch"
              class="page-search-input"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-select
              v-model="pageDataSource"
              placeholder="选择数据源"
              @change="handlePageDataSourceChange"
              style="width: 150px; margin-left: 12px;"
            >
              <el-option label="全部" value="all" />
              <el-option label="DJI错误码" value="dji" />
              <el-option label="HMS错误码" value="hms" />
            </el-select>
            <el-button 
              type="primary" 
              size="default"
              @click="showAllErrorCodes"
              style="margin-left: 12px;"
            >
              <el-icon><View /></el-icon>
              查看全部
            </el-button>
          </div>
          <div class="page-stats">
            <el-tag type="info">DJI: {{ djiCodes.length }}</el-tag>
            <el-tag type="success" style="margin-left: 8px;">HMS: {{ hmsCodes.length }}</el-tag>
            <el-tag type="warning" style="margin-left: 8px;">显示: {{ pageFilteredCodes.length }}</el-tag>
          </div>
        </div>

        <!-- 错误码列表 -->
        <div class="error-codes-page-list">
          <div 
            v-for="errorCode in paginatedCodes" 
            :key="errorCode.id"
            class="error-code-page-item"
            @click="showErrorCodePageDetail(errorCode)"
          >
            <div class="error-code-page-header">
              <span class="error-code-page-code">{{ errorCode.code }}</span>
              <div class="error-code-page-tags">
                <el-tag :type="errorCode.type === 'dji' ? getErrorCodeType(errorCode.level) : 'info'" size="small">
                  {{ errorCode.type === 'dji' ? errorCode.level : 'HMS' }}
                </el-tag>
                <el-tag :type="errorCode.type === 'dji' ? 'primary' : 'success'" size="small" style="margin-left: 4px;">
                  {{ errorCode.type === 'dji' ? 'DJI' : 'HMS' }}
                </el-tag>
              </div>
            </div>
            <div class="error-code-page-content">
              <div v-if="errorCode.type === 'dji'" class="error-code-page-description">
                {{ errorCode.description }}
              </div>
              <div v-else class="error-code-page-description">
                <div class="error-code-page-zh">{{ errorCode.zh }}</div>
                <div class="error-code-page-en" v-if="errorCode.en">{{ errorCode.en }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div class="pagination-container" v-if="pageFilteredCodes.length > pageSize">
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            :total="pageFilteredCodes.length"
            layout="prev, pager, next, total"
            @current-change="handlePageChange"
          />
        </div>
      </div>
    </el-drawer>


     <!-- 聊天式主题消息对话框 -->
     <el-dialog 
       v-model="topicDialogVisible" 
       :title="activeTopicDialog + ' (QoS ' + getTopicQos(activeTopicDialog) + ') - 消息对话'" 
       width="90%"
       :before-close="closeTopicDialog"
       class="chat-dialog"
     >
       <div v-if="activeTopicDialog" class="chat-container">
         <!-- 聊天头部 -->
         <div class="chat-header">
           <div class="chat-topic-info">
             <el-tag type="primary" size="large">{{ activeTopicDialog }}</el-tag>
             <el-tag type="warning" size="large" style="margin-left: 10px;">
               QoS {{ getTopicQos(activeTopicDialog) }}
             </el-tag>
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
                style="cursor: pointer;"
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
                      <el-tag 
                        type="info" 
                        size="small"
                        style="margin-right: 8px;"
                      >
                        QoS {{ message.qos || 0 }}
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
          <div class="service-sender" style="display:flex; align-items:flex-start; gap:8px; margin-bottom:8px;">
            <el-input 
              v-model="serviceMethod" 
              size="small" 
              placeholder="方法名，点击右侧按钮选择" 
              readonly
              style="width: 220px; cursor: pointer;"
              @click="openMethodPanel"
            >
              <template #append>
                <el-button @click="openMethodPanel" :icon="Setting" />
              </template>
            </el-input>
            <el-input 
              v-model="serviceParams" 
              size="small" 
              type="textarea" 
              :rows="2" 
              placeholder='JSON参数，例如 {"module_list":["0","3"]}' 
              style="flex:1;"
            />
            <el-button 
              type="primary" 
              size="small" 
              :loading="sendingService" 
              @click="handleSendService"
            >发送</el-button>
          </div>
           <div class="chat-status">
             <el-tag 
               :type="mqttProxyStore.isConnected ? 'success' : 'danger'"
               size="small"
             >
               {{ mqttProxyStore.isConnected ? '已连接' : '未连接' }}
             </el-tag>
             <span class="status-text">
               {{ mqttProxyStore.isConnected ? '正在监听消息...' : '连接已断开' }}
             </span>
           </div>
         </div>
       </div>
     </el-dialog>
   
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

   <!-- 方法选择弹出面板 -->
   <el-dialog 
     v-model="methodPanelVisible" 
     title="选择方法" 
     width="900px"
     :close-on-click-modal="false"
     @close="closeMethodPanel"
   >
     <div class="method-panel">
       <!-- 搜索框和自定义输入 -->
       <div style="margin-bottom: 20px;">
         <el-input
           v-model="methodSearchQuery"
           placeholder="搜索方法或输入自定义方法名"
           clearable
           size="default"
           @input="handleMethodSearch"
         >
           <template #prefix>
             <el-icon><Search /></el-icon>
           </template>
         </el-input>
       </div>

       <!-- 左右布局：左侧方法列表，右侧参数编辑 -->
       <div class="method-panel-layout">
         <!-- 左侧：方法列表 -->
         <div class="method-list-container">
           <!-- 调试方法区域 -->
           <div class="method-section">
             <div class="method-section-header">
               <el-tag type="warning" size="default">调试</el-tag>
             </div>
             <div class="method-list">
               <div
                 v-for="method in filteredDebugMethods"
                 :key="method.value"
                 class="method-item"
                 :class="{ 'method-item-selected': selectedMethodTemp === method.value }"
                 @click="selectMethod(method)"
               >
                 <div class="method-item-content">
                   <div class="method-name">{{ method.value }}</div>
                   <div class="method-label">{{ method.label }}</div>
                 </div>
                 <el-icon v-if="selectedMethodTemp === method.value" class="method-check-icon">
                   <Check />
                 </el-icon>
               </div>
             </div>
           </div>

           <!-- 业务方法区域（仅在有业务方法时显示） -->
           <div v-if="filteredBusinessMethods.length > 0" class="method-section">
             <div class="method-section-header">
               <el-tag type="success" size="default">业务方法</el-tag>
             </div>
             <div class="method-list">
               <div
                 v-for="method in filteredBusinessMethods"
                 :key="method.value"
                 class="method-item"
                 :class="{ 'method-item-selected': selectedMethodTemp === method.value }"
                 @click="selectMethod(method)"
               >
                 <div class="method-item-content">
                   <div class="method-name">{{ method.value }}</div>
                   <div class="method-label">{{ method.label }}</div>
                 </div>
                 <el-icon v-if="selectedMethodTemp === method.value" class="method-check-icon">
                   <Check />
                 </el-icon>
               </div>
             </div>
           </div>

           <!-- 自定义方法名 -->
           <div v-if="methodSearchQuery && !isMethodInList(methodSearchQuery)" class="custom-method-section">
             <div class="method-section-header">
               <el-tag type="info" size="default">自定义方法</el-tag>
             </div>
             <div 
               class="method-item method-item-custom"
               @click="selectCustomMethod(methodSearchQuery)"
             >
               <div class="method-item-content">
                 <div class="method-name">{{ methodSearchQuery }}</div>
                 <div class="method-label">点击使用此方法名</div>
               </div>
               <el-icon class="method-check-icon">
                 <Plus />
               </el-icon>
             </div>
           </div>
         </div>

         <!-- 右侧：参数编辑区域 -->
         <div v-if="selectedMethodTemp" class="params-editor-container">
           <div class="params-editor-section">
             <div class="method-section-header">
               <el-tag type="primary" size="default">方法参数</el-tag>
               <span style="margin-left: 12px; color: #999; font-size: 12px;">{{ getSelectedMethodLabel() }}</span>
             </div>
             <div class="params-editor">
               <el-input
                 v-model="serviceParamsTemp"
                 type="textarea"
                 :rows="12"
                 placeholder='请输入JSON格式的参数，例如：{"module_list":["0","3"]}'
                 style="font-family: monospace; font-size: 13px;"
                 @blur="validateParams"
               />
               <div class="params-actions">
                 <el-button size="small" @click="loadDefaultParams" v-if="hasDefaultParams()">
                   加载默认参数
                 </el-button>
                 <el-button size="small" @click="formatParams">
                   格式化JSON
                 </el-button>
                 <el-button size="small" @click="clearParams">
                   清空
                 </el-button>
                 <span v-if="paramsError" style="color: var(--el-color-error); font-size: 12px; margin-left: 12px;">
                   {{ paramsError }}
                 </span>
                 <span v-else-if="serviceParamsTemp" style="color: var(--el-color-success); font-size: 12px; margin-left: 12px;">
                   ✓ JSON格式正确
                 </span>
               </div>
             </div>
           </div>
         </div>

         <!-- 未选择方法时的提示 -->
         <div v-else class="params-editor-container params-editor-empty">
           <div class="empty-hint">
             <el-icon style="font-size: 48px; color: #ccc; margin-bottom: 12px;"><Document /></el-icon>
             <div style="color: #999; font-size: 14px;">请选择一个方法以编辑参数</div>
           </div>
         </div>
       </div>
     </div>
     <template #footer>
       <div style="display: flex; justify-content: space-between; align-items: center;">
         <span style="color: #999; font-size: 12px;">
           当前选择: <strong>{{ selectedMethodTemp || '未选择' }}</strong>
         </span>
         <div>
           <el-button @click="closeMethodPanel">取消</el-button>
           <el-button type="primary" @click="confirmMethodSelection" :disabled="!selectedMethodTemp">确定</el-button>
         </div>
       </div>
     </template>
   </el-dialog>

   <!-- 消息详情抽屉 -->
   <el-drawer
     v-model="messageDetailVisible"
     title="消息详情"
     direction="rtl"
     size="50%"
     class="message-detail-drawer"
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
         <div class="detail-row">
           <span class="detail-label">QoS:</span>
           <el-tag type="info" size="small">
             {{ selectedMessage.qos || 0 }}
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
         <div class="json-section-header">
           <h4>原始数据</h4>
           <el-button 
             type="primary" 
             size="small" 
             @click="copyJsonPayload"
             :icon="copySuccess ? 'Check' : 'CopyDocument'"
             :class="{ 'copy-success': copySuccess }"
           >
             {{ copySuccess ? '已复制' : '复制' }}
           </el-button>
         </div>
         <div class="json-viewer">
           <pre class="json-content" v-html="formatJsonPayload(selectedMessage.payload)"></pre>
         </div>
       </div>
     </div>
   </el-drawer>
  </div>
 </template>

 <script setup>
 import { ref, onMounted, onUnmounted, nextTick, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useMqttProxyStore } from '@/stores/mqtt-proxy'
import dayjs from 'dayjs'
import { Rank, ArrowUp, ArrowDown, ArrowRight, Setting, Search, Monitor, View, VideoCamera, Document, CopyDocument, Check, Plus } from '@element-plus/icons-vue'
import { FullScreen } from '@element-plus/icons-vue'
import '@/styles/dashboard.css'
import ConnectionManager from '../components/mqtt/ConnectionManager.vue'
import DeviceSelector from '../components/DeviceSelector.vue'
import { useRouter, useRoute } from 'vue-router'
import { useDeviceStore } from '../stores/device'

 const mqttProxyStore = useMqttProxyStore()
 const deviceStore = useDeviceStore()
 
// 侧边栏折叠状态
const sidebarCollapsed = ref(false)

// 设备管理抽屉状态
const deviceManagerVisible = ref(false)
 
const topicInput = ref(`thing/product/${(deviceStore.currentDevice?.airport_sn || deviceStore.currentDeviceSn) || 'DEVICE_SN'}/events`)

// 监听设备变化，更新输入框（不再自动重新订阅主题）
watch(() => deviceStore.currentDeviceSn, (newSn, oldSn) => {
  // 延迟获取设备信息，确保currentDevice已更新
  setTimeout(() => {
    const currentDevice = deviceStore.currentDevice
    const airportSn = currentDevice?.airport_sn || newSn
    topicInput.value = `thing/product/${airportSn || 'DEVICE_SN'}/events`
    console.log('设备切换，更新主题输入框:', {
      newSn,
      airportSn,
      currentDevice: currentDevice?.name,
      topicInput: topicInput.value
    })
  }, 100) // 延迟100ms确保设备信息已更新
  
  // 设备切换时不再自动重新订阅主题
  // 用户需要手动订阅新设备的主题
  if (newSn !== oldSn && oldSn) {
    console.log('设备已切换，请手动订阅新设备的主题')
  }
})
 
// 错误码查询相关
const errorCodeSearchQuery = ref('')
const errorCodes = ref([])
const filteredErrorCodes = ref([])
const errorCodeSource = ref('all') // 'all', 'dji', 'hms'

// DJI 错误码
const djiSearchQuery = ref('')
const djiCodes = ref([])
const filteredDjiCodes = ref([])

// HMS 错误码
const hmsSearchQuery = ref('')
const hmsCodes = ref([])
const filteredHmsCodes = ref([])

// 错误码详情页面相关
const errorCodePageVisible = ref(false)
const pageSearchQuery = ref('')
const pageDataSource = ref('all')
const pageFilteredCodes = ref([])
const currentPage = ref(1)
const pageSize = ref(50)

// 计算属性：是否有搜索结果
const hasSearchResults = computed(() => {
  if (errorCodeSource.value === 'all') {
    return filteredDjiCodes.value.length > 0 || filteredHmsCodes.value.length > 0
  } else if (errorCodeSource.value === 'dji') {
    return filteredDjiCodes.value.length > 0
  } else if (errorCodeSource.value === 'hms') {
    return filteredHmsCodes.value.length > 0
  }
  return false
})

// 分页数据
const paginatedCodes = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return pageFilteredCodes.value.slice(start, end)
})
 // 使用 store 中的 subscribedTopics，转换为数组格式
const subscribedTopics = computed(() => {
  const topics = Array.from(mqttProxyStore.subscribedTopics.keys())
  console.log('已订阅主题列表:', topics)
  console.log('Store中的subscribedTopics:', mqttProxyStore.subscribedTopics)
  return topics
})

// 获取主题的 QoS 信息
const getTopicQos = (topic) => {
  const topicInfo = mqttProxyStore.subscribedTopics.get(topic)
  return topicInfo ? topicInfo.qos : 0
}

// 获取所有需要显示的主题（包括已订阅但还没有消息的主题）
const allDisplayTopics = computed(() => {
  const topics = new Map()
  
  // 添加所有已订阅的主题
  subscribedTopics.value.forEach(topic => {
    const topicInfo = mqttProxyStore.subscribedTopics.get(topic)
    topics.set(topic, {
      topic,
      qos: topicInfo ? topicInfo.qos : 0,
      timestamp: topicInfo ? topicInfo.timestamp : new Date().toISOString(),
      hasData: false
    })
  })
  
  // 添加有设备状态的主题
  Object.entries(mqttProxyStore.deviceStatus).forEach(([deviceId, device]) => {
    if (device.topic) {
      const topicInfo = mqttProxyStore.subscribedTopics.get(device.topic)
      topics.set(device.topic, {
        topic: device.topic,
        qos: topicInfo ? topicInfo.qos : 0,
        timestamp: device.lastUpdate || new Date().toISOString(),
        hasData: true,
        deviceId,
        device
      })
    }
  })
  
  return topics
})

 const topicDialogVisible = ref(false) // 主题消息对话框显示状态
 const activeTopicDialog = ref('') // 当前激活的主题对话框
 const chatMessagesRef = ref(null) // 聊天消息区域引用
 const topicCategoriesRef = ref(null) // 主题分类容器引用
 const connMgrRef = ref(null)
 const router = useRouter()
 const route = useRoute()
 const gotoRedis = () => { router.push({ name: 'RedisExplorer' }) }

// 切换侧边栏折叠状态
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

// 打开连接管理器抽屉
const toggleConnectionManager = () => {
  if (connMgrRef.value) {
    connMgrRef.value.open()
  }
}

// 打开设备管理器抽屉
const toggleDeviceManager = () => {
  deviceManagerVisible.value = true
}

// 关闭设备管理器抽屉
const handleDeviceManagerClose = (done) => {
  deviceManagerVisible.value = false
  done()
}


 onMounted(async () => {
   if (route.query && (route.query.openMqtt === '1' || route.query.openMqtt === 1)) {
     openConnMgr()
   }
   
   // 自动连接MQTT和订阅设备
   await autoConnectAndSubscribe()
 })

// 自动连接MQTT和订阅设备
const autoConnectAndSubscribe = async () => {
  try {
    console.log('Dashboard页面开始自动连接和订阅')
    
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
        return
      }
    }
    
    // 等待MQTT连接稳定后自动订阅首个机场和飞机
    setTimeout(async () => {
      await autoSubscribeFirstDevices()
    }, 1000)
    
  } catch (error) {
    console.error('Dashboard自动连接和订阅失败:', error)
  }
}

// 自动订阅首个设备
const autoSubscribeFirstDevices = async () => {
  try {
    // 等待MQTT连接完成
    if (!mqttProxyStore.isConnected) {
      console.log('等待MQTT连接...')
      return
    }
    
    // 获取设备列表
    const deviceList = deviceStore.deviceList
    const airportDevices = deviceList.filter(device => device.type === 'airport' || device.airport_sn)
    const aircraftDevices = deviceList.filter(device => device.type === 'aircraft' && device.airport_sn)
    
    // 获取首个机场设备
    const firstAirport = airportDevices[0]
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
    const firstAircraft = aircraftDevices[0]
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
    
    // 显示成功消息
    if (firstAirport || firstAircraft) {
      ElMessage.success('已自动订阅首个设备数据')
    }
    
  } catch (error) {
    console.error('自动订阅设备失败:', error)
    ElMessage.warning('自动订阅设备失败，请手动订阅')
  }
}
 
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
      // deviceDrcTopics: { name: 'DRC协议', data: deviceDrcTopics.value }
    }
    
    return order.map(categoryKey => ({
      key: categoryKey,
      name: categories[categoryKey]?.name || categoryKey,
      data: categories[categoryKey]?.data || []
    }))
  })

  // 计算属性：所有主题合并到一个数组
  const allTopics = computed(() => {
    const topics = []
    orderedTopicCategories.value.forEach(category => {
      topics.push(...category.data)
    })
    return topics
  })

             // DJI主题定义 - 核心主题
    const deviceOsdTopics = ref([
      {
        name: '设备OSD',
        topic: 'thing/product/{device_sn}/osd'
      },
      {
        name: '设备状态',
        topic: 'thing/product/{device_sn}/state'
      }
    ])

    const deviceServiceTopics = ref([
      {
        name: '服务调用',
        topic: 'thing/product/{device_sn}/services'
      },
      {
        name: '服务响应',
        topic: 'thing/product/{device_sn}/services_reply'
      }
    ])

    const deviceEventTopics = ref([
      {
        name: '事件通知',
        topic: 'thing/product/{device_sn}/events'
      },
      {
        name: '事件响应',
        topic: 'thing/product/{device_sn}/events_reply'
      }
    ])

    const deviceRequestTopics = ref([
      {
        name: '设备请求',
        topic: 'thing/product/{device_sn}/requests'
      },
      {
        name: '请求响应',
        topic: 'thing/product/{device_sn}/requests_reply'
      }
    ])

    const systemStatusTopics = ref([
      {
        name: '系统状态',
        topic: 'sys/product/{device_sn}/status'
      },
      {
        name: '状态响应',
        topic: 'sys/product/{device_sn}/status_reply'
      }
    ])

    const devicePropertyTopics = ref([
      {
        name: '属性设置',
        topic: 'thing/product/{device_sn}/property/set'
      },
      {
        name: '属性设置响应',
        topic: 'thing/product/{device_sn}/property/set_reply'
      },
      {
        name: '设备定频数据',
        topic: 'thing/product/{device_sn}/osd'
      },
      {
        name: '设备状态数据',
        topic: 'thing/product/{device_sn}/state'
      }
    ])

    const deviceDrcTopics = ref([
      {
        name: 'DRC上行',
        topic: 'thing/product/{device_sn}/drc/up'
      },
      {
        name: 'DRC下行',
        topic: 'thing/product/{device_sn}/drc/down'
      }
    ])

// 记录开关：哪些主题开启了消息记录
const recordingTopics = ref(new Set())

// 监听订阅主题变化，自动开启记录
watch(subscribedTopics, (newTopics) => {
  newTopics.forEach(topic => {
    if (!recordingTopics.value.has(topic)) {
      recordingTopics.value.add(topic)
    }
  })
}, { immediate: true })

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

// 处理设备切换
const handleDeviceChange = async (deviceSn) => {
  if (deviceSn) {
    try {
      // 记录切换前的连接状态和订阅数量
      const wasConnected = mqttProxyStore.isConnected
      const previousSubscriptions = Array.from(subscribedTopics.value)
      
      // 如果之前已连接MQTT，先取消所有订阅
      if (wasConnected && previousSubscriptions.length > 0) {
        try {
          console.log('设备切换前，取消所有之前的订阅...')
          for (const topic of previousSubscriptions) {
            console.log('取消订阅:', topic)
            await mqttProxyStore.unsubscribeTopics(topic)
          }
          ElMessage.info(`已取消 ${previousSubscriptions.length} 个之前的订阅`)
        } catch (unsubscribeError) {
          console.error('取消订阅失败:', unsubscribeError)
          ElMessage.warning('取消部分订阅失败，但继续切换设备')
        }
      }
      
      // 找到对应的设备（通过airport_sn或sn匹配）
      const targetDevice = deviceStore.deviceList.find(device => 
        device.airport_sn === deviceSn || device.sn === deviceSn
      )
      
      if (!targetDevice) {
        ElMessage.error('未找到对应的设备')
        return
      }
      
      // 切换设备（使用设备内部的sn字段）
      await deviceStore.setCurrentDevice(targetDevice.sn)
      ElMessage.success(`已切换到设备: ${targetDevice.name} (${targetDevice.airport_sn || targetDevice.sn})`)
      
      // 如果之前已连接MQTT，自动重新连接
      if (wasConnected) {
        try {
          console.log('设备切换后自动重新连接MQTT...')
          await mqttProxyStore.connect()
          ElMessage.success('MQTT已自动重新连接，请手动订阅新设备的主题')
        } catch (mqttError) {
          console.error('自动重连MQTT失败:', mqttError)
          ElMessage.warning('设备切换成功，但MQTT重连失败，请手动连接')
        }
      }
    } catch (error) {
      ElMessage.error(`切换设备失败: ${error.message}`)
    }
  } else {
    // 清空当前设备
    deviceStore.currentDeviceSn = ''
    ElMessage.info('已清空当前设备选择')
  }
}

// 连接MQTT
const handleConnect = async () => {
  try {
    await mqttProxyStore.connect()
    ElMessage.success('MQTT连接成功')
  } catch (error) {
    ElMessage.error(`连接失败: ${error.message}`)
  }
}

// 清空所有订阅
const clearAllSubscriptions = async () => {
  try {
    const currentSubscriptions = Array.from(subscribedTopics.value)
    if (currentSubscriptions.length === 0) {
      ElMessage.info('当前没有订阅任何主题')
      return
    }
    
    // 取消所有订阅
    for (const topic of currentSubscriptions) {
      console.log('取消订阅:', topic)
      await mqttProxyStore.unsubscribeTopics(topic)
    }
    
    ElMessage.success(`已清空 ${currentSubscriptions.length} 个订阅`)
  } catch (error) {
    console.error('清空订阅失败:', error)
    ElMessage.error(`清空订阅失败: ${error.message}`)
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
              // 如果已订阅，直接打开对话框
              openTopicDialog(topic)
              topicInput.value = ''
              return
            }
            
            await mqttProxyStore.subscribeToTopics(topic)
            recordingTopics.value.add(topic)
            topicInput.value = '' // 清空输入框
            ElMessage.success(`成功订阅主题: ${topic}`)
            
            // 订阅成功后自动打开对话框
            setTimeout(() => {
              openTopicDialog(topic)
            }, 500) // 延迟500ms确保订阅完成
          } catch (error) {
            ElMessage.error(`订阅失败: ${error.message}`)
          }
        }

               // 断开连接
        const handleDisconnect = () => {
          mqttProxyStore.disconnect()
          ElMessage.info('MQTT连接已断开')
        }
        
                 // 删除主题
         const removeTopic = async (index) => {
           const topic = subscribedTopics.value[index]
           try {
             await mqttProxyStore.unsubscribeTopics(topic)
             ElMessage.info(`已移除主题: ${topic}`)
           } catch (e) {
             console.error('取消订阅失败:', e)
           }
         }
         
                   // 订阅指定主题
          const subscribeTopic = async (topicTemplate) => {
            try {
              // 将模板中的 {device_sn} 和 {gateway_sn} 替换为实际的设备ID
              const topic = deviceStore.getTopicWithDeviceSn(topicTemplate)
              
              console.log('点击订阅主题:', topic)
              console.log('当前已订阅主题:', subscribedTopics.value)
              
              if (subscribedTopics.value.includes(topic)) {
                ElMessage.warning('该主题已订阅')
                return
              }
              
              console.log('开始订阅主题:', topic)
              await mqttProxyStore.subscribeToTopics(topic)
              console.log('订阅请求已发送:', topic)
              
              // 等待一下让状态更新
              await new Promise(resolve => setTimeout(resolve, 100))
              
              // 检查订阅是否成功
              if (subscribedTopics.value.includes(topic)) {
                ElMessage.success(`成功订阅主题: ${topic}`)
              } else {
                ElMessage.warning(`订阅请求已发送，请稍后查看状态: ${topic}`)
              }
            } catch (error) {
              console.error('订阅失败:', error)
              ElMessage.error(`订阅失败: ${error.message}`)
            }
          }

          // 取消订阅指定主题
          const unsubscribeTopic = async (topic) => {
            try {
              console.log('点击取消订阅主题:', topic)
              
              if (!subscribedTopics.value.includes(topic)) {
                ElMessage.warning('该主题未订阅')
                return
              }
              
              console.log('开始取消订阅主题:', topic)
              await mqttProxyStore.unsubscribeTopics(topic)
              console.log('取消订阅请求已发送:', topic)
              
              // 等待一下让状态更新
              await new Promise(resolve => setTimeout(resolve, 100))
              
              // 检查取消订阅是否成功
              if (!subscribedTopics.value.includes(topic)) {
                ElMessage.success(`成功取消订阅主题: ${topic}`)
              } else {
                ElMessage.warning(`取消订阅请求已发送，请稍后查看状态: ${topic}`)
              }
            } catch (error) {
              console.error('取消订阅失败:', error)
              ElMessage.error(`取消订阅失败: ${error.message}`)
            }
          }

// 清空历史记录
const clearHistory = () => {
  mqttProxyStore.clearMessageHistory()
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

 
// 获取指定主题的消息
const getTopicMessages = (topic) => {
  const messages = mqttProxyStore.messageHistory.filter(message => message.topic === topic)
  // console.log(`主题 ${topic} 的消息数量:`, messages.length)
  // console.log('所有消息历史:', mqttProxyStore.messageHistory.length)
  return messages
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
 
  // ===== 服务端发送（方法+参数） =====
  const serviceMethods = ref({
    debug: [
      { value: 'debug_mode_open', label: '开启调试', params: '{}' },
      { value: 'debug_mode_close', label: '关闭调试', params: '{}' },
      { value: 'cover_open', label: '打开舱盖', params: '{}' },
      { value: 'cover_close', label: '关闭舱盖', params: '{}' },
      { value: 'supplement_light_open', label: '打开补光灯', params: '{}' },
      { value: 'supplement_light_close', label: '关闭补光灯', params: '{}' },
      { value: 'device_reboot', label: '机场重启', params: '{}' }
    ],
    business: [
      { value: 'live_set_params', label: '设置直播参数', params: '{"resolution":"1920x1080","bitrate":2000,"fps":25,"quality":"high","audio_enabled":false}' },
      { value: 'live_start', label: '开始直播', params: '{"data":{"url":"channel=1ZNDH1D0010098_39-0-7&sn=1ZNDH1D0010098&token=006dca67721582a48768ec4d817b7b25a86IADk%2Fcm%2Fdv%2BHY6qT%2FAKM6y7TcUe4lXNvZpycH7vUMAlM6pFALUKF2zyCIgA82pQE8cCoYAQAAQDxwKhgAgDxwKhgAwDxwKhgBADxwKhg&uid=50000","url_type":0,"video_id":"1ZNDH1D0010098/39-0-7/normal-0","video_quality":0}}' },
      { value: 'live_stop', label: '停止直播', params: '{"video_id":"1ZNDH1D0010098/42-0-0/zoom-0"}' },
      { value: 'live_lens_change', label: '设置直播摄像头', params: '{"video_id":"1581F5BMD228Q00A82XX/39-0-7/zoom-0","video_type":"zoom"}' }
    ]
  })
  
  const serviceMethod = ref('debug_mode_open')
  const serviceParams = ref('{}')
  const sendingService = ref(false)
  
  // 方法选择面板相关
  const methodPanelVisible = ref(false)
  const methodSearchQuery = ref('')
  const selectedMethodTemp = ref(null)
  const serviceParamsTemp = ref('')
  const paramsError = ref('')
  
  // 打开方法选择面板
  const openMethodPanel = () => {
    methodPanelVisible.value = true
    selectedMethodTemp.value = serviceMethod.value
    serviceParamsTemp.value = serviceParams.value
    methodSearchQuery.value = ''
    paramsError.value = ''
    // 如果有选中的方法，加载默认参数
    if (selectedMethodTemp.value) {
      loadParamsForMethod(selectedMethodTemp.value)
    }
  }
  
  // 确认方法选择
  const confirmMethodSelection = () => {
    // 验证参数格式
    if (!validateParams()) {
      ElMessage.warning('请修正参数格式错误')
      return
    }
    if (selectedMethodTemp.value) {
      serviceMethod.value = selectedMethodTemp.value
      // 如果参数为空或只包含空白，设置为空对象
      if (!serviceParamsTemp.value || !serviceParamsTemp.value.trim()) {
        serviceParams.value = '{}'
      } else {
        serviceParams.value = serviceParamsTemp.value
      }
    }
    methodPanelVisible.value = false
    methodSearchQuery.value = ''
    paramsError.value = ''
    ElMessage.success('方法和参数已更新')
  }
  
  // 关闭方法选择面板（恢复原值）
  const closeMethodPanel = () => {
    selectedMethodTemp.value = serviceMethod.value
    serviceParamsTemp.value = serviceParams.value
    methodPanelVisible.value = false
    methodSearchQuery.value = ''
    paramsError.value = ''
  }
  
  // 选择方法
  const selectMethod = (method) => {
    selectedMethodTemp.value = method.value
    loadParamsForMethod(method.value)
  }
  
  // 选择自定义方法
  const selectCustomMethod = (methodName) => {
    selectedMethodTemp.value = methodName
    // 自定义方法没有默认参数，保持当前参数或设为空对象
    if (!serviceParamsTemp.value || serviceParamsTemp.value.trim() === '') {
      serviceParamsTemp.value = '{}'
    }
    paramsError.value = ''
  }
  
  // 为方法加载参数
  const loadParamsForMethod = (methodValue) => {
    const allMethods = [...serviceMethods.value.debug, ...serviceMethods.value.business]
    const method = allMethods.find(m => m.value === methodValue)
    if (method && method.params) {
      serviceParamsTemp.value = method.params
      paramsError.value = ''
    } else {
      // 如果没有默认参数，设为空对象
      if (!serviceParamsTemp.value || serviceParamsTemp.value.trim() === '') {
        serviceParamsTemp.value = '{}'
      }
    }
    validateParams()
  }
  
  // 获取选中方法的标签
  const getSelectedMethodLabel = () => {
    if (!selectedMethodTemp.value) return ''
    const allMethods = [...serviceMethods.value.debug, ...serviceMethods.value.business]
    const method = allMethods.find(m => m.value === selectedMethodTemp.value)
    return method ? method.label : ''
  }
  
  // 检查是否有默认参数
  const hasDefaultParams = () => {
    if (!selectedMethodTemp.value) return false
    const allMethods = [...serviceMethods.value.debug, ...serviceMethods.value.business]
    const method = allMethods.find(m => m.value === selectedMethodTemp.value)
    return method && method.params && method.params !== '{}'
  }
  
  // 加载默认参数
  const loadDefaultParams = () => {
    if (selectedMethodTemp.value) {
      loadParamsForMethod(selectedMethodTemp.value)
    }
  }
  
  // 格式化参数JSON
  const formatParams = () => {
    try {
      if (serviceParamsTemp.value && serviceParamsTemp.value.trim()) {
        const parsed = JSON.parse(serviceParamsTemp.value)
        serviceParamsTemp.value = JSON.stringify(parsed, null, 2)
        paramsError.value = ''
        ElMessage.success('JSON格式化成功')
      }
    } catch (e) {
      paramsError.value = 'JSON格式错误: ' + e.message
      ElMessage.error('JSON格式错误: ' + e.message)
    }
  }
  
  // 清空参数
  const clearParams = () => {
    serviceParamsTemp.value = '{}'
    paramsError.value = ''
  }
  
  // 验证参数格式
  const validateParams = () => {
    if (!serviceParamsTemp.value || serviceParamsTemp.value.trim() === '') {
      paramsError.value = ''
      return true
    }
    try {
      JSON.parse(serviceParamsTemp.value)
      paramsError.value = ''
      return true
    } catch (e) {
      paramsError.value = 'JSON格式错误: ' + e.message
      return false
    }
  }
  
  // 检查方法是否在列表中
  const isMethodInList = (methodName) => {
    const allMethods = [...serviceMethods.value.debug, ...serviceMethods.value.business]
    return allMethods.some(m => m.value === methodName)
  }
  
  // 过滤调试方法
  const filteredDebugMethods = computed(() => {
    if (!methodSearchQuery.value) {
      return serviceMethods.value.debug
    }
    const query = methodSearchQuery.value.toLowerCase()
    return serviceMethods.value.debug.filter(m => 
      m.value.toLowerCase().includes(query) || 
      m.label.toLowerCase().includes(query)
    )
  })
  
  // 过滤业务方法
  const filteredBusinessMethods = computed(() => {
    if (!methodSearchQuery.value) {
      return serviceMethods.value.business
    }
    const query = methodSearchQuery.value.toLowerCase()
    return serviceMethods.value.business.filter(m => 
      m.value.toLowerCase().includes(query) || 
      m.label.toLowerCase().includes(query)
    )
  })
  
  // 方法搜索处理
  const handleMethodSearch = () => {
    // 搜索时自动选择匹配的第一个方法（如果有）
    const allFiltered = [...filteredDebugMethods.value, ...filteredBusinessMethods.value]
    if (allFiltered.length > 0 && !isMethodInList(methodSearchQuery.value)) {
      // 如果有搜索匹配的方法，不自动选择自定义方法
    }
  }
  
  // 方法切换处理
  const handleMethodChange = (methodValue) => {
    // 从所有方法列表中查找
    const allMethods = [...serviceMethods.value.debug, ...serviceMethods.value.business]
    const method = allMethods.find(m => m.value === methodValue)
    if (method && method.params) {
      serviceParams.value = method.params
    } else if (!method) {
      // 自定义方法名，不更新参数，保持当前参数不变
      // 用户可以手动编辑参数
    }
  }
 
  const handleSendService = async () => {
    try {
      if (!serviceMethod.value) {
        ElMessage.warning('请选择方法')
        return
      }
      let data = {}
      if (serviceParams.value && serviceParams.value.trim()) {
        try {
          data = JSON.parse(serviceParams.value)
        } catch (e) {
          ElMessage.error('参数必须是合法的JSON')
          return
        }
      }
      const gatewaySn = deviceStore.currentDevice?.airport_sn || deviceStore.currentDeviceSn
      if (!gatewaySn) {
        ElMessage.error('缺少网关SN（机场SN）')
        return
      }
      sendingService.value = true
      const { tid } = await mqttProxyStore.sendService(gatewaySn, serviceMethod.value, data, { autoSubscribeReply: true })
      ElMessage.success(`已发送服务调用，tid=${tid}`)
    } catch (err) {
      ElMessage.error(err?.message || '发送失败')
    } finally {
      sendingService.value = false
    }
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
   mqttProxyStore.messageHistory = mqttProxyStore.messageHistory.filter(
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

 // 展示设备卡片的条件
const shouldShowDeviceCard = (device) => {
  if (!device || !device.topic) return false
  const t = device.topic
  
  // 如果主题在已订阅列表中，则显示
  if (subscribedTopics.value.includes(t)) {
    return true
  }
  
  // 显示所有有消息的主题（无论是否订阅）
  return true
}

// 获取可见设备卡片数量
const getVisibleDeviceCount = () => {
  return allDisplayTopics.value.size
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

// 打开飞机详情页
const openAircraftDetail = (deviceSn) => {
  if (!deviceSn) {
    ElMessage.warning('设备SN为空，无法打开详情页')
    return
  }
  
  // 获取当前设备信息
  const currentDevice = deviceStore.currentDevice
  const airportSn = currentDevice?.airport_sn
  const aircraftSn = deviceSn
  
  if (airportSn && airportSn !== aircraftSn) {
    // 使用新的双参数路由：/aircraft/:airportSn/:aircraftSn
    console.log('使用新路由格式:', { airportSn, aircraftSn })
    router.push({ 
      name: 'AircraftDetail', 
      params: { 
        airportSn: airportSn,
        aircraftSn: aircraftSn 
      } 
    })
  } else {
    // 使用旧的单参数路由：/aircraft/:deviceSn
    console.log('使用旧路由格式:', { deviceSn })
    router.push({ 
      name: 'AircraftDetailLegacy', 
      params: { deviceSn } 
    })
  }
}

// 打开摄像头直播页面
const openCameraLive = () => {
  router.push({ name: 'CameraLive' })
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

// 消息详情抽屉相关
const messageDetailVisible = ref(false)
const selectedMessage = ref(null)
const copySuccess = ref(false)

const openErrorDoc = () => {
  errorDocTab.value = 'dji'
  errorDocVisible.value = true
}

// 格式化JSON数据
const formatJsonPayload = (payload) => {
  try {
    let jsonString
    if (typeof payload === 'string') {
      // 尝试解析字符串为JSON
      const parsed = JSON.parse(payload)
      jsonString = JSON.stringify(parsed, null, 2)
    } else if (typeof payload === 'object' && payload !== null) {
      // 直接格式化对象
      jsonString = JSON.stringify(payload, null, 2)
    } else {
      // 其他类型直接返回
      return String(payload)
    }
    
    // 添加语法高亮
    return jsonString
      .replace(/(".*?")\s*:/g, '<span class="json-key">$1</span>:')
      .replace(/:\s*(".*?")/g, ': <span class="json-string">$1</span>')
      .replace(/:\s*(\d+)/g, ': <span class="json-number">$1</span>')
      .replace(/:\s*(true|false)/g, ': <span class="json-boolean">$1</span>')
      .replace(/:\s*(null)/g, ': <span class="json-null">$1</span>')
      .replace(/([{}[\]])/g, '<span class="json-bracket">$1</span>')
      .replace(/(,)/g, '<span class="json-comma">$1</span>')
  } catch (error) {
    // 如果解析失败，返回原始字符串
    return String(payload)
  }
}

// 复制JSON数据
const copyJsonPayload = async () => {
  if (!selectedMessage.value) return
  
  try {
    let jsonString
    if (typeof selectedMessage.value.payload === 'string') {
      // 尝试解析字符串为JSON
      const parsed = JSON.parse(selectedMessage.value.payload)
      jsonString = JSON.stringify(parsed, null, 2)
    } else if (typeof selectedMessage.value.payload === 'object' && selectedMessage.value.payload !== null) {
      // 直接格式化对象
      jsonString = JSON.stringify(selectedMessage.value.payload, null, 2)
    } else {
      // 其他类型直接返回
      jsonString = String(selectedMessage.value.payload)
    }
    
    await navigator.clipboard.writeText(jsonString)
    copySuccess.value = true
    
    // 2秒后重置状态
    setTimeout(() => {
      copySuccess.value = false
    }, 2000)
    
    ElMessage.success('JSON数据已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败，请手动复制')
  }
}

// 显示消息详情
const showMessageDetail = (message) => {
  selectedMessage.value = message
  messageDetailVisible.value = true
}

// 错误码查询方法
const loadErrorCodes = async () => {
  // 加载 DJI 错误码
  await loadDjiErrorCodes()
  // 加载 HMS 错误码
  await loadHmsErrorCodes()
}

// 加载 DJI 错误码
const loadDjiErrorCodes = async () => {
  try {
    const response = await fetch('/docs/dji-error-codes.md')
    const text = await response.text()
    // 解析制表符分隔的错误码格式
    const lines = text.split('\n')
    const codes = []
    
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i].trim()
      if (line && !line.includes('错误码') && !line.includes('错误码描述')) {
        // 使用制表符或空格分割
        const parts = line.split(/\t|\s{2,}/).filter(p => p.trim())
        if (parts.length >= 2) {
          const code = parts[0].trim()
          const description = parts[1].trim()
          
          // 根据错误码范围判断级别
          let level = 'INFO'
          const codeNum = parseInt(code)
          if (codeNum >= 312000 && codeNum < 313000) {
            level = 'WARN'
          } else if (codeNum >= 314000) {
            level = 'ERROR'
          }
          
          codes.push({
            code: code,
            level: level,
            description: description
          })
        }
      }
    }
    
    djiCodes.value = codes
    filteredDjiCodes.value = codes
  } catch (error) {
    console.error('加载DJI错误码失败:', error)
    // 使用示例数据
    djiCodes.value = [
      { code: '312014', level: 'WARN', description: '设备升级中，请勿重复操作' },
      { code: '312015', level: 'WARN', description: '机场业务繁忙无法进行设备升级' },
      { code: '314000', level: 'ERROR', description: '设备当前无法支持该操作' }
    ]
    filteredDjiCodes.value = djiCodes.value
  }
}

// 加载 HMS 错误码
const loadHmsErrorCodes = async () => {
  try {
    const response = await fetch('/docs/hms.json')
    const data = await response.json()
    const codes = []
    
    for (const [key, value] of Object.entries(data)) {
      // 从 key 中提取错误码
      const codeMatch = key.match(/0x[0-9a-fA-F]+/)
      if (codeMatch && value) {
        codes.push({
          code: codeMatch[0],
          key: key,
          zh: value.zh || '',
          en: value.en || ''
        })
      }
    }
    
    hmsCodes.value = codes
    filteredHmsCodes.value = codes
  } catch (error) {
    console.error('加载HMS错误码失败:', error)
    // 使用示例数据
    hmsCodes.value = [
      { code: '0x12040000', key: 'dock_tip_0x12040000', zh: '机场RTK设备断连', en: 'Dock RTK device disconnected' },
      { code: '0x12040001', key: 'dock_tip_0x12040001', zh: '机场位置被移动', en: 'Dock moved' }
    ]
    filteredHmsCodes.value = hmsCodes.value
  }
}

// DJI 错误码搜索
const handleDjiSearch = () => {
  const query = djiSearchQuery.value.toLowerCase()
  console.log('DJI搜索查询:', query)
  
  if (!query) {
    filteredDjiCodes.value = djiCodes.value
    return
  }
  
  const filtered = djiCodes.value.filter(code => {
    const codeMatch = code.code.toLowerCase().includes(query)
    const descMatch = code.description.toLowerCase().includes(query)
    const levelMatch = code.level.toLowerCase().includes(query)
    
    return codeMatch || descMatch || levelMatch
  })
  
  console.log('DJI搜索结果数量:', filtered.length)
  filteredDjiCodes.value = filtered
}

// HMS 错误码搜索
const handleHmsSearch = () => {
  const query = hmsSearchQuery.value.toLowerCase()
  console.log('HMS搜索查询:', query)
  console.log('HMS总错误码数量:', hmsCodes.value.length)
  
  if (!query) {
    filteredHmsCodes.value = hmsCodes.value
    console.log('清空HMS搜索，显示所有错误码')
    return
  }
  
  const filtered = hmsCodes.value.filter(code => {
    const codeMatch = code.code && code.code.toLowerCase().includes(query)
    const zhMatch = code.zh && code.zh.toLowerCase().includes(query)
    const enMatch = code.en && code.en.toLowerCase().includes(query)
    
    if (codeMatch || zhMatch || enMatch) {
      console.log('HMS匹配到错误码:', code)
    }
    
    return codeMatch || zhMatch || enMatch
  })
  
  console.log('HMS搜索结果数量:', filtered.length)
  console.log('HMS搜索结果:', filtered.slice(0, 3))
  filteredHmsCodes.value = filtered
}

// 数据源切换处理
const handleSourceChange = () => {
  console.log('切换到数据源:', errorCodeSource.value)
  // 清空搜索框
  djiSearchQuery.value = ''
  hmsSearchQuery.value = ''
  // 重置搜索结果
  filteredDjiCodes.value = djiCodes.value
  filteredHmsCodes.value = hmsCodes.value
}

const getErrorCodeType = (level) => {
  switch (level) {
    case 'ERROR': return 'danger'
    case 'WARN': return 'warning'
    case 'INFO': return 'success'
    default: return 'info'
  }
}

const showErrorCodeDetail = (errorCode) => {
  ElMessage({
    message: `错误码: ${errorCode.code}\n级别: ${errorCode.level}\n描述: ${errorCode.description}`,
    type: 'info',
    duration: 5000,
    showClose: true
  })
}

const showHmsErrorDetail = (errorCode) => {
  const zhText = errorCode.zh || '无中文描述'
  const enText = errorCode.en || '无英文描述'
  ElMessage({
    message: `错误码: ${errorCode.code}\n中文: ${zhText}\n英文: ${enText}`,
    type: 'info',
    duration: 5000,
    showClose: true
  })
}

// 错误码详情页面相关方法
const openErrorCodePage = () => {
  errorCodePageVisible.value = true
  pageSearchQuery.value = ''
  pageDataSource.value = 'all'
  currentPage.value = 1
  updatePageFilteredCodes()
}

const handleCloseErrorCodePage = () => {
  errorCodePageVisible.value = false
}

const updatePageFilteredCodes = () => {
  let allCodes = []
  
  // 合并所有错误码
  if (pageDataSource.value === 'all' || pageDataSource.value === 'dji') {
    const djiCodesWithType = djiCodes.value.map(code => ({
      ...code,
      type: 'dji',
      id: `dji_${code.code}`
    }))
    allCodes = allCodes.concat(djiCodesWithType)
  }
  
  if (pageDataSource.value === 'all' || pageDataSource.value === 'hms') {
    const hmsCodesWithType = hmsCodes.value.map(code => ({
      ...code,
      type: 'hms',
      id: `hms_${code.code}`
    }))
    allCodes = allCodes.concat(hmsCodesWithType)
  }
  
  // 搜索过滤
  if (pageSearchQuery.value) {
    const query = pageSearchQuery.value.toLowerCase()
    allCodes = allCodes.filter(code => {
      if (code.type === 'dji') {
        return (code.code && code.code.toLowerCase().includes(query)) ||
               (code.description && code.description.toLowerCase().includes(query)) ||
               (code.level && code.level.toLowerCase().includes(query))
      } else {
        return (code.code && code.code.toLowerCase().includes(query)) ||
               (code.zh && code.zh.toLowerCase().includes(query)) ||
               (code.en && code.en.toLowerCase().includes(query))
      }
    })
  }
  
  pageFilteredCodes.value = allCodes
  console.log('页面错误码总数:', allCodes.length)
}

const handlePageSearch = () => {
  currentPage.value = 1
  updatePageFilteredCodes()
}

const handlePageDataSourceChange = () => {
  currentPage.value = 1
  updatePageFilteredCodes()
}

const handlePageChange = (page) => {
  currentPage.value = page
}

const showErrorCodePageDetail = (errorCode) => {
  if (errorCode.type === 'dji') {
    showErrorCodeDetail(errorCode)
  } else {
    showHmsErrorDetail(errorCode)
  }
}
const closeErrorDoc = () => {
  errorDocVisible.value = false
}

const openConnMgr = () => { connMgrRef.value && connMgrRef.value.open() }

// 查看全部错误码
const showAllErrorCodes = () => {
  errorCodePageVisible.value = true
  pageSearchQuery.value = ''
  pageDataSource.value = 'all'
  currentPage.value = 1
  updatePageFilteredCodes()
}

// 打开错误码查询抽屉
const openErrorCodeDialog = () => {
  errorCodePageVisible.value = true
}

// 组件挂载时自动连接
onMounted(() => {
  // 加载保存的主题顺序
  loadTopicCategoryOrder()
  // 加载错误码数据
  loadErrorCodes()
})

// 组件卸载时断开连接
onUnmounted(() => {
  // 可以选择是否在组件卸载时断开连接
  // mqttProxyStore.disconnect()
})
 </script>

<style scoped>
.dashboard {
  height: 100vh;
  overflow: hidden;
}

/* 设备选择器样式 */
.current-device-selector {
  display: flex;
  align-items: center;
  gap: 15px;
}

.current-device-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--el-fill-color-light);
  border-radius: 6px;
  border: 1px solid var(--el-border-color);
}

.device-name {
  font-weight: 500;
  font-size: 12px;
  color: var(--el-text-color-primary);
}

.device-sn {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: var(--el-text-color-regular);
  background: var(--el-fill-color);
  padding: 2px 6px;
  border-radius: 3px;
}

.device-option {
  display: flex;
  align-items: center;
  width: 100%;
}

/* 卡片头部操作按钮样式 */
.card-header-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 侧边栏收缩状态指示器 */
.sidebar-collapsed-indicator {
  position: absolute;
  top: 20px;
  left: 10px;
  z-index: 1000;
}

.expand-button {
  font-size: 14px;
  font-weight: bold;
  padding: 8px 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  border-radius: 6px;
  transition: all 0.3s ease;
}

.expand-button:hover {
  transform: translateX(2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

/* MQTT连接信息样式 */
.mqtt-connection-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--el-fill-color-lighter);
  border-radius: 6px;
  border: 1px solid var(--el-border-color);
  margin-left: 10px;
}

.mqtt-profile-name {
  font-weight: 500;
  font-size: 12px;
  color: var(--el-text-color-primary);
}

.mqtt-server-info {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: var(--el-text-color-regular);
  background: var(--el-fill-color);
  padding: 2px 6px;
  border-radius: 3px;
}

.dashboard-aside {
  border-right: 1px solid var(--el-border-color);
  background: var(--el-fill-color-blank);
  transition: width 0.3s ease;
  overflow: hidden;
}

.dashboard-aside.collapsed {
  width: 0 !important;
  border-right: none;
}

.sidebar-content {
  width: 100%;
  height: 100%;
  overflow-y: auto;
  padding: 8px;
}


.topic-list-card {
  margin-bottom: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.topic-list-card .el-card__body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dashboard-header {
  border-bottom: 1px solid var(--el-border-color);
  background: var(--el-fill-color-blank);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.collapse-btn {
  flex-shrink: 0;
  font-weight: bold;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  margin-right: 12px;
}

.collapse-btn:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.dashboard-main {
  padding: 16px;
  background: var(--el-fill-color-lighter);
  overflow-y: auto;
}

.device-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.device-card {
  min-height: 200px;
}

.device-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.device-id {
  font-size: 12px;
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.device-content {
  padding: 8px 0;
}

.device-info p {
  margin: 4px 0;
  font-size: 14px;
}

.topic-categories {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.topic-category {
  margin-bottom: 12px;
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  padding: 8px;
  background: var(--el-fill-color-blank);
  transition: all 0.3s;
}

.topic-category:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.category-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.drag-handle {
  cursor: move;
  color: var(--el-text-color-secondary);
}

.topic-item {
  margin-bottom: 8px;
  padding: 8px;
  background: var(--el-fill-color-light);
  border-radius: 4px;
  border: 1px solid transparent;
  transition: all 0.2s;
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.topic-item:hover {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

.topic-name {
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 4px;
}

.topic-path {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
}

.topic-desc {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.4;
}

.error-message {
  margin-bottom: 16px;
}

.topic-messages-container {
  margin-top: 16px;
}

/* 错误码查询抽屉样式 */
.error-code-drawer .el-drawer__body {
  padding: 0;
}

.error-code-page {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.page-header {
  padding: 20px;
  border-bottom: 1px solid var(--el-border-color-light);
  background: var(--el-bg-color);
}

.search-controls {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.page-stats {
  display: flex;
  align-items: center;
}

.error-codes-page-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.error-code-page-item {
  padding: 16px;
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.error-code-page-item:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.error-code-page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.error-code-page-code {
  font-size: 16px;
  font-weight: bold;
  color: var(--el-color-primary);
}

.error-code-page-tags {
  display: flex;
  gap: 4px;
}

.error-code-page-description {
  color: var(--el-text-color-regular);
  line-height: 1.5;
}

.error-code-page-zh {
  font-weight: 500;
}

.error-code-page-en {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}

.pagination-container {
  padding: 20px;
  border-top: 1px solid var(--el-border-color-light);
  background: var(--el-bg-color);
}

/* 消息详情抽屉样式 */
.message-detail-drawer .el-drawer__body {
  padding: 0;
}

.message-detail {
  padding: 20px;
  height: 100%;
  overflow-y: auto;
}

.detail-section {
  margin-bottom: 24px;
  padding: 16px;
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  background: var(--el-bg-color);
}

.detail-section h4 {
  margin: 0 0 16px 0;
  color: var(--el-color-primary);
  font-size: 16px;
  font-weight: 600;
}

.json-section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.json-section-header h4 {
  margin: 0;
}

.copy-success {
  background-color: var(--el-color-success) !important;
  border-color: var(--el-color-success) !important;
  color: white !important;
}

.detail-row {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.detail-row:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.detail-label {
  font-weight: 600;
  color: var(--el-text-color-regular);
  min-width: 80px;
  margin-right: 12px;
}

.detail-value {
  color: var(--el-text-color-primary);
  word-break: break-all;
}

.detail-payload {
  background: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color-light);
  border-radius: 6px;
  padding: 12px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  max-height: 300px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-word;
}

/* 主题列表样式 */
.topic-list {
  max-height: 1000px;
  overflow-y: auto;
}

.topic-count {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-left: 8px;
}

.empty-topics {
  padding: 20px;
  text-align: center;
}

.topic-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.topic-item:last-child {
  border-bottom: none;
}

.topic-info {
  flex: 1;
  min-width: 0;
}

.topic-path {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: var(--el-text-color-primary);
  word-break: break-all;
  line-height: 1.4;
  font-weight: bold;
}

/* JSON查看器样式 */
.json-viewer {
  background: #1e1e1e;
  border: 1px solid #3c3c3c;
  border-radius: 8px;
  padding: 16px;
  max-height: 400px;
  overflow-y: auto;
  position: relative;
}

.json-content {
  margin: 0;
  font-family: 'Fira Code', 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #d4d4d4;
  white-space: pre-wrap;
  word-break: break-word;
  background: none;
  border: none;
  padding: 0;
}

/* JSON语法高亮 */
.json-key {
  color: #9cdcfe;
  font-weight: 600;
}

.json-string {
  color: #ce9178;
}

.json-number {
  color: #b5cea8;
}

.json-boolean {
  color: #569cd6;
  font-weight: 600;
}

.json-null {
  color: #569cd6;
  font-style: italic;
}

.json-bracket {
  color: #d4d4d4;
  font-weight: bold;
}

.json-comma {
  color: #d4d4d4;
}
</style>
