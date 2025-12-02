// File: internal/handler/chat_handler.go

package handler

import (
	"encoding/json"
	"net/http"

	"chatbot/internal/service"
	"log/slog"
)

// ChatRequest 定义了客户端发送过来的 JSON 数据结构
// `json:"message"` 告诉 Go：当收到 JSON 里的 "message" 字段时，把它填到 Message 变量里
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse 定义了我们要返回给客户端的 JSON 数据结构
type ChatResponse struct {
	Response string `json:"response"`
}

// ChatHandler 是我们的 HTTP 处理器
// 它持有 ChatService 的引用，这样它就能调用核心逻辑了
type ChatHandler struct {
	svc service.ChatService
}

// HandleChat 是真正处理 HTTP 请求的方法
// 它符合 http.HandlerFunc 的标准签名
func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Received raw request", "remote_addr", r.RemoteAddr, "method", r.Method)
	// 1. 只允许 POST 方法
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. 解析请求体 (JSON -> Struct)
	var req ChatRequest
	//NewDecoder 用于从请求体中读取并解码 JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. 调用业务服务 (Service Layer)
	// 注意：Handler 层不负责业务逻辑，只负责搬运数据
	responseMsg := h.svc.GetResponse(req.Message)

	// 4. 构建响应 (Struct -> JSON)
	resp := ChatResponse{
		Response: responseMsg,
	}

	// 5. 发送响应
	w.Header().Set("Content-Type", "application/json") // 告诉客户端返回的是 JSON
	w.WriteHeader(http.StatusOK)                       // 设置状态码 200 OK
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode response", "error", err)
		// 这里通常不需要再写 http.Error，因为 Header 已经发出去了
	}
}

// NewChatHandler 是一个构造函数，用于创建 ChatHandler 实例
func NewChatHandler(svc service.ChatService) *ChatHandler {
	return &ChatHandler{
		svc: svc,
	}
}
