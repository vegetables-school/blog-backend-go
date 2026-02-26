package handlers

import (
	"blog/middleware"
	"blog/models"
	"blog/services"
	"blog/utils"
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
		utils.SendError(w, http.StatusBadRequest, "无效请求")
		return
	}
	blogID, err := primitive.ObjectIDFromHex(req.BlogID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "无效博客ID")
		return
	}
	userIDHex := middleware.GetUserID(r)
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "无效用户ID")
		return
	}
	username := middleware.GetUsername(r)
	comment, err := h.commentService.CreateComment(blogID, userID, username, req.Content)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "评论失败")
		return
	}
	utils.SendJSON(w, http.StatusCreated, comment)
}

// DeleteCommentHandler 删除评论（只能本人或管理员）
// @Summary 删除评论
// @Description 删除指定评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param data body models.DeleteCommentRequest true "评论ID"
// @Success 204 {string} string "删除成功"
// @Failure 400 {string} string "无效请求"
// @Failure 401 {string} string "未认证"
// @Failure 403 {string} string "无权限"
// @Failure 404 {string} string "未找到"
// @Router /api/comment [delete]
func (h *CommentHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "无效请求")
		return
	}
	commentIDHex := req.ID
	commentID, err := primitive.ObjectIDFromHex(commentIDHex)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "无效评论ID")
		return
	}

	// 获取当前用户ID和角色
	currentUserID := middleware.GetUserID(r)
	currentUserRole := middleware.GetUserRole(r)
	if currentUserID == "" {
		utils.SendError(w, http.StatusUnauthorized, "未认证")
		return
	}

	// 获取评论信息
	comment, err := h.commentService.GetCommentByID(commentID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "评论未找到")
		return
	}

	// 检查权限：只有评论作者或管理员可以删除
	commentUserID := comment.UserID.Hex()
	if commentUserID != currentUserID && currentUserRole != "admin" {
		utils.SendError(w, http.StatusForbidden, "无权限删除他人评论")
		return
	}

	err = h.commentService.DeleteComment(commentID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "删除失败")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetCommentUserID 用于权限中间件，获取评论的 user_id
// 注意：此方法已废弃，权限验证已集成到 DeleteCommentHandler 中
func (h *CommentHandler) GetCommentUserID(r *http.Request) (primitive.ObjectID, error) {
	return primitive.NilObjectID, nil
}
