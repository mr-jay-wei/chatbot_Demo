// File: internal/service/chat_service.go

package service

import (
	"log/slog" // 引入标准日志库
)

// ChatService defines the interface for the core chat logic.
type ChatService interface {
	GetResponse(message string) string
}

// ============echoService==================
type echoService struct{}

// GetResponse implements the ChatService interface.
func (s *echoService) GetResponse(message string) string {
	// 使用结构化日志记录业务行为
	// 这里的 "content" 是一个键值对，方便后续查询
	slog.Info("Processing message", "message_length", len(message), "content", message)

	return "You said: " + message
}

// ===========NewChatService=================
func NewChatService() ChatService {
	return &echoService{}
}
