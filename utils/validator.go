package utils

import (
	"regexp"
	"unicode"
)

// ValidatePassword 验证密码强度
// 要求：至少8位，包含大小写字母、数字
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return NewValidationError("密码长度至少8位")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	requiredCount := 0
	if hasUpper {
		requiredCount++
	}
	if hasLower {
		requiredCount++
	}
	if hasNumber {
		requiredCount++
	}
	if hasSpecial {
		requiredCount++
	}

	if requiredCount < 3 {
		return NewValidationError("密码必须包含大写字母、小写字母、数字中的至少三种，建议包含特殊字符")
	}

	return nil
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return NewValidationError("邮箱格式不正确")
	}
	return nil
}

// SanitizeSearchKeyword 清理搜索关键字，防止正则注入
func SanitizeSearchKeyword(keyword string) string {
	// 转义正则特殊字符
	return regexp.QuoteMeta(keyword)
}

