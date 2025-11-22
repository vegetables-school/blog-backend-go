package handlers

import (
	"encoding/json"
	"net/http"

	"blog/models"
	"blog/services"
)

// AuthHandler 处理用户认证的HTTP请求
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler 创建新的AuthHandler实例
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// 暂时禁止注册功能
	// http.Error(w, "注册功能暂时关闭", http.StatusForbidden)
	// return

	var req models.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 基本验证
	if req.Username == "" || req.Password == "" || req.Email == "" {
		http.Error(w, "用户名、密码和邮箱不能为空", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "密码长度至少6位", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthUserResponse{Data: user})
}

// Login 用户登录
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 基本验证
	if req.Username == "" || req.Password == "" {
		http.Error(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}

	authResponse, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var loginResp LoginResponse
	loginResp.Data.Token = authResponse.Token
	loginResp.Data.User = &authResponse.User
	json.NewEncoder(w).Encode(loginResp)
}
