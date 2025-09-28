package database

import (
	"database/sql"
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
	`

	// 执行创建表语句
	if _, err := db.Exec(createMQTTProfilesTable); err != nil {
		return err
	}

	if _, err := db.Exec(createDevicesTable); err != nil {
		return err
	}

	log.Println("Database tables created successfully")
	return nil
}
