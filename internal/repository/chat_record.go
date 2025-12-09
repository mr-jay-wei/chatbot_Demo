// File: internal/repository/chat_record.go

package repository

import (
	"time"
)

// ChatRecord 代表一条聊天记录
type ChatRecord struct {
	UserMessage string    `bson:"user_message"` // 用户说的话
	AIMessage   string    `bson:"ai_message"`   // AI 回复的话
	CreatedAt   time.Time `bson:"created_at"`   // 时间
}
