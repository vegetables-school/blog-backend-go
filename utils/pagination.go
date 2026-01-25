package utils

import (
	"net/http"
	"strconv"
)

// ParsePagination 解析分页参数
func ParsePagination(r *http.Request) (page, limit int64) {
	page = 1
	limit = 10

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil && p > 0 {
			page = p
		}
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return
}

// CalculateSkip 计算跳过的文档数
func CalculateSkip(page, limit int64) int64 {
	return (page - 1) * limit
}
