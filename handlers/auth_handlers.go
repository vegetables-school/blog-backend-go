package handlers

import (
	"encoding/json"
	"net/http"

	"blog/models"
	"blog/services"
	"blog/utils"
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
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body models.UserRegisterRequest true "注册信息"
// @Success 201 {object} AuthUserResponse
// @Failure 400 {string} string "无效的请求数据"
// @Failure 403 {string} string "注册功能暂时关闭"
// @Router /api/admin/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 基本验证
	if req.Username == "" || req.Password == "" || req.Email == "" {
		utils.SendError(w, http.StatusBadRequest, "用户名、密码和邮箱不能为空")
		return
	}

	// 密码强度验证
	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 邮箱格式验证
	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendJSON(w, http.StatusCreated, user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body models.UserLoginRequest true "登录信息"
// @Success 200 {object} LoginResponse
// @Failure 400 {string} string "无效的请求数据"
// @Router /api/admin/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 基本验证
	if req.Username == "" || req.Password == "" {
		utils.SendError(w, http.StatusBadRequest, "用户名和密码不能为空")
		return
	}

	authResponse, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	loginResp := LoginResponse{
		Data: LoginData{
			Token: authResponse.Token,
			User:  &authResponse.User,
		},
	}
	utils.SendJSON(w, http.StatusOK, loginResp)
}
