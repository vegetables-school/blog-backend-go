package routes

import (
	"blog/handlers"
	"blog/middleware"

	"github.com/gorilla/mux"
)

// RegisterRoutes 聚合调用前端(public)与后台(admin)路由注册，保持向后兼容
func RegisterRoutes(r *mux.Router, blogHandler *handlers.BlogHandler, authHandler *handlers.AuthHandler, jwtMiddleware *middleware.JWTMiddleware) {
	RegisterPublicRoutes(r, blogHandler, authHandler)
	RegisterFrontRoutes(r, authHandler)
	RegisterAdminRoutes(r, blogHandler, authHandler, jwtMiddleware)
}
