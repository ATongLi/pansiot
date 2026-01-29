package record

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"pansiot-device/internal/core"
)

// setupTestStorage 创建测试存储
func setupTestStorage(t *testing.T) (*JSONFileStorage, string) {
	t.Helper()

	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "alarm_storage_test")
	if err := os.RemoveAll(tempDir); err != nil {
		t.Fatalf("清理临时目录失败: %v", err)
	}

	storage, err := NewJSONFileStorage(tempDir)
	if err != nil {
		t.Fatalf("创建存储失败: %v", err)
	}

	return storage, tempDir
}

// teardownTestStorage 清理测试存储
func teardownTestStorage(t *testing.T, tempDir string) {
	t.Helper()
	if err := os.RemoveAll(tempDir); err != nil {
		t.Logf("清理临时目录失败: %v", err)
	}
}

// createTestRecord 创建测试记录
func createTestRecord(recordID string, ruleID string, triggerTime time.Time) *AlarmRecord {
	return &AlarmRecord{
		RecordID:     recordID,
		RuleID:       ruleID,
		RuleName:     "测试规则",
		EventType:    EventTrigger,
		Level:        3,
		State:        core.AlarmStateActive,
		Category:     "TEMP",
		TriggerTime:  triggerTime,
		AlarmMessage: "测试报警消息",
		Threshold:    80.0,
		TriggerValue: 85.0,
		ResponsibleUsers: []string{"user1", "user2"},
		CloudReported:   false,
		StorageLocation: StorageLocal,
		CreatedAt:       time.Now(),
	}
}

// TestJSONFileStorageSave 测试保存记录
func TestJSONFileStorageSave(t *testing.T) {
	storage, tempDir := setupTestStorage(t)
	defer teardownTestStorage(t, tempDir)

	// 创建测试记录
	record := createTestRecord(
		"REC_TEST_001",
		"RULE_001",
		time.Now(),
	)

	// 保存记录
	err := storage.Save(record)
	if err != nil {
		t.Fatalf("保存记录失败: %v", err)
	}

	// 验证文件已创建
	date := record.TriggerTime.Format("2006-01-02")
	expectedPath := filepath.Join(tempDir, date, record.RecordID+".json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("记录文件未创建: %s", expectedPath)
	}

	// 验证索引
	storage.indexMu.RLock()
	entry, exists := storage.index.ByID[record.RecordID]
	storage.indexMu.RUnlock()

	if !exists {
		t.Error("记录未添加到索引")
	} else {
		if entry.RuleID != record.RuleID {
			t.Errorf("索引中的规则ID不匹配: 期望 %s, 实际 %s", record.RuleID, entry.RuleID)
		}
		if entry.Level != record.Level {
			t.Errorf("索引中的级别不匹配: 期望 %d, 实际 %d", record.Level, entry.Level)
		}
	}
}

// TestJSONFileStorageQuery 测试查询记录
func TestJSONFileStorageQuery(t *testing.T) {
	storage, tempDir := setupTestStorage(t)
	defer teardownTestStorage(t, tempDir)

	// 创建多个测试记录
	now := time.Now()
	records := []*AlarmRecord{
		createTestRecord("REC_001", "RULE_001", now.Add(-2*time.Hour)),
		createTestRecord("REC_002", "RULE_001", now.Add(-1*time.Hour)),
		createTestRecord("REC_003", "RULE_002", now),
		createTestRecord("REC_004", "RULE_003", now),
	}

	// 修改部分记录的属性
	records[1].Level = 2
	records[2].Category = "PRESSURE"
	records[3].Level = 4

	// 保存所有记录
	for _, record := range records {
		if err := storage.Save(record); err != nil {
			t.Fatalf("保存记录失败: %v", err)
		}
	}

	tests := []struct {
		name     string
		query    *RecordQuery
		expected int // 期望的记录数
	}{
		{
			name:     "查询所有记录",
			query:    NewRecordQuery(),
			expected: 4,
		},
		{
			name: "按规则ID查询",
			query: NewRecordQuery().WithRuleIDs("RULE_001"),
			expected: 2,
		},
		{
			name: "按级别查询",
			query: NewRecordQuery().WithLevels(3),
			expected: 2, // records[0] and records[2] have level 3
		},
		{
			name: "按类别查询",
			query: NewRecordQuery().WithCategories("TEMP"),
			expected: 3,
		},
		{
			name: "按时间范围查询",
			query: NewRecordQuery().
				WithStartTime(now.Add(-90 * time.Minute)).
				WithEndTime(now.Add(30 * time.Minute)),
			expected: 3,
		},
		{
			name: "组合查询",
			query: NewRecordQuery().
				WithRuleIDs("RULE_001").
				WithLevels(3).
				WithLimit(10),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := storage.Query(tt.query)
			if err != nil {
				t.Fatalf("查询失败: %v", err)
			}

			if len(result) != tt.expected {
				t.Errorf("期望 %d 条记录, 实际 %d 条", tt.expected, len(result))
			}
		})
	}
}

// TestJSONFileStorageDelete 测试删除记录
func TestJSONFileStorageDelete(t *testing.T) {
	storage, tempDir := setupTestStorage(t)
	defer teardownTestStorage(t, tempDir)

	// 创建并保存记录
	record := createTestRecord("REC_DELETE_001", "RULE_001", time.Now())
	err := storage.Save(record)
	if err != nil {
		t.Fatalf("保存记录失败: %v", err)
	}

	// 验证记录存在
	if _, exists := storage.index.ByID[record.RecordID]; !exists {
		t.Error("记录未保存到索引")
	}

	// 删除记录
	err = storage.Delete(record.RecordID)
	if err != nil {
		t.Fatalf("删除记录失败: %v", err)
	}

	// 验证记录已从索引中删除
	storage.indexMu.RLock()
	_, exists := storage.index.ByID[record.RecordID]
	storage.indexMu.RUnlock()

	if exists {
		t.Error("记录未从索引中删除")
	}

	// 验证文件已删除
	date := record.TriggerTime.Format("2006-01-02")
	filePath := filepath.Join(tempDir, date, record.RecordID+".json")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("记录文件未删除")
	}
}

// TestJSONFileStorageCleanupBefore 测试清理过期记录
func TestJSONFileStorageCleanupBefore(t *testing.T) {
	storage, tempDir := setupTestStorage(t)
	defer teardownTestStorage(t, tempDir)

	now := time.Now()

	// 创建不同时间的记录
	records := []*AlarmRecord{
		createTestRecord("REC_OLD_001", "RULE_001", now.Add(-10*24*time.Hour)), // 10天前
		createTestRecord("REC_OLD_002", "RULE_002", now.Add(-5*24*time.Hour)),  // 5天前
		createTestRecord("REC_NEW_001", "RULE_003", now.Add(-2*time.Hour)),     // 2小时前
		createTestRecord("REC_NEW_002", "RULE_004", now),                       // 现在
	}

	for _, record := range records {
		if err := storage.Save(record); err != nil {
			t.Fatalf("保存记录失败: %v", err)
		}
	}

	// 清理3天前的记录
	cutoffTime := now.Add(-3 * 24 * time.Hour)
	err := storage.CleanupBefore(cutoffTime)
	if err != nil {
		t.Fatalf("清理记录失败: %v", err)
	}

	// 验证旧记录已删除
	storage.indexMu.RLock()
	_, exists1 := storage.index.ByID["REC_OLD_001"]
	_, exists2 := storage.index.ByID["REC_OLD_002"]
	storage.indexMu.RUnlock()

	if exists1 || exists2 {
		t.Error("过期记录未删除")
	}

	// 验证新记录仍存在
	storage.indexMu.RLock()
	_, exists3 := storage.index.ByID["REC_NEW_001"]
	_, exists4 := storage.index.ByID["REC_NEW_002"]
	storage.indexMu.RUnlock()

	if !exists3 || !exists4 {
		t.Error("新记录被误删")
	}
}

// TestJSONFileStorageGetStats 测试获取统计信息
func TestJSONFileStorageGetStats(t *testing.T) {
	storage, tempDir := setupTestStorage(t)
	defer teardownTestStorage(t, tempDir)

	now := time.Now()

	// 创建测试记录
	records := []*AlarmRecord{
		createTestRecord("REC_001", "RULE_001", now.Add(-2*time.Hour)),
		createTestRecord("REC_002", "RULE_001", now.Add(-1*time.Hour)),
		createTestRecord("REC_003", "RULE_002", now),
	}

	records[0].Level = 2
	records[1].Level = 3
	records[2].Category = "PRESSURE"

	for _, record := range records {
		if err := storage.Save(record); err != nil {
			t.Fatalf("保存记录失败: %v", err)
		}
	}

	// 获取统计信息
	stats, err := storage.GetStats()
	if err != nil {
		t.Fatalf("获取统计信息失败: %v", err)
	}

	// 验证总数
	if stats.TotalRecords != 3 {
		t.Errorf("期望总数 3, 实际 %d", stats.TotalRecords)
	}

	// 验证级别统计
	if stats.GetLevelCount(2) != 1 {
		t.Errorf("期望级别2有1条记录, 实际 %d", stats.GetLevelCount(2))
	}
	if stats.GetLevelCount(3) != 2 {
		t.Errorf("期望级别3有2条记录, 实际 %d", stats.GetLevelCount(3))
	}

	// 验证类别统计
	if stats.GetCategoryCount("TEMP") != 2 {
		t.Errorf("期望TEMP类别有2条记录, 实际 %d", stats.GetCategoryCount("TEMP"))
	}
}

// TestRecordIndex 测试索引功能
func TestRecordIndex(t *testing.T) {
	storage, tempDir := setupTestStorage(t)
	defer teardownTestStorage(t, tempDir)

	now := time.Now()
	records := []*AlarmRecord{
		createTestRecord("REC_001", "RULE_001", now.Add(-2*time.Hour)),
		createTestRecord("REC_002", "RULE_001", now.Add(-1*time.Hour)),
		createTestRecord("REC_003", "RULE_002", now),
	}

	for _, record := range records {
		if err := storage.Save(record); err != nil {
			t.Fatalf("保存记录失败: %v", err)
		}
	}

	// 验证规则索引
	storage.indexMu.RLock()
	rule001Records := storage.index.ByRule["RULE_001"]
	rule002Records := storage.index.ByRule["RULE_002"]
	storage.indexMu.RUnlock()

	if len(rule001Records) != 2 {
		t.Errorf("期望 RULE_001 有2条记录, 实际 %d", len(rule001Records))
	}
	if len(rule002Records) != 1 {
		t.Errorf("期望 RULE_002 有1条记录, 实际 %d", len(rule002Records))
	}

	// 验证时间索引（应该按时间降序排列，最新的在前）
	storage.indexMu.RLock()
	timeIndex := storage.index.ByTime
	storage.indexMu.RUnlock()

	if len(timeIndex) != 3 {
		t.Errorf("期望时间索引有3条记录, 实际 %d", len(timeIndex))
	}

	// 验证时间顺序（降序）
	for i := 1; i < len(timeIndex); i++ {
		if timeIndex[i-1].TriggerTime.Before(timeIndex[i].TriggerTime) {
			t.Error("时间索引未按时间降序排列")
		}
	}
}

// TestRebuildIndex 测试重建索引
func TestRebuildIndex(t *testing.T) {
	storage, tempDir := setupTestStorage(t)
	defer teardownTestStorage(t, tempDir)

	now := time.Now()
	records := []*AlarmRecord{
		createTestRecord("REC_001", "RULE_001", now.Add(-2*time.Hour)),
		createTestRecord("REC_002", "RULE_001", now.Add(-1*time.Hour)),
		createTestRecord("REC_003", "RULE_002", now),
	}

	// 保存记录（会自动更新索引）
	for _, record := range records {
		if err := storage.Save(record); err != nil {
			t.Fatalf("保存记录失败: %v", err)
		}
	}

	// 清空索引
	storage.indexMu.Lock()
	storage.index = &RecordIndex{
		ByDate:     make(map[string][]string),
		ByRule:     make(map[string][]string),
		ByLevel:    make(map[core.AlarmLevel][]string),
		ByCategory: make(map[string][]string),
		ByID:       make(map[string]*IndexEntry),
		ByTime:     make([]TimeEntry, 0),
	}
	storage.indexMu.Unlock()

	// 重建索引
	err := storage.rebuildIndex()
	if err != nil {
		t.Fatalf("重建索引失败: %v", err)
	}

	// 验证索引已重建
	storage.indexMu.RLock()
	total := len(storage.index.ByID)
	storage.indexMu.RUnlock()

	if total != 3 {
		t.Errorf("重建索引后期望3条记录, 实际 %d", total)
	}
}

// TestConcurrentOperations 测试并发操作
func TestConcurrentOperations(t *testing.T) {
	storage, tempDir := setupTestStorage(t)
	defer teardownTestStorage(t, tempDir)

	done := make(chan bool)

	// 启动多个goroutine并发写入
	for i := 0; i < 10; i++ {
		go func(id int) {
			record := createTestRecord(
				fmt.Sprintf("REC_%03d", id),
				fmt.Sprintf("RULE_%03d", id%3),
				time.Now(),
			)
			storage.Save(record)
			done <- true
		}(i)
	}

	// 等待所有写入完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有记录都已保存
	stats, _ := storage.GetStats()
	if stats.TotalRecords != 10 {
		t.Errorf("期望10条记录, 实际 %d", stats.TotalRecords)
	}
}

// BenchmarkSave 性能测试：保存记录
func BenchmarkSave(b *testing.B) {
	tempDir := filepath.Join(os.TempDir(), "alarm_storage_bench")
	os.RemoveAll(tempDir)
	storage, _ := NewJSONFileStorage(tempDir)
	defer os.RemoveAll(tempDir)

	record := createTestRecord("REC_BENCH", "RULE_001", time.Now())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		record.RecordID = fmt.Sprintf("REC_%d", i)
		storage.Save(record)
	}
}

// BenchmarkQuery 性能测试：查询记录
func BenchmarkQuery(b *testing.B) {
	tempDir := filepath.Join(os.TempDir(), "alarm_storage_bench")
	os.RemoveAll(tempDir)
	storage, _ := NewJSONFileStorage(tempDir)
	defer os.RemoveAll(tempDir)

	// 预先保存1000条记录
	now := time.Now()
	for i := 0; i < 1000; i++ {
		record := createTestRecord(
			fmt.Sprintf("REC_%d", i),
			fmt.Sprintf("RULE_%d", i%10),
			now.Add(time.Duration(i) * time.Minute),
		)
		storage.Save(record)
	}

	query := NewRecordQuery().WithLimit(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Query(query)
	}
}
