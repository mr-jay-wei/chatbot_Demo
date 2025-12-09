// File: cmd/chatbot/main.go

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"chatbot/internal/config"
	"chatbot/internal/handler"
	"chatbot/internal/repository" // 引入 repository
	"chatbot/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 1. 加载配置
	cfg := config.LoadConfig()
	setupLogger(cfg.LogLevel, cfg.Env)

	// 2. 连接 MongoDB (新增步骤)
	slog.Info("Connecting to MongoDB...")
	// 设置连接超时 10 秒
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 连接数据库
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		slog.Error("Failed to create MongoDB client", "error", err)
		os.Exit(1)
	}

	// 检查连接是否成功 (Ping)
	if err := client.Ping(ctx, nil); err != nil {
		slog.Error("Failed to ping MongoDB", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to MongoDB successfully")

	// 程序退出时断开数据库连接
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			slog.Error("Failed to disconnect MongoDB", "error", err)
		}
	}()

	// 获取数据库实例 (数据库名 libraryDB 来自你的连接串)
	db := client.Database("libraryDB")

	// 3. 初始化层级依赖
	// Layer 1: Repository (数据层)
	chatRepo := repository.NewMongoChatRepo(db)

	// Layer 2: Service (业务层) - 注入 Repo
	if cfg.AIKey == "" {
		slog.Warn("AI_API_KEY is not set.")
	}
	chatSvc := service.NewAIService(cfg, chatRepo)

	// Layer 3: Handler (接口层)
	chatHdl := handler.NewChatHandler(chatSvc)

	// 4. 配置路由与服务器
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", chatHdl.HandleChat)

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
