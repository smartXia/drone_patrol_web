import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { DEVICE_SN, GATEWAY_SN, API_BASE_URL } from '../config.js'
import axios from 'axios'

export const useDeviceStore = defineStore('device', () => {
  // 状态
  const currentDeviceSn = ref('')
  const currentGatewaySn = ref('')
  const deviceList = ref([])

  // 计算属性
  const currentDevice = computed(() => {
    return deviceList.value.find(device => device.sn === currentDeviceSn.value) || null
  })

  const currentGateway = computed(() => {
    return deviceList.value.find(device => device.sn === currentGatewaySn.value) || null
  })

  const onlineDevices = computed(() => {
    return deviceList.value.filter(device => device.status === 'online')
  })

  const offlineDevices = computed(() => {
    return deviceList.value.filter(device => device.status === 'offline')
  })

  // 方法
  const setCurrentDevice = async (deviceSn) => {
    const device = deviceList.value.find(d => d.sn === deviceSn)
    if (device) {
      try {
        const response = await axios.post(`${API_BASE_URL}/api/devices/${device.id}/set-current`)
        if (response.data.code === 0) {
          currentDeviceSn.value = deviceSn
          // 如果是网关设备，同时更新网关SN
          if (device.type === 'gateway') {
            currentGatewaySn.value = deviceSn
          }
          console.log('切换到设备:', device.name, deviceSn)
          // 重新加载设备列表以更新状态
          await loadDevices()
        } else {
          console.error('设置当前设备失败:', response.data.message)
        }
      } catch (error) {
        console.error('设置当前设备失败:', error)
      }
    } else {
      console.warn('设备不存在:', deviceSn)
    }
  }

  const setCurrentGateway = async (gatewaySn) => {
    const gateway = deviceList.value.find(d => d.sn === gatewaySn)
    if (gateway) {
      try {
        const response = await axios.post(`${API_BASE_URL}/api/devices/${gateway.id}/set-gateway`)
        if (response.data.code === 0) {
          currentGatewaySn.value = gatewaySn
          console.log('切换到网关:', gateway.name, gatewaySn)
          // 重新加载设备列表以更新状态
          await loadDevices()
        } else {
          console.error('设置当前网关失败:', response.data.message)
        }
      } catch (error) {
        console.error('设置当前网关失败:', error)
      }
    } else {
      console.warn('网关不存在:', gatewaySn)
    }
  }

  const addDevice = async (device) => {
    try {
      const response = await axios.post(`${API_BASE_URL}/api/devices`, {
        name: device.name,
        sn: device.sn,
        type: device.type || 'drone',
        status: device.status || 'offline',
        airport_sn: device.airport_sn || ''
      })
      
      if (response.data.code === 0) {
        console.log('添加设备成功:', device.name, device.sn)
        // 重新加载设备列表
        await loadDevices()
        return response.data.data
      } else {
        console.error('添加设备失败:', response.data.message)
        throw new Error(response.data.message)
      }
    } catch (error) {
      console.error('添加设备失败:', error)
      throw error
    }
  }

  const updateDevice = async (deviceId, deviceData) => {
    try {
      const response = await axios.put(`${API_BASE_URL}/api/devices/${deviceId}`, {
        name: deviceData.name,
        sn: deviceData.sn,
        type: deviceData.type,
        status: deviceData.status,
        airport_sn: deviceData.airport_sn || ''
      })
      
      if (response.data.code === 0) {
        console.log('更新设备成功:', deviceData.name, deviceData.sn)
        // 重新加载设备列表
        await loadDevices()
        return response.data
      } else {
        console.error('更新设备失败:', response.data.message)
        throw new Error(response.data.message)
      }
    } catch (error) {
      console.error('更新设备失败:', error)
      throw error
    }
  }

  const removeDevice = async (deviceSn) => {
    const device = deviceList.value.find(d => d.sn === deviceSn)
    if (!device) {
      console.error('设备不存在:', deviceSn)
      throw new Error('设备不存在')
    }
    
    try {
      console.log('正在删除设备:', device.name, 'ID:', device.id, 'SN:', device.sn)
      const response = await axios.delete(`${API_BASE_URL}/api/devices/${device.id}`)
      if (response.data.code === 0) {
        console.log('删除设备成功:', device.name, deviceSn)
        // 如果删除的是当前设备，清空当前设备选择
        if (currentDeviceSn.value === deviceSn) {
          currentDeviceSn.value = ''
        }
        // 重新加载设备列表
        await loadDevices()
      } else {
        console.error('删除设备失败:', response.data.message)
        throw new Error(response.data.message)
      }
    } catch (error) {
      console.error('删除设备失败:', error)
      throw error
    }
  }

  const updateDeviceStatus = async (deviceSn, status) => {
    const device = deviceList.value.find(d => d.sn === deviceSn)
    if (device) {
      try {
        const response = await axios.put(`${API_BASE_URL}/api/devices/${device.id}`, {
          status: status
        })
        if (response.data.code === 0) {
          device.status = status
          device.lastSeen = new Date().toISOString()
          console.log('更新设备状态成功:', device.name, status)
        } else {
          console.error('更新设备状态失败:', response.data.message)
        }
      } catch (error) {
        console.error('更新设备状态失败:', error)
      }
    }
  }

  const getDeviceBySn = (deviceSn) => {
    return deviceList.value.find(d => d.sn === deviceSn)
  }

  const getTopicWithDeviceSn = (topicTemplate) => {
    // 获取当前设备信息
    const currentDevice = deviceList.value.find(device => device.sn === currentDeviceSn.value)
    const airportSn = currentDevice?.airport_sn || currentDeviceSn.value
    
    return topicTemplate
      .replace(/\{device_sn\}/g, airportSn || 'DEVICE_SN')
      .replace(/\{gateway_sn\}/g, currentGatewaySn.value || 'GATEWAY_SN')
  }


  // 从 SQLite 数据库加载设备列表
  const loadDevices = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/api/devices`)
      if (response.data.code === 0) {
        deviceList.value = response.data.data
        console.log('从数据库加载设备列表:', deviceList.value.length, '个设备')
        console.log('设备列表详情:', deviceList.value)
        
        // 检查是否有默认设备
        const defaultDevices = deviceList.value.filter(device => device.name === '默认设备')
        if (defaultDevices.length > 0) {
          console.warn('发现默认设备，建议使用"删除默认设备"按钮清理:', defaultDevices)
        }
        
        // 更新当前设备信息
        const currentResponse = await axios.get(`${API_BASE_URL}/api/devices/current`)
        if (currentResponse.data.code === 0) {
          const currentData = currentResponse.data.data
          if (currentData.device) {
            currentDeviceSn.value = currentData.device.sn
          } else {
            currentDeviceSn.value = ''
          }
          if (currentData.gateway) {
            currentGatewaySn.value = currentData.gateway.sn
          } else {
            currentGatewaySn.value = ''
          }
        } else {
          // 如果没有当前设备，清空选择
          currentDeviceSn.value = ''
          currentGatewaySn.value = ''
        }
      } else {
        console.error('加载设备列表失败:', response.data.message)
      }
    } catch (error) {
      console.error('加载设备列表失败:', error)
    }
  }


  // 清空所有设备
  const clearAllDevices = async () => {
    try {
      const response = await axios.delete(`${API_BASE_URL}/api/devices/clear`)
      if (response.data.code === 0) {
        console.log('所有设备已清空')
        // 重新加载设备列表
        await loadDevices()
        // 清空当前设备选择
        currentDeviceSn.value = ''
        currentGatewaySn.value = ''
        return response.data
      } else {
        console.error('清空设备失败:', response.data.message)
        throw new Error(response.data.message)
      }
    } catch (error) {
      console.error('清空设备失败:', error)
      throw error
    }
  }

  // 删除所有默认设备
  const removeDefaultDevices = async () => {
    try {
      const response = await axios.delete(`${API_BASE_URL}/api/devices/remove-defaults`)
      if (response.data.code === 0) {
        console.log('默认设备已删除')
        // 重新加载设备列表
        await loadDevices()
        return response.data
      } else {
        console.error('删除默认设备失败:', response.data.message)
        throw new Error(response.data.message)
      }
    } catch (error) {
      console.error('删除默认设备失败:', error)
      throw error
    }
  }

  // 自动清理默认设备
  const autoCleanDefaultDevices = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/api/devices`)
      if (response.data.code === 0) {
        const devices = response.data.data
        const defaultDevices = devices.filter(device => device.name === '默认设备')
        if (defaultDevices.length > 0) {
          console.log('发现默认设备，自动清理中...')
          await removeDefaultDevices()
        }
      }
    } catch (error) {
      console.error('自动清理默认设备失败:', error)
    }
  }

  return {
    // 状态
    currentDeviceSn,
    currentGatewaySn,
    deviceList,
    
    // 计算属性
    currentDevice,
    currentGateway,
    onlineDevices,
    offlineDevices,
    
    // 方法
    setCurrentDevice,
    setCurrentGateway,
    addDevice,
    updateDevice,
    removeDevice,
    updateDeviceStatus,
    getDeviceBySn,
    getTopicWithDeviceSn,
    loadDevices,
    clearAllDevices,
    removeDefaultDevices,
    autoCleanDefaultDevices
  }
})
