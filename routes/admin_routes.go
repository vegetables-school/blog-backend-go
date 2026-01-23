package routes

import (
	"blog/handlers"
	"blog/middleware"

	"github.com/gorilla/mux"
)

// RegisterAdminRoutes 注册后台管理相关路由：需要鉴权的写操作与认证
func RegisterAdminRoutes(r *mux.Router, blogHandler *handlers.BlogHandler, authHandler *handlers.AuthHandler, jwtMiddleware *middleware.JWTMiddleware) {

		// 评论删除接口（仅本人或管理员）
		// 需在主入口注册 CommentHandler 并传入
	}
	// 认证端点（登录/注册）
	r.HandleFunc("/api/admin/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/admin/auth/login", authHandler.Login).Methods("POST")

	// 博客管理端点（需要鉴权的写操作）
	r.HandleFunc("/api/admin/blog", jwtMiddleware.Authenticate(middleware.RequireAdmin(blogHandler.CreateBlog))).Methods("POST")
	r.HandleFunc("/api/admin/blog/{id}", jwtMiddleware.Authenticate(middleware.RequireAdmin(blogHandler.DeleteBlog))).Methods("DELETE")
	r.HandleFunc("/api/admin/blog/{id}", jwtMiddleware.Authenticate(middleware.RequireAdmin(blogHandler.UpdateBlog))).Methods("PUT")

	// 标签管理（需要鉴权）
	r.HandleFunc("/api/admin/tag", jwtMiddleware.Authenticate(middleware.RequireAdmin(blogHandler.AddTagHandler))).Methods("POST")
	r.HandleFunc("/api/admin/tag", jwtMiddleware.Authenticate(middleware.RequireAdmin(blogHandler.DeleteTagHandler))).Methods("DELETE")
}
