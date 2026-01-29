package collector

import (
	"context"
	"fmt"
	"testing"
	"time"

	"pansiot-device/internal/adapter"
	"pansiot-device/internal/core"
	"pansiot-device/internal/storage"
)

// TestNewCollector 测试采集器创建
func TestNewCollector(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()

	collector := NewCollector(factory, storage)

	if collector == nil {
		t.Fatal("采集器创建失败")
	}

	if collector.tasks == nil {
		t.Error("tasks map 未初始化")
	}

	if collector.taskRunners == nil {
		t.Error("taskRunners map 未初始化")
	}

	if collector.running.Load() {
		t.Error("新创建的采集器不应处于运行状态")
	}
}

// TestCollectorStartStop 测试采集器启动和停止
func TestCollectorStartStop(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	ctx := context.Background()

	// 测试启动
	err := collector.Start(ctx)
	if err != nil {
		t.Fatalf("启动采集器失败: %v", err)
	}

	if !collector.running.Load() {
		t.Error("采集器启动后 running 应为 true")
	}

	// 测试重复启动
	err = collector.Start(ctx)
	if err == nil {
		t.Error("重复启动应该返回错误")
	}

	// 测试停止
	err = collector.Stop()
	if err != nil {
		t.Fatalf("停止采集器失败: %v", err)
	}

	// 测试重复停止
	err = collector.Stop()
	if err == nil {
		t.Error("重复停止应该返回错误")
	}
}

// TestValidateTask 测试任务验证
func TestValidateTask(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	tests := []struct {
		name    string
		task    *core.CollectionTask
		wantErr bool
	}{
		{
			name: "有效任务",
			task: &core.CollectionTask{
				ID:           "TASK_001",
				Frequency:    1000,
				DeviceID:     "DEVICE_001",
				ProtocolType: "mock",
				VariableIDs:  []uint64{1, 2, 3},
				Priority:     5,
				Timeout:      5000,
				Enable:       true,
			},
			wantErr: false,
		},
		{
			name: "空任务ID",
			task: &core.CollectionTask{
				ID:           "",
				Frequency:    1000,
				DeviceID:     "DEVICE_001",
				ProtocolType: "mock",
				VariableIDs:  []uint64{1, 2, 3},
				Priority:     5,
				Timeout:      5000,
			},
			wantErr: true,
		},
		{
			name: "频率为0",
			task: &core.CollectionTask{
				ID:           "TASK_002",
				Frequency:    0,
				DeviceID:     "DEVICE_001",
				ProtocolType: "mock",
				VariableIDs:  []uint64{1, 2, 3},
				Priority:     5,
				Timeout:      5000,
			},
			wantErr: true,
		},
		{
			name: "空设备ID",
			task: &core.CollectionTask{
				ID:           "TASK_003",
				Frequency:    1000,
				DeviceID:     "",
				ProtocolType: "mock",
				VariableIDs:  []uint64{1, 2, 3},
				Priority:     5,
				Timeout:      5000,
			},
			wantErr: true,
		},
		{
			name: "空协议类型",
			task: &core.CollectionTask{
				ID:           "TASK_004",
				Frequency:    1000,
				DeviceID:     "DEVICE_001",
				ProtocolType: "",
				VariableIDs:  []uint64{1, 2, 3},
				Priority:     5,
				Timeout:      5000,
			},
			wantErr: true,
		},
		{
			name: "空变量列表",
			task: &core.CollectionTask{
				ID:           "TASK_005",
				Frequency:    1000,
				DeviceID:     "DEVICE_001",
				ProtocolType: "mock",
				VariableIDs:  []uint64{},
				Priority:     5,
				Timeout:      5000,
			},
			wantErr: true,
		},
		{
			name: "优先级过低",
			task: &core.CollectionTask{
				ID:           "TASK_006",
				Frequency:    1000,
				DeviceID:     "DEVICE_001",
				ProtocolType: "mock",
				VariableIDs:  []uint64{1, 2, 3},
				Priority:     0,
				Timeout:      5000,
			},
			wantErr: true,
		},
		{
			name: "优先级过高",
			task: &core.CollectionTask{
				ID:           "TASK_007",
				Frequency:    1000,
				DeviceID:     "DEVICE_001",
				ProtocolType: "mock",
				VariableIDs:  []uint64{1, 2, 3},
				Priority:     11,
				Timeout:      5000,
			},
			wantErr: true,
		},
		{
			name: "超时为0",
			task: &core.CollectionTask{
				ID:           "TASK_008",
				Frequency:    1000,
				DeviceID:     "DEVICE_001",
				ProtocolType: "mock",
				VariableIDs:  []uint64{1, 2, 3},
				Priority:     5,
				Timeout:      0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := collector.validateTask(tt.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAddTask 测试添加任务
func TestAddTask(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	task := &core.CollectionTask{
		ID:           "TASK_001",
		Frequency:    1000,
		DeviceID:     "MOCK-DEVICE-001",
		ProtocolType: "mock",
		VariableIDs:  []uint64{100001, 100002},
		Priority:     5,
		Timeout:      5000,
		Enable:       true,
	}

	// 测试添加任务
	err := collector.AddTask(task)
	if err != nil {
		t.Fatalf("添加任务失败: %v", err)
	}

	// 验证任务已添加
	if _, exists := collector.tasks[task.ID]; !exists {
		t.Error("任务未添加到 tasks map")
	}

	// 测试重复添加
	err = collector.AddTask(task)
	if err == nil {
		t.Error("重复添加任务应该返回错误")
	}
}

// TestAddTaskAfterStart 测试启动后添加任务
func TestAddTaskAfterStart(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	ctx := context.Background()

	// 启动采集器
	err := collector.Start(ctx)
	if err != nil {
		t.Fatalf("启动采集器失败: %v", err)
	}
	defer collector.Stop()

	// 添加任务
	task := &core.CollectionTask{
		ID:           "TASK_001",
		Frequency:    200, // 200ms
		DeviceID:     "MOCK-DEVICE-001",
		ProtocolType: "mock",
		VariableIDs:  []uint64{100001, 100002},
		Priority:     5,
		Timeout:      100,
		Enable:       true,
	}

	err = collector.AddTask(task)
	if err != nil {
		t.Fatalf("启动后添加任务失败: %v", err)
	}

	// 验证任务运行器已创建
	time.Sleep(500 * time.Millisecond) // 等待任务执行

	collector.mu.RLock()
	runner, exists := collector.taskRunners[task.ID]
	collector.mu.RUnlock()

	if !exists {
		t.Error("任务运行器未创建")
	} else if !runner.IsRunning() {
		t.Error("任务运行器未启动")
	}

	// 验证统计数据
	stats := collector.GetStats()
	if stats.TotalCollections == 0 {
		t.Error("任务未执行采集")
	}
}

// TestRemoveTask 测试移除任务
func TestRemoveTask(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	task := &core.CollectionTask{
		ID:           "TASK_001",
		Frequency:    1000,
		DeviceID:     "MOCK-DEVICE-001",
		ProtocolType: "mock",
		VariableIDs:  []uint64{100001, 100002},
		Priority:     5,
		Timeout:      5000,
		Enable:       true,
	}

	// 添加任务
	err := collector.AddTask(task)
	if err != nil {
		t.Fatalf("添加任务失败: %v", err)
	}

	// 移除任务
	err = collector.RemoveTask(task.ID)
	if err != nil {
		t.Fatalf("移除任务失败: %v", err)
	}

	// 验证任务已移除
	if _, exists := collector.tasks[task.ID]; exists {
		t.Error("任务未从 tasks map 中移除")
	}

	// 测试移除不存在的任务
	err = collector.RemoveTask("NON_EXISTENT")
	if err == nil {
		t.Error("移除不存在的任务应该返回错误")
	}
}

// TestUpdateTask 测试更新任务
func TestUpdateTask(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	task := &core.CollectionTask{
		ID:           "TASK_001",
		Frequency:    1000,
		DeviceID:     "MOCK-DEVICE-001",
		ProtocolType: "mock",
		VariableIDs:  []uint64{100001, 100002},
		Priority:     5,
		Timeout:      5000,
		Enable:       true,
	}

	// 添加任务
	err := collector.AddTask(task)
	if err != nil {
		t.Fatalf("添加任务失败: %v", err)
	}

	// 更新任务
	task.Priority = 8
	task.Timeout = 3000

	err = collector.UpdateTask(task)
	if err != nil {
		t.Fatalf("更新任务失败: %v", err)
	}

	// 验证任务已更新
	updatedTask, err := collector.GetTask(task.ID)
	if err != nil {
		t.Fatalf("获取任务失败: %v", err)
	}

	if updatedTask.Priority != 8 {
		t.Errorf("任务优先级未更新，期望 8，实际 %d", updatedTask.Priority)
	}

	if updatedTask.Timeout != 3000 {
		t.Errorf("任务超时未更新，期望 3000，实际 %d", updatedTask.Timeout)
	}

	// 测试更新不存在的任务
	nonExistentTask := &core.CollectionTask{
		ID:           "NON_EXISTENT",
		Frequency:    1000,
		DeviceID:     "MOCK-DEVICE-001",
		ProtocolType: "mock",
		VariableIDs:  []uint64{100001},
		Priority:     5,
		Timeout:      5000,
	}

	err = collector.UpdateTask(nonExistentTask)
	if err == nil {
		t.Error("更新不存在的任务应该返回错误")
	}
}

// TestGetTask 测试获取任务
func TestGetTask(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	task := &core.CollectionTask{
		ID:           "TASK_001",
		Frequency:    1000,
		DeviceID:     "MOCK-DEVICE-001",
		ProtocolType: "mock",
		VariableIDs:  []uint64{100001, 100002},
		Priority:     5,
		Timeout:      5000,
		Enable:       true,
	}

	// 测试获取不存在的任务
	_, err := collector.GetTask("NON_EXISTENT")
	if err == nil {
		t.Error("获取不存在的任务应该返回错误")
	}

	// 添加任务
	err = collector.AddTask(task)
	if err != nil {
		t.Fatalf("添加任务失败: %v", err)
	}

	// 测试获取存在的任务
	retrievedTask, err := collector.GetTask(task.ID)
	if err != nil {
		t.Fatalf("获取任务失败: %v", err)
	}

	if retrievedTask.ID != task.ID {
		t.Errorf("获取的任务ID不匹配，期望 %s，实际 %s", task.ID, retrievedTask.ID)
	}
}

// TestListTasks 测试列出任务
func TestListTasks(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	// 初始状态应为空
	tasks := collector.ListTasks()
	if len(tasks) != 0 {
		t.Errorf("初始任务列表应为空，实际长度 %d", len(tasks))
	}

	// 添加多个任务
	task1 := &core.CollectionTask{
		ID:           "TASK_001",
		Frequency:    1000,
		DeviceID:     "DEVICE_001",
		ProtocolType: "mock",
		VariableIDs:  []uint64{1},
		Priority:     5,
		Timeout:      5000,
	}

	task2 := &core.CollectionTask{
		ID:           "TASK_002",
		Frequency:    500,
		DeviceID:     "DEVICE_002",
		ProtocolType: "mock",
		VariableIDs:  []uint64{2},
		Priority:     3,
		Timeout:      3000,
	}

	collector.AddTask(task1)
	collector.AddTask(task2)

	// 验证任务列表
	tasks = collector.ListTasks()
	if len(tasks) != 2 {
		t.Errorf("任务列表长度应为 2，实际 %d", len(tasks))
	}
}

// TestCollectorStats 测试统计信息
func TestCollectorStats(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	ctx := context.Background()

	// 启动采集器
	err := collector.Start(ctx)
	if err != nil {
		t.Fatalf("启动采集器失败: %v", err)
	}
	defer collector.Stop()

	// 添加并运行任务
	task := &core.CollectionTask{
		ID:           "TASK_001",
		Frequency:    100, // 100ms
		DeviceID:     "MOCK-DEVICE-001",
		ProtocolType: "mock",
		VariableIDs:  []uint64{100001, 100002, 100003},
		Priority:     5,
		Timeout:      50,
		Enable:       true,
	}

	err = collector.AddTask(task)
	if err != nil {
		t.Fatalf("添加任务失败: %v", err)
	}

	// 等待一些采集完成
	time.Sleep(500 * time.Millisecond)

	// 获取统计信息
	stats := collector.GetStats()

	// 验证统计信息
	if stats.TotalCollections == 0 {
		t.Error("总采集次数应大于0")
	}

	if stats.SuccessCount == 0 {
		t.Error("成功次数应大于0")
	}

	if stats.FailureCount != 0 {
		t.Errorf("失败次数应为0，实际 %d", stats.FailureCount)
	}

	if stats.LastCollectTime.IsZero() {
		t.Error("最后采集时间不应为零")
	}

	if stats.AvgDuration == 0 {
		t.Error("平均采集耗时应大于0")
	}

	t.Logf("统计信息: 总采集=%d, 成功=%d, 失败=%d, 平均耗时=%v",
		stats.TotalCollections, stats.SuccessCount, stats.FailureCount, stats.AvgDuration)
}

// TestMultipleTasks 测试多任务并发
func TestMultipleTasks(t *testing.T) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	ctx := context.Background()

	// 启动采集器
	err := collector.Start(ctx)
	if err != nil {
		t.Fatalf("启动采集器失败: %v", err)
	}
	defer collector.Stop()

	// 添加多个不同频率的任务
	tasks := []*core.CollectionTask{
		{
			ID:           "TASK_100",
			Frequency:    100,
			DeviceID:     "MOCK-DEVICE-001",
			ProtocolType: "mock",
			VariableIDs:  []uint64{100001},
			Priority:     8,
			Timeout:      50,
			Enable:       true,
		},
		{
			ID:           "TASK_200",
			Frequency:    200,
			DeviceID:     "MOCK-DEVICE-002",
			ProtocolType: "mock",
			VariableIDs:  []uint64{200001},
			Priority:     5,
			Timeout:      100,
			Enable:       true,
		},
		{
			ID:           "TASK_500",
			Frequency:    500,
			DeviceID:     "MOCK-DEVICE-003",
			ProtocolType: "mock",
			VariableIDs:  []uint64{300001},
			Priority:     3,
			Timeout:      200,
			Enable:       true,
		},
	}

	for _, task := range tasks {
		err = collector.AddTask(task)
		if err != nil {
			t.Fatalf("添加任务 %s 失败: %v", task.ID, err)
		}
	}

	// 等待任务执行
	time.Sleep(2 * time.Second)

	// 验证所有任务都在运行
	collector.mu.RLock()
	taskCount := len(collector.taskRunners)
	collector.mu.RUnlock()

	if taskCount != len(tasks) {
		t.Errorf("运行中的任务数应为 %d，实际 %d", len(tasks), taskCount)
	}

	// 验证统计信息
	stats := collector.GetStats()
	t.Logf("多任务统计: 总采集=%d, 成功=%d, 失败=%d",
		stats.TotalCollections, stats.SuccessCount, stats.FailureCount)

	if stats.TotalCollections < 10 {
		t.Errorf("总采集次数应至少为10，实际 %d", stats.TotalCollections)
	}
}

// TestConfig 测试配置管理
func TestConfig(t *testing.T) {
	// 测试默认配置
	defaultConfig := DefaultConfig()
	if defaultConfig.MaxConcurrentTasks != 100 {
		t.Errorf("默认MaxConcurrentTasks应为100，实际 %d", defaultConfig.MaxConcurrentTasks)
	}

	if defaultConfig.DefaultTimeout != 30*time.Second {
		t.Errorf("默认DefaultTimeout应为30s，实际 %v", defaultConfig.DefaultTimeout)
	}

	if !defaultConfig.EnableStatistics {
		t.Error("默认EnableStatistics应为true")
	}

	// 测试高性能配置
	hpConfig := HighPerformanceConfig()
	if hpConfig.MaxConcurrentTasks != 500 {
		t.Errorf("高性能配置MaxConcurrentTasks应为500，实际 %d", hpConfig.MaxConcurrentTasks)
	}

	// 测试低资源配置
	lrConfig := LowResourceConfig()
	if lrConfig.MaxConcurrentTasks != 20 {
		t.Errorf("低资源配置MaxConcurrentTasks应为20，实际 %d", lrConfig.MaxConcurrentTasks)
	}

	if lrConfig.EnableStatistics {
		t.Error("低资源配置EnableStatistics应为false")
	}
}

// BenchmarkCollectorAddTask 性能测试：添加任务
func BenchmarkCollectorAddTask(b *testing.B) {
	storage := storage.NewMemoryStorage()
	factory := adapter.NewAdapterFactory()
	collector := NewCollector(factory, storage)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := &core.CollectionTask{
			ID:           fmt.Sprintf("TASK_%d", i),
			Frequency:    1000,
			DeviceID:     "DEVICE_001",
			ProtocolType: "mock",
			VariableIDs:  []uint64{1, 2, 3},
			Priority:     5,
			Timeout:      5000,
		}
		collector.AddTask(task)
	}
}
