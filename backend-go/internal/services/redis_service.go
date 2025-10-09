package services

import (
	"context"
	"fmt"
	"time"

	"drone-patrol-backend/internal/models"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
	config *models.RedisConnectPayload
}

func NewRedisService() *RedisService {
	return &RedisService{}
}

// 测试Redis连接
func (s *RedisService) TestConnection(payload *models.RedisConnectPayload) (*models.APIResponse, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", payload.Host, payload.Port),
		Password: payload.Password,
		DB:       payload.DB,
	})
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("连接错误: %v", err),
			Data:    map[string]bool{"connected": false},
		}, nil
	}

	// 连接成功后，保存配置并创建持久连接
	s.config = payload
	s.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", payload.Host, payload.Port),
		Password: payload.Password,
		DB:       payload.DB,
	})

	return &models.APIResponse{
		Code:    0,
		Message: "连接成功",
		Data:    map[string]bool{"connected": true},
	}, nil
}

// 扫描Redis键
func (s *RedisService) Scan(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()

	// 使用正确的扫描参数
	var cursor uint64 = 0
	if payload.Cursor != "" {
		fmt.Sscanf(payload.Cursor, "%d", &cursor)
	}

	pattern := payload.Key
	if pattern == "" {
		pattern = "*"
	}

	keys, newCursor, err := client.Scan(ctx, cursor, pattern, 100).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("扫描键失败: %v", err),
		}, err
	}

	// 为每个键获取类型和TTL
	var keyItems []map[string]interface{}
	for _, key := range keys {
		keyType, _ := client.Type(ctx, key).Result()
		ttl, _ := client.TTL(ctx, key).Result()
		keyItems = append(keyItems, map[string]interface{}{
			"key":  key,
			"type": keyType,
			"ttl":  int64(ttl.Seconds()),
		})
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data: map[string]interface{}{
			"keys":   keyItems,
			"cursor": fmt.Sprintf("%d", newCursor),
		},
	}, nil
}

// 获取键类型
func (s *RedisService) GetType(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	keyType, err := client.Type(ctx, payload.Key).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取键类型失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    map[string]string{"type": keyType},
	}, nil
}

// 获取键TTL
func (s *RedisService) GetTTL(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	ttl, err := client.TTL(ctx, payload.Key).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取TTL失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    map[string]int64{"ttl": int64(ttl.Seconds())},
	}, nil
}

// 获取键元数据（TTL、类型、内存使用等）
func (s *RedisService) GetKeyMetadata(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()

	// 获取键类型
	keyType, err := client.Type(ctx, payload.Key).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取键类型失败: %v", err),
		}, err
	}

	// 获取TTL
	ttl, err := client.TTL(ctx, payload.Key).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取TTL失败: %v", err),
		}, err
	}

	// 获取内存使用情况（如果支持）
	memoryUsage := int64(0)
	if keyType != "none" {
		// 使用MEMORY USAGE命令获取内存使用（Redis 4.0+）
		memUsage, err := client.MemoryUsage(ctx, payload.Key).Result()
		if err == nil {
			memoryUsage = memUsage
		}
	}

	// 获取键长度
	var keyLength int64
	switch keyType {
	case "string":
		keyLength, _ = client.StrLen(ctx, payload.Key).Result()
	case "list":
		keyLength, _ = client.LLen(ctx, payload.Key).Result()
	case "set":
		keyLength, _ = client.SCard(ctx, payload.Key).Result()
	case "zset":
		keyLength, _ = client.ZCard(ctx, payload.Key).Result()
	case "hash":
		keyLength, _ = client.HLen(ctx, payload.Key).Result()
	}

	// 检查键是否存在
	exists, err := client.Exists(ctx, payload.Key).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("检查键存在性失败: %v", err),
		}, err
	}

	metadata := map[string]interface{}{
		"type":        keyType,
		"ttl":         int64(ttl.Seconds()),
		"memoryUsage": memoryUsage,
		"length":      keyLength,
		"exists":      exists > 0,
	}

	// 如果键存在，添加更多信息
	if exists > 0 {
		// 获取键的编码类型
		encoding, err := client.ObjectEncoding(ctx, payload.Key).Result()
		if err == nil {
			metadata["encoding"] = encoding
		}

		// 获取键的引用计数
		refCount, err := client.ObjectRefCount(ctx, payload.Key).Result()
		if err == nil {
			metadata["refCount"] = refCount
		}

		// 获取键的空闲时间（idle time）
		idleTime, err := client.ObjectIdleTime(ctx, payload.Key).Result()
		if err == nil {
			metadata["idleTime"] = int64(idleTime.Seconds())
		}
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    metadata,
	}, nil
}

// 设置键过期
func (s *RedisService) SetExpire(payload *models.RedisKeyPayload, seconds int) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.Expire(ctx, payload.Key, time.Duration(seconds)*time.Second).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("设置过期失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    map[string]bool{"result": result},
	}, nil
}

// 持久化键
func (s *RedisService) Persist(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.Persist(ctx, payload.Key).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("持久化键失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    map[string]bool{"result": result},
	}, nil
}

// 重命名键
func (s *RedisService) Rename(payload *models.RedisKeyPayload, newKey string) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	err := client.Rename(ctx, payload.Key, newKey).Err()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("重命名键失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "重命名成功",
	}, nil
}

// 删除键
func (s *RedisService) Delete(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.Del(ctx, payload.Key).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("删除键失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "删除成功",
		Data:    map[string]int64{"deleted": result},
	}, nil
}

// 获取字符串值
func (s *RedisService) Get(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	value, err := client.Get(ctx, payload.Key).Result()
	if err != nil {
		if err == redis.Nil {
			return &models.APIResponse{
				Code:    1,
				Message: "键不存在",
			}, nil
		}
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取值失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    map[string]string{"value": value},
	}, nil
}

// 设置字符串值
func (s *RedisService) Set(payload *models.RedisValuePayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	err := client.Set(ctx, payload.Key, payload.Value, 0).Err()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("设置值失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "设置成功",
	}, nil
}

// 获取哈希所有字段
func (s *RedisService) HashGetAll(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	fields, err := client.HGetAll(ctx, payload.Key).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取哈希字段失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    fields,
	}, nil
}

// 设置哈希字段
func (s *RedisService) HashSet(payload *models.RedisHashPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	err := client.HMSet(ctx, payload.Key, payload.Value).Err()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("设置哈希字段失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "设置成功",
	}, nil
}

// 删除哈希字段
func (s *RedisService) HashDel(payload *models.RedisHashPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.HDel(ctx, payload.Key, payload.Field).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("删除哈希字段失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "删除成功",
		Data:    map[string]int64{"deleted": result},
	}, nil
}

// 获取列表范围
func (s *RedisService) ListRange(payload *models.RedisListPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	values, err := client.LRange(ctx, payload.Key, int64(payload.Start), int64(payload.Stop)).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取列表范围失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    values,
	}, nil
}

// 左推入列表
func (s *RedisService) ListLPush(payload *models.RedisListPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.LPush(ctx, payload.Key, payload.Value).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("左推入列表失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "推入成功",
		Data:    map[string]int64{"length": result},
	}, nil
}

// 右推入列表
func (s *RedisService) ListRPush(payload *models.RedisListPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.RPush(ctx, payload.Key, payload.Value).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("右推入列表失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "推入成功",
		Data:    map[string]int64{"length": result},
	}, nil
}

// 左弹出列表
func (s *RedisService) ListLPop(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	value, err := client.LPop(ctx, payload.Key).Result()
	if err != nil {
		if err == redis.Nil {
			return &models.APIResponse{
				Code:    1,
				Message: "列表为空",
			}, nil
		}
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("左弹出列表失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "弹出成功",
		Data:    map[string]string{"value": value},
	}, nil
}

// 右弹出列表
func (s *RedisService) ListRPop(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	value, err := client.RPop(ctx, payload.Key).Result()
	if err != nil {
		if err == redis.Nil {
			return &models.APIResponse{
				Code:    1,
				Message: "列表为空",
			}, nil
		}
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("右弹出列表失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "弹出成功",
		Data:    map[string]string{"value": value},
	}, nil
}

// 设置列表元素
func (s *RedisService) ListSet(payload *models.RedisListPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	err := client.LSet(ctx, payload.Key, int64(payload.Index), payload.Value).Err()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("设置列表元素失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "设置成功",
	}, nil
}

// 删除列表元素
func (s *RedisService) ListLRem(payload *models.RedisListPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.LRem(ctx, payload.Key, int64(payload.Count), payload.Value).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("删除列表元素失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "删除成功",
		Data:    map[string]int64{"removed": result},
	}, nil
}

// 扫描集合
func (s *RedisService) SetScan(payload *models.RedisKeyPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	members, cursor, err := client.SScan(ctx, payload.Key, 0, "*", 100).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("扫描集合失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data: map[string]interface{}{
			"members": members,
			"cursor":  cursor,
		},
	}, nil
}

// 添加集合成员
func (s *RedisService) SetSAdd(payload *models.RedisSetPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.SAdd(ctx, payload.Key, payload.Member).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("添加集合成员失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "添加成功",
		Data:    map[string]int64{"added": result},
	}, nil
}

// 删除集合成员
func (s *RedisService) SetSRem(payload *models.RedisSetPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.SRem(ctx, payload.Key, payload.Member).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("删除集合成员失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "删除成功",
		Data:    map[string]int64{"removed": result},
	}, nil
}

// 获取有序集合范围
func (s *RedisService) ZSetRange(payload *models.RedisZSetPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	members, err := client.ZRangeWithScores(ctx, payload.Key, 0, -1).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取有序集合范围失败: %v", err),
		}, err
	}

	var result []map[string]interface{}
	for _, member := range members {
		result = append(result, map[string]interface{}{
			"member": member.Member,
			"score":  member.Score,
		})
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    result,
	}, nil
}

// 添加有序集合成员
func (s *RedisService) ZSetZAdd(payload *models.RedisZSetPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.ZAdd(ctx, payload.Key, &redis.Z{
		Score:  payload.Score,
		Member: payload.Member,
	}).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("添加有序集合成员失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "添加成功",
		Data:    map[string]int64{"added": result},
	}, nil
}

// 删除有序集合成员
func (s *RedisService) ZSetZRem(payload *models.RedisZSetPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.ZRem(ctx, payload.Key, payload.Member).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("删除有序集合成员失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "删除成功",
		Data:    map[string]int64{"removed": result},
	}, nil
}

// 增加有序集合分数
func (s *RedisService) ZSetZIncrBy(payload *models.RedisZSetPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()
	result, err := client.ZIncrBy(ctx, payload.Key, payload.Score, payload.Member).Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("增加有序集合分数失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "增加成功",
		Data:    map[string]float64{"score": result},
	}, nil
}

// 执行Redis命令
func (s *RedisService) ExecuteCommand(payload *models.RedisCommandPayload) (*models.APIResponse, error) {
	client := s.getClient()
	if client == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "Redis客户端未初始化",
		}, nil
	}

	ctx := context.Background()

	// 解析命令
	args := parseCommand(payload.Command)
	if len(args) == 0 {
		return &models.APIResponse{
			Code:    1,
			Message: "无效的命令",
		}, nil
	}

	// 执行命令
	cmd := client.Do(ctx, args...)
	result, err := cmd.Result()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("命令执行失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "命令执行成功",
		Data: map[string]interface{}{
			"command": payload.Command,
			"result":  result,
		},
	}, nil
}

// 解析命令字符串
func parseCommand(command string) []interface{} {
	var args []interface{}
	var current string
	var inQuotes bool
	var quoteChar rune

	for _, char := range command {
		if char == '"' || char == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				args = append(args, current)
				current = ""
			} else {
				current += string(char)
			}
		} else if char == ' ' && !inQuotes {
			if current != "" {
				args = append(args, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		args = append(args, current)
	}

	return args
}

// 获取Redis客户端
func (s *RedisService) getClient() *redis.Client {
	return s.client
}
