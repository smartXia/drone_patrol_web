package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"drone-patrol-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该更严格
	},
}

// WebSocketHandler 处理WebSocket连接
func (h *Handlers) WebSocketHandler(c *gin.Context) {
	log.Printf("WebSocket connection attempt from %s", c.Request.RemoteAddr)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer conn.Close()

	log.Println("WebSocket client connected successfully")

	// 生成客户端ID
	clientID := generateClientID()

	// 添加到MQTT代理服务
	h.MQTTProxy.AddClient(clientID, conn)
	defer h.MQTTProxy.RemoveClient(clientID)

	// 发送欢迎消息
	welcomeMsg := map[string]interface{}{
		"type":      "welcome",
		"message":   "WebSocket connection established",
		"clientId":  clientID,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Failed to send welcome message: %v", err)
	}

	// 处理WebSocket消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		if messageType == websocket.TextMessage {
			var wsMessage services.WebSocketMessage
			if err := json.Unmarshal(message, &wsMessage); err != nil {
				log.Printf("Failed to unmarshal WebSocket message: %v", err)
				continue
			}

			log.Printf("Received WebSocket message: %+v", wsMessage)

			// 处理消息
			if err := h.MQTTProxy.HandleWebSocketMessage(clientID, wsMessage); err != nil {
				log.Printf("Failed to handle WebSocket message: %v", err)

				// 发送错误消息
				errorMsg := services.WebSocketMessage{
					Type:    "error",
					Payload: err.Error(),
				}
				if err := conn.WriteJSON(errorMsg); err != nil {
					log.Printf("Failed to send error message: %v", err)
				}
			}
		}
	}

	log.Println("WebSocket client disconnected")
}

// generateClientID 生成客户端ID
func generateClientID() string {
	return "client_" + time.Now().Format("20060102150405") + "_" + randomString(6)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
