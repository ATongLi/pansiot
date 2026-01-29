package collector

import "time"

// Config 采集器配置
type Config struct {
	// MaxConcurrentTasks 最大并发任务数
	// 限制同时运行的任务数量，避免资源耗尽
	MaxConcurrentTasks int

	// DefaultTimeout 默认超时时间
	// 单次采集操作的超时时间
	DefaultTimeout time.Duration

	// EnableStatistics 是否启用统计
	// 启用后会记录采集次数、成功率、耗时等统计信息
	EnableStatistics bool

	// TaskStartupDelay 任务启动延迟
	// 采集器启动后，延迟多久再启动各个任务（避免同时启动造成压力）
	TaskStartupDelay time.Duration

	// StatisticsUpdateInterval 统计信息更新间隔
	// 统计信息的更新频率
	StatisticsUpdateInterval time.Duration
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		MaxConcurrentTasks:       100,               // 默认支持100个并发任务
		DefaultTimeout:           30 * time.Second,  // 默认30秒超时
		EnableStatistics:         true,              // 默认启用统计
		TaskStartupDelay:         100 * time.Millisecond, // 任务启动延迟100ms
		StatisticsUpdateInterval: 1 * time.Second,   // 每秒更新统计信息
	}
}

// HighPerformanceConfig 高性能配置
// 适用于高频采集场景（100ms级别）
func HighPerformanceConfig() Config {
	return Config{
		MaxConcurrentTasks:       500,               // 支持更多并发任务
		DefaultTimeout:           5 * time.Second,   // 更短的超时时间
		EnableStatistics:         true,
		TaskStartupDelay:         10 * time.Millisecond, // 更短的启动延迟
		StatisticsUpdateInterval: 5 * time.Second,   // 降低统计更新频率
	}
}

// LowResourceConfig 低资源配置
// 适用于资源受限的环境
func LowResourceConfig() Config {
	return Config{
		MaxConcurrentTasks:       20,                // 限制并发任务数
		DefaultTimeout:           60 * time.Second,  // 更长的超时时间
		EnableStatistics:         false,             // 禁用统计以节省资源
		TaskStartupDelay:         500 * time.Millisecond, // 更长的启动延迟
		StatisticsUpdateInterval: 10 * time.Second,
	}
}
