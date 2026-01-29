package record

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/core"
)

// RecordConfig 报警记录管理器配置
type RecordConfig struct {
	// 存储路径
	BasePath string        // 本地存储基础路径 (默认: ./data/records/alarm)
	USBPath  string        // U盘存储路径 (可选)

	// 缓存配置
	CacheSize int           // 缓存大小 (默认: 1000)
	CacheTTL  time.Duration // 缓存过期时间 (默认: 1小时)

	// 异步写入配置
	BatchSize      int           // 批量写入大小 (默认: 10)
	BatchTimeout   time.Duration // 批量写入超时 (默认: 5秒)
	QueueSize      int           // 待写入队列大小 (默认: 1000)

	// 导出配置
	ExportDir string // 导出文件目录 (默认: ./data/exports)

	// 清理配置
	MaxRecords      int           // 最大记录数 (0=不限制)
	MaxDays         int           // 最大保留天数 (0=不限制)
	CleanupInterval time.Duration // 清理间隔 (默认: 1小时)

	// 自动上报
	AutoCloudReport bool          // 是否自动上报云端 (默认: false)
}

// DefaultRecordConfig 返回默认配置
func DefaultRecordConfig() *RecordConfig {
	return &RecordConfig{
		BasePath:       "./data/records/alarm",
		CacheSize:      1000,
		CacheTTL:       1 * time.Hour,
		BatchSize:      10,
		BatchTimeout:   5 * time.Second,
		QueueSize:      1000,
		ExportDir:      "./data/exports",
		MaxDays:        30,
		CleanupInterval: 1 * time.Hour,
		AutoCloudReport: false,
	}
}

// RecordManager 报警记录管理器
// 核心职责：
// 1. 接收报警事件（触发/恢复）并持久化
// 2. 提供查询功能（支持多条件过滤、分页、排序）
// 3. 提供导出功能（CSV/JSON）
// 4. 自动清理过期记录
type RecordManager struct {
	storage    RecordStorage          // 存储接口
	config     *RecordConfig          // 配置
	stats      *RecordStats           // 统计信息

	// 缓存层
	cache      []*AlarmRecord         // 环形缓存
	cacheIndex int                    // 当前缓存写入位置
	cacheMu    sync.RWMutex

	// 异步写入
	pendingChan chan *AlarmRecord     // 待写入记录队列
	batch       []*AlarmRecord        // 批量写入缓冲区
	batchMu     sync.Mutex
	batchTimer  *time.Timer           // 批量写入定时器

	// 生命周期
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	running    atomic.Bool

	// 防止重复记录
	lastTriggerTime sync.Map // ruleID -> time.Time (防止短时间内重复触发)
}

// NewRecordManager 创建报警记录管理器
func NewRecordManager(storage RecordStorage, config *RecordConfig) *RecordManager {
	if config == nil {
		config = DefaultRecordConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	rm := &RecordManager{
		storage:     storage,
		config:      config,
		stats:       NewRecordStats(),
		cache:       make([]*AlarmRecord, config.CacheSize),
		pendingChan: make(chan *AlarmRecord, config.QueueSize),
		batch:       make([]*AlarmRecord, 0, config.BatchSize),
		batchTimer:  time.NewTimer(config.BatchTimeout),
		ctx:         ctx,
		cancel:      cancel,
	}

	// 启动后台处理协程
	rm.startBackgroundWorkers()

	return rm
}

// Start 启动记录管理器
func (rm *RecordManager) Start() error {
	if !rm.running.CompareAndSwap(false, true) {
		return fmt.Errorf("记录管理器已经在运行")
	}

	// 启动清理协程
	if rm.config.MaxDays > 0 || rm.config.MaxRecords > 0 {
		rm.wg.Add(1)
		go rm.cleanupLoop()
	}

	return nil
}

// Stop 停止记录管理器
func (rm *RecordManager) Stop() error {
	if !rm.running.CompareAndSwap(true, false) {
		return fmt.Errorf("记录管理器未在运行")
	}

	// 停止后台协程
	rm.cancel()

	// 刷新剩余批次
	rm.flushBatch()

	// 等待所有协程结束
	rm.wg.Wait()

	rm.batchTimer.Stop()

	return nil
}

// RecordTrigger 记录报警触发事件
// 当报警从 Inactive → Active 时调用
func (rm *RecordManager) RecordTrigger(alarm *engine.ActiveAlarm) error {
	if !rm.running.Load() {
		return fmt.Errorf("记录管理器未启动")
	}

	// 防止重复记录：同一规则在1秒内只记录一次触发
	if lastTime, ok := rm.lastTriggerTime.Load(alarm.RuleID); ok {
		if time.Since(lastTime.(time.Time)) < 1*time.Second {
			return nil // 跳过重复记录
		}
	}

	// 创建报警记录
	record := &AlarmRecord{
		RecordID:     rm.generateRecordID(),
		RuleID:       alarm.Rule.ID,
		RuleName:     alarm.Rule.Name,
		EventType:    EventTrigger,
		Level:        alarm.Rule.Level,
		State:        core.AlarmStateActive,
		Category:     alarm.Rule.Category,
		TriggerTime:  alarm.TriggerTime,
		AlarmMessage: rm.renderAlarmMessage(alarm),
		Threshold:    rm.extractThreshold(alarm),
		TriggerValue: alarm.TriggerValue,
		ResponsibleUsers: alarm.Rule.Responsible,
		TriggerVariables:  rm.extractTriggerVariables(alarm),
		CloudReported:   false,
		StorageLocation: StorageLocal,
		CreatedAt:       time.Now(),
	}

	// 异步写入
	return rm.writeRecord(record)
}

// RecordRecover 记录报警恢复事件
// 当报警从 Active → Cleared 时调用
func (rm *RecordManager) RecordRecover(ruleID string) error {
	if !rm.running.Load() {
		return fmt.Errorf("记录管理器未启动")
	}

	// 查找最近的触发记录（用于计算持续时间）
	triggerRecord, err := rm.findLatestTriggerRecord(ruleID)
	if err != nil {
		return fmt.Errorf("查找触发记录失败: %w", err)
	}

	// 计算持续时间
	now := time.Now()
	duration := now.Sub(triggerRecord.TriggerTime)

	// 创建恢复记录
	record := &AlarmRecord{
		RecordID:     rm.generateRecordID(),
		RuleID:       triggerRecord.RuleID,
		RuleName:     triggerRecord.RuleName,
		EventType:    EventRecover,
		Level:        triggerRecord.Level,
		State:        core.AlarmStateCleared,
		Category:     triggerRecord.Category,
		TriggerTime:  triggerRecord.TriggerTime,
		RecoverTime:  &now,
		Duration:     &duration,
		AlarmMessage: triggerRecord.AlarmMessage,
		CloudReported:   false,
		StorageLocation: StorageLocal,
		CreatedAt:       time.Now(),
	}

	// 如果有确认信息，复制过来
	if triggerRecord.AckTime != nil {
		record.AckTime = triggerRecord.AckTime
		record.AckUser = triggerRecord.AckUser
	}

	// 异步写入
	return rm.writeRecord(record)
}

// RecordShielded 记录被屏蔽的报警
// 当报警被屏蔽规则阻止时调用
func (rm *RecordManager) RecordShielded(alarm *engine.ActiveAlarm) error {
	if !rm.running.Load() {
		return fmt.Errorf("记录管理器未启动")
	}

	// 创建屏蔽记录（使用备注字段标注被屏蔽）
	record := &AlarmRecord{
		RecordID:     rm.generateRecordID(),
		RuleID:       alarm.Rule.ID,
		RuleName:     alarm.Rule.Name,
		EventType:    EventTrigger,
		Level:        alarm.Rule.Level,
		State:        core.AlarmStateActive, // 仍然记录为激活状态
		Category:     alarm.Rule.Category,
		TriggerTime:  alarm.TriggerTime,
		AlarmMessage: "[已屏蔽] " + rm.renderAlarmMessage(alarm),
		Threshold:    rm.extractThreshold(alarm),
		TriggerValue: alarm.TriggerValue,
		CloudReported:   false,
		StorageLocation: StorageLocal,
		CreatedAt:       time.Now(),
	}

	return rm.writeRecord(record)
}

// Query 查询报警记录
func (rm *RecordManager) Query(query *RecordQuery) (*RecordQueryResult, error) {
	if query == nil {
		query = NewRecordQuery()
	}

	// 先从缓存查询（如果查询条件简单）
	if rm.isSimpleQuery(query) {
		records := rm.queryFromCache(query)
		if len(records) > 0 {
			return rm.buildQueryResult(records, query), nil
		}
	}

	// 从存储层查询
	records, err := rm.storage.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}

	return rm.buildQueryResult(records, query), nil
}

// Export 导出报警记录
func (rm *RecordManager) Export(options *RecordExportOptions) error {
	if options == nil {
		return fmt.Errorf("导出选项不能为空")
	}

	// 构建查询条件
	query := NewRecordQuery()
	if options.StartTime != nil {
		query.StartTime = options.StartTime
	}
	if options.EndTime != nil {
		query.EndTime = options.EndTime
	}
	query.Limit = 1000000 // 导出时不限制数量

	// 查询记录
	result, err := rm.Query(query)
	if err != nil {
		return fmt.Errorf("查询记录失败: %w", err)
	}

	// 根据格式导出
	switch options.Format {
	case ExportFormatJSON:
		return rm.exportToJSON(result.Records, options.OutputPath)
	case ExportFormatCSV:
		return rm.exportToCSV(result.Records, options)
	default:
		return fmt.Errorf("不支持的导出格式: %d", options.Format)
	}
}

// GetStats 获取统计信息
func (rm *RecordManager) GetStats() (*RecordStats, error) {
	// 从存储层获取最新统计
	storageStats, err := rm.storage.GetStats()
	if err != nil {
		return nil, fmt.Errorf("获取统计信息失败: %w", err)
	}

	return storageStats, nil
}

// CleanupOldRecords 清理过期记录
func (rm *RecordManager) CleanupOldRecords(maxDays int) error {
	cutoffTime := time.Now().AddDate(0, 0, -maxDays)
	return rm.storage.CleanupBefore(cutoffTime)
}

// ============ 内部方法 ============

// startBackgroundWorkers 启动后台处理协程
func (rm *RecordManager) startBackgroundWorkers() {
	// 启动批量写入协程
	rm.wg.Add(1)
	go rm.batchWriter()
}

// batchWriter 批量写入协程
func (rm *RecordManager) batchWriter() {
	defer rm.wg.Done()

	for {
		select {
		case record := <-rm.pendingChan:
			rm.batchMu.Lock()
			rm.batch = append(rm.batch, record)

			// 达到批次大小，立即写入
			if len(rm.batch) >= rm.config.BatchSize {
				rm.flushBatchLocked()
			}
			rm.batchMu.Unlock()

		case <-rm.batchTimer.C:
			// 超时，写入当前批次
			rm.batchMu.Lock()
			if len(rm.batch) > 0 {
				rm.flushBatchLocked()
			}
			rm.batchTimer.Reset(rm.config.BatchTimeout)
			rm.batchMu.Unlock()

		case <-rm.ctx.Done():
			// 退出前刷新剩余批次
			rm.batchMu.Lock()
			rm.flushBatchLocked()
			rm.batchMu.Unlock()
			return
		}
	}
}

// flushBatchLocked 刷新批次（需要持有 batchMu）
func (rm *RecordManager) flushBatchLocked() {
	if len(rm.batch) == 0 {
		return
	}

	// 批量保存到存储
	for _, record := range rm.batch {
		if err := rm.storage.Save(record); err != nil {
			// 记录错误，但不中断批次
			fmt.Printf("[错误] 保存记录失败: %v\n", err)
		} else {
			// 更新缓存
			rm.addToCache(record)
			// 更新统计
			rm.updateStats(record)
		}
	}

	// 清空批次
	rm.batch = rm.batch[:0]
}

// flushBatch 刷新批次（公开方法，用于Stop时调用）
func (rm *RecordManager) flushBatch() {
	rm.batchMu.Lock()
	defer rm.batchMu.Unlock()
	rm.flushBatchLocked()
}

// writeRecord 写入单条记录（异步）
func (rm *RecordManager) writeRecord(record *AlarmRecord) error {
	// 记录最后触发时间
	if record.EventType == EventTrigger {
		rm.lastTriggerTime.Store(record.RuleID, time.Now())
	}

	// 非阻塞写入
	select {
	case rm.pendingChan <- record:
		return nil
	default:
		// 队列已满，同步写入
		if err := rm.storage.Save(record); err != nil {
			return fmt.Errorf("同步写入失败: %w", err)
		}
		rm.addToCache(record)
		rm.updateStats(record)
		return nil
	}
}

// addToCache 添加到缓存
func (rm *RecordManager) addToCache(record *AlarmRecord) {
	rm.cacheMu.Lock()
	defer rm.cacheMu.Unlock()

	// 环形缓存：覆盖旧记录
	rm.cache[rm.cacheIndex] = record
	rm.cacheIndex = (rm.cacheIndex + 1) % rm.config.CacheSize
}

// queryFromCache 从缓存查询
func (rm *RecordManager) queryFromCache(query *RecordQuery) []*AlarmRecord {
	rm.cacheMu.RLock()
	defer rm.cacheMu.RUnlock()

	var result []*AlarmRecord

	for _, record := range rm.cache {
		if record == nil {
			continue
		}

		if rm.matchQuery(record, query) {
			result = append(result, record)
		}
	}

	return result
}

// isSimpleQuery 判断是否为简单查询（可从缓存满足）
func (rm *RecordManager) isSimpleQuery(query *RecordQuery) bool {
	// 简单查询：只按规则ID查询，且数量<100
	return len(query.RuleIDs) > 0 &&
		len(query.Levels) == 0 &&
		len(query.Categories) == 0 &&
		query.StartTime == nil &&
		query.EndTime == nil &&
		query.Limit <= 100
}

// matchQuery 判断记录是否匹配查询条件
func (rm *RecordManager) matchQuery(record *AlarmRecord, query *RecordQuery) bool {
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

// buildQueryResult 构建查询结果
func (rm *RecordManager) buildQueryResult(records []*AlarmRecord, query *RecordQuery) *RecordQueryResult {
	// 排序
	rm.sortRecords(records, query.SortBy, query.SortDesc)

	// 分页
	total := len(records)
	offset := query.Offset
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		return &RecordQueryResult{
			Records: []*AlarmRecord{},
			Total:   total,
			Offset:  offset,
			Limit:   query.Limit,
			HasMore: false,
		}
	}

	end := offset + query.Limit
	if end > total {
		end = total
	}

	return &RecordQueryResult{
		Records: records[offset:end],
		Total:   total,
		Offset:  offset,
		Limit:   query.Limit,
		HasMore: end < total,
	}
}

// sortRecords 排序记录
func (rm *RecordManager) sortRecords(records []*AlarmRecord, sortBy string, sortDesc bool) {
	// 简单实现：按触发时间排序
	// TODO: 支持更多排序字段
}

// findLatestTriggerRecord 查找最新的触发记录
func (rm *RecordManager) findLatestTriggerRecord(ruleID string) (*AlarmRecord, error) {
	query := NewRecordQuery().
		WithRuleIDs(ruleID).
		WithEventTypes(EventTrigger).
		WithSort("trigger_time", true).
		WithLimit(1)

	result, err := rm.Query(query)
	if err != nil {
		return nil, err
	}

	if len(result.Records) == 0 {
		return nil, fmt.Errorf("未找到规则 %s 的触发记录", ruleID)
	}

	return result.Records[0], nil
}

// updateStats 更新统计信息
func (rm *RecordManager) updateStats(record *AlarmRecord) {
	if record.EventType == EventTrigger {
		atomic.AddInt64(&rm.stats.TriggerCount, 1)
		atomic.AddInt64(&rm.stats.TotalRecords, 1)
	} else if record.EventType == EventRecover {
		atomic.AddInt64(&rm.stats.RecoverCount, 1)
		atomic.AddInt64(&rm.stats.TotalRecords, 1)
	}

	if record.IsActive() {
		atomic.AddInt64(&rm.stats.ActiveCount, 1)
	}
}

// cleanupLoop 清理循环
func (rm *RecordManager) cleanupLoop() {
	defer rm.wg.Done()

	ticker := time.NewTicker(rm.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if rm.config.MaxDays > 0 {
				if err := rm.CleanupOldRecords(rm.config.MaxDays); err != nil {
					fmt.Printf("[错误] 清理过期记录失败: %v\n", err)
				}
			}
		case <-rm.ctx.Done():
			return
		}
	}
}

// generateRecordID 生成记录ID
func (rm *RecordManager) generateRecordID() string {
	// 格式: REC_<timestamp>_<counter>
	// 例如: REC_20250115150405_0001
	nano := time.Now().UnixNano()
	counter := atomic.AddInt64(&rm.stats.TotalRecords, 1)
	return fmt.Sprintf("REC_%d_%04d", nano, counter%10000)
}

// renderAlarmMessage 渲染报警消息（简单实现）
func (rm *RecordManager) renderAlarmMessage(alarm *engine.ActiveAlarm) string {
	// TODO: 集成 content 模块的完整渲染逻辑
	return alarm.Rule.Name
}

// extractThreshold 提取阈值
func (rm *RecordManager) extractThreshold(alarm *engine.ActiveAlarm) interface{} {
	// TODO: 从条件中提取阈值
	return nil
}

// extractTriggerVariables 提取触发变量信息
func (rm *RecordManager) extractTriggerVariables(alarm *engine.ActiveAlarm) []TriggerVariableInfo {
	// TODO: 从报警中提取变量信息
	return nil
}
