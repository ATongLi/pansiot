package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"pansiot-device/internal/core"
)

// MockConsumer Mock消费者
// 用于测试订阅机制是否正常工作
type MockConsumer struct {
	*BaseConsumer                  // 嵌入基类
	subscribeVariableIDs []uint64  // 要订阅的变量ID列表
	receivedData         []VariableData // 接收到的数据
	mu                   sync.RWMutex   // 保护 receivedData
}

// VariableData 变量数据记录
type VariableData struct {
	VariableID uint64
	Value      interface{}
	Quality    core.QualityCode
	Timestamp  time.Time
}

// MockConsumerConfig Mock消费者配置
type MockConsumerConfig struct {
	SubscribeVariableIDs []uint64 // 要订阅的变量ID列表
}

// NewMockConsumer 创建Mock消费者
func NewMockConsumer(id string, storage core.Storage, config MockConsumerConfig) *MockConsumer {
	mock := &MockConsumer{
		BaseConsumer:         NewBaseConsumer(id, "mock", storage),
		subscribeVariableIDs: config.SubscribeVariableIDs,
		receivedData:         make([]VariableData, 0),
	}

	mock.SetConfig(config)
	return mock
}

// Start 启动Mock消费者
func (mc *MockConsumer) Start(ctx context.Context) error {
	if err := mc.BaseConsumer.Start(ctx); err != nil {
		return err
	}

	// 订阅变量
	if len(mc.subscribeVariableIDs) > 0 {
		if err := mc.storage.Subscribe(mc.GetID(), mc.subscribeVariableIDs, mc.onVariableUpdate); err != nil {
			mc.BaseConsumer.Stop()
			return fmt.Errorf("订阅变量失败: %v", err)
		}
		log.Printf("[MockConsumer] 消费者[%s]已订阅 %d 个变量: %v",
			mc.GetID(), len(mc.subscribeVariableIDs), mc.subscribeVariableIDs)
	}

	log.Printf("[MockConsumer] 消费者[%s]已启动", mc.GetID())
	return nil
}

// Stop 停止Mock消费者
func (mc *MockConsumer) Stop() error {
	if !mc.IsRunning() {
		return fmt.Errorf("消费者[%s]未在运行", mc.GetID())
	}

	// 取消订阅
	if len(mc.subscribeVariableIDs) > 0 {
		if err := mc.storage.UnsubscribeAll(mc.GetID()); err != nil {
			log.Printf("[MockConsumer] 取消订阅失败: %v", err)
		}
	}

	if err := mc.BaseConsumer.Stop(); err != nil {
		return err
	}

	log.Printf("[MockConsumer] 消费者[%s]已停止", mc.GetID())
	return nil
}

// onVariableUpdate 变量更新回调函数
func (mc *MockConsumer) onVariableUpdate(update core.VariableUpdate) {
	mc.mu.Lock()
	// 记录接收到的数据
	data := VariableData{
		VariableID: update.VariableID,
		Value:      update.Value,
		Quality:    update.Quality,
		Timestamp:  update.Timestamp,
	}
	mc.receivedData = append(mc.receivedData, data)
	receiveCount := len(mc.receivedData)
	mc.mu.Unlock()

	// 更新统计
	mc.IncrementSuccess()

	log.Printf("[MockConsumer] 消费者[%s]收到更新: 变量ID=%d, 值=%v, 质量=%d, 时间=%v (累计接收%d条)",
		mc.GetID(), update.VariableID, update.Value, update.Quality,
		update.Timestamp.Format("15:04:05.000"), receiveCount)
}

// GetReceivedData 获取接收到的所有数据
func (mc *MockConsumer) GetReceivedData() []VariableData {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// 返回副本，避免外部修改
	result := make([]VariableData, len(mc.receivedData))
	copy(result, mc.receivedData)
	return result
}

// GetReceivedDataByVariableID 获取指定变量的接收数据
func (mc *MockConsumer) GetReceivedDataByVariableID(variableID uint64) []VariableData {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	result := make([]VariableData, 0)
	for _, data := range mc.receivedData {
		if data.VariableID == variableID {
			result = append(result, data)
		}
	}
	return result
}

// GetReceiveCount 获取总接收次数
func (mc *MockConsumer) GetReceiveCount() int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return len(mc.receivedData)
}

// GetReceiveCountByVariableID 获取指定变量的接收次数
func (mc *MockConsumer) GetReceiveCountByVariableID(variableID uint64) int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	count := 0
	for _, data := range mc.receivedData {
		if data.VariableID == variableID {
			count++
		}
	}
	return count
}

// ClearReceivedData 清空接收到的数据
func (mc *MockConsumer) ClearReceivedData() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.receivedData = make([]VariableData, 0)
	log.Printf("[MockConsumer] 消费者[%s]已清空接收数据", mc.GetID())
}

// PrintSummary 打印接收数据摘要
func (mc *MockConsumer) PrintSummary() {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	stats := mc.GetStats()

	fmt.Printf("\n=== Mock消费者[%s]统计摘要 ===\n", mc.GetID())
	fmt.Printf("订阅变量数: %d\n", len(mc.subscribeVariableIDs))
	fmt.Printf("总接收次数: %d\n", len(mc.receivedData))
	fmt.Printf("总处理次数: %d\n", stats.TotalProcessed)
	fmt.Printf("成功次数: %d\n", stats.SuccessCount)
	fmt.Printf("失败次数: %d\n", stats.FailureCount)
	fmt.Printf("最后处理时间: %v\n", stats.LastProcessTime.Format("2006-01-02 15:04:05.000"))

	// 统计每个变量的接收次数
	variableCount := make(map[uint64]int)
	for _, data := range mc.receivedData {
		variableCount[data.VariableID]++
	}

	fmt.Printf("\n各变量接收统计:\n")
	for _, varID := range mc.subscribeVariableIDs {
		count := variableCount[varID]
		fmt.Printf("  变量 %d: %d 次\n", varID, count)
	}
	fmt.Println("================================\n")
}

// 确保MockConsumer实现了core.Consumer接口
var _ core.Consumer = (*MockConsumer)(nil)
