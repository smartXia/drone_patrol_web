package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

// MQTTProxyService MQTT代理服务
type MQTTProxyService struct {
	clients map[string]*MQTTClient
	mutex   sync.RWMutex
}

// MQTTClient MQTT客户端包装
type MQTTClient struct {
	ID          string
	Client      mqtt.Client
	WSConn      *websocket.Conn
	Config      MQTTConfig
	IsConnected bool
	mutex       sync.RWMutex
}

// MQTTConfig MQTT配置
type MQTTConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"clientId"`
}

// WebSocketMessage WebSocket消息
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Topic   string      `json:"topic,omitempty"`
	Payload string      `json:"payload,omitempty"`
	QoS     int         `json:"qos,omitempty"`
	Retain  bool        `json:"retain,omitempty"`
	Config  *MQTTConfig `json:"config,omitempty"`
}

// NewMQTTProxyService 创建MQTT代理服务
func NewMQTTProxyService() *MQTTProxyService {
	return &MQTTProxyService{
		clients: make(map[string]*MQTTClient),
	}
}

// AddClient 添加客户端
func (s *MQTTProxyService) AddClient(clientID string, wsConn *websocket.Conn) *MQTTClient {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &MQTTClient{
		ID:          clientID,
		WSConn:      wsConn,
		IsConnected: false,
	}

	s.clients[clientID] = client
	return client
}

// RemoveClient 移除客户端
func (s *MQTTProxyService) RemoveClient(clientID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[clientID]; exists {
		if client.Client != nil && client.Client.IsConnected() {
			client.Client.Disconnect(250)
		}
		delete(s.clients, clientID)
	}
}

// GetClient 获取客户端
func (s *MQTTProxyService) GetClient(clientID string) (*MQTTClient, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	client, exists := s.clients[clientID]
	return client, exists
}

// HandleWebSocketMessage 处理WebSocket消息
func (s *MQTTProxyService) HandleWebSocketMessage(clientID string, message WebSocketMessage) error {
	client, exists := s.GetClient(clientID)
	if !exists {
		return fmt.Errorf("client not found: %s", clientID)
	}

	switch message.Type {
	case "connect":
		return s.handleConnect(client, message.Config)
	case "disconnect":
		return s.handleDisconnect(client)
	case "subscribe":
		return s.handleSubscribe(client, message.Topic, message.QoS)
	case "unsubscribe":
		return s.handleUnsubscribe(client, message.Topic)
	case "publish":
		return s.handlePublish(client, message.Topic, message.Payload, message.QoS, message.Retain)
	default:
		return fmt.Errorf("unknown message type: %s", message.Type)
	}
}

// handleConnect 处理连接
func (s *MQTTProxyService) handleConnect(client *MQTTClient, config *MQTTConfig) error {
	if config == nil {
		return fmt.Errorf("config is required for connect")
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	// 如果已经连接，先断开
	if client.Client != nil && client.Client.IsConnected() {
		client.Client.Disconnect(250)
	}

	client.Config = *config

	// 创建MQTT客户端选项
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Host, config.Port))
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(false)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetKeepAlive(60 * time.Second)

	// 设置连接处理器
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Printf("MQTT client %s connected to %s:%d", client.ID, config.Host, config.Port)
		client.IsConnected = true
		s.sendWebSocketMessage(client, WebSocketMessage{
			Type: "mqtt_connected",
		})
	})

	// 设置连接丢失处理器
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Printf("MQTT client %s connection lost: %v", client.ID, err)
		client.IsConnected = false
		s.sendWebSocketMessage(client, WebSocketMessage{
			Type: "mqtt_disconnected",
		})
	})

	// 创建客户端
	client.Client = mqtt.NewClient(opts)

	// 连接
	if token := client.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("MQTT client %s connection failed: %v", client.ID, token.Error())
		s.sendWebSocketMessage(client, WebSocketMessage{
			Type:    "mqtt_error",
			Payload: token.Error().Error(),
		})
		return token.Error()
	}

	return nil
}

// handleDisconnect 处理断开连接
func (s *MQTTProxyService) handleDisconnect(client *MQTTClient) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	if client.Client != nil && client.Client.IsConnected() {
		client.Client.Disconnect(250)
	}

	client.IsConnected = false
	return nil
}

// handleSubscribe 处理订阅
func (s *MQTTProxyService) handleSubscribe(client *MQTTClient, topic string, qos int) error {
	if !client.IsConnected || client.Client == nil {
		return fmt.Errorf("MQTT client not connected")
	}

	if token := client.Client.Subscribe(topic, byte(qos), func(c mqtt.Client, msg mqtt.Message) {
		// 处理接收到的消息
		s.handleMQTTMessage(client, msg.Topic(), string(msg.Payload()), int(msg.Qos()), msg.Retained())
	}); token.Wait() && token.Error() != nil {
		log.Printf("MQTT subscribe failed for client %s, topic %s: %v", client.ID, topic, token.Error())
		s.sendWebSocketMessage(client, WebSocketMessage{
			Type:    "subscribe_result",
			Topic:   topic,
			Payload: token.Error().Error(),
		})
		return token.Error()
	}

	log.Printf("MQTT client %s subscribed to topic: %s", client.ID, topic)
	s.sendWebSocketMessage(client, WebSocketMessage{
		Type:  "subscription_success",
		Topic: topic,
		QoS:   qos,
	})

	return nil
}

// handleUnsubscribe 处理取消订阅
func (s *MQTTProxyService) handleUnsubscribe(client *MQTTClient, topic string) error {
	if !client.IsConnected || client.Client == nil {
		return fmt.Errorf("MQTT client not connected")
	}

	if token := client.Client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Printf("MQTT unsubscribe failed for client %s, topic %s: %v", client.ID, topic, token.Error())
		return token.Error()
	}

	log.Printf("MQTT client %s unsubscribed from topic: %s", client.ID, topic)
	return nil
}

// handlePublish 处理发布
func (s *MQTTProxyService) handlePublish(client *MQTTClient, topic, payload string, qos int, retain bool) error {
	if !client.IsConnected || client.Client == nil {
		return fmt.Errorf("MQTT client not connected")
	}

	if token := client.Client.Publish(topic, byte(qos), retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("MQTT publish failed for client %s, topic %s: %v", client.ID, topic, token.Error())
		s.sendWebSocketMessage(client, WebSocketMessage{
			Type:    "publish_result",
			Topic:   topic,
			Payload: token.Error().Error(),
		})
		return token.Error()
	}

	log.Printf("MQTT client %s published to topic: %s", client.ID, topic)
	s.sendWebSocketMessage(client, WebSocketMessage{
		Type:    "publish_result",
		Topic:   topic,
		Payload: "success",
	})

	return nil
}

// handleMQTTMessage 处理MQTT消息
func (s *MQTTProxyService) handleMQTTMessage(client *MQTTClient, topic, payload string, qos int, retain bool) {
	log.Printf("MQTT message received for client %s: %s -> %s", client.ID, topic, payload)

	s.sendWebSocketMessage(client, WebSocketMessage{
		Type:    "mqtt_message",
		Topic:   topic,
		Payload: payload,
		QoS:     qos,
		Retain:  retain,
	})
}

// sendWebSocketMessage 发送WebSocket消息
func (s *MQTTProxyService) sendWebSocketMessage(client *MQTTClient, message WebSocketMessage) {
	if client.WSConn == nil {
		return
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v", err)
		return
	}

	if err := client.WSConn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		log.Printf("Failed to send WebSocket message: %v", err)
	}
}
