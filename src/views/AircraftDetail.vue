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
          <div v-if="hasReceivedData && osdData" class="data-update-info">
            <el-tag type="info" size="small" style="margin-right: 8px;">
              最后更新: {{ formatTime(osdData.timestamp) }}
            </el-tag>
            <el-tag 
              :type="getDataStatusType()" 
              size="small" 
              style="margin-right: 8px;"
            >
              {{ getDataStatusText() }}
            </el-tag>
          </div>
          <el-button 
            v-if="hasReceivedData" 
            type="primary" 
            size="small" 
            @click="refreshData"
            style="margin-right: 8px;"
          >
            刷新数据
          </el-button>
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
        <!-- 连接状态提示 -->
        <div v-if="!mqttProxyStore.isConnected" class="connection-warning">
          <el-alert
            title="MQTT未连接"
            description="请先在仪表板页面连接MQTT服务器以获取实时数据"
            type="warning"
            show-icon
            :closable="false"
          />
        </div>

        <!-- 数据加载状态 -->
        <div v-if="mqttProxyStore.isConnected && !hasReceivedData" class="loading-state">
          <el-empty description="等待OSD数据...">
            <el-button type="primary" @click="subscribeToOSD">
              手动订阅OSD主题
            </el-button>
            <el-button type="default" @click="refreshData" style="margin-left: 10px;">
              刷新数据
            </el-button>
          </el-empty>
        </div>

        <!-- 实时数据展示 -->
        <div v-if="hasReceivedData" class="data-container">
          <!-- 机场信息 (左侧) -->
          <el-row :gutter="8" class="data-grid">
            <el-col :span="12">
              <h3 style="margin-bottom: 15px; color: var(--el-text-color-primary); border-left: 4px solid var(--el-color-primary); padding-left: 10px;">
                🏢 机场信息
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
          </el-row>

          <!-- 机场基本信息 -->
          <el-row :gutter="8" class="data-grid">
            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Location /></el-icon>
                    <span>机场位置</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">纬度:</span>
                    <span class="value">{{ osdData?.data?.latitude?.toFixed(6) || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">经度:</span>
                    <span class="value">{{ osdData?.data?.longitude?.toFixed(6) || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">高度:</span>
                    <span class="value">{{ osdData?.data?.height?.toFixed(2) || '--' }}m</span>
                  </div>
                  <div class="info-item">
                    <span class="label">航向:</span>
                    <span class="value">{{ osdData?.data?.heading?.toFixed(1) || '--' }}°</span>
                  </div>
                  <div class="info-item">
                    <span class="label">家点有效:</span>
                    <el-tag :type="osdData?.data?.home_position_is_valid === 1 ? 'success' : 'warning'" size="small">
                      {{ osdData?.data?.home_position_is_valid === 1 ? '有效' : '无效' }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Cloudy /></el-icon>
                    <span>环境监控</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">环境温度:</span>
                    <span class="value">{{ osdData?.data?.environment_temperature?.toFixed(1) || '--' }}°C</span>
                  </div>
                  <div class="info-item">
                    <span class="label">舱内温度:</span>
                    <span class="value">{{ osdData?.data?.temperature?.toFixed(1) || '--' }}°C</span>
                  </div>
                  <div class="info-item">
                    <span class="label">湿度:</span>
                    <span class="value">{{ osdData?.data?.humidity || '--' }}%</span>
                  </div>
                  <div class="info-item">
                    <span class="label">风速:</span>
                    <span class="value">{{ osdData?.data?.wind_speed || '--' }}m/s</span>
                  </div>
                  <div class="info-item">
                    <span class="label">降雨量:</span>
                    <el-tag :type="getRainfallType(osdData?.data?.rainfall)" size="small">
                      {{ getRainfallText(osdData?.data?.rainfall) }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Lightning /></el-icon>
                    <span>电源系统</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">供电电压:</span>
                    <span class="value">{{ osdData?.data?.electric_supply_voltage || '--' }}V</span>
                  </div>
                  <div class="info-item">
                    <span class="label">工作电压:</span>
                    <span class="value">{{ (osdData?.data?.working_voltage / 1000)?.toFixed(2) || '--' }}V</span>
                  </div>
                  <div class="info-item">
                    <span class="label">工作电流:</span>
                    <span class="value">{{ (osdData?.data?.working_current / 1000)?.toFixed(2) || '--' }}A</span>
                  </div>
                  <div class="info-item">
                    <span class="label">ACDC功率:</span>
                    <span class="value">{{ osdData?.data?.acdc_power_input?.toFixed(2) || '--' }}W</span>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Connection /></el-icon>
                    <span>网络通信</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">网络类型:</span>
                    <el-tag :type="getNetworkTypeColor(osdData?.data?.network_state?.type)" size="small">
                      {{ getNetworkTypeText(osdData?.data?.network_state?.type) }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">网络质量:</span>
                    <div style="display: flex; align-items: center; gap: 8px;">
                      <el-tag :type="getNetworkQualityColor(osdData?.data?.network_state?.quality)" size="small">
                        {{ getNetworkQualityText(osdData?.data?.network_state?.quality) }}
                      </el-tag>
                      <el-progress 
                        :percentage="(osdData?.data?.network_state?.quality || 0) * 20"
                        :color="getNetworkQualityProgressColor(osdData?.data?.network_state?.quality)"
                        :show-text="false"
                        :stroke-width="4"
                        style="flex: 1;"
                      />
                    </div>
                  </div>
                  <div class="info-item">
                    <span class="label">传输速率:</span>
                    <span class="value">{{ getNetworkRateText(osdData?.data?.network_state?.rate) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">空中回传:</span>
                    <el-tag :type="osdData?.data?.air_transfer_enable ? 'success' : 'info'" size="small">
                      {{ osdData?.data?.air_transfer_enable ? '开启' : '关闭' }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Setting /></el-icon>
                    <span>机场状态</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">机场模式:</span>
                    <el-tag :type="getModeCodeType(osdData?.data?.mode_code)" size="small">
                      {{ getModeCodeText(osdData?.data?.mode_code) }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">任务状态:</span>
                    <el-tag :type="getFlightTaskStepType(osdData?.data?.flighttask_step_code)" size="small">
                      {{ getFlightTaskStepText(osdData?.data?.flighttask_step_code) }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">在机库:</span>
                    <el-tag :type="osdData?.data?.drone_in_dock ? 'success' : 'info'" size="small">
                      {{ osdData?.data?.drone_in_dock ? '是' : '否' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">静音模式:</span>
                    <el-tag :type="osdData?.data?.silent_mode ? 'warning' : 'success'" size="small">
                      {{ osdData?.data?.silent_mode ? '开启' : '关闭' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">紧急停止:</span>
                    <el-tag :type="osdData?.data?.emergency_stop_state ? 'danger' : 'success'" size="small">
                      {{ osdData?.data?.emergency_stop_state ? '激活' : '正常' }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Folder /></el-icon>
                    <span>存储信息</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">总容量:</span>
                    <span class="value">{{ formatStorage(osdData?.data?.storage?.total) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">已使用:</span>
                    <span class="value">{{ formatStorage(osdData?.data?.storage?.used) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">使用率:</span>
                    <el-progress 
                      :percentage="getStorageUsage(osdData?.data?.storage)"
                      :color="getStorageUsageColor(osdData?.data?.storage)"
                      :show-text="true"
                      :stroke-width="6"
                    />
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Monitor /></el-icon>
                    <span>子设备信息</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">飞机序列号:</span>
                    <span class="value">{{ osdData?.data?.sub_device?.device_sn || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">设备型号:</span>
                    <span class="value">{{ osdData?.data?.sub_device?.device_model_key || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">在线状态:</span>
                    <el-tag :type="osdData?.data?.sub_device?.device_online_status ? 'success' : 'danger'" size="small">
                      {{ osdData?.data?.sub_device?.device_online_status ? '在线' : '离线' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">配对状态:</span>
                    <el-tag :type="osdData?.data?.sub_device?.device_paired ? 'success' : 'warning'" size="small">
                      {{ osdData?.data?.sub_device?.device_paired ? '已配对' : '未配对' }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Location /></el-icon>
                    <span>位置状态</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">校准状态:</span>
                    <el-tag :type="osdData?.data?.position_state?.is_calibration ? 'warning' : 'success'" size="small">
                      {{ osdData?.data?.position_state?.is_calibration ? '校准中' : '已校准' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">定位状态:</span>
                    <el-tag :type="getPositionFixedType(osdData?.data?.position_state?.is_fixed)" size="small">
                      {{ getPositionFixedText(osdData?.data?.position_state?.is_fixed) }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">GPS卫星数:</span>
                    <span class="value">{{ osdData?.data?.position_state?.gps_number || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">RTK卫星数:</span>
                    <span class="value">{{ osdData?.data?.position_state?.rtk_number || '--' }}</span>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>

          <!-- 飞机信息 (右侧) -->
          <el-row :gutter="8" class="data-grid" style="margin-top: 30px;">
            <el-col :span="12">
              <h3 style="margin-bottom: 15px; color: var(--el-text-color-primary); border-left: 4px solid var(--el-color-success); padding-left: 10px;">
                ✈️ 飞机信息
              </h3>
            </el-col>
          </el-row>

          <!-- 备降点和工作状态信息 -->
          <el-row :gutter="8" class="data-grid" style="margin-top: 15px;">
            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Location /></el-icon>
                    <span>备降点信息</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">备降点纬度:</span>
                    <span class="value">{{ osdData?.data?.alternate_land_point?.latitude?.toFixed(6) || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">备降点经度:</span>
                    <span class="value">{{ osdData?.data?.alternate_land_point?.longitude?.toFixed(6) || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">备降点高度:</span>
                    <span class="value">{{ osdData?.data?.alternate_land_point?.height?.toFixed(2) || '--' }}m</span>
                  </div>
                  <div class="info-item">
                    <span class="label">安全降落高度:</span>
                    <span class="value">{{ osdData?.data?.alternate_land_point?.safe_land_height || '--' }}m</span>
                  </div>
                  <div class="info-item">
                    <span class="label">配置状态:</span>
                    <el-tag :type="osdData?.data?.alternate_land_point?.is_configured ? 'success' : 'warning'" size="small">
                      {{ osdData?.data?.alternate_land_point?.is_configured ? '已配置' : '未配置' }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Setting /></el-icon>
                    <span>工作状态</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">任务编号:</span>
                    <span class="value">{{ osdData?.data?.job_number || '--' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">累计时间:</span>
                    <span class="value">{{ formatAccTime(osdData?.data?.acc_time) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">激活时间:</span>
                    <span class="value">{{ formatTimestamp(osdData?.data?.activation_time) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">部署模式:</span>
                    <el-tag :type="osdData?.data?.deployment_mode ? 'success' : 'info'" size="small">
                      {{ osdData?.data?.deployment_mode ? '已部署' : '未部署' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">任务状态:</span>
                    <el-tag :type="getFlightTaskStepType(osdData?.data?.flighttask_step_code)" size="small">
                      {{ getFlightTaskStepText(osdData?.data?.flighttask_step_code) }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">机场状态:</span>
                    <el-tag :type="getModeCodeType(osdData?.data?.mode_code)" size="small">
                      {{ getModeCodeText(osdData?.data?.mode_code) }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Lightning /></el-icon>
                    <span>电力供应</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">供电电压:</span>
                    <span class="value">{{ osdData?.data?.electric_supply_voltage || '--' }}V</span>
                  </div>
                  <div class="info-item">
                    <span class="label">工作电压:</span>
                    <span class="value">{{ (osdData?.data?.working_voltage / 1000)?.toFixed(2) || '--' }}V</span>
                  </div>
                  <div class="info-item">
                    <span class="label">工作电流:</span>
                    <span class="value">{{ (osdData?.data?.working_current / 1000)?.toFixed(2) || '--' }}A</span>
                  </div>
                  <div class="info-item">
                    <span class="label">ACDC功率:</span>
                    <span class="value">{{ osdData?.data?.acdc_power_input?.toFixed(2) || '--' }}W</span>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Lightning /></el-icon>
                    <span>备用电池</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">电压:</span>
                    <span class="value">{{ (osdData?.data?.backup_battery?.voltage / 1000)?.toFixed(2) || '--' }}V</span>
                  </div>
                  <div class="info-item">
                    <span class="label">温度:</span>
                    <span class="value">{{ osdData?.data?.backup_battery?.temperature || '--' }}°C</span>
                  </div>
                  <div class="info-item">
                    <span class="label">开关状态:</span>
                    <el-tag :type="osdData?.data?.backup_battery?.switch ? 'success' : 'info'" size="small">
                      {{ osdData?.data?.backup_battery?.switch ? '开启' : '关闭' }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Setting /></el-icon>
                    <span>云台和POE</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">云台支架状态:</span>
                    <el-tag :type="osdData?.data?.gimbal_holder_state ? 'warning' : 'success'" size="small">
                      {{ osdData?.data?.gimbal_holder_state ? '异常' : '正常' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">POE链路状态:</span>
                    <el-tag :type="osdData?.data?.poe_link_status ? 'success' : 'danger'" size="small">
                      {{ osdData?.data?.poe_link_status ? '已连接' : '未连接' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">POE功率输出:</span>
                    <span class="value">{{ osdData?.data?.poe_power_output || '--' }}W</span>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Lightning /></el-icon>
                    <span>电池维护</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">维护状态:</span>
                    <el-tag :type="osdData?.data?.drone_battery_maintenance_info?.maintenance_state ? 'warning' : 'success'" size="small">
                      {{ osdData?.data?.drone_battery_maintenance_info?.maintenance_state ? '维护中' : '正常' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">剩余维护时间:</span>
                    <span class="value">{{ osdData?.data?.drone_battery_maintenance_info?.maintenance_time_left || '--' }}s</span>
                  </div>
                  <div class="info-item">
                    <span class="label">加热状态:</span>
                    <el-tag :type="osdData?.data?.drone_battery_maintenance_info?.heat_state ? 'warning' : 'info'" size="small">
                      {{ osdData?.data?.drone_battery_maintenance_info?.heat_state ? '加热中' : '正常' }}
                    </el-tag>
                  </div>
                  <!-- 电池详情 -->
                  <div v-if="osdData?.data?.drone_battery_maintenance_info?.batteries?.length" style="margin-top: 8px;">
                    <div style="margin-bottom: 6px; color: var(--el-text-color-primary); font-size: 12px; font-weight: 500;">电池详情</div>
                    <div v-for="(battery, index) in osdData.data.drone_battery_maintenance_info.batteries" :key="index" class="battery-info">
                      <div class="info-item">
                        <span class="label">{{ getBatteryIndexText(battery.index) }}:</span>
                        <el-progress 
                          :percentage="getBatteryCapacity(battery.capacity_percent)"
                          :color="getBatteryColor(getBatteryCapacity(battery.capacity_percent))"
                          :show-text="true"
                          :stroke-width="4"
                        />
                      </div>
                      <div class="info-item">
                        <span class="label">电压:</span>
                        <span class="value">{{ getBatteryVoltage(battery.voltage) }}</span>
                      </div>
                      <div class="info-item">
                        <span class="label">温度:</span>
                        <span class="value">{{ getBatteryTemperature(battery.temperature) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </el-card>
            </el-col>

            <el-col :span="3">
              <el-card class="info-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><Setting /></el-icon>
                    <span>其他状态</span>
                  </div>
                </template>
                <div class="info-content">
                  <div class="info-item">
                    <span class="label">机盖状态:</span>
                    <el-tag :type="osdData?.data?.cover_state ? 'warning' : 'success'" size="small">
                      {{ osdData?.data?.cover_state ? '开启' : '关闭' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">补光灯:</span>
                    <el-tag :type="osdData?.data?.supplement_light_state ? 'success' : 'info'" size="small">
                      {{ osdData?.data?.supplement_light_state ? '开启' : '关闭' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">空调状态:</span>
                    <el-tag :type="osdData?.data?.air_conditioner?.air_conditioner_state ? 'success' : 'info'" size="small">
                      {{ osdData?.data?.air_conditioner?.air_conditioner_state ? '开启' : '关闭' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">告警状态:</span>
                    <el-tag :type="osdData?.data?.alarm_state ? 'danger' : 'success'" size="small">
                      {{ osdData?.data?.alarm_state ? '告警' : '正常' }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">用户体验改善:</span>
                    <el-tag :type="getUserExperienceType(osdData?.data?.user_experience_improvement)" size="small">
                      {{ getUserExperienceText(osdData?.data?.user_experience_improvement) }}
                    </el-tag>
                  </div>
                  <div class="info-item">
                    <span class="label">固件一致性:</span>
                    <el-tag :type="osdData?.data?.compatible_status ? 'warning' : 'success'" size="small">
                      {{ osdData?.data?.compatible_status ? '需要升级' : '正常' }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>





          <!-- 原始数据展示 -->
          <el-card class="raw-data-card" shadow="hover" style="margin-top: 15px;">
            <template #header>
              <div class="card-header">
                <el-icon><Document /></el-icon>
                <span>原始OSD数据</span>
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="copyRawData"
                  :icon="copySuccess ? 'Check' : 'CopyDocument'"
                  :class="{ 'copy-success': copySuccess }"
                  style="margin-left: auto;"
                >
                  {{ copySuccess ? '已复制' : '复制' }}
                </el-button>
              </div>
            </template>
            <div class="json-viewer">
              <pre class="json-content" v-html="formatJsonPayload(osdData)"></pre>
            </div>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  ArrowLeft, 
  Location, 
  Lightning, 
  Cloudy, 
  Connection, 
  Setting, 
  Folder, 
  Monitor, 
  Document 
} from '@element-plus/icons-vue'
import { useMqttProxyStore } from '@/stores/mqtt-proxy'
import { useDeviceStore } from '@/stores/device'

const route = useRoute()
const router = useRouter()
const mqttProxyStore = useMqttProxyStore()
const deviceStore = useDeviceStore()

// 响应式数据
const deviceSn = ref('')
const osdData = ref(null)
const hasReceivedData = ref(false)
const copySuccess = ref(false)
const messageListener = ref(null)
const lastProcessedMessageId = ref(null)
const updateTimeout = ref(null)
const dataStableCount = ref(0) // 数据稳定计数器
const lastDataUpdateTime = ref(null) // 最后数据更新时间
const dataTimeout = ref(null) // 数据超时定时器
const retryCount = ref(0) // 重试次数
const maxRetries = 3 // 最大重试次数

// 计算属性
const osdTopic = computed(() => `thing/product/${deviceSn.value}/osd`)

// 页面初始化
onMounted(() => {
  deviceSn.value = route.params.deviceSn || ''
  if (deviceSn.value) {
    // 清理缓存数据
    clearCachedData()
    subscribeToOSD()
    // 启动数据超时监控
    startDataTimeoutMonitor()
  } else {
    ElMessage.error('未找到设备SN参数')
    goBack()
  }
})

// 页面卸载时取消订阅
onUnmounted(() => {
  stopMessageListener()
  stopDataTimeoutMonitor()
  if (updateTimeout.value) {
    clearTimeout(updateTimeout.value)
  }
  if (osdTopic.value) {
    mqttProxyStore.unsubscribeTopics(osdTopic.value)
  }
})

// 订阅OSD主题
const subscribeToOSD = async () => {
  if (!mqttProxyStore.isConnected) {
    ElMessage.warning('MQTT未连接，无法订阅OSD数据')
    return
  }

  try {
    await mqttProxyStore.subscribeToTopics(osdTopic.value, 1)
    ElMessage.success(`已订阅OSD主题: ${osdTopic.value}`)
    
    // 开始监听消息
    startMessageListener()
  } catch (error) {
    console.error('订阅OSD主题失败:', error)
    ElMessage.error('订阅OSD主题失败')
  }
}

// 开始消息监听
const startMessageListener = () => {
  if (messageListener.value) {
    clearInterval(messageListener.value)
  }
  
  messageListener.value = setInterval(() => {
    handleOSDMessage()
  }, 2000) // 每2秒检查一次新消息，进一步减少频率
}

// 停止消息监听
const stopMessageListener = () => {
  if (messageListener.value) {
    clearInterval(messageListener.value)
    messageListener.value = null
  }
}

// 处理OSD消息
const handleOSDMessage = () => {
  // 从消息历史中查找最新的OSD消息
  const osdMessages = mqttProxyStore.messageHistory.filter(msg => msg.topic === osdTopic.value)
  if (osdMessages.length > 0) {
    const latestMessage = osdMessages[0] // 最新的消息在数组开头
    
    // 检查是否是新消息（通过消息ID比较）
    if (latestMessage.id !== lastProcessedMessageId.value) {
      try {
        const newData = JSON.parse(latestMessage.payload)
        
        // 验证数据完整性
        if (isValidOSDData(newData)) {
          // 进一步检查数据是否真的有变化（比较关键字段）
          if (!osdData.value || hasDataChanged(osdData.value, newData)) {
            // 使用防抖更新，避免频繁更新
            if (updateTimeout.value) {
              clearTimeout(updateTimeout.value)
            }
            
            updateTimeout.value = setTimeout(() => {
              // 实现局部更新机制
              if (osdData.value) {
                // 合并数据，保留已有数据，只更新新数据
                osdData.value = mergeData(osdData.value, newData)
              } else {
                // 如果没有现有数据，直接使用新数据
                osdData.value = newData
              }
              
              hasReceivedData.value = true
              lastProcessedMessageId.value = latestMessage.id
              lastDataUpdateTime.value = Date.now()
              
              // 重置重试计数
              retryCount.value = 0
              
              console.log('数据已更新:', new Date().toLocaleTimeString(), '消息ID:', latestMessage.id)
              console.log('数据内容:', newData)
            }, 200) // 增加防抖延迟到200ms
          } else {
            // 数据没有实质性变化，只更新消息ID
            lastProcessedMessageId.value = latestMessage.id
            console.log('OSD数据无变化，跳过更新')
          }
        } else {
          console.warn('收到无效的OSD数据，跳过更新:', newData)
          // 即使数据无效，也更新消息ID，避免重复处理
          lastProcessedMessageId.value = latestMessage.id
          // 如果当前没有数据，尝试使用部分有效数据
          if (!osdData.value && newData && typeof newData === 'object') {
            console.log('尝试使用部分有效数据')
            osdData.value = newData
            hasReceivedData.value = true
            lastDataUpdateTime.value = Date.now()
          }
        }
      } catch (error) {
        console.error('解析OSD数据失败:', error, '原始数据:', latestMessage.payload)
        // 解析失败也更新消息ID
        lastProcessedMessageId.value = latestMessage.id
      }
    }
  } else {
    // 如果没有OSD消息，但已经有数据，保持现有数据不变
    if (osdData.value && hasReceivedData.value) {
      console.log('暂无新OSD消息，保持现有数据显示')
      // 确保数据状态标志保持为true
      if (!hasReceivedData.value) {
        hasReceivedData.value = true
      }
    }
  }
}

// 验证数据是否有效 - 支持多种数据格式
const isValidOSDData = (data) => {
  if (!data || typeof data !== 'object') return false
  
  // 检查是否有data字段
  if (data.data && typeof data.data === 'object') {
    // 检查是否包含已知的数据格式字段
    const dataContent = data.data
    
    // 第一种格式：OSD数据（位置、通信等）
    const osdFields = ['departure_point', 'sdr', 'wireless_link', 'tilt_angle', 'drc_state']
    const hasOSDFields = osdFields.some(field => field in dataContent)
    
    // 第二种格式：设备状态数据（工作状态、电力等）
    const deviceFields = ['job_number', 'acc_time', 'electric_supply_voltage', 'backup_battery', 'drone_battery_maintenance_info']
    const hasDeviceFields = deviceFields.some(field => field in dataContent)
    
    // 第三种格式：实时状态数据（位置、环境等）
    const realtimeFields = ['latitude', 'longitude', 'height', 'drone_charge_state', 'network_state', 'temperature']
    const hasRealtimeFields = realtimeFields.some(field => field in dataContent)
    
    return hasOSDFields || hasDeviceFields || hasRealtimeFields
  }
  
  return false
}

// 检查数据是否真的有变化 - 优化变化检测
const hasDataChanged = (oldData, newData) => {
  if (!oldData || !newData) return true
  
  // 比较时间戳
  const oldTimestamp = oldData.timestamp || 0
  const newTimestamp = newData.timestamp || 0
  
  // 如果时间戳不同，说明是新数据
  if (oldTimestamp !== newTimestamp) {
    return true
  }
  
  // 比较关键数据字段，降低变化检测的严格度
  const keyFields = ['latitude', 'longitude', 'height', 'drone_charge_state', 'environment_temperature', 'network_state', 'temperature']
  for (const field of keyFields) {
    const oldValue = oldData.data?.[field]
    const newValue = newData.data?.[field]
    
    // 对于数值类型，允许小的误差
    if (typeof oldValue === 'number' && typeof newValue === 'number') {
      if (Math.abs(oldValue - newValue) > 0.001) {
        return true
      }
    } else if (oldValue !== newValue) {
      return true
    }
  }
  
  return false
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 数据合并函数 - 实现局部更新
const mergeData = (oldData, newData) => {
  if (!oldData) return newData
  if (!newData) return oldData
  
  // 深度合并数据
  const merged = { ...oldData }
  
  // 合并顶层字段
  Object.keys(newData).forEach(key => {
    if (key === 'data' && newData.data && typeof newData.data === 'object') {
      // 合并data字段
      merged.data = { ...oldData.data, ...newData.data }
    } else {
      // 合并其他字段
      merged[key] = newData[key]
    }
  })
  
  return merged
}


// 数据超时监控
const startDataTimeoutMonitor = () => {
  if (dataTimeout.value) {
    clearInterval(dataTimeout.value)
  }
  
  dataTimeout.value = setInterval(() => {
    const now = Date.now()
    const lastUpdate = lastDataUpdateTime.value
    
    // 如果超过30秒没有数据更新，显示警告
    if (lastUpdate && now - lastUpdate > 30 * 1000) {
      console.warn('数据更新超时，尝试重新订阅')
      ElMessage.warning('数据更新超时，正在尝试重新连接...')
      
      // 尝试重新订阅
      if (retryCount.value < maxRetries) {
        retryCount.value++
        subscribeToOSD()
      } else {
        ElMessage.error('多次重试失败，请检查网络连接')
        retryCount.value = 0
      }
    }
  }, 10000) // 每10秒检查一次
}

const stopDataTimeoutMonitor = () => {
  if (dataTimeout.value) {
    clearInterval(dataTimeout.value)
    dataTimeout.value = null
  }
}

// 刷新数据
const refreshData = () => {
  console.log('手动刷新数据')
  ElMessage.info('正在刷新数据...')
  
  // 重置状态
  hasReceivedData.value = false
  lastProcessedMessageId.value = null
  retryCount.value = 0
  
  // 清理缓存数据
  clearCachedData()
  
  // 重新订阅
  subscribeToOSD()
}

// 清理缓存数据
const clearCachedData = () => {
  try {
    // 清理当前设备的缓存
    localStorage.removeItem(`osd_data_${deviceSn.value}`)
    console.log('已清理缓存数据')
  } catch (error) {
    console.error('清理缓存数据失败:', error)
  }
}

// 获取数据状态类型
const getDataStatusType = () => {
  if (!lastDataUpdateTime.value) return 'info'
  
  const now = Date.now()
  const timeDiff = now - lastDataUpdateTime.value
  
  if (timeDiff < 10000) return 'success' // 10秒内
  if (timeDiff < 30000) return 'warning' // 30秒内
  return 'danger' // 超过30秒
}

// 获取数据状态文本
const getDataStatusText = () => {
  if (!lastDataUpdateTime.value) return '无数据'
  
  const now = Date.now()
  const timeDiff = now - lastDataUpdateTime.value
  
  if (timeDiff < 10000) return '实时'
  if (timeDiff < 30000) return '延迟'
  return '超时'
}

// 格式化存储大小
const formatStorage = (bytes) => {
  if (!bytes) return '--'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  
  return `${size.toFixed(2)} ${units[unitIndex]}`
}

// 获取存储使用率
const getStorageUsage = (storage) => {
  if (!storage || !storage.total) return 0
  return Math.round((storage.used / storage.total) * 100)
}

// 获取存储使用率颜色
const getStorageUsageColor = (storage) => {
  const usage = getStorageUsage(storage)
  if (usage > 90) return '#f56c6c'
  if (usage > 70) return '#e6a23c'
  return '#67c23a'
}

// 获取电池颜色
const getBatteryColor = (percentage) => {
  if (percentage > 50) return '#67c23a'
  if (percentage > 20) return '#e6a23c'
  return '#f56c6c'
}

// 获取充电状态类型
const getChargeStateType = (state) => {
  const stateMap = {
    0: 'success', // 正常
    1: 'warning', // 充电中
    2: 'danger'   // 异常
  }
  return stateMap[state] || 'info'
}

// 获取充电状态文本
const getChargeStateText = (state) => {
  const stateMap = {
    0: '正常',
    1: '充电中',
    2: '异常'
  }
  return stateMap[state] || '未知'
}

// 获取电池存储模式类型
const getBatteryStoreModeType = (mode) => {
  const modeMap = {
    0: 'success', // 正常模式
    1: 'warning', // 存储模式
    2: 'info'     // 其他
  }
  return modeMap[mode] || 'info'
}

// 获取电池存储模式文本
const getBatteryStoreModeText = (mode) => {
  const modeMap = {
    0: '正常模式',
    1: '存储模式',
    2: '其他模式'
  }
  return modeMap[mode] || '未知'
}

// 获取网络类型颜色
const getNetworkTypeColor = (type) => {
  const typeMap = {
    1: 'primary', // 4G
    2: 'success'  // 以太网
  }
  return typeMap[type] || 'info'
}

// 获取网络类型文本
const getNetworkTypeText = (type) => {
  const typeMap = {
    1: '4G',
    2: '以太网'
  }
  return typeMap[type] || '未知'
}

// 获取网络质量颜色
const getNetworkQualityColor = (quality) => {
  if (quality >= 4) return 'success'  // 较好、好
  if (quality >= 2) return 'warning'  // 较差、一般
  if (quality >= 1) return 'danger'   // 差
  return 'info'  // 无信号
}

// 获取网络质量进度条颜色
const getNetworkQualityProgressColor = (quality) => {
  if (quality >= 4) return '#67c23a'  // 较好、好
  if (quality >= 2) return '#e6a23c'  // 较差、一般
  if (quality >= 1) return '#f56c6c'  // 差
  return '#909399'  // 无信号
}

// 获取机场序列号
const getAirportSn = () => {
  // 优先级：设备存储中的机场序列号 > MQTT数据中的gateway > 当前设备SN
  return deviceStore.currentDevice?.airport_sn || 
         osdData.value?.gateway || 
         deviceSn.value || 
         '--'
}

// 获取网络质量文本
const getNetworkQualityText = (quality) => {
  const qualityMap = {
    0: '无信号',
    1: '差',
    2: '较差',
    3: '一般',
    4: '较好',
    5: '好'
  }
  return qualityMap[quality] || '未知'
}

// 获取网络传输速率文本
const getNetworkRateText = (rate) => {
  if (rate === null || rate === undefined) return '--'
  return `${rate} KB/s`
}

// 获取定位状态类型
const getPositionFixedType = (fixed) => {
  const fixedMap = {
    0: 'danger',  // 未定位
    1: 'warning', // 单点定位
    2: 'success'  // 差分定位
  }
  return fixedMap[fixed] || 'info'
}

// 获取定位状态文本
const getPositionFixedText = (fixed) => {
  const fixedMap = {
    0: '未定位',
    1: '单点定位',
    2: '差分定位'
  }
  return fixedMap[fixed] || '未知'
}

// 获取定位质量颜色
const getPositionQualityColor = (quality) => {
  if (quality >= 4) return '#67c23a'
  if (quality >= 2) return '#e6a23c'
  return '#f56c6c'
}

// 获取模式代码类型
const getModeCodeType = (code) => {
  const codeMap = {
    0: 'success', // 空闲中
    1: 'info',    // 现场调试
    2: 'info',    // 远程调试
    3: 'warning', // 固件升级中
    4: 'primary', // 作业中
    5: 'warning'  // 待标定
  }
  return codeMap[code] || 'info'
}

// 获取模式代码文本
const getModeCodeText = (code) => {
  const codeMap = {
    0: '空闲中',
    1: '现场调试',
    2: '远程调试',
    3: '固件升级中',
    4: '作业中',
    5: '待标定'
  }
  return codeMap[code] || '未知'
}

// 获取电池序号文本
const getBatteryIndexText = (index) => {
  const indexMap = {
    0: '左电池',
    1: '右电池'
  }
  return indexMap[index] || `电池${index + 1}`
}

// 获取电池电量（处理异常值32767）
const getBatteryCapacity = (capacity) => {
  if (capacity === null || capacity === undefined) return 0
  if (capacity === 32767) return 0 // 异常值
  return Math.max(0, Math.min(100, capacity)) // 确保在0-100范围内
}

// 获取电池电压（处理异常值32767）
const getBatteryVoltage = (voltage) => {
  if (voltage === null || voltage === undefined) return '--'
  if (voltage === 32767) return '异常'
  return `${(voltage / 1000).toFixed(2)}V`
}

// 获取电池温度（处理异常值32767）
const getBatteryTemperature = (temperature) => {
  if (temperature === null || temperature === undefined) return '--'
  if (temperature === 32767) return '异常'
  return `${temperature.toFixed(1)}°C`
}

// 获取降雨量类型
const getRainfallType = (rainfall) => {
  const typeMap = {
    0: 'success', // 无雨
    1: 'info',    // 小雨
    2: 'warning', // 中雨
    3: 'danger'   // 大雨
  }
  return typeMap[rainfall] || 'info'
}

// 获取降雨量文本
const getRainfallText = (rainfall) => {
  const textMap = {
    0: '无雨',
    1: '小雨',
    2: '中雨',
    3: '大雨'
  }
  return textMap[rainfall] || '未知'
}

// 获取用户体验改善类型
const getUserExperienceType = (status) => {
  const typeMap = {
    0: 'info',    // 初始状态
    1: 'warning', // 拒绝加入
    2: 'success'  // 同意加入
  }
  return typeMap[status] || 'info'
}

// 获取用户体验改善文本
const getUserExperienceText = (status) => {
  const textMap = {
    0: '初始状态',
    1: '拒绝加入',
    2: '同意加入'
  }
  return textMap[status] || '未知'
}

// 获取机场任务状态类型
const getFlightTaskStepType = (code) => {
  const typeMap = {
    0: 'info',    // 作业准备中
    1: 'primary', // 飞行作业中
    2: 'warning', // 作业后状态恢复
    3: 'info',    // 自定义飞行区更新中
    4: 'info',    // 地形障碍物更新中
    5: 'success', // 任务空闲
    255: 'danger', // 飞行器异常
    256: 'info'   // 未知状态
  }
  return typeMap[code] || 'info'
}

// 获取机场任务状态文本
const getFlightTaskStepText = (code) => {
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
  return textMap[code] || '未知'
}

// 格式化JSON数据
const formatJsonPayload = (payload) => {
  try {
    let jsonString
    if (typeof payload === 'string') {
      const parsed = JSON.parse(payload)
      jsonString = JSON.stringify(parsed, null, 2)
    } else if (typeof payload === 'object' && payload !== null) {
      jsonString = JSON.stringify(payload, null, 2)
    } else {
      return String(payload)
    }
    
    return jsonString
      .replace(/(".*?")\s*:/g, '<span class="json-key">$1</span>:')
      .replace(/:\s*(".*?")/g, ': <span class="json-string">$1</span>')
      .replace(/:\s*(\d+)/g, ': <span class="json-number">$1</span>')
      .replace(/:\s*(true|false)/g, ': <span class="json-boolean">$1</span>')
      .replace(/:\s*(null)/g, ': <span class="json-null">$1</span>')
      .replace(/([{}[\]])/g, '<span class="json-bracket">$1</span>')
      .replace(/(,)/g, '<span class="json-comma">$1</span>')
  } catch (error) {
    return String(payload)
  }
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '--'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 格式化累计时间（毫秒转换为小时）
const formatAccTime = (accTime) => {
  if (!accTime) return '--'
  const hours = Math.floor(accTime / (1000 * 60 * 60))
  const minutes = Math.floor((accTime % (1000 * 60 * 60)) / (1000 * 60))
  return `${hours}小时${minutes}分钟`
}

// 格式化时间戳
const formatTimestamp = (timestamp) => {
  if (!timestamp) return '--'
  const date = new Date(timestamp * 1000) // 假设是秒级时间戳
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 复制原始数据
const copyRawData = async () => {
  if (!osdData.value) return
  
  try {
    const jsonString = JSON.stringify(osdData.value, null, 2)
    await navigator.clipboard.writeText(jsonString)
    copySuccess.value = true
    
    setTimeout(() => {
      copySuccess.value = false
    }, 2000)
    
    ElMessage.success('OSD数据已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败，请手动复制')
  }
}
</script>

<style scoped>
.aircraft-detail {
  height: 100vh;
  background: var(--el-bg-color-page);
}

.detail-header {
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-left h2 {
  margin: 0;
  color: var(--el-text-color-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.detail-main {
  padding: 15px 10px;
  background: var(--el-bg-color-page);
}

.connection-warning {
  margin-bottom: 20px;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
}

.data-container {
  width: 100%;
  margin: 0;
}

.data-grid {
  margin-bottom: 0;
}

.info-card {
  height: auto;
  margin-bottom: 15px;
  min-height: 180px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.info-content {
  padding: 6px 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  padding: 6px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
  gap: 10px;
}

.info-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.label {
  font-weight: 500;
  color: var(--el-text-color-regular);
  min-width: 60px;
  font-size: 13px;
  flex-shrink: 0;
}

.value {
  color: var(--el-text-color-primary);
  font-weight: 600;
  font-size: 13px;
  flex: 1;
  text-align: right;
}

.raw-data-card {
  margin-top: 20px;
}

.json-viewer {
  background: #1e1e1e;
  border: 1px solid #3c3c3c;
  border-radius: 8px;
  padding: 16px;
  max-height: 400px;
  overflow-y: auto;
}

.json-content {
  margin: 0;
  font-family: 'Fira Code', 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #d4d4d4;
  white-space: pre-wrap;
  word-break: break-word;
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

.copy-success {
  background-color: var(--el-color-success) !important;
  border-color: var(--el-color-success) !important;
  color: white !important;
}

.battery-info {
  background: var(--el-bg-color-page);
  border: 1px solid var(--el-border-color-light);
  border-radius: 6px;
  padding: 8px;
  margin-bottom: 8px;
}

.battery-info .info-item {
  margin-bottom: 6px;
  padding: 3px 0;
}

.battery-info .info-item:last-child {
  margin-bottom: 0;
}
</style>
