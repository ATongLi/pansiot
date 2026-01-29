package alarm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"pansiot-device/internal/consumer"
	"pansiot-device/internal/alarm/action"
	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/alarm/rule"
	"pansiot-device/internal/alarm/record"
	"pansiot-device/internal/core"
)

// AlarmConsumer 报警消费者
// 整合所有报警组件，实现完整的报警处理流程
type AlarmConsumer struct {
	*consumer.BaseConsumer // 继承基类
	mu       sync.RWMutex   // 保护activeAlarms的读写锁

	// 核心组件
	evaluator     *engine.Evaluator
	stateMachine  *engine.AlarmStateMachine
	ruleManager   *rule.RuleManager
	actionExecutor *action.Executor             // 动作执行器
	recordManager *record.RecordManager        // 记录管理器

	// 活跃报警
	activeAlarms map[string]*engine.ActiveAlarm // ruleID -> 报警实例

	// 配置
	config *AlarmConsumerConfig

	// 工作协程池
	evalChan chan *evalTask    // 评估任务队列
	workers  int               // 工作协程数
	ctx      context.Context    // 上下文
	cancel   context.CancelFunc // 取消函数
}

// evalTask 评估任务
type evalTask struct {
	rule     *rule.AlarmRule
	variable *core.Variable
}

// AlarmConsumerConfig 报警消费者配置
type AlarmConsumerConfig struct {
	EvalWorkers   int           // 评估工作协程数（默认10）
	EvalTimeout   time.Duration // 评估超时（默认100ms）
	RecoverDelay  time.Duration // 恢复确认延迟（默认3s）
	AutoSubscribe bool          // 自动订阅变量
}

// DefaultConfig 默认配置
func DefaultConfig() *AlarmConsumerConfig {
	return &AlarmConsumerConfig{
		EvalWorkers:   10,
		EvalTimeout:   100 * time.Millisecond,
		RecoverDelay:  3 * time.Second,
		AutoSubscribe: true,
	}
}

// NewAlarmConsumer 创建报警消费者
func NewAlarmConsumer(storage core.Storage, cfg *AlarmConsumerConfig) *AlarmConsumer {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// 创建基类（不传递config参数）
	baseConsumer := consumer.NewBaseConsumer(
		"alarm-consumer",
		"alarm",
		storage,
	)

	// 创建评估引擎
	evaluator := engine.NewEvaluator(storage)

	// 创建状态机
	stateMachine := engine.NewAlarmStateMachine(nil)

	// 创建规则管理器
	ruleManager := rule.NewRuleManager(storage)

	// 创建动作执行器
	actionExecutor := action.NewActionExecutor(storage)

	// 创建记录存储
	recordStorage, err := record.NewJSONFileStorage("./data/records/alarm")
	if err != nil {
		log.Printf("创建记录存储失败: %v (将禁用记录功能)", err)
	}

	// 创建记录管理器
	var recordManager *record.RecordManager
	if recordStorage != nil {
		recordConfig := record.DefaultRecordConfig()
		recordManager = record.NewRecordManager(recordStorage, recordConfig)
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())

	ac := &AlarmConsumer{
		BaseConsumer:   baseConsumer,
		evaluator:      evaluator,
		stateMachine:   stateMachine,
		ruleManager:    ruleManager,
		actionExecutor: actionExecutor,
		recordManager:  recordManager,
		activeAlarms:   make(map[string]*engine.ActiveAlarm),
		config:         cfg,
		evalChan:       make(chan *evalTask, cfg.EvalWorkers*2),
		workers:        cfg.EvalWorkers,
		ctx:            ctx,
		cancel:         cancel,
	}

	// 设置配置
	ac.SetConfig(cfg)

	return ac
}

// Start 启动消费者（实现BaseConsumer接口）
func (ac *AlarmConsumer) Start(ctx context.Context) error {
	log.Printf("启动报警消费者")

	// 调用基类Start
	if err := ac.BaseConsumer.Start(ctx); err != nil {
		return err
	}

	// 设置变量订阅
	if ac.config.AutoSubscribe {
		if err := ac.setupSubscriptions(); err != nil {
			return fmt.Errorf("设置变量订阅失败: %w", err)
		}
	}

	// 启动工作协程
	for i := 0; i < ac.workers; i++ {
		ac.GetWaitGroup().Add(1)
		go ac.evalWorker(i)
	}

	// 启动动作执行器
	if err := ac.actionExecutor.Start(ctx); err != nil {
		return fmt.Errorf("启动动作执行器失败: %w", err)
	}

	// 启动记录管理器
	if ac.recordManager != nil {
		if err := ac.recordManager.Start(); err != nil {
			return fmt.Errorf("启动记录管理器失败: %w", err)
		}
		log.Printf("记录管理器启动成功")
	}

	log.Printf("报警消费者启动成功，工作协程数: %d", ac.workers)
	return nil
}

// Stop 停止消费者（实现BaseConsumer接口）
func (ac *AlarmConsumer) Stop() error {
	log.Printf("停止报警消费者")

	// 取消上下文
	ac.cancel()

	// 关闭评估队列
	close(ac.evalChan)

	// 等待工作协程退出
	ac.GetWaitGroup().Wait()

	// 取消订阅
	if ac.config.AutoSubscribe {
		ac.GetStorage().UnsubscribeAll(ac.GetID())
	}

	// 停止动作执行器
	if err := ac.actionExecutor.Stop(); err != nil {
		log.Printf("停止动作执行器失败: %v", err)
	}

	// 停止记录管理器
	if ac.recordManager != nil {
		if err := ac.recordManager.Stop(); err != nil {
			log.Printf("停止记录管理器失败: %v", err)
		}
	}

	// 调用基类Stop
	if err := ac.BaseConsumer.Stop(); err != nil {
		return err
	}

	log.Printf("报警消费者已停止")
	return nil
}

// setupSubscriptions 设置变量订阅
func (ac *AlarmConsumer) setupSubscriptions() error {
	// 获取所有启用的规则
	rules := ac.ruleManager.ListEnabledRules()
	if len(rules) == 0 {
		log.Printf("没有启用的规则，跳过订阅设置")
		return nil
	}

	// 方案选择：规则数少时单独订阅，规则数多时批量订阅
	if len(rules) < 100 {
		// 方案1：每个规则单独订阅
		for _, r := range rules {
			variableIDs := r.GetVariableIDs()
			subscriberID := fmt.Sprintf("alarm-%s", r.ID)
			if err := ac.GetStorage().Subscribe(subscriberID, variableIDs, ac.onVariableUpdate); err != nil {
				return fmt.Errorf("订阅规则 %s 的变量失败: %w", r.ID, err)
			}
		}
		log.Printf("已设置 %d 个规则的订阅", len(rules))
	} else {
		// 方案2：按变量聚合订阅
		// 获取所有涉及的变量ID
		allVariableIDs := make(map[uint64]bool)
		for _, r := range rules {
			for _, vid := range r.GetVariableIDs() {
				allVariableIDs[vid] = true
			}
		}

		// 批量订阅
		variableIDs := make([]uint64, 0, len(allVariableIDs))
		for vid := range allVariableIDs {
			variableIDs = append(variableIDs, vid)
		}

		if err := ac.GetStorage().Subscribe(ac.GetID(), variableIDs, ac.onVariableUpdate); err != nil {
			return fmt.Errorf("批量订阅变量失败: %w", err)
		}
		log.Printf("已批量订阅 %d 个变量", len(variableIDs))
	}

	return nil
}

// onVariableUpdate 变量更新回调（订阅触发）
func (ac *AlarmConsumer) onVariableUpdate(update core.VariableUpdate) {
	// 查找引用此变量的所有规则
	ruleIDs := ac.ruleManager.GetVariableIndex().GetRulesByVariable(update.VariableID)

	if len(ruleIDs) == 0 {
		return
	}

	// 异步评估每个规则
	for _, ruleID := range ruleIDs {
		rule, exists := ac.ruleManager.GetRule(ruleID)
		if !exists || !rule.Enabled {
			continue
		}

		// 将评估任务发送到队列
		task := &evalTask{
			rule:     rule,
			variable: &core.Variable{
				ID:      update.VariableID,
				Value:   update.Value,
				Quality: update.Quality,
			},
		}

		select {
		case ac.evalChan <- task:
			// 任务已入队
		default:
			// 队列满，记录警告
			log.Printf("评估队列已满，丢弃规则 %s 的评估任务", rule.ID)
		}
	}
}

// evalWorker 评估工作协程
func (ac *AlarmConsumer) evalWorker(workerID int) {
	defer ac.GetWaitGroup().Done()

	log.Printf("评估工作协程 %d 启动", workerID)

	for {
		select {
		case <-ac.ctx.Done():
			log.Printf("评估工作协程 %d 退出", workerID)
			return

		case task, ok := <-ac.evalChan:
			if !ok {
				// 队列已关闭
				return
			}

			// 执行评估
			ac.evaluateRule(task.rule)
		}
	}
}

// evaluateRule 评估单个规则
func (ac *AlarmConsumer) evaluateRule(rule *rule.AlarmRule) {
	// 使用评估引擎评估规则
	triggered, err := ac.evaluator.EvaluateRule(rule)
	if err != nil {
		// 检查是否是"使能条件不满足"的错误
		if strings.Contains(err.Error(), "使能条件不满足") {
			// 记录被屏蔽的报警
			ac.recordShieldedAlarm(rule, err.Error())
			return
		}
		log.Printf("评估规则 %s 失败: %v", rule.ID, err)
		return
	}

	// 获取当前状态
	currentState := ac.stateMachine.GetState(rule.ID)

	// 处理评估结果
	if triggered {
		// 条件满足
		if currentState == core.AlarmStateInactive || currentState == core.AlarmStateCleared {
			// 触发报警
			ac.triggerAlarm(rule)
		}
	} else {
		// 条件不满足
		if currentState == core.AlarmStateActive || currentState == core.AlarmStateAcknowledged {
			// 恢复报警
			ac.recoverAlarm(rule)
		}
	}
}

// triggerAlarm 触发报警
func (ac *AlarmConsumer) triggerAlarm(rule *rule.AlarmRule) {
	// 获取触发值（从条件中提取第一个变量的值）
	var triggerValue interface{}
	// 注意：这里使用完整的类型断言路径
	if singleCond, ok := rule.Condition.(interface{ GetVariableID() uint64 }); ok {
		if variable, err := ac.GetStorage().ReadVar(singleCond.GetVariableID()); err == nil {
			triggerValue = variable.Value
		}
	}

	// 转换状态
	err := ac.stateMachine.TransitionTo(rule.ID, core.AlarmStateActive, "trigger", "")
	if err != nil {
		log.Printf("状态转换失败: %v", err)
		return
	}

	// 创建活跃报警实例
	activeAlarm := &engine.ActiveAlarm{
		RuleID:       rule.ID,
		Rule:         rule,
		State:        core.AlarmStateActive,
		TriggerTime:  time.Now(),
		TriggerValue: triggerValue,
	}

	// 存储活跃报警
	ac.mu.Lock()
	ac.activeAlarms[rule.ID] = activeAlarm
	ac.mu.Unlock()

	// 记录日志
	log.Printf("报警触发: 规则=%s, 级别=%d, 值=%v", rule.Name, rule.Level, triggerValue)

	// 记录触发事件
	if ac.recordManager != nil {
		if err := ac.recordManager.RecordTrigger(activeAlarm); err != nil {
			log.Printf("记录触发事件失败: %v", err)
		}
	}

	// 执行触发动作（异步）
	if len(rule.TriggerActions) > 0 {
		go func() {
			if err := ac.actionExecutor.ExecuteBatch(context.Background(), rule.TriggerActions, activeAlarm); err != nil {
				log.Printf("执行触发动作失败: %v", err)
			}
		}()
	}

	// TODO: 发送通知
	// ac.sendNotification(rule, activeAlarm)
}

// recoverAlarm 恢复报警
func (ac *AlarmConsumer) recoverAlarm(rule *rule.AlarmRule) {
	// 启动恢复确认延迟
	go func() {
		// 等待恢复确认延迟
		time.Sleep(ac.config.RecoverDelay)

		// 再次检查条件是否仍然不满足
		triggered, err := ac.evaluator.EvaluateRule(rule)
		if err != nil {
			log.Printf("恢复验证失败: %v", err)
			return
		}

		if triggered {
			// 条件重新满足，取消恢复
			return
		}

		// 转换状态
		err = ac.stateMachine.TransitionTo(rule.ID, core.AlarmStateCleared, "recover", "")
		if err != nil {
			log.Printf("状态转换失败: %v", err)
			return
		}

		// 更新活跃报警
		ac.mu.Lock()
		if alarm, exists := ac.activeAlarms[rule.ID]; exists {
			now := time.Now()
			alarm.State = core.AlarmStateCleared
			alarm.RecoverTime = &now
		}
		ac.mu.Unlock()

		// 记录日志
		log.Printf("报警恢复: 规则=%s", rule.Name)

		// 记录恢复事件
		if ac.recordManager != nil {
			if err := ac.recordManager.RecordRecover(rule.ID); err != nil {
				log.Printf("记录恢复事件失败: %v", err)
			}
		}

		// 执行恢复动作（异步）
		if len(rule.RecoverActions) > 0 {
			go func() {
				// 获取活跃报警
				ac.mu.RLock()
				alarm := ac.activeAlarms[rule.ID]
				ac.mu.RUnlock()

				if alarm != nil {
					if err := ac.actionExecutor.ExecuteBatch(context.Background(), rule.RecoverActions, alarm); err != nil {
						log.Printf("执行恢复动作失败: %v", err)
					}
				}
			}()
		}

		// TODO: 发送恢复通知
		// ac.sendRecoveryNotification(rule)
	}()
}

// recordShieldedAlarm 记录被屏蔽的报警
// 当规则的使能条件不满足时调用此方法，记录屏蔽事件但不执行报警动作
func (ac *AlarmConsumer) recordShieldedAlarm(rule *rule.AlarmRule, reason string) {
	if ac.recordManager == nil {
		return
	}

	// 提取触发值（尝试读取第一个变量的值）
	var triggerValue interface{}
	if singleCond, ok := rule.Condition.(interface{ GetVariableID() uint64 }); ok {
		if variable, err := ac.GetStorage().ReadVar(singleCond.GetVariableID()); err == nil {
			triggerValue = variable.Value
		}
	}

	// 创建 ActiveAlarm 对象
	activeAlarm := &engine.ActiveAlarm{
		RuleID:       rule.ID,
		Rule:         rule,
		State:        core.AlarmStateActive, // 仍标记为 Active
		TriggerTime:  time.Now(),
		TriggerValue: triggerValue,
	}

	// 记录被屏蔽的报警
	if err := ac.recordManager.RecordShielded(activeAlarm); err != nil {
		log.Printf("[报警屏蔽] 记录失败: 规则=%s, 错误=%v", rule.ID, err)
	} else {
		log.Printf("[报警屏蔽] 规则=%s, 级别=%d, 原因=%s", rule.Name, rule.Level, reason)
	}
}

// acknowledgeAlarm 确认报警
func (ac *AlarmConsumer) acknowledgeAlarm(ruleID string, userID string) error {
	// 检查当前状态
	currentState := ac.stateMachine.GetState(ruleID)
	if currentState != core.AlarmStateActive {
		return fmt.Errorf("只能确认激活状态的报警，当前状态: %d", currentState)
	}

	// 转换状态
	err := ac.stateMachine.TransitionTo(ruleID, core.AlarmStateAcknowledged, "acknowledge", userID)
	if err != nil {
		return fmt.Errorf("状态转换失败: %w", err)
	}

	// 更新活跃报警
	ac.mu.Lock()
	if alarm, exists := ac.activeAlarms[ruleID]; exists {
		alarm.State = core.AlarmStateAcknowledged
		alarm.AckUser = userID
		now := time.Now()
		alarm.AckTime = &now
	}
	ac.mu.Unlock()

	// 记录日志
	log.Printf("报警已确认: 规则=%s, 用户=%s", ruleID, userID)

	return nil
}

// getActiveAlarm 获取活跃报警
func (ac *AlarmConsumer) getActiveAlarm(ruleID string) (*engine.ActiveAlarm, bool) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	alarm, exists := ac.activeAlarms[ruleID]
	if !exists {
		return nil, false
	}

	// 返回副本
	alarmCopy := *alarm
	return &alarmCopy, true
}

// GetActiveAlarms 获取所有活跃报警
func (ac *AlarmConsumer) GetActiveAlarms() []*engine.ActiveAlarm {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	alarms := make([]*engine.ActiveAlarm, 0, len(ac.activeAlarms))
	for _, alarm := range ac.activeAlarms {
		// 返回副本
		alarmCopy := *alarm
		alarms = append(alarms, &alarmCopy)
	}

	return alarms
}

// GetStateMachine 获取状态机（供外部使用）
func (ac *AlarmConsumer) GetStateMachine() *engine.AlarmStateMachine {
	return ac.stateMachine
}

// GetRuleManager 获取规则管理器（供外部使用）
func (ac *AlarmConsumer) GetRuleManager() *rule.RuleManager {
	return ac.ruleManager
}

// GetEvaluator 获取评估引擎（供外部使用）
func (ac *AlarmConsumer) GetEvaluator() *engine.Evaluator {
	return ac.evaluator
}

// GetRecordManager 获取记录管理器（供外部使用）
func (ac *AlarmConsumer) GetRecordManager() *record.RecordManager {
	return ac.recordManager
}

// AddRule 添加规则
func (ac *AlarmConsumer) AddRule(rule *rule.AlarmRule) error {
	if err := ac.ruleManager.AddRule(rule); err != nil {
		return err
	}

	// 如果设置了自动订阅，重新设置订阅
	if ac.config.AutoSubscribe && ac.IsRunning() {
		ac.GetStorage().UnsubscribeAll(ac.GetID())
		if err := ac.setupSubscriptions(); err != nil {
			return fmt.Errorf("重新设置订阅失败: %w", err)
		}
	}

	return nil
}

// RemoveRule 删除规则
func (ac *AlarmConsumer) RemoveRule(ruleID string) error {
	if err := ac.ruleManager.RemoveRule(ruleID); err != nil {
		return err
	}

	// 清除状态
	ac.stateMachine.Clear(ruleID)

	// 删除活跃报警
	ac.mu.Lock()
	delete(ac.activeAlarms, ruleID)
	ac.mu.Unlock()

	// 如果设置了自动订阅，重新设置订阅
	if ac.config.AutoSubscribe && ac.IsRunning() {
		ac.GetStorage().UnsubscribeAll(ac.GetID())
		if err := ac.setupSubscriptions(); err != nil {
			return fmt.Errorf("重新设置订阅失败: %w", err)
		}
	}

	return nil
}
