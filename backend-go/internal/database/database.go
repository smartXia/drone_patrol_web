package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func InitDB(dbPath string) (*DB, error) {
	// 确保数据目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// 打开数据库连接
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 创建表
	if err := createTables(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createTables(db *sql.DB) error {
	// 创建MQTT配置表
	createMQTTProfilesTable := `
	CREATE TABLE IF NOT EXISTS mqtt_profiles (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		config TEXT NOT NULL,
		is_default INTEGER DEFAULT 0,
		updated_at INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_mqtt_profiles_default ON mqtt_profiles(is_default);
	`

	// 创建设备表
	createDevicesTable := `
	CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		sn TEXT UNIQUE NOT NULL,
		type TEXT NOT NULL DEFAULT 'drone',
		status TEXT NOT NULL DEFAULT 'offline',
		airport_sn TEXT DEFAULT '',
		last_seen INTEGER,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		is_current INTEGER DEFAULT 0,
		is_gateway INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_devices_sn ON devices(sn);
	CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
	CREATE INDEX IF NOT EXISTS idx_devices_current ON devices(is_current);
	CREATE INDEX IF NOT EXISTS idx_devices_gateway ON devices(is_gateway);
	CREATE INDEX IF NOT EXISTS idx_devices_airport_sn ON devices(airport_sn);
	`

	// 执行创建表语句
	if _, err := db.Exec(createMQTTProfilesTable); err != nil {
		return err
	}

	if _, err := db.Exec(createDevicesTable); err != nil {
		return err
	}

	// 检查并添加 airport_sn 字段到现有表
	if err := addAirportSnColumnIfNotExists(db); err != nil {
		log.Printf("Airport SN column migration failed: %v", err)
		return err
	}

	log.Println("Database tables created successfully")
	return nil
}

// 检查并添加 airport_sn 列（如果不存在）
func addAirportSnColumnIfNotExists(db *sql.DB) error {
	// 检查列是否已存在
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('devices') WHERE name='airport_sn'").Scan(&count)
	if err != nil {
		return fmt.Errorf("检查列是否存在失败: %v", err)
	}

	// 如果列不存在，则添加
	if count == 0 {
		log.Println("添加 airport_sn 列到 devices 表")
		_, err := db.Exec("ALTER TABLE devices ADD COLUMN airport_sn TEXT DEFAULT ''")
		if err != nil {
			return fmt.Errorf("添加 airport_sn 列失败: %v", err)
		}
	} else {
		log.Println("airport_sn 列已存在，跳过添加")
	}

	// 创建索引（如果不存在）
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_devices_airport_sn ON devices(airport_sn)")
	if err != nil {
		return fmt.Errorf("创建 airport_sn 索引失败: %v", err)
	}

	return nil
}
