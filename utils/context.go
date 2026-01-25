package utils

import (
	"context"
	"time"

	"blog/config"
)

// NewContext 创建带超时的 context
func NewContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// NewQuickContext 创建快速查询 context
func NewQuickContext() (context.Context, context.CancelFunc) {
	return NewContext(config.QuickQueryTimeout)
}

// NewDefaultContext 创建默认查询 context
func NewDefaultContext() (context.Context, context.CancelFunc) {
	return NewContext(config.DefaultTimeout)
}

// NewWriteContext 创建写操作 context
func NewWriteContext() (context.Context, context.CancelFunc) {
	return NewContext(config.WriteTimeout)
}

// NewSlowQueryContext 创建慢查询 context
func NewSlowQueryContext() (context.Context, context.CancelFunc) {
	return NewContext(config.GetSlowQueryTimeout())
}
