package engine

import (
	"fmt"
	"math"

	"pansiot-device/internal/alarm/rule"
)

// DeadbandFilter 死区过滤器
// 防止变量在阈值附近抖动时频繁触发报警
type DeadbandFilter struct {
	tracker *StateTracker
}

// NewDeadbandFilter 创建死区过滤器
func NewDeadbandFilter(tracker *StateTracker) *DeadbandFilter {
	return &DeadbandFilter{
		tracker: tracker,
	}
}

// ApplyDeadband 应用死区逻辑
// 参数：
//   - cond: 单个条件
//   - currentValue: 当前变量值
//   - threshold: 阈值
// 返回：
//   - satisfied: 是否满足条件
//   - inDeadband: 是否在死区内
func (df *DeadbandFilter) ApplyDeadband(cond *rule.SingleCondition, currentValue float64, threshold float64) (satisfied bool, inDeadband bool) {
	// 如果没有设置死区，直接比较
	if cond.Deadband <= 0 {
		return df.compareDirectly(cond, currentValue, threshold), false
	}

	// 获取或创建状态
	state := df.tracker.GetOrCreateState(cond.VariableID)

	// 计算死区范围
	lowerBound := threshold - cond.Deadband
	upperBound := threshold + cond.Deadband

	// 根据操作符类型，智能判断是否在死区内
	var inRange bool
	switch cond.Operator {
	case rule.OpGTE:
		// >= 操作符：只有当值在 [lowerBound, threshold) 时才在死区
		// 达到 threshold 及以上应该触发
		inRange = currentValue >= lowerBound && currentValue < threshold
	case rule.OpGT:
		// > 操作符：只有当值在 [lowerBound, threshold] 时才在死区
		// 超过 threshold 才触发
		inRange = currentValue >= lowerBound && currentValue <= threshold
	case rule.OpLTE:
		// <= 操作符：只有当值在 (threshold, upperBound] 时才在死区
		// 达到 threshold 及以下应该触发
		inRange = currentValue > threshold && currentValue <= upperBound
	case rule.OpLT:
		// < 操作符：只有当值在 [threshold, upperBound] 时才在死区
		// 低于 threshold 才触发
		inRange = currentValue >= threshold && currentValue <= upperBound
	default:
		// 其他操作符使用标准死区判断
		inRange = currentValue >= lowerBound && currentValue <= upperBound
	}

	if inRange {
		// 在死区内，保持上次状态
		df.tracker.SetDeadbandStatus(cond.VariableID, true)
		return state.LastState, true
	}

	// 超出死区，执行实际比较
	satisfied = df.compareDirectly(cond, currentValue, threshold)
	df.tracker.SetDeadbandStatus(cond.VariableID, false)

	return satisfied, false
}

// compareDirectly 直接比较（不考虑死区）
func (df *DeadbandFilter) compareDirectly(cond *rule.SingleCondition, currentValue float64, threshold float64) bool {
	switch cond.Operator {
	case rule.OpGT:
		return currentValue > threshold
	case rule.OpLT:
		return currentValue < threshold
	case rule.OpGTE:
		return currentValue >= threshold
	case rule.OpLTE:
		return currentValue <= threshold
	case rule.OpEQ:
		return math.Abs(currentValue-threshold) < 1e-9 // 浮点数相等判断
	case rule.OpNEQ:
		return math.Abs(currentValue-threshold) >= 1e-9
	default:
		return false
	}
}

// ApplyDeadbandWithValueRef 应用死区逻辑（支持动态阈值）
func (df *DeadbandFilter) ApplyDeadbandWithValueRef(cond *rule.SingleCondition, currentValue float64, thresholdRefValue float64) (satisfied bool, inDeadband bool, err error) {
	// 如果没有设置死区，直接比较
	if cond.Deadband <= 0 {
		result := df.compareDirectly(cond, currentValue, thresholdRefValue)
		return result, false, nil
	}

	// 获取或创建状态
	state := df.tracker.GetOrCreateState(cond.VariableID)

	// 计算死区范围
	lowerBound := thresholdRefValue - cond.Deadband
	upperBound := thresholdRefValue + cond.Deadband

	// 检查是否在死区内
	if currentValue >= lowerBound && currentValue <= upperBound {
		// 在死区内，保持上次状态
		df.tracker.SetDeadbandStatus(cond.VariableID, true)
		return state.LastState, true, nil
	}

	// 超出死区，执行实际比较
	satisfied = df.compareDirectly(cond, currentValue, thresholdRefValue)
	df.tracker.SetDeadbandStatus(cond.VariableID, false)

	return satisfied, false, nil
}

// Reset 重置死区状态
func (df *DeadbandFilter) Reset(variableID uint64) {
	df.tracker.SetDeadbandStatus(variableID, false)
}

// GetDeadbandStatus 获取死区状态
func (df *DeadbandFilter) GetDeadbandStatus(variableID uint64) bool {
	return df.tracker.IsInDeadband(variableID)
}

// IsInDeadband 检查值是否在死区内
func (df *DeadbandFilter) IsInDeadband(cond *rule.SingleCondition, value float64, threshold float64) bool {
	if cond.Deadband <= 0 {
		return false
	}

	lowerBound := threshold - cond.Deadband
	upperBound := threshold + cond.Deadband

	return value >= lowerBound && value <= upperBound
}

// CalculateDeadbandBounds 计算死区边界
func (df *DeadbandFilter) CalculateDeadbandBounds(threshold float64, deadband float64) (lowerBound, upperBound float64, err error) {
	if deadband < 0 {
		return 0, 0, fmt.Errorf("死区不能为负数: %f", deadband)
	}

	return threshold - deadband, threshold + deadband, nil
}
