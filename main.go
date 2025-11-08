package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"blog/handlers"
	"blog/middleware"
	"blog/routes"
	"blog/services"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB 连接字符串
	mongoURI := "mongodb://localhost:27017"
	// MongoDB 数据库名称
	dbName := "blogs-db-test"
	// MongoDB 集合名称
	collectionName := "blogs-db"

	// 连接到 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("连接 MongoDB 失败:", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("断开 MongoDB 连接失败:", err)
		}
	}()

	// 验证连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB 连接验证失败:", err)
	}
	fmt.Println("成功连接到 MongoDB")

	// 初始化服务
	blogService := services.NewBlogService(client, dbName, collectionName)

	// 初始化认证服务
	jwtSecret := "your-secret-key" // 在生产环境中应该从环境变量读取
	authService := services.NewAuthService(client, dbName, "users", jwtSecret)

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
	port := ":88"
	fmt.Printf("博客服务器启动在 http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
