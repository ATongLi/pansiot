package core

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// ID范围常量
const (
	MinSystemID = 1         // 系统变量最小ID
	MaxSystemID = 99999     // 系统变量最大ID
	MinCustomID = 100000    // 自定义变量最小ID
	MaxCustomID = 999999999 // 自定义变量最大ID
)

// IDGenerator ID生成器
type IDGenerator struct {
	mu           sync.Mutex
	lastSystemID uint64
	lastCustomID uint64
	// 已使用的ID集合，用于快速查找
	usedIDs map[uint64]bool
}

// NewIDGenerator 创建新的ID生成器
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{
		lastSystemID: MinSystemID - 1,
		lastCustomID: MinCustomID - 1,
		usedIDs:      make(map[uint64]bool),
	}
}

// GenerateSystemID 生成系统变量ID
func (g *IDGenerator) GenerateSystemID() (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 查找下一个可用的系统ID
	for i := g.lastSystemID + 1; i <= MaxSystemID; i++ {
		if !g.usedIDs[i] {
			g.usedIDs[i] = true
			g.lastSystemID = i
			return i, nil
		}
	}

	return 0, fmt.Errorf("system ID pool exhausted")
}

// GenerateCustomID 生成自定义变量ID
func (g *IDGenerator) GenerateCustomID() (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 查找下一个可用的自定义ID
	for i := g.lastCustomID + 1; i <= MaxCustomID; i++ {
		if !g.usedIDs[i] {
			g.usedIDs[i] = true
			g.lastCustomID = i
			return i, nil
		}
	}

	return 0, fmt.Errorf("custom ID pool exhausted")
}

// MarkIDUsed 标记ID为已使用
func (g *IDGenerator) MarkIDUsed(id uint64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.usedIDs[id] = true

	// 更新last指针
	if id >= MinSystemID && id <= MaxSystemID && id > g.lastSystemID {
		g.lastSystemID = id
	}
	if id >= MinCustomID && id <= MaxCustomID && id > g.lastCustomID {
		g.lastCustomID = id
	}
}

// ReleaseID 释放ID
func (g *IDGenerator) ReleaseID(id uint64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.usedIDs, id)
}

// IsIDUsed 检查ID是否已使用
func (g *IDGenerator) IsIDUsed(id uint64) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.usedIDs[id]
}

// ParseStringID 解析字符串ID
// 字符串ID格式: <type>-<deviceId>-<variableId>
// 例如: DV-PLC001-TEMP01 表示 设备变量-设备PLC001-变量TEMP01
func ParseStringID(stringID string) (idType string, deviceID string, variableID string, err error) {
	parts := strings.Split(stringID, "-")
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid string ID format: %s, expected format: type-deviceId-variableId", stringID)
	}

	idType = parts[0]
	deviceID = parts[1]
	variableID = strings.Join(parts[2:], "-") // 支持variableID中包含"-"

	return idType, deviceID, variableID, nil
}

// BuildStringID 构建字符串ID
func BuildStringID(idType, deviceID, variableID string) string {
	return fmt.Sprintf("%s-%s-%s", idType, deviceID, variableID)
}

// ParseNumericStringID 从字符串ID中提取数字部分
// 例如: "DV-PLC001-TEMP01" -> 1 (类型: DV, 设备: PLC001, 变量: TEMP01)
func ParseNumericStringID(stringID string) (uint64, error) {
	// 这是一个简化的实现，实际应用中可能需要更复杂的解析逻辑
	// 这里可以根据需要实现从字符串ID到数字ID的映射
	h := fnvHash(stringID)
	return h, nil
}

// fnvHash FNV哈希算法，用于将字符串转换为数字
func fnvHash(s string) uint64 {
	const (
		offset64 uint64 = 14695981039346656037
		prime64  uint64 = 1099511628211
	)

	h := offset64
	for _, c := range []byte(s) {
		h ^= uint64(c)
		h *= prime64
	}
	return h
}

// StringIDComponents 字符串ID的组成部分
type StringIDComponents struct {
	Type       string // 类型: SV(系统变量), DV(设备变量), CV(计算变量) 等
	DeviceID   string // 设备ID
	VariableID string // 变量标识
	Index      int    // 可选的索引（用于数组类型变量）
}

// ParseStringIDV2 解析字符串ID（增强版）
func ParseStringIDV2(stringID string) (*StringIDComponents, error) {
	parts := strings.Split(stringID, "-")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid string ID format: %s", stringID)
	}

	comp := &StringIDComponents{
		Type:       parts[0],
		DeviceID:   parts[1],
		VariableID: parts[2],
	}

	// 如果有索引部分
	if len(parts) > 3 {
		idx, err := strconv.Atoi(parts[3])
		if err == nil {
			comp.Index = idx
		}
	}

	return comp, nil
}

// BuildStringIDV2 构建字符串ID（增强版）
func BuildStringIDV2(comp *StringIDComponents) string {
	if comp.Index > 0 {
		return fmt.Sprintf("%s-%s-%s-%d", comp.Type, comp.DeviceID, comp.VariableID, comp.Index)
	}
	return fmt.Sprintf("%s-%s-%s", comp.Type, comp.DeviceID, comp.VariableID)
}

// ValidateStringID 验证字符串ID格式
func ValidateStringID(stringID string) error {
	_, _, _, err := ParseStringID(stringID)
	return err
}

// GetIDType 从数字ID获取类型
func GetIDType(id uint64) string {
	if id >= MinSystemID && id <= MaxSystemID {
		return "system"
	}
	if id >= MinCustomID && id <= MaxCustomID {
		return "custom"
	}
	return "unknown"
}

// IsValidSystemID 检查是否为有效的系统ID
func IsValidSystemID(id uint64) bool {
	return id >= MinSystemID && id <= MaxSystemID
}

// IsValidCustomID 检查是否为有效的自定义ID
func IsValidCustomID(id uint64) bool {
	return id >= MinCustomID && id <= MaxCustomID
}

// IsValidID 检查是否为有效的ID
func IsValidID(id uint64) bool {
	return IsValidSystemID(id) || IsValidCustomID(id)
}
