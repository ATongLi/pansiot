package action

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/alarm/rule"
)

// Scheduler 动作调度器（处理延迟和循环执行）
type Scheduler struct {
	executor        *Executor
	stateMachine    engine.StateMachine
	pendingActions  chan *ScheduledAction
	workers         int
	wg              sync.WaitGroup
	ctx             context.Context
	cancel          context.CancelFunc
	timers          map[string]*time.Timer
	timersMu        sync.Mutex
}

// ScheduledAction 调度的动作
type ScheduledAction struct {
	Action    *rule.Action
	Alarm     *engine.ActiveAlarm
	Scheduled time.Time
}

// NewScheduler 创建动作调度器
func NewScheduler(executor *Executor, stateMachine engine.StateMachine, workers int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		executor:       executor,
		stateMachine:   stateMachine,
		pendingActions: make(chan *ScheduledAction, workers*2),
		workers:        workers,
		ctx:            ctx,
		cancel:         cancel,
		timers:         make(map[string]*time.Timer),
	}
}

// Schedule 调度动作执行
func (s *Scheduler) Schedule(action *rule.Action, alarm *engine.ActiveAlarm) error {
	switch action.When {
	case rule.WhenTrigger:
		// 触发时立即执行或延迟执行
		if action.Delay > 0 {
			s.scheduleDelayed(action, alarm, action.Delay)
		} else {
			return s.executeNow(action, alarm)
		}

	case rule.WhenAfterTrigger:
		// 触发后延迟执行
		s.scheduleDelayed(action, alarm, action.Delay)

	case rule.WhenRecover:
		// 恢复时执行
		if action.Delay > 0 {
			s.scheduleDelayed(action, alarm, action.Delay)
		} else {
			return s.executeNow(action, alarm)
		}

	case rule.WhenAfterRecover:
		// 恢复后延迟执行
		s.scheduleDelayed(action, alarm, action.Delay)

	default:
		log.Printf("[调度失败] 未知的执行时机: %d", action.When)
		return fmt.Errorf("未知的执行时机: %d", action.When)
	}

	return nil
}

// ScheduleBatch 批量调度动作
func (s *Scheduler) ScheduleBatch(actions []rule.Action, alarm *engine.ActiveAlarm) error {
	for i := range actions {
		if err := s.Schedule(&actions[i], alarm); err != nil {
			return err
		}
	}
	return nil
}

// scheduleDelayed 延迟调度
func (s *Scheduler) scheduleDelayed(action *rule.Action, alarm *engine.ActiveAlarm, delay time.Duration) {
	timer := time.AfterFunc(delay, func() {
		select {
		case <-s.ctx.Done():
			return
		default:
			if action.Mode == rule.ModeLoop {
				s.executeWithLoop(action, alarm)
			} else {
				s.executeNow(action, alarm)
			}
		}
	})

	// 保存 timer 引用，以便可以取消
	s.timersMu.Lock()
	s.timers[action.ID] = timer
	s.timersMu.Unlock()

	log.Printf("[延迟调度] 动作ID=%s, 延迟=%v", action.ID, delay)
}

// executeNow 立即执行
func (s *Scheduler) executeNow(action *rule.Action, alarm *engine.ActiveAlarm) error {
	log.Printf("[立即执行] 动作ID=%s, 类型=%s", action.ID, action.Type)
	return s.executor.Execute(s.ctx, action, alarm)
}

// executeWithLoop 带循环的执行
func (s *Scheduler) executeWithLoop(action *rule.Action, alarm *engine.ActiveAlarm) {
	log.Printf("[循环开始] 动作ID=%s, 间隔=%v", action.ID, action.LoopDelay)

	ticker := time.NewTicker(action.LoopDelay)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-s.ctx.Done():
			log.Printf("[循环停止-上下文取消] 动作ID=%s", action.ID)
			return

		case <-ticker.C:
			// 执行动作
			if err := s.executor.Execute(s.ctx, action, alarm); err != nil {
				log.Printf("[循环执行失败] 动作ID=%s: %v", action.ID, err)
			}

			count++

			// 检查次数限制
			if action.LoopCount > 0 && count >= action.LoopCount {
				log.Printf("[循环完成-次数达标] 动作ID=%s, 次数=%d", action.ID, count)
				return
			}

			// 检查状态停止条件
			if action.LoopUntil > 0 {
				currentState := s.stateMachine.GetState(alarm.RuleID)
				if currentState == action.LoopUntil {
					log.Printf("[循环停止-状态匹配] 动作ID=%s, 状态=%d",
						action.ID, currentState)
					return
				}
			}
		}
	}
}

// CancelAction 取消动作（停止定时器）
func (s *Scheduler) CancelAction(actionID string) {
	s.timersMu.Lock()
	defer s.timersMu.Unlock()

	if timer, exists := s.timers[actionID]; exists {
		timer.Stop()
		delete(s.timers, actionID)
		log.Printf("[取消动作] 动作ID=%s", actionID)
	}
}

// CancelAllActions 取消所有待执行动作
func (s *Scheduler) CancelAllActions() {
	s.timersMu.Lock()
	defer s.timersMu.Unlock()

	for id, timer := range s.timers {
		timer.Stop()
		log.Printf("[取消动作] 动作ID=%s", id)
	}

	s.timers = make(map[string]*time.Timer)
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	log.Printf("[动作调度器启动] 工作协程数=%d", s.workers)

	// 启动工作协程
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() error {
	log.Printf("[动作调度器停止]")

	// 取消上下文
	s.cancel()

	// 取消所有定时器
	s.CancelAllActions()

	// 等待工作协程退出
	s.wg.Wait()

	return nil
}

// worker 工作协程
func (s *Scheduler) worker(workerID int) {
	defer s.wg.Done()

	log.Printf("[调度工作协程%d启动]", workerID)

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("[调度工作协程%d退出]", workerID)
			return

		case scheduledAction := <-s.pendingActions:
			// 执行动作
			action := scheduledAction.Action
			alarm := scheduledAction.Alarm

			if action.Mode == rule.ModeLoop {
				s.executeWithLoop(action, alarm)
			} else {
				s.executeNow(action, alarm)
			}
		}
	}
}
