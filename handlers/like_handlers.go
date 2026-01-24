package handlers

import (
	"blog/middleware"
	"blog/models"
	"blog/services"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeHandler struct {
	likeService *services.LikeService
}

func NewLikeHandler(likeService *services.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// AddLikeHandler 新增点赞
// @Summary 点赞
// @Description 给博客点赞
// @Tags 点赞
// @Accept json
// @Produce json
// @Param data body models.AddLikeRequest true "点赞内容"
// @Success 201 {object} models.Like
// @Failure 400 {string} string "参数错误"
// @Router /api/like [post]
func (h *LikeHandler) AddLikeHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AddLikeRequest
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
	if err := h.likeService.AddLike(blogID, userID); err != nil {
		http.Error(w, "点赞失败", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// RemoveLikeHandler 取消点赞（只能本人或管理员）
// @Summary 取消点赞
// @Description 取消对博客的点赞
// @Tags 点赞
// @Accept json
// @Produce json
// @Param data body models.RemoveLikeRequest true "点赞内容"
// @Success 204 {string} string "取消成功"
// @Failure 404 {string} string "未找到"
// @Router /api/like [delete]
func (h *LikeHandler) RemoveLikeHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RemoveLikeRequest
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
	if err := h.likeService.RemoveLike(blogID, userID); err != nil {
		http.Error(w, "取消点赞失败", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetLikeUserID 用于权限中间件，获取点赞的 user_id
func (h *LikeHandler) GetLikeUserID(r *http.Request) (primitive.ObjectID, error) {
	var req models.RemoveLikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return primitive.NilObjectID, err
	}
	blogID, err := primitive.ObjectIDFromHex(req.BlogID)
	if err != nil {
		return primitive.NilObjectID, err
	}
	userIDHex := middleware.GetUserID(r)
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return primitive.NilObjectID, err
	}
	// 查询点赞记录，确保该用户对该博客有点赞
	like, err := h.likeService.GetLikeByBlogAndUser(blogID, userID)
	if err != nil {
		return primitive.NilObjectID, err // 没有点赞记录
	}
	return like.UserID, nil
}
