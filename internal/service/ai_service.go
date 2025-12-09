// File: internal/service/ai_service.go

package service

import (
	"context"
	"log/slog"
	"time"

	"chatbot/internal/config"
	"chatbot/internal/repository" // 引入 repository 包

	"github.com/sashabaranov/go-openai"
)

type aiService struct {
	client *openai.Client
	model  string
	repo   repository.ChatRepository // 新增：持有 Repo 的引用
}

// 修改构造函数，传入 repo
func NewAIService(cfg *config.Config, repo repository.ChatRepository) ChatService {
	clientConfig := openai.DefaultConfig(cfg.AIKey)
	clientConfig.BaseURL = cfg.AIBaseURL
	client := openai.NewClientWithConfig(clientConfig)

	return &aiService{
		client: client,
		model:  cfg.AIModel,
		repo:   repo, // 注入 repo
	}
}

func (s *aiService) GetResponse(message string) string {
	// 1. 调用 AI (保持原有逻辑)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: message},
		},
	}

	slog.Info("Sending request to AI", "model", s.model)

	var aiResponse string
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		slog.Error("AI API call failed", "error", err)
		aiResponse = "Sorry, I'm having trouble thinking right now."
	} else if len(resp.Choices) > 0 {
		aiResponse = resp.Choices[0].Message.Content
	} else {
		aiResponse = "I received an empty response."
	}

	// 2. 异步保存到数据库 (新增逻辑)
	// 我们使用 go func 开启一个协程来保存，这样不会阻塞给用户的回复速度
	go func(userMsg, aiMsg string) {
		// 创建一个独立的上下文，超时 5 秒
		saveCtx, saveCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer saveCancel()

		record := &repository.ChatRecord{
			UserMessage: userMsg,
			AIMessage:   aiMsg,
			CreatedAt:   time.Now(),
		}

		if err := s.repo.Save(saveCtx, record); err != nil {
			// 这里的错误只记录日志，不影响主流程
			slog.Error("Async save failed", "error", err)
		} else {
			slog.Debug("Chat record saved to MongoDB")
		}
	}(message, aiResponse)

	return aiResponse
}
