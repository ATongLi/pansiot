package engine

import (
	"fmt"
	"sync"

	"pansiot-device/internal/alarm/rule"
)

// EdgeDetector 边沿检测器
// 检测变量的跳变边沿（0→1上升沿，1→0下降沿）
type EdgeDetector struct {
	tracker *StateTracker
	mu      sync.RWMutex
}

// NewEdgeDetector 创建边沿检测器
func NewEdgeDetector(tracker *StateTracker) *EdgeDetector {
	return &EdgeDetector{
		tracker: tracker,
	}
}

// DetectEdge 检测边沿
// 参数：
//   - cond: 单个条件
//   - currentValue: 当前变量值
// 返回：
//   - edgeDetected: 是否检测到边沿
//   - edgeType: 边沿类型（"rising", "falling", "none"）
func (ed *EdgeDetector) DetectEdge(cond *rule.SingleCondition, currentValue interface{}) (edgeDetected bool, edgeType string) {
	// 只有边沿操作符才需要边沿检测
	if cond.Operator < rule.OpRise || cond.Operator > rule.OpFall {
		return false, "none"
	}

	// 转换当前值为bool
	currentBool, err := toBool(currentValue)
	if err != nil {
		return false, "none"
	}

	// 获取上次的值
	state := ed.tracker.GetOrCreateState(cond.VariableID)
	lastValue := state.LastValue

	// 如果没有历史值，记录当前值并返回false
	if lastValue == nil {
		ed.tracker.UpdateState(cond.VariableID, currentValue, false)
		return false, "none"
	}

	// 转换上次值为bool
	lastBool, err := toBool(lastValue)
	if err != nil {
		ed.tracker.UpdateState(cond.VariableID, currentValue, false)
		return false, "none"
	}

	// 检测边沿
	var detected bool
	var edge string

	switch cond.Operator {
	case rule.OpRise:
		// 0→1 上升沿：上次为false，当前为true
		detected = !lastBool && currentBool
		edge = "rising"

	case rule.OpFall:
		// 1→0 下降沿：上次为true，当前为false
		detected = lastBool && !currentBool
		edge = "falling"

	default:
		detected = false
		edge = "none"
	}

	// 更新状态
	ed.tracker.UpdateState(cond.VariableID, currentValue, detected)

	return detected, edge
}

// DetectRisingEdge 检测上升沿（0→1）
func (ed *EdgeDetector) DetectRisingEdge(variableID uint64, currentValue bool) bool {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	state := ed.tracker.GetOrCreateState(variableID)
	lastValue := state.LastValue

	// 如果没有历史值，记录当前值并返回false
	if lastValue == nil {
		ed.tracker.UpdateState(variableID, currentValue, false)
		return false
	}

	lastBool, ok := lastValue.(bool)
	if !ok {
		ed.tracker.UpdateState(variableID, currentValue, false)
		return false
	}

	// 0→1 上升沿
	if !lastBool && currentValue {
		ed.tracker.UpdateState(variableID, currentValue, true)
		return true
	}

	ed.tracker.UpdateState(variableID, currentValue, false)
	return false
}

// DetectFallingEdge 检测下降沿（1→0）
func (ed *EdgeDetector) DetectFallingEdge(variableID uint64, currentValue bool) bool {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	state := ed.tracker.GetOrCreateState(variableID)
	lastValue := state.LastValue

	// 如果没有历史值，记录当前值并返回false
	if lastValue == nil {
		ed.tracker.UpdateState(variableID, currentValue, false)
		return false
	}

	lastBool, ok := lastValue.(bool)
	if !ok {
		ed.tracker.UpdateState(variableID, currentValue, false)
		return false
	}

	// 1→0 下降沿
	if lastBool && !currentValue {
		ed.tracker.UpdateState(variableID, currentValue, true)
		return true
	}

	ed.tracker.UpdateState(variableID, currentValue, false)
	return false
}

// DetectAnyEdge 检测任意边沿（上升或下降）
func (ed *EdgeDetector) DetectAnyEdge(variableID uint64, currentValue bool) bool {
	return ed.DetectRisingEdge(variableID, currentValue) ||
		ed.DetectFallingEdge(variableID, currentValue)
}

// HasHistory 检查是否有历史值
func (ed *EdgeDetector) HasHistory(variableID uint64) bool {
	state := ed.tracker.GetState(variableID)
	return state != nil && state.LastValue != nil
}

// Reset 重置边沿检测状态
func (ed *EdgeDetector) Reset(variableID uint64) {
	ed.tracker.ClearVariable(variableID)
}

// ResetAll 重置所有边沿检测状态
func (ed *EdgeDetector) ResetAll() {
	ed.tracker.Clear()
}

// toBool 转换为bool（辅助函数）
func toBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case int, int8, int16, int32, int64:
		return value != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		return value != 0, nil
	case float32:
		return v != 0, nil
	case float64:
		return v != 0, nil
	default:
		return false, fmt.Errorf("不支持的类型: %T", value)
	}
}
