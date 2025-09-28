package models

import "time"

// 通用响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MQTT配置相关
type MQTTProfile struct {
	ID        string                 `json:"id" db:"id"`
	Name      string                 `json:"name" db:"name"`
	Config    map[string]interface{} `json:"config" db:"config"`
	IsDefault bool                   `json:"isDefault" db:"is_default"`
	UpdatedAt int64                  `json:"updatedAt" db:"updated_at"`
}

type MQTTProfilePayload struct {
	Name      string                 `json:"name" binding:"required"`
	Config    map[string]interface{} `json:"config" binding:"required"`
	IsDefault bool                   `json:"isDefault"`
}

// 设备相关
type Device struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	SN        string     `json:"sn" db:"sn"`
	Type      string     `json:"type" db:"type"`
	Status    string     `json:"status" db:"status"`
	LastSeen  *time.Time `json:"lastSeen" db:"last_seen"`
	CreatedAt int64      `json:"createdAt" db:"created_at"`
	UpdatedAt int64      `json:"updatedAt" db:"updated_at"`
	IsCurrent bool       `json:"isCurrent" db:"is_current"`
	IsGateway bool       `json:"isGateway" db:"is_gateway"`
}

type DevicePayload struct {
	Name   string `json:"name" binding:"required"`
	SN     string `json:"sn" binding:"required"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type DeviceUpdatePayload struct {
	Name   *string `json:"name,omitempty"`
	SN     *string `json:"sn,omitempty"`
	Type   *string `json:"type,omitempty"`
	Status *string `json:"status,omitempty"`
}

// 错误码相关
type ErrorCode struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
}

// Redis相关
type RedisConnectPayload struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type RedisKeyPayload struct {
	Key    string `json:"key" binding:"required"`
	Cursor string `json:"cursor,omitempty"`
}

type RedisValuePayload struct {
	Key   string      `json:"key" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

type RedisHashPayload struct {
	Key   string            `json:"key" binding:"required"`
	Field string            `json:"field,omitempty"`
	Value map[string]string `json:"value,omitempty"`
}

type RedisListPayload struct {
	Key   string `json:"key" binding:"required"`
	Start int    `json:"start,omitempty"`
	Stop  int    `json:"stop,omitempty"`
	Value string `json:"value,omitempty"`
	Index int    `json:"index,omitempty"`
	Count int    `json:"count,omitempty"`
}

type RedisSetPayload struct {
	Key     string   `json:"key" binding:"required"`
	Member  string   `json:"member,omitempty"`
	Members []string `json:"members,omitempty"`
}

type RedisZSetPayload struct {
	Key    string  `json:"key" binding:"required"`
	Member string  `json:"member,omitempty"`
	Score  float64 `json:"score,omitempty"`
	Min    string  `json:"min,omitempty"`
	Max    string  `json:"max,omitempty"`
}

// 网络测试
type PingPayload struct {
	Host string `json:"host" binding:"required"`
	Port int    `json:"port" binding:"required"`
}
