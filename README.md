
# Go AI Chatbot (Production Grade)

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

这是一个基于 Go 语言构建的、生产级别的 AI 聊天机器人后端服务。本项目旨在展示如何从零开始构建一个高可扩展、无状态、分层架构的现代 Go 应用。

它不仅实现了与大语言模型（LLM, 如 DeepSeek/OpenAI）的对接，还包含了一套完整的工程化最佳实践，包括依赖注入、结构化日志、多环境配置管理和单元测试。

## 🏗 核心架构与设计原则

本项目严格遵循 **Standard Go Project Layout** 和 **Clean Architecture** 原则：

*   **分层架构 (Layered Architecture)**:
    *   **Handler 层 (`internal/handler`)**: 负责 HTTP 请求的解析、验证和响应构建。不包含业务逻辑。
    *   **Service 层 (`internal/service`)**: 核心业务逻辑层。定义了 `ChatService` 接口，实现了业务解耦。
    *   **Config 层 (`internal/config`)**: 负责环境感知和配置加载。
*   **依赖注入 (Dependency Injection)**: 通过接口 (`Interface`) 隔离依赖，使得从 "Echo 模式" 切换到 "AI 模式" 无需修改调用方代码，同时也便于单元测试。
*   **无状态设计 (Stateless)**: 服务本身不保存会话状态，符合云原生 (Cloud Native) 扩展要求。
*   **可观测性 (Observability)**: 集成 `log/slog` 实现结构化日志，支持不同环境（Dev/Prod）的日志级别和格式切换。

## 🚀 开发历程 (Development Journey)

本项目经历了从基础骨架到完整 AI 服务的演进过程：

### Phase 1: 基础设施搭建
1.  **项目初始化**: 使用 `go mod` 管理依赖，建立标准的 `cmd/`, `internal/`, `pkg/` 目录结构。
2.  **核心抽象**: 定义 `ChatService` 接口，确立了面向接口编程的基础。
3.  **CLI 原型**: 最初实现了一个命令行交互版本，用于快速验证核心逻辑。

### Phase 2: 服务化与工程化
4.  **Web 服务化**: 引入 `net/http`，将 CLI 改造为 RESTful API 服务，实现了 HTTP Handler。
5.  **结构化日志**: 引入 Go 1.21+ 标准库 `slog`，替代 `fmt.Println`，实现了日志与业务数据的分离。
6.  **单元测试**: 采用表格驱动测试 (Table-Driven Tests) 模式，确保核心逻辑的稳定性。

### Phase 3: 智能化与生产就绪
7.  **配置管理**: 引入 `godotenv` 和多环境配置策略 (`.env.dev`, `.env.prod`)，通过 `APP_ENV` 环境变量动态切换。
8.  **AI 集成**: 接入 `go-openai` SDK，实现了真正的 LLM 调用能力，并配置了超时控制 (`Context Timeout`) 以防止服务雪崩。
9.  **安全加固**: 通过 `.gitignore` 策略严格防止 API Key 泄露。

## 📂 项目结构

```text
chatbot/
├── cmd/
│   └── chatbot/
│       └── main.go           # 程序入口，负责依赖组装和服务器启动
├── internal/
│   ├── config/               # 配置管理 (Env, Flags)
│   ├── handler/              # HTTP 处理器 (Request/Response)
│   └── service/              # 业务逻辑 (AI Integration, Echo Logic)
├── .env.dev                  # 开发环境配置 (Git 忽略)
├── .env.prod                 # 生产环境配置 (Git 忽略)
├── .gitignore                # Git 忽略规则
├── go.mod                    # 依赖定义
└── README.md                 # 项目文档
```

## 🛠️ 快速开始 (Getting Started)

### 1. 环境要求
*   Go 1.21 或更高版本
*   Git

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 配置环境
在项目根目录创建 `.env.dev` 文件（参考 `env.example`）：

```env
AI_BASE_URL=https://api.deepseek.com/v1
AI_API_KEY=sk-your-api-key-here
AI_MODEL=deepseek-chat
LOG_LEVEL=debug
```

### 4. 运行项目

**开发环境 (默认):**
```bash
go run ./cmd/chatbot
```

**测试环境：**
```powershell
# PowerShell
$env:APP_ENV="test"; go run ./cmd/chatbot
```

**生产环境：**
```powershell
# PowerShell
$env:APP_ENV="prod"; go run ./cmd/chatbot
```

### 5. API 测试
启动服务后，发送 POST 请求：

```bash
curl -X POST http://localhost:8080/chat \
     -H "Content-Type: application/json" \
     -d '{"message": "你好，Go语言有什么特点？"}'
```

## 🧪 测试

运行单元测试：
```bash
go test -v ./internal/service
```

## 🔮 未来规划 (Roadmap)

*   [ ] **Docker 化**: 编写 Dockerfile 和 docker-compose.yml。
*   [ ] **流式响应 (Streaming)**: 支持 Server-Sent Events (SSE) 以实现打字机效果。
*   [ ] **上下文记忆**: 引入 Redis 存储对话历史，实现多轮对话。
*   [ ] **API 网关**: 引入 Gin 或 Echo 框架以支持更复杂的路由和中间件。

---
*Built with ❤️ by a Go Developer.*