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
	if err := h.likeService.AddLike(blogID, userID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "点赞失败")
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
// @Failure 400 {string} string "无效请求"
// @Failure 401 {string} string "未认证"
// @Failure 404 {string} string "点赞记录未找到"
// @Router /api/like [delete]
func (h *LikeHandler) RemoveLikeHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RemoveLikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "无效请求")
		return
	}

	// 验证博客ID
	blogID, err := primitive.ObjectIDFromHex(req.BlogID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "无效博客ID")
		return
	}

	// 从JWT token获取当前用户ID（不使用请求体中的user_id）
	currentUserIDHex := middleware.GetUserID(r)
	if currentUserIDHex == "" {
		utils.SendError(w, http.StatusUnauthorized, "未认证")
		return
	}

	currentUserID, err := primitive.ObjectIDFromHex(currentUserIDHex)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "无效用户ID")
		return
	}

	// 验证点赞记录是否存在
	_, err = h.likeService.GetLikeByBlogAndUser(blogID, currentUserID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "点赞记录未找到")
		return
	}

	if err := h.likeService.RemoveLike(blogID, currentUserID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "取消点赞失败")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetLikeUserID 用于权限中间件，获取点赞的 user_id
// 注意：此方法已废弃，权限验证已集成到 RemoveLikeHandler 中
func (h *LikeHandler) GetLikeUserID(r *http.Request) (primitive.ObjectID, error) {
	return primitive.NilObjectID, nil
}
