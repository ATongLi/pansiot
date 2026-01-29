package action

import (
	"context"
	"fmt"
	"log"
	"time"

	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/alarm/rule"
	"pansiot-device/internal/core"
)

// VariableWriterHandler 变量写值处理器（真实功能）
type VariableWriterHandler struct {
	storage core.Storage
}

// Handle 执行变量写值
func (h *VariableWriterHandler) Handle(ctx context.Context, action *rule.Action, alarm *engine.ActiveAlarm) error {
	// 参数验证
	variableID, ok := action.Params["variable_id"].(uint64)
	if !ok {
		return fmt.Errorf("缺少 variable_id 参数")
	}

	value, ok := action.Params["value"]
	if !ok {
		return fmt.Errorf("缺少 value 参数")
	}

	// 构造 Variable 对象
	variable := &core.Variable{
		ID:        variableID,
		Value:     value,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
	}

	// 写入存储层
	if err := h.storage.WriteVar(variable); err != nil {
		return fmt.Errorf("写入变量失败: %w", err)
	}

	log.Printf("[变量写值成功] 变量ID=%d, 值=%v, 规则=%s",
		variableID, value, alarm.RuleID)
	return nil
}

// Validate 验证动作参数
func (h *VariableWriterHandler) Validate(action *rule.Action) error {
	// 检查 variable_id 参数
	variableID, ok := action.Params["variable_id"]
	if !ok {
		return fmt.Errorf("缺少必需参数: variable_id")
	}

	_, ok = variableID.(uint64)
	if !ok {
		return fmt.Errorf("参数 variable_id 类型错误，期望 uint64")
	}

	// 检查 value 参数
	_, ok = action.Params["value"]
	if !ok {
		return fmt.Errorf("缺少必需参数: value")
	}

	return nil
}
