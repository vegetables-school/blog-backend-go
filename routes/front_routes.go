package routes

import (
	"blog/handlers"
	"blog/middleware"

	"github.com/gorilla/mux"
)

// RegisterFrontRoutes 注册普通用户专属接口（需登录）
func RegisterFrontRoutes(r *mux.Router, authHandler *handlers.AuthHandler, jwtMiddleware *middleware.JWTMiddleware, commentHandler *handlers.CommentHandler, likeHandler *handlers.LikeHandler) {
	// 评论接口（需登录）
	r.HandleFunc("/api/comment", jwtMiddleware.Authenticate(commentHandler.CreateCommentHandler)).Methods("POST")
	r.HandleFunc("/api/comment", jwtMiddleware.Authenticate(commentHandler.DeleteCommentHandler)).Methods("DELETE")

	// 点赞接口（需登录）
	r.HandleFunc("/api/like", jwtMiddleware.Authenticate(likeHandler.AddLikeHandler)).Methods("POST")
	r.HandleFunc("/api/like", jwtMiddleware.Authenticate(likeHandler.RemoveLikeHandler)).Methods("DELETE")

	// 用户个人信息（需登录）
	// r.HandleFunc("/api/user/profile", jwtMiddleware.Authenticate(authHandler.GetProfile)).Methods("GET")
	// r.HandleFunc("/api/user/update", jwtMiddleware.Authenticate(authHandler.UpdateProfile)).Methods("PUT")
	// 可继续添加更多普通用户专属接口
}
