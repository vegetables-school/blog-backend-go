package handlers

import (
	"blog/middleware"
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
func (h *LikeHandler) AddLikeHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BlogID string `json:"blog_id"`
	}
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
func (h *LikeHandler) RemoveLikeHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BlogID string `json:"blog_id"`
	}
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
	var req struct {
		BlogID string `json:"blog_id"`
	}
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
	return userID, nil
}
