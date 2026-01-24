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

// NewBlogHandler 创建新的 BlogHandler 实例（确保导出）
func NewBlogHandler(blogService *services.BlogService) *BlogHandler {
	return &BlogHandler{blogService: blogService}
}

// 标签相关请求体
type TagRequest struct {
	Tag string `json:"tag"`
}

// GetBlogsByTagsHandler 分页获取包含所有指定标签的文章
// @Summary 分页获取包含所有指定标签的文章
// @Description 根据标签数组分页获取所有包含这些标签的博客文章
// @Tags 博客
// @Param tag query string true "标签（可重复）"
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} BlogListResponse
// @Failure 500 {string} string "获取文章失败"
// @Router /api/blogs/by-tags [get]
func (h *BlogHandler) GetBlogsByTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query()["tag"] // 支持多个 tag=xxx
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page := int64(1)
	limit := int64(10)
	if pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	blogs, total, err := h.blogService.GetBlogsByTagsWithPagination(tags, page, limit)
	if err != nil {
		http.Error(w, "获取文章失败", http.StatusInternalServerError)
		return
	}
	resp := BlogListResponse{
		Data:       blogs,
		Pagination: Pagination{Page: page, Limit: limit, Total: total},
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}

// SearchBlogsHandler 支持模糊搜索和标签筛选
// @Summary 博客模糊搜索与标签筛选
// @Description 支持关键字和标签数组的分页搜索
// @Tags 博客
// @Param keyword query string false "关键字"
// @Param tag query string false "标签（可重复）"
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} BlogListResponse
// @Failure 500 {string} string "搜索失败"
// @Router /api/blogs/search [get]
func (h *BlogHandler) SearchBlogsHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	tagsParam := r.URL.Query()["tag"] // 支持多个 tag=xxx
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page := int64(1)
	limit := int64(10)
	if pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	blogs, total, err := h.blogService.SearchBlogs(keyword, tagsParam, page, limit)
	if err != nil {
		http.Error(w, "搜索失败", http.StatusInternalServerError)
		return
	}
	resp := BlogListResponse{
		Data:       blogs,
		Pagination: Pagination{Page: page, Limit: limit, Total: total},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetAllTagsHandler 获取所有标签
// @Summary 获取所有标签
// @Description 获取所有博客标签
// @Tags 博客
// @Success 200 {object} map[string][]string
// @Failure 500 {string} string "获取标签失败"
// @Router /api/tags [get]
func (h *BlogHandler) GetAllTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := h.blogService.GetAllTags()
	if err != nil {
		http.Error(w, "获取标签失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Data []string `json:"data"`
	}{Data: tags})
}

// AddTagHandler 新增标签（如有专门标签集合时用）
// @Summary 新增标签
// @Description 新增一个标签（如有专门标签集合时用）
// @Tags 博客
// @Accept json
// @Produce json
// @Param data body TagRequest true "标签内容"
// @Success 201 {string} string "创建成功"
// @Failure 400 {string} string "无效的标签"
// @Failure 500 {string} string "添加标签失败"
// @Router /api/admin/tag [post]
func (h *BlogHandler) AddTagHandler(w http.ResponseWriter, r *http.Request) {
	var req TagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Tag == "" {
		http.Error(w, "无效的标签", http.StatusBadRequest)
		return
	}
	err := h.blogService.AddTag(req.Tag)
	if err != nil {
		http.Error(w, "添加标签失败", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// DeleteTagHandler 删除标签（会从所有博客中移除该标签）
// @Summary 删除标签
// @Description 删除标签（会从所有博客中移除该标签）
// @Tags 博客
// @Accept json
// @Produce json
// @Param data body TagRequest true "标签内容"
// @Success 204 {string} string "删除成功"
// @Failure 400 {string} string "无效的标签"
// @Failure 500 {string} string "删除标签失败"
// @Router /api/admin/tag [delete]
func (h *BlogHandler) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	var req TagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Tag == "" {
		http.Error(w, "无效的标签", http.StatusBadRequest)
		return
	}
	err := h.blogService.DeleteTag(req.Tag)
	if err != nil {
		http.Error(w, "删除标签失败", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetBlogsPaginated 分页获取博客文章
// @Summary 分页获取博客文章
// @Description 分页获取博客文章
// @Tags 博客
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} BlogListResponse
// @Failure 500 {string} string "获取文章失败"
// @Router /api/blogs [get]
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
// @Summary 获取单篇博客文章
// @Description 根据 ID 获取博客详情
// @Tags 博客
// @Param id path string true "博客ID"
// @Success 200 {object} BlogResponse
// @Failure 404 {string} string "文章未找到"
// @Router /api/blog/{id} [get]
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
// @Summary 创建新博客文章
// @Description 创建新博客文章
// @Tags 博客
// @Accept json
// @Produce json
// @Param data body models.CreateBlogRequest true "博客内容"
// @Success 201 {object} BlogResponse
// @Failure 400 {string} string "无效的请求数据"
// @Failure 401 {string} string "未认证用户"
// @Failure 500 {string} string "创建文章失败"
// @Router /api/admin/blog [post]
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
// @Summary 更新博客文章
// @Description 更新博客文章
// @Tags 博客
// @Accept json
// @Produce json
// @Param id path string true "博客ID"
// @Param data body models.UpdateBlogRequest true "更新内容"
// @Success 200 {object} BlogResponse
// @Failure 400 {string} string "无效的请求数据"
// @Failure 401 {string} string "未认证用户"
// @Failure 404 {string} string "文章未找到"
// @Router /api/admin/blog/{id} [put]
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
// @Summary 删除博客文章
// @Description 删除博客文章
// @Tags 博客
// @Param id path string true "博客ID"
// @Success 204 {string} string "删除成功"
// @Failure 401 {string} string "未认证用户"
// @Failure 404 {string} string "文章未找到"
// @Router /api/admin/blog/{id} [delete]
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
