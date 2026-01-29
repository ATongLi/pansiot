package script

import (
	"time"
)

// Script 脚本定义
type Script struct {
	ID          string        // 脚本唯一标识
	Name        string        // 脚本名称
	Content     string        // 脚本内容（JavaScript 代码）
	Description string        // 脚本描述
	Enabled     bool          // 是否启用
	Triggers    []ScriptTrigger // 触发器配置
	Variables   []string      // 依赖的变量 ID（用于订阅）
	Timeout     time.Duration // 执行超时（默认 5s）
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ScriptTrigger 脚本触发器
type ScriptTrigger struct {
	ID     string      // 触发器唯一标识
	Type   TriggerType // 触发类型
	Config interface{} // 触发配置（根据类型不同而不同）
	Enabled bool       // 是否启用
}

// TriggerType 触发器类型
type TriggerType int

const (
	TriggerTypeVariable  TriggerType = iota // 变量变化触发
	TriggerTypePeriodic                     // 周期执行
	TriggerTypeSystem                       // 系统触发
	TriggerTypeAlarm                        // 报警触发
	TriggerTypeUI                           // UI 触发（按钮点击等）
)

// VariableTriggerConfig 变量触发配置
type VariableTriggerConfig struct {
	VariableID uint64      // 变量 ID
	Condition  string      // 条件：">", "<", "==", "!=", ">=", "<="
	Threshold  interface{} // 阈值
	EdgeOnly   bool        // 仅边沿触发（true=只触发一次）
}

// PeriodicTriggerConfig 周期触发配置
type PeriodicTriggerConfig struct {
	Interval    time.Duration // 执行间隔（固定间隔模式）
	CronExpr    string        // Cron表达式（Cron模式，格式：秒 分 时 日 月 周）
	StartTime   string        // 开始时间（格式："HH:MM:SS"，用于时间窗口）
	EndTime     string        // 结束时间（格式："HH:MM:SS"，用于时间窗口）
	DaysOfWeek  []int         // 星期几（1=周一，...，7=周日，用于时间窗口）
	TimeZone    string        // 时区（如："Asia/Shanghai"，默认为系统本地时区）
}

// SystemTriggerConfig 系统触发配置
type SystemTriggerConfig struct {
	Event SystemEvent // 系统事件类型
}

// SystemEvent 系统事件
type SystemEvent int

const (
	SystemEventStartup  SystemEvent = iota // 开机启动
	SystemEventShutdown                    // 关机前
)

// AlarmTriggerConfig 报警触发配置
type AlarmTriggerConfig struct {
	RuleID    string         // 报警规则 ID
	EventType AlarmEventType // 事件类型
}

// AlarmEventType 报警事件类型
type AlarmEventType int

const (
	AlarmEventTrigger  AlarmEventType = iota // 报警触发
	AlarmEventRecover                        // 报警恢复
	AlarmEventAcknowledge                    // 报警确认
)

// ScriptConfig 脚本消费者配置
type ScriptConfig struct {
	// VM 池配置
	VMPoolSize    int           // VM 池大小（默认 10）
	VMMaxIdle     time.Duration // VM 最大空闲时间（默认 5m）
	VMMaxLifetime time.Duration // VM 最大生命周期（默认 30m）

	// 执行配置
	DefaultTimeout time.Duration // 默认执行超时（默认 5s）
	MaxConcurrent  int           // 最大并发执行数（默认 100）
	QueueSize      int           // 执行队列大小（默认 1000）

	// 安全配置
	EnableSandbox bool  // 是否启用沙箱（默认 true）
	MemoryLimit   int64 // 内存限制（字节，默认 10MB）
	MaxExecutions int   // 单脚本最大执行次数/分钟（默认 60）

	// 日志配置
	EnableLog bool   // 是否启用脚本日志（默认 false）
	LogPath   string // 日志路径（默认 "./data/logs/script/"）
}

// DefaultScriptConfig 默认配置
func DefaultScriptConfig() *ScriptConfig {
	return &ScriptConfig{
		VMPoolSize:     10,
		VMMaxIdle:      5 * time.Minute,
		VMMaxLifetime:  30 * time.Minute,
		DefaultTimeout: 5 * time.Second,
		MaxConcurrent:  100,
		QueueSize:      1000,
		EnableSandbox:  true,
		MemoryLimit:    10 * 1024 * 1024, // 10MB
		MaxExecutions:  60,
		EnableLog:      false,
		LogPath:        "./data/logs/script/",
	}
}

// ScriptStatus 脚本状态
type ScriptStatus struct {
	ScriptID      string       // 脚本 ID
	Loaded        bool         // 是否已加载
	Enabled       bool         // 是否启用
	LastExecution time.Time    // 最后执行时间
	ExecCount     int64        // 执行次数
	ErrorCount    int64        // 错误次数
	LastError     string       // 最后错误信息
	State         ScriptState  // 当前状态
}

// ScriptState 脚本状态
type ScriptState int

const (
	ScriptStateIdle      ScriptState = iota // 空闲
	ScriptStateRunning                      // 运行中
	ScriptStateError                        // 错误
	ScriptStateDisabled                     // 已禁用
)
