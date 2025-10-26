package models

import "time"

// Blog 表示一篇博客文章
type Blog struct {
	ID        int       `json:"id"`        // 博客文章的唯一标识符
	Title     string    `json:"title"`     // 博客文章的标题
	Content   string    `json:"content"`   // 博客文章的内容
	Author    string    `json:"author"`    // 博客文章的作者
	CreatedAt time.Time `json:"created_at"` // 博客文章的创建时间
	UpdatedAt time.Time `json:"updated_at"` // 博客文章的更新时间
}