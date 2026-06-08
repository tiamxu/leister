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
			BaseURL: mustGetEnv("LEISTER_API_URL"),
			Timeout: getEnvInt("LEISTER_API_TIMEOUT", 30),
		},
	}
}

// mustGetEnv 获取必填环境变量，缺失则 panic
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("环境变量 " + key + " 未设置，请检查配置")
	}
	return value
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
