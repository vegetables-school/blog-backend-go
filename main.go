package main

import (
	"fmt"
	"log"
	"net/http"

	"blog/handlers"
	"blog/services"

	"github.com/gorilla/mux"
)

func main() {
	// 初始化服务
	postService := services.NewPostService()
	postService.InitializeSampleData()

	// 初始化处理器
	postHandler := handlers.NewPostHandler(postService)

	// 创建路由
	r := mux.NewRouter()

	// 定义 API 端点
	r.HandleFunc("/api/posts", postHandler.GetPosts).Methods("GET")
	r.HandleFunc("/api/posts/{id}", postHandler.GetPost).Methods("GET")
	r.HandleFunc("/api/posts", postHandler.CreatePost).Methods("POST")
	r.HandleFunc("/api/posts/{id}", postHandler.UpdatePost).Methods("PUT")
	r.HandleFunc("/api/posts/{id}", postHandler.DeletePost).Methods("DELETE")

	// 健康检查端点
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "服务运行正常")
	}).Methods("GET")

	// 启动服务器
	port := ":8080"
	fmt.Printf("博客服务器启动在端口 %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
