package handlers

import (
	"drone-patrol-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取设备列表
func (h *Handlers) GetDevices(c *gin.Context) {
	response, err := h.deviceService.GetDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 获取当前设备信息
func (h *Handlers) GetCurrentDevices(c *gin.Context) {
	response, err := h.deviceService.GetCurrentDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 创建设备
func (h *Handlers) CreateDevice(c *gin.Context) {
	var payload models.DevicePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.deviceService.CreateDevice(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 更新设备
func (h *Handlers) UpdateDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	var payload models.DeviceUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.deviceService.UpdateDevice(deviceID, &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 删除设备
func (h *Handlers) DeleteDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	response, err := h.deviceService.DeleteDevice(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 设置当前设备
func (h *Handlers) SetCurrentDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	response, err := h.deviceService.SetCurrentDevice(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 设置网关设备
func (h *Handlers) SetGatewayDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	response, err := h.deviceService.SetGatewayDevice(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 清空所有设备
func (h *Handlers) ClearAllDevices(c *gin.Context) {
	response, err := h.deviceService.ClearAllDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 删除默认设备
func (h *Handlers) RemoveDefaultDevices(c *gin.Context) {
	response, err := h.deviceService.RemoveDefaultDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
