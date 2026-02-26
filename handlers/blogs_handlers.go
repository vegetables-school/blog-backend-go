package handlers

import (
	"encoding/json"
	"net/http"

	"blog/middleware"
	"blog/services"
	"blog/utils"

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
	tags := r.URL.Query()["tag"]
	page, limit := utils.ParsePagination(r)

	blogs, total, err := h.blogService.GetBlogsByTagsWithPagination(tags, page, limit)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "获取文章失败")
		return
	}

	resp := BlogListResponse{
		Data:       blogs,
		Pagination: Pagination{Page: page, Limit: limit, Total: total},
	}
	utils.SendJSON(w, http.StatusOK, resp)
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
	tagsParam := r.URL.Query()["tag"]
	page, limit := utils.ParsePagination(r)

	blogs, total, err := h.blogService.SearchBlogs(keyword, tagsParam, page, limit)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "搜索失败")
		return
	}

	resp := BlogListResponse{
		Data:       blogs,
		Pagination: Pagination{Page: page, Limit: limit, Total: total},
	}
	utils.SendJSON(w, http.StatusOK, resp)
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
		utils.SendError(w, http.StatusInternalServerError, "获取标签失败")
		return
	}
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{"data": tags})
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
// @Failure 501 {string} string "功能未实现"
// @Router /api/admin/tag [post]
func (h *BlogHandler) AddTagHandler(w http.ResponseWriter, r *http.Request) {
	var req TagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Tag == "" {
		utils.SendError(w, http.StatusBadRequest, "无效的标签")
		return
	}
	// 标签管理已集成在博客创建/编辑中，此端点暂不实现单独的标签集合管理
	utils.SendError(w, http.StatusNotImplemented, "标签请直接在创建或编辑博客时管理")
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
		utils.SendError(w, http.StatusBadRequest, "无效的标签")
		return
	}
	err := h.blogService.DeleteTag(req.Tag)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "删除标签失败")
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
	page, limit := utils.ParsePagination(r)

	blogs, total, err := h.blogService.GetBlogsWithPagination(page, limit)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "获取文章失败")
		return
	}

	response := BlogListResponse{
		Data:       blogs,
		Pagination: Pagination{Page: page, Limit: limit, Total: total},
	}

	utils.SendJSON(w, http.StatusOK, response)
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
		utils.SendError(w, http.StatusNotFound, "文章未找到")
		return
	}

	utils.SendJSON(w, http.StatusOK, BlogResponse{Data: blog})
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
		utils.SendError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	// 从认证上下文中获取作者信息
	author := middleware.GetUsername(r)
	if author == "" {
		utils.SendError(w, http.StatusUnauthorized, "未认证用户")
		return
	}

	showVal := true
	if req.Show != nil {
		showVal = *req.Show
	}

	blog, err := h.blogService.CreateBlog(req.Title, req.Content, author, req.Tags, showVal)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "创建文章失败")
		return
	}

	utils.SendJSON(w, http.StatusCreated, BlogResponse{Data: blog})
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
// @Failure 403 {string} string "无权限操作他人博客"
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
		utils.SendError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	// 从认证上下文中获取用户信息（用于权限校验）
	username := middleware.GetUsername(r)
	role := middleware.GetUserRole(r)
	if username == "" {
		utils.SendError(w, http.StatusUnauthorized, "未认证用户")
		return
	}

	// 检查权限：只有管理员可以修改作者字段，普通用户只能修改自己的博客
	if req.Author != nil && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "无权限修改作者信息")
		return
	}

	// 检查是否是文章作者或管理员
	existingBlog, err := h.blogService.GetBlogByID(id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "文章未找到")
		return
	}

	if existingBlog.Author != username && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "无权限操作他人博客")
		return
	}

	blog, err := h.blogService.UpdateBlog(id, req.Title, req.Content, req.Author, req.Tags, req.Show, req.Views)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "文章未找到")
		return
	}

	utils.SendJSON(w, http.StatusOK, BlogResponse{Data: blog})
}

// DeleteBlog 删除博客文章
// @Summary 删除博客文章
// @Description 删除博客文章
// @Tags 博客
// @Param id path string true "博客ID"
// @Success 204 {string} string "删除成功"
// @Failure 401 {string} string "未认证用户"
// @Failure 403 {string} string "无权限操作他人博客"
// @Failure 404 {string} string "文章未找到"
// @Router /api/admin/blog/{id} [delete]
func (h *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// 从认证上下文中获取用户信息
	username := middleware.GetUsername(r)
	role := middleware.GetUserRole(r)
	if username == "" {
		utils.SendError(w, http.StatusUnauthorized, "未认证用户")
		return
	}

	// 检查是否是文章作者或管理员
	existingBlog, err := h.blogService.GetBlogByID(id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "文章未找到")
		return
	}

	if existingBlog.Author != username && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "无权限操作他人博客")
		return
	}

	if err := h.blogService.DeleteBlog(id); err != nil {
		utils.SendError(w, http.StatusNotFound, "文章未找到")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
