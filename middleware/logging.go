package middleware

import (
	"net/http"
	"time"

	"blog/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// responseWriterWrapper 包装 ResponseWriter 以捕获状态码
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// RequestLogger 记录 HTTP 请求
func RequestLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 生成请求 ID
		requestID := uuid.New().String()

		// 包装 ResponseWriter 以捕获状态码
		wrappedWriter := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(wrappedWriter, r)

		duration := time.Since(start)

		utils.Info("HTTP Request",
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.Int("status", wrappedWriter.statusCode),
			zap.Duration("duration", duration),
			zap.String("user_agent", r.UserAgent()),
			zap.String("remote_addr", r.RemoteAddr),
		)
	}
}
