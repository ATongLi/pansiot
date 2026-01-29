package record

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"pansiot-device/internal/core"
)

// RecordStorage 报警记录存储接口
type RecordStorage interface {
	// Save 保存单条记录
	Save(record *AlarmRecord) error

	// Query 查询记录
	Query(query *RecordQuery) ([]*AlarmRecord, error)

	// Delete 删除指定记录
	Delete(recordID string) error

	// CleanupBefore 清理指定时间之前的记录
	CleanupBefore(time.Time) error

	// GetStats 获取统计信息
	GetStats() (*RecordStats, error)
}

// JSONFileStorage JSON文件存储实现
// 特点：
// - 按日期分区存储：/records/alarm/2025-01-15/REC_xxx.json
// - 内存索引加速查询
// - 支持本地和U盘存储
// - 原子写入保证数据完整性
type JSONFileStorage struct {
	basePath string        // 本地存储基础路径
	usbPath  string        // U盘存储路径（可选）
	index    *RecordIndex  // 内存索引
	indexMu  sync.RWMutex  // 索引锁
	mu       sync.RWMutex  // 读写锁
	stats    *RecordStats  // 统计信息
	statsMu  sync.Mutex    // 统计锁
}

// RecordIndex 内存索引
// 用于加速查询，避免全文件扫描
type RecordIndex struct {
	ByDate     map[string][]string              // 日期 -> [recordID]
	ByRule     map[string][]string              // 规则ID -> [recordID]
	ByLevel    map[core.AlarmLevel][]string     // 级别 -> [recordID]
	ByCategory map[string][]string              // 类别 -> [recordID]
	ByID       map[string]*IndexEntry           // recordID -> IndexEntry
	ByTime     []TimeEntry                      // 时间排序的记录索引
}

// IndexEntry 索引条目
type IndexEntry struct {
	RecordID    string
	FilePath    string
	Date        string
	RuleID      string
	Level       core.AlarmLevel
	Category    string
	TriggerTime time.Time
	EventType   EventType
}

// TimeEntry 时间索引条目
type TimeEntry struct {
	RecordID    string
	TriggerTime time.Time
}

// NewJSONFileStorage 创建JSON文件存储
func NewJSONFileStorage(basePath string) (*JSONFileStorage, error) {
	if basePath == "" {
		basePath = "./data/records/alarm"
	}

	// 确保基础路径存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	storage := &JSONFileStorage{
		basePath: basePath,
		index: &RecordIndex{
			ByDate:     make(map[string][]string),
			ByRule:     make(map[string][]string),
			ByLevel:    make(map[core.AlarmLevel][]string),
			ByCategory: make(map[string][]string),
			ByID:       make(map[string]*IndexEntry),
			ByTime:     make([]TimeEntry, 0),
		},
		stats: NewRecordStats(),
	}

	// 加载现有索引
	if err := storage.loadIndex(); err != nil {
		// 索引加载失败不影响启动，重建索引
		storage.rebuildIndex()
	}

	return storage, nil
}

// Save 保存单条记录
func (s *JSONFileStorage) Save(record *AlarmRecord) error {
	if record == nil {
		return fmt.Errorf("记录不能为空")
	}

	// 确定存储路径
	date := record.TriggerTime.Format("2006-01-02")
	dateDir := filepath.Join(s.basePath, date)
	fileName := fmt.Sprintf("%s.json", record.RecordID)
	filePath := filepath.Join(dateDir, fileName)

	// 确保日期目录存在
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		return fmt.Errorf("创建日期目录失败: %w", err)
	}

	// 原子写入：先写临时文件，再重命名
	tempPath := filePath + ".tmp"
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化记录失败: %w", err)
	}

	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		os.Remove(tempPath) // 清理临时文件
		return fmt.Errorf("重命名文件失败: %w", err)
	}

	// 更新索引
	s.updateIndex(record, filePath)

	return nil
}

// Query 查询记录
func (s *JSONFileStorage) Query(query *RecordQuery) ([]*AlarmRecord, error) {
	if query == nil {
		query = NewRecordQuery()
	}

	// 使用索引快速定位候选记录ID
	candidateIDs := s.findCandidates(query)

	// 加载记录并过滤
	var records []*AlarmRecord
	for _, recordID := range candidateIDs {
		record, err := s.loadRecordByID(recordID)
		if err != nil {
			// 记录加载失败（文件可能被删除），跳过
			continue
		}

		// 精确过滤
		if s.matchRecord(record, query) {
			records = append(records, record)
		}
	}

	return records, nil
}

// Delete 删除指定记录
func (s *JSONFileStorage) Delete(recordID string) error {
	s.indexMu.RLock()
	entry, exists := s.index.ByID[recordID]
	s.indexMu.RUnlock()

	if !exists {
		return fmt.Errorf("记录不存在: %s", recordID)
	}

	// 删除文件
	if err := os.Remove(entry.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	// 从索引中删除
	s.removeFromIndex(recordID)

	return nil
}

// CleanupBefore 清理指定时间之前的记录
func (s *JSONFileStorage) CleanupBefore(cutoffTime time.Time) error {
	s.indexMu.RLock()

	// 找出所有需要删除的记录ID
	// ByTime 是降序排列（最新的在前），所以需要遍历所有记录
	var toDelete []string
	for _, entry := range s.index.ByTime {
		if entry.TriggerTime.Before(cutoffTime) {
			toDelete = append(toDelete, entry.RecordID)
		}
		// 不要 break，因为后面的记录可能更旧
	}
	s.indexMu.RUnlock()

	// 删除记录和文件
	for _, recordID := range toDelete {
		if err := s.Delete(recordID); err != nil {
			// 继续删除其他记录
			fmt.Printf("[警告] 删除记录 %s 失败: %v\n", recordID, err)
		}
	}

	// 清理空目录
	s.cleanEmptyDirs()

	return nil
}

// GetStats 获取统计信息
func (s *JSONFileStorage) GetStats() (*RecordStats, error) {
	s.statsMu.Lock()
	defer s.statsMu.Unlock()

	// 更新统计信息
	s.statsMu.Unlock()
	s.updateStats()
	s.statsMu.Lock()

	return s.stats, nil
}

// ============ 索引管理 ============

// loadIndex 加载索引
func (s *JSONFileStorage) loadIndex() error {
	indexPath := filepath.Join(s.basePath, "records_index.json")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 索引文件不存在，重建索引
			return s.rebuildIndex()
		}
		return err
	}

	var index RecordIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return err
	}

	s.indexMu.Lock()
	s.index = &index
	s.indexMu.Unlock()

	return nil
}

// saveIndex 保存索引
func (s *JSONFileStorage) saveIndex() error {
	s.indexMu.RLock()
	defer s.indexMu.RUnlock()

	indexPath := filepath.Join(s.basePath, "records_index.json")
	data, err := json.MarshalIndent(s.index, "", "  ")
	if err != nil {
		return err
	}

	// 原子写入
	tempPath := indexPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tempPath, indexPath)
}

// rebuildIndex 重建索引
func (s *JSONFileStorage) rebuildIndex() error {
	// 清空当前索引
	s.indexMu.Lock()
	s.index = &RecordIndex{
		ByDate:     make(map[string][]string),
		ByRule:     make(map[string][]string),
		ByLevel:    make(map[core.AlarmLevel][]string),
		ByCategory: make(map[string][]string),
		ByID:       make(map[string]*IndexEntry),
		ByTime:     make([]TimeEntry, 0),
	}
	s.indexMu.Unlock()

	// 遍历所有日期目录
	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 跳过索引文件
		if entry.Name() == "records_index.json" {
			continue
		}

		// 解析日期
		date := entry.Name()
		if _, err := time.Parse("2006-01-02", date); err != nil {
			continue // 不是有效的日期目录
		}

		// 遍历日期目录下的所有记录文件
		dateDir := filepath.Join(s.basePath, date)
		files, err := os.ReadDir(dateDir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			filePath := filepath.Join(dateDir, file.Name())

			// 加载记录
			record, err := s.loadRecordFromFile(filePath)
			if err != nil {
				continue // 加载失败，跳过
			}

			// 更新索引
			s.updateIndex(record, filePath)
		}
	}

	// 保存索引
	return s.saveIndex()
}

// updateIndex 更新索引
func (s *JSONFileStorage) updateIndex(record *AlarmRecord, filePath string) {
	s.indexMu.Lock()
	defer s.indexMu.Unlock()

	date := record.TriggerTime.Format("2006-01-02")

	// 创建索引条目
	entry := &IndexEntry{
		RecordID:    record.RecordID,
		FilePath:    filePath,
		Date:        date,
		RuleID:      record.RuleID,
		Level:       record.Level,
		Category:    record.Category,
		TriggerTime: record.TriggerTime,
		EventType:   record.EventType,
	}

	// 更新各维度索引
	s.index.ByID[record.RecordID] = entry
	s.index.ByRule[record.RuleID] = append(s.index.ByRule[record.RuleID], record.RecordID)
	s.index.ByLevel[record.Level] = append(s.index.ByLevel[record.Level], record.RecordID)
	s.index.ByCategory[record.Category] = append(s.index.ByCategory[record.Category], record.RecordID)
	s.index.ByDate[date] = append(s.index.ByDate[date], record.RecordID)

	// 更新时间索引（有序插入 - 降序，最新的在前）
	timeEntry := TimeEntry{
		RecordID:    record.RecordID,
		TriggerTime: record.TriggerTime,
	}
	// 使用二分查找插入位置（按时间降序）
	idx := sort.Search(len(s.index.ByTime), func(i int) bool {
		return s.index.ByTime[i].TriggerTime.Before(timeEntry.TriggerTime)
	})
	s.index.ByTime = append(s.index.ByTime, TimeEntry{})
	copy(s.index.ByTime[idx+1:], s.index.ByTime[idx:])
	s.index.ByTime[idx] = timeEntry

	// 定期保存索引（每100次更新保存一次）
	if len(s.index.ByID)%100 == 0 {
		go s.saveIndex()
	}
}

// removeFromIndex 从索引中删除
func (s *JSONFileStorage) removeFromIndex(recordID string) {
	s.indexMu.Lock()
	defer s.indexMu.Unlock()

	entry, exists := s.index.ByID[recordID]
	if !exists {
		return
	}

	// 从各维度索引中删除
	delete(s.index.ByID, recordID)

	// 从 ByRule 删除
	s.removeFromSlice(s.index.ByRule[entry.RuleID], recordID)

	// 从 ByLevel 删除
	s.removeFromSlice(s.index.ByLevel[entry.Level], recordID)

	// 从 ByCategory 删除
	s.removeFromSlice(s.index.ByCategory[entry.Category], recordID)

	// 从 ByDate 删除
	s.removeFromSlice(s.index.ByDate[entry.Date], recordID)

	// 从 ByTime 删除
	for i, timeEntry := range s.index.ByTime {
		if timeEntry.RecordID == recordID {
			s.index.ByTime = append(s.index.ByTime[:i], s.index.ByTime[i+1:]...)
			break
		}
	}

	// 保存索引
	go s.saveIndex()
}

// findCandidates 使用索引快速查找候选记录
func (s *JSONFileStorage) findCandidates(query *RecordQuery) []string {
	s.indexMu.RLock()
	defer s.indexMu.RUnlock()

	// 如果没有任何过滤条件，返回所有记录ID
	if query == nil ||
		(len(query.RuleIDs) == 0 &&
			len(query.Levels) == 0 &&
			len(query.Categories) == 0 &&
			query.StartTime == nil &&
			query.EndTime == nil) {
		result := make([]string, len(s.index.ByTime))
		for i, entry := range s.index.ByTime {
			result[i] = entry.RecordID
		}
		return result
	}

	// 使用最有效的索引
	var candidates []string

	// 优先使用规则ID索引（最快）
	if len(query.RuleIDs) > 0 {
		for _, ruleID := range query.RuleIDs {
			if ids, ok := s.index.ByRule[ruleID]; ok {
				candidates = append(candidates, ids...)
			}
		}
	}

	// 其次使用级别索引
	if len(candidates) == 0 && len(query.Levels) > 0 {
		for _, level := range query.Levels {
			if ids, ok := s.index.ByLevel[level]; ok {
				candidates = append(candidates, ids...)
			}
		}
	}

	// 其次使用类别索引
	if len(candidates) == 0 && len(query.Categories) > 0 {
		for _, category := range query.Categories {
			if ids, ok := s.index.ByCategory[category]; ok {
				candidates = append(candidates, ids...)
			}
		}
	}

	// 如果还是没有候选记录，使用时间范围
	if len(candidates) == 0 && (query.StartTime != nil || query.EndTime != nil) {
		for _, entry := range s.index.ByTime {
			if query.StartTime != nil && entry.TriggerTime.Before(*query.StartTime) {
				continue
			}
			if query.EndTime != nil && entry.TriggerTime.After(*query.EndTime) {
				continue
			}
			candidates = append(candidates, entry.RecordID)
		}
	}

	// 如果还是没有，返回所有记录
	if len(candidates) == 0 {
		candidates = make([]string, len(s.index.ByTime))
		for i, entry := range s.index.ByTime {
			candidates[i] = entry.RecordID
		}
	}

	return candidates
}

// matchRecord 判断记录是否匹配查询条件
func (s *JSONFileStorage) matchRecord(record *AlarmRecord, query *RecordQuery) bool {
	// 时间范围
	if query.StartTime != nil && record.TriggerTime.Before(*query.StartTime) {
		return false
	}
	if query.EndTime != nil && record.TriggerTime.After(*query.EndTime) {
		return false
	}

	// 规则ID
	if len(query.RuleIDs) > 0 {
		found := false
		for _, ruleID := range query.RuleIDs {
			if record.RuleID == ruleID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 级别
	if len(query.Levels) > 0 {
		found := false
		for _, level := range query.Levels {
			if record.Level == level {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 类别
	if len(query.Categories) > 0 {
		found := false
		for _, category := range query.Categories {
			if record.Category == category {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 状态
	if len(query.States) > 0 {
		found := false
		for _, state := range query.States {
			if record.State == state {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 事件类型
	if len(query.EventTypes) > 0 {
		found := false
		for _, eventType := range query.EventTypes {
			if record.EventType == eventType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// loadRecordByID 根据ID加载记录
func (s *JSONFileStorage) loadRecordByID(recordID string) (*AlarmRecord, error) {
	s.indexMu.RLock()
	entry, exists := s.index.ByID[recordID]
	s.indexMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("记录不存在: %s", recordID)
	}

	return s.loadRecordFromFile(entry.FilePath)
}

// loadRecordFromFile 从文件加载记录
func (s *JSONFileStorage) loadRecordFromFile(filePath string) (*AlarmRecord, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var record AlarmRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

// updateStats 更新统计信息
func (s *JSONFileStorage) updateStats() {
	s.indexMu.RLock()
	defer s.indexMu.RUnlock()

	// 记录总数
	s.stats.TotalRecords = int64(len(s.index.ByID))

	// 触发和恢复计数
	triggerCount := int64(0)
	recoverCount := int64(0)
	activeCount := int64(0)

	for _, entry := range s.index.ByID {
		if entry.EventType == EventTrigger {
			triggerCount++
		} else if entry.EventType == EventRecover {
			recoverCount++
		}

		// 统计激活状态（简化处理：假设最近触发的是激活状态）
		if entry.EventType == EventTrigger {
			activeCount++
		}

		// 按级别统计
		s.stats.IncrementLevel(entry.Level)

		// 按类别统计
		s.stats.IncrementCategory(entry.Category)

		// 按规则统计
		s.stats.IncrementRule(entry.RuleID)
	}

	s.stats.TriggerCount = triggerCount
	s.stats.RecoverCount = recoverCount
	s.stats.ActiveCount = activeCount

	// 首次和末次记录时间
	if len(s.index.ByTime) > 0 {
		firstTime := s.index.ByTime[0].TriggerTime
		s.stats.FirstRecordTime = &firstTime

		lastTime := s.index.ByTime[len(s.index.ByTime)-1].TriggerTime
		s.stats.LastRecordTime = &lastTime
	}
}

// removeFromSlice 从切片中删除元素
func (s *JSONFileStorage) removeFromSlice(slice []string, item string) []string {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// cleanEmptyDirs 清理空目录
func (s *JSONFileStorage) cleanEmptyDirs() {
	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(s.basePath, entry.Name())

		// 检查目录是否为空
		files, err := os.ReadDir(dirPath)
		if err != nil || len(files) > 0 {
			continue
		}

		// 删除空目录
		os.Remove(dirPath)
	}
}
