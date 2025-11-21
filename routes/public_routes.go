package routes

import (
	"blog/handlers"

	"github.com/gorilla/mux"
)

// RegisterPublicRoutes 注册前端可访问的公开路由：内容读取
func RegisterPublicRoutes(r *mux.Router, blogHandler *handlers.BlogHandler, _ *handlers.AuthHandler) {
	// 获取单篇博客（公开访问）
	r.HandleFunc("/api/blog/{id}", blogHandler.GetBlog).Methods("GET")
}
