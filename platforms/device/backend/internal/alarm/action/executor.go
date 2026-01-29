package action

import (
	"context"
	"fmt"
	"log"
	"sync"

	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/alarm/rule"
	"pansiot-device/internal/core"
)

// ActionExecutor 动作执行器接口
type ActionExecutor interface {
	// Execute 执行单个动作
	Execute(ctx context.Context, action *rule.Action, alarm *engine.ActiveAlarm) error

	// ExecuteBatch 批量执行动作
	ExecuteBatch(ctx context.Context, actions []rule.Action, alarm *engine.ActiveAlarm) error

	// Start 启动执行器
	Start(ctx context.Context) error

	// Stop 停止执行器
	Stop() error
}

// ActionHandler 具体动作处理器接口
type ActionHandler interface {
	// Handle 处理动作
	Handle(ctx context.Context, action *rule.Action, alarm *engine.ActiveAlarm) error

	// Validate 验证动作参数
	Validate(action *rule.Action) error
}

// Executor 动作执行器实现
type Executor struct {
	mu       sync.RWMutex
	handlers map[rule.ActionType]ActionHandler // 动作类型 -> 处理器
	storage  core.Storage
	running  bool
}

// NewActionExecutor 创建动作执行器
func NewActionExecutor(storage core.Storage) *Executor {
	e := &Executor{
		handlers: make(map[rule.ActionType]ActionHandler),
		storage:  storage,
	}

	// 注册所有动作处理器
	e.RegisterHandler(rule.ActSound, &SoundPlayerHandler{})
	e.RegisterHandler(rule.ActJumpPage, &PageJumperHandler{})
	e.RegisterHandler(rule.ActWriteVar, &VariableWriterHandler{storage: storage})
	e.RegisterHandler(rule.ActPopup, &PopupNotifierHandler{})

	return e
}

// RegisterHandler 注册动作处理器
func (e *Executor) RegisterHandler(actionType rule.ActionType, handler ActionHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[actionType] = handler
}

// Execute 执行单个动作
func (e *Executor) Execute(ctx context.Context, action *rule.Action, alarm *engine.ActiveAlarm) error {
	e.mu.RLock()
	handler, exists := e.handlers[action.Type]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("不支持的动作类型: %d", action.Type)
	}

	// 验证动作参数
	if err := handler.Validate(action); err != nil {
		return fmt.Errorf("动作参数验证失败: %w", err)
	}

	// 执行动作（带 panic 恢复）
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[动作执行panic] 动作ID=%s: %v", action.ID, r)
		}
	}()

	log.Printf("[执行动作] 类型=%s, 规则=%s", action.Type, alarm.RuleID)
	return handler.Handle(ctx, action, alarm)
}

// ExecuteBatch 批量执行动作（并发执行）
func (e *Executor) ExecuteBatch(ctx context.Context, actions []rule.Action, alarm *engine.ActiveAlarm) error {
	if len(actions) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(actions))

	for _, action := range actions {
		wg.Add(1)
		go func(act rule.Action) {
			defer wg.Done()
			if err := e.Execute(ctx, &act, alarm); err != nil {
				errChan <- fmt.Errorf("动作 %s 执行失败: %w", act.ID, err)
			}
		}(action)
	}

	wg.Wait()
	close(errChan)

	// 收集错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("批量执行完成，但 %d 个动作失败: %v", len(errors), errors)
	}

	log.Printf("[批量执行完成] 成功执行 %d 个动作", len(actions))
	return nil
}

// Start 启动执行器
func (e *Executor) Start(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return fmt.Errorf("执行器已在运行")
	}

	e.running = true
	log.Printf("[动作执行器启动]")
	return nil
}

// Stop 停止执行器
func (e *Executor) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.running {
		return nil
	}

	e.running = false
	log.Printf("[动作执行器停止]")
	return nil
}
