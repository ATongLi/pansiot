package storage

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"pansiot-device/internal/core"
)

const (
	shardCount         = 64  // 分片数量
	maxConcurrentReads = 100 // 最大并发读取数
)

// shard 数据分片
type shard struct {
	mu         sync.RWMutex
	vars       map[uint64]*core.Variable // 数字ID索引
	strIDs     map[string]*core.Variable // 字符串ID索引
	deviceVars map[string][]uint64       // 设备->变量ID列表
}

// MemoryStorage 内存存储实现
type MemoryStorage struct {
	shards    [shardCount]*shard
	pubsub    *PubSubManager
	idGen     *core.IDGenerator
	stats     core.StorageStats
	startTime time.Time

	// 读取限流
	readSem chan struct{}

	// 原子计数器
	readCount  atomic.Int64
	writeCount atomic.Int64
}

// NewMemoryStorage 创建内存存储实例
func NewMemoryStorage() *MemoryStorage {
	ms := &MemoryStorage{
		pubsub:    NewPubSubManager(100), // 100 个通知工作协程
		idGen:     core.NewIDGenerator(),
		startTime: time.Now(),
		readSem:   make(chan struct{}, maxConcurrentReads),
	}

	// 初始化分片
	for i := 0; i < shardCount; i++ {
		ms.shards[i] = &shard{
			vars:       make(map[uint64]*core.Variable),
			strIDs:     make(map[string]*core.Variable),
			deviceVars: make(map[string][]uint64),
		}
	}

	return ms
}

// getShard 根据变量ID获取分片
func (ms *MemoryStorage) getShard(variableID uint64) *shard {
	return ms.shards[variableID&(shardCount-1)]
}

// getShardByStringID 根据字符串ID获取分片
func (ms *MemoryStorage) getShardByStringID(stringID string) *shard {
	hash := fnvHash(stringID)
	return ms.shards[hash&(shardCount-1)]
}

// ReadVar 读取单个变量
func (ms *MemoryStorage) ReadVar(variableID uint64) (*core.Variable, error) {
	shard := ms.getShard(variableID)

	shard.mu.RLock()
	variable := shard.vars[variableID]
	shard.mu.RUnlock()

	if variable == nil {
		return nil, fmt.Errorf("variable not found: %d", variableID)
	}

	ms.readCount.Add(1)
	return ms.cloneVariable(variable), nil
}

// ReadVars 批量读取变量
func (ms *MemoryStorage) ReadVars(variableIDs []uint64) ([]*core.Variable, error) {
	ms.readSem <- struct{}{} // 获取读令牌
	defer func() { <-ms.readSem }()

	if len(variableIDs) == 0 {
		return []*core.Variable{}, nil
	}

	// 按分片分组
	shardGroups := make(map[int][]uint64)
	for _, vid := range variableIDs {
		shardIdx := int(vid & (shardCount - 1))
		shardGroups[shardIdx] = append(shardGroups[shardIdx], vid)
	}

	results := make([]*core.Variable, 0, len(variableIDs))
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 并行读取各分片
	for shardIdx, vids := range shardGroups {
		wg.Add(1)
		go func(idx int, ids []uint64) {
			defer wg.Done()

			shard := ms.shards[idx]
			shard.mu.RLock()
			defer shard.mu.RUnlock()

			mu.Lock()
			for _, vid := range ids {
				if variable := shard.vars[vid]; variable != nil {
					results = append(results, ms.cloneVariable(variable))
				}
			}
			mu.Unlock()
		}(shardIdx, vids)
	}

	wg.Wait()
	ms.readCount.Add(int64(len(variableIDs)))
	return results, nil
}

// ReadVarByStringID 通过字符串ID读取变量
func (ms *MemoryStorage) ReadVarByStringID(stringID string) (*core.Variable, error) {
	shard := ms.getShardByStringID(stringID)

	shard.mu.RLock()
	variable := shard.strIDs[stringID]
	shard.mu.RUnlock()

	if variable == nil {
		return nil, fmt.Errorf("variable not found: %s", stringID)
	}

	ms.readCount.Add(1)
	return ms.cloneVariable(variable), nil
}

// WriteVar 写入单个变量
func (ms *MemoryStorage) WriteVar(variable *core.Variable) error {
	shard := ms.getShard(variable.ID)

	shard.mu.Lock()

	existingVar := shard.vars[variable.ID]
	if existingVar != nil {
		// 更新热数据
		existingVar.Value = variable.Value
		existingVar.Quality = variable.Quality
		existingVar.Timestamp = variable.Timestamp
	} else {
		// 新变量
		shard.vars[variable.ID] = variable
		shard.strIDs[variable.StringID] = variable
		shard.deviceVars[variable.DeviceID] = append(
			shard.deviceVars[variable.DeviceID], variable.ID)
	}

	shard.mu.Unlock()

	// 异步通知订阅者
	update := core.VariableUpdate{
		VariableID: variable.ID,
		Value:      variable.Value,
		Quality:    variable.Quality,
		Timestamp:  variable.Timestamp,
	}
	ms.pubsub.Publish(variable.ID, variable.StringID, variable.DeviceID, update)

	ms.writeCount.Add(1)
	return nil
}

// WriteVars 批量写入变量
func (ms *MemoryStorage) WriteVars(variables []*core.Variable) error {
	if len(variables) == 0 {
		return nil
	}

	// 按分片分组
	shardGroups := make(map[int][]*core.Variable)
	for _, variable := range variables {
		shardIdx := int(variable.ID & (shardCount - 1))
		shardGroups[shardIdx] = append(shardGroups[shardIdx], variable)
	}

	// 并行写入各分片
	var wg sync.WaitGroup
	for shardIdx, vars := range shardGroups {
		wg.Add(1)
		go func(idx int, shardVars []*core.Variable) {
			defer wg.Done()

			shard := ms.shards[idx]
			shard.mu.Lock()
			defer shard.mu.Unlock()

			for _, variable := range shardVars {
				existingVar := shard.vars[variable.ID]
				if existingVar != nil {
					existingVar.Value = variable.Value
					existingVar.Quality = variable.Quality
					existingVar.Timestamp = variable.Timestamp
				} else {
					shard.vars[variable.ID] = variable
					shard.strIDs[variable.StringID] = variable
					shard.deviceVars[variable.DeviceID] = append(
						shard.deviceVars[variable.DeviceID], variable.ID)
				}
			}
		}(shardIdx, vars)
	}

	wg.Wait()

	// 批量通知订阅者
	for _, variable := range variables {
		update := core.VariableUpdate{
			VariableID: variable.ID,
			Value:      variable.Value,
			Quality:    variable.Quality,
			Timestamp:  variable.Timestamp,
		}
		ms.pubsub.Publish(variable.ID, variable.StringID, variable.DeviceID, update)
	}

	ms.writeCount.Add(int64(len(variables)))
	return nil
}

// Subscribe 订阅变量更新
func (ms *MemoryStorage) Subscribe(subscriberID string, variableIDs []uint64, callback func(core.VariableUpdate)) error {
	return ms.pubsub.Subscribe(subscriberID, variableIDs, callback)
}

// SubscribeByDevice 按设备订阅
func (ms *MemoryStorage) SubscribeByDevice(subscriberID, deviceID string, callback func(core.VariableUpdate)) error {
	return ms.pubsub.SubscribeByDevice(subscriberID, deviceID, callback)
}

// SubscribeByPattern 按模式订阅
func (ms *MemoryStorage) SubscribeByPattern(subscriberID, pattern string, callback func(core.VariableUpdate)) error {
	return ms.pubsub.SubscribeByPattern(subscriberID, pattern, callback)
}

// Unsubscribe 取消订阅
func (ms *MemoryStorage) Unsubscribe(subscriberID string, variableIDs []uint64) error {
	return ms.pubsub.Unsubscribe(subscriberID, variableIDs)
}

// UnsubscribeAll 取消所有订阅
func (ms *MemoryStorage) UnsubscribeAll(subscriberID string) error {
	return ms.pubsub.UnsubscribeAll(subscriberID)
}

// CreateVariable 创建新变量
func (ms *MemoryStorage) CreateVariable(variable *core.Variable) error {
	if variable.ID == 0 {
		// 自动生成ID
		id, err := ms.idGen.GenerateCustomID()
		if err != nil {
			return err
		}
		variable.ID = id
	}

	shard := ms.getShard(variable.ID)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, exists := shard.vars[variable.ID]; exists {
		return fmt.Errorf("variable already exists: %d", variable.ID)
	}

	shard.vars[variable.ID] = variable
	shard.strIDs[variable.StringID] = variable
	shard.deviceVars[variable.DeviceID] = append(
		shard.deviceVars[variable.DeviceID], variable.ID)

	ms.stats.TotalVariables++
	return nil
}

// DeleteVariable 删除变量
func (ms *MemoryStorage) DeleteVariable(variableID uint64) error {
	shard := ms.getShard(variableID)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	variable := shard.vars[variableID]
	if variable == nil {
		return fmt.Errorf("variable not found: %d", variableID)
	}

	delete(shard.vars, variableID)
	delete(shard.strIDs, variable.StringID)

	// 从设备索引中删除（惰性清理）
	// 定期清理任务会处理残留

	ms.stats.TotalVariables--
	return nil
}

// ListVariables 列出所有变量
func (ms *MemoryStorage) ListVariables() []*core.Variable {
	allVars := make([]*core.Variable, 0, ms.stats.TotalVariables)

	for _, shard := range ms.shards {
		shard.mu.RLock()
		for _, variable := range shard.vars {
			allVars = append(allVars, ms.cloneVariable(variable))
		}
		shard.mu.RUnlock()
	}

	return allVars
}

// ListVariablesByDevice 列出指定设备的所有变量
func (ms *MemoryStorage) ListVariablesByDevice(deviceID string) []*core.Variable {
	// 遍历所有分片的设备索引
	var varIDs []uint64
	for _, shard := range ms.shards {
		shard.mu.RLock()
		if ids, ok := shard.deviceVars[deviceID]; ok {
			varIDs = append(varIDs, ids...)
		}
		shard.mu.RUnlock()
	}

	// 批量读取变量
	variables, _ := ms.ReadVars(varIDs)
	return variables
}

// GetStats 获取存储统计信息
func (ms *MemoryStorage) GetStats() core.StorageStats {
	return core.StorageStats{
		TotalVariables:     ms.stats.TotalVariables,
		TotalSubscribers:   ms.pubsub.GetSubscriberCount(),
		TotalSubscriptions: ms.pubsub.GetSubscriptionCount(),
		ReadCount:          ms.readCount.Load(),
		WriteCount:         ms.writeCount.Load(),
		StartTime:          ms.startTime,
	}
}

// cloneVariable 深拷贝变量
func (ms *MemoryStorage) cloneVariable(variable *core.Variable) *core.Variable {
	// 优化：浅拷贝 + 值类型拷贝
	clone := *variable
	return &clone
}

// fnvHash FNV哈希算法
func fnvHash(s string) uint64 {
	const (
		offset64 uint64 = 14695981039346656037
		prime64  uint64 = 1099511628211
	)

	h := offset64
	for _, c := range []byte(s) {
		h ^= uint64(c)
		h *= prime64
	}
	return h
}
