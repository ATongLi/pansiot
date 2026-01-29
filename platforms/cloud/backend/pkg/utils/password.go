package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordArgon2Config Argon2配置参数
type PasswordArgon2Config struct {
	Time    uint32 // 迭代次数
	Memory  uint32 // 内存使用（KB）
	Threads uint8  // 线程数
	KeyLen  uint32 // 密钥长度
	SaltLen uint32 // 盐长度
}

// DefaultPasswordConfig 默认密码配置
var DefaultPasswordConfig = &PasswordArgon2Config{
	Time:    3,      // 3次迭代
	Memory:  64 * 1024, // 64MB内存
	Threads: 4,      // 4个线程
	KeyLen:  32,     // 32字节密钥
	SaltLen: 16,     // 16字节盐
}

// PasswordStrength 密码强度等级
type PasswordStrength int

const (
	StrengthWeak   PasswordStrength = iota // 弱密码
	StrengthMedium                          // 中等密码
	StrengthStrong                          // 强密码
	StrengthVeryStrong                      // 非常强密码
)

// PasswordValidator 密码验证器
type PasswordValidator struct {
	minLength       int
	requireUpper    bool
	requireLower    bool
	requireNumber   bool
	requireSpecial  bool
	specialChars    string
}

// NewPasswordValidator 创建密码验证器
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		minLength:      8,
		requireUpper:   true,
		requireLower:   true,
		requireNumber:  true,
		requireSpecial: true,
		specialChars:   "!@#$%^&*()_+-=[]{}|;:,.<>?",
	}
}

// ValidatePassword 验证密码强度
func (pv *PasswordValidator) ValidatePassword(password string) error {
	if len(password) < pv.minLength {
		return fmt.Errorf("密码长度至少需要%d位", pv.minLength)
	}

	if pv.requireUpper {
		hasUpper := false
		for _, c := range password {
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return fmt.Errorf("密码必须包含至少一个大写字母")
		}
	}

	if pv.requireLower {
		hasLower := false
		for _, c := range password {
			if c >= 'a' && c <= 'z' {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return fmt.Errorf("密码必须包含至少一个小写字母")
		}
	}

	if pv.requireNumber {
		hasNumber := false
		for _, c := range password {
			if c >= '0' && c <= '9' {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			return fmt.Errorf("密码必须包含至少一个数字")
		}
	}

	if pv.requireSpecial {
		hasSpecial := false
		for _, c := range password {
			if strings.ContainsRune(pv.specialChars, c) {
				hasSpecial = true
				break
			}
		}
		if !hasSpecial {
			return fmt.Errorf("密码必须包含至少一个特殊字符: %s", pv.specialChars)
		}
	}

	return nil
}

// GetPasswordStrength 获取密码强度等级
func GetPasswordStrength(password string) PasswordStrength {
	score := 0

	// 长度评分
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}

	// 字符类型评分
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", c):
			hasSpecial = true
		}
	}

	typeCount := 0
	if hasUpper {
		typeCount++
	}
	if hasLower {
		typeCount++
	}
	if hasNumber {
		typeCount++
	}
	if hasSpecial {
		typeCount++
	}

	score += typeCount

	// 计算强度
	if score >= 6 {
		return StrengthVeryStrong
	} else if score >= 5 {
		return StrengthStrong
	} else if score >= 3 {
		return StrengthMedium
	}
	return StrengthWeak
}

// HashPassword 使用Argon2id哈希密码
func HashPassword(password string) (string, error) {
	return HashPasswordWithConfig(password, DefaultPasswordConfig)
}

// HashPasswordWithConfig 使用指定配置哈希密码
func HashPasswordWithConfig(password string, config *PasswordArgon2Config) (string, error) {
	if len(password) == 0 {
		return "", fmt.Errorf("密码不能为空")
	}

	// 生成随机盐
	salt := make([]byte, config.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("生成盐失败: %w", err)
	}

	// 使用Argon2id哈希
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		config.Time,
		config.Memory,
		config.Threads,
		config.KeyLen,
	)

	// 格式: $argon2id$v=19$m=65536,t=3,p=4$<base64 salt>$<base64 hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		config.Memory,
		config.Time,
		config.Threads,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, encodedHash string) (bool, error) {
	// 解析哈希字符串
	_, _, _, salt, hash, err := parseArgon2Hash(encodedHash)
	if err != nil {
		return false, err
	}

	// 使用相同的盐哈希输入密码
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		DefaultPasswordConfig.Time,
		DefaultPasswordConfig.Memory,
		DefaultPasswordConfig.Threads,
		DefaultPasswordConfig.KeyLen,
	)

	// 常量时间比较，防止时序攻击
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// parseArgon2Hash 解析Argon2id哈希字符串
func parseArgon2Hash(encodedHash string) (time, memory uint32, threads uint8, salt, hash []byte, err error) {
	// 格式: $argon2id$v=19$m=65536,t=3,p=4$<base64 salt>$<base64 hash>
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		err = fmt.Errorf("无效的哈希格式")
		return
	}

	// 解析参数
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return
	}

	// 解码盐和哈希
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return
	}

	return
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}

	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789!@#$%^&*"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

// ShouldRehash 检查是否需要重新哈希（当配置参数更新时）
func ShouldRehash(encodedHash string) bool {
	_, _, _, salt, _, err := parseArgon2Hash(encodedHash)
	if err != nil {
		return true
	}

	// 如果盐长度不同，说明配置已更新
	return uint32(len(salt)) != DefaultPasswordConfig.SaltLen
}
