package script

import (
	"fmt"
	"log"
	"sync"
	"time"

	"pansiot-device/internal/core"
)

// TriggerManager 触发器管理器
// 管理所有触发器的注册、评估和执行
type TriggerManager struct {
	mu               sync.RWMutex
	variableTriggers map[uint64][]*Trigger      // 变量ID -> 触发器列表
	scriptTriggers   map[string][]*Trigger      // 脚本ID -> 触发器列表
	allTriggers      map[string]*Trigger        // 触发器ID -> 触发器
	executor         ScriptExecutor             // 脚本执行器
}

// Trigger 触发器定义
type Trigger struct {
	ID            string             // 触发器唯一ID
	Type          TriggerType        // 触发器类型
	ScriptID      string             // 关联的脚本ID
	Condition     TriggerCondition   // 触发条件
	Enabled       bool               // 是否启用
	LastTriggered time.Time          // 上次触发时间
	TriggerCount  int64              // 触发次数统计
	mu            sync.Mutex         // 保护统计字段
}

// TriggerCondition 触发条件
type TriggerCondition struct {
	// 变量触发条件
	VariableID uint64      // 变量ID
	Operator   string      // ==, !=, >, <, >=, <=
	Threshold  interface{} // 阈值

	// 周期触发条件
	PeriodicConfig *PeriodicTriggerConfig // 周期配置
}

// NewTriggerManager 创建触发器管理器
func NewTriggerManager(executor ScriptExecutor) *TriggerManager {
	return &TriggerManager{
		variableTriggers: make(map[uint64][]*Trigger),
		scriptTriggers:   make(map[string][]*Trigger),
		allTriggers:      make(map[string]*Trigger),
		executor:         executor,
	}
}

// RegisterTrigger 注册触发器
func (tm *TriggerManager) RegisterTrigger(trigger *Trigger) error {
	if trigger.ID == "" {
		return fmt.Errorf("触发器ID不能为空")
	}

	if trigger.ScriptID == "" {
		return fmt.Errorf("脚本ID不能为空")
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 检查触发器ID是否已存在
	if _, exists := tm.allTriggers[trigger.ID]; exists {
		return fmt.Errorf("触发器已存在: %s", trigger.ID)
	}

	// 添加到全局索引
	tm.allTriggers[trigger.ID] = trigger

	// 添加到脚本索引
	tm.scriptTriggers[trigger.ScriptID] = append(tm.scriptTriggers[trigger.ScriptID], trigger)

	// 根据类型添加到对应索引
	switch trigger.Type {
	case TriggerTypeVariable:
		// 变量触发器：添加到变量索引
		if trigger.Condition.VariableID == 0 {
			return fmt.Errorf("变量触发器必须指定变量ID")
		}
		tm.variableTriggers[trigger.Condition.VariableID] = append(
			tm.variableTriggers[trigger.Condition.VariableID],
			trigger,
		)

	case TriggerTypePeriodic:
		// 周期触发器：由调度器管理，这里只记录
		log.Printf("[TriggerManager] 注册周期触发器: %s", trigger.ID)

	default:
		log.Printf("[TriggerManager] 未知的触发器类型: %d", trigger.Type)
	}

	log.Printf("[TriggerManager] 触发器已注册: %s (类型: %d, 脚本: %s)",
		trigger.ID, trigger.Type, trigger.ScriptID)

	return nil
}

// UnregisterTrigger 注销触发器
func (tm *TriggerManager) UnregisterTrigger(triggerID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 查找触发器
	trigger, exists := tm.allTriggers[triggerID]
	if !exists {
		return fmt.Errorf("触发器不存在: %s", triggerID)
	}

	// 从全局索引删除
	delete(tm.allTriggers, triggerID)

	// 从脚本索引删除
	triggers := tm.scriptTriggers[trigger.ScriptID]
	for i, t := range triggers {
		if t.ID == triggerID {
			tm.scriptTriggers[trigger.ScriptID] = append(triggers[:i], triggers[i+1:]...)
			break
		}
	}

	// 根据类型从对应索引删除
	switch trigger.Type {
	case TriggerTypeVariable:
		varTriggers := tm.variableTriggers[trigger.Condition.VariableID]
		for i, t := range varTriggers {
			if t.ID == triggerID {
				tm.variableTriggers[trigger.Condition.VariableID] = append(varTriggers[:i], varTriggers[i+1:]...)
				break
			}
		}

	case TriggerTypePeriodic:
		log.Printf("[TriggerManager] 注销周期触发器: %s", triggerID)
	}

	log.Printf("[TriggerManager] 触发器已注销: %s", triggerID)
	return nil
}

// onVariableChanged 变量变化回调
// 当存储层变量变化时调用
func (tm *TriggerManager) onVariableChanged(update core.VariableUpdate) {
	tm.mu.RLock()
	triggers := tm.variableTriggers[update.VariableID]
	tm.mu.RUnlock()

	if len(triggers) == 0 {
		return
	}

	// 检查每个触发器
	for _, trigger := range triggers {
		if !trigger.Enabled {
			continue
		}

		// 评估条件
		if tm.evaluateCondition(trigger.Condition, update.Value) {
			// 异步执行脚本，不阻塞主流程
			go tm.trigger(trigger, map[string]interface{}{
				"trigger_type": "variable",
				"trigger_id":   trigger.ID,
				"variable_id":  update.VariableID,
				"value":        update.Value,
				"timestamp":    update.Timestamp,
			})
		}
	}
}

// evaluateCondition 评估触发条件
func (tm *TriggerManager) evaluateCondition(condition TriggerCondition, value interface{}) bool {
	threshold := condition.Threshold

	// 比较操作
	result := tm.compare(value, threshold, condition.Operator)

	log.Printf("[TriggerManager] 条件评估: 值=%v (%T) 阈值=%v (%T) 操作=%s 结果=%v",
		value, value, threshold, threshold, condition.Operator, result)

	return result
}

// compare 比较两个值
func (tm *TriggerManager) compare(value, threshold interface{}, operator string) bool {
	// 尝试转换为float64进行比较
	var vFloat, tFloat float64
	var vOk, tOk bool

	// 类型转换
	switch v := value.(type) {
	case int:
		vFloat, vOk = float64(v), true
	case int8:
		vFloat, vOk = float64(v), true
	case int16:
		vFloat, vOk = float64(v), true
	case int32:
		vFloat, vOk = float64(v), true
	case int64:
		vFloat, vOk = float64(v), true
	case uint:
		vFloat, vOk = float64(v), true
	case uint8:
		vFloat, vOk = float64(v), true
	case uint16:
		vFloat, vOk = float64(v), true
	case uint32:
		vFloat, vOk = float64(v), true
	case uint64:
		vFloat, vOk = float64(v), true
	case float32:
		vFloat, vOk = float64(v), true
	case float64:
		vFloat, vOk = v, true
	case bool:
		vFloat, vOk = func() (float64, bool) {
			if v {
				return 1, true
			}
			return 0, true
		}()
	}

	switch t := threshold.(type) {
	case int:
		tFloat, tOk = float64(t), true
	case int8:
		tFloat, tOk = float64(t), true
	case int16:
		tFloat, tOk = float64(t), true
	case int32:
		tFloat, tOk = float64(t), true
	case int64:
		tFloat, tOk = float64(t), true
	case uint:
		tFloat, tOk = float64(t), true
	case uint8:
		tFloat, tOk = float64(t), true
	case uint16:
		tFloat, tOk = float64(t), true
	case uint32:
		tFloat, tOk = float64(t), true
	case uint64:
		tFloat, tOk = float64(t), true
	case float32:
		tFloat, tOk = float64(t), true
	case float64:
		tFloat, tOk = t, true
	case bool:
		tFloat, tOk = func() (float64, bool) {
			if t {
				return 1, true
			}
			return 0, true
		}()
	}

	// 如果都能转换为数字，使用数字比较
	if vOk && tOk {
		switch operator {
		case "==":
			return vFloat == tFloat
		case "!=":
			return vFloat != tFloat
		case ">":
			return vFloat > tFloat
		case "<":
			return vFloat < tFloat
		case ">=":
			return vFloat >= tFloat
		case "<=":
			return vFloat <= tFloat
		default:
			log.Printf("[TriggerManager] 未知的操作符: %s", operator)
			return false
		}
	}

	// 无法转换为数字，尝试直接比较
	switch operator {
	case "==":
		return value == threshold
	case "!=":
		return value != threshold
	default:
		log.Printf("[TriggerManager] 不支持的操作符（非数字类型）: %s", operator)
		return false
	}
}

// trigger 触发脚本执行
func (tm *TriggerManager) trigger(trigger *Trigger, inputData map[string]interface{}) error {
	// 更新触发统计
	trigger.mu.Lock()
	trigger.LastTriggered = time.Now()
	trigger.TriggerCount++
	count := trigger.TriggerCount
	trigger.mu.Unlock()

	log.Printf("[TriggerManager] 触发器触发: %s (脚本: %s, 次数: %d)",
		trigger.ID, trigger.ScriptID, count)

	// 执行脚本
	err := tm.executor.ExecuteScriptAsync(trigger.ScriptID, inputData)
	if err != nil {
		log.Printf("[TriggerManager] 脚本执行失败: %s, 错误: %v",
			trigger.ScriptID, err)
		return err
	}

	return nil
}

// GetTriggerInfo 获取触发器信息
func (tm *TriggerManager) GetTriggerInfo(triggerID string) (*Trigger, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	trigger, exists := tm.allTriggers[triggerID]
	if !exists {
		return nil, fmt.Errorf("触发器不存在: %s", triggerID)
	}

	// 返回副本
	triggerCopy := *trigger
	triggerCopy.Condition = trigger.Condition
	return &triggerCopy, nil
}

// GetScriptTriggers 获取脚本的所有触发器
func (tm *TriggerManager) GetScriptTriggers(scriptID string) []*Trigger {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	triggers := tm.scriptTriggers[scriptID]
	if triggers == nil {
		return nil
	}

	// 返回副本
	result := make([]*Trigger, len(triggers))
	for i, t := range triggers {
		triggerCopy := *t
		triggerCopy.Condition = t.Condition
		result[i] = &triggerCopy
	}

	return result
}

// EnableTrigger 启用触发器
func (tm *TriggerManager) EnableTrigger(triggerID string) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	trigger, exists := tm.allTriggers[triggerID]
	if !exists {
		return fmt.Errorf("触发器不存在: %s", triggerID)
	}

	trigger.Enabled = true
	log.Printf("[TriggerManager] 触发器已启用: %s", triggerID)
	return nil
}

// DisableTrigger 禁用触发器
func (tm *TriggerManager) DisableTrigger(triggerID string) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	trigger, exists := tm.allTriggers[triggerID]
	if !exists {
		return fmt.Errorf("触发器不存在: %s", triggerID)
	}

	trigger.Enabled = false
	log.Printf("[TriggerManager] 触发器已禁用: %s", triggerID)
	return nil
}

// GetTriggerStats 获取触发器统计信息
func (tm *TriggerManager) GetTriggerStats(triggerID string) (lastTriggered time.Time, count int64, err error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	trigger, exists := tm.allTriggers[triggerID]
	if !exists {
		return time.Time{}, 0, fmt.Errorf("触发器不存在: %s", triggerID)
	}

	trigger.mu.Lock()
	lastTriggered = trigger.LastTriggered
	count = trigger.TriggerCount
	trigger.mu.Unlock()

	return lastTriggered, count, nil
}

// ListAllTriggers 列出所有触发器
func (tm *TriggerManager) ListAllTriggers() []*Trigger {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	result := make([]*Trigger, 0, len(tm.allTriggers))
	for _, t := range tm.allTriggers {
		triggerCopy := *t
		triggerCopy.Condition = t.Condition
		result = append(result, &triggerCopy)
	}

	return result
}
