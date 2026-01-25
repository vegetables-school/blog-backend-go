package utils

import "errors"

// 自定义错误类型
var (
	ErrBadRequest     = errors.New("BAD_REQUEST")
	ErrUnauthorized   = errors.New("UNAUTHORIZED")
	ErrForbidden      = errors.New("FORBIDDEN")
	ErrNotFound       = errors.New("NOT_FOUND")
	ErrConflict       = errors.New("CONFLICT")
	ErrInternalServer = errors.New("INTERNAL_SERVER_ERROR")
)

// AppError 应用错误类型
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建应用错误
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    ErrBadRequest.Error(),
		Message: message,
	}
}

// NewAuthError 创建认证错误
func NewAuthError(message string) *AppError {
	return &AppError{
		Code:    ErrUnauthorized.Error(),
		Message: message,
	}
}

// NewForbiddenError 创建权限错误
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    ErrForbidden.Error(),
		Message: message,
	}
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    ErrNotFound.Error(),
		Message: message,
	}
}

// NewConflictError 创建冲突错误
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    ErrConflict.Error(),
		Message: message,
	}
}

// NewInternalError 创建内部错误
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrInternalServer.Error(),
		Message: message,
		Err:     err,
	}
}
