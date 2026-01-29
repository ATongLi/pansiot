package alarm

import (
	"testing"
	"time"

	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/alarm/rule"
	"pansiot-device/internal/core"
)

// TestAlarmConsumerIntegration 测试报警消费者的完整流程
func TestAlarmConsumerIntegration(t *testing.T) {
	// 创建模拟存储
	storage := NewMockStorage()

	// 设置初始变量值
	storage.WriteVar(&core.Variable{
		ID:      100001,
		Value:   75.0,
		Quality: core.QualityGood,
	})

	// 创建报警消费者
	config := &AlarmConsumerConfig{
		EvalWorkers:   2,
		EvalTimeout:   100 * time.Millisecond,
		RecoverDelay:  1 * time.Second,
		AutoSubscribe: false, // 测试时不自动订阅
	}
	consumer := NewAlarmConsumer(storage, config)

	// 验证组件已正确初始化
	if consumer.GetEvaluator() == nil {
		t.Fatal("评估引擎未初始化")
	}
	if consumer.GetStateMachine() == nil {
		t.Fatal("状态机未初始化")
	}
	if consumer.GetRuleManager() == nil {
		t.Fatal("规则管理器未初始化")
	}

	t.Log("报警消费者初始化成功")
}

// TestStateMachineTransitions 测试状态机状态转换
func TestStateMachineTransitions(t *testing.T) {
	sm := engine.NewAlarmStateMachine(nil)
	ruleID := "TEST_RULE_001"

	// 测试1: Inactive -> Active (合法)
	err := sm.TransitionTo(ruleID, core.AlarmStateActive, "trigger", "")
	if err != nil {
		t.Fatalf("Inactive -> Active 转换失败: %v", err)
	}

	currentState := sm.GetState(ruleID)
	if currentState != core.AlarmStateActive {
		t.Errorf("期望状态 Active，实际 %d", currentState)
	}

	// 测试2: Active -> Acknowledged (合法)
	err = sm.TransitionTo(ruleID, core.AlarmStateAcknowledged, "acknowledge", "user001")
	if err != nil {
		t.Fatalf("Active -> Acknowledged 转换失败: %v", err)
	}

	currentState = sm.GetState(ruleID)
	if currentState != core.AlarmStateAcknowledged {
		t.Errorf("期望状态 Acknowledged，实际 %d", currentState)
	}

	// 测试3: Acknowledged -> Cleared (合法)
	err = sm.TransitionTo(ruleID, core.AlarmStateCleared, "recover", "")
	if err != nil {
		t.Fatalf("Acknowledged -> Cleared 转换失败: %v", err)
	}

	currentState = sm.GetState(ruleID)
	if currentState != core.AlarmStateCleared {
		t.Errorf("期望状态 Cleared，实际 %d", currentState)
	}

	// 测试4: Cleared -> Active (重新触发，合法)
	err = sm.TransitionTo(ruleID, core.AlarmStateActive, "trigger", "")
	if err != nil {
		t.Fatalf("Cleared -> Active 转换失败: %v", err)
	}

	// 验证转换历史
	history := sm.GetTransitionHistory(ruleID)
	if len(history) != 4 {
		t.Errorf("期望 4 次状态转换，实际 %d", len(history))
	}

	t.Logf("状态转换测试通过，转换历史: %d 次", len(history))
}

// TestStateMachineInvalidTransition 测试非法状态转换
func TestStateMachineInvalidTransition(t *testing.T) {
	sm := engine.NewAlarmStateMachine(nil)
	ruleID := "TEST_RULE_002"

	// 设置初始状态为 Active
	sm.ForceTransition(ruleID, core.AlarmStateActive)

	// 测试: Active -> Inactive (非法)
	err := sm.TransitionTo(ruleID, core.AlarmStateInactive, "invalid", "")
	if err == nil {
		t.Fatal("Active -> Inactive 应该是非法转换，但没有返回错误")
	}

	t.Log("非法状态转换被正确拒绝")
}

// TestVariableIndex 测试变量索引
func TestVariableIndex(t *testing.T) {
	index := rule.NewVariableRuleIndex()

	// 创建测试规则
	rule1 := &rule.AlarmRule{
		ID:       "RULE_001",
		Name:     "测试规则1",
		Category: "TEMP",
		Level:    3,
		Enabled:  true,
		Condition: &rule.SingleCondition{
			VariableID: 100001,
			Operator:   rule.OpGTE,
			Value:      80.0,
		},
		TriggerMsg: rule.AlarmMessage{
			Type:    rule.ContentStatic,
			Content: "温度过高",
		},
	}

	rule2 := &rule.AlarmRule{
		ID:       "RULE_002",
		Name:     "测试规则2",
		Category: "TEMP",
		Level:    2,
		Enabled:  true,
		Condition: &rule.SingleCondition{
			VariableID: 100001, // 与 rule1 共享变量
			Operator:   rule.OpGT,
			Value:      90.0,
		},
		TriggerMsg: rule.AlarmMessage{
			Type:    rule.ContentStatic,
			Content: "温度严重过高",
		},
	}

	// 添加规则到索引
	index.AddRule(rule1)
	index.AddRule(rule2)

	// 测试1: 获取引用变量的规则
	ruleIDs := index.GetRulesByVariable(100001)
	if len(ruleIDs) != 2 {
		t.Errorf("期望 2 个规则引用变量 100001，实际 %d", len(ruleIDs))
	}

	// 测试2: 获取规则涉及的变量
	vars := index.GetVariablesByRule("RULE_001")
	if len(vars) != 1 {
		t.Errorf("期望 RULE_001 涉及 1 个变量，实际 %d", len(vars))
	}

	// 测试3: 查找共享变量
	shared := index.FindSharedVariables("RULE_001", "RULE_002")
	if len(shared) != 1 {
		t.Errorf("期望 1 个共享变量，实际 %d", len(shared))
	}

	// 测试4: 删除规则
	index.RemoveRule("RULE_001")
	ruleIDs = index.GetRulesByVariable(100001)
	if len(ruleIDs) != 1 {
		t.Errorf("删除 RULE_001 后，期望 1 个规则引用变量 100001，实际 %d", len(ruleIDs))
	}

	t.Log("变量索引测试通过")
}

// TestRuleManager 测试规则管理器
func TestRuleManager(t *testing.T) {
	storage := NewMockStorage()
	rm := rule.NewRuleManager(storage)

	// 创建测试规则
	testRule := &rule.AlarmRule{
		ID:       "TEST_RULE",
		Name:     "温度报警",
		Category: "TEMP",
		Level:    3,
		Enabled:  true,
		Condition: &rule.SingleCondition{
			VariableID: 100001,
			Operator:   rule.OpGTE,
			Value:      80.0,
		},
		TriggerMsg: rule.AlarmMessage{
			Type:    rule.ContentStatic,
			Content: "温度过高",
		},
	}

	// 测试1: 添加规则
	err := rm.AddRule(testRule)
	if err != nil {
		t.Fatalf("添加规则失败: %v", err)
	}

	// 测试2: 获取规则
	retrieved, exists := rm.GetRule("TEST_RULE")
	if !exists {
		t.Fatal("规则不存在")
	}
	if retrieved.Name != "温度报警" {
		t.Errorf("规则名称不匹配，期望 '温度报警'，实际 '%s'", retrieved.Name)
	}

	// 测试3: 获取统计
	stats := rm.GetStats()
	if stats.TotalRules != 1 {
		t.Errorf("期望 1 个规则，实际 %d", stats.TotalRules)
	}
	if stats.EnabledRules != 1 {
		t.Errorf("期望 1 个启用的规则，实际 %d", stats.EnabledRules)
	}

	// 测试4: 禁用规则
	err = rm.DisableRule("TEST_RULE")
	if err != nil {
		t.Fatalf("禁用规则失败: %v", err)
	}

	stats = rm.GetStats()
	if stats.EnabledRules != 0 {
		t.Errorf("禁用后期望 0 个启用的规则，实际 %d", stats.EnabledRules)
	}
	if stats.DisabledRules != 1 {
		t.Errorf("期望 1 个禁用的规则，实际 %d", stats.DisabledRules)
	}

	// 测试5: 删除规则
	err = rm.RemoveRule("TEST_RULE")
	if err != nil {
		t.Fatalf("删除规则失败: %v", err)
	}

	_, exists = rm.GetRule("TEST_RULE")
	if exists {
		t.Error("规则仍然存在，删除失败")
	}

	t.Log("规则管理器测试通过")
}

// TestVariableIndexWithComplexCondition 测试复杂条件的变量索引
func TestVariableIndexWithComplexCondition(t *testing.T) {
	index := rule.NewVariableRuleIndex()

	// 创建嵌套条件: A AND (B OR C)
	complexRule := &rule.AlarmRule{
		ID:       "COMPLEX_RULE",
		Name:     "复合条件报警",
		Category: "COMPLEX",
		Level:    4,
		Enabled:  true,
		Condition: &rule.ConditionGroup{
			Logic: rule.LogicAND,
			Conditions: []rule.Condition{
				&rule.SingleCondition{
					VariableID: 100001,
					Operator:   rule.OpGT,
					Value:      50.0,
				},
				&rule.ConditionGroup{
					Logic: rule.LogicOR,
					Conditions: []rule.Condition{
						&rule.SingleCondition{
							VariableID: 100002,
							Operator:   rule.OpGT,
							Value:      90.0,
						},
						&rule.SingleCondition{
							VariableID: 100003,
							Operator:   rule.OpLT,
							Value:      20.0,
						},
					},
				},
			},
		},
		TriggerMsg: rule.AlarmMessage{
			Type:    rule.ContentStatic,
			Content: "复合条件触发",
		},
	}

	// 添加到索引
	index.AddRule(complexRule)

	// 验证所有变量都被正确索引
	for _, varID := range []uint64{100001, 100002, 100003} {
		ruleIDs := index.GetRulesByVariable(varID)
		if len(ruleIDs) != 1 {
			t.Errorf("变量 %d 应该被 1 个规则引用，实际 %d", varID, len(ruleIDs))
		}
	}

	// 验证规则涉及所有变量
	vars := index.GetVariablesByRule("COMPLEX_RULE")
	if len(vars) != 3 {
		t.Errorf("规则应该涉及 3 个变量，实际 %d", len(vars))
	}

	t.Log("复杂条件变量索引测试通过")
}

// TestStateStats 测试状态统计
func TestStateStats(t *testing.T) {
	sm := engine.NewAlarmStateMachine(nil)

	// 添加多个规则并设置不同状态
	rules := []string{"RULE_001", "RULE_002", "RULE_003", "RULE_004"}

	sm.ForceTransition(rules[0], core.AlarmStateInactive)
	sm.ForceTransition(rules[1], core.AlarmStateActive)
	sm.ForceTransition(rules[2], core.AlarmStateAcknowledged)
	sm.ForceTransition(rules[3], core.AlarmStateCleared)

	// 验证状态查询
	if !sm.IsActive(rules[1]) {
		t.Error("RULE_001 应该是激活状态")
	}
	if !sm.IsAcknowledged(rules[2]) {
		t.Error("RULE_002 应该是已确认状态")
	}
	if !sm.IsCleared(rules[3]) {
		t.Error("RULE_003 应该是已清除状态")
	}

	// 获取转换历史
	for _, ruleID := range rules {
		history := sm.GetTransitionHistory(ruleID)
		if len(history) != 1 {
			t.Errorf("规则 %s 应该有 1 次转换，实际 %d", ruleID, len(history))
		}
	}

	t.Log("状态统计测试通过")
}

// MockStorage 模拟存储（用于测试）
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
	return &core.Variable{
		ID:      id,
		Value:   0.0,
		Quality: core.QualityGood,
	}, nil
}

func (m *MockStorage) WriteVar(variable *core.Variable) error {
	m.variables[variable.ID] = variable
	return nil
}

// 实现其他 Storage 接口方法（测试用简化版）
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
	return nil, nil
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

func (m *MockStorage) SubscribeByDevice(subscriberID string, deviceID string, callback func(core.VariableUpdate)) error {
	return nil
}

func (m *MockStorage) SubscribeByPattern(subscriberID string, pattern string, callback func(core.VariableUpdate)) error {
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
