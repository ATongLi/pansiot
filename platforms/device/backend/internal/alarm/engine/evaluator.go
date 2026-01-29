package engine

import (
	"fmt"
	"sync"

	"pansiot-device/internal/core"
	"pansiot-device/internal/alarm/rule"
)

// Evaluator 评估引擎
// 整合状态跟踪、死区过滤、边沿检测、延迟跟踪
type Evaluator struct {
	storage         core.Storage
	stateTracker    *StateTracker
	deadbandFilter  *DeadbandFilter
	edgeDetector    *EdgeDetector
	delayTracker    *DelayTracker
	alarmStateTracker *AlarmStateTracker
	mu              sync.RWMutex
}

// NewEvaluator 创建评估引擎
func NewEvaluator(storage core.Storage) *Evaluator {
	st := NewStateTracker()
	return &Evaluator{
		storage:          storage,
		stateTracker:     st,
		deadbandFilter:   NewDeadbandFilter(st),
		edgeDetector:     NewEdgeDetector(st),
		delayTracker:     NewDelayTracker(),
		alarmStateTracker: NewAlarmStateTracker(),
	}
}

// EvaluateWithState 评估条件（带状态管理）
// 完整的评估流程：读取变量 → 应用死区 → 执行比较 → 检测边沿 → 处理延迟
func (e *Evaluator) EvaluateWithState(cond rule.Condition) (satisfied bool, err error) {
	switch c := cond.(type) {
	case *rule.SingleCondition:
		return e.evaluateSingleCondition(c)
	case *rule.ConditionGroup:
		return e.evaluateConditionGroup(c)
	default:
		return false, fmt.Errorf("不支持的条件类型: %T", cond)
	}
}

// evaluateSingleCondition 评估单个条件（完整版）
func (e *Evaluator) evaluateSingleCondition(cond *rule.SingleCondition) (bool, error) {
	// 1. 读取变量值
	variable, err := e.storage.ReadVar(cond.VariableID)
	if err != nil {
		return false, fmt.Errorf("读取变量[%d]失败: %w", cond.VariableID, err)
	}

	// 2. 检查数据质量
	if variable.Quality != core.QualityGood {
		return false, fmt.Errorf("变量[%d]数据质量不佳: %v", cond.VariableID, variable.Quality)
	}

	// 3. 边沿检测（如果是边沿操作符）- 在阈值处理之前
	if cond.Operator >= rule.OpRise && cond.Operator <= rule.OpFall {
		return e.evaluateEdgeCondition(cond, variable.Value)
	}

	// 4. 获取阈值
	threshold := cond.Value
	if cond.ValueVarID != nil {
		// 动态阈值：从变量读取
		refVar, err := e.storage.ReadVar(*cond.ValueVarID)
		if err != nil {
			return false, fmt.Errorf("读取阈值变量[%d]失败: %w", *cond.ValueVarID, err)
		}
		threshold = refVar.Value
	}

	// 5. 转换为float64（用于数值比较）
	currentValue, err := toFloat64(variable.Value)
	if err != nil {
		return false, fmt.Errorf("变量值转换失败: %w", err)
	}

	thresholdFloat, err := toFloat64(threshold)
	if err != nil {
		return false, fmt.Errorf("阈值转换失败: %w", err)
	}

	// 6. 应用死区逻辑
	satisfied, _ := e.deadbandFilter.ApplyDeadband(cond, currentValue, thresholdFloat)

	// 7. 更新状态
	e.stateTracker.UpdateState(cond.VariableID, variable.Value, satisfied)

	return satisfied, nil
}

// evaluateConditionGroup 评估条件组
func (e *Evaluator) evaluateConditionGroup(group *rule.ConditionGroup) (bool, error) {
	if len(group.Conditions) == 0 {
		return false, fmt.Errorf("条件组不能为空")
	}

	switch group.Logic {
	case rule.LogicAND:
		// AND逻辑：所有条件都为true才返回true
		for _, cond := range group.Conditions {
			result, err := e.EvaluateWithState(cond)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil // 短路
			}
		}
		return true, nil

	case rule.LogicOR:
		// OR逻辑：任意条件为true就返回true
		for _, cond := range group.Conditions {
			result, err := e.EvaluateWithState(cond)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil // 短路
			}
		}
		return false, nil

	default:
		return false, fmt.Errorf("未知的逻辑操作符: %d", group.Logic)
	}
}

// evaluateEdgeCondition 评估边沿条件
func (e *Evaluator) evaluateEdgeCondition(cond *rule.SingleCondition, currentValue interface{}) (bool, error) {
	// 使用边沿检测器
	edgeDetected, edgeType := e.edgeDetector.DetectEdge(cond, currentValue)

	switch cond.Operator {
	case rule.OpRise:
		// 0→1 上升沿
		return edgeDetected && edgeType == "rising", nil

	case rule.OpFall:
		// 1→0 下降沿
		return edgeDetected && edgeType == "falling", nil

	default:
		return false, fmt.Errorf("不支持的边沿操作符: %d", cond.Operator)
	}
}

// EvaluateRule 评估报警规则（完整流程）
// 包含使能条件检查、主条件评估、延迟处理
func (e *Evaluator) EvaluateRule(alarmRule *rule.AlarmRule) (triggered bool, err error) {
	// 1. 检查规则是否启用
	if !alarmRule.Enabled {
		return false, nil
	}

	// 2. 评估使能条件（如果有）
	if alarmRule.EnableCond != nil {
		enabled, err := e.EvaluateWithState(alarmRule.EnableCond)
		if err != nil {
			return false, fmt.Errorf("使能条件评估失败: %w", err)
		}
		if !enabled {
			return false, fmt.Errorf("使能条件不满足: 规则=%s", alarmRule.ID)
		}
	}

	// 3. 评估主条件
	satisfied, err := e.EvaluateWithState(alarmRule.Condition)
	if err != nil {
		return false, fmt.Errorf("主条件评估失败: %w", err)
	}

	// 4. 检查是否有延迟
	hasDelay := false
	if singleCond, ok := alarmRule.Condition.(*rule.SingleCondition); ok && singleCond.Delay > 0 {
		hasDelay = true
	}

	if !hasDelay {
		// 无延迟，直接返回结果
		return satisfied, nil
	}

	// 5. 处理延迟逻辑
	return e.handleDelay(alarmRule, satisfied)
}

// handleDelay 处理延迟逻辑
func (e *Evaluator) handleDelay(alarmRule *rule.AlarmRule, satisfied bool) (bool, error) {
	// 只处理单个条件的延迟
	singleCond, ok := alarmRule.Condition.(*rule.SingleCondition)
	if !ok || singleCond.Delay <= 0 {
		return satisfied, nil
	}

	variableID := singleCond.VariableID

	if satisfied {
		// 条件满足，启动延迟定时器
		if !e.delayTracker.IsActive(alarmRule.ID) {
			// 启动新的延迟定时器
			err := e.delayTracker.StartDelay(alarmRule.ID, singleCond, variableID, nil, func() {
				// 延迟到期，触发报警
				e.alarmStateTracker.SetState(alarmRule.ID, core.AlarmStateActive)
			})
			if err != nil {
				return false, fmt.Errorf("启动延迟定时器失败: %w", err)
			}
		}
		return false, nil // 延迟期间不触发
	} else {
		// 条件不满足，取消延迟定时器
		e.delayTracker.Cancel(alarmRule.ID)
		return false, nil
	}
}

// GetAlarmState 获取报警状态
func (e *Evaluator) GetAlarmState(ruleID string) core.AlarmState {
	return e.alarmStateTracker.GetState(ruleID)
}

// SetAlarmState 设置报警状态
func (e *Evaluator) SetAlarmState(ruleID string, state core.AlarmState) {
	e.alarmStateTracker.SetState(ruleID, state)
}

// GetConditionState 获取条件状态
func (e *Evaluator) GetConditionState(variableID uint64) *ConditionState {
	return e.stateTracker.GetState(variableID)
}

// ClearRule 清除规则相关的所有状态
func (e *Evaluator) ClearRule(ruleID string) {
	e.alarmStateTracker.Clear(ruleID)
	e.delayTracker.Cancel(ruleID)
	// 注意：条件状态由多个规则共享，不应轻易清除
	// 如需清除，可手动调用 stateTracker.ClearVariable(variableID)
}

// GetStateTracker 获取状态跟踪器
func (e *Evaluator) GetStateTracker() *StateTracker {
	return e.stateTracker
}

// GetDeadbandFilter 获取死区过滤器
func (e *Evaluator) GetDeadbandFilter() *DeadbandFilter {
	return e.deadbandFilter
}

// GetEdgeDetector 获取边沿检测器
func (e *Evaluator) GetEdgeDetector() *EdgeDetector {
	return e.edgeDetector
}

// GetDelayTracker 获取延迟跟踪器
func (e *Evaluator) GetDelayTracker() *DelayTracker {
	return e.delayTracker
}

// GetAlarmStateTracker 获取报警状态跟踪器
func (e *Evaluator) GetAlarmStateTracker() *AlarmStateTracker {
	return e.alarmStateTracker
}

// Reset 重置所有状态
func (e *Evaluator) Reset() {
	e.stateTracker.Clear()
	e.alarmStateTracker.ClearAll()
	e.delayTracker.Clear()
}

// toFloat64 转换为float64（辅助函数）
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("不支持的类型: %T", value)
	}
}
