package engine

import (
	"fmt"
	"sync"
	"time"

	"pansiot-device/internal/alarm/rule"
	"pansiot-device/internal/core"
)

// AlarmStateMachine 报警状态机
// 集中管理报警规则的状态转换规则和历史记录
type AlarmStateMachine struct {
	mu          sync.RWMutex
	tracker     *AlarmStateTracker            // 复用现有状态跟踪器
	transitions map[string][]StateTransition  // ruleID -> 转换历史
}

// StateTransition 状态转换记录
type StateTransition struct {
	From      core.AlarmState
	To        core.AlarmState
	Timestamp time.Time
	Reason    string // 转换原因（trigger/ack/recover/force等）
	UserID    string // 操作用户（确认操作时）
}

// ActiveAlarm 活跃报警实例
type ActiveAlarm struct {
	RuleID       string
	Rule         *rule.AlarmRule
	State        core.AlarmState
	TriggerTime  time.Time
	TriggerValue interface{}
	AckUser      string
	AckTime      *time.Time
	RecoverTime  *time.Time
}

// StateMachine 状态机接口（用于动作调度器检查状态）
type StateMachine interface {
	GetState(ruleID string) core.AlarmState
}

// NewAlarmStateMachine 创建报警状态机
func NewAlarmStateMachine(tracker *AlarmStateTracker) *AlarmStateMachine {
	if tracker == nil {
		tracker = NewAlarmStateTracker()
	}
	return &AlarmStateMachine{
		tracker:     tracker,
		transitions: make(map[string][]StateTransition),
	}
}

// TransitionTo 尝试转换到新状态
// 参数：
//   - ruleID: 规则ID
//   - newState: 目标状态
//   - reason: 转换原因
//   - userID: 操作用户（可选，用于确认操作）
// 返回：
//   - error: 转换失败时返回错误
func (sm *AlarmStateMachine) TransitionTo(ruleID string, newState core.AlarmState, reason string, userID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 获取当前状态
	currentState := sm.tracker.GetState(ruleID)

	// 验证状态转换是否合法
	if !sm.canTransition(currentState, newState) {
		return fmt.Errorf("非法的状态转换: %d -> %d (ruleID=%s)",
			currentState, newState, ruleID)
	}

	// 记录转换历史
	transition := StateTransition{
		From:      currentState,
		To:        newState,
		Timestamp: time.Now(),
		Reason:    reason,
		UserID:    userID,
	}
	sm.transitions[ruleID] = append(sm.transitions[ruleID], transition)

	// 更新状态
	sm.tracker.SetState(ruleID, newState)

	return nil
}

// CanTransition 检查是否允许转换
func (sm *AlarmStateMachine) CanTransition(ruleID string, newState core.AlarmState) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	currentState := sm.tracker.GetState(ruleID)
	return sm.canTransition(currentState, newState)
}

// canTransition 内部方法：检查状态转换是否合法
func (sm *AlarmStateMachine) canTransition(from, to core.AlarmState) bool {
	// 允许的转换规则
	validTransitions := map[core.AlarmState][]core.AlarmState{
		core.AlarmStateInactive: {
			core.AlarmStateActive, // Inactive → Active
		},
		core.AlarmStateActive: {
			core.AlarmStateAcknowledged, // Active → Acknowledged
			core.AlarmStateCleared,      // Active → Cleared
		},
		core.AlarmStateAcknowledged: {
			core.AlarmStateCleared, // Acknowledged → Cleared
		},
		core.AlarmStateCleared: {
			core.AlarmStateActive, // Cleared → Active (重新触发)
		},
	}

	allowedStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowed := range allowedStates {
		if allowed == to {
			return true
		}
	}
	return false
}

// ForceTransition 强制转换到指定状态（用于初始化或恢复）
func (sm *AlarmStateMachine) ForceTransition(ruleID string, newState core.AlarmState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	currentState := sm.tracker.GetState(ruleID)

	// 记录转换历史
	transition := StateTransition{
		From:      currentState,
		To:        newState,
		Timestamp: time.Now(),
		Reason:    "force",
	}
	sm.transitions[ruleID] = append(sm.transitions[ruleID], transition)

	// 更新状态
	sm.tracker.SetState(ruleID, newState)
}

// GetState 获取当前状态
func (sm *AlarmStateMachine) GetState(ruleID string) core.AlarmState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.tracker.GetState(ruleID)
}

// GetTransitionHistory 获取转换历史
func (sm *AlarmStateMachine) GetTransitionHistory(ruleID string) []StateTransition {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	transitions := sm.transitions[ruleID]
	// 返回副本，避免外部修改
	result := make([]StateTransition, len(transitions))
	copy(result, transitions)
	return result
}

// GetLastTransition 获取最后一次转换
func (sm *AlarmStateMachine) GetLastTransition(ruleID string) *StateTransition {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	transitions := sm.transitions[ruleID]
	if len(transitions) == 0 {
		return nil
	}

	// 返回副本
	last := transitions[len(transitions)-1]
	return &last
}

// IsActive 检查报警是否处于激活状态（Active 或 Acknowledged）
func (sm *AlarmStateMachine) IsActive(ruleID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.tracker.IsActive(ruleID)
}

// IsAcknowledged 检查报警是否已确认
func (sm *AlarmStateMachine) IsAcknowledged(ruleID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.tracker.IsAcknowledged(ruleID)
}

// IsCleared 检查报警是否已清除
func (sm *AlarmStateMachine) IsCleared(ruleID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.tracker.IsCleared(ruleID)
}

// GetHistory 获取报警历史记录（包含时间信息）
func (sm *AlarmStateMachine) GetHistory(ruleID string) *AlarmHistory {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.tracker.GetHistory(ruleID)
}

// Clear 清除指定规则的状态
func (sm *AlarmStateMachine) Clear(ruleID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.tracker.Clear(ruleID)
	delete(sm.transitions, ruleID)
}

// ClearAll 清除所有状态
func (sm *AlarmStateMachine) ClearAll() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.tracker.ClearAll()
	sm.transitions = make(map[string][]StateTransition)
}

// GetAllStates 获取所有规则的状态
func (sm *AlarmStateMachine) GetAllStates() map[string]core.AlarmState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// 注意：这里需要从 AlarmStateTracker 的内部状态获取
	// 由于 AlarmStateTracker 没有提供 GetAllStates 方法，
	// 我们需要维护一个额外的映射或修改 AlarmStateTracker
	// 暂时返回空映射，后续可优化
	return make(map[string]core.AlarmState)
}

// GetStats 获取状态统计信息
type StateStats struct {
	TotalRules      int
	InactiveCount   int
	ActiveCount     int
	AcknowledgedCount int
	ClearedCount    int
}

func (sm *AlarmStateMachine) GetStats() StateStats {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// 由于无法直接访问所有规则ID，这里返回默认值
	// 实际使用时需要从外部传入规则列表
	return StateStats{}
}

// ValidateTransition 验证状态转换并返回详细的错误信息
func (sm *AlarmStateMachine) ValidateTransition(ruleID string, newState core.AlarmState) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	currentState := sm.tracker.GetState(ruleID)

	if currentState == newState {
		return fmt.Errorf("状态未改变: %d", currentState)
	}

	if !sm.canTransition(currentState, newState) {
		return fmt.Errorf("不允许的状态转换: %d -> %d", currentState, newState)
	}

	return nil
}

// GetAllowedTransitions 获取从当前状态允许的所有转换
func (sm *AlarmStateMachine) GetAllowedTransitions(ruleID string) []core.AlarmState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	currentState := sm.tracker.GetState(ruleID)

	validTransitions := map[core.AlarmState][]core.AlarmState{
		core.AlarmStateInactive:     {core.AlarmStateActive},
		core.AlarmStateActive:       {core.AlarmStateAcknowledged, core.AlarmStateCleared},
		core.AlarmStateAcknowledged: {core.AlarmStateCleared},
		core.AlarmStateCleared:      {core.AlarmStateActive},
	}

	allowed, exists := validTransitions[currentState]
	if !exists {
		return []core.AlarmState{}
	}

	// 返回副本
	result := make([]core.AlarmState, len(allowed))
	copy(result, allowed)
	return result
}
