// @title 博客后端 API 文档
// @version 1.0
// @description 这是一个用 Go 编写的博客后端 API 示例，支持用户、博客、评论、点赞等功能。
// @contact.name API Support
// @contact.email youremail@example.com
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"blog/config"
	"blog/handlers"
	"blog/middleware"
	"blog/routes"
	"blog/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	// 匿名 import handlers，确保 dummy handler 被 swag 收录
	_ "blog/handlers"
)

// 服务器配置常量
const (
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
	idleTimeout     = 60 * time.Second
	shutdownTimeout = 10 * time.Second
)

func main() {
	// 加载挂载的 .env 文件（仅开发环境使用，加载失败不影响程序运行）
	if err := godotenv.Load("./.env"); err != nil {
		log.Printf("提示: 未找到或无法加载 .env 文件，将使用环境变量或默认值")
	}

	// 脱敏打印 MongoDB 配置信息
	mongoURI := getEnv("MONGODB_URI", "")
	if mongoURI != "" {
		if len(mongoURI) > 20 {
			mongoURI = mongoURI[:20] + "***"
		}
	}
	log.Printf("MongoDB 配置: URI=%s, 数据库=%s, 集合=%s\n", mongoURI, os.Getenv("MONGODB_DATABASE"), os.Getenv("COLLECTION_NAME"))

	// 初始化 MongoDB（自动配置）
	mongoCfg, err := config.InitMongo()
	if err != nil {
		log.Fatal("MongoDB 初始化失败:", err)
	}
	log.Println("成功连接到 MongoDB")

	// 初始化服务
	blogService := services.NewBlogService(mongoCfg.Client, mongoCfg.Database, getEnv("COLLECTION_NAME", "blogs-dev"))
	commentService := services.NewCommentService(mongoCfg.Client, mongoCfg.Database, "comments")
	likeService := services.NewLikeService(mongoCfg.Client, mongoCfg.Database, "likes")

	// 初始化认证服务
	jwtSecret := getEnv("JWT_SECRET", "default-jwt-secret")
	authService := services.NewAuthService(mongoCfg.Client, mongoCfg.Database, "users", jwtSecret)

	// 安全检查：警告生产环境使用了不安全的默认 JWT 密钥
	if jwtSecret == "default-jwt-secret" {
		log.Printf("警告: 正在使用默认的 JWT 密钥，请在生产环境通过 JWT_SECRET 环境变量设置安全的密钥")
	}

	// 初始化处理器
	blogHandler := handlers.NewBlogHandler(blogService)
	authHandler := handlers.NewAuthHandler(authService)
	commentHandler := handlers.NewCommentHandler(commentService)
	likeHandler := handlers.NewLikeHandler(likeService)

	// 初始化中间件
	jwtMiddleware := middleware.NewJWTMiddleware(authService)

	// 创建路由
	r := mux.NewRouter()

	// 注册 /docs/ 静态文件路由（用于 swagger.json, swagger.yaml 等）
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	// 注册 swagger-ui 路由
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
		httpSwagger.DocExpansion("list"),
	))

	// 注册路由（集中管理）
	routes.RegisterRoutes(r, blogHandler, authHandler, jwtMiddleware, commentHandler, likeHandler)

	// 健康检查端点
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET", "HEAD")

	// 启动服务器（带超时配置）
	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  readTimeout,  // 读取请求超时
		WriteTimeout: writeTimeout, // 写入响应超时
		IdleTimeout:  idleTimeout,  // 空闲连接超时
	}

	// 启动服务器（监听信号实现优雅关闭）
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("博客服务器启动在 http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// 监听中断信号，实现优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 等待信号或服务器错误
	select {
	case <-sigChan:
		log.Println("收到关闭信号，开始优雅关闭...")
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Printf("服务器启动失败: %v", err)
			return
		}
	}

	// 关闭 HTTP 服务器
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("服务器关闭失败: %v", err)
	}

	// 显式关闭 MongoDB 连接
	if err := config.CloseMongo(); err != nil {
		log.Printf("关闭 MongoDB 连接失败: %v", err)
	} else {
		log.Println("MongoDB 连接已关闭")
	}

	log.Println("服务器已关闭")

}

// getEnv 从环境变量读取值，若为空则返回默认值
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
