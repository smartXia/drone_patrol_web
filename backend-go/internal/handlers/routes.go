package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handlers) RegisterRoutes(r *gin.Engine) {
	// 健康检查和工具API
	r.GET("/api/health", h.Health)
	r.GET("/api/error-codes", h.GetErrorCodes)
	r.GET("/api/network/ping", h.Ping)

	// WebSocket支持
	r.GET("/ws/mqtt", h.WebSocketHandler)

	// MQTT配置管理API
	mqtt := r.Group("/api/mqtt")
	{
		mqtt.GET("/profiles", h.GetMQTTProfiles)
		mqtt.POST("/profiles", h.CreateMQTTProfile)
		mqtt.GET("/profiles/:pid", h.GetMQTTProfile)
		mqtt.PUT("/profiles/:pid", h.UpdateMQTTProfile)
		mqtt.DELETE("/profiles/:pid", h.DeleteMQTTProfile)
		mqtt.POST("/profiles/:pid/default", h.SetDefaultMQTTProfile)
		mqtt.POST("/test", h.TestMQTTConnection)
	}

	// 设备管理API
	devices := r.Group("/api/devices")
	{
		devices.GET("", h.GetDevices)
		devices.GET("/current", h.GetCurrentDevices)
		devices.POST("", h.CreateDevice)
		devices.PUT("/:device_id", h.UpdateDevice)
		devices.DELETE("/:device_id", h.DeleteDevice)
		devices.POST("/:device_id/set-current", h.SetCurrentDevice)
		devices.POST("/:device_id/set-gateway", h.SetGatewayDevice)
		devices.DELETE("/clear", h.ClearAllDevices)
		devices.DELETE("/remove-defaults", h.RemoveDefaultDevices)
	}

	// 摄像头管理API
	cameras := r.Group("/api/cameras")
	{
		cameras.GET("", h.GetCameras)
		cameras.POST("", h.CreateCamera)
		cameras.PUT("/:camera_id", h.UpdateCamera)
		cameras.DELETE("/:camera_id", h.DeleteCamera)
		cameras.POST("/:camera_id/test", h.TestCameraConnection)
		cameras.POST("/:camera_id/start", h.StartCameraStream)
		cameras.POST("/:camera_id/stop", h.StopCameraStream)
	}

	// Redis代理API
	redis := r.Group("/api/redis")
	{
		redis.POST("/connect/test", h.TestRedisConnection)
	}

	// Redis操作API (兼容原有路径)
	r.POST("/connect/test", h.TestRedisConnection)
	r.POST("/scan", h.ScanKeys)
	r.POST("/type", h.GetKeyType)
	r.POST("/ttl", h.GetKeyTTL)
	r.POST("/metadata", h.GetKeyMetadata)
	r.POST("/expire", h.SetKeyExpire)
	r.POST("/persist", h.PersistKey)
	r.POST("/rename", h.RenameKey)
	r.POST("/del", h.DeleteKey)
	r.POST("/get", h.GetStringValue)
	r.POST("/set", h.SetStringValue)

	// 哈希操作
	r.POST("/hash/getall", h.GetHashAll)
	r.POST("/hash/set", h.SetHashFields)
	r.POST("/hash/del", h.DeleteHashField)

	// 列表操作
	r.POST("/list/range", h.GetListRange)
	r.POST("/list/lpush", h.ListLPush)
	r.POST("/list/rpush", h.ListRPush)
	r.POST("/list/lpop", h.ListLPop)
	r.POST("/list/rpop", h.ListRPop)
	r.POST("/list/set", h.ListSet)
	r.POST("/list/lrem", h.ListLRem)

	// 集合操作
	r.POST("/set/scan", h.SetScan)
	r.POST("/set/sadd", h.SetSAdd)
	r.POST("/set/srem", h.SetSRem)

	// 有序集合操作
	r.POST("/zset/range", h.ZSetRange)
	r.POST("/zset/zadd", h.ZSetZAdd)
	r.POST("/zset/zrem", h.ZSetZRem)
	r.POST("/zset/zincrby", h.ZSetZIncrBy)
}
