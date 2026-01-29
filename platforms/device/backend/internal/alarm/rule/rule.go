package rule

import (
	"fmt"
	"time"

	"pansiot-device/internal/core"
)

// RuleType 报警规则类型
type RuleType int

const (
	RuleTypeUserDefined RuleType = 1 // 用户自定义报警
	RuleTypeSystem      RuleType = 2 // 系统报警
	RuleTypePredictive  RuleType = 3 // 预测性报警
)

// String 返回规则类型的字符串表示
func (t RuleType) String() string {
	switch t {
	case RuleTypeUserDefined:
		return "user"
	case RuleTypeSystem:
		return "system"
	case RuleTypePredictive:
		return "predictive"
	default:
		return "unknown"
	}
}

// AlarmRule 报警规则
type AlarmRule struct {
	// ========== 基础信息 ==========
	ID          string     // 规则ID，全局唯一
	Name        string     // 规则名称
	Type        RuleType   // 规则类型：User/System/Predictive
	Category    string     // 报警类别ID（报警组）
	Level       core.AlarmLevel // 报警等级：1-低 2-中 3-高 4-严重
	Enabled     bool       // 是否启用

	// ========== 报警条件 ==========
	Condition   Condition   // 报警条件（支持嵌套AND/OR）
	EnableCond  Condition   // 使能条件（可选）：nil表示无使能条件

	// ========== 报警内容 ==========
	TriggerMsg   AlarmMessage  // 触发时的报警内容
	RecoverMsg   *AlarmMessage  // 恢复时的报警内容（可选）

	// ========== 固化动作配置（性能优先） ==========
	// 这些动作直接执行，无需遍历和类型判断
	EnableRecord     bool   // 是否记录报警事件到数据库（默认true）
	EnableCloudPush  bool   // 是否推送到云端（默认false）

	// ========== 自定义动作（可选，需遍历执行） ==========
	TriggerActions []Action    // 触发时执行的自定义动作列表
	RecoverActions []Action    // 恢复时执行的自定义动作列表（可选）

	// ========== 报警通知 ==========
	Sound       *SoundConfig   // 报警音配置（可选）
	Responsible []string       // 责任人用户ID列表

	// ========== 元数据 ==========
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   string
}

// AlarmMessage 报警内容
type AlarmMessage struct {
	Type      ContentType // 内容类型
	Content   string      // 内容文本（静态文本或动态模板）
	Variables []uint64    // 涉及的变量ID列表（用于动态内容）
}

// ContentType 内容类型
type ContentType int

const (
	ContentStatic   ContentType = iota // 静态内容：手动输入的固定文本
	ContentDynamic                     // 动态内容：可插入变量值，格式 {var:ID}
	ContentLibrary                     // 文本库：从文本库读取，支持多语言
)

// String 返回内容类型的字符串表示
func (ct ContentType) String() string {
	switch ct {
	case ContentStatic:
		return "static"
	case ContentDynamic:
		return "dynamic"
	case ContentLibrary:
		return "library"
	default:
		return "unknown"
	}
}

// SoundConfig 报警音配置
type SoundConfig struct {
	Enabled    bool          // 是否启用报警音
	File       string        // 音频文件路径（用户上传或系统内置）
	Continuous bool          // 是否持续播放
	Interval   time.Duration // 播放间隔（持续播放时有效）
	Count      int           // 播放次数（非持续播放时，0=无限）
}

// Action 报警动作
type Action struct {
	ID     string      // 动作ID（可选，用于状态跟踪）
	Type   ActionType  // 动作类型
	When   TriggerTime // 执行时机
	Mode   ExecuteMode // 执行模式

	// 参数（根据动作类型不同而不同）
	Params map[string]interface{}

	// 循环执行相关（Mode = Loop时使用）
	LoopCount  int           // 循环次数（0=无限循环）
	LoopUntil  core.AlarmState // 循环直到某个状态（0=不使用）
	LoopDelay  time.Duration // 循环间隔

	// 延迟执行相关（When = After时使用）
	Delay time.Duration // 延迟时长
}

// ActionType 动作类型
type ActionType int

const (
	ActSound       ActionType = iota // 播放声音
	ActJumpPage                      // 跳转页面
	ActScript                        // 执行脚本
	ActWriteVar                      // 写变量值
	ActPopup                         // 弹窗提示
)

// String 返回动作类型的字符串表示
func (at ActionType) String() string {
	switch at {
	case ActSound:
		return "sound"
	case ActJumpPage:
		return "jump_page"
	case ActScript:
		return "script"
	case ActWriteVar:
		return "write_var"
	case ActPopup:
		return "popup"
	default:
		return "unknown"
	}
}

// TriggerTime 执行时机
type TriggerTime int

const (
	WhenTrigger       TriggerTime = iota // 触发时立即执行
	WhenAfterTrigger                     // 触发后延迟执行
	WhenRecover                           // 恢复时立即执行
	WhenAfterRecover                     // 恢复后延迟执行
)

// String 返回执行时机的字符串表示
func (tt TriggerTime) String() string {
	switch tt {
	case WhenTrigger:
		return "on_trigger"
	case WhenAfterTrigger:
		return "after_trigger"
	case WhenRecover:
		return "on_recover"
	case WhenAfterRecover:
		return "after_recover"
	default:
		return "unknown"
	}
}

// ExecuteMode 执行模式
type ExecuteMode int

const (
	ModeOnce  ExecuteMode = iota // 执行一次
	ModeLoop                     // 循环执行
)

// String 返回执行模式的字符串表示
func (em ExecuteMode) String() string {
	switch em {
	case ModeOnce:
		return "once"
	case ModeLoop:
		return "loop"
	default:
		return "unknown"
	}
}

// AlarmCategory 报警类别（报警组）
type AlarmCategory struct {
	ID          string        // 类别ID
	Name        string        // 类别名称
	Description string        // 描述
	Shield      *ShieldConfig // 屏蔽配置（可选）
	CloudPush   bool          // 是否启用云推送
}

// ShieldConfig 屏蔽配置
type ShieldConfig struct {
	Enabled    bool              // 是否启用屏蔽
	Levels     []core.AlarmLevel // 屏蔽的报警等级列表
	Types      []RuleType        // 屏蔽的规则类型列表
	KeepRecord bool              // 屏蔽后是否保留记录
}

// Validate 验证报警规则配置
func (r *AlarmRule) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("规则ID不能为空")
	}

	if r.Name == "" {
		return fmt.Errorf("规则名称不能为空")
	}

	if r.Condition == nil {
		return fmt.Errorf("报警条件不能为空")
	}

	if err := r.Condition.Validate(); err != nil {
		return fmt.Errorf("报警条件验证失败: %w", err)
	}

	if r.Level < core.AlarmLevelLow || r.Level > core.AlarmLevelCritical {
		return fmt.Errorf("报警等级无效: %d", r.Level)
	}

	// 验证触发内容
	if err := r.TriggerMsg.Validate(); err != nil {
		return fmt.Errorf("触发内容验证失败: %w", err)
	}

	// 验证恢复内容（可选）
	if r.RecoverMsg != nil {
		if err := r.RecoverMsg.Validate(); err != nil {
			return fmt.Errorf("恢复内容验证失败: %w", err)
		}
	}

	// 验证触发动作
	for i, action := range r.TriggerActions {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("触发动作[%d]验证失败: %w", i, err)
		}
	}

	// 验证恢复动作（可选）
	for i, action := range r.RecoverActions {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("恢复动作[%d]验证失败: %w", i, err)
		}
	}

	return nil
}

// Validate 验证报警内容
func (m *AlarmMessage) Validate() error {
	switch m.Type {
	case ContentStatic:
		if m.Content == "" {
			return fmt.Errorf("静态内容不能为空")
		}
	case ContentDynamic:
		if m.Content == "" {
			return fmt.Errorf("动态模板不能为空")
		}
	case ContentLibrary:
		if m.Content == "" {
			return fmt.Errorf("文本库ID不能为空")
		}
	default:
		return fmt.Errorf("无效的内容类型: %d", m.Type)
	}
	return nil
}

// Validate 验证报警动作
func (a *Action) Validate() error {
	// 验证动作类型
	switch a.Type {
	case ActSound, ActJumpPage, ActScript, ActWriteVar, ActPopup:
		// 有效类型
	default:
		return fmt.Errorf("无效的动作类型: %d", a.Type)
	}

	// 验证执行模式
	if a.Mode == ModeLoop {
		if a.LoopCount == 0 && a.LoopUntil == 0 {
			return fmt.Errorf("循环执行模式必须设置次数或停止条件")
		}
	}

	// 验证延迟执行
	if a.When == WhenAfterTrigger || a.When == WhenAfterRecover {
		if a.Delay <= 0 {
			return fmt.Errorf("延迟执行必须设置延迟时长")
		}
	}

	return nil
}

// GetVariableIDs 获取报警规则涉及的所有变量ID
func (r *AlarmRule) GetVariableIDs() []uint64 {
	varIDSet := make(map[uint64]bool)

	// 添加条件变量
	for _, vid := range r.Condition.GetVariables() {
		varIDSet[vid] = true
	}

	// 添加使能条件变量
	if r.EnableCond != nil {
		for _, vid := range r.EnableCond.GetVariables() {
			varIDSet[vid] = true
		}
	}

	// 添加触发内容变量
	for _, vid := range r.TriggerMsg.Variables {
		varIDSet[vid] = true
	}

	// 添加恢复内容变量
	if r.RecoverMsg != nil {
		for _, vid := range r.RecoverMsg.Variables {
			varIDSet[vid] = true
		}
	}

	// 转换为切片
	vids := make([]uint64, 0, len(varIDSet))
	for vid := range varIDSet {
		vids = append(vids, vid)
	}
	return vids
}

// DeepCopy 创建规则的深拷贝
func (r *AlarmRule) DeepCopy() *AlarmRule {
	if r == nil {
		return nil
	}

	// 创建拷贝并复制基本字段
	result := &AlarmRule{
		ID:              r.ID,
		Name:            r.Name,
		Type:            r.Type,
		Category:        r.Category,
		Level:           r.Level,
		Enabled:         r.Enabled,
		EnableRecord:    r.EnableRecord,
		EnableCloudPush: r.EnableCloudPush,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
		CreatedBy:       r.CreatedBy,
	}

	// 拷贝条件（注意：Condition 接口类型需要特殊处理）
	if r.Condition != nil {
		// 对于已知的条件类型，进行类型断言和深拷贝
		if singleCond, ok := r.Condition.(*SingleCondition); ok {
			condCopy := *singleCond
			result.Condition = &condCopy
		} else if groupCond, ok := r.Condition.(*ConditionGroup); ok {
			// 深拷贝条件组
			condsCopy := make([]Condition, len(groupCond.Conditions))
			for i, cond := range groupCond.Conditions {
				if single, ok := cond.(*SingleCondition); ok {
					singleCopy := *single
					condsCopy[i] = &singleCopy
				} else if group, ok := cond.(*ConditionGroup); ok {
					// 递归拷贝嵌套的条件组
					groupCopy := ConditionGroup{
						Logic:      group.Logic,
						Conditions: []Condition{},
					}
					// 简化处理：嵌套条件组的深拷贝
					for _, c := range group.Conditions {
						if s, ok := c.(*SingleCondition); ok {
							sCopy := *s
							groupCopy.Conditions = append(groupCopy.Conditions, &sCopy)
						}
					}
					condsCopy[i] = &groupCopy
				}
			}
			groupCondCopy := ConditionGroup{
				Logic:      groupCond.Logic,
				Conditions: condsCopy,
			}
			result.Condition = &groupCondCopy
		}
	}

	// 拷贝使能条件
	if r.EnableCond != nil {
		if singleCond, ok := r.EnableCond.(*SingleCondition); ok {
			condCopy := *singleCond
			result.EnableCond = &condCopy
		}
	}

	// 拷贝报警内容
	result.TriggerMsg = r.TriggerMsg
	if r.TriggerMsg.Variables != nil {
		varsCopy := make([]uint64, len(r.TriggerMsg.Variables))
		copy(varsCopy, r.TriggerMsg.Variables)
		result.TriggerMsg.Variables = varsCopy
	}

	if r.RecoverMsg != nil {
		recoverMsgCopy := *r.RecoverMsg
		if r.RecoverMsg.Variables != nil {
			varsCopy := make([]uint64, len(r.RecoverMsg.Variables))
			copy(varsCopy, r.RecoverMsg.Variables)
			recoverMsgCopy.Variables = varsCopy
		}
		result.RecoverMsg = &recoverMsgCopy
	}

	// 拷贝动作列表
	if r.TriggerActions != nil {
		actionsCopy := make([]Action, len(r.TriggerActions))
		for i, a := range r.TriggerActions {
			actionsCopy[i] = a
		}
		result.TriggerActions = actionsCopy
	}

	if r.RecoverActions != nil {
		actionsCopy := make([]Action, len(r.RecoverActions))
		for i, a := range r.RecoverActions {
			actionsCopy[i] = a
		}
		result.RecoverActions = actionsCopy
	}

	// 拷贝责任人列表
	if r.Responsible != nil {
		responsibleCopy := make([]string, len(r.Responsible))
		for i, u := range r.Responsible {
			responsibleCopy[i] = u
		}
		result.Responsible = responsibleCopy
	}

	// 拷贝报警音配置
	if r.Sound != nil {
		soundCopy := *r.Sound
		result.Sound = &soundCopy
	}

	return result
}
