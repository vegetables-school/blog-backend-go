package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"blog/handlers"
	"blog/services"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB 连接字符串
	mongoURI := "mongodb://localhost:27017"
	dbName := "blogs-db"
	collectionName := "blogs-test"

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
	postService := services.NewPostService(client, dbName, collectionName)

	// 初始化示例数据
	if err := postService.InitializeSampleData(); err != nil {
		log.Fatal("初始化示例数据失败:", err)
	}

	// 初始化处理器
	postHandler := handlers.NewPostHandler(postService)

	// 创建路由
	r := mux.NewRouter()

	// 定义 API 端点
	r.HandleFunc("/api/posts", postHandler.GetPosts).Methods("GET")
	r.HandleFunc("/api/posts/{id}", postHandler.GetPost).Methods("GET")
	r.HandleFunc("/api/createBlog", postHandler.CreateBlog).Methods("POST")
	r.HandleFunc("/api/posts/{id}", postHandler.UpdatePost).Methods("PUT")
	r.HandleFunc("/api/posts/{id}", postHandler.DeletePost).Methods("DELETE")

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
