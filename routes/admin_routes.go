package routes

import (
	"blog/handlers"
	"blog/middleware"

	"github.com/gorilla/mux"
)

// RegisterAdminRoutes 注册后台管理相关路由：需要鉴权的写操作与认证
func RegisterAdminRoutes(r *mux.Router, blogHandler *handlers.BlogHandler, authHandler *handlers.AuthHandler, jwtMiddleware *middleware.JWTMiddleware) {
	// 认证端点（登录/注册）
	r.HandleFunc("/api/admin/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/admin/auth/login", authHandler.Login).Methods("POST")

	// 博客管理端点（需要鉴权的写操作）
	r.HandleFunc("/api/admin/blog", jwtMiddleware.Authenticate(blogHandler.CreateBlog)).Methods("POST")
	r.HandleFunc("/api/admin/blog/{id}", jwtMiddleware.Authenticate(blogHandler.UpdateBlog)).Methods("PUT")
	r.HandleFunc("/api/admin/blog/{id}", jwtMiddleware.Authenticate(blogHandler.DeleteBlog)).Methods("DELETE")
}
