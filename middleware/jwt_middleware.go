package middleware

import (
	"context"
	"net/http"
	"strings"

	"blog/services"

	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware JWT 认证中间件
type JWTMiddleware struct {
	authService *services.AuthService
}

// NewJWTMiddleware 创建新的JWTMiddleware实例
func NewJWTMiddleware(authService *services.AuthService) *JWTMiddleware {
	return &JWTMiddleware{
		authService: authService,
	}
}

// Authenticate 验证 JWT token 的中间件
func (m *JWTMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "缺少认证令牌", http.StatusUnauthorized)
			return
		}

		// 提取 Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "无效的认证令牌格式", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// 验证 token
		token, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			// 检查是否是token过期错误
			if strings.Contains(err.Error(), "token is expired") {
				http.Error(w, "认证令牌已过期，请重新登录", http.StatusUnauthorized)
				return
			}
			http.Error(w, "无效的认证令牌: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// 从 token 中提取用户信息
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["user_id"].(string)
			username := claims["username"].(string)

			// 将用户信息添加到请求上下文
			ctx := context.WithValue(r.Context(), "user_id", userID)
			ctx = context.WithValue(ctx, "username", username)
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "无效的认证令牌", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// GetUserID 从请求上下文中获取用户ID
func GetUserID(r *http.Request) string {
	if userID, ok := r.Context().Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetUsername 从请求上下文中获取用户名
func GetUsername(r *http.Request) string {
	if username, ok := r.Context().Value("username").(string); ok {
		return username
	}
	return ""
}
