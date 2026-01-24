package routes

import (
	"blog/handlers"
	"blog/middleware"

	"github.com/gorilla/mux"
)

// RegisterPublicRoutes 注册前端可访问的公开路由：内容读取
func RegisterPublicRoutes(
	r *mux.Router,
	blogHandler *handlers.BlogHandler,
	_ *handlers.AuthHandler,
	_ *middleware.JWTMiddleware,
	_ *handlers.CommentHandler,
	_ *handlers.LikeHandler,
) {
	// 只保留真正公开的接口
	r.HandleFunc("/api/blog/{id}", blogHandler.GetBlog).Methods("GET")
	r.HandleFunc("/api/blogs", blogHandler.GetBlogsPaginated).Methods("GET")
	r.HandleFunc("/api/blogs/search", blogHandler.SearchBlogsHandler).Methods("GET")
	r.HandleFunc("/api/tags", blogHandler.GetAllTagsHandler).Methods("GET")
	r.HandleFunc("/api/blogs/by-tags", blogHandler.GetBlogsByTagsHandler).Methods("GET")
}
