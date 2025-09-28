package handlers

import (
	"drone-patrol-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 健康检查
func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "ok",
		Data: map[string]string{
			"service": "unified-go",
			"version": "1.0.0",
		},
	})
}

// 获取错误码列表
func (h *Handlers) GetErrorCodes(c *gin.Context) {
	response, err := h.errorCodeService.GetErrorCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 网络ping测试
func (h *Handlers) Ping(c *gin.Context) {
	var payload models.PingPayload
	if err := c.ShouldBindQuery(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 这里应该实现实际的ping逻辑
	// 暂时返回成功响应
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "ping成功",
		Data: map[string]interface{}{
			"host": payload.Host,
			"port": payload.Port,
			"ping": true,
		},
	})
}
