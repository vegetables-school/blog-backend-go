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
	jwtMiddleware *middleware.JWTMiddleware,
	commentHandler *handlers.CommentHandler,
	likeHandler *handlers.LikeHandler,
) {
	// 评论接口（需登录）
	r.HandleFunc("/api/comment", jwtMiddleware.Authenticate(commentHandler.CreateCommentHandler)).Methods("POST")
	r.HandleFunc("/api/comment", jwtMiddleware.Authenticate(middleware.OnlySelfOrAdmin(commentHandler.GetCommentUserID)(commentHandler.DeleteCommentHandler))).Methods("DELETE")

	// 点赞接口（需登录）
	r.HandleFunc("/api/like", jwtMiddleware.Authenticate(likeHandler.AddLikeHandler)).Methods("POST")
	r.HandleFunc("/api/like", jwtMiddleware.Authenticate(middleware.OnlySelfOrAdmin(likeHandler.GetLikeUserID)(likeHandler.RemoveLikeHandler))).Methods("DELETE")

	// 获取单篇博客（公开访问）
	r.HandleFunc("/api/blog/{id}", blogHandler.GetBlog).Methods("GET")
	// 分页获取博客列表（公开访问）
	r.HandleFunc("/api/blogs", blogHandler.GetBlogsPaginated).Methods("GET")

	// 博客模糊搜索与标签筛选（公开访问）
	r.HandleFunc("/api/blogs/search", blogHandler.SearchBlogsHandler).Methods("GET")

	// 获取所有标签（公开访问）
	r.HandleFunc("/api/tags", blogHandler.GetAllTagsHandler).Methods("GET")

	// 分页获取标签下的文章（公开访问）
	r.HandleFunc("/api/blogs/by-tags", blogHandler.GetBlogsByTagsHandler).Methods("GET")
}
