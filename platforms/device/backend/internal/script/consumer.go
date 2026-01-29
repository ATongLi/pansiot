package script

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"pansiot-device/internal/consumer"
	"pansiot-device/internal/core"
)

// ScriptConsumer 脚本消费者
type ScriptConsumer struct {
	*consumer.BaseConsumer

	// 核心组件
	engine        *GojaEngine       // Goja 脚本引擎
	vmPool        *VMPool           // VM 池
	sandbox       *Sandbox          // 沙箱管理器
	triggerMgr    *TriggerManager   // 触发器管理器（Phase 2）
	scheduler     *ScriptScheduler  // 周期执行调度器（Phase 2）

	// 脚本管理
	scripts map[string]*Script // scriptID -> Script
	mu      sync.RWMutex

	// 脚本状态
	statuses map[string]*ScriptStatus // scriptID -> status

	// 配置
	config *ScriptConfig

	// 执行队列
	execQueue chan *execTask
}

// execTask 执行任务
type execTask struct {
	scriptID string
	input    map[string]interface{}
	resultCh chan *execResult
	async    bool
}

// execResult 执行结果
type execResult struct {
	output map[string]interface{}
	err    error
}

// NewScriptConsumer 创建脚本消费者
func NewScriptConsumer(id string, storage core.Storage, config *ScriptConfig) *ScriptConsumer {
	if config == nil {
		config = DefaultScriptConfig()
	}

	// 创建基类
	baseConsumer := consumer.NewBaseConsumer(id, "script", storage)

	// 创建 VM 池
	sandbox := NewSandbox(nil) // 暂时传 nil，等 vmPool 创建后再设置
	vmPool := NewVMPool(config.VMPoolSize, config.VMMaxIdle, config.VMMaxLifetime, sandbox)
	vmPool.SetStorage(storage) // 设置存储层
	sandbox.vmPool = vmPool    // 设置双向引用

	// 创建脚本引擎
	engine := NewGojaEngine(vmPool, sandbox)

	return &ScriptConsumer{
		BaseConsumer: baseConsumer,
		engine:       engine,
		vmPool:       vmPool,
		sandbox:      sandbox,
		scripts:      make(map[string]*Script),
		statuses:     make(map[string]*ScriptStatus),
		config:       config,
		execQueue:    make(chan *execTask, config.QueueSize),
	}
}

// Start 启动脚本消费者
func (sc *ScriptConsumer) Start(ctx context.Context) error {
	log.Printf("[脚本消费者] 启动脚本消费者: %s", sc.GetID())

	// 调用基类 Start
	if err := sc.BaseConsumer.Start(ctx); err != nil {
		return err
	}

	// 初始化触发器管理器
	sc.triggerMgr = NewTriggerManager(sc) // sc 实现了 ScriptExecutor 接口
	log.Printf("[脚本消费者] 触发器管理器已初始化")

	// 初始化调度器
	sc.scheduler = NewScriptScheduler(sc) // sc 实现了 ScriptExecutor 接口
	if err := sc.scheduler.Start(); err != nil {
		return fmt.Errorf("启动调度器失败: %v", err)
	}
	log.Printf("[脚本消费者] 调度器已启动")

	// 订阅脚本依赖的变量
	if err := sc.subscribeToScriptVariables(); err != nil {
		log.Printf("[脚本消费者] 警告: 变量订阅失败: %v", err)
	}

	// 启动执行工作协程
	for i := 0; i < sc.config.MaxConcurrent; i++ {
		sc.GetWaitGroup().Add(1)
		go sc.execWorker(i)
	}

	log.Printf("[脚本消费者] 脚本消费者启动成功: %s", sc.GetID())
	return nil
}

// Stop 停止脚本消费者
func (sc *ScriptConsumer) Stop() error {
	log.Printf("[脚本消费者] 停止脚本消费者: %s", sc.GetID())

	// 停止调度器
	if sc.scheduler != nil {
		if err := sc.scheduler.Stop(); err != nil {
			log.Printf("[脚本消费者] 警告: 停止调度器失败: %v", err)
		}
	}

	// 关闭执行队列
	close(sc.execQueue)

	// 等待工作协程退出
	sc.GetWaitGroup().Wait()

	// 关闭 VM 池
	sc.vmPool.Close()

	// 调用基类 Stop
	if err := sc.BaseConsumer.Stop(); err != nil {
		return err
	}

	log.Printf("[脚本消费者] 脚本消费者已停止: %s", sc.GetID())
	return nil
}

// execWorker 执行工作协程
func (sc *ScriptConsumer) execWorker(workerID int) {
	defer sc.GetWaitGroup().Done()

	log.Printf("[脚本消费者] 执行工作协程 %d 启动", workerID)

	for task := range sc.execQueue {
		// 执行脚本
		result := sc.executeTask(task)

		// 返回结果
		if !task.async && task.resultCh != nil {
			task.resultCh <- result
		}

		// 更新统计
		if result.err != nil {
			sc.IncrementFailure()
		} else {
			sc.IncrementSuccess()
		}
	}

	log.Printf("[脚本消费者] 执行工作协程 %d 退出", workerID)
}

// executeTask 执行任务
func (sc *ScriptConsumer) executeTask(task *execTask) *execResult {
	// 获取脚本
	sc.mu.RLock()
	script, exists := sc.scripts[task.scriptID]
	sc.mu.RUnlock()

	if !exists {
		return &execResult{
			err: fmt.Errorf("脚本不存在: %s", task.scriptID),
		}
	}

	// 检查脚本是否启用
	if !script.Enabled {
		return &execResult{
			err: fmt.Errorf("脚本未启用: %s", task.scriptID),
		}
	}

	// 编译脚本
	program, err := sc.engine.Compile(script.ID, script.Content)
	if err != nil {
		return &execResult{
			err: fmt.Errorf("脚本编译失败: %w", err),
		}
	}

	// 执行脚本
	timeout := script.Timeout
	if timeout == 0 {
		timeout = sc.config.DefaultTimeout
	}

	output, err := sc.engine.Execute(script.ID, program, task.input, timeout)

	// 更新状态
	sc.updateStatus(script.ID, err)

	return &execResult{
		output: output,
		err:    err,
	}
}

// updateStatus 更新脚本状态
func (sc *ScriptConsumer) updateStatus(scriptID string, err error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	status, exists := sc.statuses[scriptID]
	if !exists {
		status = &ScriptStatus{
			ScriptID: scriptID,
		}
		sc.statuses[scriptID] = status
	}

	status.LastExecution = time.Now()
	status.ExecCount++

	if err != nil {
		status.ErrorCount++
		status.LastError = err.Error()
		status.State = ScriptStateError
	} else {
		status.State = ScriptStateIdle
	}
}

// LoadScript 加载脚本
func (sc *ScriptConsumer) LoadScript(script *Script) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// 设置创建时间
	script.CreatedAt = time.Now()
	script.UpdatedAt = time.Now()

	// 设置默认超时
	if script.Timeout == 0 {
		script.Timeout = sc.config.DefaultTimeout
	}

	// 保存脚本
	sc.scripts[script.ID] = script

	// 初始化状态
	sc.statuses[script.ID] = &ScriptStatus{
		ScriptID: script.ID,
		Loaded:   true,
		Enabled:  script.Enabled,
		State:    ScriptStateIdle,
	}

	log.Printf("[脚本消费者] 脚本已加载: %s", script.ID)

	// 注册触发器（Phase 2）
	sc.mu.Unlock() // 解锁以避免在registerTriggers中死锁
	err := sc.registerTriggers(script)
	sc.mu.Lock()
	if err != nil {
		log.Printf("[脚本消费者] 警告: 注册触发器失败: %v", err)
	}

	// 重新订阅变量（如果需要）
	if err := sc.subscribeToScriptVariables(); err != nil {
		log.Printf("[脚本消费者] 警告: 变量订阅失败: %v", err)
	}

	return nil
}

// UnloadScript 卸载脚本
func (sc *ScriptConsumer) UnloadScript(scriptID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// 检查脚本是否存在
	if _, exists := sc.scripts[scriptID]; !exists {
		return fmt.Errorf("脚本不存在: %s", scriptID)
	}

	// 注销触发器（Phase 2）
	sc.mu.Unlock() // 解锁以避免在unregisterTriggers中死锁
	err := sc.unregisterTriggers(scriptID)
	sc.mu.Lock()
	if err != nil {
		log.Printf("[脚本消费者] 警告: 注销触发器失败: %v", err)
	}

	// 删除脚本
	delete(sc.scripts, scriptID)
	delete(sc.statuses, scriptID)

	// 删除编译缓存
	sc.engine.RemoveProgram(scriptID)

	log.Printf("[脚本消费者] 脚本已卸载: %s", scriptID)
	return nil
}

// ExecuteScript 执行脚本（同步）
func (sc *ScriptConsumer) ExecuteScript(scriptID string, input map[string]interface{}) (map[string]interface{}, error) {
	// 创建执行任务
	task := &execTask{
		scriptID: scriptID,
		input:    input,
		resultCh: make(chan *execResult, 1),
		async:    false,
	}

	// 提交任务
	select {
	case sc.execQueue <- task:
		// 任务已入队
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("执行队列已满")
	}

	// 等待结果
	result := <-task.resultCh
	return result.output, result.err
}

// ExecuteScriptAsync 执行脚本（异步）
func (sc *ScriptConsumer) ExecuteScriptAsync(scriptID string, input map[string]interface{}) error {
	// 创建执行任务
	task := &execTask{
		scriptID: scriptID,
		input:    input,
		async:    true,
	}

	// 提交任务
	select {
	case sc.execQueue <- task:
		return nil
	default:
		return fmt.Errorf("执行队列已满")
	}
}

// GetScriptStatus 获取脚本状态
func (sc *ScriptConsumer) GetScriptStatus(scriptID string) *ScriptStatus {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	status, exists := sc.statuses[scriptID]
	if !exists {
		return nil
	}

	// 返回副本
	statusCopy := *status
	return &statusCopy
}

// ListScripts 列出所有脚本
func (sc *ScriptConsumer) ListScripts() []*Script {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	scripts := make([]*Script, 0, len(sc.scripts))
	for _, script := range sc.scripts {
		// 返回副本
		scriptCopy := *script
		scripts = append(scripts, &scriptCopy)
	}

	return scripts
}

// GetEngine 获取脚本引擎（供外部使用）
func (sc *ScriptConsumer) GetEngine() *GojaEngine {
	return sc.engine
}

// GetVMPool 获取 VM 池（供外部使用）
func (sc *ScriptConsumer) GetVMPool() *VMPool {
	return sc.vmPool
}

// GetSandbox 获取沙箱（供外部使用）
func (sc *ScriptConsumer) GetSandbox() *Sandbox {
	return sc.sandbox
}

// ============ Phase 2: 触发器和调度器集成 ============

// registerTriggers 注册脚本的所有触发器
func (sc *ScriptConsumer) registerTriggers(script *Script) error {
	if sc.triggerMgr == nil || sc.scheduler == nil {
		return fmt.Errorf("触发器管理器或调度器未初始化")
	}

	for _, triggerConfig := range script.Triggers {
		// 创建触发器
		trigger := &Trigger{
			ID:       triggerConfig.ID,
			Type:     triggerConfig.Type,
			ScriptID: script.ID,
			Enabled:  triggerConfig.Enabled,
		}

		// 根据类型设置条件
		switch triggerConfig.Type {
		case TriggerTypeVariable:
			// 变量触发器
			config, ok := triggerConfig.Config.(*VariableTriggerConfig)
			if !ok {
				return fmt.Errorf("变量触发器配置类型错误")
			}
			trigger.Condition = TriggerCondition{
				VariableID: config.VariableID,
				Operator:   config.Condition,
				Threshold:  config.Threshold,
			}

			// 注册到触发器管理器
			if err := sc.triggerMgr.RegisterTrigger(trigger); err != nil {
				return fmt.Errorf("注册变量触发器失败: %w", err)
			}

		case TriggerTypePeriodic:
			// 周期触发器
			config, ok := triggerConfig.Config.(*PeriodicTriggerConfig)
			if !ok {
				return fmt.Errorf("周期触发器配置类型错误")
			}
			trigger.Condition = TriggerCondition{
				PeriodicConfig: config,
			}

			// 注册到调度器
			if err := sc.scheduler.AddTrigger(trigger); err != nil {
				return fmt.Errorf("注册周期触发器失败: %w", err)
			}

		default:
			log.Printf("[脚本消费者] 警告: 不支持的触发器类型: %d", triggerConfig.Type)
		}

		log.Printf("[脚本消费者] 触发器已注册: %s (脚本: %s, 类型: %d)",
			triggerConfig.ID, script.ID, triggerConfig.Type)
	}

	return nil
}

// unregisterTriggers 注销脚本的所有触发器
func (sc *ScriptConsumer) unregisterTriggers(scriptID string) error {
	if sc.triggerMgr == nil || sc.scheduler == nil {
		return nil // 如果未初始化，直接返回
	}

	// 从触发器管理器注销
	triggers := sc.triggerMgr.GetScriptTriggers(scriptID)
	for _, trigger := range triggers {
		if trigger.Type == TriggerTypeVariable {
			if err := sc.triggerMgr.UnregisterTrigger(trigger.ID); err != nil {
				log.Printf("[脚本消费者] 警告: 注销触发器失败: %v", err)
			}
		}
	}

	// 从调度器注销
	if sc.scheduler != nil {
		schedulerTriggers := sc.scheduler.ListTriggers()
		for _, trigger := range schedulerTriggers {
			if trigger.ScriptID == scriptID {
				if err := sc.scheduler.RemoveTrigger(trigger.ID); err != nil {
					log.Printf("[脚本消费者] 警告: 移除周期触发器失败: %v", err)
				}
			}
		}
	}

	return nil
}

// subscribeToScriptVariables 订阅脚本依赖的变量
func (sc *ScriptConsumer) subscribeToScriptVariables() error {
	if sc.triggerMgr == nil {
		return nil // 如果触发器管理器未初始化，跳过
	}

	// 收集所有需要订阅的变量ID
	variableSet := make(map[uint64]bool)

	sc.mu.RLock()
	for _, script := range sc.scripts {
		// 添加脚本显式声明的变量
		for _, varIDStr := range script.Variables {
			// 将字符串转换为uint64
			var varID uint64
			if _, err := fmt.Sscanf(varIDStr, "%d", &varID); err == nil {
				variableSet[varID] = true
			}
		}

		// 添加触发器中的变量
		triggers := sc.triggerMgr.GetScriptTriggers(script.ID)
		for _, trigger := range triggers {
			if trigger.Type == TriggerTypeVariable {
				variableSet[trigger.Condition.VariableID] = true
			}
		}
	}
	sc.mu.RUnlock()

	// 转换为数组
	variableIDs := make([]uint64, 0, len(variableSet))
	for id := range variableSet {
		variableIDs = append(variableIDs, id)
	}

	if len(variableIDs) == 0 {
		return nil // 没有变量需要订阅
	}

	// 订阅存储层
	storage := sc.GetStorage()
	if err := storage.Subscribe(sc.GetID(), variableIDs, sc.onVariableUpdate); err != nil {
		return fmt.Errorf("订阅变量失败: %w", err)
	}

	log.Printf("[脚本消费者] 已订阅 %d 个变量: %v", len(variableIDs), variableIDs)
	return nil
}

// onVariableUpdate 变量更新回调
func (sc *ScriptConsumer) onVariableUpdate(update core.VariableUpdate) {
	// 通知触发器管理器
	if sc.triggerMgr != nil {
		sc.triggerMgr.onVariableChanged(update)
	}
}

// GetTriggerManager 获取触发器管理器（供外部使用）
func (sc *ScriptConsumer) GetTriggerManager() *TriggerManager {
	return sc.triggerMgr
}

// GetScheduler 获取调度器（供外部使用）
func (sc *ScriptConsumer) GetScheduler() *ScriptScheduler {
	return sc.scheduler
}
