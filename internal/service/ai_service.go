// File: internal/service/ai_service.go

package service

import (
	"context"
	"log/slog"
	"time"

	"chatbot/internal/config"

	"github.com/sashabaranov/go-openai"
)

// aiService 是 ChatService 的 AI 实现版本
type aiService struct {
	client *openai.Client
	model  string
}

// GetResponse 调用 LLM API 获取回复
func (s *aiService) GetResponse(message string) string {
	// 1. 创建上下文，设置 30 秒超时，防止 AI 卡死导致服务器资源耗尽
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// 2. 构建请求
	req := openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
	}

	slog.Info("Sending request to AI", "model", s.model)

	// 3. 发起调用
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		slog.Error("AI API call failed", "error", err)
		return "Sorry, I'm having trouble thinking right now. Please try again later."
	}

	// 4. 提取回复
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}

	return "I received an empty response."
}

// NewAIService 初始化 AI 服务
func NewAIService(cfg *config.Config) ChatService {
	// 配置 OpenAI 客户端
	clientConfig := openai.DefaultConfig(cfg.AIKey)
	clientConfig.BaseURL = cfg.AIBaseURL

	client := openai.NewClientWithConfig(clientConfig)

	return &aiService{
		client: client,
		model:  cfg.AIModel,
	}
}
