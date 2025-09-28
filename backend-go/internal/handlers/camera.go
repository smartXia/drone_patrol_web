package handlers

import (
	"net/http"

	"drone-patrol-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCameras 获取所有摄像头
func (h *Handlers) GetCameras(c *gin.Context) {
	cameras, err := h.cameraService.GetCameras()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "获取摄像头列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取摄像头列表成功",
		"data":    cameras,
	})
}

// CreateCamera 创建摄像头
func (h *Handlers) CreateCamera(c *gin.Context) {
	var camera services.Camera
	if err := c.ShouldBindJSON(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 生成摄像头ID
	camera.ID = uuid.New().String()

	// 设置默认状态
	if camera.Status == "" {
		camera.Status = "offline"
	}

	// 设置默认分辨率
	if camera.Resolution == "" {
		camera.Resolution = "1280x720"
	}

	// 设置默认帧率
	if camera.FPS == 0 {
		camera.FPS = 25
	}

	err := h.cameraService.CreateCamera(&camera)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "创建摄像头失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建摄像头成功",
		"data":    camera,
	})
}

// UpdateCamera 更新摄像头
func (h *Handlers) UpdateCamera(c *gin.Context) {
	cameraID := c.Param("camera_id")
	if cameraID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "摄像头ID不能为空",
		})
		return
	}

	var camera services.Camera
	if err := c.ShouldBindJSON(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	err := h.cameraService.UpdateCamera(cameraID, &camera)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "更新摄像头失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新摄像头成功",
	})
}

// DeleteCamera 删除摄像头
func (h *Handlers) DeleteCamera(c *gin.Context) {
	cameraID := c.Param("camera_id")
	if cameraID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "摄像头ID不能为空",
		})
		return
	}

	err := h.cameraService.DeleteCamera(cameraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "删除摄像头失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除摄像头成功",
	})
}

// TestCameraConnection 测试摄像头连接
func (h *Handlers) TestCameraConnection(c *gin.Context) {
	cameraID := c.Param("camera_id")
	if cameraID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "摄像头ID不能为空",
		})
		return
	}

	// 获取摄像头信息
	camera, err := h.cameraService.GetCamera(cameraID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "摄像头不存在",
			"error":   err.Error(),
		})
		return
	}

	// 测试连接
	err = h.cameraService.TestCameraConnection(camera)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "摄像头连接测试失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "摄像头连接测试成功",
	})
}

// StartCameraStream 开始摄像头流
func (h *Handlers) StartCameraStream(c *gin.Context) {
	cameraID := c.Param("camera_id")
	if cameraID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "摄像头ID不能为空",
		})
		return
	}

	// 检查摄像头是否存在
	_, err := h.cameraService.GetCamera(cameraID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "摄像头不存在",
			"error":   err.Error(),
		})
		return
	}

	// 开始流
	err = h.cameraService.StartCameraStream(cameraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "开始摄像头流失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "摄像头流开始成功",
	})
}

// StopCameraStream 停止摄像头流
func (h *Handlers) StopCameraStream(c *gin.Context) {
	cameraID := c.Param("camera_id")
	if cameraID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "摄像头ID不能为空",
		})
		return
	}

	// 检查摄像头是否存在
	_, err := h.cameraService.GetCamera(cameraID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "摄像头不存在",
			"error":   err.Error(),
		})
		return
	}

	// 停止流
	err = h.cameraService.StopCameraStream(cameraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "停止摄像头流失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "摄像头流停止成功",
	})
}
