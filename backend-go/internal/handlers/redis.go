package handlers

import (
	"drone-patrol-backend/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 测试Redis连接
func (h *Handlers) TestRedisConnection(c *gin.Context) {
	var payload models.RedisConnectPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.TestConnection(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 扫描Redis键
func (h *Handlers) ScanKeys(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.Scan(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 获取键类型
func (h *Handlers) GetKeyType(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.GetType(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 获取键TTL
func (h *Handlers) GetKeyTTL(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.GetTTL(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 获取键元数据
func (h *Handlers) GetKeyMetadata(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.GetKeyMetadata(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 设置键过期
func (h *Handlers) SetKeyExpire(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	secondsStr := c.Query("seconds")
	seconds, err := strconv.Atoi(secondsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: seconds必须是整数",
		})
		return
	}

	response, err := h.redisService.SetExpire(&payload, seconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 持久化键
func (h *Handlers) PersistKey(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.Persist(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 重命名键
func (h *Handlers) RenameKey(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	newKey := c.Query("newKey")
	if newKey == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: newKey不能为空",
		})
		return
	}

	response, err := h.redisService.Rename(&payload, newKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 删除键
func (h *Handlers) DeleteKey(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.Delete(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 获取字符串值
func (h *Handlers) GetStringValue(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.Get(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 设置字符串值
func (h *Handlers) SetStringValue(c *gin.Context) {
	var payload models.RedisValuePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.Set(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 获取哈希所有字段
func (h *Handlers) GetHashAll(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.HashGetAll(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 设置哈希字段
func (h *Handlers) SetHashFields(c *gin.Context) {
	var payload models.RedisHashPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.HashSet(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 删除哈希字段
func (h *Handlers) DeleteHashField(c *gin.Context) {
	var payload models.RedisHashPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.HashDel(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 获取列表范围
func (h *Handlers) GetListRange(c *gin.Context) {
	var payload models.RedisListPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ListRange(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 左推入列表
func (h *Handlers) ListLPush(c *gin.Context) {
	var payload models.RedisListPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ListLPush(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 右推入列表
func (h *Handlers) ListRPush(c *gin.Context) {
	var payload models.RedisListPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ListRPush(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 左弹出列表
func (h *Handlers) ListLPop(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ListLPop(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 右弹出列表
func (h *Handlers) ListRPop(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ListRPop(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 设置列表元素
func (h *Handlers) ListSet(c *gin.Context) {
	var payload models.RedisListPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ListSet(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 删除列表元素
func (h *Handlers) ListLRem(c *gin.Context) {
	var payload models.RedisListPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ListLRem(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 扫描集合
func (h *Handlers) SetScan(c *gin.Context) {
	var payload models.RedisKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.SetScan(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 添加集合成员
func (h *Handlers) SetSAdd(c *gin.Context) {
	var payload models.RedisSetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.SetSAdd(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 删除集合成员
func (h *Handlers) SetSRem(c *gin.Context) {
	var payload models.RedisSetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.SetSRem(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 获取有序集合范围
func (h *Handlers) ZSetRange(c *gin.Context) {
	var payload models.RedisZSetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ZSetRange(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 添加有序集合成员
func (h *Handlers) ZSetZAdd(c *gin.Context) {
	var payload models.RedisZSetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ZSetZAdd(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 删除有序集合成员
func (h *Handlers) ZSetZRem(c *gin.Context) {
	var payload models.RedisZSetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ZSetZRem(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 增加有序集合分数
func (h *Handlers) ZSetZIncrBy(c *gin.Context) {
	var payload models.RedisZSetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ZSetZIncrBy(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 执行Redis命令
func (h *Handlers) ExecuteRedisCommand(c *gin.Context) {
	var payload models.RedisCommandPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.redisService.ExecuteCommand(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
