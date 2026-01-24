package handlers

import (
	"blog/middleware"
	"blog/models"
	"blog/services"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateCommentHandler 新增评论
// @Summary 创建评论
// @Description 创建新的评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param data body models.CreateCommentRequest true "评论内容"
// @Success 201 {object} models.Comment
// @Failure 400 {string} string "参数错误"
// @Router /api/comment [post]
func (h *CommentHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效请求", http.StatusBadRequest)
		return
	}
	blogID, err := primitive.ObjectIDFromHex(req.BlogID)
	if err != nil {
		http.Error(w, "无效博客ID", http.StatusBadRequest)
		return
	}
	userIDHex := middleware.GetUserID(r)
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		http.Error(w, "无效用户ID", http.StatusUnauthorized)
		return
	}
	username := middleware.GetUsername(r)
	comment, err := h.commentService.CreateComment(blogID, userID, username, req.Content)
	if err != nil {
		http.Error(w, "评论失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// DeleteCommentHandler 删除评论（只能本人或管理员）
// @Summary 删除评论
// @Description 删除指定评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param data body models.DeleteCommentRequest true "评论ID"
// @Success 204 {string} string "删除成功"
// @Failure 404 {string} string "未找到"
// @Router /api/comment [delete]
func (h *CommentHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效请求", http.StatusBadRequest)
		return
	}
	commentIDHex := req.ID
	commentID, err := primitive.ObjectIDFromHex(commentIDHex)
	if err != nil {
		http.Error(w, "无效评论ID", http.StatusBadRequest)
		return
	}
	err = h.commentService.DeleteComment(commentID)
	if err != nil {
		http.Error(w, "删除失败", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetCommentUserID 用于权限中间件，获取评论的 user_id
func (h *CommentHandler) GetCommentUserID(r *http.Request) (primitive.ObjectID, error) {
	var req models.DeleteCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return primitive.NilObjectID, err
	}
	commentIDHex := req.ID
	commentID, err := primitive.ObjectIDFromHex(commentIDHex)
	if err != nil {
		return primitive.NilObjectID, err
	}
	comment, err := h.commentService.GetCommentByID(commentID)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return comment.UserID, nil
}
