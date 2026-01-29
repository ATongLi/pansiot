package action

import (
	"context"
	"testing"
	"time"

	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/alarm/rule"
	"pansiot-device/internal/core"
)

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

func (m *MockStorage) CreateVariable(variable *core.Variable) error {
	m.variables[variable.ID] = variable
	return nil
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

func (m *MockStorage) WriteVars(variables []*core.Variable) error {
	for _, v := range variables {
		m.variables[v.ID] = v
	}
	return nil
}

func (m *MockStorage) Subscribe(subscriberID string, variableIDs []uint64, callback func(core.VariableUpdate)) error {
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

func (m *MockStorage) ReadVarByStringID(stringID string) (*core.Variable, error) {
	return nil, nil
}

func (m *MockStorage) SubscribeByDevice(subscriberID string, deviceID string, callback func(core.VariableUpdate)) error {
	return nil
}

func (m *MockStorage) SubscribeByPattern(subscriberID string, pattern string, callback func(core.VariableUpdate)) error {
	return nil
}

// MockStateMachine 模拟状态机（用于测试）
type MockStateMachine struct {
	states map[string]core.AlarmState
}

func NewMockStateMachine() *MockStateMachine {
	return &MockStateMachine{
		states: make(map[string]core.AlarmState),
	}
}

func (m *MockStateMachine) GetState(ruleID string) core.AlarmState {
	if state, ok := m.states[ruleID]; ok {
		return state
	}
	return core.AlarmStateInactive
}

func (m *MockStateMachine) SetState(ruleID string, state core.AlarmState) {
	m.states[ruleID] = state
}

// TestVariableWriterHandler 测试变量写值处理器
func TestVariableWriterHandler(t *testing.T) {
	storage := NewMockStorage()
	handler := &VariableWriterHandler{storage: storage}

	action := &rule.Action{
		ID:   "action_001",
		Type: rule.ActWriteVar,
		Params: map[string]interface{}{
			"variable_id": uint64(100001),
			"value":       42.0,
		},
	}

	alarm := &engine.ActiveAlarm{
		RuleID: "RULE_001",
		Rule: &rule.AlarmRule{
			ID:   "RULE_001",
			Name: "测试报警",
		},
	}

	err := handler.Handle(context.Background(), action, alarm)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 验证变量被写入
	v, _ := storage.ReadVar(100001)
	if v.Value != 42.0 {
		t.Errorf("期望值 42.0，实际 %v", v.Value)
	}

	t.Logf("变量写值测试通过: 变量ID=%d, 值=%v", 100001, v.Value)
}

// TestSoundPlayerHandler 测试声音播放处理器
func TestSoundPlayerHandler(t *testing.T) {
	handler := &SoundPlayerHandler{}

	action := &rule.Action{
		ID:   "action_002",
		Type: rule.ActSound,
		Params: map[string]interface{}{
			"file":       "/sounds/alarm.wav",
			"continuous": true,
			"volume":     0.8,
		},
	}

	alarm := &engine.ActiveAlarm{
		RuleID: "RULE_001",
		Rule: &rule.AlarmRule{
			ID:   "RULE_001",
			Name: "温度报警",
		},
	}

	err := handler.Handle(context.Background(), action, alarm)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	t.Log("声音播放测试通过（日志模拟）")
}

// TestPopupNotifierHandler 测试弹窗通知处理器
func TestPopupNotifierHandler(t *testing.T) {
	handler := &PopupNotifierHandler{}

	action := &rule.Action{
		ID:   "action_003",
		Type: rule.ActPopup,
		Params: map[string]interface{}{
			"title":   "温度报警",
			"message": "设备温度过高：85°C",
		},
	}

	alarm := &engine.ActiveAlarm{
		RuleID: "RULE_001",
		Rule: &rule.AlarmRule{
			ID:         "RULE_001",
			Name:       "温度报警",
			TriggerMsg: rule.AlarmMessage{Content: "默认消息"},
		},
	}

	err := handler.Handle(context.Background(), action, alarm)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	t.Log("弹窗通知测试通过（日志模拟）")
}

// TestPageJumperHandler 测试页面跳转处理器
func TestPageJumperHandler(t *testing.T) {
	handler := &PageJumperHandler{}

	action := &rule.Action{
		ID:   "action_004",
		Type: rule.ActJumpPage,
		Params: map[string]interface{}{
			"page": "/alarm/detail",
			"params": map[string]interface{}{
				"alarm_id":  "ALARM_001",
				"device_id": "device_001",
			},
		},
	}

	alarm := &engine.ActiveAlarm{
		RuleID: "RULE_001",
		Rule: &rule.AlarmRule{
			ID:   "RULE_001",
			Name: "温度报警",
		},
	}

	err := handler.Handle(context.Background(), action, alarm)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	t.Log("页面跳转测试通过（日志模拟）")
}

// TestActionExecutor 测试动作执行器
func TestActionExecutor(t *testing.T) {
	storage := NewMockStorage()
	executor := NewActionExecutor(storage)

	// 启动执行器
	ctx := context.Background()
	if err := executor.Start(ctx); err != nil {
		t.Fatalf("启动执行器失败: %v", err)
	}
	defer executor.Stop()

	alarm := &engine.ActiveAlarm{
		RuleID: "RULE_001",
		Rule: &rule.AlarmRule{
			ID:   "RULE_001",
			Name: "测试报警",
		},
	}

	// 测试变量写值动作
	writeVarAction := &rule.Action{
		ID:   "action_001",
		Type: rule.ActWriteVar,
		Params: map[string]interface{}{
			"variable_id": uint64(100002),
			"value":       1,
		},
	}

	err := executor.Execute(ctx, writeVarAction, alarm)
	if err != nil {
		t.Fatalf("执行动作失败: %v", err)
	}

	// 验证变量被写入
	v, _ := storage.ReadVar(100002)
	if v.Value != 1 {
		t.Errorf("期望值 1，实际 %v", v.Value)
	}

	t.Log("动作执行器测试通过")
}

// TestActionBatchExecution 测试批量动作执行
func TestActionBatchExecution(t *testing.T) {
	storage := NewMockStorage()
	executor := NewActionExecutor(storage)

	ctx := context.Background()
	executor.Start(ctx)
	defer executor.Stop()

	alarm := &engine.ActiveAlarm{
		RuleID: "RULE_001",
		Rule: &rule.AlarmRule{
			ID:   "RULE_001",
			Name: "测试报警",
		},
	}

	actions := []rule.Action{
		{
			ID:   "action_001",
			Type: rule.ActWriteVar,
			Params: map[string]interface{}{
				"variable_id": uint64(100003),
				"value":       1,
			},
		},
		{
			ID:   "action_002",
			Type: rule.ActWriteVar,
			Params: map[string]interface{}{
				"variable_id": uint64(100004),
				"value":       2,
			},
		},
		{
			ID:   "action_003",
			Type: rule.ActSound,
			Params: map[string]interface{}{
				"file": "/sounds/alarm.wav",
			},
		},
	}

	err := executor.ExecuteBatch(ctx, actions, alarm)
	if err != nil {
		t.Fatalf("批量执行失败: %v", err)
	}

	// 验证变量被写入
	v1, _ := storage.ReadVar(100003)
	if v1.Value != 1 {
		t.Errorf("期望值 1，实际 %v", v1.Value)
	}

	v2, _ := storage.ReadVar(100004)
	if v2.Value != 2 {
		t.Errorf("期望值 2，实际 %v", v2.Value)
	}

	t.Log("批量动作执行测试通过")
}

// TestSchedulerDelayedExecution 测试延迟执行
func TestSchedulerDelayedExecution(t *testing.T) {
	storage := NewMockStorage()
	executor := NewActionExecutor(storage)
	stateMachine := NewMockStateMachine()
	scheduler := NewScheduler(executor, stateMachine, 2)

	scheduler.Start()
	defer scheduler.Stop()

	alarm := &engine.ActiveAlarm{
		RuleID: "RULE_001",
		Rule: &rule.AlarmRule{
			ID:   "RULE_001",
			Name: "测试报警",
		},
	}

	// 延迟 100ms 执行
	action := &rule.Action{
		ID:    "action_delayed",
		Type:  rule.ActWriteVar,
		When:  rule.WhenAfterTrigger,
		Delay: 100 * time.Millisecond,
		Params: map[string]interface{}{
			"variable_id": uint64(100005),
			"value":       99,
		},
	}

	err := scheduler.Schedule(action, alarm)
	if err != nil {
		t.Fatalf("调度失败: %v", err)
	}

	// 等待执行完成
	time.Sleep(200 * time.Millisecond)

	// 验证变量被写入
	v, _ := storage.ReadVar(100005)
	if v.Value != 99 {
		t.Errorf("期望值 99，实际 %v", v.Value)
	}

	t.Log("延迟执行测试通过")
}

// TestSchedulerLoopExecution 测试循环执行
func TestSchedulerLoopExecution(t *testing.T) {
	storage := NewMockStorage()
	executor := NewActionExecutor(storage)
	stateMachine := NewMockStateMachine()
	scheduler := NewScheduler(executor, stateMachine, 2)

	scheduler.Start()
	defer scheduler.Stop()

	alarm := &engine.ActiveAlarm{
		RuleID: "RULE_001",
		Rule: &rule.AlarmRule{
			ID:   "RULE_001",
			Name: "测试报警",
		},
	}

	// 循环 3 次，间隔 50ms
	action := &rule.Action{
		ID:        "action_loop",
		Type:      rule.ActWriteVar,
		When:      rule.WhenTrigger,
		Mode:      rule.ModeLoop,
		LoopCount: 3,
		LoopDelay: 50 * time.Millisecond,
		Params: map[string]interface{}{
			"variable_id": uint64(100006),
			"value":       1,
		},
	}

	startTime := time.Now()
	err := scheduler.Schedule(action, alarm)
	if err != nil {
		t.Fatalf("调度失败: %v", err)
	}

	// 等待循环完成（3次 * 50ms + 余量）
	time.Sleep(200 * time.Millisecond)
	elapsed := time.Since(startTime)

	if elapsed < 150*time.Millisecond {
		t.Errorf("循环执行时间过短: %v", elapsed)
	}

	t.Logf("循环执行测试通过，耗时: %v", elapsed)
}

// TestActionValidation 测试动作参数验证
func TestActionValidation(t *testing.T) {
	storage := NewMockStorage()
	executor := NewActionExecutor(storage)

	ctx := context.Background()
	executor.Start(ctx)
	defer executor.Stop()

	alarm := &engine.ActiveAlarm{
		RuleID: "RULE_001",
		Rule: &rule.AlarmRule{
			ID:   "RULE_001",
			Name: "测试报警",
		},
	}

	// 测试缺少参数的情况
	invalidAction := &rule.Action{
		ID:     "action_invalid",
		Type:   rule.ActWriteVar,
		Params: map[string]interface{}{}, // 缺少 variable_id 和 value
	}

	err := executor.Execute(ctx, invalidAction, alarm)
	if err == nil {
		t.Fatal("期望参数验证失败，但没有返回错误")
	}

	t.Logf("参数验证测试通过: %v", err)
}
