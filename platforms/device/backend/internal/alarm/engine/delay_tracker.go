package engine

import (
	"fmt"
	"sync"
	"time"

	"pansiot-device/internal/alarm/rule"
)

// DelayTracker 延迟跟踪器
// 管理报警规则的延迟触发定时器
type DelayTracker struct {
	mu     sync.RWMutex
	timers map[string]*DelayTimer  // ruleID -> 延迟定时器
}

// DelayTimer 延迟定时器
type DelayTimer struct {
	ruleID       string
	condition    *rule.SingleCondition
	variableID   uint64
	targetValue  interface{}  // 目标值（需要达到的值）
	timer        *time.Timer
	startTime    time.Time
	triggerTime  time.Time // 预期触发时间
	onTimeout    func()     // 超时回调
	canceled     bool       // 是否已取消
	mu           sync.Mutex
}

// NewDelayTracker 创建延迟跟踪器
func NewDelayTracker() *DelayTracker {
	return &DelayTracker{
		timers: make(map[string]*DelayTimer),
	}
}

// StartDelay 启动延迟定时器
// 参数：
//   - ruleID: 规则ID
//   - cond: 单个条件（包含延迟时间）
//   - variableID: 变量ID
//   - targetValue: 目标值
//   - onTimeout: 延迟到期后的回调函数
// 返回：
//   - error: 错误信息
func (dt *DelayTracker) StartDelay(ruleID string, cond *rule.SingleCondition, variableID uint64, targetValue interface{}, onTimeout func()) error {
	if cond.Delay <= 0 {
		return fmt.Errorf("延迟时间必须大于0: %v", cond.Delay)
	}

	dt.mu.Lock()
	defer dt.mu.Unlock()

	// 如果已有定时器，先取消
	if existingTimer, exists := dt.timers[ruleID]; exists {
		existingTimer.Stop()
	}

	delayTimer := &DelayTimer{
		ruleID:      ruleID,
		condition:   cond,
		variableID:  variableID,
		targetValue: targetValue,
		startTime:   time.Now(),
		triggerTime: time.Now().Add(cond.Delay),
		onTimeout:   onTimeout,
	}

	// 创建定时器
	delayTimer.timer = time.AfterFunc(cond.Delay, func() {
		delayTimer.mu.Lock()
		defer delayTimer.mu.Unlock()

		if !delayTimer.canceled {
			// 执行回调
			if delayTimer.onTimeout != nil {
				delayTimer.onTimeout()
			}

			// 清理定时器
			dt.mu.Lock()
			delete(dt.timers, ruleID)
			dt.mu.Unlock()
		}
	})

	dt.timers[ruleID] = delayTimer
	return nil
}

// Cancel 取消延迟定时器
func (dt *DelayTracker) Cancel(ruleID string) bool {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	if delayTimer, exists := dt.timers[ruleID]; exists {
		delayTimer.Cancel()
		delete(dt.timers, ruleID)
		return true
	}
	return false
}

// Stop 停止延迟定时器（内部方法）
func (dt *DelayTimer) Stop() {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	dt.canceled = true
	if dt.timer != nil {
		dt.timer.Stop()
	}
}

// Cancel 取消定时器（公开方法）
func (dt *DelayTimer) Cancel() {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	dt.canceled = true
	if dt.timer != nil {
		dt.timer.Stop()
	}
}

// IsActive 检查延迟定时器是否活跃
func (dt *DelayTracker) IsActive(ruleID string) bool {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if delayTimer, exists := dt.timers[ruleID]; exists {
		delayTimer.mu.Lock()
		active := !delayTimer.canceled && time.Now().Before(delayTimer.triggerTime)
		delayTimer.mu.Unlock()
		return active
	}
	return false
}

// GetRemainingTime 获取剩余延迟时间
func (dt *DelayTracker) GetRemainingTime(ruleID string) time.Duration {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if delayTimer, exists := dt.timers[ruleID]; exists {
		delayTimer.mu.Lock()
		defer delayTimer.mu.Unlock()

		if delayTimer.canceled {
			return 0
		}

		remaining := time.Until(delayTimer.triggerTime)
		if remaining < 0 {
			return 0
		}
		return remaining
	}
	return 0
}

// GetStartTime 获取延迟开始时间
func (dt *DelayTracker) GetStartTime(ruleID string) time.Time {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if delayTimer, exists := dt.timers[ruleID]; exists {
		delayTimer.mu.Lock()
		defer delayTimer.mu.Unlock()
		return delayTimer.startTime
	}
	return time.Time{}
}

// GetTriggerTime 获取预期触发时间
func (dt *DelayTracker) GetTriggerTime(ruleID string) time.Time {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if delayTimer, exists := dt.timers[ruleID]; exists {
		delayTimer.mu.Lock()
		defer delayTimer.mu.Unlock()
		return delayTimer.triggerTime
	}
	return time.Time{}
}

// Clear 清除所有延迟定时器
func (dt *DelayTracker) Clear() {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	// 停止所有定时器
	for _, delayTimer := range dt.timers {
		delayTimer.Stop()
	}

	dt.timers = make(map[string]*DelayTimer)
}

// GetActiveCount 获取活跃的延迟定时器数量
func (dt *DelayTracker) GetActiveCount() int {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	count := 0
	for _, delayTimer := range dt.timers {
		delayTimer.mu.Lock()
		if !delayTimer.canceled && time.Now().Before(delayTimer.triggerTime) {
			count++
		}
		delayTimer.mu.Unlock()
	}
	return count
}

// GetAllTimers 获取所有延迟定时器信息
func (dt *DelayTracker) GetAllTimers() map[string]*DelayTimerInfo {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

 infos := make(map[string]*DelayTimerInfo)
	for ruleID, delayTimer := range dt.timers {
		delayTimer.mu.Lock()
		infos[ruleID] = &DelayTimerInfo{
			RuleID:      ruleID,
			VariableID:  delayTimer.variableID,
			StartTime:   delayTimer.startTime,
			TriggerTime: delayTimer.triggerTime,
			Remaining:   time.Until(delayTimer.triggerTime),
			IsActive:    !delayTimer.canceled,
		}
		delayTimer.mu.Unlock()
	}
	return infos
}

// DelayTimerInfo 延迟定时器信息
type DelayTimerInfo struct {
	RuleID      string
	VariableID  uint64
	StartTime   time.Time
	TriggerTime time.Time
	Remaining   time.Duration
	IsActive    bool
}

// CheckAndRestart 检查条件是否仍然满足，如果不满足则取消定时器
// 用于在延迟期间持续监控条件状态
func (dt *DelayTracker) CheckAndRestart(ruleID string, stillSatisfied bool) {
	if !stillSatisfied {
		// 条件不再满足，取消定时器
		dt.Cancel(ruleID)
	}
}
