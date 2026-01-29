package script

import (
	"sync"
	"testing"
	"time"

	"pansiot-device/internal/core"
)

// MockScriptConsumer 用于测试的模拟消费者
type MockScriptConsumer struct {
	mu             sync.Mutex
	executedScripts []string
	executionInputs []map[string]interface{}
}

func (m *MockScriptConsumer) ExecuteScriptAsync(scriptID string, input map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.executedScripts = append(m.executedScripts, scriptID)
	m.executionInputs = append(m.executionInputs, input)

	return nil
}

func (m *MockScriptConsumer) GetExecutedScripts() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]string{}, m.executedScripts...)
}

func (m *MockScriptConsumer) GetExecutionInputs() []map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]map[string]interface{}, len(m.executionInputs))
	copy(result, m.executionInputs)
	return result
}

func (m *MockScriptConsumer) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.executedScripts = nil
	m.executionInputs = nil
}

// 测试注册触发器
func TestTriggerManager_RegisterTrigger(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer) // 传入mock

	trigger := &Trigger{
		ID:       "TRIGGER_001",
		Type:     TriggerTypeVariable,
		ScriptID: "SCRIPT_001",
		Condition: TriggerCondition{
			VariableID: 100001,
			Operator:   ">=",
			Threshold:  80.0,
		},
		Enabled: true,
	}

	err := tm.RegisterTrigger(trigger)
	if err != nil {
		t.Fatalf("注册触发器失败: %v", err)
	}

	// 验证触发器已注册
	info, err := tm.GetTriggerInfo("TRIGGER_001")
	if err != nil {
		t.Fatalf("获取触发器信息失败: %v", err)
	}

	if info.ID != "TRIGGER_001" {
		t.Errorf("触发器ID不匹配: 期望 TRIGGER_001, 实际 %s", info.ID)
	}

	// 测试重复注册
	err = tm.RegisterTrigger(trigger)
	if err == nil {
		t.Error("期望重复注册返回错误，但没有")
	}
}

// 测试注销触发器
func TestTriggerManager_UnregisterTrigger(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer)

	trigger := &Trigger{
		ID:       "TRIGGER_002",
		Type:     TriggerTypeVariable,
		ScriptID: "SCRIPT_002",
		Condition: TriggerCondition{
			VariableID: 100002,
			Operator:   "<",
			Threshold:  10.0,
		},
		Enabled: true,
	}

	// 注册触发器
	err := tm.RegisterTrigger(trigger)
	if err != nil {
		t.Fatalf("注册触发器失败: %v", err)
	}

	// 注销触发器
	err = tm.UnregisterTrigger("TRIGGER_002")
	if err != nil {
		t.Fatalf("注销触发器失败: %v", err)
	}

	// 验证触发器已删除
	_, err = tm.GetTriggerInfo("TRIGGER_002")
	if err == nil {
		t.Error("期望触发器不存在，但找到了")
	}

	// 测试注销不存在的触发器
	err = tm.UnregisterTrigger("TRIGGER_NOT_EXISTS")
	if err == nil {
		t.Error("期望注销不存在的触发器返回错误，但没有")
	}
}

// 测试条件评估 - 所有操作符
func TestTriggerManager_EvaluateCondition(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer)

	tests := []struct {
		name      string
		value     interface{}
		threshold interface{}
		operator  string
		expected  bool
	}{
		{"大于 - 通过", 85.0, 80.0, ">", true},
		{"大于 - 不通过", 75.0, 80.0, ">", false},
		{"小于 - 通过", 75.0, 80.0, "<", true},
		{"小于 - 不通过", 85.0, 80.0, "<", false},
		{"等于 - 通过", 80.0, 80.0, "==", true},
		{"等于 - 不通过", 81.0, 80.0, "==", false},
		{"不等于 - 通过", 81.0, 80.0, "!=", true},
		{"不等于 - 不通过", 80.0, 80.0, "!=", false},
		{"大于等于 - 通过", 80.0, 80.0, ">=", true},
		{"大于等于 - 通过2", 85.0, 80.0, ">=", true},
		{"大于等于 - 不通过", 75.0, 80.0, ">=", false},
		{"小于等于 - 通过", 80.0, 80.0, "<=", true},
		{"小于等于 - 通过2", 75.0, 80.0, "<=", true},
		{"小于等于 - 不通过", 85.0, 80.0, "<=", false},
		{"整数比较", 100, 80, ">", true},
		{"布尔值比较", true, false, "==", false},
		{"布尔值比较2", true, true, "==", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := TriggerCondition{
				VariableID: 100001,
				Operator:   tt.operator,
				Threshold:  tt.threshold,
			}

			result := tm.evaluateCondition(condition, tt.value)
			if result != tt.expected {
				t.Errorf("条件评估错误: 值=%v, 阈值=%v, 操作=%s, 期望=%v, 实际=%v",
					tt.value, tt.threshold, tt.operator, tt.expected, result)
			}
		})
	}
}

// 测试变量变化回调
func TestTriggerManager_OnVariableChanged(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer)

	// 注册一个触发器
	trigger := &Trigger{
		ID:       "TRIGGER_003",
		Type:     TriggerTypeVariable,
		ScriptID: "SCRIPT_003",
		Condition: TriggerCondition{
			VariableID: 100003,
			Operator:   ">=",
			Threshold:  80.0,
		},
		Enabled: true,
	}

	err := tm.RegisterTrigger(trigger)
	if err != nil {
		t.Fatalf("注册触发器失败: %v", err)
	}

	// 模拟变量更新
	update := core.VariableUpdate{
		VariableID: 100003,
		Value:      85.0,
		Timestamp:  time.Now(),
	}

	// 调用回调
	tm.onVariableChanged(update)

	// 等待异步执行
	time.Sleep(100 * time.Millisecond)

	// 验证脚本被执行
	executed := mockConsumer.GetExecutedScripts()
	if len(executed) != 1 {
		t.Fatalf("期望执行1个脚本，实际执行了 %d 个", len(executed))
	}

	if executed[0] != "SCRIPT_003" {
		t.Errorf("期望执行脚本 SCRIPT_003，实际执行了 %s", executed[0])
	}

	// 验证输入数据
	inputs := mockConsumer.GetExecutionInputs()
	if len(inputs) != 1 {
		t.Fatalf("期望有1个输入，实际有 %d 个", len(inputs))
	}

	if inputs[0]["trigger_type"] != "variable" {
		t.Errorf("期望触发类型为 variable，实际为 %v", inputs[0]["trigger_type"])
	}

	if inputs[0]["value"] != 85.0 {
		t.Errorf("期望值为 85.0，实际为 %v", inputs[0]["value"])
	}
}

// 测试多个触发器监听同一变量
func TestTriggerManager_MultipleTriggersSameVariable(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer)

	// 注册多个触发器监听同一变量
	triggers := []*Trigger{
		{
			ID:       "TRIGGER_A", Type: TriggerTypeVariable, ScriptID: "SCRIPT_A",
			Condition: TriggerCondition{VariableID: 100004, Operator: ">=", Threshold: 80.0},
			Enabled: true,
		},
		{
			ID:       "TRIGGER_B", Type: TriggerTypeVariable, ScriptID: "SCRIPT_B",
			Condition: TriggerCondition{VariableID: 100004, Operator: ">=", Threshold: 90.0},
			Enabled: true,
		},
		{
			ID:       "TRIGGER_C", Type: TriggerTypeVariable, ScriptID: "SCRIPT_C",
			Condition: TriggerCondition{VariableID: 100004, Operator: "<", Threshold: 50.0},
			Enabled: true,
		},
	}

	for _, trigger := range triggers {
		err := tm.RegisterTrigger(trigger)
		if err != nil {
			t.Fatalf("注册触发器失败: %v", err)
		}
	}

	// 模拟变量更新到85
	update := core.VariableUpdate{
		VariableID: 100004,
		Value:      85.0,
		Timestamp:  time.Now(),
	}

	tm.onVariableChanged(update)

	// 等待异步执行
	time.Sleep(100 * time.Millisecond)

	// 验证：应该只触发 TRIGGER_A（85 >= 80 为true）
	// TRIGGER_B (85 >= 90) 和 TRIGGER_C (85 < 50) 都不应该触发
	executed := mockConsumer.GetExecutedScripts()

	if len(executed) != 1 {
		t.Logf("执行的脚本: %v", executed)
		t.Fatalf("期望执行1个脚本，实际执行了 %d 个", len(executed))
	}
}

// 测试禁用的触发器
func TestTriggerManager_DisabledTrigger(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer)

	trigger := &Trigger{
		ID:       "TRIGGER_005",
		Type:     TriggerTypeVariable,
		ScriptID: "SCRIPT_005",
		Condition: TriggerCondition{
			VariableID: 100005,
			Operator:   ">=",
			Threshold:  80.0,
		},
		Enabled: false, // 禁用
	}

	err := tm.RegisterTrigger(trigger)
	if err != nil {
		t.Fatalf("注册触发器失败: %v", err)
	}

	update := core.VariableUpdate{
		VariableID: 100005,
		Value:      85.0,
		Timestamp:  time.Now(),
	}

	tm.onVariableChanged(update)

	// 等待异步执行
	time.Sleep(100 * time.Millisecond)

	// 验证：不应该执行
	executed := mockConsumer.GetExecutedScripts()
	if len(executed) != 0 {
		t.Errorf("期望禁用的触发器不执行，但执行了 %d 个脚本", len(executed))
	}
}

// 测试启用/禁用触发器
func TestTriggerManager_EnableDisableTrigger(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer)

	trigger := &Trigger{
		ID:       "TRIGGER_006",
		Type:     TriggerTypeVariable,
		ScriptID: "SCRIPT_006",
		Condition: TriggerCondition{
			VariableID: 100006,
			Operator:   ">=",
			Threshold:  80.0,
		},
		Enabled: true,
	}

	err := tm.RegisterTrigger(trigger)
	if err != nil {
		t.Fatalf("注册触发器失败: %v", err)
	}

	// 禁用触发器
	err = tm.DisableTrigger("TRIGGER_006")
	if err != nil {
		t.Fatalf("禁用触发器失败: %v", err)
	}

	info, _ := tm.GetTriggerInfo("TRIGGER_006")
	if info.Enabled {
		t.Error("期望触发器被禁用，但仍然是启用状态")
	}

	// 启用触发器
	err = tm.EnableTrigger("TRIGGER_006")
	if err != nil {
		t.Fatalf("启用触发器失败: %v", err)
	}

	info, _ = tm.GetTriggerInfo("TRIGGER_006")
	if !info.Enabled {
		t.Error("期望触发器被启用，但仍然是禁用状态")
	}
}

// 测试获取触发器统计
func TestTriggerManager_GetTriggerStats(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer)

	trigger := &Trigger{
		ID:       "TRIGGER_007",
		Type:     TriggerTypeVariable,
		ScriptID: "SCRIPT_007",
		Condition: TriggerCondition{
			VariableID: 100007,
			Operator:   ">=",
			Threshold:  80.0,
		},
		Enabled: true,
	}

	err := tm.RegisterTrigger(trigger)
	if err != nil {
		t.Fatalf("注册触发器失败: %v", err)
	}

	// 获取初始统计
	lastTriggered, count, err := tm.GetTriggerStats("TRIGGER_007")
	if err != nil {
		t.Fatalf("获取触发器统计失败: %v", err)
	}

	if count != 0 {
		t.Errorf("期望初始触发次数为0，实际为 %d", count)
	}

	if !lastTriggered.IsZero() {
		t.Error("期望初始触发时间为零，但有值")
	}

	// 模拟触发
	update := core.VariableUpdate{
		VariableID: 100007,
		Value:      85.0,
		Timestamp:  time.Now(),
	}

	tm.onVariableChanged(update)
	time.Sleep(100 * time.Millisecond)

	// 获取更新后的统计
	lastTriggered, count, err = tm.GetTriggerStats("TRIGGER_007")
	if err != nil {
		t.Fatalf("获取触发器统计失败: %v", err)
	}

	if count != 1 {
		t.Errorf("期望触发次数为1，实际为 %d", count)
	}

	if lastTriggered.IsZero() {
		t.Error("期望触发时间有值，但是零")
	}
}

// 测试列出所有触发器
func TestTriggerManager_ListAllTriggers(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer)

	triggers := []*Trigger{
		{
			ID: "TRIGGER_A", Type: TriggerTypeVariable, ScriptID: "SCRIPT_A",
			Condition: TriggerCondition{VariableID: 100001, Operator: ">=", Threshold: 80.0},
			Enabled: true,
		},
		{
			ID: "TRIGGER_B", Type: TriggerTypeVariable, ScriptID: "SCRIPT_B",
			Condition: TriggerCondition{VariableID: 100002, Operator: "<", Threshold: 50.0},
			Enabled: true,
		},
		{
			ID: "TRIGGER_C", Type: TriggerTypePeriodic, ScriptID: "SCRIPT_C",
			Condition: TriggerCondition{
				PeriodicConfig: &PeriodicTriggerConfig{Interval: 1 * time.Minute},
			},
			Enabled: true,
		},
	}

	for _, trigger := range triggers {
		err := tm.RegisterTrigger(trigger)
		if err != nil {
			t.Fatalf("注册触发器失败: %v", err)
		}
	}

	// 列出所有触发器
	allTriggers := tm.ListAllTriggers()

	if len(allTriggers) != 3 {
		t.Errorf("期望列出3个触发器，实际列出了 %d 个", len(allTriggers))
	}
}

// 测试获取脚本的触发器
func TestTriggerManager_GetScriptTriggers(t *testing.T) {
	mockConsumer := &MockScriptConsumer{}
	tm := NewTriggerManager(mockConsumer)

	// 为同一脚本注册多个触发器
	triggers := []*Trigger{
		{ID: "TRIGGER_A", Type: TriggerTypeVariable, ScriptID: "SCRIPT_X",
			Condition: TriggerCondition{VariableID: 100001, Operator: ">=", Threshold: 80.0}},
		{ID: "TRIGGER_B", Type: TriggerTypeVariable, ScriptID: "SCRIPT_X",
			Condition: TriggerCondition{VariableID: 100002, Operator: "<", Threshold: 50.0}},
		{ID: "TRIGGER_C", Type: TriggerTypeVariable, ScriptID: "SCRIPT_Y",
			Condition: TriggerCondition{VariableID: 100003, Operator: "==", Threshold: 100.0}},
	}

	for _, trigger := range triggers {
		err := tm.RegisterTrigger(trigger)
		if err != nil {
			t.Fatalf("注册触发器失败: %v", err)
		}
	}

	// 获取 SCRIPT_X 的触发器
	scriptXTriggers := tm.GetScriptTriggers("SCRIPT_X")

	if len(scriptXTriggers) != 2 {
		t.Errorf("期望SCRIPT_X有2个触发器，实际有 %d 个", len(scriptXTriggers))
	}

	// 获取 SCRIPT_Y 的触发器
	scriptYTriggers := tm.GetScriptTriggers("SCRIPT_Y")

	if len(scriptYTriggers) != 1 {
		t.Errorf("期望SCRIPT_Y有1个触发器，实际有 %d 个", len(scriptYTriggers))
	}
}
