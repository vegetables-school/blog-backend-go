package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"blog/services"
	"github.com/gorilla/mux"
)

// PostHandler 处理博客文章的HTTP请求
type PostHandler struct {
	postService *services.PostService
}

// NewPostHandler 创建新的PostHandler实例
func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// GetPosts 获取所有博客文章
func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts := h.postService.GetAllPosts()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetPost 获取单篇博客文章
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "无效的 ID", http.StatusBadRequest)
		return
	}

	post, exists := h.postService.GetPostByID(id)
	if !exists {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// CreatePost 创建新博客文章
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	post := h.postService.CreatePost(req.Title, req.Content, req.Author)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// UpdatePost 更新博客文章
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "无效的 ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	post, exists := h.postService.UpdatePost(id, req.Title, req.Content, req.Author)
	if !exists {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// DeletePost 删除博客文章
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "无效的 ID", http.StatusBadRequest)
		return
	}

	if !h.postService.DeletePost(id) {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
