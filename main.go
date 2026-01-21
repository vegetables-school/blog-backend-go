package main

import (
	"blog/config"
	"blog/handlers"
	"blog/middleware"
	"blog/routes"
	"blog/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// 加载挂载的 .env 文件
	_ = godotenv.Load("/blog-backend-go-dev/.env")
	fmt.Printf("MongoDB 配置: URI=%s, 数据库=%s, 集合=%s\n", os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DATABASE"), os.Getenv("COLLECTION_NAME"))

	// 初始化 MongoDB（自动配置）
	client, err := config.InitMongo()
	if err != nil {
		log.Fatal("MongoDB 初始化失败:", err)
	}
	defer func() {
		if err := config.CloseMongo(); err != nil {
			log.Fatal("断开 MongoDB 连接失败:", err)
		}
	}()
	fmt.Println("成功连接到 MongoDB")

	// 初始化服务
	blogService := services.NewBlogService(client, config.MongoDatabase, getEnv("COLLECTION_NAME", "blogs-dev"))

	// 初始化认证服务
	jwtSecret := getEnv("JWT_SECRET", "default-jwt-secret") // 在生产环境中请通过环境变量注入
	authService := services.NewAuthService(client, config.MongoDatabase, "users", jwtSecret)

	// 初始化处理器
	blogHandler := handlers.NewBlogHandler(blogService)
	authHandler := handlers.NewAuthHandler(authService)

	// 初始化中间件
	jwtMiddleware := middleware.NewJWTMiddleware(authService)

	// 创建路由
	r := mux.NewRouter()

	// 注册路由（集中管理）
	routes.RegisterRoutes(r, blogHandler, authHandler, jwtMiddleware)

	// 健康检查端点
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "服务运行正常")
	}).Methods("GET")

	// 启动服务器
	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("博客服务器启动在 http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(addr, r))

}

// getEnv 从环境变量读取值，若为空则返回默认值
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
