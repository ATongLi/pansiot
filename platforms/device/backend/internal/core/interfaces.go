package core

import (
	"context"
	"time"
)

// ProtocolAdapter 协议适配器接口
// 所有协议适配器（Modbus, OPC UA, MQTT等）都必须实现此接口
type ProtocolAdapter interface {
	// Connect 连接到设备
	Connect(ctx context.Context, device *Device) error

	// Disconnect 断开设备连接
	Disconnect(ctx context.Context) error

	// IsConnected 检查连接状态
	IsConnected() bool

	// ReadVariable 读取单个变量
	ReadVariable(ctx context.Context, variableID uint64) (*Variable, error)

	// ReadVariables 批量读取变量（性能优化）
	ReadVariables(ctx context.Context, variableIDs []uint64) ([]*Variable, error)

	// WriteVariable 写入单个变量（可选，部分设备只读）
	WriteVariable(ctx context.Context, variableID uint64, value interface{}) error

	// GetProtocol 获取协议类型
	GetProtocol() string
}

// Storage 实时变量存储层接口
// 这是整个系统的核心接口，所有模块都通过此接口访问变量存储
type Storage interface {
	// ReadVar 读取单个变量
	ReadVar(variableID uint64) (*Variable, error)

	// ReadVars 批量读取变量
	ReadVars(variableIDs []uint64) ([]*Variable, error)

	// ReadVarByStringID 通过字符串ID读取变量
	ReadVarByStringID(stringID string) (*Variable, error)

	// WriteVar 写入单个变量
	WriteVar(variable *Variable) error

	// WriteVars 批量写入变量
	WriteVars(variables []*Variable) error

	// Subscribe 订阅变量更新
	Subscribe(subscriberID string, variableIDs []uint64, callback func(VariableUpdate)) error

	// SubscribeByDevice 按设备订阅（订阅该设备的所有变量）
	SubscribeByDevice(subscriberID string, deviceID string, callback func(VariableUpdate)) error

	// SubscribeByPattern 按模式订阅（支持通配符，如 "DV-PLC001-*"）
	SubscribeByPattern(subscriberID string, pattern string, callback func(VariableUpdate)) error

	// Unsubscribe 取消订阅
	Unsubscribe(subscriberID string, variableIDs []uint64) error

	// UnsubscribeAll 取消所有订阅
	UnsubscribeAll(subscriberID string) error

	// GetStats 获取存储统计信息
	GetStats() StorageStats

	// CreateVariable 创建新变量
	CreateVariable(variable *Variable) error

	// DeleteVariable 删除变量
	DeleteVariable(variableID uint64) error

	// ListVariables 列出所有变量
	ListVariables() []*Variable

	// ListVariablesByDevice 列出指定设备的所有变量
	ListVariablesByDevice(deviceID string) []*Variable
}

// StorageStats 存储统计信息
type StorageStats struct {
	TotalVariables     int       // 总变量数
	TotalSubscribers   int       // 总订阅者数
	TotalSubscriptions int       // 总订阅数
	ReadCount          int64     // 读取次数
	WriteCount         int64     // 写入次数
	StartTime          time.Time // 启动时间
}

// Collector 数据采集器接口
type Collector interface {
	// Start 启动采集器
	Start(ctx context.Context) error

	// Stop 停止采集器
	Stop() error

	// IsRunning 检查运行状态
	IsRunning() bool

	// AddTask 添加采集任务
	AddTask(task *CollectionTask) error

	// RemoveTask 移除采集任务
	RemoveTask(taskID string) error

	// UpdateTask 更新采集任务
	UpdateTask(task *CollectionTask) error

	// GetTask 获取采集任务
	GetTask(taskID string) (*CollectionTask, error)

	// ListTasks 列出所有采集任务
	ListTasks() []*CollectionTask

	// GetStats 获取采集统计信息
	GetStats() CollectorStats
}

// CollectorStats 采集器统计信息
type CollectorStats struct {
	RunningTasks     int           // 运行中的任务数
	TotalCollections int64         // 总采集次数
	SuccessCount     int64         // 成功次数
	FailureCount     int64         // 失败次数
	AvgDuration      time.Duration // 平均采集耗时
	LastCollectTime  time.Time     // 最后采集时间
}

// Consumer 消费者接口
// 所有消费模块（脚本、报警、前端通讯等）都实现此接口
type Consumer interface {
	// Start 启动消费者
	Start(ctx context.Context) error

	// Stop 停止消费者
	Stop() error

	// IsRunning 检查运行状态
	IsRunning() bool

	// GetID 获取消费者ID
	GetID() string

	// GetType 获取消费者类型
	GetType() string
}

// VariableChangeHandler 变量变化处理器接口
// 用于处理变量更新事件
type VariableChangeHandler interface {
	// OnVariableChange 变量更新时调用
	OnVariableChange(update VariableUpdate) error
}

// AlarmProcessor 报警处理器接口
type AlarmProcessor interface {
	Consumer
	// AddRule 添加报警规则
	AddRule(rule *AlarmRule) error

	// RemoveRule 移除报警规则
	RemoveRule(ruleID string) error

	// GetActiveAlarms 获取激活的报警列表
	GetActiveAlarms() []*AlarmEvent

	// AcknowledgeAlarm 确认报警
	AcknowledgeAlarm(alarmID string, user string) error

	// ClearAlarm 清除报警
	ClearAlarm(alarmID string) error
}

// ScriptEngine 脚本引擎接口
type ScriptEngine interface {
	Consumer
	// LoadScript 加载脚本
	LoadScript(scriptID, scriptContent string) error

	// UnloadScript 卸载脚本
	UnloadScript(scriptID string) error

	// Execute 执行脚本
	Execute(scriptID string, input map[string]interface{}) (map[string]interface{}, error)

	// GetScriptStatus 获取脚本状态
	GetScriptStatus(scriptID string) (*ScriptStatus, error)
}

// ScriptStatus 脚本状态
type ScriptStatus struct {
	ScriptID     string    // 脚本ID
	IsLoaded     bool      // 是否已加载
	IsEnabled    bool      // 是否启用
	LastExecute  time.Time // 最后执行时间
	ExecuteCount int64     // 执行次数
	ErrorCount   int64     // 错误次数
	LastError    string    // 最后错误信息
}

// HistoryStorage 历史数据存储接口
type HistoryStorage interface {
	Consumer
	// Store 存储历史数据
	Store(variableID uint64, value interface{}, timestamp time.Time) error

	// Query 查询历史数据
	Query(variableID uint64, startTime, endTime time.Time, limit int) ([]HistoryDataPoint, error)

	// QueryAggregate 聚合查询
	QueryAggregate(variableID uint64, startTime, endTime time.Time, aggregateType string) (*AggregateResult, error)
}

// HistoryDataPoint 历史数据点
type HistoryDataPoint struct {
	VariableID uint64
	Value      interface{}
	Timestamp  time.Time
	Quality    QualityCode
}

// AggregateResult 聚合结果
type AggregateResult struct {
	VariableID uint64
	StartTime  time.Time
	EndTime    time.Time
	Count      int64
	Min        interface{}
	Max        interface{}
	Avg        float64
	Sum        float64
}

// DataReporter 数据上报接口
type DataReporter interface {
	Consumer
	// Report 上报数据
	Report(data []*Variable) error

	// SetReportMode 设置上报模式
	SetReportMode(mode ReportMode) error

	// GetReportStatus 获取上报状态
	GetReportStatus() *ReportStatus
}

// ReportMode 上报模式
type ReportMode int

const (
	ReportModeRealtime ReportMode = iota // 实时上报
	ReportModeBatch                      // 批量上报
)

// ReportStatus 上报状态
type ReportStatus struct {
	Mode           ReportMode // 当前模式
	TotalReports   int64      // 总上报次数
	SuccessCount   int64      // 成功次数
	FailureCount   int64      // 失败次数
	LastReportTime time.Time  // 最后上报时间
	CacheSize      int        // 缓存队列大小
}

// FrontendCommunicator 前端通讯接口
type FrontendCommunicator interface {
	Consumer
	// HandleConnection 处理新的WebSocket连接
	HandleConnection(connID string, conn interface{}) error

	// CloseConnection 关闭连接
	CloseConnection(connID string) error

	// SendToConnection 发送数据到指定连接
	SendToConnection(connID string, data interface{}) error

	// Broadcast 广播数据到所有连接
	Broadcast(data interface{}) error

	// GetConnections 获取所有连接
	GetConnections() []string
}

// ConfigLoader 配置加载器接口
type ConfigLoader interface {
	// Load 加载配置
	Load(configPath string) error

	// Save 保存配置
	Save(configPath string) error

	// Validate 验证配置
	Validate() error

	// Get 获取配置
	Get(key string) (interface{}, error)

	// Set 设置配置
	Set(key string, value interface{}) error
}

// Logger 日志接口
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}
