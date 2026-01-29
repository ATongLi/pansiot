package adapter

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"pansiot-device/internal/core"
)

// MockAdapterConfig Mock适配器配置
type MockAdapterConfig struct {
	DeviceID     string
	Protocol     string
	DataDelay    time.Duration // 数据返回延迟（模拟网络延迟）
	ErrorRate    float64       // 错误率（0.0-1.0，用于测试错误处理）
	ValueRange   [2]float64    // 值范围（最小值，最大值）
	AutoIncrement bool          // 值是否自动递增
}

// MockAdapter Mock适配器，用于测试和演示
type MockAdapter struct {
	*BaseAdapter
	config    MockAdapterConfig
	mu        sync.RWMutex
	variables  map[uint64]*core.Variable // 模拟的设备变量
	counter   uint64                    // 自动递增计数器
}

// NewMockAdapter 创建Mock适配器
func NewMockAdapter(config MockAdapterConfig) *MockAdapter {
	return &MockAdapter{
		BaseAdapter: NewBaseAdapter(&core.Device{
		ID:       config.DeviceID,
		Protocol: config.Protocol,
	}),
		config:   config,
		variables: make(map[uint64]*core.Variable),
		counter:  0,
	}
}

// Connect 连接到Mock设备
func (ma *MockAdapter) Connect(ctx context.Context, device *core.Device) error {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	ma.connected = true
	ma.reconnectCount = 0

	log.Printf("[MockAdapter] Connected to mock device: %s", device.ID)

	return nil
}

// Disconnect 断开连接
func (ma *MockAdapter) Disconnect(ctx context.Context) error {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	if !ma.connected {
		return nil
	}

	ma.connected = false
	close(ma.stopChan)

	log.Printf("[MockAdapter] Disconnected from mock device: %s", ma.device.ID)

	return nil
}

// ReadVariable 读取单个变量
func (ma *MockAdapter) ReadVariable(ctx context.Context, variableID uint64) (*core.Variable, error) {
	// 检查连接状态
	if !ma.IsConnected() {
		return nil, fmt.Errorf("not connected to device")
	}

	// 模拟网络延迟
	if ma.config.DataDelay > 0 {
		time.Sleep(ma.config.DataDelay)
	}

	// 模拟错误
	if shouldInjectError(ma.config.ErrorRate) {
		ma.SetLastError(fmt.Errorf("mock read error"))
		ma.IncrementReadCount(false)
		return nil, fmt.Errorf("mock read error")
	}

	ma.mu.RLock()
	defer ma.mu.RUnlock()

	// 生成或返回模拟值
	variable, exists := ma.variables[variableID]
	if !exists {
		// 自动创建变量
		variable = &core.Variable{
			ID:          variableID,
			StringID:    fmt.Sprintf("MOCK-%s-VAR%d", ma.config.DeviceID, variableID),
			Name:        fmt.Sprintf("Mock Variable %d", variableID),
			DataType:    core.DataTypeFloat64,
			DeviceID:    ma.config.DeviceID,
			Value:       ma.generateValue(),
			Quality:     core.QualityGood,
			Timestamp:   time.Now(),
		}
		ma.variables[variableID] = variable
	}

	// 如果配置了自动递增
	if ma.config.AutoIncrement {
		ma.counter++
		variable.Value = float64(ma.counter)
		variable.Timestamp = time.Now()
	}

	ma.IncrementReadCount(true)
	return variable, nil
}

// ReadVariables 批量读取变量
func (ma *MockAdapter) ReadVariables(ctx context.Context, variableIDs []uint64) ([]*core.Variable, error) {
	variables := make([]*core.Variable, 0, len(variableIDs))

	for _, vid := range variableIDs {
		variable, err := ma.ReadVariable(ctx, vid)
		if err != nil {
			return variables, err
		}
		variables = append(variables, variable)
	}

	return variables, nil
}

// WriteVariable 写入单个变量
func (ma *MockAdapter) WriteVariable(ctx context.Context, variableID uint64, value interface{}) error {
	// 检查连接状态
	if !ma.IsConnected() {
		return fmt.Errorf("not connected to device")
	}

	// 模拟网络延迟
	if ma.config.DataDelay > 0 {
		time.Sleep(ma.config.DataDelay)
	}

	// 模拟错误
	if shouldInjectError(ma.config.ErrorRate) {
		ma.SetLastError(fmt.Errorf("mock write error"))
		ma.IncrementWriteCount(false)
		return fmt.Errorf("mock write error")
	}

	ma.mu.Lock()
	defer ma.mu.Unlock()

	// 存储写入的值
	variable := &core.Variable{
		ID:        variableID,
		StringID:  fmt.Sprintf("MOCK-%s-VAR%d", ma.config.DeviceID, variableID),
		DataType:  core.DataTypeFloat64,
		DeviceID:  ma.config.DeviceID,
		Value:     value,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}

	ma.variables[variableID] = variable
	ma.IncrementWriteCount(true)

	return nil
}

// GetProtocol 获取协议类型
func (ma *MockAdapter) GetProtocol() string {
	return "mock"
}

// generateValue 生成模拟值
func (ma *MockAdapter) generateValue() interface{} {
	minVal := ma.config.ValueRange[0]
	maxVal := ma.config.ValueRange[1]

	if minVal == maxVal {
		return minVal
	}

	// 生成范围内的随机值
	return minVal + (maxVal-minVal)*0.5 // 简化：返回中间值
}

// shouldInjectError 判断是否应该注入错误
func shouldInjectError(errorRate float64) bool {
	if errorRate <= 0 {
		return false
	}
	if errorRate >= 1.0 {
		return true
	}
	// 简化实现：这里应该用随机数
	return false
}

// AddVariable 添加模拟变量
func (ma *MockAdapter) AddVariable(variable *core.Variable) {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	ma.variables[variable.ID] = variable
}

// SetVariable 设置模拟变量的值
func (ma *MockAdapter) SetVariable(variableID uint64, value interface{}) error {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	variable, exists := ma.variables[variableID]
	if !exists {
		return fmt.Errorf("variable not found: %d", variableID)
	}

	variable.Value = value
	variable.Timestamp = time.Now()

	return nil
}

// ListVariables 列出所有模拟变量
func (ma *MockAdapter) ListVariables() []*core.Variable {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	variables := make([]*core.Variable, 0, len(ma.variables))
	for _, v := range ma.variables {
		variables = append(variables, v)
	}

	return variables
}
