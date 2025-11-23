package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	// 返回强类型响应，主数据放在 data 下
	resp := BlogListResponse{
		Data:       blogs,
		Pagination: Pagination{Page: 1, Limit: int64(len(blogs)), Total: int64(len(blogs))},
	}
	json.NewEncoder(w).Encode(resp)
}

// GetBlogsPaginated 分页获取博客文章
func (h *BlogHandler) GetBlogsPaginated(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// 默认值
	page := int64(1)
	limit := int64(10)

	if pageStr != "" {
		if parsedPage, err := strconv.ParseInt(pageStr, 10, 64); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr != "" {
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 64); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	blogs, total, err := h.blogService.GetBlogsWithPagination(page, limit)
	if err != nil {
		http.Error(w, "获取文章失败", http.StatusInternalServerError)
		return
	}

	response := BlogListResponse{
		Data:       blogs,
		Pagination: Pagination{Page: page, Limit: limit, Total: total},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
	json.NewEncoder(w).Encode(BlogResponse{Data: blog})
}

// CreateBlog 创建新博客文章
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags,omitempty"`
		Show    *bool    `json:"show,omitempty"`
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

	showVal := true
	if req.Show != nil {
		showVal = *req.Show
	}

	blog, err := h.blogService.CreateBlog(req.Title, req.Content, author, req.Tags, showVal)
	if err != nil {
		http.Error(w, "创建文章失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(BlogResponse{Data: blog})
}

// UpdateBlog 更新博客文章
func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Title   *string  `json:"title,omitempty"`
		Content *string  `json:"content,omitempty"`
		Author  *string  `json:"author,omitempty"`
		Tags    []string `json:"tags,omitempty"`
		Show    *bool    `json:"show,omitempty"`
		Views   *int64   `json:"views,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 从认证上下文中获取用户信息（用于权限校验）
	username := middleware.GetUsername(r)
	if username == "" {
		http.Error(w, "未认证用户", http.StatusUnauthorized)
		return
	}

	// 检查是否是文章作者（这里简化了，实际应该从数据库检查）
	// TODO: 添加权限检查

	blog, err := h.blogService.UpdateBlog(id, req.Title, req.Content, req.Author, req.Tags, req.Show, req.Views)
	if err != nil {
		http.Error(w, "文章未找到", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(BlogResponse{Data: blog})
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
