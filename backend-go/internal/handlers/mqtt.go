package handlers

import (
	"drone-patrol-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取MQTT配置列表
func (h *Handlers) GetMQTTProfiles(c *gin.Context) {
	response, err := h.mqttService.GetProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 创建MQTT配置
func (h *Handlers) CreateMQTTProfile(c *gin.Context) {
	var payload models.MQTTProfilePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.mqttService.CreateProfile(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 获取单个MQTT配置
func (h *Handlers) GetMQTTProfile(c *gin.Context) {
	profileID := c.Param("pid")
	response, err := h.mqttService.GetProfile(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 更新MQTT配置
func (h *Handlers) UpdateMQTTProfile(c *gin.Context) {
	profileID := c.Param("pid")
	var payload models.MQTTProfilePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.mqttService.UpdateProfile(profileID, &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 删除MQTT配置
func (h *Handlers) DeleteMQTTProfile(c *gin.Context) {
	profileID := c.Param("pid")
	response, err := h.mqttService.DeleteProfile(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 设置默认MQTT配置
func (h *Handlers) SetDefaultMQTTProfile(c *gin.Context) {
	profileID := c.Param("pid")
	response, err := h.mqttService.SetDefaultProfile(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 测试MQTT连接
func (h *Handlers) TestMQTTConnection(c *gin.Context) {
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.mqttService.TestConnection(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
