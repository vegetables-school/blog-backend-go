package handlers

import (
	"blog/models"
)

// Pagination 分页信息
type Pagination struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
	Total int64 `json:"total"`
}

// BlogListResponse 公共博客列表响应，主数据在 data 字段
type BlogListResponse struct {
	Data       []*models.Blog `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

// BlogResponse 单篇博客响应
type BlogResponse struct {
	Data *models.Blog `json:"data"`
}

// AuthUserResponse 注册返回用户数据
type AuthUserResponse struct {
	Data *models.User `json:"data"`
}

// LoginResponse 为登录返回的结构（token + user）
type LoginResponse struct {
	Data struct {
		Token string       `json:"token"`
		User  *models.User `json:"user"`
	} `json:"data"`
}
