package handlers

import (
	"drone-patrol-backend/internal/services"
)

type Handlers struct {
	deviceService    *services.DeviceService
	mqttService      *services.MQTTService
	redisService     *services.RedisService
	errorCodeService *services.ErrorCodeService
	MQTTProxy        *services.MQTTProxyService
	cameraService    *services.CameraService
}

func NewHandlers(
	deviceService *services.DeviceService,
	mqttService *services.MQTTService,
	redisService *services.RedisService,
	errorCodeService *services.ErrorCodeService,
	mqttProxy *services.MQTTProxyService,
	cameraService *services.CameraService,
) *Handlers {
	return &Handlers{
		deviceService:    deviceService,
		mqttService:      mqttService,
		redisService:     redisService,
		errorCodeService: errorCodeService,
		MQTTProxy:        mqttProxy,
		cameraService:    cameraService,
	}
}
