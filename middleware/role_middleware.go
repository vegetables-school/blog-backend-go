package middleware

import (
	"net/http"
)

// RequireAdmin 只允许管理员访问的中间件
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := GetUserRole(r)
		if role != "admin" {
			http.Error(w, "无权限，管理员专用接口", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
