package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog 表示一篇博客文章
type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 博客文章的唯一标识符
	Title     string             `bson:"title" json:"title"`                 // 博客文章的标题
	Content   string             `bson:"content" json:"content"`             // 博客文章的内容
	Author    string             `bson:"author" json:"author"`               // 博客文章的作者
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`       // 博客文章的创建时间
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`       // 博客文章的更新时间
}