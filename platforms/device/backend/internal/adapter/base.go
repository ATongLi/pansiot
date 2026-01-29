package adapter

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"pansiot-device/internal/core"
)

// AdapterStats 适配器统计信息
type AdapterStats struct {
	TotalReads    int64     // 总读取次数
	SuccessReads  int64     // 成功读取次数
	FailedReads    int64     // 失败读取次数
	TotalWrites   int64     // 总写入次数
	SuccessWrites int64     // 成功写入次数
	FailedWrites  int64     // 失败写入次数
	LastError     error     // 最后一次错误
	LastErrorTime time.Time // 最后错误时间
}

// BaseAdapter 适配器基类，提供公共功能
type BaseAdapter struct {
	device         *core.Device
	mu             sync.RWMutex
	connected      bool
	lastError      error
	reconnectMax   int
	reconnectCount int
	stats          AdapterStats
	stopChan       chan struct{}
}

// NewBaseAdapter 创建适配器基类实例
func NewBaseAdapter(device *core.Device) *BaseAdapter {
	return &BaseAdapter{
		device:       device,
		reconnectMax: 5, // 默认最多重连5次
		stopChan:     make(chan struct{}),
	}
}

// Connect 连接到设备（基类实现，子类可覆盖）
func (ba *BaseAdapter) Connect(ctx context.Context, device *core.Device) error {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	// 子类应该实现实际的连接逻辑
	ba.connected = true
	ba.reconnectCount = 0
	log.Printf("[Adapter] Connected to device: %s (%s)", device.ID, device.Protocol)
	return nil
}

// Disconnect 断开设备连接
func (ba *BaseAdapter) Disconnect(ctx context.Context) error {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	if !ba.connected {
		return nil
	}

	ba.connected = false
	close(ba.stopChan)
	log.Printf("[Adapter] Disconnected from device: %s", ba.device.ID)
	return nil
}

// IsConnected 检查连接状态
func (ba *BaseAdapter) IsConnected() bool {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.connected
}

// GetDevice 获取设备信息
func (ba *BaseAdapter) GetDevice() *core.Device {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.device
}

// GetStats 获取统计信息
func (ba *BaseAdapter) GetStats() AdapterStats {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.stats
}

// Reconnect 重连设备
func (ba *BaseAdapter) Reconnect(ctx context.Context) error {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	if ba.connected {
		return nil
	}

	if ba.reconnectCount >= ba.reconnectMax {
		return fmt.Errorf("max reconnect attempts (%d) reached", ba.reconnectMax)
	}

	ba.reconnectCount++
	delay := time.Duration(ba.reconnectCount) * time.Second

	log.Printf("[Adapter] Reconnecting to device %s (attempt %d/%d)...",
		ba.device.ID, ba.reconnectCount, ba.reconnectMax)

	time.Sleep(delay)

	// 子类应该实现实际的重连逻辑
	ba.connected = true
	return nil
}

// SetLastError 记录最后错误
func (ba *BaseAdapter) SetLastError(err error) {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	ba.lastError = err
	ba.stats.LastError = err
	ba.stats.LastErrorTime = time.Now()

	if err != nil {
		log.Printf("[Adapter] Error on device %s: %v", ba.device.ID, err)
	}
}

// GetLastError 获取最后错误
func (ba *BaseAdapter) GetLastError() error {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.lastError
}

// IncrementReadCount 增加读取计数
func (ba *BaseAdapter) IncrementReadCount(success bool) {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	ba.stats.TotalReads++
	if success {
		ba.stats.SuccessReads++
	} else {
		ba.stats.FailedReads++
	}
}

// IncrementWriteCount 增加写入计数
func (ba *BaseAdapter) IncrementWriteCount(success bool) {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	ba.stats.TotalWrites++
	if success {
		ba.stats.SuccessWrites++
	} else {
		ba.stats.FailedWrites++
	}
}

// IsStopped 检查是否已停止
func (ba *BaseAdapter) IsStopped() bool {
	select {
	case <-ba.stopChan:
		return true
	default:
		return false
	}
}
