package collector

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"pansiot-device/internal/adapter"
	"pansiot-device/internal/core"
)

// Collector 数据采集器
// 负责按配置频率调度协议适配器，批量采集设备数据并写入实时存储层
type Collector struct {
	mu             sync.RWMutex
	tasks          map[string]*core.CollectionTask // 任务列表
	taskRunners    map[string]*TaskRunner         // 运行中的任务
	adapterFactory  *adapter.AdapterFactory        // 适配器工厂
	storage        core.Storage                   // 实时存储层
	config         Config                          // 配置
	stopChan       chan struct{}
	running        atomic.Bool
	wg             sync.WaitGroup
	stats          CollectorStats                  // 统计信息
}

// CollectorStats 采集器统计信息
type CollectorStats struct {
	TotalCollections int64         // 总采集次数
	SuccessCount     int64         // 成功次数
	FailureCount     int64         // 失败次数
	LastCollectTime  time.Time     // 最后采集时间
	AvgDuration      time.Duration // 平均采集耗时
}

// NewCollector 创建采集器
func NewCollector(adapterFactory *adapter.AdapterFactory, storage core.Storage) *Collector {
	return &Collector{
		tasks:         make(map[string]*core.CollectionTask),
		taskRunners:   make(map[string]*TaskRunner),
		adapterFactory: adapterFactory,
		storage:       storage,
		config:        DefaultConfig(),
		stopChan:      make(chan struct{}),
		stats:         CollectorStats{},
	}
}

// Start 启动采集器
func (c *Collector) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running.Load() {
		return fmt.Errorf("采集器已在运行")
	}

	c.running.Store(true)
	log.Printf("[Collector] 采集器已启动")

	return nil
}

// Stop 停止采集器
func (c *Collector) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running.Load() {
		return fmt.Errorf("采集器未在运行")
	}

	log.Printf("[Collector] 正在停止采集器...")
	c.running.Store(false)
	close(c.stopChan)

	// 停止所有任务
	for _, runner := range c.taskRunners {
		runner.Stop()
	}

	// 等待所有任务停止
	c.wg.Wait()

	log.Printf("[Collector] 采集器已停止")
	return nil
}

// AddTask 添加采集任务
func (c *Collector) AddTask(task *core.CollectionTask) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.tasks[task.ID]; exists {
		return fmt.Errorf("任务已存在: %s", task.ID)
	}

	// 验证任务配置
	if err := c.validateTask(task); err != nil {
		return fmt.Errorf("任务配置无效: %v", err)
	}

	c.tasks[task.ID] = task
	log.Printf("[Collector] 已添加任务: %s (频率: %dms, 变量数: %d)",
		task.ID, task.Frequency, len(task.VariableIDs))

	// 如果采集器正在运行，立即启动任务
	if c.running.Load() {
		return c.startTask(task)
	}

	return nil
}

// RemoveTask 移除采集任务
func (c *Collector) RemoveTask(taskID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.tasks[taskID]; !exists {
		return fmt.Errorf("任务不存在: %s", taskID)
	}

	// 停止任务运行器
	if runner, exists := c.taskRunners[taskID]; exists {
		runner.Stop()
		delete(c.taskRunners, taskID)
	}

	delete(c.tasks, taskID)
	log.Printf("[Collector] 已移除任务: %s", taskID)

	return nil
}

// UpdateTask 更新采集任务
func (c *Collector) UpdateTask(task *core.CollectionTask) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.tasks[task.ID]; !exists {
		return fmt.Errorf("任务不存在: %s", task.ID)
	}

	// 验证任务配置
	if err := c.validateTask(task); err != nil {
		return fmt.Errorf("任务配置无效: %v", err)
	}

	// 如果任务正在运行，需要重启
	if runner, exists := c.taskRunners[task.ID]; exists {
		runner.Stop()
		delete(c.taskRunners, task.ID)
	}

	c.tasks[task.ID] = task
	log.Printf("[Collector] 已更新任务: %s", task.ID)

	// 如果采集器正在运行，重新启动任务
	if c.running.Load() && task.Enable {
		return c.startTask(task)
	}

	return nil
}

// GetTask 获取采集任务
func (c *Collector) GetTask(taskID string) (*core.CollectionTask, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	task, exists := c.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("任务不存在: %s", taskID)
	}

	return task, nil
}

// ListTasks 列出所有采集任务
func (c *Collector) ListTasks() []*core.CollectionTask {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tasks := make([]*core.CollectionTask, 0, len(c.tasks))
	for _, task := range c.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// GetStats 获取采集器统计信息
func (c *Collector) GetStats() CollectorStats {
	return c.stats
}

// validateTask 验证任务配置
func (c *Collector) validateTask(task *core.CollectionTask) error {
	if task.ID == "" {
		return fmt.Errorf("任务ID不能为空")
	}
	if task.Frequency <= 0 {
		return fmt.Errorf("采集频率必须大于0")
	}
	if task.DeviceID == "" {
		return fmt.Errorf("设备ID不能为空")
	}
	if task.ProtocolType == "" {
		return fmt.Errorf("协议类型不能为空")
	}
	if len(task.VariableIDs) == 0 {
		return fmt.Errorf("变量ID列表不能为空")
	}
	if task.Priority < 1 || task.Priority > 10 {
		return fmt.Errorf("优先级必须在1-10之间")
	}
	if task.Timeout <= 0 {
		return fmt.Errorf("超时时间必须大于0")
	}

	return nil
}

// startTask 启动任务
func (c *Collector) startTask(task *core.CollectionTask) error {
	// 检查是否已存在运行器
	if _, exists := c.taskRunners[task.ID]; exists {
		return fmt.Errorf("任务已在运行: %s", task.ID)
	}

	// 创建设备配置
	device := &core.Device{
		ID:       task.DeviceID,
		Protocol: task.ProtocolType,
	}

	// 通过工厂创建适配器
	protocolAdapter, err := c.adapterFactory.Create(device)
	if err != nil {
		return fmt.Errorf("创建协议适配器失败: %v", err)
	}

	// 连接设备
	ctx := context.Background()
	if err := protocolAdapter.Connect(ctx, device); err != nil {
		return fmt.Errorf("连接设备失败: %v", err)
	}

	// 创建任务运行器
	runner := NewTaskRunner(task, protocolAdapter, c.storage, &c.stats)
	c.taskRunners[task.ID] = runner

	// 启动任务
	if task.Enable {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			runner.Start(c.stopChan)
		}()

		log.Printf("[Collector] 已启动任务: %s", task.ID)
	}

	return nil
}
