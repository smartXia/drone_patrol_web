package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/tencentyun/tls-sig-api-v2-golang/tencentyun"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// TencentLiveConfig 腾讯云直播配置
type TencentLiveConfig struct {
	SDKAppID  int64  `json:"sdk_app_id"`
	SecretKey string `json:"secret_key"`
	LiveURL   string `json:"live_url"`
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

// TencentLiveService 腾讯云直播服务
type TencentLiveService struct {
	config *TencentLiveConfig
}

// NewTencentLiveService 创建腾讯云直播服务实例
func NewTencentLiveService() *TencentLiveService {
	return &TencentLiveService{
		config: &TencentLiveConfig{
			SDKAppID:  1600093496,                                                         // 您的SDK App ID
			SecretKey: "ba4030ba7f43d949c4d085912d2777d4f43bc5d01425b38c4dee67658d04171b", // 您的Secret Key
			LiveURL:   "rtmp://rtmp.rtc.qq.com/push/",
		},
	}
}

// GenerateUserSig 生成用户签名
func (s *TencentLiveService) GenerateUserSig(userID string) (string, error) {
	// 使用腾讯云TLS签名API生成UserSig
	userSig, err := tencentyun.GenUserSig(1600093496, s.config.SecretKey, userID, 86400*180) // 180天有效期
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

// GeneratePlayURL 生成播放地址
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

	// 生成推流地址
	pushURL := s.GeneratePushURL(streamID)

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
