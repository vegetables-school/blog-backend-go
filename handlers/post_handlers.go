package handlers

import (
	"encoding/json"
	"net/http"

	"blog/middleware"
	"blog/services"

	"github.com/gorilla/mux"
)

// BlogHandler 处理博客文章的HTTP请求
type BlogHandler struct {
	blogService *services.BlogService
}

// NewBlogHandler 创建新的BlogHandler实例
func NewBlogHandler(blogService *services.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
	}
}

// GetBlogs 获取所有博客文章
func (h *BlogHandler) GetBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.blogService.GetAllBlogs()
	if err != nil {
		http.Error(w, "获取文章失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

// GetBlog 获取单篇博客文章
func (h *BlogHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	blog, err := h.blogService.GetBlogByID(id)
	if err != nil {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

// CreateBlog 创建新博客文章
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 从认证上下文中获取作者信息
	author := middleware.GetUsername(r)
	if author == "" {
		http.Error(w, "未认证用户", http.StatusUnauthorized)
		return
	}

	blog, err := h.blogService.CreateBlog(req.Title, req.Content, author)
	if err != nil {
		http.Error(w, "创建文章失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

// UpdateBlog 更新博客文章
func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 从认证上下文中获取用户信息
	username := middleware.GetUsername(r)
	if username == "" {
		http.Error(w, "未认证用户", http.StatusUnauthorized)
		return
	}

	// 检查是否是文章作者（这里简化了，实际应该从数据库检查）
	// TODO: 添加权限检查

	blog, err := h.blogService.UpdateBlog(id, req.Title, req.Content, username)
	if err != nil {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

// DeleteBlog 删除博客文章
func (h *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// 从认证上下文中获取用户信息
	username := middleware.GetUsername(r)
	if username == "" {
		http.Error(w, "未认证用户", http.StatusUnauthorized)
		return
	}

	// 检查是否是文章作者（这里简化了，实际应该从数据库检查）
	// TODO: 添加权限检查

	if err := h.blogService.DeleteBlog(id); err != nil {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
