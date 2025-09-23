import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import { generateClientId } from '../utils/mqtt'

export const useMqttStore = defineStore('mqtt', () => {
  // 状态
  const client = ref(null)
  const isConnected = ref(false)
  const isConnecting = ref(false)
  const connectionError = ref(null)
  const deviceStatus = ref({})
  const messageHistory = ref([])

  // 动态配置与配置集
  const currentConfig = ref(null)
  const profiles = ref([])
  const currentProfileId = ref(null)

  // 后端本地服务地址
  const apiBase = import.meta.env?.VITE_MQTT_API_BASE || 'http://127.0.0.1:8080'

  // 计算属性
  const connectionStatus = computed(() => {
    if (isConnected.value) return 'connected'
    if (isConnecting.value) return 'connecting'
    return 'disconnected'
  })

  // 默认配置（用于首次使用或未选择配置时）
  const defaultConfig = {
    protocol: 'ws',
    host: '',
    port: undefined,
    path: '/mqtt',
    username: '',
    password: '',
    clientId: generateClientId('web_client'),
    clean: true,
    reconnectPeriod: 0,
    connectTimeout: 30000,
    keepalive: 60,
    rejectUnauthorized: false
  }

  // 连接MQTT
  const connect = async () => {
    if (isConnecting.value || isConnected.value) {
      console.log('MQTT已连接或正在连接中，跳过重复连接')
      return
    }

    try {
      isConnecting.value = true
      connectionError.value = null

      // 如果没有配置，先加载配置
      if (!currentConfig.value) {
        await loadProfiles()
      }

      const cfg = currentConfig.value || defaultConfig
      const url = buildMqttUrl(cfg)
      const options = {
        clientId: cfg.clientId || generateClientId('web_client'),
        username: cfg.username,
        password: cfg.password,
        clean: cfg.clean !== false,
        reconnectPeriod: cfg.reconnectPeriod ?? 0,
        connectTimeout: cfg.connectTimeout ?? 30000,
        keepalive: cfg.keepalive ?? 60,
        protocolVersion: 4,
        rejectUnauthorized: !!cfg.rejectUnauthorized
      }
      
      console.log('正在连接MQTT:', url)
      console.log('连接配置:', {
        host: cfg.host,
        port: cfg.port,
        protocol: cfg.protocol,
        path: cfg.path,
        username: cfg.username ? '***' : '未设置',
        clientId: options.clientId
      })
      
      client.value = mqtt.connect(url, options)

            client.value.on('connect', () => {
        console.log('MQTT连接成功')
        isConnected.value = true
        isConnecting.value = false
      })

      client.value.on('message', (topic, message) => {
        console.log('收到MQTT消息:', topic, message.toString())
        handleMessage(topic, message)
      })

      client.value.on('error', (error) => {
        console.error('MQTT连接错误:', error)
        console.error('错误详情:', {
          message: error.message,
          code: error.code,
          stack: error.stack,
          url: url,
          config: {
            host: cfg.host,
            port: cfg.port,
            protocol: cfg.protocol,
            path: cfg.path
          }
        })
        
        // 提供更友好的错误信息
        let errorMessage = error.message
        if (error.code === 'ECONNREFUSED') {
          errorMessage = `连接被拒绝 - 请检查服务器地址 ${cfg.host}:${cfg.port} 是否正确，以及服务器是否运行`
        } else if (error.code === 'ENOTFOUND') {
          errorMessage = `无法解析主机名 ${cfg.host} - 请检查网络连接和DNS设置`
        } else if (error.message.includes('WebSocket connection failed')) {
          errorMessage = `WebSocket连接失败 - 请检查服务器是否支持WebSocket协议，或尝试使用TCP协议`
        } else if (error.message.includes('timeout')) {
          errorMessage = `连接超时 - 请检查网络连接和防火墙设置`
        }
        
        connectionError.value = errorMessage
        isConnected.value = false
        isConnecting.value = false
      })

      client.value.on('close', () => {
        console.log('MQTT连接关闭')
        isConnected.value = false
        isConnecting.value = false
      })

      client.value.on('reconnect', () => {
        console.log('MQTT重新连接中...')
        isConnecting.value = true
      })

    } catch (error) {
      console.error('MQTT连接失败:', error)
      connectionError.value = error.message
      isConnecting.value = false
    }
  }

  // 使用指定配置连接
  const connectWithConfig = async (config) => {
    currentConfig.value = { ...config, clientId: config.clientId || generateClientId('web_client') }
    return connect()
  }

  // 断开连接
  const disconnect = () => {
    if (client.value) {
      console.log('正在断开MQTT连接...')
      client.value.end(true) // 强制断开
      client.value = null
    }
    isConnected.value = false
    isConnecting.value = false
    connectionError.value = null
  }

  // ===== 配置管理（对接本地 Node 服务） =====
  const loadProfiles = async () => {
    const { data } = await axios.get(`${apiBase}/api/mqtt/profiles`)
    if (data && data.code === 0) {
      profiles.value = data.data || []
      const d = profiles.value.find(p => p.isDefault)
      if (d) {
        currentProfileId.value = d.id
        currentConfig.value = d.config
      }
    }
  }

  const saveProfile = async (payload) => {
    const { data } = await axios.post(`${apiBase}/api/mqtt/profiles`, payload)
    if (data && data.code === 0) {
      await loadProfiles()
      return data.data?.id
    }
    throw new Error(data?.message || '保存失败')
  }

  const updateProfile = async (id, payload) => {
    const { data } = await axios.put(`${apiBase}/api/mqtt/profiles/${id}`, payload)
    if (data && data.code === 0) {
      await loadProfiles()
      return true
    }
    throw new Error(data?.message || '更新失败')
  }

  const deleteProfile = async (id) => {
    const { data } = await axios.delete(`${apiBase}/api/mqtt/profiles/${id}`)
    if (data && data.code === 0) {
      await loadProfiles()
      if (currentProfileId.value === id) {
        currentProfileId.value = null
        currentConfig.value = null
      }
      return true
    }
    throw new Error(data?.message || '删除失败')
  }

  const setDefaultProfile = async (id) => {
    const { data } = await axios.post(`${apiBase}/api/mqtt/profiles/${id}/default`)
    if (data && data.code === 0) {
      await loadProfiles()
      return true
    }
    throw new Error(data?.message || '设置默认失败')
  }

  const selectProfile = async (id) => {
    const prof = profiles.value.find(p => p.id === id)
    if (!prof) throw new Error('配置不存在')
    currentProfileId.value = id
    currentConfig.value = prof.config
  }

  const setConfig = (config) => {
    currentConfig.value = { ...config }
  }

  const testConnection = async (payload) => {
    const { data } = await axios.post(`${apiBase}/api/mqtt/test`, payload)
    if (data && data.code === 0) return true
    throw new Error(data?.message || '测试连接失败')
  }

  // 检查网络连通性
  const checkNetworkConnectivity = async (host, port) => {
    try {
      const response = await axios.get(`${apiBase}/api/network/ping`, {
        params: { host, port },
        timeout: 5000
      })
      
      // 转换后端API返回的数据格式
      const { code, data, message } = response.data
      if (code === 0) {
        return {
          success: data.success,
          duration: data.duration,
          host: data.host,
          port: data.port
        }
      } else {
        return {
          success: false,
          error: message || '网络检查失败',
          duration: data?.duration || 0
        }
      }
    } catch (error) {
      console.error('网络连通性检查失败:', error)
      return { 
        success: false, 
        error: error.response?.data?.message || error.message || '网络检查失败',
        duration: 0
      }
    }
  }

  // 订阅主题
  const subscribeToTopics = (topics) => {
    if (!client.value || !isConnected.value) {
      throw new Error('MQTT未连接')
    }
    
    if (!topics || (Array.isArray(topics) && topics.length === 0)) {
      throw new Error('主题不能为空')
    }

    // 将单个主题转换为数组
    const topicArray = Array.isArray(topics) ? topics : [topics]
    
    console.log('开始订阅主题...')
    console.log(`尝试订阅主题: ${topicArray.join(', ')}`)
    
    return new Promise((resolve, reject) => {
      // 订阅多个主题
      client.value.subscribe(topicArray, { qos: 0 }, (err, granted) => {
        if (err) {
          console.error(`订阅主题失败:`, err)
          reject(err)
        } else {
          console.log(`成功订阅主题:`, granted)
          resolve(granted)
        }
      })
    })
  }

  // 取消订阅主题
  const unsubscribeTopics = (topics) => {
    if (!client.value || !isConnected.value) {
      throw new Error('MQTT未连接')
    }

    if (!topics || (Array.isArray(topics) && topics.length === 0)) {
      throw new Error('主题不能为空')
    }

    const topicArray = Array.isArray(topics) ? topics : [topics]

    return new Promise((resolve, reject) => {
      client.value.unsubscribe(topicArray, (err) => {
        if (err) {
          console.error('取消订阅失败:', err)
          reject(err)
        } else {
          console.log('已取消订阅:', topicArray)
          resolve(true)
        }
      })
    })
  }

  // 处理接收到的消息
  const handleMessage = (topic, message) => {
    try {
      console.log('处理消息:', topic, message.toString())
      console.log('消息原始内容:', message)
      
      let payload
      try {
        payload = JSON.parse(message.toString())
        console.log('解析为JSON成功:', payload)
      } catch (e) {
        // 如果不是JSON格式，直接使用字符串
        payload = { data: message.toString() }
        console.log('解析为字符串:', payload)
      }
      
      const timestamp = new Date().toISOString()
      
      // 解析主题获取设备ID
      const topicParts = topic.split('/')
      // 对于 thing/product/8UUXN2B00A00ST/events 格式，设备ID是第3部分
      const deviceId = topicParts[2] || topicParts[1] || 'unknown'
      console.log('解析的设备ID:', deviceId)
      console.log('主题分割结果:', topicParts)
      
      // 更新设备状态
      if (!deviceStatus.value[deviceId]) {
        deviceStatus.value[deviceId] = {}
      }
      
      deviceStatus.value[deviceId] = {
        ...deviceStatus.value[deviceId],
        ...payload,
        lastUpdate: timestamp,
        topic: topic
      }

      // 添加到消息历史
      messageHistory.value.push({
        timestamp,
        topic,
        deviceId,
        payload
      })

      console.log('消息历史记录数量:', messageHistory.value.length)
      console.log('设备状态数量:', Object.keys(deviceStatus.value).length)

      // 限制历史记录数量
      if (messageHistory.value.length > 1000) {
        messageHistory.value = messageHistory.value.slice(-1000)
      }

    } catch (error) {
      console.error('解析MQTT消息失败:', error)
    }
  }

  // 发布消息
  const publish = (topic, message) => {
    if (!client.value || !isConnected.value) {
      throw new Error('MQTT未连接')
    }
    
    client.value.publish(topic, JSON.stringify(message))
  }

  // 获取设备状态
  const getDeviceStatus = (deviceId) => {
    return deviceStatus.value[deviceId] || null
  }

  // 获取所有设备状态
  const getAllDeviceStatus = () => {
    return deviceStatus.value
  }

  // 清空消息历史
  const clearMessageHistory = () => {
    messageHistory.value = []
  }

  return {
    // 状态
    isConnected,
    isConnecting,
    connectionError,
    deviceStatus,
    messageHistory,
    currentConfig,
    profiles,
    currentProfileId,
    
    // 计算属性
    connectionStatus,
    
    // 方法
    connect,
    connectWithConfig,
    disconnect,
    subscribeToTopics,
    unsubscribeTopics,
    publish,
    getDeviceStatus,
    getAllDeviceStatus,
    clearMessageHistory,
    // 配置管理
    loadProfiles,
    saveProfile,
    updateProfile,
    deleteProfile,
    setDefaultProfile,
    selectProfile,
    setConfig,
    testConnection,
    checkNetworkConnectivity
  }
})
