package routes

import (
	"blog/handlers"
	"blog/middleware"

	"github.com/gorilla/mux"
)

// RegisterRoutes 将所有路由集中注册到传入的 mux.Router
func RegisterRoutes(r *mux.Router, blogHandler *handlers.BlogHandler, authHandler *handlers.AuthHandler, jwtMiddleware *middleware.JWTMiddleware) {

	// 认证端点
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// 定义 API 端点
	r.HandleFunc("/api/blog", jwtMiddleware.Authenticate(blogHandler.CreateBlog)).Methods("POST")
	r.HandleFunc("/api/blog/{id}", jwtMiddleware.Authenticate(blogHandler.UpdateBlog)).Methods("PUT")
	r.HandleFunc("/api/blog/{id}", jwtMiddleware.Authenticate(blogHandler.DeleteBlog)).Methods("DELETE")
	r.HandleFunc("/api/blog/{id}", blogHandler.GetBlog).Methods("GET")
}
