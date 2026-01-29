package record

import (
	"encoding/json"
	"time"

	"pansiot-device/internal/core"
)

// AlarmRecord 报警记录
type AlarmRecord struct {
	// 基础信息
	RecordID string // 记录ID（全局唯一）
	RuleID   string // 规则ID
	RuleName string // 规则名称

	// 报警事件信息
	EventType EventType       // 事件类型
	Level     core.AlarmLevel // 报警级别
	State     core.AlarmState // 报警状态
	Category  string          // 报警类别

	// 时间信息
	TriggerTime  time.Time      // 触发时间
	AckTime      *time.Time     // 确认时间
	RecoverTime  *time.Time     // 恢复时间
	Duration     *time.Duration // 持续时间（恢复时计算）

	// 内容信息
	AlarmMessage string // 报警内容（已渲染）

	// 数值信息
	Threshold     interface{} // 阈值
	TriggerValue  interface{} // 触发值

	// 责任人信息
	ResponsibleUsers []string // 责任人用户ID列表
	AckUser          string    // 确认用户

	// 变量信息
	TriggerVariables []TriggerVariableInfo // 触发变量信息

	// 云上报信息
	CloudReported   bool        // 是否已上报云端
	CloudReportTime *time.Time  // 云上报时间

	// 存储信息
	StorageLocation StorageLocation // 存储位置
	CreatedAt       time.Time        // 记录创建时间
}

// EventType 事件类型
type EventType int

const (
	EventTrigger  EventType = iota // 报警触发
	EventRecover                   // 报警恢复
)

// String 返回事件类型的字符串表示
func (et EventType) String() string {
	switch et {
	case EventTrigger:
		return "trigger"
	case EventRecover:
		return "recover"
	default:
		return "unknown"
	}
}

// StorageLocation 存储位置
type StorageLocation int

const (
	StorageLocal  StorageLocation = iota // 本设备
	StorageUSB                            // U盘
)

// String 返回存储位置的字符串表示
func (sl StorageLocation) String() string {
	switch sl {
	case StorageLocal:
		return "local"
	case StorageUSB:
		return "usb"
	default:
		return "unknown"
	}
}

// TriggerVariableInfo 触发变量信息
type TriggerVariableInfo struct {
	VariableID   uint64           // 变量ID
	VariableName string           // 变量名称
	Value        interface{}      // 触发时的值
	Quality      core.QualityCode // 质量码
}

// ToJSON 转换为JSON字符串
func (r *AlarmRecord) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON 从JSON字符串解析
func (r *AlarmRecord) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

// GetDurationString 获取持续时间的字符串表示
func (r *AlarmRecord) GetDurationString() string {
	if r.Duration == nil {
		return "-"
	}
	return r.Duration.String()
}

// IsTriggerEvent 是否为触发事件
func (r *AlarmRecord) IsTriggerEvent() bool {
	return r.EventType == EventTrigger
}

// IsRecoverEvent 是否为恢复事件
func (r *AlarmRecord) IsRecoverEvent() bool {
	return r.EventType == EventRecover
}

// IsActive 报警是否处于激活状态
func (r *AlarmRecord) IsActive() bool {
	return r.State == core.AlarmStateActive || r.State == core.AlarmStateAcknowledged
}

// IsAcknowledged 报警是否已确认
func (r *AlarmRecord) IsAcknowledged() bool {
	return r.State == core.AlarmStateAcknowledged
}

// IsCleared 报警是否已清除
func (r *AlarmRecord) IsCleared() bool {
	return r.State == core.AlarmStateCleared
}

// GetResponsibleUsersJSON 获取责任人列表的JSON表示
func (r *AlarmRecord) GetResponsibleUsersJSON() string {
	if r.ResponsibleUsers == nil {
		return "[]"
	}
	data, _ := json.Marshal(r.ResponsibleUsers)
	return string(data)
}

// GetTriggerVariablesJSON 获取触发变量列表的JSON表示
func (r *AlarmRecord) GetTriggerVariablesJSON() string {
	if r.TriggerVariables == nil {
		return "[]"
	}
	data, _ := json.Marshal(r.TriggerVariables)
	return string(data)
}

// RecordQuery 报警记录查询条件
type RecordQuery struct {
	StartTime *time.Time        // 开始时间
	EndTime   *time.Time        // 结束时间
	RuleIDs   []string          // 规则ID列表
	Levels    []core.AlarmLevel // 报警级别列表
	Categories []string         // 报警类别列表
	States    []core.AlarmState // 报警状态列表
	EventTypes []EventType      // 事件类型列表
	Limit     int               // 限制数量
	Offset    int               // 偏移量
	SortBy    string            // 排序字段
	SortDesc  bool              // 是否降序
}

// RecordQueryResult 报警记录查询结果
type RecordQueryResult struct {
	Records  []*AlarmRecord // 记录列表
	Total    int            // 总记录数
	Offset   int            // 当前偏移量
	Limit    int            // 当前限制数量
	HasMore  bool           // 是否有更多记录
}

// NewRecordQuery 创建新的查询条件
func NewRecordQuery() *RecordQuery {
	return &RecordQuery{
		Limit:    100,  // 默认限制100条
		Offset:   0,    // 默认从0开始
		SortBy:   "trigger_time", // 默认按触发时间排序
		SortDesc: true, // 默认降序
	}
}

// WithStartTime 设置开始时间
func (q *RecordQuery) WithStartTime(startTime time.Time) *RecordQuery {
	q.StartTime = &startTime
	return q
}

// WithEndTime 设置结束时间
func (q *RecordQuery) WithEndTime(endTime time.Time) *RecordQuery {
	q.EndTime = &endTime
	return q
}

// WithRuleIDs 设置规则ID列表
func (q *RecordQuery) WithRuleIDs(ruleIDs ...string) *RecordQuery {
	q.RuleIDs = ruleIDs
	return q
}

// WithLevels 设置报警级别列表
func (q *RecordQuery) WithLevels(levels ...core.AlarmLevel) *RecordQuery {
	q.Levels = levels
	return q
}

// WithCategories 设置报警类别列表
func (q *RecordQuery) WithCategories(categories ...string) *RecordQuery {
	q.Categories = categories
	return q
}

// WithStates 设置报警状态列表
func (q *RecordQuery) WithStates(states ...core.AlarmState) *RecordQuery {
	q.States = states
	return q
}

// WithEventTypes 设置事件类型列表
func (q *RecordQuery) WithEventTypes(eventTypes ...EventType) *RecordQuery {
	q.EventTypes = eventTypes
	return q
}

// WithLimit 设置限制数量
func (q *RecordQuery) WithLimit(limit int) *RecordQuery {
	q.Limit = limit
	return q
}

// WithOffset 设置偏移量
func (q *RecordQuery) WithOffset(offset int) *RecordQuery {
	q.Offset = offset
	return q
}

// WithSort 设置排序
func (q *RecordQuery) WithSort(sortBy string, sortDesc bool) *RecordQuery {
	q.SortBy = sortBy
	q.SortDesc = sortDesc
	return q
}

// RecordExportOptions 报警记录导出选项
type RecordExportOptions struct {
	Format        ExportFormat // 导出格式
	StartTime     *time.Time   // 开始时间
	EndTime       *time.Time   // 结束时间
	IncludeHeader bool         // 是否包含表头
	OutputPath    string       // 输出路径
}

// ExportFormat 导出格式
type ExportFormat int

const (
	ExportFormatJSON ExportFormat = iota // JSON格式
	ExportFormatCSV                      // CSV格式
	ExportFormatExcel                    // Excel格式（暂未实现）
)

// String 返回导出格式的字符串表示
func (ef ExportFormat) String() string {
	switch ef {
	case ExportFormatJSON:
		return "json"
	case ExportFormatCSV:
		return "csv"
	case ExportFormatExcel:
		return "excel"
	default:
		return "unknown"
	}
}

// RecordStats 报警记录统计信息
type RecordStats struct {
	TotalRecords    int64            // 总记录数
	TriggerCount    int64            // 触发事件数
	RecoverCount    int64            // 恢复事件数
	ActiveCount     int64            // 当前激活的报警数
	RecordsByLevel  map[core.AlarmLevel]int64 // 按级别统计
	RecordsByCategory map[string]int64        // 按类别统计
	RecordsByRule   map[string]int64          // 按规则统计
	FirstRecordTime *time.Time       // 最早记录时间
	LastRecordTime  *time.Time       // 最晚记录时间
}

// NewRecordStats 创建新的统计信息
func NewRecordStats() *RecordStats {
	return &RecordStats{
		RecordsByLevel:   make(map[core.AlarmLevel]int64),
		RecordsByCategory: make(map[string]int64),
		RecordsByRule:    make(map[string]int64),
	}
}

// IncrementLevel 增加指定级别的计数
func (s *RecordStats) IncrementLevel(level core.AlarmLevel) {
	s.RecordsByLevel[level]++
}

// IncrementCategory 增加指定类别的计数
func (s *RecordStats) IncrementCategory(category string) {
	s.RecordsByCategory[category]++
}

// IncrementRule 增加指定规则的计数
func (s *RecordStats) IncrementRule(ruleID string) {
	s.RecordsByRule[ruleID]++
}

// GetLevelCount 获取指定级别的计数
func (s *RecordStats) GetLevelCount(level core.AlarmLevel) int64 {
	return s.RecordsByLevel[level]
}

// GetCategoryCount 获取指定类别的计数
func (s *RecordStats) GetCategoryCount(category string) int64 {
	return s.RecordsByCategory[category]
}

// GetRuleCount 获取指定规则的计数
func (s *RecordStats) GetRuleCount(ruleID string) int64 {
	return s.RecordsByRule[ruleID]
}
