import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import { generateClientId } from '../utils/mqtt'
import { API_BASE_URL, WS_BASE_URL, TIMEOUT } from '../config.js'

export const useMqttProxyStore = defineStore('mqtt-proxy', () => {
  // 状态
  const websocket = ref(null)
  const isConnected = ref(false)
  const isConnecting = ref(false)
  const connectionError = ref(null)
  const deviceStatus = ref({})
  const messageHistory = ref([])
  const subscribedTopics = ref(new Map()) // 存储 topic -> {qos, timestamp}
  const currentConfig = ref(null)
  const profiles = ref([])
  const currentProfileId = ref(null)

  // 后端本地服务地址
  const apiBase = API_BASE_URL

  // 计算属性
  const connectionStatus = computed(() => {
    if (isConnected.value) return 'connected'
    if (isConnecting.value) return 'connecting'
    return 'disconnected'
  })

  const visibleDeviceCount = computed(() => {
    return Object.keys(deviceStatus.value).length
  })

  // 当前MQTT配置名称
  const currentMQTTProfileName = computed(() => {
    if (currentProfileId.value && profiles.value.length > 0) {
      const profile = profiles.value.find(p => p.id === currentProfileId.value)
      return profile ? profile.name : '默认配置'
    }
    return '默认配置'
  })

  // 默认配置
  const defaultConfig = {
    protocol: 'tcp',
    host: '',
    port: 1883,
    path: '',
    username: '',
    password: '',
    clientId: generateClientId('web_client'),
    clean: true,
    reconnectPeriod: 0,
    connectTimeout: 30000,
    keepalive: 60,
    rejectUnauthorized: false
  }

  // 连接MQTT（通过后端代理）
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
      
      console.log('正在连接MQTT代理:', cfg.host, cfg.port)
      console.log('连接配置:', {
        host: cfg.host,
        port: cfg.port,
        protocol: cfg.protocol,
        username: cfg.username ? '***' : '未设置',
        clientId: cfg.clientId
      })

      // 如果已有WebSocket连接，先关闭
      if (websocket.value) {
        websocket.value.close()
        websocket.value = null
      }

      // 创建 WebSocket 连接到后端代理
      const wsUrl = `${WS_BASE_URL}/ws/mqtt`
      websocket.value = new WebSocket(wsUrl)

      websocket.value.onopen = () => {
        console.log('WebSocket 代理连接成功')
        
        // 检查 WebSocket 是否仍然有效
        if (!websocket.value) {
          console.error('WebSocket 在 onopen 时已被关闭')
          return
        }
        
        // 发送 MQTT 连接请求
        const connectMessage = {
          type: 'connect',
          config: {
            host: cfg.host,
            port: cfg.port,
            username: cfg.username,
            password: cfg.password,
            clientId: cfg.clientId
          }
        }
        
        try {
          websocket.value.send(JSON.stringify(connectMessage))
          console.log('MQTT 连接请求已发送')
        } catch (error) {
          console.error('发送 MQTT 连接请求失败:', error)
        }
      }

      websocket.value.onmessage = (event) => {
        if (!websocket.value) {
          console.warn('WebSocket 在 onmessage 时已被关闭')
          return
        }
        const message = JSON.parse(event.data)
        handleProxyMessage(message)
      }

      websocket.value.onclose = () => {
        console.log('WebSocket 代理连接关闭')
        isConnected.value = false
        isConnecting.value = false
        websocket.value = null
      }

      websocket.value.onerror = (error) => {
        console.error('WebSocket 代理连接错误:', error)
        connectionError.value = 'WebSocket 代理连接失败'
        isConnecting.value = false
        websocket.value = null
      }

    } catch (error) {
      console.error('MQTT连接失败:', error)
      connectionError.value = error.message
      isConnecting.value = false
    }
  }

  const connectWithConfig = async (config) => {
    if (isConnecting.value || isConnected.value) {
      console.log('MQTT已连接或正在连接中，跳过重复连接')
      return
    }

    try {
      isConnecting.value = true
      connectionError.value = null

      // 设置当前配置
      currentConfig.value = config
      
      console.log('正在使用指定配置连接MQTT代理:', config.host, config.port)
      console.log('连接配置:', {
        host: config.host,
        port: config.port,
        protocol: config.protocol,
        username: config.username ? '***' : '未设置',
        clientId: config.clientId
      })

      // 如果已有WebSocket连接，先关闭
      if (websocket.value) {
        websocket.value.close()
        websocket.value = null
      }

      // 创建 WebSocket 连接到后端代理
      const wsUrl = `${WS_BASE_URL}/ws/mqtt`
      websocket.value = new WebSocket(wsUrl)

      websocket.value.onopen = () => {
        console.log('WebSocket 代理连接成功')
        // 检查 WebSocket 是否仍然有效
        if (!websocket.value) {
          console.error('WebSocket 在 onopen 时已被关闭')
          return
        }
        // 发送连接命令到后端
        try {
          websocket.value.send(JSON.stringify({
            type: 'connect',
            config: config
          }))
          console.log('MQTT 连接请求已发送')
        } catch (error) {
          console.error('发送 MQTT 连接请求失败:', error)
        }
      }

      websocket.value.onmessage = (event) => {
        if (!websocket.value) {
          console.warn('WebSocket 在 onmessage 时已被关闭')
          return
        }
        const message = JSON.parse(event.data)
        handleProxyMessage(message)
      }

      websocket.value.onclose = () => {
        console.log('WebSocket 代理连接关闭')
        isConnected.value = false
        isConnecting.value = false
        websocket.value = null
      }

      websocket.value.onerror = (error) => {
        console.error('WebSocket 代理连接错误:', error)
        connectionError.value = 'WebSocket 代理连接失败'
        isConnecting.value = false
        websocket.value = null
      }

    } catch (error) {
      console.error('MQTT连接失败:', error)
      connectionError.value = error.message
      isConnecting.value = false
    }
  }

  // 处理代理消息
  const handleProxyMessage = (message) => {
    switch (message.type) {
      case 'connect_result':
        if (message.success) {
          console.log('MQTT连接成功')
          isConnected.value = true
          isConnecting.value = false
        } else {
          console.error('MQTT连接失败')
          connectionError.value = 'MQTT连接失败'
          isConnecting.value = false
        }
        break
        
      case 'mqtt_connected':
        console.log('MQTT连接成功')
        isConnected.value = true
        isConnecting.value = false
        break
        
      case 'mqtt_error':
        console.error('MQTT错误:', message.message)
        connectionError.value = message.message
        isConnecting.value = false
        break
        
      case 'mqtt_disconnected':
        console.log('MQTT连接断开')
        isConnected.value = false
        break
        
      case 'mqtt_message':
        handleMqttMessage(message.topic, message.payload, message.qos, message.retain)
        break
        
      case 'subscription_success':
        console.log('订阅成功:', message.topic)
        subscribedTopics.value.set(message.topic, {
          qos: message.qos || 0,
          timestamp: new Date().toISOString()
        })
        console.log('当前已订阅主题:', Array.from(subscribedTopics.value.keys()))
        break
        
      case 'subscribe_result':
        // 检查订阅结果，优先使用result返回码判断
        let subscribeSuccess = false
        let subscribeResultCode = null
        
        // 尝试从不同位置获取result字段
        if (typeof message.result === 'number') {
          subscribeResultCode = message.result
          subscribeSuccess = message.result === 0
        } else if (message.data && typeof message.data.result === 'number') {
          subscribeResultCode = message.data.result
          subscribeSuccess = message.data.result === 0
        } else if (message.success === true || message.payload === 'success') {
          // 兼容旧格式
          subscribeSuccess = true
          subscribeResultCode = 0
        } else {
          // 默认情况下，如果没有明确的失败标识，认为成功
          subscribeSuccess = true
          subscribeResultCode = 0
        }
        
        if (subscribeSuccess) {
          console.log('✅ 订阅成功:', message.topic, subscribeResultCode !== null ? `(result=${subscribeResultCode})` : '')
          subscribedTopics.value.set(message.topic, {
            qos: message.qos || 0,
            timestamp: new Date().toISOString()
          })
          console.log('当前已订阅主题:', Array.from(subscribedTopics.value.keys()))
        } else {
          console.error('❌ 订阅失败:', message.topic, subscribeResultCode !== null ? `(result=${subscribeResultCode})` : '')
        }
        break
        
      case 'publish_result':
        // 检查发布结果，优先使用result返回码判断
        let isSuccess = false
        let resultCode = null
        
        // 尝试从不同位置获取result字段
        if (typeof message.result === 'number') {
          resultCode = message.result
          isSuccess = message.result === 0
        } else if (message.data && typeof message.data.result === 'number') {
          resultCode = message.data.result
          isSuccess = message.data.result === 0
        } else if (message.success === true || message.payload === 'success') {
          // 兼容旧格式 - payload为'success'表示成功
          isSuccess = true
          resultCode = 0
        } else {
          // 默认情况下，如果没有明确的失败标识，认为成功
          isSuccess = true
          resultCode = 0
        }
        
        if (isSuccess) {
          console.log('✅ 发布成功:', message.topic, resultCode !== null ? `(result=${resultCode})` : '')
        } else {
          console.error('❌ 发布失败:', message.topic, resultCode !== null ? `(result=${resultCode})` : '')
          console.error('发布失败详情:', message)
          if (message.error) {
            console.error('发布错误信息:', message.error)
          }
        }
        break
    }
  }

  // 处理MQTT消息
  const handleMqttMessage = (topic, payload, qos, retain) => {
    // console.log('前端收到MQTT消息:', topic, payload)  
    const message = {
      id: Date.now() + Math.random(),
      topic,
      payload,
      qos,
      retain,
      timestamp: new Date().toISOString()
    }

    // 添加到消息历史
    messageHistory.value.unshift(message)
    // console.log('消息历史长度:', messageHistory.value.length)
    if (messageHistory.value.length > 1000) {
      messageHistory.value = messageHistory.value.slice(0, 1000)
    }

    // 解析设备状态
    try {
      const data = JSON.parse(payload)
      const deviceId = extractDeviceIdFromTopic(topic)
      
      if (deviceId) {
        deviceStatus.value[deviceId] = {
          ...deviceStatus.value[deviceId],
          ...data,
          lastUpdate: message.timestamp,
          topic
        }
      }
    } catch (e) {
      // 非JSON消息，忽略
    }
  }

  // 从主题提取设备ID
  const extractDeviceIdFromTopic = (topic) => {
    const parts = topic.split('/')
    return parts[1] || null
  }

  // 断开连接
  const disconnect = () => {
    if (websocket.value) {
      // 发送断开连接请求
      const disconnectMessage = {
        type: 'disconnect'
      }
      websocket.value.send(JSON.stringify(disconnectMessage))
      
      websocket.value.close()
      websocket.value = null
    }
    
    isConnected.value = false
    isConnecting.value = false
    connectionError.value = null
    subscribedTopics.value.clear()
  }

  // 订阅主题
  const subscribe = (topic, qos = 0) => {
    if (websocket.value && isConnected.value) {
      const subscribeMessage = {
        type: 'subscribe',
        topic,
        qos
      }
      websocket.value.send(JSON.stringify(subscribeMessage))
    }
  }

  // 订阅主题（兼容旧接口）
  const subscribeToTopics = async (topic, qos = 0) => {
    console.log('=== MQTT订阅请求 ===')
    console.log('Topic:', topic)
    console.log('QoS:', qos)
    console.log('连接状态:', isConnected.value)
    
    if (!isConnected.value) {
      console.error('❌ MQTT未连接，无法订阅')
      throw new Error('MQTT未连接')
    }
    
    if (websocket.value) {
      const subscribeMessage = {
        type: 'subscribe',
        topic,
        qos
      }
      
      console.log('发送订阅请求:', subscribeMessage)
      websocket.value.send(JSON.stringify(subscribeMessage))
      
      // 等待后端确认订阅成功
      return new Promise((resolve, reject) => {
        const timeout = setTimeout(() => {
          reject(new Error('订阅超时'))
        }, 5000)
        
        const handleMessage = (event) => {
          const message = JSON.parse(event.data)
          if (message.type === 'subscription_success' && message.topic === topic) {
            clearTimeout(timeout)
            websocket.value.removeEventListener('message', handleMessage)
            console.log('订阅确认成功:', topic)
            resolve()
          } else if (message.type === 'subscribe_result' && message.topic === topic) {
            clearTimeout(timeout)
            websocket.value.removeEventListener('message', handleMessage)
            
            // 检查订阅结果，优先使用result返回码判断
            let confirmSuccess = false
            if (typeof message.result === 'number') {
              confirmSuccess = message.result === 0
            } else if (message.data && typeof message.data.result === 'number') {
              confirmSuccess = message.data.result === 0
            } else if (message.success === true || message.payload === 'success') {
              confirmSuccess = true
            } else {
              // 默认情况下，如果没有明确的失败标识，认为成功
              confirmSuccess = true
            }
            
            if (confirmSuccess) {
              console.log('✅ 订阅确认成功:', topic)
              resolve()
            } else {
              console.error('❌ 订阅确认失败:', topic)
              reject(new Error('订阅失败'))
            }
          }
        }
        
        websocket.value.addEventListener('message', handleMessage)
      })
    } else {
      throw new Error('WebSocket连接未建立')
    }
  }

  // 发布消息
  const publish = (topic, payload, qos = 0, retain = false) => {
    if (!websocket.value) {
      console.error('WebSocket连接不存在，无法发布消息:', topic)
      return
    }
    
    if (!isConnected.value) {
      console.error('MQTT未连接，无法发布消息:', topic)
      return
    }
    
    try {
      const publishMessage = {
        type: 'publish',
        topic,
        payload,
        qos,
        retain
      }
      
      console.log('准备发布消息:', {
        topic,
        payloadLength: payload?.length || 0,
        qos,
        retain
      })
      
      websocket.value.send(JSON.stringify(publishMessage))
    } catch (error) {
      console.error('发布消息时发生错误:', error, 'Topic:', topic)
    }
  }

  // 发布消息（兼容旧接口）
  const publishMessage = async (topic, payload, qos = 0, retain = false) => {
    if (!isConnected.value) {
      throw new Error('MQTT未连接')
    }
    
    if (websocket.value) {
      const publishMessage = {
        type: 'publish',
        topic,
        payload,
        qos,
        retain
      }
      websocket.value.send(JSON.stringify(publishMessage))
      
      return Promise.resolve()
    } else {
      throw new Error('WebSocket连接未建立')
    }
  }

  // 取消订阅
  const unsubscribe = (topic) => {
    subscribedTopics.value.delete(topic)
    // 注意：这里需要后端支持取消订阅功能
  }

  // 取消订阅主题（兼容旧接口）
  const unsubscribeTopics = async (topic) => {
    if (!isConnected.value) {
      throw new Error('MQTT未连接')
    }
    
    if (websocket.value) {
      const unsubscribeMessage = {
        type: 'unsubscribe',
        topic
      }
      websocket.value.send(JSON.stringify(unsubscribeMessage))
      
      // 从已订阅主题列表中移除
      subscribedTopics.value.delete(topic)
      
      return Promise.resolve()
    } else {
      throw new Error('WebSocket连接未建立')
    }
  }

  // 清空消息历史
  const clearMessageHistory = () => {
    messageHistory.value = []
  }

  // 清空设备状态
  const clearDeviceStatus = () => {
    deviceStatus.value = {}
  }

  // 重置连接错误
  const resetConnectionError = () => {
    connectionError.value = null
  }

  // ===== 配置管理 =====
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
    }
  }

  const updateProfile = async (id, payload) => {
    const { data } = await axios.put(`${apiBase}/api/mqtt/profiles/${id}`, payload)
    if (data && data.code === 0) {
      await loadProfiles()
    }
  }

  const deleteProfile = async (id) => {
    const { data } = await axios.delete(`${apiBase}/api/mqtt/profiles/${id}`)
    if (data && data.code === 0) {
      await loadProfiles()
    }
  }

  const setDefaultProfile = async (id) => {
    const { data } = await axios.post(`${apiBase}/api/mqtt/profiles/${id}/default`)
    if (data && data.code === 0) {
      await loadProfiles()
    }
  }

  const selectProfile = async (id) => {
    const { data } = await axios.get(`${apiBase}/api/mqtt/profiles/${id}`)
    if (data && data.code === 0) {
      currentConfig.value = data.data.config
    }
    return data
  }

  // 测试连接
  const testConnection = async (config) => {
    const { data } = await axios.post(`${apiBase}/api/mqtt/test`, config)
    return data
  }

  // 网络连通性检查
  const checkNetworkConnectivity = async (host, port) => {
    try {
      const { data } = await axios.get(`${apiBase}/api/network/ping`, {
        params: { host, port }
      })
      
      if (data && data.code === 0) {
        return {
          success: data.data.success,
          error: data.data.error,
          duration: data.data.duration
        }
      } else {
        return {
          success: false,
          error: data?.message || '网络检查失败',
          duration: 0
        }
      }
    } catch (error) {
      return {
        success: false,
        error: error.message,
        duration: 0
      }
    }
  }

  // ===== DJI 服务调用辅助方法 =====
  const genId = (prefix = 'id') => `${prefix}_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`

  /**
   * 发送 DJI 服务调用
   * @param {string} gatewaySn 网关SN（机场SN）
   * @param {string} method 方法名，例如 'fileupload_list'
   * @param {object} data 负载数据
   * @param {object} options { qos?: number, retain?: boolean, autoSubscribeReply?: boolean }
   * @returns {{tid:string,bid:string,topic:string,replyTopic:string,payload:object}}
   */
  const sendService = async (gatewaySn, method, data = {}, options = {}) => {
    if (!isConnected.value) {
      throw new Error('MQTT未连接')
    }
    if (!gatewaySn) {
      throw new Error('缺少 gatewaySn')
    }
    const qos = options.qos ?? 0
    const retain = options.retain ?? false
    const autoSubscribeReply = options.autoSubscribeReply ?? true

    const topic = `thing/product/${gatewaySn}/services`
    const replyTopic = `thing/product/${gatewaySn}/services_reply`
    const tid = genId('tid')
    const bid = genId('bid')
    const payload = {
      tid,
      bid,
      timestamp: Date.now(),
      method,
      data
    }

    if (autoSubscribeReply) {
      try { 
        await subscribeToTopics(replyTopic, qos)
        console.log('已订阅回复Topic:', replyTopic)
      } catch (error) { 
        console.error('订阅回复Topic失败:', replyTopic, error)
      }
    }

    console.log('发送服务调用:', {
      method,
      topic,
      replyTopic,
      tid,
      bid,
      dataKeys: Object.keys(data)
    })

    publish(topic, JSON.stringify(payload), qos, retain)
    return { tid, bid, topic, replyTopic, payload }
  }

  /**
   * 获取设备可上传的文件列表（fileupload_list）
   * @param {string} gatewaySn 网关SN（机场SN）
   * @param {Array<string|number>} moduleList 模块过滤列表，如 ['0','3'] 或 [0,3]
   * @param {object} options 见 sendService
   */
  const requestFileUploadList = async (gatewaySn, moduleList = [], options = {}) => {
    const modules = Array.isArray(moduleList) ? moduleList.map(v => String(v)) : []
    return await sendService(gatewaySn, 'fileupload_list', { module_list: modules }, options)
  }

  /**
   * 请求文件上传开始
   * @param {string} gatewaySn 网关SN
   * @param {object} uploadData 上传数据
   * @param {object} options 见 sendService
   */
  const requestFileUploadStart = async (gatewaySn, uploadData, options = {}) => {
    return await sendService(gatewaySn, 'fileupload_start', uploadData, options)
  }

  /**
   * 请求文件上传状态更新
   * @param {string} gatewaySn 网关SN
   * @param {object} statusData 状态数据
   * @param {object} options 见 sendService
   */
  const requestFileUploadUpdate = async (gatewaySn, statusData, options = {}) => {
    return await sendService(gatewaySn, 'fileupload_update', statusData, options)
  }

  return {
    // 状态
    websocket,
    isConnected,
    isConnecting,
    connectionError,
    deviceStatus,
    messageHistory,
    subscribedTopics,
    currentConfig,
    profiles,
    currentProfileId,
    
    // 计算属性
    connectionStatus,
    visibleDeviceCount,
    currentMQTTProfileName,
    
    // 连接方法
    connect,
    connectWithConfig,
    disconnect,
    subscribe,
    publish,
    unsubscribe,
    
    // 兼容旧接口的方法
    subscribeToTopics,
    unsubscribeTopics,
    publishMessage,
    
    // 工具方法
    clearMessageHistory,
    clearDeviceStatus,
    resetConnectionError,
    
    // 配置管理
    loadProfiles,
    saveProfile,
    updateProfile,
    deleteProfile,
    setDefaultProfile,
    selectProfile,
    testConnection,
    checkNetworkConnectivity
    ,
    // DJI 服务调用
    sendService,
    requestFileUploadList,
    requestFileUploadStart,
    requestFileUploadUpdate
  }
})
