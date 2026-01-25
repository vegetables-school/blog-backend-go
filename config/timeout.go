package config

import "time"

// 数据库操作超时配置
const (
	DefaultTimeout    = 10 * time.Second // 默认查询超时
	QuickQueryTimeout = 5 * time.Second  // 快速查询超时
	SlowQueryTimeout  = 30 * time.Second // 慢查询超时
	WriteTimeout      = 15 * time.Second // 写操作超时
)

// GetTimeout 根据操作类型获取超时时间
func GetTimeout(isSlow bool, isWrite bool) time.Duration {
	switch {
	case isSlow:
		return SlowQueryTimeout
	case isWrite:
		return WriteTimeout
	default:
		return DefaultTimeout
	}
}

// GetSlowQueryTimeout 获取慢查询超时时间
func GetSlowQueryTimeout() time.Duration {
	return SlowQueryTimeout
}
