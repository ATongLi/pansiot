package storage

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"pansiot-device/internal/core"
)

// NotifyTask 通知任务
type NotifyTask struct {
	Subscription *Subscription
	Update       core.VariableUpdate
}

// Subscription 订阅信息（内部使用）
type Subscription struct {
	ID         string
	Mode       core.SubscriptionMode
	MatchExpr  string                    // 匹配表达式
	Callback   func(core.VariableUpdate) // 回调函数
	BufferSize int
	CreatedAt  time.Time
	LastNotify time.Time
}

// SubscriberMeta 订阅者元数据
type SubscriberMeta struct {
	ID         string
	TotalSubs  int
	CreatedAt  time.Time
	LastActive time.Time
}

// PubSubManager 发布订阅管理器
type PubSubManager struct {
	mu sync.RWMutex

	// 精确订阅：variableID -> []subscription
	exactSubs map[uint64][]*Subscription

	// 前缀订阅：prefix -> []subscription
	prefixSubs map[string][]*Subscription

	// 通配符订阅：pattern -> []subscription
	wildcardSubs map[string][]*Subscription

	// 设备订阅：deviceID -> []subscription
	deviceSubs map[string][]*Subscription

	// 订阅者元数据
	subscribers map[string]*SubscriberMeta

	// 协程池
	workerPool chan struct{}
	notifyChan chan NotifyTask
}

// NewPubSubManager 创建发布订阅管理器
func NewPubSubManager(workerSize int) *PubSubManager {
	pm := &PubSubManager{
		exactSubs:    make(map[uint64][]*Subscription),
		prefixSubs:   make(map[string][]*Subscription),
		wildcardSubs: make(map[string][]*Subscription),
		deviceSubs:   make(map[string][]*Subscription),
		subscribers:  make(map[string]*SubscriberMeta),
		workerPool:   make(chan struct{}, workerSize),
		notifyChan:   make(chan NotifyTask, 10000),
	}

	// 启动工作协程
	for i := 0; i < workerSize; i++ {
		pm.workerPool <- struct{}{}
		go pm.notifyWorker()
	}

	return pm
}

// Subscribe 精确订阅
func (pm *PubSubManager) Subscribe(subscriberID string, variableIDs []uint64, callback func(core.VariableUpdate)) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	sub := &Subscription{
		ID:         fmt.Sprintf("%s-exact-%v", subscriberID, variableIDs),
		Mode:       core.SubscriptionModeExact,
		Callback:   callback,
		BufferSize: 1000,
		CreatedAt:  time.Now(),
	}

	for _, vid := range variableIDs {
		if pm.exactSubs[vid] == nil {
			pm.exactSubs[vid] = make([]*Subscription, 0)
		}
		pm.exactSubs[vid] = append(pm.exactSubs[vid], sub)
	}

	pm.updateSubscriberMeta(subscriberID)
	return nil
}

// SubscribeByDevice 按设备订阅
func (pm *PubSubManager) SubscribeByDevice(subscriberID, deviceID string, callback func(core.VariableUpdate)) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	sub := &Subscription{
		ID:         fmt.Sprintf("%s-device-%s", subscriberID, deviceID),
		Mode:       core.SubscriptionModeDevice,
		MatchExpr:  deviceID,
		Callback:   callback,
		BufferSize: 1000,
		CreatedAt:  time.Now(),
	}

	if pm.deviceSubs[deviceID] == nil {
		pm.deviceSubs[deviceID] = make([]*Subscription, 0)
	}
	pm.deviceSubs[deviceID] = append(pm.deviceSubs[deviceID], sub)

	pm.updateSubscriberMeta(subscriberID)
	return nil
}

// SubscribeByPattern 按模式订阅
func (pm *PubSubManager) SubscribeByPattern(subscriberID, pattern string, callback func(core.VariableUpdate)) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	mode := core.SubscriptionModeExact
	if strings.Contains(pattern, "*") {
		mode = core.SubscriptionModeWildcard
	} else if strings.HasSuffix(pattern, "-") {
		mode = core.SubscriptionModePrefix
	}

	sub := &Subscription{
		ID:         fmt.Sprintf("%s-pattern-%s", subscriberID, pattern),
		Mode:       mode,
		MatchExpr:  pattern,
		Callback:   callback,
		BufferSize: 1000,
		CreatedAt:  time.Now(),
	}

	key := pattern
	switch mode {
	case core.SubscriptionModeWildcard:
		if pm.wildcardSubs[key] == nil {
			pm.wildcardSubs[key] = make([]*Subscription, 0)
		}
		pm.wildcardSubs[key] = append(pm.wildcardSubs[key], sub)
	case core.SubscriptionModePrefix:
		if pm.prefixSubs[key] == nil {
			pm.prefixSubs[key] = make([]*Subscription, 0)
		}
		pm.prefixSubs[key] = append(pm.prefixSubs[key], sub)
	}

	pm.updateSubscriberMeta(subscriberID)
	return nil
}

// Unsubscribe 取消订阅
func (pm *PubSubManager) Unsubscribe(subscriberID string, variableIDs []uint64) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, vid := range variableIDs {
		if subs, ok := pm.exactSubs[vid]; ok {
			filtered := make([]*Subscription, 0)
			for _, sub := range subs {
				if !strings.HasPrefix(sub.ID, subscriberID) {
					filtered = append(filtered, sub)
				}
			}
			pm.exactSubs[vid] = filtered
		}
	}

	pm.updateSubscriberMeta(subscriberID)
	return nil
}

// UnsubscribeAll 取消所有订阅
func (pm *PubSubManager) UnsubscribeAll(subscriberID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 清理精确订阅
	for vid := range pm.exactSubs {
		filtered := make([]*Subscription, 0)
		for _, sub := range pm.exactSubs[vid] {
			if !strings.HasPrefix(sub.ID, subscriberID) {
				filtered = append(filtered, sub)
			}
		}
		pm.exactSubs[vid] = filtered
	}

	// 清理设备订阅
	for deviceID := range pm.deviceSubs {
		filtered := make([]*Subscription, 0)
		for _, sub := range pm.deviceSubs[deviceID] {
			if !strings.HasPrefix(sub.ID, subscriberID) {
				filtered = append(filtered, sub)
			}
		}
		pm.deviceSubs[deviceID] = filtered
	}

	// 清理前缀订阅
	for prefix := range pm.prefixSubs {
		filtered := make([]*Subscription, 0)
		for _, sub := range pm.prefixSubs[prefix] {
			if !strings.HasPrefix(sub.ID, subscriberID) {
				filtered = append(filtered, sub)
			}
		}
		pm.prefixSubs[prefix] = filtered
	}

	// 清理通配符订阅
	for pattern := range pm.wildcardSubs {
		filtered := make([]*Subscription, 0)
		for _, sub := range pm.wildcardSubs[pattern] {
			if !strings.HasPrefix(sub.ID, subscriberID) {
				filtered = append(filtered, sub)
			}
		}
		pm.wildcardSubs[pattern] = filtered
	}

	delete(pm.subscribers, subscriberID)
	return nil
}

// Publish 发布更新
func (pm *PubSubManager) Publish(variableID uint64, stringID string, deviceID string, update core.VariableUpdate) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// 1. 精确匹配订阅
	if subs, ok := pm.exactSubs[variableID]; ok {
		for _, sub := range subs {
			pm.enqueueNotify(sub, update)
		}
	}

	// 2. 设备订阅
	if subs, ok := pm.deviceSubs[deviceID]; ok {
		for _, sub := range subs {
			pm.enqueueNotify(sub, update)
		}
	}

	// 3. 前缀匹配订阅
	for prefix, subs := range pm.prefixSubs {
		if strings.HasPrefix(stringID, prefix) {
			for _, sub := range subs {
				pm.enqueueNotify(sub, update)
			}
		}
	}

	// 4. 通配符匹配订阅
	for pattern, subs := range pm.wildcardSubs {
		if matchWildcard(pattern, stringID) {
			for _, sub := range subs {
				pm.enqueueNotify(sub, update)
			}
		}
	}
}

// enqueueNotify 入队通知任务
func (pm *PubSubManager) enqueueNotify(sub *Subscription, update core.VariableUpdate) {
	select {
	case pm.notifyChan <- NotifyTask{Subscription: sub, Update: update}:
	default:
		log.Printf("[WARN] notification queue full, dropping update for subscriber: %s", sub.ID)
	}
}

// notifyWorker 通知工作协程
func (pm *PubSubManager) notifyWorker() {
	for task := range pm.notifyChan {
		<-pm.workerPool // 获取令牌

		func() {
			defer func() { pm.workerPool <- struct{}{} }() // 归还令牌

			// 恢复panic，避免单个回调异常导致整个协程崩溃
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[ERROR] subscription callback panic: %v, subscriber: %s", r, task.Subscription.ID)
				}
			}()

			done := make(chan struct{})
			go func() {
				defer close(done)
				task.Subscription.Callback(task.Update)
			}()

			select {
			case <-done:
				task.Subscription.LastNotify = time.Now()
			case <-time.After(100 * time.Millisecond):
				log.Printf("[WARN] subscription callback timeout: %s", task.Subscription.ID)
			}
		}()
	}
}

// GetSubscriberCount 获取订阅者数量
func (pm *PubSubManager) GetSubscriberCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.subscribers)
}

// GetSubscriptionCount 获取总订阅数
func (pm *PubSubManager) GetSubscriptionCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	count := 0
	for _, subs := range pm.exactSubs {
		count += len(subs)
	}
	for _, subs := range pm.deviceSubs {
		count += len(subs)
	}
	for _, subs := range pm.prefixSubs {
		count += len(subs)
	}
	for _, subs := range pm.wildcardSubs {
		count += len(subs)
	}

	return count
}

// updateSubscriberMeta 更新订阅者元数据
func (pm *PubSubManager) updateSubscriberMeta(subscriberID string) {
	meta := pm.subscribers[subscriberID]
	if meta == nil {
		meta = &SubscriberMeta{
			ID:        subscriberID,
			CreatedAt: time.Now(),
		}
		pm.subscribers[subscriberID] = meta
	}
	meta.LastActive = time.Now()
}

// matchWildcard 通配符匹配
func matchWildcard(pattern, s string) bool {
	// 简化版实现，支持 * 通配符
	// 生产环境建议使用 github.com/gobwas/glob
	patternParts := strings.Split(pattern, "*")
	if len(patternParts) == 1 {
		return pattern == s
	}

	idx := 0
	for i, part := range patternParts {
		if part == "" {
			continue
		}
		found := strings.Index(s[idx:], part)
		if found == -1 {
			return false
		}
		idx += found + len(part)
		if i == len(patternParts)-1 && idx != len(s) {
			return false
		}
	}
	return true
}
