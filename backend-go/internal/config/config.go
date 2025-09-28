package config

import (
	"os"
)

type Config struct {
	Environment  string
	DatabasePath string
	Port         string
}

func Load() *Config {
	return &Config{
		Environment:  getEnv("ENV", "development"),
		DatabasePath: getEnv("DATABASE_PATH", "./data/backend.db"),
		Port:         getEnv("PORT", "18080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
