// WHIP WebRTC推流服务
export class WhipService {
  constructor(config) {
    this.apiBaseUrl = config.apiBaseUrl
    this.room = null
    this.pc = null
    this.stream = null
  }

  // 获取WHIP认证信息
  async getWhipAuth(room) {
    try {
      const response = await fetch(`${this.apiBaseUrl}/api/whip/auth/${room}`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error('获取WHIP认证信息失败:', error)
      throw error
    }
  }

  // 开始WebRTC推流（从大疆机场设备）
  async startStream(room, videoElement) {
    try {
      console.log('开始WHIP推流，房间号:', room)

      // 创建RTCPeerConnection
      this.pc = new RTCPeerConnection({
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' },
          { urls: 'stun:stun1.l.google.com:19302' }
        ]
      })

      // 注意：这里不再获取本地摄像头，而是从大疆机场设备获取视频流
      // 大疆机场的视频流将通过WHIP协议从设备端推送到腾讯云
      // 然后通过WebRTC接收并显示
      
      // 设置视频元素以接收来自大疆机场的流
      if (videoElement) {
        // 视频元素将显示从大疆机场设备传来的视频流
        videoElement.srcObject = null // 初始为空，等待远程流
      }

      // 创建offer - 接收大疆机场的视频流
      const offer = await this.pc.createOffer({
        offerToReceiveAudio: true,
        offerToReceiveVideo: true
      })

      await this.pc.setLocalDescription(offer)

      // 发送WHIP请求 - 只传递房间号，后端自动生成认证信息
      const whipUrl = `${this.apiBaseUrl}/api/whip/stream?room=${room}`
      
      console.log('WHIP请求URL:', whipUrl)
      console.log('WHIP请求参数:', { room })
      
      const response = await fetch(whipUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/sdp'
        },
        body: offer.sdp
      })

      if (!response.ok) {
        throw new Error(`WHIP请求失败: ${response.status}`)
      }

      // 获取answer
      const answerSdp = await response.text()
      const answer = new RTCSessionDescription({
        type: 'answer',
        sdp: answerSdp
      })

      await this.pc.setRemoteDescription(answer)

      // 监听远程流（来自大疆机场设备）
      this.pc.ontrack = (event) => {
        console.log('接收到大疆机场视频流:', event.streams[0])
        if (videoElement && event.streams[0]) {
          videoElement.srcObject = event.streams[0]
        }
        this.stream = event.streams[0]
      }

      // 监听连接状态
      this.pc.onconnectionstatechange = () => {
        console.log('WebRTC连接状态:', this.pc.connectionState)
      }

      this.room = room
      console.log('WHIP接收大疆机场视频流开始成功')
      
      return {
        success: true,
        room: room,
        stream: this.stream
      }

    } catch (error) {
      console.error('WHIP推流失败:', error)
      throw error
    }
  }

  // 停止接收流
  async stopStream() {
    try {
      // 停止接收大疆机场视频流
      if (this.pc) {
        this.pc.close()
        this.pc = null
      }

      if (this.stream) {
        this.stream = null
      }

      this.room = null
      console.log('WHIP接收大疆机场视频流已停止')
      
      return { success: true }
    } catch (error) {
      console.error('停止WHIP接收流失败:', error)
      throw error
    }
  }

  // 获取推流状态
  getStreamStatus() {
    return {
      isStreaming: this.pc && this.pc.connectionState === 'connected',
      room: this.room,
      connectionState: this.pc ? this.pc.connectionState : 'disconnected'
    }
  }
}
