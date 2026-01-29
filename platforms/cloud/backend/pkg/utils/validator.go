package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Validator 验证器
type Validator struct {
	emailRegex *regexp.Regexp
	phoneRegex *regexp.Regexp
}

// NewValidator 创建验证器
func NewValidator() *Validator {
	return &Validator{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
		phoneRegex: regexp.MustCompile(`^\+?[1-9]\d{1,14}$`),
	}
}

// ValidateEmail 验证邮箱格式
func (v *Validator) ValidateEmail(email string) error {
	if len(email) == 0 {
		return fmt.Errorf("邮箱不能为空")
	}

	if !v.emailRegex.MatchString(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	return nil
}

// ValidatePhone 验证手机号（支持国际区号）
func (v *Validator) ValidatePhone(phone string) error {
	if len(phone) == 0 {
		return fmt.Errorf("手机号不能为空")
	}

	if !v.phoneRegex.MatchString(phone) {
		return fmt.Errorf("手机号格式不正确（支持国际区号）")
	}

	return nil
}

// ValidateUsername 验证用户名
func (v *Validator) ValidateUsername(username string) error {
	if len(username) == 0 {
		return fmt.Errorf("用户名不能为空")
	}

	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("用户名长度必须在3-32位之间")
	}

	// 只允许字母、数字、下划线
	for _, c := range username {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' {
			return fmt.Errorf("用户名只能包含字母、数字和下划线")
		}
	}

	return nil
}

// ValidateTenantName 验证组织名称
func (v *Validator) ValidateTenantName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("组织名称不能为空")
	}

	if len(name) < 2 || len(name) > 100 {
		return fmt.Errorf("组织名称长度必须在2-100位之间")
	}

	return nil
}

// ValidateIndustry 验证所属行业
func (v *Validator) ValidateIndustry(industry string) error {
	if len(industry) == 0 {
		return fmt.Errorf("所属行业不能为空")
	}

	validIndustries := []string{
		"制造业",
		"能源",
		"交通",
		"建筑",
		"农业",
		"医疗",
		"教育",
		"金融",
		"零售",
		"物流",
		"其他",
	}

	for _, valid := range validIndustries {
		if industry == valid {
			return nil
		}
	}

	return fmt.Errorf("所属行业无效，必须是以下之一: %s", strings.Join(validIndustries, "、"))
}

// ValidateSerialNumber 验证企业序列号
func (v *Validator) ValidateSerialNumber(serialNumber string) error {
	if len(serialNumber) == 0 {
		return fmt.Errorf("企业序列号不能为空")
	}

	if !ValidateSerialNumber(serialNumber) {
		return fmt.Errorf("企业序列号格式不正确（应为8位: 4位字符+4位数字）")
	}

	return nil
}

// ValidateVerificationCode 验证验证码
func (v *Validator) ValidateVerificationCode(code string) error {
	if len(code) == 0 {
		return fmt.Errorf("验证码不能为空")
	}

	// 验证码通常为4-6位数字
	if len(code) < 4 || len(code) > 6 {
		return fmt.Errorf("验证码长度必须在4-6位之间")
	}

	for _, c := range code {
		if !unicode.IsDigit(c) {
			return fmt.Errorf("验证码只能包含数字")
		}
	}

	return nil
}

// ValidateRoleName 验证角色名称
func (v *Validator) ValidateRoleName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("角色名称不能为空")
	}

	if len(name) < 2 || len(name) > 50 {
		return fmt.Errorf("角色名称长度必须在2-50位之间")
	}

	return nil
}

// ValidatePermissionCode 验证权限代码
func (v *Validator) ValidatePermissionCode(code string) error {
	if len(code) == 0 {
		return fmt.Errorf("权限代码不能为空")
	}

	// 权限代码格式: MODULE_ACTION (如: DEVICE_CREATE)
	parts := strings.Split(code, "_")
	if len(parts) < 2 {
		return fmt.Errorf("权限代码格式不正确（应为: MODULE_ACTION）")
	}

	return nil
}

// ValidateQuotaValue 验证配额值
func (v *Validator) ValidateQuotaValue(value int) error {
	if value < 0 {
		return fmt.Errorf("配额值不能为负数")
	}

	return nil
}

// SanitizeString 清理字符串（去除首尾空格和特殊字符）
func (v *Validator) SanitizeString(input string) string {
	// 去除首尾空格
	input = strings.TrimSpace(input)

	// 去除控制字符
	var result strings.Builder
	for _, r := range input {
		if !unicode.IsControl(r) || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ValidateRequired 验证必填字段
func (v *Validator) ValidateRequired(fieldName, value string) error {
	if len(value) == 0 {
		return fmt.Errorf("%s不能为空", fieldName)
	}
	return nil
}

// ValidateLength 验证字段长度
func (v *Validator) ValidateLength(fieldName, value string, minLen, maxLen int) error {
	length := len(value)
	if length < minLen || length > maxLen {
		return fmt.Errorf("%s长度必须在%d-%d位之间", fieldName, minLen, maxLen)
	}
	return nil
}

// ValidateEnum 验证枚举值
func (v *Validator) ValidateEnum(fieldName, value string, validValues []string) error {
	for _, valid := range validValues {
		if value == valid {
			return nil
		}
	}
	return fmt.Errorf("%s无效，必须是以下之一: %s", fieldName, strings.Join(validValues, "、"))
}
