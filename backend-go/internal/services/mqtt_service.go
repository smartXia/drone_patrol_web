package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"drone-patrol-backend/internal/database"
	"drone-patrol-backend/internal/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type MQTTService struct {
	db *database.DB
}

func NewMQTTService(db *database.DB) *MQTTService {
	return &MQTTService{db: db}
}

// 获取MQTT配置列表
func (s *MQTTService) GetProfiles() (*models.APIResponse, error) {
	query := `SELECT id, name, config, is_default, updated_at FROM mqtt_profiles ORDER BY is_default DESC, updated_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取MQTT配置列表失败: %v", err),
		}, err
	}
	defer rows.Close()

	var profiles []models.MQTTProfile
	for rows.Next() {
		var profile models.MQTTProfile
		var configJSON string

		err := rows.Scan(&profile.ID, &profile.Name, &configJSON, &profile.IsDefault, &profile.UpdatedAt)
		if err != nil {
			return &models.APIResponse{
				Code:    1,
				Message: fmt.Sprintf("扫描MQTT配置数据失败: %v", err),
			}, err
		}

		// 解析配置JSON
		err = json.Unmarshal([]byte(configJSON), &profile.Config)
		if err != nil {
			return &models.APIResponse{
				Code:    1,
				Message: fmt.Sprintf("解析MQTT配置失败: %v", err),
			}, err
		}

		profiles = append(profiles, profile)
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    profiles,
	}, nil
}

// 创建MQTT配置
func (s *MQTTService) CreateProfile(payload *models.MQTTProfilePayload) (*models.APIResponse, error) {
	profileID := uuid.New().String()
	currentTime := time.Now().UnixMilli()

	// 如果设置为默认，先清除其他默认配置
	if payload.IsDefault {
		_, err := s.db.Exec("UPDATE mqtt_profiles SET is_default = 0")
		if err != nil {
			return &models.APIResponse{
				Code:    1,
				Message: fmt.Sprintf("清除默认配置失败: %v", err),
			}, err
		}
	}

	// 序列化配置
	configJSON, err := json.Marshal(payload.Config)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("序列化配置失败: %v", err),
		}, err
	}

	// 插入新配置
	query := `INSERT INTO mqtt_profiles (id, name, config, is_default, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, profileID, payload.Name, string(configJSON), payload.IsDefault, currentTime)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("创建MQTT配置失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "MQTT配置创建成功",
		Data:    map[string]string{"id": profileID},
	}, nil
}

// 获取单个MQTT配置
func (s *MQTTService) GetProfile(profileID string) (*models.APIResponse, error) {
	query := `SELECT id, name, config, is_default, updated_at FROM mqtt_profiles WHERE id = ?`

	row := s.db.QueryRow(query, profileID)
	var profile models.MQTTProfile
	var configJSON string

	err := row.Scan(&profile.ID, &profile.Name, &configJSON, &profile.IsDefault, &profile.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.APIResponse{
				Code:    1,
				Message: "MQTT配置不存在",
			}, nil
		}
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取MQTT配置失败: %v", err),
		}, err
	}

	// 解析配置JSON
	err = json.Unmarshal([]byte(configJSON), &profile.Config)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("解析MQTT配置失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    profile,
	}, nil
}

// 更新MQTT配置
func (s *MQTTService) UpdateProfile(profileID string, payload *models.MQTTProfilePayload) (*models.APIResponse, error) {
	// 检查配置是否存在
	var existingID string
	err := s.db.QueryRow("SELECT id FROM mqtt_profiles WHERE id = ?", profileID).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.APIResponse{
				Code:    1,
				Message: "MQTT配置不存在",
			}, nil
		}
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("检查MQTT配置失败: %v", err),
		}, err
	}

	// 如果设置为默认，先清除其他默认配置
	if payload.IsDefault {
		_, err = s.db.Exec("UPDATE mqtt_profiles SET is_default = 0")
		if err != nil {
			return &models.APIResponse{
				Code:    1,
				Message: fmt.Sprintf("清除默认配置失败: %v", err),
			}, err
		}
	}

	// 序列化配置
	configJSON, err := json.Marshal(payload.Config)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("序列化配置失败: %v", err),
		}, err
	}

	// 更新配置
	query := `UPDATE mqtt_profiles SET name = ?, config = ?, is_default = ?, updated_at = ? WHERE id = ?`
	_, err = s.db.Exec(query, payload.Name, string(configJSON), payload.IsDefault, time.Now().UnixMilli(), profileID)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("更新MQTT配置失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "MQTT配置更新成功",
	}, nil
}

// 删除MQTT配置
func (s *MQTTService) DeleteProfile(profileID string) (*models.APIResponse, error) {
	result, err := s.db.Exec("DELETE FROM mqtt_profiles WHERE id = ?", profileID)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("删除MQTT配置失败: %v", err),
		}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取删除结果失败: %v", err),
		}, err
	}

	if rowsAffected == 0 {
		return &models.APIResponse{
			Code:    1,
			Message: "MQTT配置不存在",
		}, nil
	}

	return &models.APIResponse{
		Code:    0,
		Message: "MQTT配置删除成功",
	}, nil
}

// 设置默认MQTT配置
func (s *MQTTService) SetDefaultProfile(profileID string) (*models.APIResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("开始事务失败: %v", err),
		}, err
	}
	defer tx.Rollback()

	// 清除所有默认配置
	_, err = tx.Exec("UPDATE mqtt_profiles SET is_default = 0")
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("清除默认配置失败: %v", err),
		}, err
	}

	// 设置指定配置为默认
	result, err := tx.Exec("UPDATE mqtt_profiles SET is_default = 1 WHERE id = ?", profileID)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("设置默认配置失败: %v", err),
		}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取设置结果失败: %v", err),
		}, err
	}

	if rowsAffected == 0 {
		return &models.APIResponse{
			Code:    1,
			Message: "MQTT配置不存在",
		}, nil
	}

	err = tx.Commit()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("提交事务失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "设置默认配置成功",
	}, nil
}

// 测试MQTT连接
func (s *MQTTService) TestConnection(config map[string]interface{}) (*models.APIResponse, error) {
	// 解析配置
	broker, ok := config["broker"].(string)
	if !ok {
		return &models.APIResponse{
			Code:    1,
			Message: "缺少broker配置",
		}, nil
	}

	port, ok := config["port"].(float64)
	if !ok {
		return &models.APIResponse{
			Code:    1,
			Message: "缺少port配置",
		}, nil
	}

	username, _ := config["username"].(string)
	password, _ := config["password"].(string)
	clientID, _ := config["clientId"].(string)

	// 创建MQTT客户端选项
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, int(port)))
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetCleanSession(true)

	// 创建客户端
	client := mqtt.NewClient(opts)

	// 尝试连接
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("连接错误: %v", token.Error()),
			Data:    map[string]bool{"connected": false},
		}, nil
	}

	// 连接成功，断开连接
	client.Disconnect(250)

	return &models.APIResponse{
		Code:    0,
		Message: "连接成功",
		Data:    map[string]bool{"connected": true},
	}, nil
}
