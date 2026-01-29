package script

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// ScriptExecutor 脚本执行器接口
// 用于解耦调度器和具体消费者实现
type ScriptExecutor interface {
	ExecuteScriptAsync(scriptID string, input map[string]interface{}) error
}

// ScriptScheduler 脚本调度器
// 管理周期性执行的脚本调度（支持固定间隔和Cron表达式）
type ScriptScheduler struct {
	mu              sync.RWMutex
	triggers        map[string]*Trigger         // 所有触发器 (ID -> Trigger)
	intervalTickers map[string]*time.Ticker     // 固定间隔触发器 (ID -> Ticker)
	intervalStops   map[string]chan struct{}    // 固定间隔停止信号 (ID -> StopChan)
	cronJobs        map[string]cron.EntryID      // Cron表达式触发器 (ID -> EntryID)
	cron            *cron.Cron                   // Cron调度器实例
	executor        ScriptExecutor               // 脚本执行器
	running         bool                         // 运行状态
	muRunning       sync.Mutex                   // 保护 running 字段
}

// NewScriptScheduler 创建脚本调度器
func NewScriptScheduler(executor ScriptExecutor) *ScriptScheduler {
	return &ScriptScheduler{
		triggers:        make(map[string]*Trigger),
		intervalTickers: make(map[string]*time.Ticker),
		intervalStops:   make(map[string]chan struct{}),
		cronJobs:        make(map[string]cron.EntryID),
		cron:            cron.New(cron.WithSeconds()), // 支持秒级精度（6段式）
		executor:        executor,
		running:         false,
	}
}

// Start 启动调度器
func (s *ScriptScheduler) Start() error {
	s.muRunning.Lock()
	defer s.muRunning.Unlock()

	if s.running {
		return fmt.Errorf("调度器已在运行")
	}

	// 启动Cron调度器
	s.cron.Start()

	s.running = true
	log.Printf("[ScriptScheduler] 调度器已启动 (支持Interval和Cron模式)")
	return nil
}

// Stop 停止调度器
func (s *ScriptScheduler) Stop() error {
	s.muRunning.Lock()
	defer s.muRunning.Unlock()

	if !s.running {
		return fmt.Errorf("调度器未运行")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 停止所有Interval Tickers
	for triggerID, stopChan := range s.intervalStops {
		close(stopChan)
		delete(s.intervalStops, triggerID)
	}

	for triggerID, ticker := range s.intervalTickers {
		ticker.Stop()
		delete(s.intervalTickers, triggerID)
	}

	// 停止Cron调度器
	ctx := s.cron.Stop()
	// 等待Cron停止（最多等待5秒）
	select {
	case <-ctx.Done():
		log.Printf("[ScriptScheduler] Cron调度器已停止")
	case <-time.After(5 * time.Second):
		log.Printf("[ScriptScheduler] Cron调度器停止超时")
	}

	// 清空
	s.triggers = make(map[string]*Trigger)

	s.running = false
	log.Printf("[ScriptScheduler] 调度器已停止")
	return nil
}

// AddTrigger 添加周期触发器
func (s *ScriptScheduler) AddTrigger(trigger *Trigger) error {
	if trigger.Type != TriggerTypePeriodic {
		return fmt.Errorf("只支持周期触发器")
	}

	if trigger.Condition.PeriodicConfig == nil {
		return fmt.Errorf("周期触发器必须配置PeriodicConfig")
	}

	s.muRunning.Lock()
	defer s.muRunning.Unlock()

	if !s.running {
		return fmt.Errorf("调度器未运行")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否已存在
	if _, exists := s.triggers[trigger.ID]; exists {
		return fmt.Errorf("触发器已存在: %s", trigger.ID)
	}

	config := trigger.Condition.PeriodicConfig

	// 判断调度模式
	if config.CronExpr != "" {
		// Cron表达式模式
		return s.addCronTrigger(trigger, config)
	} else if config.Interval > 0 {
		// 固定间隔模式
		return s.addIntervalTrigger(trigger, config)
	} else {
		return fmt.Errorf("必须指定CronExpr或Interval之一")
	}
}

// addCronTrigger 添加Cron表达式触发器
func (s *ScriptScheduler) addCronTrigger(trigger *Trigger, config *PeriodicTriggerConfig) error {
	// 验证Cron表达式
	if !s.cronScheduleIsValid(config.CronExpr) {
		return fmt.Errorf("无效的Cron表达式: %s", config.CronExpr)
	}

	// 创建包装函数，集成时间窗口检查
	jobFunc := func() {
		// 检查时间窗口（复用现有逻辑）
		if s.shouldExecuteByTimeWindow(config) {
			log.Printf("[ScriptScheduler] Cron触发器触发: %s", trigger.ID)

			// 异步执行脚本
			s.executeScript(trigger, "periodic")
		}
	}

	// 添加到Cron调度器
	entryID, err := s.cron.AddFunc(config.CronExpr, jobFunc)
	if err != nil {
		return fmt.Errorf("添加Cron任务失败: %v", err)
	}

	// 保存映射关系
	s.triggers[trigger.ID] = trigger
	s.cronJobs[trigger.ID] = entryID

	log.Printf("[ScriptScheduler] 已添加Cron触发器: %s, 表达式: %s",
		trigger.ID, config.CronExpr)
	return nil
}

// addIntervalTrigger 添加固定间隔触发器
func (s *ScriptScheduler) addIntervalTrigger(trigger *Trigger, config *PeriodicTriggerConfig) error {
	// 验证间隔
	if config.Interval <= 0 {
		return fmt.Errorf("执行间隔必须大于0")
	}

	// 创建停止信号通道
	stopChan := make(chan struct{})

	// 创建定时器
	ticker := time.NewTicker(config.Interval)

	// 启动监听goroutine
	go func(triggerID string) {
		log.Printf("[ScriptScheduler] Interval触发器启动: %s (间隔: %v)",
			triggerID, config.Interval)

		for {
			select {
			case <-ticker.C:
				// 检查时间窗口
				if s.shouldExecuteByTimeWindow(config) {
					log.Printf("[ScriptScheduler] Interval触发器触发: %s", triggerID)
					s.executeScript(trigger, "periodic")
				}

			case <-stopChan:
				log.Printf("[ScriptScheduler] Interval触发器停止: %s", triggerID)
				return
			}
		}
	}(trigger.ID)

	// 保存映射关系
	s.triggers[trigger.ID] = trigger
	s.intervalTickers[trigger.ID] = ticker
	s.intervalStops[trigger.ID] = stopChan

	log.Printf("[ScriptScheduler] 已添加Interval触发器: %s, 间隔: %v",
		trigger.ID, config.Interval)
	return nil
}

// RemoveTrigger 移除周期触发器
func (s *ScriptScheduler) RemoveTrigger(triggerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trigger, exists := s.triggers[triggerID]
	if !exists {
		return fmt.Errorf("触发器不存在: %s", triggerID)
	}

	config := trigger.Condition.PeriodicConfig

	// 根据调度类型移除
	if config.CronExpr != "" {
		// 移除Cron触发器
		entryID, ok := s.cronJobs[triggerID]
		if ok {
			s.cron.Remove(entryID)
			delete(s.cronJobs, triggerID)
		}
	} else {
		// 移除Interval触发器
		if stopChan, exists := s.intervalStops[triggerID]; exists {
			close(stopChan)
			delete(s.intervalStops, triggerID)
		}

		if ticker, exists := s.intervalTickers[triggerID]; exists {
			ticker.Stop()
			delete(s.intervalTickers, triggerID)
		}
	}

	// 从触发器列表中删除
	delete(s.triggers, triggerID)

	log.Printf("[ScriptScheduler] 已移除触发器: %s", triggerID)
	return nil
}

// shouldExecuteByTimeWindow 检查时间窗口（提取现有逻辑）
func (s *ScriptScheduler) shouldExecuteByTimeWindow(config *PeriodicTriggerConfig) bool {
	now := time.Now()

	// 检查星期几
	if len(config.DaysOfWeek) > 0 {
		currentDay := int(now.Weekday())

		// 转换：Go的周日=0，我们的配置中周日=7
		if currentDay == 0 {
			currentDay = 7
		}

		found := false
		for _, day := range config.DaysOfWeek {
			if day == currentDay {
				found = true
				break
			}
		}

		if !found {
			log.Printf("[ScriptScheduler] 跳过执行: 不在指定的星期几 (当前: %d, 配置: %v)",
				currentDay, config.DaysOfWeek)
			return false
		}
	}

	// 检查时间窗口
	if config.StartTime != "" || config.EndTime != "" {
		currentTime := now.Format("15:04:05")

		// 检查开始时间
		if config.StartTime != "" && currentTime < config.StartTime {
			log.Printf("[ScriptScheduler] 跳过执行: 未到开始时间 (当前: %s, 开始: %s)",
				currentTime, config.StartTime)
			return false
		}

		// 检查结束时间
		if config.EndTime != "" && currentTime > config.EndTime {
			log.Printf("[ScriptScheduler] 跳过执行: 已过结束时间 (当前: %s, 结束: %s)",
				currentTime, config.EndTime)
			return false
		}
	}

	return true
}

// cronScheduleIsValid 验证Cron表达式
func (s *ScriptScheduler) cronScheduleIsValid(expr string) bool {
	_, err := cron.ParseStandard(expr)
	return err == nil
}

// executeScript 执行脚本（提取公共逻辑）
func (s *ScriptScheduler) executeScript(trigger *Trigger, triggerType string) {
	err := s.executor.ExecuteScriptAsync(trigger.ScriptID, map[string]interface{}{
		"trigger_type": triggerType,
		"trigger_id":   trigger.ID,
		"timestamp":    time.Now(),
	})

	if err != nil {
		log.Printf("[ScriptScheduler] 脚本执行失败: %s, 错误: %v",
			trigger.ScriptID, err)
	}
}

// IsRunning 检查调度器是否运行中
func (s *ScriptScheduler) IsRunning() bool {
	s.muRunning.Lock()
	defer s.muRunning.Unlock()
	return s.running
}

// GetTriggerCount 获取当前触发器数量
func (s *ScriptScheduler) GetTriggerCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.triggers)
}

// ListTriggers 列出所有周期触发器
func (s *ScriptScheduler) ListTriggers() []*Trigger {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]*Trigger, 0, len(s.triggers))
	for _, trigger := range s.triggers {
		result = append(result, trigger)
	}

	return result
}

// UpdateTriggerInterval 更新触发器执行间隔（仅Interval模式）
func (s *ScriptScheduler) UpdateTriggerInterval(triggerID string, newInterval time.Duration) error {
	if newInterval <= 0 {
		return fmt.Errorf("执行间隔必须大于0")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否存在
	trigger, exists := s.triggers[triggerID]
	if !exists {
		return fmt.Errorf("触发器不存在: %s", triggerID)
	}

	config := trigger.Condition.PeriodicConfig

	// 只有Interval触发器才能更新间隔
	if config.CronExpr != "" {
		return fmt.Errorf("Cron触发器不支持更新间隔，请使用RemoveTrigger后重新添加")
	}

	// 停止旧的ticker
	if stopChan, exists := s.intervalStops[triggerID]; exists {
		close(stopChan)
		delete(s.intervalStops, triggerID)
	}

	if ticker, exists := s.intervalTickers[triggerID]; exists {
		ticker.Stop()
		delete(s.intervalTickers, triggerID)
	}

	// 更新间隔
	config.Interval = newInterval

	// 重新添加Interval触发器
	return s.addIntervalTrigger(trigger, config)
}
