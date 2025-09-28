package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Camera 摄像头结构体
type Camera struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	URL         string    `json:"url" db:"url"`
	Username    string    `json:"username" db:"username"`
	Password    string    `json:"password" db:"password"`
	Resolution  string    `json:"resolution" db:"resolution"`
	FPS         int       `json:"fps" db:"fps"`
	Status      string    `json:"status" db:"status"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CameraService 摄像头服务
type CameraService struct {
	db *sql.DB
}

// NewCameraService 创建摄像头服务
func NewCameraService(db *sql.DB) *CameraService {
	return &CameraService{db: db}
}

// CreateCamera 创建摄像头
func (s *CameraService) CreateCamera(camera *Camera) error {
	query := `
		INSERT INTO cameras (id, name, url, username, password, resolution, fps, status, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err := s.db.Exec(query,
		camera.ID,
		camera.Name,
		camera.URL,
		camera.Username,
		camera.Password,
		camera.Resolution,
		camera.FPS,
		camera.Status,
		camera.Description,
		now,
		now,
	)

	if err != nil {
		log.Printf("创建摄像头失败: %v", err)
		return err
	}

	log.Printf("摄像头创建成功: %s", camera.Name)
	return nil
}

// GetCameras 获取所有摄像头
func (s *CameraService) GetCameras() ([]Camera, error) {
	query := `
		SELECT id, name, url, username, password, resolution, fps, status, description, created_at, updated_at
		FROM cameras
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("查询摄像头列表失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var cameras []Camera
	for rows.Next() {
		var camera Camera
		err := rows.Scan(
			&camera.ID,
			&camera.Name,
			&camera.URL,
			&camera.Username,
			&camera.Password,
			&camera.Resolution,
			&camera.FPS,
			&camera.Status,
			&camera.Description,
			&camera.CreatedAt,
			&camera.UpdatedAt,
		)
		if err != nil {
			log.Printf("扫描摄像头数据失败: %v", err)
			continue
		}
		cameras = append(cameras, camera)
	}

	log.Printf("查询到 %d 个摄像头", len(cameras))
	return cameras, nil
}

// GetCamera 根据ID获取摄像头
func (s *CameraService) GetCamera(id string) (*Camera, error) {
	query := `
		SELECT id, name, url, username, password, resolution, fps, status, description, created_at, updated_at
		FROM cameras
		WHERE id = ?
	`

	var camera Camera
	err := s.db.QueryRow(query, id).Scan(
		&camera.ID,
		&camera.Name,
		&camera.URL,
		&camera.Username,
		&camera.Password,
		&camera.Resolution,
		&camera.FPS,
		&camera.Status,
		&camera.Description,
		&camera.CreatedAt,
		&camera.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("摄像头不存在")
		}
		log.Printf("查询摄像头失败: %v", err)
		return nil, err
	}

	return &camera, nil
}

// UpdateCamera 更新摄像头
func (s *CameraService) UpdateCamera(id string, camera *Camera) error {
	query := `
		UPDATE cameras 
		SET name = ?, url = ?, username = ?, password = ?, resolution = ?, fps = ?, status = ?, description = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	result, err := s.db.Exec(query,
		camera.Name,
		camera.URL,
		camera.Username,
		camera.Password,
		camera.Resolution,
		camera.FPS,
		camera.Status,
		camera.Description,
		now,
		id,
	)

	if err != nil {
		log.Printf("更新摄像头失败: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取影响行数失败: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("摄像头不存在")
	}

	log.Printf("摄像头更新成功: %s", camera.Name)
	return nil
}

// DeleteCamera 删除摄像头
func (s *CameraService) DeleteCamera(id string) error {
	query := `DELETE FROM cameras WHERE id = ?`

	result, err := s.db.Exec(query, id)
	if err != nil {
		log.Printf("删除摄像头失败: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取影响行数失败: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("摄像头不存在")
	}

	log.Printf("摄像头删除成功: %s", id)
	return nil
}

// UpdateCameraStatus 更新摄像头状态
func (s *CameraService) UpdateCameraStatus(id string, status string) error {
	query := `UPDATE cameras SET status = ?, updated_at = ? WHERE id = ?`

	now := time.Now()
	_, err := s.db.Exec(query, status, now, id)
	if err != nil {
		log.Printf("更新摄像头状态失败: %v", err)
		return err
	}

	log.Printf("摄像头状态更新成功: %s -> %s", id, status)
	return nil
}

// TestCameraConnection 测试摄像头连接
func (s *CameraService) TestCameraConnection(camera *Camera) error {
	// 这里可以实现实际的摄像头连接测试逻辑
	// 目前返回模拟结果
	log.Printf("测试摄像头连接: %s (%s)", camera.Name, camera.URL)

	// 模拟连接测试延迟
	time.Sleep(1 * time.Second)

	// 模拟连接成功
	return nil
}

// StartCameraStream 开始摄像头流
func (s *CameraService) StartCameraStream(id string) error {
	// 更新摄像头状态为在线
	err := s.UpdateCameraStatus(id, "online")
	if err != nil {
		return err
	}

	log.Printf("摄像头流开始: %s", id)
	return nil
}

// StopCameraStream 停止摄像头流
func (s *CameraService) StopCameraStream(id string) error {
	// 更新摄像头状态为离线
	err := s.UpdateCameraStatus(id, "offline")
	if err != nil {
		return err
	}

	log.Printf("摄像头流停止: %s", id)
	return nil
}

// CreateCameraTable 创建摄像头表
func (s *CameraService) CreateCameraTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS cameras (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			username TEXT,
			password TEXT,
			resolution TEXT DEFAULT '1280x720',
			fps INTEGER DEFAULT 25,
			status TEXT DEFAULT 'offline',
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Printf("创建摄像头表失败: %v", err)
		return err
	}

	log.Println("摄像头表创建成功")
	return nil
}

// InsertDefaultCameras 插入默认摄像头数据
func (s *CameraService) InsertDefaultCameras() error {
	// 检查是否已有摄像头数据
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM cameras").Scan(&count)
	if err != nil {
		log.Printf("查询摄像头数量失败: %v", err)
		return err
	}

	if count > 0 {
		log.Println("摄像头表已有数据，跳过默认数据插入")
		return nil
	}

	// 插入默认摄像头数据
	defaultCameras := []Camera{
		{
			ID:          "camera_1",
			Name:        "机场主摄像头",
			URL:         "rtsp://192.168.1.100:554/stream1",
			Username:    "admin",
			Password:    "admin123",
			Resolution:  "1920x1080",
			FPS:         25,
			Status:      "offline",
			Description: "机场主要监控摄像头",
		},
		{
			ID:          "camera_2",
			Name:        "跑道摄像头",
			URL:         "rtsp://192.168.1.101:554/stream1",
			Username:    "admin",
			Password:    "admin123",
			Resolution:  "1280x720",
			FPS:         30,
			Status:      "offline",
			Description: "跑道监控摄像头",
		},
		{
			ID:          "camera_3",
			Name:        "塔台摄像头",
			URL:         "rtsp://192.168.1.102:554/stream1",
			Username:    "admin",
			Password:    "admin123",
			Resolution:  "1920x1080",
			FPS:         25,
			Status:      "offline",
			Description: "塔台监控摄像头",
		},
	}

	for _, camera := range defaultCameras {
		err := s.CreateCamera(&camera)
		if err != nil {
			log.Printf("插入默认摄像头失败: %v", err)
			continue
		}
	}

	log.Println("默认摄像头数据插入完成")
	return nil
}
