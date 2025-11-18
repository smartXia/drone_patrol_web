<template>
  <div class="airport-info-tab">
    <!-- 订阅状态提示 -->
    <div v-if="!hasReceivedData" class="waiting-container">
      <el-alert
        title="等待接收机场OSD数据"
        type="info"
        :closable="false"
        show-icon
      >
        <template #default>
          <div style="margin-top: 10px;">
            <p><strong>订阅状态:</strong> {{ subscriptionStatus }}</p>
            <p><strong>订阅主题:</strong> <code>{{ airportOsdTopic || 'thing/product/' + getAirportSn() + '/osd' }}</code></p>
            <p><strong>等待数据中...</strong> 请确保设备已连接并正在发送OSD数据</p>
            <div style="margin-top: 15px;">
              <el-button 
                type="primary" 
                size="small" 
                :loading="subscribing"
                @click="subscribeToAirportOSD"
              >
                {{ subscribing ? '订阅中...' : '手动订阅' }}
              </el-button>
              <el-button 
                size="small" 
                @click="refreshStatus"
                style="margin-left: 10px;"
              >
                刷新状态
              </el-button>
              <el-button 
                size="small" 
                @click="checkStateMessages"
                type="info"
                style="margin-left: 10px;"
              >
                检查状态消息
              </el-button>
            </div>
            
            <!-- 调试信息 -->
            <div style="margin-top: 15px; padding: 10px; background: var(--el-color-info-light-9); border-radius: 4px; font-size: 12px;">
              <div><strong>调试信息:</strong></div>
              <div>MQTT连接: {{ mqttProxyStore.isConnected ? '✅ 已连接' : '❌ 未连接' }}</div>
              <div>消息历史数量: {{ mqttProxyStore.messageHistory.length }}</div>
              <div>机场OSD主题: <code>{{ airportOsdTopic }}</code></div>
              <div>设备状态主题: <code>thing/product/{{ getAirportSn() }}/state</code></div>
              <div>机场序列号: <code>{{ getAirportSn() }}</code></div>
              <div>已接收数据: {{ hasReceivedData ? '✅ 是' : '❌ 否' }}</div>
              <div>状态数据连接: {{ stateDataStatus === 'connected' ? '✅ 已连接' : '❌ 未连接' }}</div>
              <div>负载属性: {{ gimbalPitch !== null ? '✅ 有数据' : '❌ 无数据' }}</div>
              <div v-if="mqttProxyStore.messageHistory.length > 0">
                <div>最新消息主题: <code>{{ mqttProxyStore.messageHistory[0]?.topic }}</code></div>
                <div>最新消息时间: {{ new Date(mqttProxyStore.messageHistory[0]?.timestamp).toLocaleString() }}</div>
              </div>
            </div>
          </div>
        </template>
      </el-alert>
    </div>

    <!-- 机场信息展示 -->
    <div class="data-container">
      <!-- 机场信息 -->
      <el-row :gutter="8" class="data-grid">
        <el-col :span="12">
          <h3 style="margin-bottom: 15px; color: var(--el-text-color-primary); border-left: 4px solid var(--el-color-primary); padding-left: 10px;">
            🏢 机场信息
            <span v-if="osdData?.timestamp" style="font-size: 12px; color: var(--el-color-success); margin-left: 10px;">
              (数据已更新: {{ formatUpdateTime(osdData.timestamp) }})
            </span>
          </h3>
        </el-col>
        <el-col :span="12" style="text-align: right;">
          <div style="margin-bottom: 15px;">
            <span style="color: var(--el-text-color-regular); font-size: 14px;">机场序列号: </span>
            <span style="color: var(--el-text-color-primary); font-weight: 500;">{{ getAirportSn() }}</span>
            <!-- 调试信息 -->
            <div style="font-size: 10px; color: #999; margin-top: 2px;">
              调试: gateway={{ osdData?.gateway }}, deviceSn={{ deviceSn }}, currentDevice.airport_sn={{ deviceStore.currentDevice?.airport_sn }}
            </div>
          </div>
        </el-col>
        
        <!-- 环境状态信息 -->
        <el-col :span="4">
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Sunny /></el-icon>
                <span>环境状态</span>
                <el-tag size="small" :type="environmentData ? (isDataStale(environmentData.timestamp) ? 'warning' : 'success') : 'info'">
                  {{ environmentData ? formatRelativeTime(environmentData.timestamp) : '等待' }}
                  <span v-if="environmentData && isDataStale(environmentData.timestamp)" style="margin-left: 4px;">⚠️</span>
                </el-tag>
              </div>
            </template>
            <div class="info-content">
              <div class="info-item">
                <span class="label">舱内温度:</span>
                <span class="value">{{ environmentData?.data?.temperature || osdData?.data?.temperature || '--' }}°C</span>
              </div>
              <div class="info-item">
                <span class="label">环境温度:</span>
                <span class="value">{{ environmentData?.data?.environment_temperature || osdData?.data?.environment_temperature || '--' }}°C</span>
              </div>
              <div class="info-item">
                <span class="label">湿度:</span>
                <span class="value">{{ environmentData?.data?.humidity || osdData?.data?.humidity || '--' }}%</span>
              </div>
              <div class="info-item">
                <span class="label">风速:</span>
                <span class="value">{{ environmentData?.data?.wind_speed || osdData?.data?.wind_speed || '--' }}m/s</span>
              </div>
              <div class="info-item">
                <span class="label">降雨量:</span>
                <el-tag :type="getRainfallType(environmentData?.data?.rainfall || osdData?.data?.rainfall)" size="small">
                  {{ getRainfallText(environmentData?.data?.rainfall || osdData?.data?.rainfall) }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">飞行器在舱:</span>
                <el-tag :type="getDroneInDockType(environmentData?.data?.drone_in_dock || osdData?.data?.drone_in_dock)" size="small">
                  {{ getDroneInDockText(environmentData?.data?.drone_in_dock || osdData?.data?.drone_in_dock) }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">网络类型:</span>
                <el-tag :type="(environmentData?.data?.network_state?.type || osdData?.data?.network_state?.type) === 1 ? 'success' : 'info'" size="small">
                  {{ (environmentData?.data?.network_state?.type || osdData?.data?.network_state?.type) === 1 ? '4G' : '以太网' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">网络质量:</span>
                <el-tag :type="getNetworkQualityType(environmentData?.data?.network_state?.quality || osdData?.data?.network_state?.quality)" size="small">
                  {{ getNetworkQualityText(environmentData?.data?.network_state?.quality || osdData?.data?.network_state?.quality) }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">网络速率:</span>
                <span class="value">{{ (environmentData?.data?.network_state?.rate || osdData?.data?.network_state?.rate) || '--' }} KB/s</span>
              </div>
              <div class="info-item">
                <span class="label">补光灯:</span>
                <el-tag :type="(environmentData?.data?.supplement_light_state || osdData?.data?.supplement_light_state) === 1 ? 'success' : 'info'" size="small">
                  {{ (environmentData?.data?.supplement_light_state || osdData?.data?.supplement_light_state) === 1 ? '开启' : '关闭' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">舱盖状态:</span>
                <el-tag :type="getCoverStateType(environmentData?.data?.cover_state || osdData?.data?.cover_state)" size="small">
                  {{ getCoverStateText(environmentData?.data?.cover_state || osdData?.data?.cover_state) }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">纬度:</span>
                <span class="value">{{ environmentData?.data?.latitude || osdData?.data?.latitude || '--' }}</span>
              </div>
              <div class="info-item">
                <span class="label">经度:</span>
                <span class="value">{{ environmentData?.data?.longitude || osdData?.data?.longitude || '--' }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <!-- 飞机信息 -->
        <el-col :span="4">
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Location /></el-icon>
                <span>飞机信息</span>
                <el-tag size="small" :type="(aircraftData || osdData) ? (isDataStale(aircraftData?.timestamp || osdData?.timestamp) ? 'warning' : 'success') : 'info'">
                  {{ (aircraftData || osdData) ? formatRelativeTime(aircraftData?.timestamp || osdData?.timestamp) : '等待' }}
                  <span v-if="(aircraftData || osdData) && isDataStale(aircraftData?.timestamp || osdData?.timestamp)" style="margin-left: 4px;">⚠️</span>
                  <span v-if="(aircraftData || osdData) && isDataStale(aircraftData?.timestamp || osdData?.timestamp)" style="margin-left: 4px; font-size: 10px;">(历史数据)</span>
                </el-tag>
              </div>
            </template>
            <div class="info-content">
              <!-- 历史数据提示 -->
              <div v-if="(aircraftData || osdData) && isDataStale(aircraftData?.timestamp || osdData?.timestamp)" style="background: var(--el-color-warning-light-9); padding: 5px; border-radius: 4px; margin-bottom: 8px; font-size: 12px; color: var(--el-color-warning);">
                📋 显示历史数据，最后更新: {{ formatUpdateTime(aircraftData?.timestamp || osdData?.timestamp) }}
              </div>
              
              <div class="info-item">
                <span class="label">飞机序列号:</span>
                <span class="value">{{ aircraftData?.data?.sub_device?.device_sn || osdData?.data?.sub_device?.device_sn || '--' }}</span>
              </div>
              <div class="info-item">
                <span class="label">飞机型号:</span>
                <span class="value" :title="`格式: {domain-type-subtype}`">{{ aircraftData?.data?.sub_device?.device_model_key || osdData?.data?.sub_device?.device_model_key || '--' }}</span>
              </div>
              <div class="info-item">
                <span class="label">飞机开机状态:</span>
                <el-tag :type="(aircraftData?.data?.sub_device?.device_online_status || osdData?.data?.sub_device?.device_online_status) === 1 ? 'success' : 'info'" size="small">
                  {{ (aircraftData?.data?.sub_device?.device_online_status || osdData?.data?.sub_device?.device_online_status) === 1 ? '开机' : '关机' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">飞机对频状态:</span>
                <el-tag :type="(aircraftData?.data?.sub_device?.device_paired || osdData?.data?.sub_device?.device_paired) === 1 ? 'success' : 'warning'" size="small">
                  {{ (aircraftData?.data?.sub_device?.device_paired || osdData?.data?.sub_device?.device_paired) === 1 ? '已对频' : '未对频' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">固件版本:</span>
                <span class="value">{{ aircraftData?.data?.firmware_version || osdData?.data?.firmware_version || '--' }}</span>
              </div>
              <div class="info-item">
                <span class="label">电池电量:</span>
                <span class="value">{{ aircraftData?.data?.drone_battery_maintenance_info?.batteries?.[0]?.capacity_percent || osdData?.data?.drone_battery_maintenance_info?.batteries?.[0]?.capacity_percent || '--' }}%</span>
              </div>
              <div class="info-item">
                <span class="label">电池数量:</span>
                <span class="value">{{ aircraftData?.data?.drone_battery_maintenance_info?.batteries?.length || osdData?.data?.drone_battery_maintenance_info?.batteries?.length || '--' }}个</span>
              </div>
              <div class="info-item">
                <span class="label">电池电压:</span>
                <span class="value">{{ aircraftData?.data?.drone_battery_maintenance_info?.batteries?.[0]?.voltage || osdData?.data?.drone_battery_maintenance_info?.batteries?.[0]?.voltage || '--' }}mV</span>
              </div>
              <div class="info-item">
                <span class="label">电池温度:</span>
                <span class="value">{{ aircraftData?.data?.drone_battery_maintenance_info?.batteries?.[0]?.temperature || osdData?.data?.drone_battery_maintenance_info?.batteries?.[0]?.temperature || '--' }}°C</span>
              </div>
              <div class="info-item">
                <span class="label">电池维护状态:</span>
                <el-tag :type="(aircraftData?.data?.drone_battery_maintenance_info?.maintenance_state || osdData?.data?.drone_battery_maintenance_info?.maintenance_state) === 0 ? 'success' : 'warning'" size="small">
                  {{ (aircraftData?.data?.drone_battery_maintenance_info?.maintenance_state || osdData?.data?.drone_battery_maintenance_info?.maintenance_state) === 0 ? '正常' : '需要维护' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">自收敛坐标:</span>
                <span class="value">{{ osdData?.data?.self_converge_coordinate ? `${osdData.data.self_converge_coordinate.longitude.toFixed(6)}, ${osdData.data.self_converge_coordinate.latitude.toFixed(6)}` : '--' }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <!-- 机场基础信息 -->
        <el-col :span="4">
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Connection /></el-icon>
                <span>机场基础</span>
                <el-tag size="small" :type="airportBasicData ? (isDataStale(airportBasicData.timestamp) ? 'warning' : 'success') : 'info'">
                  {{ airportBasicData ? formatRelativeTime(airportBasicData.timestamp) : '等待' }}
                  <span v-if="airportBasicData && isDataStale(airportBasicData.timestamp)" style="margin-left: 4px;">⚠️</span>
                </el-tag>
              </div>
            </template>
            <div class="info-content">
              <div class="info-item">
                <span class="label">机场序列号:</span>
                <span class="value">{{ getAirportSn() }}</span>
              </div>
              <div class="info-item">
                <span class="label">累计作业次数:</span>
                <span class="value">{{ airportBasicData?.data?.job_number || osdData?.data?.job_number || '--' }} 次</span>
              </div>
              <div class="info-item">
                <span class="label">累计运行时长:</span>
                <span class="value">{{ formatDuration(airportBasicData?.data?.acc_time || osdData?.data?.acc_time) }}</span>
              </div>
              <div class="info-item">
                <span class="label">工作电压:</span>
                <span class="value">{{ airportBasicData?.data?.working_voltage || osdData?.data?.working_voltage || '--' }}mV</span>
              </div>
              <div class="info-item">
                <span class="label">工作电流:</span>
                <span class="value">{{ airportBasicData?.data?.working_current || osdData?.data?.working_current || '--' }}mA</span>
              </div>
              <div class="info-item">
                <span class="label">备用电池:</span>
                <el-tag :type="(airportBasicData?.data?.backup_battery?.switch || osdData?.data?.backup_battery?.switch) === 1 ? 'success' : 'info'" size="small">
                  {{ (airportBasicData?.data?.backup_battery?.switch || osdData?.data?.backup_battery?.switch) === 1 ? '开启' : '关闭' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">备用电池电压:</span>
                <span class="value">{{ airportBasicData?.data?.backup_battery?.voltage || osdData?.data?.backup_battery?.voltage || '--' }}mV</span>
              </div>
              <div class="info-item">
                <span class="label">激活时间:</span>
                <span class="value">{{ formatTime(airportBasicData?.data?.activation_time || osdData?.data?.activation_time) }}</span>
              </div>
              <div class="info-item">
                <span class="label">机场状态:</span>
                <el-tag :type="getModeCodeType(airportBasicData?.data?.mode_code || osdData?.data?.mode_code)" size="small">
                  {{ getModeCodeText(airportBasicData?.data?.mode_code || osdData?.data?.mode_code) }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">Home点有效性:</span>
                <el-tag :type="getHomePositionValidType(airportBasicData?.data?.home_position_is_valid || osdData?.data?.home_position_is_valid)" size="small">
                  {{ getHomePositionValidText(airportBasicData?.data?.home_position_is_valid || osdData?.data?.home_position_is_valid) }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">机场朝向角:</span>
                <span class="value">{{ airportBasicData?.data?.heading || osdData?.data?.heading || '--' }}°</span>
              </div>
              <div class="info-item">
                <span class="label">空调状态:</span>
                <el-tag :type="getAirConditionerType(airportBasicData?.data?.air_conditioner?.air_conditioner_state || osdData?.data?.air_conditioner?.air_conditioner_state)" size="small">
                  {{ getAirConditionerText(airportBasicData?.data?.air_conditioner?.air_conditioner_state || osdData?.data?.air_conditioner?.air_conditioner_state) }}
                </el-tag>
              </div>
              <div class="info-item" v-if="(airportBasicData?.data?.air_conditioner?.switch_time || osdData?.data?.air_conditioner?.switch_time) !== undefined">
                <span class="label">空调切换时间:</span>
                <span class="value">{{ airportBasicData?.data?.air_conditioner?.switch_time || osdData?.data?.air_conditioner?.switch_time || '--' }} 秒</span>
              </div>
              <div class="info-item">
                <span class="label">空中回传:</span>
                <el-tag :type="(airportBasicData?.data?.air_transfer_enable || osdData?.data?.air_transfer_enable) ? 'success' : 'info'" size="small">
                  {{ (airportBasicData?.data?.air_transfer_enable || osdData?.data?.air_transfer_enable) ? '开启' : '关闭' }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <!-- 图传链路信息 -->
        <el-col :span="4">
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Monitor /></el-icon>
                <span>图传链路</span>
                <el-tag size="small" :type="wirelessLinkData ? (isDataStale(wirelessLinkData.timestamp) ? 'warning' : 'success') : 'info'">
                  {{ wirelessLinkData ? formatRelativeTime(wirelessLinkData.timestamp) : '等待' }}
                  <span v-if="wirelessLinkData && isDataStale(wirelessLinkData.timestamp)" style="margin-left: 4px;">⚠️</span>
                </el-tag>
              </div>
            </template>
            <div class="info-content">
              <div class="info-item">
                <span class="label">4G连接状态:</span>
                <el-tag :type="wirelessLinkData?.data?.wireless_link?.['4g_link_state'] === 1 ? 'success' : 'danger'" size="small">
                  {{ wirelessLinkData?.data?.wireless_link?.['4g_link_state'] === 1 ? '已连接' : '未连接' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">SDR连接状态:</span>
                <el-tag :type="wirelessLinkData?.data?.wireless_link?.sdr_link_state === 1 ? 'success' : 'danger'" size="small">
                  {{ wirelessLinkData?.data?.wireless_link?.sdr_link_state === 1 ? '已连接' : '未连接' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">4G信号质量:</span>
                <span class="value">{{ getSignalQuality(wirelessLinkData?.data?.wireless_link?.['4g_quality']) }}/5</span>
              </div>
              <div class="info-item">
                <span class="label">SDR信号质量:</span>
                <span class="value">{{ getSignalQuality(wirelessLinkData?.data?.wireless_link?.sdr_quality) }}/5</span>
              </div>
              <div class="info-item">
                <span class="label">4G频段:</span>
                <span class="value">{{ wirelessLinkData?.data?.wireless_link?.['4g_freq_band'] || '--' }}MHz</span>
              </div>
              <div class="info-item">
                <span class="label">SDR频段:</span>
                <span class="value">{{ wirelessLinkData?.data?.wireless_link?.sdr_freq_band || '--' }}MHz</span>
              </div>
              <div class="info-item">
                <span class="label">天端4G质量:</span>
                <span class="value">{{ getSignalQuality(wirelessLinkData?.data?.wireless_link?.['4g_uav_quality']) }}/5</span>
              </div>
              <div class="info-item">
                <span class="label">地端4G质量:</span>
                <span class="value">{{ getSignalQuality(wirelessLinkData?.data?.wireless_link?.['4g_gnd_quality']) }}/5</span>
              </div>
              <div class="info-item">
                <span class="label">Dongle数量:</span>
                <span class="value">{{ wirelessLinkData?.data?.wireless_link?.dongle_number || '--' }}个</span>
              </div>
              <div class="info-item">
                <span class="label">链路模式:</span>
                <el-tag :type="wirelessLinkData?.data?.wireless_link?.link_workmode === 0 ? 'info' : 'success'" size="small">
                  {{ wirelessLinkData?.data?.wireless_link?.link_workmode === 0 ? 'SDR模式' : '4G融合模式' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">待上传文件:</span>
                <span class="value">{{ wirelessLinkData?.data?.media_file_detail?.remain_upload || '--' }}个</span>
              </div>
              <div class="info-item">
                <span class="label">SDR下行质量:</span>
                <span class="value">{{ wirelessLinkData?.data?.sdr?.down_quality || '--' }}%</span>
              </div>
              <div class="info-item">
                <span class="label">SDR上行质量:</span>
                <span class="value">{{ wirelessLinkData?.data?.sdr?.up_quality || '--' }}%</span>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <!-- 直播信息 -->
        <el-col :span="4">
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Camera /></el-icon>
                <span>直播信息</span>
                <el-tag size="small" :type="(airportBasicData?.data?.live_capacity || osdData?.data?.live_capacity) ? 'success' : 'info'">
                  {{ (airportBasicData?.data?.live_capacity || osdData?.data?.live_capacity) ? '已连接' : '等待' }}
                </el-tag>
              </div>
            </template>
            <div class="info-content">
              <!-- 调试信息 -->
              <div class="info-item" style="background: var(--el-color-warning-light-9); padding: 5px; border-radius: 4px; margin-bottom: 8px; font-size: 12px;">
                <span class="label">数据状态:</span>
                <div style="margin-top: 2px;">
                  <div>airportBasicData: {{ airportBasicData ? '有数据' : '无数据' }}</div>
                  <div>osdData: {{ osdData ? '有数据' : '无数据' }}</div>
                  <div>wirelessLinkData: {{ wirelessLinkData ? '有数据' : '无数据' }}</div>
                  <div>live_capacity: {{ (airportBasicData?.data?.live_capacity || osdData?.data?.live_capacity || wirelessLinkData?.data?.live_capacity) ? '有数据' : '无数据' }}</div>
                  <div>live_status: {{ (wirelessLinkData?.data?.live_status || osdData?.data?.live_status)?.length > 0 ? '有数据' : '无数据' }}</div>
                </div>
              </div>
              
              <!-- 直播能力信息 -->
              <div class="info-item" v-if="(airportBasicData?.data?.live_capacity || osdData?.data?.live_capacity || wirelessLinkData?.data?.live_capacity)">
                <span class="label">可推流码流数:</span>
                <span class="value">{{ airportBasicData?.data?.live_capacity?.available_video_number || osdData?.data?.live_capacity?.available_video_number || wirelessLinkData?.data?.live_capacity?.available_video_number || '--' }}</span>
              </div>
              <div class="info-item" v-if="(airportBasicData?.data?.live_capacity || osdData?.data?.live_capacity || wirelessLinkData?.data?.live_capacity)">
                <span class="label">最大同时推流:</span>
                <span class="value">{{ airportBasicData?.data?.live_capacity?.coexist_video_number_max || osdData?.data?.live_capacity?.coexist_video_number_max || wirelessLinkData?.data?.live_capacity?.coexist_video_number_max || '--' }}</span>
              </div>
              
              <!-- 视频源设备列表 -->
              <div class="info-item" v-if="(airportBasicData?.data?.live_capacity?.device_list || osdData?.data?.live_capacity?.device_list || wirelessLinkData?.data?.live_capacity?.device_list)?.length > 0">
                <span class="label">视频源设备:</span>
                <div style="margin-top: 5px;">
                  <div v-for="(device, index) in (airportBasicData?.data?.live_capacity?.device_list || osdData?.data?.live_capacity?.device_list || wirelessLinkData?.data?.live_capacity?.device_list)" :key="index" style="margin-bottom: 8px; padding: 5px; background: var(--el-color-info-light-9); border-radius: 4px; font-size: 12px;">
                    <div style="font-weight: bold; color: var(--el-color-primary);">{{ device.sn }}</div>
                    <div style="margin-top: 2px;">
                      <span style="color: var(--el-color-success);">可推流: {{ device.available_video_number }}</span>
                      <span style="margin-left: 10px; color: var(--el-color-warning);">最大同时: {{ device.coexist_video_number_max }}</span>
                    </div>
                    <div v-if="device.camera_list?.length > 0" style="margin-top: 5px;">
                      <div v-for="(camera, camIndex) in device.camera_list" :key="camIndex" style="margin-left: 10px; margin-bottom: 3px; padding: 3px; background: var(--el-color-success-light-9); border-radius: 3px;">
                        <div style="font-weight: bold;">相机: {{ camera.camera_index }}</div>
                        <div style="font-size: 11px; color: var(--el-color-text-regular);">
                          可推流: {{ camera.available_video_number }} | 最大同时: {{ camera.coexist_video_number_max }}
                        </div>
                        <div v-if="camera.video_list?.length > 0" style="margin-top: 3px;">
                          <div v-for="(video, vidIndex) in camera.video_list" :key="vidIndex" style="margin-left: 10px; font-size: 10px; color: var(--el-color-text-secondary);">
                            {{ video.video_index }} ({{ video.video_type }})
                            <span v-if="video.switchable_video_types?.length > 0" style="color: var(--el-color-primary);">
                              [可切换: {{ video.switchable_video_types.join(', ') }}]
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- 当前直播状态 -->
              <div class="info-item" v-if="(wirelessLinkData?.data?.live_status || osdData?.data?.live_status)?.length > 0">
                <span class="label">当前直播状态:</span>
                <div style="margin-top: 5px;">
                  <div v-for="(stream, index) in (wirelessLinkData?.data?.live_status || osdData?.data?.live_status)" :key="index" style="margin-bottom: 3px; font-size: 12px;">
                    <el-tag :type="stream.status === 1 ? 'success' : 'info'" size="small">
                      {{ stream.video_id }} - {{ getVideoQualityText(stream.video_quality) }}
                    </el-tag>
                  </div>
                </div>
              </div>
              
              <!-- 无数据时的提示 -->
              <div v-if="!(airportBasicData?.data?.live_capacity || osdData?.data?.live_capacity) && !(wirelessLinkData?.data?.live_status || osdData?.data?.live_status)?.length" style="text-align: center; padding: 20px; color: var(--el-color-text-secondary);">
                <div style="font-size: 14px; margin-bottom: 8px;">📹 等待直播数据</div>
                <div style="font-size: 12px;">请确保设备已连接并发送直播相关信息</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <!-- 设备属性推送 -->
        <el-col :span="4">
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Connection /></el-icon>
                <span>设备属性</span>
                <el-tag size="small" type="success">实时</el-tag>
              </div>
            </template>
            <div class="info-content">
              <div class="info-item">
                <span class="label">状态数据:</span>
                <el-tag :type="stateDataStatus === 'connected' ? 'success' : 'info'" size="small">
                  {{ stateDataStatus === 'connected' ? '已连接' : '未连接' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">定频数据:</span>
                <el-tag :type="osdDataStatus === 'connected' ? 'success' : 'info'" size="small">
                  {{ osdDataStatus === 'connected' ? '已连接' : '未连接' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">状态主题:</span>
                <span class="value" style="font-size: 10px;">thing/product/{{ getAirportSn() }}/state</span>
              </div>
              <div class="info-item">
                <span class="label">OSD主题:</span>
                <span class="value" style="font-size: 10px;">thing/product/{{ getAirportSn() }}/osd</span>
              </div>
              <div class="info-item">
                <span class="label">上报频率:</span>
                <span class="value">0.5Hz</span>
              </div>
              <div class="info-item">
                <span class="label">方向:</span>
                <el-tag type="warning" size="small">up</el-tag>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <!-- 负载属性上报 -->
        <el-col :span="4">
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Camera /></el-icon>
                <span>负载属性</span>
                <el-tag size="small" type="info">相机</el-tag>
              </div>
            </template>
            <div class="info-content">
              <div class="info-item">
                <span class="label">云台俯仰:</span>
                <span class="value">{{ gimbalPitch || '--' }}°</span>
              </div>
              <div class="info-item">
                <span class="label">云台偏航:</span>
                <span class="value">{{ gimbalYaw || '--' }}°</span>
              </div>
              <div class="info-item">
                <span class="label">云台横滚:</span>
                <span class="value">{{ gimbalRoll || '--' }}°</span>
              </div>
              <div class="info-item">
                <span class="label">负载索引:</span>
                <span class="value">{{ formattedPayloadIndex }}</span>
              </div>
              <div class="info-item">
                <span class="label">产品类型:</span>
                <span class="value">{{ payloadType || '--' }}</span>
              </div>
              <div class="info-item">
                <span class="label">子类型:</span>
                <span class="value">{{ payloadSubtype || '--' }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
        
        
        <!-- 电池信息 -->
        <el-col :span="4">
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Lightning /></el-icon>
                <span>电池信息</span>
                <el-tag size="small" type="warning">监控</el-tag>
              </div>
            </template>
            <div class="info-content">
              <div class="info-item">
                <span class="label">工作电压:</span>
                <span class="value">{{ osdData?.data?.working_voltage || '--' }}mV</span>
              </div>
              <div class="info-item">
                <span class="label">工作电流:</span>
                <span class="value">{{ osdData?.data?.working_current || '--' }}mA</span>
              </div>
              <div class="info-item">
                <span class="label">电池模式:</span>
                <el-tag :type="getBatteryStoreModeType(osdData?.data?.battery_store_mode)" size="small">
                  {{ getBatteryStoreModeText(osdData?.data?.battery_store_mode) }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">备用电池:</span>
                <el-tag :type="osdData?.data?.backup_battery?.switch === 1 ? 'success' : 'info'" size="small">
                  {{ osdData?.data?.backup_battery?.switch === 1 ? '开启' : '关闭' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">备用电池电压:</span>
                <span class="value">{{ osdData?.data?.backup_battery?.voltage || '--' }}mV</span>
              </div>
              <div class="info-item">
                <span class="label">备用电池温度:</span>
                <span class="value">{{ osdData?.data?.backup_battery?.temperature || '--' }}°C</span>
              </div>
              <div class="info-item">
                <span class="label">供电电压:</span>
                <span class="value">{{ osdData?.data?.electric_supply_voltage || '--' }}V</span>
              </div>
              <div class="info-item">
                <span class="label">ACDC功率:</span>
                <span class="value">{{ osdData?.data?.acdc_power_input ? osdData.data.acdc_power_input.toFixed(1) : '--' }}W</span>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <!-- 其他机场信息卡片... -->
      </el-row>
    </div>

    <!-- 原始数据展示按钮 -->
    <div style="position: fixed; top: 50%; right: 20px; transform: translateY(-50%); z-index: 1000;">
      <el-button 
        type="primary" 
        :icon="DataLine" 
        circle 
        size="large"
        @click="showRawDataDrawer = true"
        title="查看原始数据"
      />
    </div>
    
    <!-- 调试按钮 -->
    <div style="position: fixed; top: 60%; right: 20px; transform: translateY(-50%); z-index: 1000;">
      <el-button 
        type="warning" 
        circle 
        size="large"
        @click="handleOSDMessage"
        title="手动处理OSD消息"
      >
        🔄
      </el-button>
    </div>

    <!-- 原始数据抽屉 -->
    <el-drawer
      v-model="showRawDataDrawer"
      title="原始数据展示"
      direction="rtl"
      size="60%"
      :before-close="handleDrawerClose"
    >
      <div class="info-content">
        <el-tabs type="border-card">
          <el-tab-pane label="OSD数据" name="osd">
            <pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; max-height: 500px; overflow-y: auto; font-size: 12px;">{{ JSON.stringify(osdData, null, 2) }}</pre>
          </el-tab-pane>
          <el-tab-pane label="机场基础信息" name="basic">
            <pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; max-height: 500px; overflow-y: auto; font-size: 12px;">{{ JSON.stringify(airportBasicData, null, 2) }}</pre>
          </el-tab-pane>
          <el-tab-pane label="图传链路信息" name="wireless">
            <pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; max-height: 500px; overflow-y: auto; font-size: 12px;">{{ JSON.stringify(wirelessLinkData, null, 2) }}</pre>
          </el-tab-pane>
          <el-tab-pane label="环境状态信息" name="environment">
            <pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; max-height: 500px; overflow-y: auto; font-size: 12px;">{{ JSON.stringify(environmentData, null, 2) }}</pre>
          </el-tab-pane>
          <el-tab-pane label="飞机信息" name="aircraft">
            <pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; max-height: 500px; overflow-y: auto; font-size: 12px;">{{ JSON.stringify(aircraftData, null, 2) }}</pre>
          </el-tab-pane>
          <el-tab-pane label="设备状态" name="state">
            <pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; max-height: 500px; overflow-y: auto; font-size: 12px;">{{ JSON.stringify(stateData, null, 2) }}</pre>
          </el-tab-pane>
          <el-tab-pane label="数据更新信息" name="updates">
                <div style="padding: 10px;">
                  <div style="margin-bottom: 15px;">
                    <h4 style="color: var(--el-color-primary); margin-bottom: 10px;">📊 数据更新状态</h4>
                    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 10px;">
                      <div style="padding: 8px; background: var(--el-color-info-light-9); border-radius: 4px;">
                        <strong>OSD数据:</strong> 
                        <span v-if="osdData?.timestamp" :style="isDataStale(osdData.timestamp) ? 'color: var(--el-color-warning);' : 'color: var(--el-color-success);'">
                          {{ formatRelativeTime(osdData.timestamp) }}
                          <span v-if="isDataStale(osdData.timestamp)" style="margin-left: 4px;">⚠️</span>
                        </span>
                        <span v-else style="color: var(--el-color-warning);">等待数据</span>
                      </div>
                      <div style="padding: 8px; background: var(--el-color-info-light-9); border-radius: 4px;">
                        <strong>机场基础信息:</strong> 
                        <span v-if="airportBasicData?.timestamp" :style="isDataStale(airportBasicData.timestamp) ? 'color: var(--el-color-warning);' : 'color: var(--el-color-success);'">
                          {{ formatRelativeTime(airportBasicData.timestamp) }}
                          <span v-if="isDataStale(airportBasicData.timestamp)" style="margin-left: 4px;">⚠️</span>
                        </span>
                        <span v-else style="color: var(--el-color-warning);">等待数据</span>
                        <div v-if="airportBasicData?.timestamp" style="font-size: 10px; color: #666; margin-top: 2px;">
                          原始值: {{ airportBasicData.timestamp }} (类型: {{ typeof airportBasicData.timestamp }})
                        </div>
                      </div>
                      <div style="padding: 8px; background: var(--el-color-info-light-9); border-radius: 4px;">
                        <strong>图传链路信息:</strong> 
                        <span v-if="wirelessLinkData?.timestamp" :style="isDataStale(wirelessLinkData.timestamp) ? 'color: var(--el-color-warning);' : 'color: var(--el-color-success);'">
                          {{ formatRelativeTime(wirelessLinkData.timestamp) }}
                          <span v-if="isDataStale(wirelessLinkData.timestamp)" style="margin-left: 4px;">⚠️</span>
                        </span>
                        <span v-else style="color: var(--el-color-warning);">等待数据</span>
                      </div>
                      <div style="padding: 8px; background: var(--el-color-info-light-9); border-radius: 4px;">
                        <strong>环境状态信息:</strong> 
                        <span v-if="environmentData?.timestamp" :style="isDataStale(environmentData.timestamp) ? 'color: var(--el-color-warning);' : 'color: var(--el-color-success);'">
                          {{ formatRelativeTime(environmentData.timestamp) }}
                          <span v-if="isDataStale(environmentData.timestamp)" style="margin-left: 4px;">⚠️</span>
                        </span>
                        <span v-else style="color: var(--el-color-warning);">等待数据</span>
                      </div>
                      <div style="padding: 8px; background: var(--el-color-info-light-9); border-radius: 4px;">
                        <strong>飞机信息:</strong> 
                        <span v-if="aircraftData?.timestamp" :style="isDataStale(aircraftData.timestamp) ? 'color: var(--el-color-warning);' : 'color: var(--el-color-success);'">
                          {{ formatRelativeTime(aircraftData.timestamp) }}
                          <span v-if="isDataStale(aircraftData.timestamp)" style="margin-left: 4px;">⚠️</span>
                        </span>
                        <span v-else style="color: var(--el-color-warning);">等待数据</span>
                      </div>
                      <div style="padding: 8px; background: var(--el-color-info-light-9); border-radius: 4px;">
                        <strong>设备状态:</strong> 
                        <span v-if="stateData?.timestamp" :style="isDataStale(stateData.timestamp) ? 'color: var(--el-color-warning);' : 'color: var(--el-color-success);'">
                          {{ formatRelativeTime(stateData.timestamp) }}
                          <span v-if="isDataStale(stateData.timestamp)" style="margin-left: 4px;">⚠️</span>
                        </span>
                        <span v-else style="color: var(--el-color-warning);">等待数据</span>
                      </div>
                    </div>
                  </div>
                  <div>
                    <h4 style="color: var(--el-color-primary); margin-bottom: 10px;">🔄 最近更新字段</h4>
                    <div v-if="osdData?.data" style="background: var(--el-color-success-light-9); padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                      <strong>OSD数据字段:</strong> {{ Object.keys(osdData.data).join(', ') }}
                    </div>
                    <div v-if="airportBasicData?.data" style="background: var(--el-color-info-light-9); padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                      <strong>机场基础信息字段:</strong> {{ Object.keys(airportBasicData.data).join(', ') }}
                    </div>
                    <div v-if="wirelessLinkData?.data" style="background: var(--el-color-warning-light-9); padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                      <strong>图传链路信息字段:</strong> {{ Object.keys(wirelessLinkData.data).join(', ') }}
                    </div>
                    <div v-if="environmentData?.data" style="background: var(--el-color-danger-light-9); padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                      <strong>环境状态信息字段:</strong> {{ Object.keys(environmentData.data).join(', ') }}
                    </div>
                    <div v-if="aircraftData?.data" style="background: var(--el-color-success-light-9); padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                      <strong>飞机信息字段:</strong> {{ Object.keys(aircraftData.data).join(', ') }}
                    </div>
                    <div v-if="stateData?.data" style="background: var(--el-color-primary-light-9); padding: 10px; border-radius: 4px;">
                      <strong>设备状态字段:</strong> {{ Object.keys(stateData.data).join(', ') }}
                    </div>
                  </div>
                </div>
              </el-tab-pane>
          <el-tab-pane label="所有消息历史" name="all">
            <pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; max-height: 500px; overflow-y: auto; font-size: 12px;">{{ JSON.stringify(mqttProxyStore.messageHistory, null, 2) }}</pre>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-drawer>

    <!-- 等待数据状态 -->
    <div v-if="!hasReceivedData" class="waiting-data" style="margin: 20px 0; padding: 20px; text-align: center; background: var(--el-color-warning-light-9); border-radius: 8px;">
      <el-icon style="color: var(--el-color-warning); font-size: 48px; margin-bottom: 16px;"><Loading /></el-icon>
      <h3 style="color: var(--el-color-warning-dark-2); margin-bottom: 8px;">等待接收机场OSD数据...</h3>
      <p style="color: var(--el-text-color-regular); font-size: 14px;">请确保MQTT连接正常，设备正在上报数据</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Connection, Camera, Loading, Lightning, Sunny, Cloudy, Location, Monitor, Setting, DataLine } from '@element-plus/icons-vue'
import { useMqttProxyStore } from '@/stores/mqtt-proxy'
import { useDeviceStore } from '@/stores/device'

// Props
const props = defineProps({
  deviceSn: String,
  airportSn: String
})

// Emits
const emit = defineEmits(['dataReceived', 'subscriptionStatus'])

// Stores
const mqttProxyStore = useMqttProxyStore()
const deviceStore = useDeviceStore()

// 响应式数据
const osdData = ref(null)
const hasReceivedData = ref(false)
const subscriptionStatus = ref('未订阅')
const subscribing = ref(false)
const stateDataStatus = ref('disconnected')
const osdDataStatus = ref('disconnected')
const stateData = ref(null)
const gimbalPitch = ref(null)
const showRawDataDrawer = ref(false)
const gimbalYaw = ref(null)
const gimbalRoll = ref(null)
const payloadIndex = ref(null)
const payloadType = ref(null)
const payloadSubtype = ref(null)

// 分类存储不同类型的数据
const airportBasicData = ref(null)      // 机场基础信息
const wirelessLinkData = ref(null)      // 图传链路信息
const environmentData = ref(null)       // 环境状态信息
const aircraftData = ref(null)          // 飞机信息

// 计算属性
const osdTopic = computed(() => {
  const airportSn = props.airportSn || deviceStore.currentDevice?.airport_sn || props.deviceSn
  return `thing/product/${airportSn}/osd`
})

const airportOsdTopic = computed(() => {
  return osdTopic.value
})

// 格式化的负载索引显示
const formattedPayloadIndex = computed(() => {
  if (payloadType.value && payloadSubtype.value && payloadIndex.value) {
    return `${payloadType.value}-${payloadSubtype.value}-${payloadIndex.value}`
  }
  return payloadIndex.value || '--'
})

// 获取机场序列号
const getAirportSn = () => {
  return props.airportSn || 
         deviceStore.currentDevice?.airport_sn || 
         props.deviceSn || 
         '--'
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '--'
  return new Date(timestamp).toLocaleString()
}

// 数据分类存储函数
const classifyAndStoreData = (data, message) => {
  if (!data) return
  
  // 机场基础信息：包含作业次数、运行时间、电压电流等
  if (data.job_number !== undefined || data.acc_time !== undefined || 
      data.working_voltage !== undefined || data.working_current !== undefined ||
      data.backup_battery !== undefined) {
    airportBasicData.value = { ...message, data }
    console.log('📊 机场基础信息已更新:', data)
  }
  
  // 图传链路信息：包含4G/SDR连接状态、频率等
  if (data.wireless_link !== undefined || data.flighttask_prepare_capacity !== undefined ||
      data.media_file_detail !== undefined) {
    wirelessLinkData.value = { ...message, data }
    console.log('📡 图传链路信息已更新:', data)
  }
  
  // 环境状态信息：包含温度湿度、位置、网络状态等
  if (data.temperature !== undefined || data.humidity !== undefined || 
      data.latitude !== undefined || data.longitude !== undefined ||
      data.network_state !== undefined || data.drone_in_dock !== undefined) {
    environmentData.value = { ...message, data }
    console.log('🌡️ 环境状态信息已更新:', data)
  }
  
  // 飞机信息：包含电池、位置、维护状态等
  if (data.battery !== undefined || data.firmware_version !== undefined ||
      data.maintain_status !== undefined || data.position_state !== undefined) {
    aircraftData.value = { ...message, data }
    console.log('✈️ 飞机信息已更新:', data)
  }
}

// 格式化更新时间
const formatUpdateTime = (timestamp) => {
  if (!timestamp) return '--'
  const date = new Date(timestamp)
  return date.toLocaleTimeString()
}

// 格式化相对时间
const formatRelativeTime = (timestamp) => {
  if (!timestamp) return '--'
  
  // 确保时间戳是数字格式
  let timestampNum
  if (typeof timestamp === 'string') {
    timestampNum = new Date(timestamp).getTime()
  } else {
    timestampNum = timestamp
  }
  
  // 检查时间戳是否有效
  if (isNaN(timestampNum) || timestampNum <= 0) {
    return '无效时间'
  }
  
  const now = Date.now()
  const diff = Math.floor((now - timestampNum) / 1000)
  
  // 检查差值是否合理（不能是负数或过大）
  if (diff < 0) {
    return '未来时间'
  } else if (diff > 365 * 24 * 3600) { // 超过1年
    return '很久以前'
  }
  
  if (diff < 60) {
    return `${diff}秒前`
  } else if (diff < 3600) {
    return `${Math.floor(diff / 60)}分钟前`
  } else if (diff < 86400) {
    return `${Math.floor(diff / 3600)}小时前`
  } else {
    return `${Math.floor(diff / 86400)}天前`
  }
}

// 检查数据是否过期（超过5分钟认为过期）
const isDataStale = (timestamp) => {
  if (!timestamp) return true
  
  // 确保时间戳是数字格式
  let timestampNum
  if (typeof timestamp === 'string') {
    timestampNum = new Date(timestamp).getTime()
  } else {
    timestampNum = timestamp
  }
  
  // 检查时间戳是否有效
  if (isNaN(timestampNum) || timestampNum <= 0) {
    return true // 无效时间戳认为过期
  }
  
  const now = Date.now()
  const diff = Math.floor((now - timestampNum) / 1000)
  return diff > 300 // 5分钟
}

// 格式化持续时间
const formatDuration = (seconds) => {
  if (!seconds) return '--'
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${hours}小时${minutes}分钟`
}

// 抽屉关闭处理
const handleDrawerClose = (done) => {
  showRawDataDrawer.value = false
  done()
}

// 处理信号质量异常值
const getSignalQuality = (quality) => {
  if (!quality || quality === 65535 || quality < 0 || quality > 5) {
    return '--'
  }
  return quality
}

// 机场状态相关方法
const getModeCodeType = (modeCode) => {
  const typeMap = {
    0: 'success',  // 空闲中
    1: 'info',     // 现场调试
    2: 'info',     // 远程调试
    3: 'warning',  // 固件升级中
    4: 'primary',  // 作业中
    5: 'warning'   // 待标定
  }
  return typeMap[modeCode] || 'info'
}

const getModeCodeText = (modeCode) => {
  const textMap = {
    0: '空闲中',
    1: '现场调试',
    2: '远程调试',
    3: '固件升级中',
    4: '作业中',
    5: '待标定'
  }
  return textMap[modeCode] || '未知状态'
}

// Home点有效性相关方法
const getHomePositionValidType = (validCode) => {
  const typeMap = {
    0: 'danger',   // 航向和经纬度坐标都无效
    1: 'success',  // 航向和经纬度坐标都有效
    2: 'warning',  // 航向有效，经纬度无效
    3: 'warning'   // 经纬度有效，航向无效
  }
  return typeMap[validCode] || 'info'
}

const getHomePositionValidText = (validCode) => {
  const textMap = {
    0: '航向和经纬度都无效',
    1: '航向和经纬度都有效',
    2: '航向有效，经纬度无效',
    3: '经纬度有效，航向无效'
  }
  return textMap[validCode] || '未知状态'
}

// 空调状态相关方法
const getAirConditionerType = (stateCode) => {
  const typeMap = {
    0: 'info',      // 空闲模式
    1: 'primary',   // 制冷模式
    2: 'warning',   // 制热模式
    3: 'success',   // 除湿模式
    4: 'info',      // 制冷退出模式
    5: 'info',      // 制热退出模式
    6: 'info',      // 除湿退出模式
    7: 'primary',   // 制冷准备模式
    8: 'warning',   // 制热准备模式
    9: 'success',   // 除湿准备模式
    10: 'primary',  // 风冷准备中
    11: 'primary',  // 风冷中
    12: 'info',     // 风冷退出中
    13: 'success',  // 除雾准备中
    14: 'success',  // 除雾中
    15: 'info'      // 除雾退出中
  }
  return typeMap[stateCode] || 'info'
}

const getAirConditionerText = (stateCode) => {
  const textMap = {
    0: '空闲模式',
    1: '制冷模式',
    2: '制热模式',
    3: '除湿模式',
    4: '制冷退出模式',
    5: '制热退出模式',
    6: '除湿退出模式',
    7: '制冷准备模式',
    8: '制热准备模式',
    9: '除湿准备模式',
    10: '风冷准备中',
    11: '风冷中',
    12: '风冷退出中',
    13: '除雾准备中',
    14: '除雾中',
    15: '除雾退出中'
  }
  return textMap[stateCode] || '未知状态'
}

// 视频质量相关方法
const getVideoQualityText = (qualityCode) => {
  const qualityMap = {
    0: '自适应',
    1: '流畅',
    2: '标清',
    3: '高清',
    4: '超清'
  }
  return qualityMap[qualityCode] || '未知质量'
}


const getFlightTaskStepType = (stepCode) => {
  const typeMap = {
    0: 'info',     // 作业准备中
    1: 'primary',  // 飞行作业中
    2: 'info',     // 作业后状态恢复
    3: 'warning',  // 自定义飞行区更新中
    4: 'warning',  // 地形障碍物更新中
    5: 'success',  // 任务空闲
    255: 'danger', // 飞行器异常
    256: 'info'    // 未知状态
  }
  return typeMap[stepCode] || 'info'
}

const getFlightTaskStepText = (stepCode) => {
  const textMap = {
    0: '作业准备中',
    1: '飞行作业中',
    2: '作业后状态恢复',
    3: '自定义飞行区更新中',
    4: '地形障碍物更新中',
    5: '任务空闲',
    255: '飞行器异常',
    256: '未知状态'
  }
  return textMap[stepCode] || '未知状态'
}

const getDroneInDockType = (inDock) => {
  return inDock === 1 ? 'success' : 'danger'
}

const getDroneInDockText = (inDock) => {
  return inDock === 1 ? '舱内' : '舱外'
}

const getCoverStateType = (coverState) => {
  const typeMap = {
    0: 'info',      // 关闭
    1: 'success',   // 打开
    2: 'warning',   // 半开
    3: 'danger'     // 舱盖状态异常
  }
  return typeMap[coverState] || 'info'
}

const getCoverStateText = (coverState) => {
  const textMap = {
    0: '关闭',
    1: '打开',
    2: '半开',
    3: '舱盖状态异常'
  }
  return textMap[coverState] || '未知状态'
}

// 环境监控相关方法
const getRainfallType = (rainfall) => {
  const typeMap = {
    0: 'success',  // 无雨
    1: 'info',     // 小雨
    2: 'warning',  // 中雨
    3: 'danger'    // 大雨
  }
  return typeMap[rainfall] || 'info'
}

const getRainfallText = (rainfall) => {
  const textMap = {
    0: '无雨',
    1: '小雨',
    2: '中雨',
    3: '大雨'
  }
  return textMap[rainfall] || '未知'
}

// 电池信息相关方法
const getBatteryStoreModeType = (mode) => {
  const typeMap = {
    1: 'success',  // 计划模式
    2: 'warning'   // 待命模式
  }
  return typeMap[mode] || 'info'
}

const getBatteryStoreModeText = (mode) => {
  const textMap = {
    1: '计划模式',
    2: '待命模式'
  }
  return textMap[mode] || '未知模式'
}

// 获取网络质量类型
const getNetworkQualityType = (quality) => {
  if (quality === 0) return 'danger'
  if (quality === 1) return 'danger'
  if (quality === 2) return 'warning'
  if (quality === 3) return 'warning'
  if (quality === 4) return 'success'
  if (quality === 5) return 'success'
  return 'info'
}

// 获取网络质量文本
const getNetworkQualityText = (quality) => {
  if (quality === 0) return '无信号'
  if (quality === 1) return '差'
  if (quality === 2) return '较差'
  if (quality === 3) return '一般'
  if (quality === 4) return '较好'
  if (quality === 5) return '好'
  return '未知'
}

// 订阅机场OSD数据
const subscribeToOSD = async () => {
  if (!mqttProxyStore.isConnected) {
    ElMessage.warning('MQTT未连接，无法订阅机场OSD数据')
    subscriptionStatus.value = 'MQTT未连接'
    return
  }

  subscribing.value = true
  subscriptionStatus.value = '订阅中...'

  try {
    await mqttProxyStore.subscribeToTopics(osdTopic.value, 1)
    osdDataStatus.value = 'connected'
    subscriptionStatus.value = '已订阅'
    ElMessage.success(`已订阅机场OSD主题: ${osdTopic.value}`)
    console.log('机场OSD订阅成功:', osdTopic.value)
    emit('subscriptionStatus', 'success')
  } catch (error) {
    subscriptionStatus.value = '订阅失败'
    console.error('订阅机场OSD主题失败:', error)
    ElMessage.error('订阅机场OSD主题失败')
    osdDataStatus.value = 'disconnected'
    emit('subscriptionStatus', 'error')
  } finally {
    subscribing.value = false
  }
}

// 订阅机场OSD数据（别名方法）
const subscribeToAirportOSD = subscribeToOSD

// 订阅设备状态数据
const subscribeToDeviceState = async () => {
  if (!mqttProxyStore.isConnected) {
    ElMessage.warning('MQTT未连接，无法订阅设备状态数据')
    return
  }

  try {
    const airportSn = getAirportSn()
    const stateTopic = `thing/product/${airportSn}/state`
    
    await mqttProxyStore.subscribeToTopics(stateTopic, 1)
    stateDataStatus.value = 'connected'
    ElMessage.success(`已订阅设备状态主题: ${stateTopic}`)
    
    console.log('设备状态订阅成功:', {
      stateTopic,
      airportSn,
      direction: 'up',
      frequency: '状态变化时上报'
    })
  } catch (error) {
    console.error('订阅设备状态主题失败:', error)
    ElMessage.error('订阅设备状态主题失败')
    stateDataStatus.value = 'disconnected'
  }
}

// 处理OSD消息
const handleOSDMessage = () => {
  try {
    console.log('🔍 检查机场OSD消息:', {
      osdTopic: osdTopic.value,
      messageHistoryLength: mqttProxyStore.messageHistory.length,
      hasReceivedData: hasReceivedData.value,
      currentTimestamp: osdData.value?.timestamp,
      allTopics: mqttProxyStore.messageHistory.map(msg => msg.topic)
    })
    
    if (!osdTopic.value) {
      console.log('❌ 机场OSD主题为空')
      return
    }
    
    // 查找所有匹配主题的消息
    const matchingMessages = mqttProxyStore.messageHistory.filter(msg => 
      msg.topic === osdTopic.value
    )
    
    console.log('📨 匹配主题的消息数量:', matchingMessages.length)
    console.log('📨 匹配的消息:', matchingMessages)
    
    if (matchingMessages.length > 0) {
      // 获取最新的消息
      const latestMessage = matchingMessages[0]
      console.log('📨 最新消息:', latestMessage)
      
      // 检查是否是新消息
      const isNewMessage = !osdData.value || 
                          latestMessage.timestamp > osdData.value.timestamp
      
      console.log('🆕 是否为新消息:', isNewMessage, {
        currentTimestamp: osdData.value?.timestamp,
        newTimestamp: latestMessage.timestamp
      })
      
      if (isNewMessage) {
        console.log('✅ 处理新的机场OSD数据:', latestMessage)
        
        // 解析消息数据
        let newData = null
        try {
          newData = JSON.parse(latestMessage.payload)
          console.log('📋 解析后的OSD数据:', newData)
        } catch (error) {
          console.error('❌ 解析OSD数据失败:', error)
          return
        }
        
        // 根据数据内容分类存储
        classifyAndStoreData(newData.data, latestMessage)
        
        // 局部更新：合并新数据到现有数据中
        if (osdData.value) {
          // 如果已有数据，进行局部更新
          const existingData = osdData.value.data || {}
          const mergedData = { ...existingData, ...newData.data }
          
          osdData.value = {
            ...latestMessage,
            data: mergedData
          }
          
          console.log('🔄 局部更新OSD数据:', {
            updatedFields: Object.keys(newData.data || {}),
            mergedData: mergedData
          })
        } else {
          // 如果没有现有数据，直接设置
          osdData.value = latestMessage
        }
        
        hasReceivedData.value = true
        subscriptionStatus.value = '已接收数据'
        emit('dataReceived', latestMessage)
        
        // 数据更新完成，时间戳已显示在标题旁边
        console.log('✅ 机场数据已更新，时间戳显示在标题旁边')
      } else {
        console.log('⏭️ 消息已处理过，跳过')
      }
    } else {
      console.log('📭 没有找到匹配主题的消息')
    }
  } catch (error) {
    console.error('❌ 处理机场OSD消息失败:', error)
  }
}

// 处理设备状态消息
const handleStateMessage = () => {
  try {
    const stateTopic = `thing/product/${getAirportSn()}/state`
    
    console.log('🔍 检查设备状态消息:', {
      stateTopic,
      messageHistoryLength: mqttProxyStore.messageHistory.length,
      stateDataStatus: stateDataStatus.value,
      currentTimestamp: stateData.value?.timestamp
    })
    
    // 显示所有消息的主题，用于调试
    console.log('📋 所有消息主题:', mqttProxyStore.messageHistory.map(msg => ({
      topic: msg.topic,
      timestamp: msg.timestamp,
      payload: msg.payload?.substring(0, 100) + '...'
    })))
    
    // 查找所有匹配主题的消息
    const matchingMessages = mqttProxyStore.messageHistory.filter(msg => 
      msg.topic === stateTopic
    )
    
    console.log('📨 匹配状态主题的消息数量:', matchingMessages.length)
    console.log('📨 匹配的状态消息:', matchingMessages)
    
    // 如果没有匹配的消息，尝试查找包含 'state' 的消息
    if (matchingMessages.length === 0) {
      const stateRelatedMessages = mqttProxyStore.messageHistory.filter(msg => 
        msg.topic.includes('state') || msg.topic.includes('thing/product')
      )
      console.log('🔍 查找包含state或thing/product的消息:', stateRelatedMessages)
    }
    
    if (matchingMessages.length > 0) {
      // 获取最新的消息
      const latestStateMessage = matchingMessages[0]
      console.log('📨 最新状态消息:', latestStateMessage)
      
      // 检查是否是新消息
      const isNewMessage = !stateData.value || 
                          latestStateMessage.timestamp > stateData.value.timestamp
      
      console.log('🆕 是否为新状态消息:', isNewMessage, {
        currentTimestamp: stateData.value?.timestamp,
        newTimestamp: latestStateMessage.timestamp
      })
      
      if (isNewMessage) {
        console.log('✅ 处理新的设备状态数据:', latestStateMessage)
        
        // 解析消息数据
        let newStateData = null
        try {
          newStateData = JSON.parse(latestStateMessage.payload)
          console.log('📋 解析后的状态数据:', newStateData)
        } catch (error) {
          console.error('❌ 解析状态数据失败:', error)
          return
        }
        
        // 局部更新：合并新数据到现有状态数据中
        if (stateData.value) {
          const existingData = stateData.value.data || {}
          const mergedData = { ...existingData, ...newStateData.data }
          
          stateData.value = {
            ...latestStateMessage,
            data: mergedData
          }
          
          console.log('🔄 局部更新状态数据:', {
            updatedFields: Object.keys(newStateData.data || {}),
            mergedData: mergedData
          })
        } else {
          stateData.value = latestStateMessage
        }
        
        stateDataStatus.value = 'connected'
        
        // 解析负载属性数据 - 根据DJI Dock3规范
        if (newStateData.data && newStateData.data.payload) {
          const payload = newStateData.data.payload
          
          // 解析云台角度信息
          gimbalPitch.value = payload.gimbal_pitch || payload.pitch
          gimbalYaw.value = payload.gimbal_yaw || payload.yaw
          gimbalRoll.value = payload.gimbal_roll || payload.roll
          
          // 解析负载索引和类型信息
          payloadIndex.value = payload.payload_index || payload.index
          payloadType.value = payload.type
          payloadSubtype.value = payload.subtype
          
          // 根据DJI规范，负载索引格式为: type-subtype-gimbalIndex
          const payloadIndexFormatted = payloadType.value && payloadSubtype.value && payloadIndex.value 
            ? `${payloadType.value}-${payloadSubtype.value}-${payloadIndex.value}`
            : payloadIndex.value
          
          console.log('📷 DJI Dock3 负载属性数据更新:', {
            // 云台角度信息
            gimbalPitch: gimbalPitch.value,
            gimbalYaw: gimbalYaw.value,
            gimbalRoll: gimbalRoll.value,
            
            // 负载信息
            payloadIndex: payloadIndexFormatted,
            payloadType: payloadType.value,
            payloadSubtype: payloadSubtype.value,
            
            // 原始数据
            rawPayload: payload,
            
            // 消息元数据
            timestamp: latestStateMessage.timestamp,
            topic: latestStateMessage.topic
          })
          
          ElMessage.success(`负载属性已更新: ${payloadIndexFormatted}`)
        } else {
          console.log('⚠️ 状态消息中没有负载属性数据')
          console.log('📋 消息结构:', {
            hasData: !!newStateData.data,
            hasPayload: !!(newStateData.data && newStateData.data.payload),
            messageKeys: Object.keys(latestStateMessage),
            dataKeys: newStateData.data ? Object.keys(newStateData.data) : []
          })
        }
      } else {
        console.log('⏭️ 状态消息已处理过，跳过')
      }
    } else {
      console.log('📭 没有找到匹配状态主题的消息')
    }
  } catch (error) {
    console.error('❌ 处理设备状态数据失败:', error)
  }
}

// 消息监听器
let messageListener = null

const startMessageListener = () => {
  if (messageListener) {
    clearInterval(messageListener)
  }
  
  messageListener = setInterval(() => {
    handleOSDMessage()
    handleStateMessage()
  }, 2000)
}

const stopMessageListener = () => {
  if (messageListener) {
    clearInterval(messageListener)
    messageListener = null
  }
}

// 检查状态消息
const checkStateMessages = () => {
  console.log('🔍 手动检查状态消息')
  console.log('📊 当前状态:', {
    mqttConnected: mqttProxyStore.isConnected,
    messageHistoryLength: mqttProxyStore.messageHistory.length,
    stateTopic: `thing/product/${getAirportSn()}/state`,
    stateDataStatus: stateDataStatus.value
  })
  
  // 显示所有消息
  console.log('📋 所有消息:', mqttProxyStore.messageHistory)
  
  // 查找状态相关消息
  const stateTopic = `thing/product/${getAirportSn()}/state`
  const stateMessages = mqttProxyStore.messageHistory.filter(msg => 
    msg.topic === stateTopic
  )
  
  console.log('📨 状态主题消息:', stateMessages)
  
  if (stateMessages.length > 0) {
    ElMessage.success(`找到 ${stateMessages.length} 条状态消息`)
    handleStateMessage()
  } else {
    ElMessage.warning('未找到状态消息')
    
    // 显示所有主题
    const allTopics = mqttProxyStore.messageHistory.map(msg => msg.topic)
    console.log('📋 所有消息主题:', allTopics)
  }
}

// 刷新状态
const refreshStatus = () => {
  console.log('🔄 刷新机场订阅状态')
  console.log('📊 当前状态:', {
    mqttConnected: mqttProxyStore.isConnected,
    messageHistoryLength: mqttProxyStore.messageHistory.length,
    airportOsdTopic: airportOsdTopic.value,
    hasReceivedData: hasReceivedData.value
  })
  
  if (mqttProxyStore.isConnected) {
    subscriptionStatus.value = 'MQTT已连接，等待数据...'
    console.log('✅ MQTT已连接')
    
    // 检查是否有消息历史
    if (mqttProxyStore.messageHistory.length > 0) {
      console.log('📨 消息历史:', mqttProxyStore.messageHistory)
      
      // 检查是否有匹配的消息
      const matchingMessages = mqttProxyStore.messageHistory.filter(msg => 
        msg.topic === airportOsdTopic.value
      )
      
      if (matchingMessages.length > 0) {
        console.log('🎯 找到匹配的消息:', matchingMessages)
        subscriptionStatus.value = `找到 ${matchingMessages.length} 条消息，正在处理...`
        // 立即处理消息
        handleOSDMessage()
      } else {
        console.log('❌ 没有找到匹配主题的消息')
        subscriptionStatus.value = 'MQTT已连接，但未找到匹配的消息'
      }
    } else {
      console.log('📭 消息历史为空')
      subscriptionStatus.value = 'MQTT已连接，消息历史为空'
    }
  } else {
    subscriptionStatus.value = 'MQTT未连接'
    console.log('❌ MQTT未连接')
  }
}

// 组件挂载
onMounted(() => {
  console.log('机场信息组件挂载，开始订阅数据')
  subscribeToOSD()
  subscribeToDeviceState()
  startMessageListener()
})

// 组件卸载
onUnmounted(() => {
  console.log('机场信息组件卸载，停止监听')
  stopMessageListener()
  
  // 取消订阅
  if (osdTopic.value) {
    mqttProxyStore.unsubscribeTopics(osdTopic.value)
  }
  
  const stateTopic = `thing/product/${getAirportSn()}/state`
  mqttProxyStore.unsubscribeTopics(stateTopic)
})
</script>

<style scoped>
.airport-info-tab {
  padding: 16px;
}

.data-container {
  margin-top: 16px;
}

.data-grid {
  margin-bottom: 16px;
}

.info-card {
  height: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.info-content {
  padding: 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.info-item:last-child {
  border-bottom: none;
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

.waiting-container {
  margin: 20px 0;
}

.waiting-container code {
  background: var(--el-color-info-light-9);
  color: var(--el-color-info);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}
</style>
