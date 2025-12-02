// File: cmd/chatbot/main.go

package main

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"chatbot/internal/config"
	"chatbot/internal/handler"
	"chatbot/internal/service"
)

func main() {
	// 1. 先加载配置 (Config First)
	// 注意：此时还没有配置好的 Logger，所以 LoadConfig 内部只能用默认 Logger
	cfg := config.LoadConfig()

	// 2. 初始化日志系统 (Logger Setup)
	setupLogger(cfg.LogLevel, cfg.Env)

	// 3. 初始化服务
	if cfg.AIKey == "" {
		slog.Warn("AI_API_KEY is not set. Chatbot might not work correctly.")
	}

	slog.Info("Initializing Service", "env", cfg.Env, "log_level", cfg.LogLevel)

	// 使用 AI 服务
	chatSvc := service.NewAIService(cfg)
	chatHdl := handler.NewChatHandler(chatSvc)

	// 4. 配置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", chatHdl.HandleChat)

	// 5. 启动服务器
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 100 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("Starting HTTP server", "addr", ":8080")
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// setupLogger 根据配置初始化全局 Logger
func setupLogger(levelStr string, env string) {
	// 1. 解析日志级别字符串 (e.g., "debug" -> slog.LevelDebug)
	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // 默认 fallback
	}

	// 2. 配置 Handler 选项
	opts := &slog.HandlerOptions{
		Level: level, // 设置级别
		// AddSource: true, // 如果设为 true，日志会显示文件名和行号（生产环境通常关闭以节省性能）
	}

	// 3. 根据环境选择格式
	var handler slog.Handler
	if env == "prod" {
		// 生产环境：使用 JSON 格式，方便机器解析 (ELK, Datadog)
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		// 开发环境：使用文本格式，方便人类阅读
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	// 4. 设置为全局默认 Logger
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
