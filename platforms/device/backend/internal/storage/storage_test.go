package storage

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"pansiot-device/internal/core"
)

// TestCreateVariable 测试创建变量
func TestCreateVariable(t *testing.T) {
	ms := NewMemoryStorage()

	variable := &core.Variable{
		ID:        100000,
		StringID:  "DV-PLC001-TEMP01",
		Name:      "温度",
		DataType:  core.DataTypeFloat32,
		DeviceID:  "PLC001",
		Value:     25.5,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}

	err := ms.CreateVariable(variable)
	if err != nil {
		t.Fatalf("Failed to create variable: %v", err)
	}

	// 验证变量已创建
	readVar, err := ms.ReadVar(variable.ID)
	if err != nil {
		t.Fatalf("Failed to read variable: %v", err)
	}

	if readVar.ID != variable.ID {
		t.Errorf("Expected ID %d, got %d", variable.ID, readVar.ID)
	}
	if readVar.Value != variable.Value {
		t.Errorf("Expected value %v, got %v", variable.Value, readVar.Value)
	}
}

// TestReadVar 测试读取变量
func TestReadVar(t *testing.T) {
	ms := setupTestStorage(100)

	variable, err := ms.ReadVar(100000)
	if err != nil {
		t.Fatalf("Failed to read variable: %v", err)
	}

	if variable.ID != 100000 {
		t.Errorf("Expected ID 100000, got %d", variable.ID)
	}
}

// TestReadVarNotFound 测试读取不存在的变量
func TestReadVarNotFound(t *testing.T) {
	ms := NewMemoryStorage()

	_, err := ms.ReadVar(999999)
	if err == nil {
		t.Error("Expected error when reading non-existent variable")
	}
}

// TestReadVars 测试批量读取
func TestReadVars(t *testing.T) {
	ms := setupTestStorage(100)

	variableIDs := make([]uint64, 10)
	for i := 0; i < 10; i++ {
		variableIDs[i] = uint64(100000 + i)
	}

	variables, err := ms.ReadVars(variableIDs)
	if err != nil {
		t.Fatalf("Failed to read variables: %v", err)
	}

	if len(variables) != 10 {
		t.Errorf("Expected 10 variables, got %d", len(variables))
	}
}

// TestWriteVar 测试写入变量
func TestWriteVar(t *testing.T) {
	ms := setupTestStorage(100)

	variable := &core.Variable{
		ID:        100000,
		StringID:  "DV-PLC001-TEMP01",
		Value:     30.5,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}

	err := ms.WriteVar(variable)
	if err != nil {
		t.Fatalf("Failed to write variable: %v", err)
	}

	// 验证值已更新
	readVar, _ := ms.ReadVar(100000)
	if readVar.Value != 30.5 {
		t.Errorf("Expected value 30.5, got %v", readVar.Value)
	}
}

// TestWriteVars 测试批量写入
func TestWriteVars(t *testing.T) {
	ms := setupTestStorage(100)

	variables := make([]*core.Variable, 10)
	for i := 0; i < 10; i++ {
		variables[i] = &core.Variable{
			ID:        uint64(100000 + i),
			StringID:  core.BuildStringID("DV", "PLC001", fmt.Sprintf("VAR%02d", i)),
			Value:     float64(i + 100),
			Quality:   core.QualityGood,
			Timestamp: time.Now(),
		}
	}

	err := ms.WriteVars(variables)
	if err != nil {
		t.Fatalf("Failed to write variables: %v", err)
	}

	// 验证值已更新
	variableIDs := make([]uint64, 10)
	for i := 0; i < 10; i++ {
		variableIDs[i] = uint64(100000 + i)
	}

	readVars, _ := ms.ReadVars(variableIDs)
	for i, v := range readVars {
		expectedValue := float64(i + 100)
		if v.Value != expectedValue {
			t.Errorf("Variable %d: expected value %v, got %v", i, expectedValue, v.Value)
		}
	}
}

// TestReadVarByStringID 测试通过字符串ID读取
func TestReadVarByStringID(t *testing.T) {
	ms := setupTestStorage(100)

	variable, err := ms.ReadVarByStringID("DV-PLC001-VAR000000")
	if err != nil {
		t.Fatalf("Failed to read variable by string ID: %v", err)
	}

	if variable.StringID != "DV-PLC001-VAR000000" {
		t.Errorf("Expected StringID DV-PLC001-VAR000000, got %s", variable.StringID)
	}
}

// TestSubscribe 测试订阅
func TestSubscribe(t *testing.T) {
	ms := setupTestStorage(100)

	var notifyCount atomic.Int32
	var receivedValue atomic.Value

	err := ms.Subscribe("test-subscriber", []uint64{100000}, func(update core.VariableUpdate) {
		notifyCount.Add(1)
		receivedValue.Store(update.Value)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 写入变量
	variable := &core.Variable{
		ID:        100000,
		StringID:  "DV-PLC001-VAR000000",
		Value:     99.9,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}
	ms.WriteVar(variable)

	// 等待通知
	time.Sleep(200 * time.Millisecond)

	if notifyCount.Load() < 1 {
		t.Error("Expected at least 1 notification")
	}

	value := receivedValue.Load()
	if value == nil {
		t.Error("Expected to receive value")
	}
}

// TestSubscribeByDevice 测试按设备订阅
func TestSubscribeByDevice(t *testing.T) {
	ms := setupTestStorage(100)

	var notifyCount atomic.Int32

	err := ms.SubscribeByDevice("test-subscriber", "PLC001", func(update core.VariableUpdate) {
		notifyCount.Add(1)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe by device: %v", err)
	}

	// 写入多个变量
	for i := 0; i < 10; i++ {
		variable := &core.Variable{
			ID:        uint64(100000 + i),
			StringID:  core.BuildStringID("DV", "PLC001", fmt.Sprintf("VAR%02d", i)),
			Value:     float64(i),
			Quality:   core.QualityGood,
			Timestamp: time.Now(),
			DeviceID:  "PLC001",
		}
		ms.WriteVar(variable)
	}

	// 等待通知
	time.Sleep(200 * time.Millisecond)

	if notifyCount.Load() < 10 {
		t.Errorf("Expected at least 10 notifications, got %d", notifyCount.Load())
	}
}

// TestSubscribeByPattern 测试按模式订阅
func TestSubscribeByPattern(t *testing.T) {
	ms := setupTestStorage(100)

	var notifyCount atomic.Int32

	err := ms.SubscribeByPattern("test-subscriber", "DV-PLC001-VAR00*", func(update core.VariableUpdate) {
		notifyCount.Add(1)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe by pattern: %v", err)
	}

	// 写入匹配的变量
	for i := 0; i < 5; i++ {
		variable := &core.Variable{
			ID:        uint64(100000 + i),
			StringID:  core.BuildStringID("DV", "PLC001", fmt.Sprintf("VAR00%d", i)),
			Value:     float64(i),
			Quality:   core.QualityGood,
			Timestamp: time.Now(),
			DeviceID:  "PLC001",
		}
		ms.WriteVar(variable)
	}

	// 等待通知
	time.Sleep(200 * time.Millisecond)

	if notifyCount.Load() < 5 {
		t.Errorf("Expected at least 5 notifications, got %d", notifyCount.Load())
	}
}

// TestUnsubscribe 测试取消订阅
func TestUnsubscribe(t *testing.T) {
	ms := setupTestStorage(100)

	var notifyCount atomic.Int32

	err := ms.Subscribe("test-subscriber", []uint64{100000}, func(update core.VariableUpdate) {
		notifyCount.Add(1)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 取消订阅
	err = ms.Unsubscribe("test-subscriber", []uint64{100000})
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// 写入变量
	variable := &core.Variable{
		ID:        100000,
		StringID:  "DV-PLC001-VAR000000",
		Value:     99.9,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}
	ms.WriteVar(variable)

	// 等待通知
	time.Sleep(200 * time.Millisecond)

	// 应该没有通知
	if notifyCount.Load() > 0 {
		t.Error("Expected no notifications after unsubscribe")
	}
}

// TestUnsubscribeAll 测试取消所有订阅
func TestUnsubscribeAll(t *testing.T) {
	ms := setupTestStorage(100)

	var notifyCount atomic.Int32

	// 订阅多个变量
	variableIDs := []uint64{100000, 100001, 100002}
	err := ms.Subscribe("test-subscriber", variableIDs, func(update core.VariableUpdate) {
		notifyCount.Add(1)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 取消所有订阅
	err = ms.UnsubscribeAll("test-subscriber")
	if err != nil {
		t.Fatalf("Failed to unsubscribe all: %v", err)
	}

	// 写入变量
	for _, vid := range variableIDs {
		variable := &core.Variable{
			ID:        vid,
			StringID:  fmt.Sprintf("DV-PLC001-VAR%06d", vid-100000),
			Value:     99.9,
			Quality:   core.QualityGood,
			Timestamp: time.Now(),
		}
		ms.WriteVar(variable)
	}

	// 等待通知
	time.Sleep(200 * time.Millisecond)

	// 应该没有通知
	if notifyCount.Load() > 0 {
		t.Error("Expected no notifications after unsubscribe all")
	}
}

// TestDeleteVariable 测试删除变量
func TestDeleteVariable(t *testing.T) {
	ms := setupTestStorage(100)

	err := ms.DeleteVariable(100000)
	if err != nil {
		t.Fatalf("Failed to delete variable: %v", err)
	}

	// 验证变量已删除
	_, err = ms.ReadVar(100000)
	if err == nil {
		t.Error("Expected error when reading deleted variable")
	}
}

// TestListVariables 测试列出所有变量
func TestListVariables(t *testing.T) {
	ms := setupTestStorage(100)

	variables := ms.ListVariables()
	if len(variables) != 100 {
		t.Errorf("Expected 100 variables, got %d", len(variables))
	}
}

// TestListVariablesByDevice 测试按设备列出变量
func TestListVariablesByDevice(t *testing.T) {
	ms := setupTestStorage(100)

	variables := ms.ListVariablesByDevice("PLC001")
	if len(variables) != 100 {
		t.Errorf("Expected 100 variables for PLC001, got %d", len(variables))
	}
}

// TestGetStats 测试获取统计信息
func TestGetStats(t *testing.T) {
	ms := setupTestStorage(100)

	stats := ms.GetStats()
	if stats.TotalVariables != 100 {
		t.Errorf("Expected 100 total variables, got %d", stats.TotalVariables)
	}
	if stats.ReadCount == 0 {
		// 可能在setupTestStorage中有读取操作
		ms.ReadVar(100000)
		stats = ms.GetStats()
		if stats.ReadCount == 0 {
			t.Error("Expected non-zero read count")
		}
	}
}

// TestConcurrentReadWrite 测试并发读写
func TestConcurrentReadWrite(t *testing.T) {
	ms := setupTestStorage(100)

	const numGoroutines = 50
	const opsPerGoroutine = 100

	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			variable := &core.Variable{
				ID:        uint64(100000 + goroutineID%100),
				StringID:  fmt.Sprintf("DV-PLC001-VAR%06d", goroutineID%100),
				Value:     float64(goroutineID),
				Quality:   core.QualityGood,
				Timestamp: time.Now(),
			}

			for j := 0; j < opsPerGoroutine; j++ {
				// 80% 读取, 20% 写入
				if j%5 == 0 {
					ms.WriteVar(variable)
				} else {
					ms.ReadVar(variable.ID)
				}
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(startTime)

	t.Logf("Concurrent test completed: %d goroutines, %d ops each in %v", numGoroutines, opsPerGoroutine, elapsed)

	// 验证数据一致性
	variable, _ := ms.ReadVar(100000)
	if variable == nil {
		t.Error("Expected variable to exist after concurrent operations")
	}
}

// TestConcurrentWriteSameVariable 测试并发写入同一变量
func TestConcurrentWriteSameVariable(t *testing.T) {
	ms := setupTestStorage(10)

	const numGoroutines = 100
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(value float64) {
			defer wg.Done()
			variable := &core.Variable{
				ID:        100000,
				StringID:  "DV-PLC001-VAR000000",
				Value:     value,
				Quality:   core.QualityGood,
				Timestamp: time.Now(),
			}
			ms.WriteVar(variable)
		}(float64(i))
	}

	wg.Wait()

	// 验证变量存在
	variable, err := ms.ReadVar(100000)
	if err != nil {
		t.Fatalf("Failed to read variable after concurrent writes: %v", err)
	}

	if variable == nil {
		t.Error("Expected variable to exist")
	}
}

// TestEmptyBatchRead 测试空批量读取
func TestEmptyBatchRead(t *testing.T) {
	ms := NewMemoryStorage()

	variables, err := ms.ReadVars([]uint64{})
	if err != nil {
		t.Fatalf("Failed to read empty batch: %v", err)
	}

	if len(variables) != 0 {
		t.Errorf("Expected 0 variables, got %d", len(variables))
	}
}

// TestEmptyBatchWrite 测试空批量写入
func TestEmptyBatchWrite(t *testing.T) {
	ms := NewMemoryStorage()

	err := ms.WriteVars([]*core.Variable{})
	if err != nil {
		t.Fatalf("Failed to write empty batch: %v", err)
	}
}

// TestDuplicateVariableCreation 测试重复创建变量
func TestDuplicateVariableCreation(t *testing.T) {
	ms := NewMemoryStorage()

	variable := &core.Variable{
		ID:        100000,
		StringID:  "DV-PLC001-TEMP01",
		Name:      "温度",
		DataType:  core.DataTypeFloat32,
		DeviceID:  "PLC001",
		Value:     25.5,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}

	err := ms.CreateVariable(variable)
	if err != nil {
		t.Fatalf("Failed to create variable: %v", err)
	}

	// 尝试再次创建
	err = ms.CreateVariable(variable)
	if err == nil {
		t.Error("Expected error when creating duplicate variable")
	}
}

// setupTestStorage 创建测试用的存储实例
func setupTestStorage(count int) *MemoryStorage {
	ms := NewMemoryStorage()

	batchSize := 100
	for i := 0; i < count; i += batchSize {
		variables := make([]*core.Variable, 0, batchSize)
		for j := 0; j < batchSize && i+j < count; j++ {
			variables = append(variables, &core.Variable{
				ID:        uint64(100000 + i + j),
				StringID:  core.BuildStringID("DV", "PLC001", fmt.Sprintf("VAR%06d", i+j)),
				Name:      fmt.Sprintf("Variable %d", i+j),
				DataType:  core.DataTypeFloat32,
				DeviceID:  "PLC001",
				Value:     float64(i + j),
				Quality:   core.QualityGood,
				Timestamp: time.Now(),
			})
		}
		ms.WriteVars(variables)
	}

	return ms
}
