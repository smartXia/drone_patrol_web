<template>
  <div class="aircraft-info-tab">
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
            <p><strong>订阅主题:</strong> <code>{{ aircraftOsdTopic || 'thing/product/' + getAircraftSn() + '/osd' }}</code></p>
            <p><strong>等待数据中...</strong> 请确保设备已连接并正在发送OSD数据</p>
            <div style="margin-top: 15px;">
              <el-button 
                type="primary" 
                size="small" 
                :loading="subscribing"
                @click="subscribeToAircraftOSD"
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
            </div>
            
            <!-- 调试信息 -->
            <div style="margin-top: 15px; padding: 10px; background: var(--el-color-info-light-9); border-radius: 4px; font-size: 12px;">
              <div><strong>调试信息:</strong></div>
              <div>MQTT连接: {{ mqttProxyStore.isConnected ? '✅ 已连接' : '❌ 未连接' }}</div>
              <div>消息历史数量: {{ mqttProxyStore.messageHistory.length }}</div>
              <div>飞机OSD主题: <code>{{ aircraftOsdTopic }}</code></div>
              <div>机场序列号: <code>{{ getAirportSn() }}</code></div>
              <div>飞机序列号: <code>{{ getAircraftSn() }}</code></div>
              <div>已接收数据: {{ hasReceivedData ? '✅ 是' : '❌ 否' }}</div>
              <div v-if="mqttProxyStore.messageHistory.length > 0">
                <div>最新消息主题: <code>{{ mqttProxyStore.messageHistory[0]?.topic }}</code></div>
                <div>最新消息时间: {{ new Date(mqttProxyStore.messageHistory[0]?.timestamp).toLocaleString() }}</div>
              </div>
            </div>
          </div>
        </template>
      </el-alert>
    </div>

    <!-- 飞机信息展示 -->
    <div v-if="hasReceivedData" class="data-container">
      <!-- 飞机信息 -->
      <el-row :gutter="8" class="data-grid">
        <el-col :span="12">
          <h3 style="margin-bottom: 15px; color: var(--el-text-color-primary); border-left: 4px solid var(--el-color-success); padding-left: 10px;">
            ✈️ 飞机信息
          </h3>
        </el-col>
        <el-col :span="12" style="text-align: right;">
          <div style="margin-bottom: 15px;">
            <div style="margin-bottom: 5px;">
              <span style="color: var(--el-text-color-regular); font-size: 14px;">飞机序列号: </span>
              <span style="color: var(--el-text-color-primary); font-weight: 500;">{{ getAircraftSn() }}</span>
            </div>
            <div style="margin-bottom: 5px;">
              <span style="color: var(--el-text-color-regular); font-size: 12px;">最新刷新: </span>
              <span style="color: var(--el-color-success); font-size: 12px; font-weight: 500;">
                {{ lastUpdateTime ? lastUpdateTime.toLocaleString() : '未知' }}
              </span>
            </div>
            <!-- 调试信息 -->
            <div style="font-size: 12px; color: #999; margin-top: 2px;">
              调试: aircraftOsdData={{ aircraftOsdData?.gateway }}, route.aircraftSn={{ getAircraftSn() }}
            </div>
          </div>
        </el-col>
      </el-row>

      <!-- 飞机基本信息 - 所有卡片放在一行 -->
      <el-row :gutter="2" class="data-grid">
        <!-- 飞机状态与飞行参数 -->
        <el-col :span="6">
          <el-card class="info-card" shadow="hover" @click="openFlightMapModal" style="cursor: pointer; margin-bottom: 4px;">
            <template #header>
              <div class="card-header" style="padding: 10px 16px; font-size: 16px;">
                <el-icon><Location /></el-icon>
                <span>飞机状态与飞行参数</span>
                <el-tag size="small" type="success">实时</el-tag>
              </div>
            </template>
            <div class="info-content" style="padding: 8px;">
              <!-- 飞机状态 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);">
                <div style="font-weight: 600; color: var(--el-color-primary); margin-bottom: 4px; font-size: 15px;">✈️ 飞机状态</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">序列号:</span>
                  <span class="value" style="font-size: 11px;">{{ getAircraftSn() }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">状态:</span>
                  <el-tag :type="getModeCodeType(aircraftOsdData?.data?.mode_code)" size="small" style="font-size: 12px;">
                    {{ getModeCodeText(aircraftOsdData?.data?.mode_code) }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">起落架:</span>
                  <el-tag :type="getGearType(aircraftOsdData?.data?.gear)" size="small" style="font-size: 12px;">
                    {{ getGearText(aircraftOsdData?.data?.gear) }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">夜航灯:</span>
                  <el-tag :type="getNightLightsType(aircraftOsdData?.data?.night_lights_state)" size="small" style="font-size: 12px;">
                    {{ getNightLightsText(aircraftOsdData?.data?.night_lights_state) }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">RID:</span>
                  <el-tag :type="aircraftOsdData?.data?.rid_state ? 'success' : 'warning'" size="small" style="font-size: 12px;">
                    {{ aircraftOsdData?.data?.rid_state ? '开启' : '关闭' }}
                  </el-tag>
                </div>
              </div>
              
              <!-- 飞行参数 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);">
                <div style="font-weight: 600; color: var(--el-color-success); margin-bottom: 4px; font-size: 15px;">📍 飞行参数</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">纬度:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.latitude?.toFixed(6) || '未知' }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">经度:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.longitude?.toFixed(6) || '未知' }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">高度:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.height?.toFixed(1) || '未知' }}m</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">航向:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.attitude_head?.toFixed(1) || '未知' }}°</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">俯仰:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.attitude_pitch?.toFixed(1) || '未知' }}°</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">横滚:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.attitude_roll?.toFixed(1) || '未知' }}°</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">水平速度:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.horizontal_speed?.toFixed(1) || '未知' }}m/s</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">垂直速度:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.vertical_speed || '未知' }}m/s</span>
                </div>
              </div>
              
              <!-- 环境信息 -->
              <div>
                <div style="font-weight: 600; color: var(--el-color-info); margin-bottom: 4px; font-size: 15px;">🌤️ 环境信息</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">风速:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.wind_speed || '未知' }}m/s</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">风向:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.wind_direction || '未知' }}°</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">海拔:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.elevation || '未知' }}m</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">国家:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.country || '未知' }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <!-- 电池与存储信息 -->
        <el-col :span="6">
          <el-card class="info-card" shadow="hover" style="margin-bottom: 4px;">
            <template #header>
              <div class="card-header" style="padding: 10px 16px; font-size: 16px;">
                <el-icon><Lightning /></el-icon>
                <span>电池与存储信息</span>
                <el-tag size="small" type="warning">电量</el-tag>
              </div>
            </template>
            <div class="info-content" style="padding: 8px;">
              <!-- 电池信息 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);">
                <div style="font-weight: 600; color: var(--el-color-warning); margin-bottom: 4px; font-size: 15px;">🔋 电池信息</div>
                
                <!-- 电池总体信息 -->
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">总电量:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.battery?.capacity_percent || '未知' }}%</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">剩余飞行:</span>
                  <span class="value" style="font-size: 11px;">{{ formatDuration(aircraftOsdData?.data?.battery?.remain_flight_time) || '未知' }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">返航电量:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.battery?.return_home_power || '未知' }}%</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">降落电量:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.battery?.landing_power || '未知' }}%</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">低电量告警:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.serious_low_battery_warning_threshold || '未知' }}%</span>
                </div>

                <!-- 电池详细信息 -->
                <div v-if="aircraftOsdData?.data?.battery?.batteries && aircraftOsdData.data.battery.batteries.length > 0">
                  <div v-for="(battery, index) in aircraftOsdData.data.battery.batteries" :key="index" 
                       style="margin-top: 6px; padding: 4px; background: var(--el-color-warning-light-9); border-radius: 4px; border-left: 3px solid var(--el-color-warning);">
                    <div style="font-weight: 500; color: var(--el-color-warning-dark-2); margin-bottom: 3px; font-size: 12px;">
                      电池 {{ battery.index !== undefined ? battery.index : index + 1 }}
                    </div>
                    <div class="info-item" style="margin-bottom: 2px; font-size: 11px;">
                      <span class="label" style="font-size: 10px;">电量:</span>
                      <span class="value" style="font-size: 10px;">{{ battery.capacity_percent || '未知' }}%</span>
                    </div>
                    <div class="info-item" style="margin-bottom: 2px; font-size: 11px;">
                      <span class="label" style="font-size: 10px;">电压:</span>
                      <span class="value" style="font-size: 10px;">{{ battery.voltage || '未知' }}mV</span>
                    </div>
                    <div class="info-item" style="margin-bottom: 2px; font-size: 11px;">
                      <span class="label" style="font-size: 10px;">温度:</span>
                      <span class="value" style="font-size: 10px;">{{ battery.temperature ? battery.temperature.toFixed(1) : '未知' }}°C</span>
                    </div>
                    <div class="info-item" style="margin-bottom: 2px; font-size: 11px;">
                      <span class="label" style="font-size: 10px;">循环次数:</span>
                      <span class="value" style="font-size: 10px;">{{ battery.loop_times || '未知' }}</span>
                    </div>
                    <div class="info-item" style="margin-bottom: 2px; font-size: 11px;">
                      <span class="label" style="font-size: 10px;">序列号:</span>
                      <span class="value" style="font-size: 10px;">{{ battery.sn || '未知' }}</span>
                    </div>
                    <div class="info-item" style="margin-bottom: 2px; font-size: 11px;">
                      <span class="label" style="font-size: 10px;">类型:</span>
                      <span class="value" style="font-size: 10px;">{{ getBatteryTypeText(battery.type) }}</span>
                    </div>
                    <div class="info-item" style="margin-bottom: 2px; font-size: 11px;">
                      <span class="label" style="font-size: 10px;">子类型:</span>
                      <span class="value" style="font-size: 10px;">{{ getBatterySubTypeText(battery.sub_type) }}</span>
                    </div>
                    <div class="info-item" style="margin-bottom: 2px; font-size: 11px;">
                      <span class="label" style="font-size: 10px;">固件版本:</span>
                      <span class="value" style="font-size: 10px;">{{ battery.firmware_version || '未知' }}</span>
                    </div>
                    <div class="info-item" style="margin-bottom: 2px; font-size: 11px;">
                      <span class="label" style="font-size: 10px;">高电压存储:</span>
                      <span class="value" style="font-size: 10px;">{{ battery.high_voltage_storage_days || '未知' }}天</span>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- 存储信息 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);">
                <div style="font-weight: 600; color: var(--el-color-info); margin-bottom: 4px; font-size: 15px;">💾 存储信息</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">总容量:</span>
                  <span class="value" style="font-size: 11px;">{{ formatStorageSize(aircraftOsdData?.data?.storage?.total) || '未知' }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">已使用:</span>
                  <span class="value" style="font-size: 11px;">{{ formatStorageSize(aircraftOsdData?.data?.storage?.used) || '未知' }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">使用率:</span>
                  <span class="value" style="font-size: 11px;">{{ getStorageUsagePercent() }}%</span>
                </div>
              </div>
              
              <!-- 飞行统计 -->
              <div>
                <div style="font-weight: 600; color: var(--el-color-success); margin-bottom: 4px; font-size: 15px;">📊 飞行统计</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">总飞行时间:</span>
                  <span class="value" style="font-size: 11px;">{{ formatFlightTime(aircraftOsdData?.data?.total_flight_time) || '未知' }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">总飞行距离:</span>
                  <span class="value" style="font-size: 11px;">{{ formatDistance(aircraftOsdData?.data?.total_flight_distance) || '未知' }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <!-- 相机与载荷信息 -->
        <el-col :span="6">
          <el-card class="info-card" shadow="hover" style="margin-bottom: 4px;">
            <template #header>
              <div class="card-header" style="padding: 10px 16px; font-size: 16px;">
                <el-icon><Camera /></el-icon>
                <span>相机与载荷信息</span>
                <el-tag size="small" type="primary">实时</el-tag>
              </div>
            </template>
            <div class="info-content" style="padding: 8px;">
              <!-- 相机信息 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);">
                <div style="font-weight: 600; color: var(--el-color-primary); margin-bottom: 4px; font-size: 15px;">📷 相机信息</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">相机模式:</span>
                  <el-tag :type="getCameraModeType(aircraftOsdData?.data?.cameras?.[0]?.camera_mode)" size="small" style="font-size: 10px;">
                    {{ getCameraModeText(aircraftOsdData?.data?.cameras?.[0]?.camera_mode) }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">录制状态:</span>
                  <el-tag :type="aircraftOsdData?.data?.cameras?.[0]?.recording_state === 1 ? 'danger' : 'info'" size="small" style="font-size: 10px;">
                    {{ aircraftOsdData?.data?.cameras?.[0]?.recording_state === 1 ? '录制中' : '未录制' }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">拍照状态:</span>
                  <el-tag :type="aircraftOsdData?.data?.cameras?.[0]?.photo_state === 1 ? 'success' : 'info'" size="small" style="font-size: 10px;">
                    {{ aircraftOsdData?.data?.cameras?.[0]?.photo_state === 1 ? '拍照中' : '待机' }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">剩余照片:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.cameras?.[0]?.remain_photo_num || '未知' }}张</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">剩余录制:</span>
                  <span class="value" style="font-size: 11px;">{{ formatDuration(aircraftOsdData?.data?.cameras?.[0]?.remain_record_duration) }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">变焦倍数:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.cameras?.[0]?.zoom_factor?.toFixed(1) || '未知' }}x</span>
                </div>
              </div>
              
              <!-- 云台信息 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);" v-if="aircraftOsdData?.data?.['99-0-0']">
                <div style="font-weight: 600; color: var(--el-color-primary); margin-bottom: 4px; font-size: 15px;">🎯 云台信息</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">载荷索引:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.['99-0-0']?.payload_index || '--' }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">云台俯仰:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.['99-0-0']?.gimbal_pitch?.toFixed(1) || '--' }}°</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">云台横滚:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.['99-0-0']?.gimbal_roll?.toFixed(1) || '--' }}°</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">云台偏航:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.['99-0-0']?.gimbal_yaw?.toFixed(1) || '--' }}°</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">变焦倍数:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.['99-0-0']?.zoom_factor?.toFixed(2) || '--' }}x</span>
                </div>
              </div>
              
              <!-- 热成像信息 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);" v-if="aircraftOsdData?.data?.['99-0-0']">
                <div style="font-weight: 600; color: var(--el-color-warning); margin-bottom: 4px; font-size: 15px;">🌡️ 热成像信息</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">最高温度:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.['99-0-0']?.thermal_global_temperature_max?.toFixed(1) || '--' }}°C</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">最低温度:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.['99-0-0']?.thermal_global_temperature_min?.toFixed(1) || '--' }}°C</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">增益模式:</span>
                  <el-tag :type="getThermalGainType(aircraftOsdData?.data?.['99-0-0']?.thermal_gain_mode)" size="small" style="font-size: 10px;">
                    {{ getThermalGainText(aircraftOsdData?.data?.['99-0-0']?.thermal_gain_mode) }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">调色板:</span>
                  <el-tag :type="getThermalPaletteType(aircraftOsdData?.data?.['99-0-0']?.thermal_current_palette_style)" size="small" style="font-size: 10px;">
                    {{ getThermalPaletteText(aircraftOsdData?.data?.['99-0-0']?.thermal_current_palette_style) }}
                  </el-tag>
                </div>
              </div>
              
              <!-- 测量信息 -->
              <div v-if="aircraftOsdData?.data?.['99-0-0']">
                <div style="font-weight: 600; color: var(--el-color-success); margin-bottom: 4px; font-size: 15px;">📏 测量信息</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">目标距离:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.['99-0-0']?.measure_target_distance?.toFixed(1) || '--' }}m</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">目标高度:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.['99-0-0']?.measure_target_altitude?.toFixed(1) || '--' }}m</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">测量状态:</span>
                  <el-tag :type="aircraftOsdData?.data?.['99-0-0']?.measure_target_error_state === 0 ? 'success' : 'danger'" size="small" style="font-size: 10px;">
                    {{ aircraftOsdData?.data?.['99-0-0']?.measure_target_error_state === 0 ? '正常' : '错误' }}
                  </el-tag>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
        <!-- 限制与安全信息 -->
        <el-col :span="6">
          <el-card class="info-card" shadow="hover" style="margin-bottom: 4px;">
            <template #header>
              <div class="card-header" style="padding: 10px 16px; font-size: 16px;">
                <el-icon><Cloudy /></el-icon>
                <span>限制与安全信息</span>
                <el-tag size="small" type="info">安全</el-tag>
              </div>
            </template>
            <div class="info-content" style="padding: 8px;">
              <!-- 飞行限制 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);">
                <div style="font-weight: 600; color: var(--el-color-warning); margin-bottom: 4px; font-size: 15px;">🚫 飞行限制</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">高度限制:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.height_limit || '--' }}m</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">距离限制:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.distance_limit_status?.distance_limit || '--' }}m</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">返航高度:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.rth_altitude || '--' }}m</span>
                </div>
              </div>
              
              <!-- 安全设置 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);">
                <div style="font-weight: 600; color: var(--el-color-danger); margin-bottom: 4px; font-size: 15px;">⚠️ 安全设置</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">失控动作:</span>
                  <el-tag :type="getRcLostActionType(aircraftOsdData?.data?.rc_lost_action)" size="small" style="font-size: 10px;">
                    {{ getRcLostActionText(aircraftOsdData?.data?.rc_lost_action) }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">接近高度限制:</span>
                  <el-tag :type="aircraftOsdData?.data?.is_near_height_limit ? 'warning' : 'success'" size="small" style="font-size: 10px;">
                    {{ aircraftOsdData?.data?.is_near_height_limit ? '是' : '否' }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">接近区域限制:</span>
                  <el-tag :type="aircraftOsdData?.data?.is_near_area_limit ? 'warning' : 'success'" size="small" style="font-size: 10px;">
                    {{ aircraftOsdData?.data?.is_near_area_limit ? '是' : '否' }}
                  </el-tag>
                </div>
              </div>
              
              <!-- 避障系统 -->
              <div style="margin-bottom: 8px; padding-bottom: 6px; border-bottom: 1px solid var(--el-border-color-lighter);">
                <div style="font-weight: 600; color: var(--el-color-info); margin-bottom: 4px; font-size: 15px;">🛡️ 避障系统</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">上方避障:</span>
                  <el-tag :type="aircraftOsdData?.data?.obstacle_avoidance?.upside ? 'success' : 'warning'" size="small" style="font-size: 10px;">
                    {{ aircraftOsdData?.data?.obstacle_avoidance?.upside ? '开启' : '关闭' }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">水平避障:</span>
                  <el-tag :type="aircraftOsdData?.data?.obstacle_avoidance?.horizon ? 'success' : 'warning'" size="small" style="font-size: 10px;">
                    {{ aircraftOsdData?.data?.obstacle_avoidance?.horizon ? '开启' : '关闭' }}
                  </el-tag>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">下方避障:</span>
                  <el-tag :type="aircraftOsdData?.data?.obstacle_avoidance?.downside ? 'success' : 'warning'" size="small" style="font-size: 10px;">
                    {{ aircraftOsdData?.data?.obstacle_avoidance?.downside ? '开启' : '关闭' }}
                  </el-tag>
                </div>
              </div>
              
              <!-- 定位信息 -->
              <div>
                <div style="font-weight: 600; color: var(--el-color-success); margin-bottom: 4px; font-size: 15px;">📡 定位信息</div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">GPS数量:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.position_state?.gps_number || '--' }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">RTK数量:</span>
                  <span class="value" style="font-size: 11px;">{{ aircraftOsdData?.data?.position_state?.rtk_number || '--' }}</span>
                </div>
                <div class="info-item" style="margin-bottom: 3px; font-size: 12px;">
                  <span class="label" style="font-size: 11px;">定位质量:</span>
                  <el-tag :type="aircraftOsdData?.data?.position_state?.quality >= 4 ? 'success' : 'warning'" size="small" style="font-size: 10px;">
                    {{ aircraftOsdData?.data?.position_state?.quality || '--' }}
                  </el-tag>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 等待数据状态 -->
    <div v-else class="waiting-data" style="margin: 20px 0; padding: 20px; text-align: center; background: var(--el-color-warning-light-9); border-radius: 8px;">
      <el-icon style="color: var(--el-color-warning); font-size: 48px; margin-bottom: 16px;"><Loading /></el-icon>
      <h3 style="color: var(--el-color-warning-dark-2); margin-bottom: 8px;">等待接收飞机OSD数据...</h3>
      <p style="color: var(--el-text-color-regular); font-size: 14px;">请确保MQTT连接正常，设备正在上报数据</p>
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
        @click="handleAircraftOSDMessage"
        title="手动处理飞机OSD消息"
      >
        🔄
      </el-button>
    </div>

    <!-- 数据调试按钮 -->
    <div style="position: fixed; top: 70%; right: 20px; transform: translateY(-50%); z-index: 1000;">
      <el-button 
        type="info" 
        circle 
        size="large"
        @click="debugData"
        title="调试数据状态"
      >
        🐛
      </el-button>
    </div>

    <!-- 飞行轨迹组件 -->
    <FlightTrajectoryMap
      v-model="showFlightMapModal"
      :aircraft-data="aircraftOsdData?.data"
      :aircraft-sn="getAircraftSn()"
    />

    <!-- 原始数据抽屉 -->
    <el-drawer
      v-model="showRawDataDrawer"
      title="飞机原始数据展示"
      direction="rtl"
      size="60%"
      :before-close="handleDrawerClose"
    >
      <div class="info-content">
        <el-tabs type="border-card">
          <el-tab-pane label="飞机OSD数据" name="aircraft-osd">
            <pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; max-height: 500px; overflow-y: auto; font-size: 12px;">{{ JSON.stringify(aircraftOsdData, null, 2) }}</pre>
          </el-tab-pane>
          <el-tab-pane label="数据更新信息" name="updates">
            <div style="padding: 10px;">
              <div style="margin-bottom: 15px;">
                <h4 style="color: var(--el-color-primary); margin-bottom: 10px;">📊 数据更新状态</h4>
                <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 10px;">
                  <div style="padding: 8px; background: var(--el-color-info-light-9); border-radius: 4px;">
                    <strong>飞机OSD数据:</strong> 
                    <span v-if="aircraftOsdData?.timestamp" :style="isDataStale(aircraftOsdData.timestamp) ? 'color: var(--el-color-warning);' : 'color: var(--el-color-success);'">
                      {{ formatRelativeTime(aircraftOsdData.timestamp) }}
                      <span v-if="isDataStale(aircraftOsdData.timestamp)" style="margin-left: 4px;">⚠️</span>
                    </span>
                    <span v-else style="color: var(--el-color-warning);">等待数据</span>
                  </div>
                </div>
              </div>
              <div>
                <h4 style="color: var(--el-color-primary); margin-bottom: 10px;">🔄 最近更新字段</h4>
                <div v-if="aircraftOsdData?.data" style="background: var(--el-color-success-light-9); padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                  <strong>飞机OSD数据字段:</strong> {{ Object.keys(aircraftOsdData.data).join(', ') }}
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Location, Lightning, Cloudy, Loading, DataLine, Sunny, Camera } from '@element-plus/icons-vue'
import { useMqttProxyStore } from '@/stores/mqtt-proxy'
import { useDeviceStore } from '@/stores/device'
import { COORDINATE_CONVERTER } from '@/config/amap'
import FlightTrajectoryMap from '@/components/FlightTrajectoryMap.vue'

// Props
const props = defineProps({
  deviceSn: String,
  aircraftSn: String
})

// Emits
const emit = defineEmits(['dataReceived', 'subscriptionStatus'])

// Stores
const mqttProxyStore = useMqttProxyStore()
const deviceStore = useDeviceStore()

// 响应式数据
const aircraftOsdData = ref(null)
const hasReceivedData = ref(false)
const subscriptionStatus = ref('未订阅')
const subscribing = ref(false)
const showRawDataDrawer = ref(false)
const showFlightMapModal = ref(false)
const lastUpdateTime = ref(null)

// 计算属性
const aircraftOsdTopic = computed(() => {
  // 使用飞机序列号来订阅飞机OSD数据
  const aircraftSn = getAircraftSn()
  return aircraftSn && aircraftSn !== '--' ? `thing/product/${aircraftSn}/osd` : null
})

// 获取飞机序列号
const getAircraftSn = () => {
  return props.aircraftSn || props.deviceSn || '--'
}

// 获取机场序列号
const getAirportSn = () => {
  return deviceStore.currentDevice?.airport_sn || deviceStore.currentDeviceSn || '--'
}

// 获取飞行模式文本
const getFlightModeText = (mode) => {
  const modeMap = {
    0: '手动模式',
    1: '自动模式',
    2: '返航模式',
    3: '悬停模式',
    4: '降落模式'
  }
  return modeMap[mode] || '未知模式'
}

// 获取飞机状态码类型
const getModeCodeType = (code) => {
  const typeMap = {
    0: 'success',   // 空闲中
    1: 'warning',   // 现场调试
    2: 'warning',   // 远程调试
    3: 'info',      // 固件升级中
    4: 'success',   // 作业中
    5: 'warning'    // 待标定
  }
  return typeMap[code] || 'info'
}

// 获取飞机状态码文本
const getModeCodeText = (code) => {
  const textMap = {
    0: '空闲中',
    1: '现场调试',
    2: '远程调试',
    3: '固件升级中',
    4: '作业中',
    5: '待标定'
  }
  return textMap[code] || '未知状态'
}

// 格式化存储大小
const formatStorageSize = (bytes) => {
  if (!bytes) return '--'
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + sizes[i]
}

// 获取存储使用率
const getStorageUsagePercent = () => {
  const total = aircraftOsdData.value?.data?.storage?.total
  const used = aircraftOsdData.value?.data?.storage?.used
  if (!total || !used) return '未知'
  return ((used / total) * 100).toFixed(1)
}

// 格式化飞行时间
const formatFlightTime = (seconds) => {
  if (!seconds) return '未知'
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)
  return `${hours}h ${minutes}m ${secs}s`
}

// 获取GCJ02纬度
const getGcj02Latitude = () => {
  if (!aircraftOsdData.value?.data?.latitude || !aircraftOsdData.value?.data?.longitude) return null
  
  try {
    const [gcjLng, gcjLat] = COORDINATE_CONVERTER.wgs84ToGcj02(
      aircraftOsdData.value.data.longitude, 
      aircraftOsdData.value.data.latitude
    )
    return gcjLat
  } catch (error) {
    console.error('❌ 坐标转换失败:', error)
    return null
  }
}

// 获取GCJ02经度
const getGcj02Longitude = () => {
  if (!aircraftOsdData.value?.data?.latitude || !aircraftOsdData.value?.data?.longitude) return null
  
  try {
    const [gcjLng, gcjLat] = COORDINATE_CONVERTER.wgs84ToGcj02(
      aircraftOsdData.value.data.longitude, 
      aircraftOsdData.value.data.latitude
    )
    return gcjLng
  } catch (error) {
    console.error('❌ 坐标转换失败:', error)
    return null
  }
}

// 获取起落架类型
const getGearType = (gear) => {
  if (gear === undefined || gear === null) return 'info'
  return gear === 1 ? 'success' : 'info'
}

// 获取起落架文本
const getGearText = (gear) => {
  if (gear === undefined || gear === null) return '未知'
  return gear === 1 ? '收起' : '放下'
}

// 获取夜航灯类型
const getNightLightsType = (state) => {
  if (state === undefined || state === null) return 'info'
  return state === 1 ? 'success' : 'info'
}

// 获取夜航灯文本
const getNightLightsText = (state) => {
  if (state === undefined || state === null) return '未知'
  return state === 1 ? '开启' : '关闭'
}

// 获取录制状态类型
const getRecordingStateType = (state) => {
  if (state === undefined || state === null) return 'info'
  return state === 1 ? 'danger' : 'info'
}

// 获取录制状态文本
const getRecordingStateText = (state) => {
  if (state === undefined || state === null) return '未知'
  return state === 1 ? '录制中' : '未录制'
}

// 热成像增益模式相关方法
const getThermalGainType = (mode) => {
  const typeMap = {
    0: 'info',      // 自动
    1: 'success',   // 手动
    2: 'warning'    // 其他
  }
  return typeMap[mode] || 'info'
}

const getThermalGainText = (mode) => {
  const textMap = {
    0: '自动',
    1: '手动',
    2: '其他'
  }
  return textMap[mode] || '未知'
}

// 热成像调色板相关方法
const getThermalPaletteType = (style) => {
  const typeMap = {
    0: 'info',      // 默认
    1: 'success',   // 彩虹
    2: 'warning',   // 铁红
    3: 'primary'    // 其他
  }
  return typeMap[style] || 'info'
}

const getThermalPaletteText = (style) => {
  const textMap = {
    0: '默认',
    1: '彩虹',
    2: '铁红',
    3: '其他'
  }
  return textMap[style] || '未知'
}

// 相机模式相关方法
const getCameraModeType = (mode) => {
  const typeMap = {
    0: 'info',      // 拍照模式
    1: 'success',   // 录像模式
    2: 'warning'    // 其他模式
  }
  return typeMap[mode] || 'info'
}

const getCameraModeText = (mode) => {
  const textMap = {
    0: '拍照模式',
    1: '录像模式',
    2: '其他模式'
  }
  return textMap[mode] || '未知'
}

// 失控动作相关方法
const getRcLostActionType = (action) => {
  const typeMap = {
    0: 'info',      // 悬停
    1: 'warning',   // 降落
    2: 'danger'     // 返航
  }
  return typeMap[action] || 'info'
}

const getRcLostActionText = (action) => {
  const textMap = {
    0: '悬停',
    1: '降落',
    2: '返航'
  }
  return textMap[action] || '未知'
}

// 格式化持续时间（秒转换为时分秒）
const formatDuration = (seconds) => {
  if (!seconds || seconds === 0) return '0秒'
  
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)
  
  let result = ''
  if (hours > 0) result += `${hours}小时`
  if (minutes > 0) result += `${minutes}分钟`
  if (secs > 0) result += `${secs}秒`
  
  return result || '0秒'
}



// 订阅飞机OSD数据
const subscribeToAircraftOSD = async () => {
  if (!mqttProxyStore.isConnected) {
    ElMessage.warning('MQTT未连接，无法订阅飞机OSD数据')
    subscriptionStatus.value = 'MQTT未连接'
    return
  }

  if (!aircraftOsdTopic.value) {
    ElMessage.warning('飞机序列号为空，无法订阅飞机OSD数据')
    subscriptionStatus.value = '飞机序列号为空'
    return
  }

  subscribing.value = true
  subscriptionStatus.value = '订阅中...'

  try {
    await mqttProxyStore.subscribeToTopics(aircraftOsdTopic.value, 1)
    subscriptionStatus.value = '已订阅'
    ElMessage.success(`已订阅飞机OSD主题: ${aircraftOsdTopic.value}`)
    console.log('飞机OSD订阅成功:', aircraftOsdTopic.value)
    emit('subscriptionStatus', 'success')
  } catch (error) {
    subscriptionStatus.value = '订阅失败'
    console.error('订阅飞机OSD主题失败:', error)
    ElMessage.error('订阅飞机OSD主题失败')
    emit('subscriptionStatus', 'error')
  } finally {
    subscribing.value = false
  }
}

// 处理飞机OSD消息
const handleAircraftOSDMessage = () => {
  try {
    console.log('🔍 检查飞机OSD消息:', {
      aircraftOsdTopic: aircraftOsdTopic.value,
      messageHistoryLength: mqttProxyStore.messageHistory.length,
      hasReceivedData: hasReceivedData.value,
      currentTimestamp: aircraftOsdData.value?.timestamp,
      allTopics: mqttProxyStore.messageHistory.map(msg => msg.topic)
    })
    
    if (!aircraftOsdTopic.value) {
      console.log('❌ 飞机OSD主题为空')
      return
    }
    
    // 查找所有匹配主题的消息
    const matchingMessages = mqttProxyStore.messageHistory.filter(msg => 
      msg.topic === aircraftOsdTopic.value
    )
    
    console.log('📨 匹配主题的消息数量:', matchingMessages.length)
    console.log('📨 匹配的消息:', matchingMessages)
    console.log('📨 期望的主题:', aircraftOsdTopic.value)
    console.log('📨 实际收到的主题:', mqttProxyStore.messageHistory.map(msg => msg.topic))
    
    if (matchingMessages.length > 0) {
      // 获取最新的消息
      const latestMessage = matchingMessages[0]
      console.log('📨 最新消息:', latestMessage)
      console.log('📨 最新消息数据:', latestMessage.data)
      
      // 检查是否是新消息
      const isNewMessage = !aircraftOsdData.value || 
                          latestMessage.timestamp > aircraftOsdData.value.timestamp
      
      console.log('🆕 是否为新消息:', isNewMessage, {
        currentTimestamp: aircraftOsdData.value?.timestamp,
        newTimestamp: latestMessage.timestamp
      })
      
      if (isNewMessage) {
        console.log('✅ 处理新的飞机OSD数据:', latestMessage)
        
        // 解析payload数据
        let parsedData = null
        try {
          if (latestMessage.payload) {
            parsedData = JSON.parse(latestMessage.payload)
            console.log('✅ 解析的payload数据:', parsedData)
          }
        } catch (error) {
          console.error('❌ 解析payload失败:', error)
        }
        
        // 创建包含解析数据的消息对象
        const processedMessage = {
          ...latestMessage,
          data: parsedData?.data || latestMessage.data
        }
        
        aircraftOsdData.value = processedMessage
        hasReceivedData.value = true
        subscriptionStatus.value = '已接收数据'
        lastUpdateTime.value = new Date()
        emit('dataReceived', processedMessage)
        
        // 不再显示弹出提示，只在页面上显示刷新时间
        console.log('✅ 飞机OSD数据已更新:', aircraftOsdData.value)
        console.log('✅ 飞机OSD数据内容:', {
          mode_code: aircraftOsdData.value?.data?.mode_code,
          battery_percent: aircraftOsdData.value?.data?.battery?.capacity_percent,
          latitude: aircraftOsdData.value?.data?.latitude,
          longitude: aircraftOsdData.value?.data?.longitude,
          height: aircraftOsdData.value?.data?.height,
          gear: aircraftOsdData.value?.data?.gear,
          night_lights_state: aircraftOsdData.value?.data?.night_lights_state
        })
        
        // 输出高精度坐标用于调试
        if (aircraftOsdData.value?.data?.latitude && aircraftOsdData.value?.data?.longitude) {
          console.log('📍 高精度坐标:', {
            latitude: aircraftOsdData.value.data.latitude,
            longitude: aircraftOsdData.value.data.longitude,
            latitude_precise: aircraftOsdData.value.data.latitude.toFixed(10),
            longitude_precise: aircraftOsdData.value.data.longitude.toFixed(10)
          })
        }
      } else {
        console.log('⏭️ 消息已处理过，跳过')
      }
    } else {
      console.log('📭 没有找到匹配主题的消息')
      console.log('📭 所有消息主题:', mqttProxyStore.messageHistory.map(msg => msg.topic))
      console.log('📭 期望的主题:', aircraftOsdTopic.value)
    }
  } catch (error) {
    console.error('❌ 处理飞机OSD消息失败:', error)
  }
}

// 消息监听器 - 使用 watch 监听消息历史变化
watch(
  () => mqttProxyStore.messageHistory,
  (newMessages, oldMessages) => {
    if (newMessages && newMessages.length > 0) {
      console.log('📨 消息历史发生变化，检查飞机OSD消息')
      handleAircraftOSDMessage()
    }
  },
  { deep: true, immediate: true }
)

// 刷新状态
const refreshStatus = () => {
  console.log('🔄 刷新订阅状态')
  console.log('📊 当前状态:', {
    mqttConnected: mqttProxyStore.isConnected,
    messageHistoryLength: mqttProxyStore.messageHistory.length,
    aircraftOsdTopic: aircraftOsdTopic.value,
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
        msg.topic === aircraftOsdTopic.value
      )
      
      if (matchingMessages.length > 0) {
        console.log('🎯 找到匹配的消息:', matchingMessages)
        subscriptionStatus.value = `找到 ${matchingMessages.length} 条消息，正在处理...`
        // 立即处理消息
        handleAircraftOSDMessage()
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
  console.log('飞机信息组件挂载，开始订阅数据')
  subscribeToAircraftOSD()
})

// 组件卸载
onUnmounted(() => {
  console.log('飞机信息组件卸载')
  
  // 取消订阅
  if (aircraftOsdTopic.value) {
    mqttProxyStore.unsubscribeTopics(aircraftOsdTopic.value)
  }
})

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

// 抽屉关闭处理
const handleDrawerClose = (done) => {
  showRawDataDrawer.value = false
  done()
}

// 调试数据状态
const debugData = () => {
  console.log('🐛 调试数据状态:')
  console.log('📊 aircraftOsdData.value:', aircraftOsdData.value)
  console.log('📊 aircraftOsdData.value?.data:', aircraftOsdData.value?.data)
  console.log('📊 hasReceivedData.value:', hasReceivedData.value)
  console.log('📊 具体数据字段:')
  console.log('  - latitude:', aircraftOsdData.value?.data?.latitude)
  console.log('  - longitude:', aircraftOsdData.value?.data?.longitude)
  console.log('  - height:', aircraftOsdData.value?.data?.height)
  console.log('  - battery.capacity_percent:', aircraftOsdData.value?.data?.battery?.capacity_percent)
  console.log('  - mode_code:', aircraftOsdData.value?.data?.mode_code)
  console.log('  - gear:', aircraftOsdData.value?.data?.gear)
  console.log('  - wind_speed:', aircraftOsdData.value?.data?.wind_speed)
  console.log('  - wind_direction:', aircraftOsdData.value?.data?.wind_direction)
  console.log('  - elevation:', aircraftOsdData.value?.data?.elevation)
  console.log('  - country:', aircraftOsdData.value?.data?.country)
}

// 打开飞行地图模态框
const openFlightMapModal = () => {
  showFlightMapModal.value = true
}

// 格式化距离
const formatDistance = (meters) => {
  if (!meters) return '未知'
  if (meters >= 1000) {
    return `${(meters / 1000).toFixed(2)}km`
  }
  return `${meters.toFixed(2)}m`
}

// 获取电池类型文本
const getBatteryTypeText = (type) => {
  if (type === undefined || type === null) return '未知'
  const typeMap = {
    0: '标准电池',
    1: '高容量电池',
    2: '智能电池',
    3: '专业电池'
  }
  return typeMap[type] || `类型${type}`
}

// 获取电池子类型文本
const getBatterySubTypeText = (subType) => {
  if (subType === undefined || subType === null) return '未知'
  const subTypeMap = {
    0: '标准版',
    1: '增强版',
    2: '专业版',
    3: '旗舰版'
  }
  return subTypeMap[subType] || `子类型${subType}`
}
</script>

<style scoped>
.aircraft-info-tab {
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