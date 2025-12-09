// File: internal/config/config.go

package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env       string // 当前环境名称
	LogLevel  string
	AIBaseURL string
	AIKey     string
	AIModel   string
	MongoURI  string
}

func LoadConfig() *Config {
	// 1. 获取当前环境标识，默认为 "dev"
	env := getEnv("APP_ENV", "dev")

	// 2. 拼接配置文件名，例如 ".env.dev"
	envFile := ".env." + env

	// 3. 尝试加载对应的配置文件
	// 如果文件不存在，我们只打印警告，不报错退出。
	// 因为在 Docker/K8s 生产环境中，我们通常直接通过系统环境变量注入配置，而不使用 .env 文件。
	if err := godotenv.Load(envFile); err != nil {
		slog.Warn("Config file not found, relying on system env vars", "file", envFile)
	} else {
		slog.Info("Config loaded from file", "file", envFile)
	}

	return &Config{
		Env:       env,
		AIBaseURL: getEnv("AI_BASE_URL", "https://api.openai.com/v1"),
		AIKey:     getEnv("AI_API_KEY", ""),
		AIModel:   getEnv("AI_MODEL", "gpt-3.5-turbo"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		MongoURI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
