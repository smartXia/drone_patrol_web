/**
 * DJI Dock3直播控制服务
 * 通过MQTT控制DJI机场摄像头直播功能
 */
import { useMqttProxyStore } from '@/stores/mqtt-proxy'

export class DJILiveService {
  constructor() {
    this.mqttStore = useMqttProxyStore()
  }

  /**
   * 开始DJI机场直播
   * @param {string} gatewaySn 机场SN
   * @param {object} liveConfig 直播配置
   * @returns {Promise<object>} 操作结果
   */
  async startDockLive(gatewaySn, liveConfig = {}) {
    try {
      if (!this.mqttStore.isConnected) {
        throw new Error('MQTT未连接')
      }

      const defaultConfig = {
        live_url: liveConfig.liveUrl || 'rtmp://rtmp.rtc.qq.com/push/',
        resolution: liveConfig.resolution || '1920x1080',
        bitrate: liveConfig.bitrate || 2000,
        fps: liveConfig.fps || 25,
        quality: liveConfig.quality || 'high'
      }

      console.log('开始DJI机场直播:', { gatewaySn, config: defaultConfig })

      // 发送DJI直播开始服务调用
      const result = await this.mqttStore.sendService(gatewaySn, 'live_start', defaultConfig)
      
      return {
        success: true,
        tid: result.tid,
        bid: result.bid,
        message: 'DJI直播开始命令已发送'
      }
    } catch (error) {
      console.error('开始DJI机场直播失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  }

  /**
   * 停止DJI机场直播
   * @param {string} gatewaySn 机场SN
   * @returns {Promise<object>} 操作结果
   */
  async stopDockLive(gatewaySn) {
    try {
      if (!this.mqttStore.isConnected) {
        throw new Error('MQTT未连接')
      }

      console.log('停止DJI机场直播:', gatewaySn)

      // 发送DJI直播停止服务调用
      const result = await this.mqttStore.sendService(gatewaySn, 'live_stop', {})
      
      return {
        success: true,
        tid: result.tid,
        bid: result.bid,
        message: 'DJI直播停止命令已发送'
      }
    } catch (error) {
      console.error('停止DJI机场直播失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  }

  /**
   * 获取DJI直播状态
   * @param {string} gatewaySn 机场SN
   * @returns {Promise<object>} 直播状态
   */
  async getDockLiveStatus(gatewaySn) {
    try {
      if (!this.mqttStore.isConnected) {
        throw new Error('MQTT未连接')
      }

      console.log('获取DJI直播状态:', gatewaySn)

      // 发送DJI直播状态查询服务调用
      const result = await this.mqttStore.sendService(gatewaySn, 'live_status', {})
      
      return {
        success: true,
        tid: result.tid,
        bid: result.bid,
        message: 'DJI直播状态查询已发送'
      }
    } catch (error) {
      console.error('获取DJI直播状态失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  }

  /**
   * 设置DJI直播参数
   * @param {string} gatewaySn 机场SN
   * @param {object} params 直播参数
   * @returns {Promise<object>} 操作结果
   */
  async setDockLiveParams(gatewaySn, params = {}) {
    try {
      if (!this.mqttStore.isConnected) {
        throw new Error('MQTT未连接')
      }

      const defaultParams = {
        resolution: params.resolution || '1920x1080',
        bitrate: params.bitrate || 2000,
        fps: params.fps || 25,
        quality: params.quality || 'high',
        audio_enabled: params.audioEnabled || false
      }

      console.log('设置DJI直播参数:', { gatewaySn, params: defaultParams })

      // 发送DJI直播参数设置服务调用
      const result = await this.mqttStore.sendService(gatewaySn, 'live_set_params', defaultParams)
      
      return {
        success: true,
        tid: result.tid,
        bid: result.bid,
        message: 'DJI直播参数设置已发送'
      }
    } catch (error) {
      console.error('设置DJI直播参数失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  }

  /**
   * 开始DJI飞机直播
   * @param {string} aircraftSn 飞机SN
   * @param {object} liveConfig 直播配置
   * @returns {Promise<object>} 操作结果
   */
  async startAircraftLive(aircraftSn, liveConfig = {}) {
    try {
      if (!this.mqttStore.isConnected) {
        throw new Error('MQTT未连接')
      }

      const defaultConfig = {
        live_url: liveConfig.liveUrl || 'rtmp://rtmp.rtc.qq.com/push/',
        resolution: liveConfig.resolution || '1920x1080',
        bitrate: liveConfig.bitrate || 2000,
        fps: liveConfig.fps || 25,
        quality: liveConfig.quality || 'high'
      }

      console.log('开始DJI飞机直播:', { aircraftSn, config: defaultConfig })

      // 发送DJI飞机直播开始服务调用
      const result = await this.mqttStore.sendService(aircraftSn, 'live_start', defaultConfig)
      
      return {
        success: true,
        tid: result.tid,
        bid: result.bid,
        message: 'DJI飞机直播开始命令已发送'
      }
    } catch (error) {
      console.error('开始DJI飞机直播失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  }

  /**
   * 停止DJI飞机直播
   * @param {string} aircraftSn 飞机SN
   * @returns {Promise<object>} 操作结果
   */
  async stopAircraftLive(aircraftSn) {
    try {
      if (!this.mqttStore.isConnected) {
        throw new Error('MQTT未连接')
      }

      console.log('停止DJI飞机直播:', aircraftSn)

      // 发送DJI飞机直播停止服务调用
      const result = await this.mqttStore.sendService(aircraftSn, 'live_stop', {})
      
      return {
        success: true,
        tid: result.tid,
        bid: result.bid,
        message: 'DJI飞机直播停止命令已发送'
      }
    } catch (error) {
      console.error('停止DJI飞机直播失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  }

  /**
   * 获取DJI飞机直播状态
   * @param {string} aircraftSn 飞机SN
   * @returns {Promise<object>} 直播状态
   */
  async getAircraftLiveStatus(aircraftSn) {
    try {
      if (!this.mqttStore.isConnected) {
        throw new Error('MQTT未连接')
      }

      console.log('获取DJI飞机直播状态:', aircraftSn)

      // 发送DJI飞机直播状态查询服务调用
      const result = await this.mqttStore.sendService(aircraftSn, 'live_status', {})
      
      return {
        success: true,
        tid: result.tid,
        bid: result.bid,
        message: 'DJI飞机直播状态查询已发送'
      }
    } catch (error) {
      console.error('获取DJI飞机直播状态失败:', error)
      return {
        success: false,
        error: error.message
      }
    }
  }

  /**
   * 监听DJI直播状态变化
   * @param {string} deviceSn 设备SN（机场或飞机）
   * @param {function} callback 状态变化回调
   */
  subscribeToLiveStatus(deviceSn, callback) {
    const statusTopic = `thing/product/${deviceSn}/live_status`
    
    // 订阅直播状态主题
    this.mqttStore.subscribeToTopics(statusTopic, 1)
    
    // 监听消息
    this.mqttStore.onMessage((topic, message) => {
      if (topic === statusTopic) {
        try {
          const data = JSON.parse(message)
          callback(data)
        } catch (error) {
          console.error('解析DJI直播状态消息失败:', error)
        }
      }
    })
  }
}

export default DJILiveService
