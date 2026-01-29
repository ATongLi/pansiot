package storage

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"pansiot-device/internal/core"
)

func BenchmarkReadVar(b *testing.B) {
	ms := setupBenchmarkStorage(100000)
	variableID := uint64(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.ReadVar(variableID)
	}
}

func BenchmarkWriteVar(b *testing.B) {
	ms := setupBenchmarkStorage(100000)
	variable := &core.Variable{
		ID:        100000,
		StringID:  "DV-PLC001-VAR000000",
		Value:     25.5,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.WriteVar(variable)
	}
}

func BenchmarkReadVars(b *testing.B) {
	ms := setupBenchmarkStorage(100000)
	variableIDs := make([]uint64, 100)
	for i := 0; i < 100; i++ {
		variableIDs[i] = uint64(100000 + i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.ReadVars(variableIDs)
	}
}

func BenchmarkWriteVars(b *testing.B) {
	ms := setupBenchmarkStorage(100000)
	variables := make([]*core.Variable, 100)
	for i := 0; i < 100; i++ {
		variables[i] = &core.Variable{
			ID:        uint64(100000 + i),
			StringID:  fmt.Sprintf("DV-PLC001-VAR%06d", i),
			Value:     float64(i),
			Quality:   core.QualityGood,
			Timestamp: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.WriteVars(variables)
	}
}

func BenchmarkConcurrentReadWrite(b *testing.B) {
	ms := setupBenchmarkStorage(100000)
	variableIDs := make([]uint64, 1000)
	for i := 0; i < 1000; i++ {
		variableIDs[i] = uint64(100000 + i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		vid := variableIDs[0]
		variable := &core.Variable{
			ID:        vid,
			StringID:  "DV-PLC001-VAR000000",
			Value:     25.5,
			Quality:   core.QualityGood,
			Timestamp: time.Now(),
		}

		for pb.Next() {
			// 80% 读取, 20% 写入
			if b.N%5 == 0 {
				ms.WriteVar(variable)
			} else {
				ms.ReadVar(vid)
			}
		}
	})
}

func BenchmarkSubscribe(b *testing.B) {
	ms := setupBenchmarkStorage(100000)
	notifyCount := 0
	callback := func(update core.VariableUpdate) {
		notifyCount++
	}

	ms.Subscribe("bench-subscriber", []uint64{100000}, callback)

	variable := &core.Variable{
		ID:        100000,
		StringID:  "DV-PLC001-VAR000000",
		Value:     25.5,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.WriteVar(variable)
	}
}

func BenchmarkReadVarByStringID(b *testing.B) {
	ms := setupBenchmarkStorage(100000)
	stringID := "DV-PLC001-VAR000000"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.ReadVarByStringID(stringID)
	}
}

func BenchmarkListVariables(b *testing.B) {
	ms := setupBenchmarkStorage(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.ListVariables()
	}
}

func BenchmarkListVariablesByDevice(b *testing.B) {
	ms := setupBenchmarkStorage(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.ListVariablesByDevice("PLC001")
	}
}

// 测试读写延迟
func TestReadWriteLatency(t *testing.T) {
	ms := setupBenchmarkStorage(100000)
	variable := &core.Variable{
		ID:        100000,
		StringID:  "DV-PLC001-VAR000000",
		Value:     25.5,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}

	// 测试写入延迟
	iterations := 1000
	writeTimes := make([]time.Duration, iterations)
	for i := 0; i < iterations; i++ {
		start := time.Now()
		ms.WriteVar(variable)
		writeTimes[i] = time.Since(start)
	}

	avgWrite := averageDuration(writeTimes)
	t.Logf("Average write latency: %v", avgWrite)
	if avgWrite > 100*time.Microsecond {
		t.Errorf("Write latency exceeds 100µs: %v", avgWrite)
	}

	// 测试读取延迟
	readTimes := make([]time.Duration, iterations)
	for i := 0; i < iterations; i++ {
		start := time.Now()
		ms.ReadVar(100000)
		readTimes[i] = time.Since(start)
	}

	avgRead := averageDuration(readTimes)
	t.Logf("Average read latency: %v", avgRead)
	if avgRead > 100*time.Microsecond {
		t.Errorf("Read latency exceeds 100µs: %v", avgRead)
	}
}

// 测试10万+变量的性能
func TestLargeScalePerformance(t *testing.T) {
	t.Logf("Creating 100000 variables...")
	start := time.Now()
	ms := setupBenchmarkStorage(100000)
	creationTime := time.Since(start)
	t.Logf("Creation time: %v", creationTime)

	// 测试单次读取
	start = time.Now()
	ms.ReadVar(100000)
	readTime := time.Since(start)
	t.Logf("Single read latency: %v", readTime)

	// 测试批量读取
	start = time.Now()
	variableIDs := make([]uint64, 10000)
	for i := 0; i < 10000; i++ {
		variableIDs[i] = uint64(100000 + i)
	}
	ms.ReadVars(variableIDs)
	batchReadTime := time.Since(start)
	t.Logf("Batch read (10k vars): %v", batchReadTime)

	// 获取统计信息
	stats := ms.GetStats()
	t.Logf("Total variables: %d", stats.TotalVariables)
	t.Logf("Total reads: %d", stats.ReadCount)
	t.Logf("Total writes: %d", stats.WriteCount)
}

// 测试并发压力
func TestConcurrentPressure(t *testing.T) {
	ms := setupBenchmarkStorage(100000)

	const numGoroutines = 100
	const opsPerGoroutine = 1000

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			variable := &core.Variable{
				ID:        uint64(100000 + goroutineID%100000),
				StringID:  fmt.Sprintf("DV-PLC001-VAR%06d", goroutineID%100000),
				Value:     float64(goroutineID),
				Quality:   core.QualityGood,
				Timestamp: time.Now(),
			}

			for j := 0; j < opsPerGoroutine; j++ {
				ms.WriteVar(variable)
				ms.ReadVar(variable.ID)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalOps := numGoroutines * opsPerGoroutine * 2 // 读 + 写
	opsPerSecond := float64(totalOps) / elapsed.Seconds()

	t.Logf("Completed %d ops in %v (%.2f ops/sec)", totalOps, elapsed, opsPerSecond)
}

// setupBenchmarkStorage 创建基准测试用的存储实例
func setupBenchmarkStorage(count int) *MemoryStorage {
	ms := NewMemoryStorage()

	batchSize := 1000
	for i := 0; i < count; i += batchSize {
		variables := make([]*core.Variable, 0, batchSize)
		for j := 0; j < batchSize && i+j < count; j++ {
			variables = append(variables, &core.Variable{
				ID:        uint64(100000 + i + j),
				StringID:  fmt.Sprintf("DV-PLC001-VAR%06d", i+j),
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

// averageDuration 计算平均时长
func averageDuration(durations []time.Duration) time.Duration {
	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	return sum / time.Duration(len(durations))
}
