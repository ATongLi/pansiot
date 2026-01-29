package engine

import (
	"sync"
	"time"

	"pansiot-device/internal/core"
)

// ConditionState 条件状态（用于状态跟踪）
type ConditionState struct {
	VariableID uint64        // 变量ID
	LastValue  interface{}    // 上次的值
	LastState  bool          // 上次是否满足条件
	LastUpdate time.Time     // 上次更新时间
	InDeadband bool          // 是否在死区内
}

// StateTracker 状态跟踪器
// 跟踪变量的历史值和状态，用于死区过滤、边沿检测等
type StateTracker struct {
	mu     sync.RWMutex
	states map[uint64]*ConditionState // variableID -> 状态
}

// NewStateTracker 创建状态跟踪器
func NewStateTracker() *StateTracker {
	return &StateTracker{
		states: make(map[uint64]*ConditionState),
	}
}

// GetState 获取变量状态
func (st *StateTracker) GetState(variableID uint64) *ConditionState {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.states[variableID]
}

// GetOrCreateState 获取或创建变量状态
func (st *StateTracker) GetOrCreateState(variableID uint64) *ConditionState {
	st.mu.Lock()
	defer st.mu.Unlock()

	state, exists := st.states[variableID]
	if !exists {
		state = &ConditionState{
			VariableID: variableID,
			LastState:  false,
			LastUpdate: time.Now(),
		}
		st.states[variableID] = state
	}
	return state
}

// UpdateState 更新变量状态
func (st *StateTracker) UpdateState(variableID uint64, value interface{}, satisfied bool) {
	st.mu.Lock()
	defer st.mu.Unlock()

	state := st.states[variableID]
	if state == nil {
		state = &ConditionState{
			VariableID: variableID,
		}
		st.states[variableID] = state
	}

	state.LastValue = value
	state.LastState = satisfied
	state.LastUpdate = time.Now()
}

// GetLastValue 获取变量上次的值
func (st *StateTracker) GetLastValue(variableID uint64) interface{} {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if state, ok := st.states[variableID]; ok {
		return state.LastValue
	}
	return nil
}

// GetLastState 获取变量上次是否满足条件
func (st *StateTracker) GetLastState(variableID uint64) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if state, ok := st.states[variableID]; ok {
		return state.LastState
	}
	return false
}

// SetDeadbandStatus 设置死区状态
func (st *StateTracker) SetDeadbandStatus(variableID uint64, inDeadband bool) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if state, ok := st.states[variableID]; ok {
		state.InDeadband = inDeadband
	}
}

// IsInDeadband 检查是否在死区内
func (st *StateTracker) IsInDeadband(variableID uint64) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if state, ok := st.states[variableID]; ok {
		return state.InDeadband
	}
	return false
}

// Clear 清除所有状态
func (st *StateTracker) Clear() {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.states = make(map[uint64]*ConditionState)
}

// ClearVariable 清除指定变量的状态
func (st *StateTracker) ClearVariable(variableID uint64) {
	st.mu.Lock()
	defer st.mu.Unlock()
	delete(st.states, variableID)
}

// AlarmStateTracker 报警状态跟踪器
// 跟踪每个报警规则的状态
type AlarmStateTracker struct {
	mu      sync.RWMutex
	states  map[string]core.AlarmState // ruleID -> 报警状态
	history map[string]*AlarmHistory   // ruleID -> 历史记录
}

// AlarmHistory 报警历史记录
type AlarmHistory struct {
	RuleID        string
	TriggerTime   time.Time
	TriggerValue  interface{}
	AckTime       *time.Time
	RecoverTime   *time.Time
	TriggerCount  int // 触发次数
}

// NewAlarmStateTracker 创建报警状态跟踪器
func NewAlarmStateTracker() *AlarmStateTracker {
	return &AlarmStateTracker{
		states:  make(map[string]core.AlarmState),
		history: make(map[string]*AlarmHistory),
	}
}

// GetState 获取报警状态
func (ast *AlarmStateTracker) GetState(ruleID string) core.AlarmState {
	ast.mu.RLock()
	defer ast.mu.RUnlock()

	if state, ok := ast.states[ruleID]; ok {
		return state
	}
	return core.AlarmStateInactive
}

// SetState 设置报警状态
func (ast *AlarmStateTracker) SetState(ruleID string, state core.AlarmState) {
	ast.mu.Lock()
	defer ast.mu.Unlock()

	ast.states[ruleID] = state

	// 更新历史记录
	history := ast.history[ruleID]
	if history == nil {
		history = &AlarmHistory{
			RuleID: ruleID,
		}
		ast.history[ruleID] = history
	}

	// 状态转换记录
	switch state {
	case core.AlarmStateActive:
		if history.TriggerTime.IsZero() {
			history.TriggerTime = time.Now()
			history.TriggerCount++
		}
	case core.AlarmStateAcknowledged:
		if history.AckTime == nil {
			now := time.Now()
			history.AckTime = &now
		}
	case core.AlarmStateCleared:
		if history.RecoverTime == nil {
			now := time.Now()
			history.RecoverTime = &now
		}
	}
}

// GetHistory 获取报警历史记录
func (ast *AlarmStateTracker) GetHistory(ruleID string) *AlarmHistory {
	ast.mu.RLock()
	defer ast.mu.RUnlock()
	return ast.history[ruleID]
}

// IsActive 检查报警是否处于激活状态
func (ast *AlarmStateTracker) IsActive(ruleID string) bool {
	state := ast.GetState(ruleID)
	return state == core.AlarmStateActive || state == core.AlarmStateAcknowledged
}

// IsAcknowledged 检查报警是否已确认
func (ast *AlarmStateTracker) IsAcknowledged(ruleID string) bool {
	return ast.GetState(ruleID) == core.AlarmStateAcknowledged
}

// IsCleared 检查报警是否已清除
func (ast *AlarmStateTracker) IsCleared(ruleID string) bool {
	return ast.GetState(ruleID) == core.AlarmStateCleared
}

// Clear 清除指定规则的状态
func (ast *AlarmStateTracker) Clear(ruleID string) {
	ast.mu.Lock()
	defer ast.mu.Unlock()
	delete(ast.states, ruleID)
	delete(ast.history, ruleID)
}

// ClearAll 清除所有状态
func (ast *AlarmStateTracker) ClearAll() {
	ast.mu.Lock()
	defer ast.mu.Unlock()
	ast.states = make(map[string]core.AlarmState)
	ast.history = make(map[string]*AlarmHistory)
}
