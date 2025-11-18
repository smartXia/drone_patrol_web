package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// WhipRequest WHIP请求结构
type WhipRequest struct {
	Room     string `form:"room" binding:"required"`
	TxSecret string // 临时字段，用于转发
	TxTime   string // 临时字段，用于转发
}

// WhipService WHIP服务
type WhipService struct {
	WhipServiceURL string
	PushDomain     string
	SecretKey      string
}

// NewWhipService 创建WHIP服务
func NewWhipService() *WhipService {
	return &WhipService{
		WhipServiceURL: "https://webrtcpush.tlivewebrtcpush.com/webrtc/v2/whip",
		PushDomain:     "214487.push.tlivecloud.com",
		SecretKey:      "6f51e7f3d1f79b60ded0e71207840c79066f8a721b90ba52802ecaab45305eab", // 与TRTC使用相同的密钥
	}
}

// GenerateTxSecret 生成txSecret
func (w *WhipService) GenerateTxSecret(room string) string {
	// 使用腾讯云WHIP的标准签名算法
	timestamp := time.Now().Unix()
	// 腾讯云WHIP签名格式: MD5(room + txTime + secretKey)
	data := fmt.Sprintf("%s%d%s", room, timestamp, w.SecretKey)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GenerateTxTime 生成txTime
func (w *WhipService) GenerateTxTime() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// GenerateAuthToken 生成认证token
func (w *WhipService) GenerateAuthToken(room string) string {
	txSecret := w.GenerateTxSecret(room)
	txTime := w.GenerateTxTime()
	// 腾讯云WHIP认证格式
	return fmt.Sprintf("webrtc://%s/live/%s?txSecret=%s&txTime=%s",
		w.PushDomain, room, txSecret, txTime)
}

// Whip 执行WHIP请求
func (w *WhipService) Whip(ctx context.Context, offer []byte, req WhipRequest) (*http.Response, error) {
	// 创建转发请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", w.WhipServiceURL, bytes.NewReader(offer))
	if err != nil {
		return nil, fmt.Errorf("creating request failed: %w", err)
	}

	// 使用腾讯云WHIP认证格式（参考代码）
	authToken := fmt.Sprintf("webrtc://%s/live/%s?txSecret=%s&txTime=%s",
		w.PushDomain, req.Room, req.TxSecret, req.TxTime)

	// 尝试不同的认证格式
	authToken2 := fmt.Sprintf("rtmp://%s/live/%s?txSecret=%s&txTime=%s",
		w.PushDomain, req.Room, req.TxSecret, req.TxTime)
	fmt.Printf("备用认证Token: %s\n", authToken2)

	// 调试输出认证信息
	fmt.Printf("WHIP认证信息 - Room: %s, TxSecret: %s, TxTime: %s\n",
		req.Room, req.TxSecret, req.TxTime)
	fmt.Printf("完整认证Token: %s\n", authToken)

	// 验证签名算法
	expectedSecret := w.GenerateTxSecret(req.Room)
	fmt.Printf("验证签名 - 期望: %s, 实际: %s\n", expectedSecret, req.TxSecret)

	// 设置请求头（参考代码格式）
	httpReq.Header.Set("Authorization", "Bearer "+authToken)
	httpReq.Header.Set("Content-Type", "application/sdp")

	// 发送请求（设置合理的超时）
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request failed: %w", err)
	}

	return resp, nil
}

// WhipHandler WHIP处理器（参考代码实现）
func (w *WhipService) WhipHandler(c *gin.Context) {
	var req WhipRequest

	// 1. 绑定查询参数
	fmt.Printf("WHIP请求URL: %s\n", c.Request.URL.String())
	fmt.Printf("WHIP请求查询参数: %+v\n", c.Request.URL.Query())
	fmt.Printf("WHIP请求Content-Type: %s\n", c.Request.Header.Get("Content-Type"))

	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Printf("WHIP查询参数绑定失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters", "details": err.Error()})
		return
	}

	fmt.Printf("WHIP请求参数解析成功: %+v\n", req)

	// 2. 读取客户端SDP offer
	offer, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read offer"})
		return
	}

	// 3. 自动生成认证信息
	txSecret := w.GenerateTxSecret(req.Room)
	txTime := w.GenerateTxTime()

	// 创建完整的认证请求
	authReq := WhipRequest{
		Room: req.Room,
	}

	// 临时设置认证信息用于转发
	authReq.TxSecret = txSecret
	authReq.TxTime = txTime

	fmt.Printf("自动生成认证信息 - TxSecret: %s, TxTime: %s\n", txSecret, txTime)

	// 4. 转发到目标服务器
	startTime := time.Now()
	resp, err := w.Whip(c.Request.Context(), offer, authReq)
	if err != nil {
		fmt.Printf("WHIP转发失败，耗时 %v: %v\n", time.Since(startTime), err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Forwarding failed",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 4. 读取目标服务器响应
	answer, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read target server response"})
		return
	}

	// 5. 将响应返回给客户端
	for k, v := range resp.Header {
		c.Header(k, v[0])
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), answer)
}
