package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 表示用户
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"` // 密码哈希，不在 JSON 中返回
	Email     string             `bson:"email" json:"email"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserLoginRequest 登录请求
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserRegisterRequest 注册请求
type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
