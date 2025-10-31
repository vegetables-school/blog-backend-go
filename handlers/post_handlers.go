package handlers

import (
	"encoding/json"
	"net/http"

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
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		http.Error(w, "获取文章失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetPost 获取单篇博客文章
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// CreateBlog 创建新博客文章
func (h *PostHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	post, err := h.postService.CreateBlog(req.Title, req.Content, req.Author)
	if err != nil {
		http.Error(w, "创建文章失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// UpdatePost 更新博客文章
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	post, err := h.postService.UpdatePost(id, req.Title, req.Content, req.Author)
	if err != nil {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// DeletePost 删除博客文章
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.postService.DeletePost(id); err != nil {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
