package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse 统一成功响应格式
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SendError 发送错误响应（支持多种调用方式）
// 方式1: SendError(w, err) - err 是 AppError
// 方式2: SendError(w, err, statusCode) - err 是普通 error
// 方式3: SendError(w, statusCode, message) - 向后兼容
func SendError(w http.ResponseWriter, args ...interface{}) {
	var statusCode int
	var message string
	var err error

	if len(args) == 1 {
		// 只有一个参数，可能是 AppError
		if appErr, ok := args[0].(*AppError); ok {
			err = appErr
			message = appErr.Message
			switch appErr.Code {
			case ErrBadRequest.Error(), "BAD_REQUEST":
				statusCode = http.StatusBadRequest
			case ErrUnauthorized.Error(), "UNAUTHORIZED":
				statusCode = http.StatusUnauthorized
			case ErrForbidden.Error(), "FORBIDDEN":
				statusCode = http.StatusForbidden
			case ErrNotFound.Error(), "NOT_FOUND":
				statusCode = http.StatusNotFound
			case ErrConflict.Error(), "CONFLICT":
				statusCode = http.StatusConflict
			default:
				statusCode = http.StatusInternalServerError
				log.Printf("Internal Error: %v", appErr)
			}
		} else if e, ok := args[0].(error); ok {
			// 普通错误
			err = e
			message = e.Error()
			statusCode = http.StatusInternalServerError
		} else if s, ok := args[0].(string); ok {
			// 纯字符串
			message = s
			statusCode = http.StatusInternalServerError
		}
	} else if len(args) == 2 {
		// 两个参数
		if sc, ok := args[0].(int); ok {
			statusCode = sc
		}
		if s, ok := args[1].(string); ok {
			message = s
		}
	}

	response := ErrorResponse{
		Error: message,
	}

	if err != nil {
		if appErr, ok := err.(*AppError); ok {
			response.Code = appErr.Code
			response.Details = appErr.Details
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// SendJSON 发送 JSON 响应
func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{Data: data})
}

// SendJSONMessage 发带消息的 JSON 响应
func SendJSONMessage(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{Data: data, Message: message})
}
