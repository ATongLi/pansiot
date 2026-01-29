package rule

import (
	"fmt"
	"time"

	"pansiot-device/internal/core"
)

// Condition 报警条件接口
// 支持嵌套的AND/OR条件组合
type Condition interface {
	// Evaluate 评估条件是否满足
	// 注意：完整评估需要状态支持（用于边沿检测和死区处理）
	// 这里提供基础实现，完整版在engine/evaluator.go中
	Evaluate(storage core.Storage) (bool, error)

	// GetVariables 获取条件涉及的所有变量ID
	GetVariables() []uint64

	// Validate 验证条件配置
	Validate() error
}

// CompareOp 比较操作符（统一所有操作）
type CompareOp int

const (
	OpGT  CompareOp = iota // 大于 >
	OpLT                  // 小于 <
	OpGTE                 // 大于等于 >=
	OpLTE                 // 小于等于 <=
	OpEQ                  // 等于 =
	OpNEQ                 // 不等于 !=
	OpRise                // 上升沿 0→1
	OpFall                // 下降沿 1→0
)

// String 返回操作符的字符串表示
func (op CompareOp) String() string {
	switch op {
	case OpGT:
		return ">"
	case OpLT:
		return "<"
	case OpGTE:
		return ">="
	case OpLTE:
		return "<="
	case OpEQ:
		return "="
	case OpNEQ:
		return "!="
	case OpRise:
		return "0→1"
	case OpFall:
		return "1→0"
	default:
		return "unknown"
	}
}

// ParseCompareOp 从字符串解析操作符
func ParseCompareOp(s string) (CompareOp, error) {
	switch s {
	case ">":
		return OpGT, nil
	case "<":
		return OpLT, nil
	case ">=", "=>":
		return OpGTE, nil
	case "<=", "=<":
		return OpLTE, nil
	case "=", "==":
		return OpEQ, nil
	case "!=", "<>":
		return OpNEQ, nil
	case "0→1", "0->1", "rising_edge":
		return OpRise, nil
	case "1→0", "1->0", "falling_edge":
		return OpFall, nil
	default:
		return OpGT, fmt.Errorf("未知的比较操作符: %s", s)
	}
}

// LogicOp 逻辑操作符
type LogicOp int

const (
	LogicAND LogicOp = iota // AND：全部满足
	LogicOR                // OR：任一满足
)

// String 返回逻辑操作符的字符串表示
func (lo LogicOp) String() string {
	switch lo {
	case LogicAND:
		return "AND"
	case LogicOR:
		return "OR"
	default:
		return "unknown"
	}
}

// ParseLogicOp 从字符串解析逻辑操作符
func ParseLogicOp(s string) (LogicOp, error) {
	switch s {
	case "AND", "and", "&", "&&":
		return LogicAND, nil
	case "OR", "or", "|", "||":
		return LogicOR, nil
	default:
		return LogicAND, fmt.Errorf("未知的逻辑操作符: %s", s)
	}
}

// SingleCondition 单个条件（叶子节点）
// 表示一个变量的简单条件判断
type SingleCondition struct {
	VariableID uint64    // 监控的变量ID
	Operator   CompareOp // 比较操作符
	Value      interface{} // 阈值（静态值）

	// 可选字段
	ValueVarID *uint64       // 从变量读取阈值（动态阈值）
	Deadband   float64       // 死区（防止抖动）
	Delay      time.Duration // 延迟时间（延迟触发）
}

// Evaluate 评估单个条件
// 注意：实际评估需要状态支持（用于边沿检测和死区处理）
// 这里提供基础实现，完整的评估在engine/evaluator.go中
func (c *SingleCondition) Evaluate(storage core.Storage) (bool, error) {
	// 读取变量值
	variable, err := storage.ReadVar(c.VariableID)
	if err != nil {
		return false, fmt.Errorf("读取变量[%d]失败: %w", c.VariableID, err)
	}

	// 检查数据质量
	if variable.Quality != core.QualityGood {
		return false, fmt.Errorf("变量[%d]数据质量不佳: %v", c.VariableID, variable.Quality)
	}

	// 获取阈值
	threshold := c.Value
	if c.ValueVarID != nil {
		// 从变量读取阈值
		refVar, err := storage.ReadVar(*c.ValueVarID)
		if err != nil {
			return false, fmt.Errorf("读取阈值变量[%d]失败: %w", *c.ValueVarID, err)
		}
		threshold = refVar.Value
	}

	// 执行比较（简化实现，完整版在evaluator.go中）
	result, err := c.compareValues(variable.Value, threshold)
	if err != nil {
		return false, fmt.Errorf("比较值失败: %w", err)
	}

	return result, nil
}

// compareValues 比较两个值
func (c *SingleCondition) compareValues(a, b interface{}) (bool, error) {
	// 位报警特殊处理
	if c.Operator >= OpRise && c.Operator <= OpFall {
		return c.compareEdge(a, b)
	}

	// 字报警比较
	return c.compareNumeric(a, b)
}

// compareNumeric 数值比较
func (c *SingleCondition) compareNumeric(a, b interface{}) (bool, error) {
	// 转换为float64进行比较
	aFloat, err := toFloat64(a)
	if err != nil {
		return false, err
	}

	bFloat, err := toFloat64(b)
	if err != nil {
		return false, err
	}

	// 应用死区
	if c.Deadband > 0 {
		// 死区逻辑：当前值在阈值±死区范围内时不触发
		// 具体实现需要状态支持
	}

	// 执行比较
	switch c.Operator {
	case OpGT:
		return aFloat > bFloat, nil
	case OpLT:
		return aFloat < bFloat, nil
	case OpGTE:
		return aFloat >= bFloat, nil
	case OpLTE:
		return aFloat <= bFloat, nil
	case OpEQ:
		return aFloat == bFloat, nil
	case OpNEQ:
		return aFloat != bFloat, nil
	default:
		return false, fmt.Errorf("不支持的操作符: %d", c.Operator)
	}
}

// compareEdge 边沿比较（位报警）
func (c *SingleCondition) compareEdge(a, b interface{}) (bool, error) {
	// 转换为bool
	aBool, err := toBool(a)
	if err != nil {
		return false, err
	}

	bBool, err := toBool(b)
	if err != nil {
		return false, err
	}

	// 边沿检测需要历史值支持
	// 这里只提供基础框架
	switch c.Operator {
	case OpRise:
		// 0→1: 需要前一个值为false，当前值为true
		return aBool && !bBool, nil
	case OpFall:
		// 1→0: 需要前一个值为true，当前值为false
		return !aBool && bBool, nil
	default:
		return false, fmt.Errorf("不支持的边沿操作符: %d", c.Operator)
	}
}

// GetVariables 获取条件涉及的所有变量ID
func (c *SingleCondition) GetVariables() []uint64 {
	ids := []uint64{c.VariableID}
	if c.ValueVarID != nil {
		ids = append(ids, *c.ValueVarID)
	}
	return ids
}

// Validate 验证条件配置
func (c *SingleCondition) Validate() error {
	if c.VariableID == 0 {
		return fmt.Errorf("变量ID不能为空")
	}

	// 验证操作符
	if c.Operator < OpGT || c.Operator > OpFall {
		return fmt.Errorf("无效的操作符: %d", c.Operator)
	}

	// 验证阈值：必须设置静态阈值或阈值变量ID
	if c.Value == nil && c.ValueVarID == nil {
		return fmt.Errorf("必须设置阈值或阈值变量ID")
	}

	// 验证延迟时间
	if c.Delay < 0 {
		return fmt.Errorf("延迟时间不能为负数")
	}

	// 验证死区
	if c.Deadband < 0 {
		return fmt.Errorf("死区不能为负数")
	}

	return nil
}

// ConditionGroup 条件组（支持嵌套）
// 表示多个条件的AND/OR组合
type ConditionGroup struct {
	Logic      LogicOp    // 逻辑关系：AND / OR
	Conditions []Condition // 子条件列表
}

// Evaluate 评估条件组
func (g *ConditionGroup) Evaluate(storage core.Storage) (bool, error) {
	if len(g.Conditions) == 0 {
		return false, fmt.Errorf("条件组不能为空")
	}

	switch g.Logic {
	case LogicAND:
		// AND逻辑：所有条件都为true才返回true
		for _, cond := range g.Conditions {
			result, err := cond.Evaluate(storage)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil // AND短路
			}
		}
		return true, nil

	case LogicOR:
		// OR逻辑：任意条件为true就返回true
		for _, cond := range g.Conditions {
			result, err := cond.Evaluate(storage)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil // OR短路
			}
		}
		return false, nil

	default:
		return false, fmt.Errorf("未知的逻辑操作符: %d", g.Logic)
	}
}

// GetVariables 获取所有变量ID（去重）
func (g *ConditionGroup) GetVariables() []uint64 {
	// 使用map去重
	idSet := make(map[uint64]bool)
	for _, cond := range g.Conditions {
		for _, vid := range cond.GetVariables() {
			idSet[vid] = true
		}
	}

	// 转换为切片
	ids := make([]uint64, 0, len(idSet))
	for vid := range idSet {
		ids = append(ids, vid)
	}
	return ids
}

// Validate 验证条件组配置
func (g *ConditionGroup) Validate() error {
	if len(g.Conditions) == 0 {
		return fmt.Errorf("条件组不能为空")
	}

	// 验证逻辑操作符
	if g.Logic != LogicAND && g.Logic != LogicOR {
		return fmt.Errorf("无效的逻辑操作符: %d", g.Logic)
	}

	// 验证每个子条件
	for i, cond := range g.Conditions {
		if err := cond.Validate(); err != nil {
			return fmt.Errorf("条件[%d]验证失败: %w", i, err)
		}
	}

	return nil
}

// 辅助函数：转换为float64
func toFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	case float64:
		return val, nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("不支持的类型: %T", v)
	}
}

// 辅助函数：转换为bool
func toBool(v interface{}) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case int, int8, int16, int32, int64:
		return val != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		return val != 0, nil
	case float32:
		return val != 0, nil
	case float64:
		return val != 0, nil
	default:
		return false, fmt.Errorf("不支持的类型: %T", v)
	}
}

// NewSimpleCondition 创建简单条件
func NewSimpleCondition(variableID uint64, operator CompareOp, value interface{}) *SingleCondition {
	return &SingleCondition{
		VariableID: variableID,
		Operator:   operator,
		Value:      value,
	}
}

// NewAndCondition 创建AND条件组
func NewAndCondition(conditions ...Condition) *ConditionGroup {
	return &ConditionGroup{
		Logic:      LogicAND,
		Conditions: conditions,
	}
}

// NewOrCondition 创建OR条件组
func NewOrCondition(conditions ...Condition) *ConditionGroup {
	return &ConditionGroup{
		Logic:      LogicOR,
		Conditions: conditions,
	}
}
