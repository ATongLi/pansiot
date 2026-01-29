package collector

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"pansiot-device/internal/core"
)

// TaskRunner 任务运行器
// 负责按配置频率执行单个采集任务
type TaskRunner struct {
	task    *core.CollectionTask // 任务配置
	adapter core.ProtocolAdapter // 协议适配器
	storage core.Storage         // 存储层
	stats   *CollectorStats      // 统计信息指针
	ticker  *time.Ticker         // 定时器
	stopChan chan struct{}       // 停止信号
	running atomic.Bool          // 运行状态
	mu      sync.RWMutex         // 保护内部状态
}

// NewTaskRunner 创建任务运行器
func NewTaskRunner(
	task *core.CollectionTask,
	adapter core.ProtocolAdapter,
	storage core.Storage,
	stats *CollectorStats,
) *TaskRunner {
	return &TaskRunner{
		task:     task,
		adapter:  adapter,
		storage:  storage,
		stats:    stats,
		stopChan: make(chan struct{}),
	}
}

// Start 启动任务
func (tr *TaskRunner) Stop() {
	if !tr.running.Load() {
		return
	}

	tr.running.Store(false)
	close(tr.stopChan)

	tr.mu.Lock()
	defer tr.mu.Unlock()

	if tr.ticker != nil {
		tr.ticker.Stop()
	}

	log.Printf("[TaskRunner] 任务已停止: %s", tr.task.ID)
}

// Start 启动任务
func (tr *TaskRunner) Start(globalStopChan chan struct{}) {
	if tr.running.Swap(true) {
		log.Printf("[TaskRunner] 任务已在运行: %s", tr.task.ID)
		return
	}

	// 创建定时器
	frequency := time.Duration(tr.task.Frequency) * time.Millisecond
	tr.ticker = time.NewTicker(frequency)
	defer tr.ticker.Stop()

	log.Printf("[TaskRunner] 任务已启动: %s (频率: %dms, 变量数: %d)",
		tr.task.ID, tr.task.Frequency, len(tr.task.VariableIDs))

	// 立即执行一次采集
	tr.collect()

	// 定时循环
	for {
		select {
		case <-tr.ticker.C:
			tr.collect()
		case <-tr.stopChan:
			log.Printf("[TaskRunner] 收到停止信号: %s", tr.task.ID)
			return
		case <-globalStopChan:
			log.Printf("[TaskRunner] 收到全局停止信号: %s", tr.task.ID)
			return
		}
	}
}

// collect 执行数据采集
func (tr *TaskRunner) collect() {
	startTime := time.Now()

	// 更新统计信息
	atomic.AddInt64(&tr.stats.TotalCollections, 1)

	// 创建超时上下文
	timeout := time.Duration(tr.task.Timeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 批量读取变量
	variables, err := tr.adapter.ReadVariables(ctx, tr.task.VariableIDs)
	if err != nil {
		atomic.AddInt64(&tr.stats.FailureCount, 1)
		log.Printf("[TaskRunner] 采集失败 [%s]: %v", tr.task.ID, err)
		return
	}

	// 验证返回的变量数量
	if len(variables) != len(tr.task.VariableIDs) {
		atomic.AddInt64(&tr.stats.FailureCount, 1)
		log.Printf("[TaskRunner] 采集变量数量不匹配 [%s]: 期望 %d, 实际 %d",
			tr.task.ID, len(tr.task.VariableIDs), len(variables))
		return
	}

	// 批量写入存储层
	writeErrors := 0
	for _, variable := range variables {
		if err := tr.storage.WriteVar(variable); err != nil {
			writeErrors++
			log.Printf("[TaskRunner] 写入变量失败 [%s] ID=%d: %v",
				tr.task.ID, variable.ID, err)
		}
	}

	// 更新统计
	if writeErrors > 0 {
		atomic.AddInt64(&tr.stats.FailureCount, 1)
		log.Printf("[TaskRunner] 采集部分失败 [%s]: %d/%d 变量写入失败",
			tr.task.ID, writeErrors, len(variables))
	} else {
		atomic.AddInt64(&tr.stats.SuccessCount, 1)
	}

	// 更新耗时统计
	elapsed := time.Since(startTime)
	tr.mu.Lock()
	if tr.stats.TotalCollections == 1 {
		tr.stats.AvgDuration = elapsed
	} else {
		// 简单移动平均
		tr.stats.AvgDuration = (tr.stats.AvgDuration + elapsed) / 2
	}
	tr.stats.LastCollectTime = startTime
	tr.mu.Unlock()

	// 记录成功日志
	if writeErrors == 0 {
		log.Printf("[TaskRunner] 采集成功 [%s]: %d个变量, 耗时: %v",
			tr.task.ID, len(variables), elapsed)
	}
}

// IsRunning 检查任务是否在运行
func (tr *TaskRunner) IsRunning() bool {
	return tr.running.Load()
}

// GetTask 获取任务配置
func (tr *TaskRunner) GetTask() *core.CollectionTask {
	return tr.task
}
