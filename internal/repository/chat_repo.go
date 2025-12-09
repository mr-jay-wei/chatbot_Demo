// File: internal/repository/chat_repo.go

package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
)

// ChatRepository 定义了操作聊天记录的接口
// 这样做的好处是：以后如果你想换成 MySQL，只需要实现这个接口，不用改 Service 层代码
type ChatRepository interface {
	Save(ctx context.Context, record *ChatRecord) error
}

// mongoChatRepo 是具体的 MongoDB 实现
type mongoChatRepo struct {
	collection *mongo.Collection
}

// NewMongoChatRepo 创建一个 MongoDB 仓储实例
func NewMongoChatRepo(db *mongo.Database) ChatRepository {
	// 我们把记录存在 "chat_history" 集合（表）中
	return &mongoChatRepo{
		collection: db.Collection("chat_history"),
	}
}

// Save 将记录保存到 MongoDB
func (r *mongoChatRepo) Save(ctx context.Context, record *ChatRecord) error {
	_, err := r.collection.InsertOne(ctx, record)
	if err != nil {
		slog.Error("Failed to save chat record to MongoDB", "error", err)
		return err
	}
	return nil
}
