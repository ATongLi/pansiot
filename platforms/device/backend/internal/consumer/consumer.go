package consumer

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"pansiot-device/internal/core"
)

// BaseConsumer 消费者基类
// 提供所有消费者的通用功能：生命周期管理、状态查询、统计信息等
type BaseConsumer struct {
	id           string              // 消费者唯一标识
	consumerType string              // 消费者类型：alarm, script, history, reporter, websocket
	storage      core.Storage        // 实时存储层引用
	running      atomic.Bool         // 运行状态
	stopChan     chan struct{}       // 停止信号通道
	wg           sync.WaitGroup      // 等待组，用于优雅关闭
	mu           sync.RWMutex        // 读写锁，保护内部状态
	config       interface{}         // 配置（具体类型由各消费者定义）
	stats        ConsumerStats       // 统计信息
}

// ConsumerStats 消费者统计信息
type ConsumerStats struct {
	TotalProcessed  int64     // 总处理次数
	SuccessCount    int64     // 成功次数
	FailureCount    int64     // 失败次数
	LastProcessTime time.Time // 最后处理时间
	StartTime       time.Time // 启动时间
}

// NewBaseConsumer 创建消费者基类实例
func NewBaseConsumer(id, consumerType string, storage core.Storage) *BaseConsumer {
	return &BaseConsumer{
		id:           id,
		consumerType: consumerType,
		storage:      storage,
		stopChan:     make(chan struct{}),
		stats:        ConsumerStats{},
	}
}

// Start 启动消费者
// 子类应该重写此方法以实现具体的启动逻辑
func (bc *BaseConsumer) Start(ctx context.Context) error {
	if !bc.running.CompareAndSwap(false, true) {
		return fmt.Errorf("消费者[%s]已在运行", bc.id)
	}

	bc.mu.Lock()
	bc.stats.StartTime = time.Now()
	bc.mu.Unlock()

	return nil
}

// Stop 停止消费者
// 子类应该重写此方法以实现具体的停止逻辑
func (bc *BaseConsumer) Stop() error {
	if !bc.running.CompareAndSwap(true, false) {
		return fmt.Errorf("消费者[%s]未在运行", bc.id)
	}

	close(bc.stopChan)
	bc.wg.Wait()

	// 重新创建 stopChan 以便下次启动
	bc.stopChan = make(chan struct{})

	return nil
}

// IsRunning 检查消费者是否在运行
func (bc *BaseConsumer) IsRunning() bool {
	return bc.running.Load()
}

// GetID 获取消费者ID
func (bc *BaseConsumer) GetID() string {
	return bc.id
}

// GetType 获取消费者类型
func (bc *BaseConsumer) GetType() string {
	return bc.consumerType
}

// GetStats 获取统计信息
func (bc *BaseConsumer) GetStats() ConsumerStats {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.stats
}

// GetStorage 获取存储层引用
func (bc *BaseConsumer) GetStorage() core.Storage {
	return bc.storage
}

// GetConfig 获取配置
func (bc *BaseConsumer) GetConfig() interface{} {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.config
}

// SetConfig 设置配置
func (bc *BaseConsumer) SetConfig(config interface{}) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.config = config
}

// IncrementSuccess 增加成功计数
func (bc *BaseConsumer) IncrementSuccess() {
	atomic.AddInt64(&bc.stats.SuccessCount, 1)
	atomic.AddInt64(&bc.stats.TotalProcessed, 1)

	bc.mu.Lock()
	bc.stats.LastProcessTime = time.Now()
	bc.mu.Unlock()
}

// IncrementFailure 增加失败计数
func (bc *BaseConsumer) IncrementFailure() {
	atomic.AddInt64(&bc.stats.FailureCount, 1)
	atomic.AddInt64(&bc.stats.TotalProcessed, 1)

	bc.mu.Lock()
	bc.stats.LastProcessTime = time.Now()
	bc.mu.Unlock()
}

// GetStopChan 获取停止信号通道
// 子类可以使用此通道来监听停止信号
func (bc *BaseConsumer) GetStopChan() <-chan struct{} {
	return bc.stopChan
}

// GetWaitGroup 获取等待组
// 子类可以使用此等待组来管理goroutine
func (bc *BaseConsumer) GetWaitGroup() *sync.WaitGroup {
	return &bc.wg
}
