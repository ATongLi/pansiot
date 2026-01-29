package engine

import (
	"fmt"
	"testing"
	"time"

	"pansiot-device/internal/alarm/rule"
	"pansiot-device/internal/core"
)

// MockStorage 模拟存储层
type MockStorage struct {
	variables map[uint64]*core.Variable
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		variables: make(map[uint64]*core.Variable),
	}
}

func (m *MockStorage) ReadVar(id uint64) (*core.Variable, error) {
	if v, ok := m.variables[id]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("变量不存在: %d", id)
}

func (m *MockStorage) ReadVars(ids []uint64) ([]*core.Variable, error) {
	result := make([]*core.Variable, 0, len(ids))
	for _, id := range ids {
		if v, ok := m.variables[id]; ok {
			result = append(result, v)
		}
	}
	return result, nil
}

func (m *MockStorage) ReadVarByStringID(stringID string) (*core.Variable, error) {
	return nil, fmt.Errorf("未实现")
}

func (m *MockStorage) WriteVar(variable *core.Variable) error {
	m.variables[variable.ID] = variable
	return nil
}

func (m *MockStorage) WriteVars(variables []*core.Variable) error {
	for _, v := range variables {
		m.variables[v.ID] = v
	}
	return nil
}

func (m *MockStorage) Subscribe(subscriberID string, variableIDs []uint64, callback func(core.VariableUpdate)) error {
	return nil
}

func (m *MockStorage) SubscribeByDevice(subscriberID, deviceID string, callback func(core.VariableUpdate)) error {
	return nil
}

func (m *MockStorage) SubscribeByPattern(subscriberID, pattern string, callback func(core.VariableUpdate)) error {
	return nil
}

func (m *MockStorage) Unsubscribe(subscriberID string, variableIDs []uint64) error {
	return nil
}

func (m *MockStorage) UnsubscribeAll(subscriberID string) error {
	return nil
}

func (m *MockStorage) GetStats() core.StorageStats {
	return core.StorageStats{}
}

func (m *MockStorage) CreateVariable(variable *core.Variable) error {
	m.variables[variable.ID] = variable
	return nil
}

func (m *MockStorage) DeleteVariable(variableID uint64) error {
	delete(m.variables, variableID)
	return nil
}

func (m *MockStorage) ListVariables() []*core.Variable {
	result := make([]*core.Variable, 0, len(m.variables))
	for _, v := range m.variables {
		result = append(result, v)
	}
	return result
}

func (m *MockStorage) ListVariablesByDevice(deviceID string) []*core.Variable {
	return []*core.Variable{}
}

// TestStateTracker 测试状态跟踪器
func TestStateTracker(t *testing.T) {
	st := NewStateTracker()

	// 测试更新状态
	st.UpdateState(100001, 85.5, true)

	state := st.GetState(100001)
	if state == nil {
		t.Fatalf("状态不应为nil")
	}

	if state.LastValue != 85.5 {
		t.Errorf("期望 LastValue=85.5, 实际=%v", state.LastValue)
	}

	if !state.LastState {
		t.Errorf("期望 LastState=true, 实际=false")
	}
}

// TestDeadbandFilter 测试死区过滤器
func TestDeadbandFilter(t *testing.T) {
	st := NewStateTracker()
	df := NewDeadbandFilter(st)

	// 创建带死区的条件
	cond := &rule.SingleCondition{
		VariableID: 100001,
		Operator:   rule.OpGTE,
		Deadband:   2.0,
	}

	tests := []struct {
		name        string
		value       float64
		threshold   float64
		shouldSatisfy bool
		inDeadband   bool
	}{
		{
			name:        "在死区内-保持上次状态",
			value:       79.0,
			threshold:   80.0,
			shouldSatisfy: false, // 首次，上次状态为false
			inDeadband:   true,
		},
		{
			name:        "超出死区-满足条件",
			value:       82.0,
			threshold:   80.0,
			shouldSatisfy: true,
			inDeadband:   false,
		},
		{
			name:        "超出死区-不满足条件",
			value:       77.0,
			threshold:   80.0,
			shouldSatisfy: false,
			inDeadband:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			satisfied, inDeadband := df.ApplyDeadband(cond, tt.value, tt.threshold)

			if satisfied != tt.shouldSatisfy {
				t.Errorf("期望 satisfied=%v, 实际=%v", tt.shouldSatisfy, satisfied)
			}

			if inDeadband != tt.inDeadband {
				t.Errorf("期望 inDeadband=%v, 实际=%v", tt.inDeadband, inDeadband)
			}
		})
	}
}

// TestEdgeDetector 测试边沿检测器
func TestEdgeDetector(t *testing.T) {
	st := NewStateTracker()
	ed := NewEdgeDetector(st)

	// 测试上升沿检测
	cond := &rule.SingleCondition{
		VariableID: 100002,
		Operator:   rule.OpRise,
	}

	// 第一次检测（无历史值）
	detected, _ := ed.DetectEdge(cond, true)
	if detected {
		t.Error("首次检测不应检测到边沿")
	}

	// 第二次检测（值仍为true，不应触发）
	detected, _ = ed.DetectEdge(cond, true)
	if detected {
		t.Error("值仍为true，不应触发上升沿")
	}

	// 第三次检测（false→true，上升沿）
	detected, edgeType := ed.DetectEdge(cond, false)
	if detected {
		t.Error("false→false不应检测到边沿")
	}

	detected, edgeType = ed.DetectEdge(cond, true)
	if !detected {
		t.Error("应检测到上升沿")
	}
	if edgeType != "rising" {
		t.Errorf("期望边沿类型=rising, 实际=%s", edgeType)
	}

	// 测试下降沿检测
	condFall := &rule.SingleCondition{
		VariableID: 100003,
		Operator:   rule.OpFall,
	}

	// 先设置为true
	ed.DetectEdge(condFall, true)

	// true→false，下降沿
	detected, edgeType = ed.DetectEdge(condFall, false)
	if !detected {
		t.Error("应检测到下降沿")
	}
	if edgeType != "falling" {
		t.Errorf("期望边沿类型=falling, 实际=%s", edgeType)
	}
}

// TestDelayTracker 测试延迟跟踪器
func TestDelayTracker(t *testing.T) {
	dt := NewDelayTracker()

	cond := &rule.SingleCondition{
		Delay: 100 * time.Millisecond,
	}

	triggered := false
	onTimeout := func() {
		triggered = true
	}

	// 启动延迟
	err := dt.StartDelay("RULE_001", cond, 100001, nil, onTimeout)
	if err != nil {
		t.Fatalf("启动延迟失败: %v", err)
	}

	// 检查是否活跃
	if !dt.IsActive("RULE_001") {
		t.Error("定时器应该是活跃状态")
	}

	// 等待延迟到期
	time.Sleep(150 * time.Millisecond)

	// 检查是否触发
	if !triggered {
		t.Error("延迟定时器应已触发")
	}
}

// TestEvaluator 测试评估引擎
func TestEvaluator(t *testing.T) {
	storage := NewMockStorage()
	evaluator := NewEvaluator(storage)

	// 设置变量值
	storage.WriteVar(&core.Variable{
		ID:      100001,
		Value:   85.5,
		Quality: core.QualityGood,
	})

	// 测试简单条件评估
	cond := &rule.SingleCondition{
		VariableID: 100001,
		Operator:   rule.OpGT,
		Value:      80.0,
	}

	satisfied, err := evaluator.EvaluateWithState(cond)
	if err != nil {
		t.Fatalf("评估失败: %v", err)
	}

	if !satisfied {
		t.Error("条件应该满足")
	}
}

// TestEvaluatorWithDeadband 测试带死区的评估
func TestEvaluatorWithDeadband(t *testing.T) {
	storage := NewMockStorage()
	evaluator := NewEvaluator(storage)

	// 设置变量值
	storage.WriteVar(&core.Variable{
		ID:      100001,
		Value:   79.0,
		Quality: core.QualityGood,
	})

	// 创建带死区的条件
	cond := &rule.SingleCondition{
		VariableID: 100001,
		Operator:   rule.OpGTE,
		Value:      80.0,
		Deadband:   2.0,
	}

	// 首次评估（在死区内）
	satisfied, err := evaluator.EvaluateWithState(cond)
	if err != nil {
		t.Fatalf("评估失败: %v", err)
	}

	if satisfied {
		t.Error("在死区内且上次状态为false，应返回false")
	}

	// 第二次评估（值超出死区）
	storage.WriteVar(&core.Variable{
		ID:      100001,
		Value:   82.0,
		Quality: core.QualityGood,
	})

	satisfied, err = evaluator.EvaluateWithState(cond)
	if err != nil {
		t.Fatalf("评估失败: %v", err)
	}

	if !satisfied {
		t.Error("值82.0应满足条件（>=80.0）")
	}
}

// TestEvaluatorWithEdgeDetection 测试边沿检测评估
func TestEvaluatorWithEdgeDetection(t *testing.T) {
	storage := NewMockStorage()
	evaluator := NewEvaluator(storage)

	// 测试上升沿条件
	cond := &rule.SingleCondition{
		VariableID: 100002,
		Operator:   rule.OpRise,
	}

	// 第一次：false（无历史值）
	storage.WriteVar(&core.Variable{
		ID:      100002,
		Value:   true,
		Quality: core.QualityGood,
	})

	satisfied, err := evaluator.EvaluateWithState(cond)
	if err != nil {
		t.Fatalf("评估失败: %v", err)
	}

	if satisfied {
		t.Error("首次检测不应触发上升沿")
	}

	// 第二次：false
	storage.WriteVar(&core.Variable{
		ID:      100002,
		Value:   true,
		Quality: core.QualityGood,
	})

	satisfied, err = evaluator.EvaluateWithState(cond)
	if err != nil {
		t.Fatalf("评估失败: %v", err)
	}

	if satisfied {
		t.Error("值仍为true，不应触发上升沿")
	}

	// 第三次：false→true（上升沿）
	storage.WriteVar(&core.Variable{
		ID:      100002,
		Value:   false,
		Quality: core.QualityGood,
	})

	evaluator.EvaluateWithState(cond) // 记录false值

	storage.WriteVar(&core.Variable{
		ID:      100002,
		Value:   true,
		Quality: core.QualityGood,
	})

	satisfied, err = evaluator.EvaluateWithState(cond)
	if err != nil {
		t.Fatalf("评估失败: %v", err)
	}

	if !satisfied {
		t.Error("false→true应触发上升沿")
	}
}

// TestConditionGroupEvaluation 测试条件组评估
func TestConditionGroupEvaluation(t *testing.T) {
	storage := NewMockStorage()
	evaluator := NewEvaluator(storage)

	// 设置变量值
	storage.WriteVar(&core.Variable{
		ID:      100001,
		Value:   85.0,
		Quality: core.QualityGood,
	})
	storage.WriteVar(&core.Variable{
		ID:      100002,
		Value:   95.0,
		Quality: core.QualityGood,
	})

	// 测试AND条件组
	andGroup := &rule.ConditionGroup{
		Logic: rule.LogicAND,
		Conditions: []rule.Condition{
			&rule.SingleCondition{
				VariableID: 100001,
				Operator:   rule.OpGT,
				Value:      80.0,
			},
			&rule.SingleCondition{
				VariableID: 100002,
				Operator:   rule.OpGT,
				Value:      90.0,
			},
		},
	}

	satisfied, err := evaluator.EvaluateWithState(andGroup)
	if err != nil {
		t.Fatalf("评估失败: %v", err)
	}

	if !satisfied {
		t.Error("两个条件都应满足")
	}

	// 测试OR条件组
	orGroup := &rule.ConditionGroup{
		Logic: rule.LogicOR,
		Conditions: []rule.Condition{
			&rule.SingleCondition{
				VariableID: 100001,
				Operator:   rule.OpLT,
				Value:      50.0,
			},
			&rule.SingleCondition{
				VariableID: 100002,
				Operator:   rule.OpGT,
				Value:      90.0,
			},
		},
	}

	satisfied, err = evaluator.EvaluateWithState(orGroup)
	if err != nil {
		t.Fatalf("评估失败: %v", err)
	}

	if !satisfied {
		t.Error("OR条件组应满足（第二个条件为true）")
	}
}

// BenchmarkSingleConditionEvaluation 基准测试：单个条件评估
func BenchmarkSingleConditionEvaluation(b *testing.B) {
	storage := NewMockStorage()
	storage.WriteVar(&core.Variable{
		ID:      100001,
		Value:   85.5,
		Quality: core.QualityGood,
	})

	evaluator := NewEvaluator(storage)
	cond := &rule.SingleCondition{
		VariableID: 100001,
		Operator:   rule.OpGT,
		Value:      80.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evaluator.EvaluateWithState(cond)
	}
}
