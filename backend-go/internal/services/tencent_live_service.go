package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/tencentyun/tls-sig-api-v2-golang/tencentyun"

	"github.com/gin-gonic/gin"
)

// TencentLiveConfig 腾讯云直播配置
type TencentLiveConfig struct {
	SDKAppID  int64  `json:"sdk_app_id"`
	SecretKey string `json:"secret_key"`
	LiveURL   string `json:"live_url"`
	TRTCAppID int64  `json:"trtc_app_id"` // TRTC应用ID
}

// LiveStreamRequest 直播流请求
type LiveStreamRequest struct {
	DeviceSN   string `json:"device_sn" binding:"required"`
	StreamType string `json:"stream_type"` // "airport" or "aircraft"
	Resolution string `json:"resolution"`  // "1280x720", "1920x1080"
	Bitrate    int    `json:"bitrate"`     // 码率
	FPS        int    `json:"fps"`         // 帧率
	Quality    string `json:"quality"`     // "low", "medium", "high"
}

// LiveStreamType 直播类型常量
const (
	LiveStreamTypeAirport  = "airport"  // 机场直播
	LiveStreamTypeAircraft = "aircraft" // 无人机直播
)

// LiveStreamResponse 直播流响应
type LiveStreamResponse struct {
	Success  bool   `json:"success"`
	StreamID string `json:"stream_id"`
	PushURL  string `json:"push_url"`
	PlayURL  string `json:"play_url"`
	UserSig  string `json:"user_sig"`
	SDKAppID int64  `json:"sdk_app_id"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

// TRTCRoomRequest TRTC房间请求
type TRTCRoomRequest struct {
	DeviceSN string `json:"device_sn" binding:"required"` // 设备SN作为房间号
	UserID   string `json:"user_id"`                      // 用户ID
}

// TRTCRoomResponse TRTC房间响应
type TRTCRoomResponse struct {
	Success  bool   `json:"success"`
	RoomID   string `json:"room_id"`    // 房间号（设备SN）
	UserID   string `json:"user_id"`    // 用户ID
	UserSig  string `json:"user_sig"`   // 用户签名
	SDKAppID int64  `json:"sdk_app_id"` // SDK应用ID
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

// TencentLiveService 腾讯云直播服务
type TencentLiveService struct {
	config *TencentLiveConfig
}

// NewTencentLiveService 创建腾讯云直播服务实例
func NewTencentLiveService() *TencentLiveService {
	return &TencentLiveService{
		config: &TencentLiveConfig{
			SDKAppID:  1600108347,                                                         // DJI设备使用的SDK App ID
			SecretKey: "6f51e7f3d1f79b60ded0e71207840c79066f8a721b90ba52802ecaab45305eab", // 您的Secret Key
			LiveURL:   "rtmp://rtmp.rtc.qq.com/push/",
			TRTCAppID: 1600108347, // TRTC应用ID（与SDK App ID相同）
		},
	}
}

// GenerateUserSig 生成用户签名
func (s *TencentLiveService) GenerateUserSig(userID string) (string, error) {
	// 使用腾讯云TLS签名API生成UserSig
	userSig, err := tencentyun.GenUserSig(int(s.config.SDKAppID), s.config.SecretKey, userID, 86400*180) // 180天有效期
	if err != nil {
		return "", fmt.Errorf("生成UserSig失败: %v", err)
	}
	return userSig, nil
}

// GenerateStreamID 生成流ID
func (s *TencentLiveService) GenerateStreamID(deviceSN string, streamType string) string {
	timestamp := time.Now().Unix()
	random := rand.Intn(10000)
	return fmt.Sprintf("%s_%s_%d_%d", deviceSN, streamType, timestamp, random)
}

// GeneratePushURL 生成推流地址
func (s *TencentLiveService) GeneratePushURL(streamID string) string {
	// 生成推流URL，包含签名
	timestamp := time.Now().Unix() + 86400 // 24小时有效期
	txTime := fmt.Sprintf("%X", timestamp)

	// 生成签名
	signature := s.generateSignature(streamID, txTime)

	return fmt.Sprintf("%s%s?txSecret=%s&txTime=%s", s.config.LiveURL, streamID, signature, txTime)
}

// GenerateDJIPushURL 生成DJI设备格式的推流地址
func (s *TencentLiveService) GenerateDJIPushURL(deviceSN string, userSig string) string {
	// DJI设备格式: rtmp://rtmp.rtc.qq.com/push/{sn}?sdkappid={sdkappid}&userid={userid}&usersig={usersig}
	return fmt.Sprintf("%s%s?sdkappid=%d&userid=%s&usersig=%s",
		s.config.LiveURL, deviceSN, s.config.SDKAppID, deviceSN, userSig)
}

// GeneratePlayURL 生成播放地址7870
func (s *TencentLiveService) GeneratePlayURL(streamID string) string {
	// 生成播放地址（这里使用腾讯云直播的播放地址格式）
	return fmt.Sprintf("https://live.rtc.qq.com/live/%s.flv", streamID)
}

// generateSignature 生成签名
func (s *TencentLiveService) generateSignature(streamID, txTime string) string {
	// 生成签名字符串
	signString := fmt.Sprintf("%s%s%s", s.config.SecretKey, streamID, txTime)

	// 计算HMAC-SHA256
	h := hmac.New(sha256.New, []byte(s.config.SecretKey))
	h.Write([]byte(signString))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

// CreateLiveStream 创建直播流
func (s *TencentLiveService) CreateLiveStream(req *LiveStreamRequest) (*LiveStreamResponse, error) {
	// 生成流ID
	streamID := s.GenerateStreamID(req.DeviceSN, req.StreamType)

	// 生成用户签名
	userSig, err := s.GenerateUserSig(req.DeviceSN)
	if err != nil {
		return &LiveStreamResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 统一使用DJI设备格式的推流地址
	pushURL := s.GenerateDJIPushURL(req.DeviceSN, userSig)

	// 生成播放地址
	playURL := s.GeneratePlayURL(streamID)

	return &LiveStreamResponse{
		Success:  true,
		StreamID: streamID,
		PushURL:  pushURL,
		PlayURL:  playURL,
		UserSig:  userSig,
		SDKAppID: s.config.SDKAppID,
		Message:  "直播流创建成功",
	}, nil
}

// GetLiveStreamStatus 获取直播流状态
func (s *TencentLiveService) GetLiveStreamStatus(streamID string) (*LiveStreamResponse, error) {
	// 这里可以调用腾讯云API查询直播流状态
	// 暂时返回模拟数据
	return &LiveStreamResponse{
		Success:  true,
		StreamID: streamID,
		Message:  "直播流状态正常",
	}, nil
}

// StopLiveStream 停止直播流
func (s *TencentLiveService) StopLiveStream(streamID string) (*LiveStreamResponse, error) {
	// 这里可以调用腾讯云API停止直播流
	// 暂时返回模拟数据
	return &LiveStreamResponse{
		Success:  true,
		StreamID: streamID,
		Message:  "直播流已停止",
	}, nil
}

// CreateLiveStreamHandler 创建直播流处理器
func CreateLiveStreamHandler(c *gin.Context) {
	var req LiveStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建腾讯云直播服务
	service := NewTencentLiveService()

	// 创建直播流
	response, err := service.CreateLiveStream(&req)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "创建直播流失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, response)
}

// GetLiveStreamStatusHandler 获取直播流状态处理器
func GetLiveStreamStatusHandler(c *gin.Context) {
	streamID := c.Param("streamId")
	if streamID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "流ID不能为空",
		})
		return
	}

	service := NewTencentLiveService()
	response, err := service.GetLiveStreamStatus(streamID)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "获取直播流状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, response)
}

// StopLiveStreamHandler 停止直播流处理器
func StopLiveStreamHandler(c *gin.Context) {
	streamID := c.Param("streamId")
	if streamID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "流ID不能为空",
		})
		return
	}

	service := NewTencentLiveService()
	response, err := service.StopLiveStream(streamID)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "停止直播流失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, response)
}

// CreateTRTCRoom 创建TRTC房间
func (s *TencentLiveService) CreateTRTCRoom(req *TRTCRoomRequest) (*TRTCRoomResponse, error) {
	// 使用设备SN作为房间号
	roomID := req.DeviceSN

	// 生成用户签名
	userSig, err := s.GenerateUserSig(req.UserID)
	if err != nil {
		return &TRTCRoomResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &TRTCRoomResponse{
		Success:  true,
		RoomID:   roomID,
		UserID:   req.UserID,
		UserSig:  userSig,
		SDKAppID: s.config.TRTCAppID,
		Message:  "TRTC房间创建成功",
	}, nil
}

// CreateTRTCRoomHandler 创建TRTC房间处理器
func CreateTRTCRoomHandler(c *gin.Context) {
	var req TRTCRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建腾讯云直播服务
	service := NewTencentLiveService()

	// 创建TRTC房间
	response, err := service.CreateTRTCRoom(&req)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "创建TRTC房间失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, response)
}

// JoinTRTCRoomHandler 加入TRTC房间处理器
func JoinTRTCRoomHandler(c *gin.Context) {
	var req TRTCRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建腾讯云直播服务
	service := NewTencentLiveService()

	// 加入TRTC房间
	response, err := service.CreateTRTCRoom(&req) // 使用相同的逻辑，因为都是生成用户签名
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "加入TRTC房间失败: " + err.Error(),
		})
		return
	}

	// 修改消息为加入房间
	response.Message = "用户已加入TRTC房间"

	c.JSON(200, response)
}
