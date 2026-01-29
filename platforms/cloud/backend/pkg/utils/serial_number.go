package utils

import (
	"crypto/rand"
	"math/big"
	"fmt"
	"sync"
)

// SerialNumberGenerator 企业序列号生成器
// 格式: {4位随机字符}{4位自增ID} (总长度8位)
// 示例: A3F20001
type SerialNumberGenerator struct {
	mu       sync.Mutex
	lastID   int32
	prefix   string
	counter  int32
}

// NewSerialNumberGenerator 创建序列号生成器
func NewSerialNumberGenerator() *SerialNumberGenerator {
	return &SerialNumberGenerator{
		lastID:  0,
		counter: 1,
	}
}

// Generate 生成企业序列号
// 格式: {4位随机字符}{4位自增ID}
func (g *SerialNumberGenerator) Generate() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 生成4位随机字符（大小写字母+数字）
	randomChars, err := g.generateRandomChars(4)
	if err != nil {
		return "", fmt.Errorf("failed to generate random chars: %w", err)
	}

	// 自增ID（4位，从0001开始）
	id := g.counter
	g.counter++

	// 格式化为8位序列号: XXXX####
	serialNumber := fmt.Sprintf("%s%04d", randomChars, id)

	return serialNumber, nil
}

// GenerateWithID 使用指定ID生成序列号（用于数据库导入）
func (g *SerialNumberGenerator) GenerateWithID(id int32) (string, error) {
	randomChars, err := g.generateRandomChars(4)
	if err != nil {
		return "", fmt.Errorf("failed to generate random chars: %w", err)
	}

	serialNumber := fmt.Sprintf("%s%04d", randomChars, id)
	return serialNumber, nil
}

// SetCounter 设置计数器（用于初始化）
func (g *SerialNumberGenerator) SetCounter(counter int32) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counter = counter
}

// generateRandomChars 生成指定长度的随机字符
func (g *SerialNumberGenerator) generateRandomChars(length int) (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 排除易混淆字符: 0,O,I,1
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

// ValidateSerialNumber 验证序列号格式
func ValidateSerialNumber(sn string) bool {
	if len(sn) != 8 {
		return false
	}

	// 前4位必须是字母或数字
	for i := 0; i < 4; i++ {
		c := sn[i]
		if !((c >= 'A' && c <= 'Z') || (c >= '2' && c <= '9')) {
			return false
		}
	}

	// 后4位必须是数字
	for i := 4; i < 8; i++ {
		c := sn[i]
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

// FormatSerialNumber 格式化序列号显示（4-4格式）
func FormatSerialNumber(sn string) string {
	if len(sn) != 8 {
		return sn
	}
	return sn[0:4] + " " + sn[4:8]
}

// ParseFormattedSerialNumber 解析格式化的序列号（去除空格）
func ParseFormattedSerialNumber(formattedSN string) string {
	result := make([]byte, 0, 8)
	for _, c := range formattedSN {
		if c != ' ' {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
