package script

import (
	"sync"
	"testing"
	"time"
)

// MockScriptConsumerForScheduler 用于测试调度器的模拟消费者
type MockScriptConsumerForScheduler struct {
	mu             sync.Mutex
	executedScripts []string
	executionInputs []map[string]interface{}
}

func (m *MockScriptConsumerForScheduler) ExecuteScriptAsync(scriptID string, input map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.executedScripts = append(m.executedScripts, scriptID)
	m.executionInputs = append(m.executionInputs, input)

	return nil
}

func (m *MockScriptConsumerForScheduler) GetExecutedScripts() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]string{}, m.executedScripts...)
}

func (m *MockScriptConsumerForScheduler) GetExecutionInputs() []map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]map[string]interface{}, len(m.executionInputs))
	copy(result, m.executionInputs)
	return result
}

func (m *MockScriptConsumerForScheduler) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.executedScripts = nil
	m.executionInputs = nil
}

// 测试启动和停止调度器
func TestScriptScheduler_StartStop(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer) // 传入mock

	// 启动调度器
	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}

	if !scheduler.IsRunning() {
		t.Error("期望调度器正在运行，但IsRunning返回false")
	}

	// 重复启动应该返回错误
	err = scheduler.Start()
	if err == nil {
		t.Error("期望重复启动返回错误，但没有")
	}

	// 停止调度器
	err = scheduler.Stop()
	if err != nil {
		t.Fatalf("停止调度器失败: %v", err)
	}

	if scheduler.IsRunning() {
		t.Error("期望调度器已停止，但IsRunning返回true")
	}

	// 重复停止应该返回错误
	err = scheduler.Stop()
	if err == nil {
		t.Error("期望重复停止返回错误，但没有")
	}
}

// 测试添加周期触发器
func TestScriptScheduler_AddTrigger(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	// 启动调度器
	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 创建周期触发器
	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_001",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_001",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval: 100 * time.Millisecond, // 100ms间隔用于测试
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加周期触发器失败: %v", err)
	}

	// 验证触发器数量
	if scheduler.GetTriggerCount() != 1 {
		t.Errorf("期望触发器数量为1，实际为 %d", scheduler.GetTriggerCount())
	}

	// 等待触发
	time.Sleep(150 * time.Millisecond)

	// 验证脚本被执行
	executed := mockConsumer.GetExecutedScripts()
	if len(executed) < 1 {
		t.Errorf("期望至少执行1次，实际执行了 %d 次", len(executed))
	}
}

// 测试重复添加触发器
func TestScriptScheduler_AddDuplicateTrigger(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_002",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_002",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval: 1 * time.Second,
			},
		},
		Enabled: true,
	}

	// 第一次添加
	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加周期触发器失败: %v", err)
	}

	// 第二次添加（应该失败）
	err = scheduler.AddTrigger(trigger)
	if err == nil {
		t.Error("期望重复添加触发器返回错误，但没有")
	}
}

// 测试移除触发器
func TestScriptScheduler_RemoveTrigger(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_003",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_003",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval: 100 * time.Millisecond,
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加周期触发器失败: %v", err)
	}

	// 移除触发器
	err = scheduler.RemoveTrigger("SCHED_TRIGGER_003")
	if err != nil {
		t.Fatalf("移除触发器失败: %v", err)
	}

	// 验证触发器已删除
	if scheduler.GetTriggerCount() != 0 {
		t.Errorf("期望触发器数量为0，实际为 %d", scheduler.GetTriggerCount())
	}

	// 移除不存在的触发器应该返回错误
	err = scheduler.RemoveTrigger("NOT_EXISTS")
	if err == nil {
		t.Error("期望移除不存在的触发器返回错误，但没有")
	}
}

// 测试时间窗口 - 开始时间
func TestScriptScheduler_ShouldExecute_StartTime(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 获取当前时间
	now := time.Now()
	futureTime := now.Add(1 * time.Hour).Format("15:04:05")

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_004",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_004",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval:  100 * time.Millisecond,
				StartTime: futureTime, // 设置为未来时间
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加周期触发器失败: %v", err)
	}

	// 等待超过间隔时间
	time.Sleep(200 * time.Millisecond)

	// 验证脚本未执行（因为还没到开始时间）
	executed := mockConsumer.GetExecutedScripts()
	if len(executed) > 0 {
		t.Errorf("期望不执行（未到开始时间），但执行了 %d 次", len(executed))
	}
}

// 测试时间窗口 - 结束时间
func TestScriptScheduler_ShouldExecute_EndTime(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 获取当前时间
	now := time.Now()
	pastTime := now.Add(-1 * time.Hour).Format("15:04:05")

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_005",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_005",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval: 100 * time.Millisecond,
				EndTime:  pastTime, // 设置为过去时间
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加周期触发器失败: %v", err)
	}

	// 等待超过间隔时间
	time.Sleep(200 * time.Millisecond)

	// 验证脚本未执行（因为已过结束时间）
	executed := mockConsumer.GetExecutedScripts()
	if len(executed) > 0 {
		t.Errorf("期望不执行（已过结束时间），但执行了 %d 次", len(executed))
	}
}

// 测试星期几限制
func TestScriptScheduler_ShouldExecute_DaysOfWeek(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 获取当前星期几
	currentDay := int(time.Now().Weekday())
	if currentDay == 0 {
		currentDay = 7 // 周日转为7
	}

	// 设置为不同的星期几（应该不执行）
	differentDay := currentDay + 1
	if differentDay > 7 {
		differentDay = 1
	}

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_006",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_006",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval:   100 * time.Millisecond,
				DaysOfWeek: []int{differentDay}, // 设置为不同的星期几
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加周期触发器失败: %v", err)
	}

	// 等待超过间隔时间
	time.Sleep(200 * time.Millisecond)

	// 验证脚本未执行（因为不是指定的星期几）
	executed := mockConsumer.GetExecutedScripts()
	if len(executed) > 0 {
		t.Logf("当前星期几: %d, 配置星期几: %v", currentDay, []int{differentDay})
		t.Errorf("期望不执行（不是指定的星期几），但执行了 %d 次", len(executed))
	}
}

// 测试星期几 - 应该执行
func TestScriptScheduler_ShouldExecute_DaysOfWeek_ShouldExecute(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 获取当前星期几
	currentDay := int(time.Now().Weekday())
	if currentDay == 0 {
		currentDay = 7 // 周日转为7
	}

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_007",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_007",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval:   100 * time.Millisecond,
				DaysOfWeek: []int{currentDay}, // 设置为当前星期几
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加周期触发器失败: %v", err)
	}

	// 等待触发
	time.Sleep(150 * time.Millisecond)

	// 验证脚本被执行
	executed := mockConsumer.GetExecutedScripts()
	if len(executed) < 1 {
		t.Logf("当前星期几: %d, 配置星期几: %v", currentDay, []int{currentDay})
		t.Errorf("期望执行1次，但执行了 %d 次", len(executed))
	}
}

// 测试多个周期触发器
func TestScriptScheduler_MultipleTriggers(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	triggers := []*Trigger{
		{
			ID: "SCHED_TRIGGER_A", Type: TriggerTypePeriodic, ScriptID: "SCRIPT_A",
			Condition: TriggerCondition{
				PeriodicConfig: &PeriodicTriggerConfig{Interval: 100 * time.Millisecond},
			},
			Enabled: true,
		},
		{
			ID: "SCHED_TRIGGER_B", Type: TriggerTypePeriodic, ScriptID: "SCRIPT_B",
			Condition: TriggerCondition{
				PeriodicConfig: &PeriodicTriggerConfig{Interval: 150 * time.Millisecond},
			},
			Enabled: true,
		},
		{
			ID: "SCHED_TRIGGER_C", Type: TriggerTypePeriodic, ScriptID: "SCRIPT_C",
			Condition: TriggerCondition{
				PeriodicConfig: &PeriodicTriggerConfig{Interval: 200 * time.Millisecond},
			},
			Enabled: true,
		},
	}

	for _, trigger := range triggers {
		err = scheduler.AddTrigger(trigger)
		if err != nil {
			t.Fatalf("添加周期触发器失败: %v", err)
		}
	}

	// 等待所有触发器至少触发一次
	time.Sleep(250 * time.Millisecond)

	// 验证脚本被执行
	executed := mockConsumer.GetExecutedScripts()

	// 至少应该执行3次（每个触发器至少1次）
	if len(executed) < 3 {
		t.Errorf("期望至少执行3次，实际执行了 %d 次", len(executed))
	}

	// 验证三个不同的脚本都被执行
	scriptSet := make(map[string]bool)
	for _, scriptID := range executed {
		scriptSet[scriptID] = true
	}

	if len(scriptSet) != 3 {
		t.Errorf("期望3个不同的脚本被执行，实际只有 %d 个", len(scriptSet))
	}
}

// 测试列出触发器
func TestScriptScheduler_ListTriggers(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	triggers := []*Trigger{
		{
			ID: "SCHED_TRIGGER_A", Type: TriggerTypePeriodic, ScriptID: "SCRIPT_A",
			Condition: TriggerCondition{
				PeriodicConfig: &PeriodicTriggerConfig{Interval: 1 * time.Second},
			},
			Enabled: true,
		},
		{
			ID: "SCHED_TRIGGER_B", Type: TriggerTypePeriodic, ScriptID: "SCRIPT_B",
			Condition: TriggerCondition{
				PeriodicConfig: &PeriodicTriggerConfig{Interval: 1 * time.Second},
			},
			Enabled: true,
		},
	}

	for _, trigger := range triggers {
		err = scheduler.AddTrigger(trigger)
		if err != nil {
			t.Fatalf("添加周期触发器失败: %v", err)
		}
	}

	// 列出触发器
	list := scheduler.ListTriggers()

	if len(list) != 2 {
		t.Errorf("期望列出2个触发器，实际列出了 %d 个", len(list))
	}
}

// 测试更新触发器间隔
func TestScriptScheduler_UpdateTriggerInterval(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_008",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_008",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval: 200 * time.Millisecond, // 初始间隔200ms
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加周期触发器失败: %v", err)
	}

	// 等待第一次触发
	time.Sleep(250 * time.Millisecond)
	initialCount := len(mockConsumer.GetExecutedScripts())

	// 更新间隔为50ms
	err = scheduler.UpdateTriggerInterval("SCHED_TRIGGER_008", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("更新触发器间隔失败: %v", err)
	}

	// 等待更多触发
	time.Sleep(200 * time.Millisecond)
	newCount := len(mockConsumer.GetExecutedScripts())

	// 新的执行次数应该明显多于初始次数（因为间隔变短了）
	if newCount <= initialCount {
		t.Logf("初始执行次数: %d, 更新后执行次数: %d", initialCount, newCount)
		t.Error("期望更新间隔后执行次数增加，但没有")
	}
}

// 测试无效的触发器间隔
func TestScriptScheduler_UpdateTriggerInterval_Invalid(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 测试零间隔
	err = scheduler.UpdateTriggerInterval("TRIGGER_NOT_EXISTS", 0)
	if err == nil {
		t.Error("期望零间隔返回错误，但没有")
	}

	// 测试负间隔
	err = scheduler.UpdateTriggerInterval("TRIGGER_NOT_EXISTS", -1*time.Second)
	if err == nil {
		t.Error("期望负间隔返回错误，但没有")
	}
}

// 测试触发器类型验证
func TestScriptScheduler_AddTrigger_InvalidType(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 尝试添加非周期触发器
	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_009",
		Type:     TriggerTypeVariable, // 错误的类型
		ScriptID: "SCRIPT_009",
		Condition: TriggerCondition{
			VariableID: 100001,
			Operator:   ">=",
			Threshold:  80.0,
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err == nil {
		t.Error("期望添加非周期触发器返回错误，但没有")
	}
}

// 测试缺少PeriodicConfig
func TestScriptScheduler_AddTrigger_NoPeriodicConfig(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_010",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_010",
		Condition: TriggerCondition{
			// 缺少PeriodicConfig
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err == nil {
		t.Error("期望缺少PeriodicConfig返回错误，但没有")
	}
}

// 测试在调度器未启动时添加触发器
func TestScriptScheduler_AddTrigger_NotRunning(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	// 不启动调度器

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_011",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_011",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval: 1 * time.Second,
			},
		},
		Enabled: true,
	}

	err := scheduler.AddTrigger(trigger)
	if err == nil {
		t.Error("期望在调度器未启动时添加触发器返回错误，但没有")
	}
}

// 测试零间隔
func TestScriptScheduler_AddTrigger_ZeroInterval(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	trigger := &Trigger{
		ID:       "SCHED_TRIGGER_012",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_012",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval: 0, // 零间隔
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err == nil {
		t.Error("期望零间隔返回错误，但没有")
	}
}

// ============ Phase 2 Enhanced: Cron 表达式测试 ============

// 测试Cron表达式触发器添加
func TestScriptScheduler_AddCronTrigger(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 创建每秒执行的Cron触发器
	trigger := &Trigger{
		ID:       "CRON_TEST_001",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_001",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				CronExpr: "* * * * * *", // 每秒
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加Cron触发器失败: %v", err)
	}

	// 验证触发器数量
	if scheduler.GetTriggerCount() != 1 {
		t.Errorf("期望触发器数量为1，实际为 %d", scheduler.GetTriggerCount())
	}

	// 等待触发
	time.Sleep(1500 * time.Millisecond)

	// 验证脚本被执行（应该至少执行1次）
	executed := mockConsumer.GetExecutedScripts()
	if len(executed) < 1 {
		t.Errorf("期望至少执行1次，实际执行了 %d 次", len(executed))
	}
}

// 测试Cron表达式执行 - 每5秒
func TestScriptScheduler_CronExecution(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 每5秒执行一次
	trigger := &Trigger{
		ID:       "CRON_TEST_002",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_002",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				CronExpr: "0/5 * * * * *", // 每5秒
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加Cron触发器失败: %v", err)
	}

	// 等待12秒（应该触发2次）
	time.Sleep(12 * time.Second)

	executed := mockConsumer.GetExecutedScripts()

	// 验证执行次数（应该约为2次，允许±1的误差）
	if len(executed) < 1 || len(executed) > 3 {
		t.Errorf("期望执行约2次，实际执行了 %d 次", len(executed))
	}
}

// 测试Cron表达式 + 时间窗口
func TestScriptScheduler_CronWithTimeWindow(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 获取当前时间
	now := time.Now()
	pastTime := now.Add(-1 * time.Hour).Format("15:04:05")

	trigger := &Trigger{
		ID:       "CRON_TEST_003",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_003",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				CronExpr:  "* * * * * *", // 每秒
				EndTime:   pastTime,     // 结束时间设为过去
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加Cron触发器失败: %v", err)
	}

	// 等待2秒
	time.Sleep(2 * time.Second)

	// 验证脚本未执行（因为已过结束时间）
	executed := mockConsumer.GetExecutedScripts()
	if len(executed) > 0 {
		t.Errorf("期望不执行（已过结束时间），但执行了 %d 次", len(executed))
	}
}

// 测试Cron触发器移除
func TestScriptScheduler_RemoveCronTrigger(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	trigger := &Trigger{
		ID:       "CRON_TEST_004",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_004",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				CronExpr: "* * * * * *",
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err != nil {
		t.Fatalf("添加Cron触发器失败: %v", err)
	}

	// 验证已添加
	if scheduler.GetTriggerCount() != 1 {
		t.Errorf("期望触发器数量为1，实际为 %d", scheduler.GetTriggerCount())
	}

	// 移除触发器
	err = scheduler.RemoveTrigger("CRON_TEST_004")
	if err != nil {
		t.Fatalf("移除Cron触发器失败: %v", err)
	}

	// 验证已删除
	if scheduler.GetTriggerCount() != 0 {
		t.Errorf("期望触发器数量为0，实际为 %d", scheduler.GetTriggerCount())
	}
}

// 测试无效Cron表达式
func TestScriptScheduler_InvalidCronExpr(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	trigger := &Trigger{
		ID:       "CRON_TEST_005",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_005",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				CronExpr: "invalid cron expression", // 无效的表达式
			},
		},
		Enabled: true,
	}

	err = scheduler.AddTrigger(trigger)
	if err == nil {
		t.Error("期望无效Cron表达式返回错误，但没有")
	}
}

// 测试Cron和Interval混合调度
func TestScriptScheduler_HybridScheduling(t *testing.T) {
	mockConsumer := &MockScriptConsumerForScheduler{}
	scheduler := NewScriptScheduler(mockConsumer)

	err := scheduler.Start()
	if err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer scheduler.Stop()

	// 添加Interval触发器
	intervalTrigger := &Trigger{
		ID:       "INTERVAL_TEST",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_INTERVAL",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				Interval: 200 * time.Millisecond, // 每200ms
			},
		},
		Enabled: true,
	}

	// 添加Cron触发器
	cronTrigger := &Trigger{
		ID:       "CRON_TEST",
		Type:     TriggerTypePeriodic,
		ScriptID: "SCRIPT_CRON",
		Condition: TriggerCondition{
			PeriodicConfig: &PeriodicTriggerConfig{
				CronExpr: "* * * * * *", // 每秒
			},
		},
		Enabled: true,
	}

	// 添加两个触发器
	err = scheduler.AddTrigger(intervalTrigger)
	if err != nil {
		t.Fatalf("添加Interval触发器失败: %v", err)
	}

	err = scheduler.AddTrigger(cronTrigger)
	if err != nil {
		t.Fatalf("添加Cron触发器失败: %v", err)
	}

	// 验证触发器数量
	if scheduler.GetTriggerCount() != 2 {
		t.Errorf("期望触发器数量为2，实际为 %d", scheduler.GetTriggerCount())
	}

	// 等待1秒
	time.Sleep(1 * time.Second)

	// 验证两个脚本都被执行
	executed := mockConsumer.GetExecutedScripts()

	// 统计不同脚本的执行次数
	intervalCount := 0
	cronCount := 0
	for _, scriptID := range executed {
		if scriptID == "SCRIPT_INTERVAL" {
			intervalCount++
		} else if scriptID == "SCRIPT_CRON" {
			cronCount++
		}
	}

	// 验证两个脚本都执行了至少一次
	if intervalCount < 1 {
		t.Errorf("期望Interval脚本至少执行1次，实际执行了 %d 次", intervalCount)
	}

	if cronCount < 1 {
		t.Errorf("期望Cron脚本至少执行1次，实际执行了 %d 次", cronCount)
	}

	t.Logf("Interval脚本执行: %d次, Cron脚本执行: %d次", intervalCount, cronCount)
}
