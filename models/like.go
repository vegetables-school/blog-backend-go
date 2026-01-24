package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Like 点赞结构
type Like struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BlogID    primitive.ObjectID `bson:"blog_id" json:"blog_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// AddLikeRequest 点赞请求
type AddLikeRequest struct {
	BlogID string `json:"blog_id"`
}

// RemoveLikeRequest 取消点赞请求
type RemoveLikeRequest struct {
	BlogID string `json:"blog_id"`
}
