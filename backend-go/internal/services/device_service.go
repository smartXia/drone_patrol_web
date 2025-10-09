package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"drone-patrol-backend/internal/database"
	"drone-patrol-backend/internal/models"

	"github.com/google/uuid"
)

type DeviceService struct {
	db *database.DB
}

func NewDeviceService(db *database.DB) *DeviceService {
	return &DeviceService{db: db}
}

// 获取设备列表
func (s *DeviceService) GetDevices() (*models.APIResponse, error) {
	query := `SELECT id, name, sn, type, status, airport_sn, last_seen, created_at, updated_at, is_current, is_gateway 
			  FROM devices ORDER BY is_current DESC, created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取设备列表失败: %v", err),
		}, err
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		var lastSeen sql.NullInt64

		err := rows.Scan(
			&device.ID, &device.Name, &device.SN, &device.Type, &device.Status, &device.AirportSN,
			&lastSeen, &device.CreatedAt, &device.UpdatedAt, &device.IsCurrent, &device.IsGateway,
		)
		if err != nil {
			return &models.APIResponse{
				Code:    1,
				Message: fmt.Sprintf("扫描设备数据失败: %v", err),
			}, err
		}

		if lastSeen.Valid {
			t := time.Unix(lastSeen.Int64/1000, (lastSeen.Int64%1000)*1000000)
			device.LastSeen = &t
		}

		devices = append(devices, device)
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    devices,
	}, nil
}

// 获取当前设备信息
func (s *DeviceService) GetCurrentDevices() (*models.APIResponse, error) {
	// 获取当前设备
	currentDeviceQuery := `SELECT id, name, sn, type, status, airport_sn, last_seen FROM devices WHERE is_current = 1`
	var currentDevice *models.Device

	row := s.db.QueryRow(currentDeviceQuery)
	var device models.Device
	var lastSeen sql.NullInt64

	err := row.Scan(&device.ID, &device.Name, &device.SN, &device.Type, &device.Status, &device.AirportSN, &lastSeen)
	if err == nil {
		if lastSeen.Valid {
			t := time.Unix(lastSeen.Int64/1000, (lastSeen.Int64%1000)*1000000)
			device.LastSeen = &t
		}
		currentDevice = &device
	} else if err != sql.ErrNoRows {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取当前设备失败: %v", err),
		}, err
	}

	// 获取当前网关
	currentGatewayQuery := `SELECT id, name, sn, type, status, airport_sn, last_seen FROM devices WHERE is_gateway = 1`
	var currentGateway *models.Device

	row = s.db.QueryRow(currentGatewayQuery)
	var gateway models.Device
	var gatewayLastSeen sql.NullInt64

	err = row.Scan(&gateway.ID, &gateway.Name, &gateway.SN, &gateway.Type, &gateway.Status, &gateway.AirportSN, &gatewayLastSeen)
	if err == nil {
		if gatewayLastSeen.Valid {
			t := time.Unix(gatewayLastSeen.Int64/1000, (gatewayLastSeen.Int64%1000)*1000000)
			gateway.LastSeen = &t
		}
		currentGateway = &gateway
	} else if err != sql.ErrNoRows {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取当前网关失败: %v", err),
		}, err
	}

	result := make(map[string]interface{})
	if currentDevice != nil {
		result["device"] = currentDevice
	}
	if currentGateway != nil {
		result["gateway"] = currentGateway
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    result,
	}, nil
}

// 创建设备
func (s *DeviceService) CreateDevice(payload *models.DevicePayload) (*models.APIResponse, error) {
	// 检查设备SN是否已存在
	var existingID string
	err := s.db.QueryRow("SELECT id FROM devices WHERE sn = ?", payload.SN).Scan(&existingID)
	if err == nil {
		return &models.APIResponse{
			Code:    1,
			Message: "设备序列号已存在",
		}, nil
	} else if err != sql.ErrNoRows {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("检查设备SN失败: %v", err),
		}, err
	}

	// 创建设备
	deviceID := uuid.New().String()
	currentTime := time.Now().UnixMilli()

	query := `INSERT INTO devices (id, name, sn, type, status, airport_sn, last_seen, created_at, updated_at, is_current, is_gateway)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 0)`

	_, err = s.db.Exec(query, deviceID, payload.Name, payload.SN, payload.Type, payload.Status, payload.AirportSN, nil, currentTime, currentTime)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("创建设备失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "设备创建成功",
		Data:    map[string]string{"id": deviceID},
	}, nil
}

// 更新设备
func (s *DeviceService) UpdateDevice(deviceID string, payload *models.DeviceUpdatePayload) (*models.APIResponse, error) {
	// 构建更新查询
	setParts := []string{}
	args := []interface{}{}

	if payload.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *payload.Name)
	}
	if payload.SN != nil {
		setParts = append(setParts, "sn = ?")
		args = append(args, *payload.SN)
	}
	if payload.Type != nil {
		setParts = append(setParts, "type = ?")
		args = append(args, *payload.Type)
	}
	if payload.Status != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *payload.Status)
	}
	if payload.AirportSN != nil {
		setParts = append(setParts, "airport_sn = ?")
		args = append(args, *payload.AirportSN)
	}

	if len(setParts) == 0 {
		return &models.APIResponse{
			Code:    1,
			Message: "没有提供要更新的字段",
		}, nil
	}

	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now().UnixMilli())
	args = append(args, deviceID)

	query := fmt.Sprintf("UPDATE devices SET %s WHERE id = ?", strings.Join(setParts, ", "))

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("更新设备失败: %v", err),
		}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取更新结果失败: %v", err),
		}, err
	}

	if rowsAffected == 0 {
		return &models.APIResponse{
			Code:    1,
			Message: "设备不存在",
		}, nil
	}

	return &models.APIResponse{
		Code:    0,
		Message: "设备更新成功",
	}, nil
}

// 删除设备
func (s *DeviceService) DeleteDevice(deviceID string) (*models.APIResponse, error) {
	result, err := s.db.Exec("DELETE FROM devices WHERE id = ?", deviceID)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("删除设备失败: %v", err),
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
			Message: "设备不存在",
		}, nil
	}

	return &models.APIResponse{
		Code:    0,
		Message: "设备删除成功",
	}, nil
}

// 设置当前设备
func (s *DeviceService) SetCurrentDevice(deviceID string) (*models.APIResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("开始事务失败: %v", err),
		}, err
	}
	defer tx.Rollback()

	// 清除所有设备的当前状态
	_, err = tx.Exec("UPDATE devices SET is_current = 0")
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("清除当前设备状态失败: %v", err),
		}, err
	}

	// 设置指定设备为当前设备
	result, err := tx.Exec("UPDATE devices SET is_current = 1 WHERE id = ?", deviceID)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("设置当前设备失败: %v", err),
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
			Message: "设备不存在",
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
		Message: "设置当前设备成功",
	}, nil
}

// 设置网关设备
func (s *DeviceService) SetGatewayDevice(deviceID string) (*models.APIResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("开始事务失败: %v", err),
		}, err
	}
	defer tx.Rollback()

	// 清除所有设备的网关状态
	_, err = tx.Exec("UPDATE devices SET is_gateway = 0")
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("清除网关设备状态失败: %v", err),
		}, err
	}

	// 设置指定设备为网关设备
	result, err := tx.Exec("UPDATE devices SET is_gateway = 1 WHERE id = ?", deviceID)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("设置网关设备失败: %v", err),
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
			Message: "设备不存在",
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
		Message: "设置网关设备成功",
	}, nil
}

// 清空所有设备
func (s *DeviceService) ClearAllDevices() (*models.APIResponse, error) {
	_, err := s.db.Exec("DELETE FROM devices")
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("清空设备失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "所有设备已清空",
	}, nil
}

// 删除默认设备
func (s *DeviceService) RemoveDefaultDevices() (*models.APIResponse, error) {
	result, err := s.db.Exec("DELETE FROM devices WHERE name = ?", "默认设备")
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("删除默认设备失败: %v", err),
		}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("获取删除结果失败: %v", err),
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: fmt.Sprintf("已删除 %d 个默认设备", rowsAffected),
	}, nil
}
