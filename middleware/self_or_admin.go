package middleware

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OnlySelfOrAdmin 只允许本人或管理员操作
func OnlySelfOrAdmin(getResourceUserID func(r *http.Request) (primitive.ObjectID, error)) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			role := GetUserRole(r)
			userID := GetUserID(r)
			if role == "admin" {
				next(w, r)
				return
			}
			resourceUserID, err := getResourceUserID(r)
			if err != nil || resourceUserID.Hex() != userID {
				http.Error(w, "无权限操作他人数据", http.StatusForbidden)
				return
			}
			next(w, r)
		}
	}
}
