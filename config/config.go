package config

import (
	"os"
	"strconv"
)

// Config 配置结构体
type Config struct {
	API APIConfig
}

// APIConfig API 配置
type APIConfig struct {
	BaseURL string
	Timeout int
}

// Load 加载配置（从环境变量或默认值）
func Load() *Config {
	return &Config{
		API: APIConfig{
			BaseURL: getEnv("LEISTER_API_URL", "http://localhost:8080"),
			Timeout: getEnvInt("LEISTER_API_TIMEOUT", 30),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取环境变量（整数类型），如果不存在或解析失败则返回默认值
func getEnvInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
