package rule

import (
	"testing"
	"time"

	"pansiot-device/internal/core"
)

// TestAlarmRuleValidation 测试报警规则验证
func TestAlarmRuleValidation(t *testing.T) {
	tests := []struct {
		name    string
		rule    *AlarmRule
		wantErr bool
	}{
		{
			name: "有效规则 - 简单条件",
			rule: &AlarmRule{
				ID:       "RULE_001",
				Name:     "测试规则",
				Type:     RuleTypeUserDefined,
				Category: "CAT_TEST",
				Condition: &SingleCondition{
					VariableID: 100001,
					Operator:   OpGT,
					Value:      80.0,
				},
				Level: core.AlarmLevelHigh,
				TriggerMsg: AlarmMessage{
					Type:    ContentStatic,
					Content: "测试报警",
				},
				Enabled:   true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "无效规则 - 空ID",
			rule: &AlarmRule{
				ID:   "",
				Name: "测试规则",
				Condition: &SingleCondition{
					VariableID: 100001,
					Operator:   OpGT,
					Value:      80.0,
				},
				Level: core.AlarmLevelHigh,
			},
			wantErr: true,
		},
		{
			name: "无效规则 - 空名称",
			rule: &AlarmRule{
				ID:   "RULE_001",
				Name: "",
				Condition: &SingleCondition{
					VariableID: 100001,
					Operator:   OpGT,
					Value:      80.0,
				},
				Level: core.AlarmLevelHigh,
			},
			wantErr: true,
		},
		{
			name: "无效规则 - 空条件",
			rule: &AlarmRule{
				ID:        "RULE_001",
				Name:      "测试规则",
				Condition: nil,
				Level:     core.AlarmLevelHigh,
			},
			wantErr: true,
		},
		{
			name: "无效规则 - 无效级别",
			rule: &AlarmRule{
				ID:   "RULE_001",
				Name: "测试规则",
				Condition: &SingleCondition{
					VariableID: 100001,
					Operator:   OpGT,
					Value:      80.0,
				},
				Level: 5, // 无效级别
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AlarmRule.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConditionValidation 测试条件验证
func TestConditionValidation(t *testing.T) {
	tests := []struct {
		name      string
		condition Condition
		wantErr   bool
	}{
		{
			name: "有效条件 - 简单条件",
			condition: &SingleCondition{
				VariableID: 100001,
				Operator:   OpGT,
				Value:      80.0,
			},
			wantErr: false,
		},
		{
			name: "有效条件 - 带变量引用",
			condition: &SingleCondition{
				VariableID: 100001,
				Operator:   OpGT,
				ValueVarID: uint64Ptr(100002),
			},
			wantErr: false,
		},
		{
			name: "无效条件 - 空变量ID",
			condition: &SingleCondition{
				VariableID: 0,
				Operator:   OpGT,
				Value:      80.0,
			},
			wantErr: true,
		},
		{
			name: "无效条件 - 无阈值且无变量引用",
			condition: &SingleCondition{
				VariableID: 100001,
				Operator:   OpGT,
				Value:      nil,
			},
			wantErr: true,
		},
		{
			name: "有效条件 - AND组",
			condition: &ConditionGroup{
				Logic: LogicAND,
				Conditions: []Condition{
					&SingleCondition{
						VariableID: 100001,
						Operator:   OpGT,
						Value:      80.0,
					},
					&SingleCondition{
						VariableID: 100002,
						Operator:   OpLT,
						Value:      100.0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "有效条件 - 嵌套AND/OR",
			condition: &ConditionGroup{
				Logic: LogicAND,
				Conditions: []Condition{
					&SingleCondition{
						VariableID: 100001,
						Operator:   OpGT,
						Value:      80.0,
					},
					&ConditionGroup{
						Logic: LogicOR,
						Conditions: []Condition{
							&SingleCondition{
								VariableID: 100002,
								Operator:   OpGT,
								Value:      90.0,
							},
							&SingleCondition{
								VariableID: 100003,
								Operator:   OpLT,
								Value:      50.0,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "无效条件 - 空条件组",
			condition: &ConditionGroup{
				Logic:      LogicAND,
				Conditions: []Condition{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.condition.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AlarmCondition.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAlarmMessageValidation 测试报警内容验证
func TestAlarmMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		message *AlarmMessage
		wantErr bool
	}{
		{
			name: "有效静态内容",
			message: &AlarmMessage{
				Type:    ContentStatic,
				Content: "测试报警消息",
			},
			wantErr: false,
		},
		{
			name: "有效动态内容",
			message: &AlarmMessage{
				Type:      ContentDynamic,
				Content:   "设备温度过高: {var:100001}°C",
				Variables: []uint64{100001},
			},
			wantErr: false,
		},
		{
			name: "有效文本库内容",
			message: &AlarmMessage{
				Type:    ContentLibrary,
				Content: "TEXT_001",
			},
			wantErr: false,
		},
		{
			name: "无效静态内容 - 空文本",
			message: &AlarmMessage{
				Type:    ContentStatic,
				Content: "",
			},
			wantErr: true,
		},
		{
			name: "无效动态内容 - 空模板",
			message: &AlarmMessage{
				Type:    ContentDynamic,
				Content: "",
			},
			wantErr: true,
		},
		{
			name: "无效文本库内容 - 空ID",
			message: &AlarmMessage{
				Type:    ContentLibrary,
				Content: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AlarmMessage.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestActionValidation 测试动作配置验证
func TestActionValidation(t *testing.T) {
	tests := []struct {
		name    string
		action  *Action
		wantErr bool
	}{
		{
			name: "有效动作 - 执行一次",
			action: &Action{
				ID:     "ACTION_001",
				Type:   ActSound,
				When:   WhenTrigger,
				Mode:   ModeOnce,
				Params: make(map[string]interface{}),
			},
			wantErr: false,
		},
		{
			name: "有效动作 - 循环执行（有次数限制）",
			action: &Action{
				ID:        "ACTION_002",
				Type:      ActWriteVar,
				When:      WhenAfterTrigger,
				Mode:      ModeLoop,
				LoopCount: 5,
				LoopDelay: time.Second * 10,
				Delay:     time.Second * 5,
				Params:    make(map[string]interface{}),
			},
			wantErr: false,
		},
		{
			name: "有效动作 - 循环执行（有状态控制）",
			action: &Action{
				ID:        "ACTION_003",
				Type:      ActPopup,
				When:      WhenRecover,
				Mode:      ModeLoop,
				LoopUntil: core.AlarmStateCleared,
				LoopDelay: time.Second * 5,
				Params:    make(map[string]interface{}),
			},
			wantErr: false,
		},
		{
			name: "无效动作 - 无效类型",
			action: &Action{
				ID:     "ACTION_001",
				Type:   99, // 无效类型
				When:   WhenTrigger,
				Mode:   ModeOnce,
				Params: make(map[string]interface{}),
			},
			wantErr: true,
		},
		{
			name: "无效动作 - 无效循环配置",
			action: &Action{
				ID:     "ACTION_004",
				Type:   ActSound,
				When:   WhenTrigger,
				Mode:   ModeLoop,
				Params: make(map[string]interface{}),
			},
			wantErr: true,
		},
		{
			name: "无效动作 - 延迟执行未设置延迟时间",
			action: &Action{
				ID:     "ACTION_005",
				Type:   ActSound,
				When:   WhenAfterTrigger,
				Mode:   ModeOnce,
				Params: make(map[string]interface{}),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.action.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Action.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCompareOp 测试比较操作符
func TestCompareOp(t *testing.T) {
	operators := []CompareOp{
		OpGT,
		OpLT,
		OpGTE,
		OpLTE,
		OpEQ,
		OpNEQ,
		OpRise,
		OpFall,
	}

	for _, op := range operators {
		t.Run(op.String(), func(t *testing.T) {
			if op.String() == "unknown" {
				t.Errorf("操作符 %d 返回了 unknown", op)
			}
		})
	}
}

// TestLogicOp 测试逻辑操作符
func TestLogicOp(t *testing.T) {
	operators := []LogicOp{LogicAND, LogicOR}

	for _, op := range operators {
		t.Run(op.String(), func(t *testing.T) {
			if op.String() == "unknown" {
				t.Errorf("逻辑操作符 %d 返回了 unknown", op)
			}
		})
	}
}

// TestConditionGetVariableIDs 测试获取变量ID
func TestConditionGetVariableIDs(t *testing.T) {
	// 单个条件
	cond1 := &SingleCondition{
		VariableID: 100001,
		Operator:   OpGT,
		Value:      80.0,
		ValueVarID: uint64Ptr(100002),
	}

	ids := cond1.GetVariables()
	if len(ids) != 2 {
		t.Errorf("期望2个变量ID，实际得到 %d", len(ids))
	}

	// 条件组
	cond2 := &ConditionGroup{
		Logic: LogicAND,
		Conditions: []Condition{
			&SingleCondition{
				VariableID: 100001,
				Operator:   OpGT,
				Value:      80.0,
			},
			&SingleCondition{
				VariableID: 100002,
				Operator:   OpLT,
				Value:      100.0,
			},
			&SingleCondition{
				VariableID: 100001, // 重复的变量ID
				Operator:   OpGT,
				Value:      50.0,
			},
		},
	}

	ids = cond2.GetVariables()
	// 应该去重，只有2个唯一的变量ID
	if len(ids) != 2 {
		t.Errorf("期望2个唯一变量ID（去重后），实际得到 %d", len(ids))
	}
}

// TestRuleType 测试报警规则类型
func TestRuleType(t *testing.T) {
	types := []RuleType{
		RuleTypeUserDefined,
		RuleTypeSystem,
		RuleTypePredictive,
	}

	expected := []string{"user", "system", "predictive"}

	for i, typ := range types {
		t.Run(typ.String(), func(t *testing.T) {
			if typ.String() != expected[i] {
				t.Errorf("期望 %s，实际得到 %s", expected[i], typ.String())
			}
		})
	}
}

// TestContentType 测试内容类型
func TestContentType(t *testing.T) {
	types := []ContentType{ContentStatic, ContentDynamic, ContentLibrary}

	expected := []string{"static", "dynamic", "library"}

	for i, typ := range types {
		t.Run(typ.String(), func(t *testing.T) {
			if typ.String() != expected[i] {
				t.Errorf("期望 %s，实际得到 %s", expected[i], typ.String())
			}
		})
	}
}

// TestActionType 测试动作类型
func TestActionType(t *testing.T) {
	types := []ActionType{
		ActSound,
		ActJumpPage,
		ActScript,
		ActWriteVar,
		ActPopup,
	}

	expected := []string{"sound", "jump_page", "script", "write_var", "popup"}

	for i, typ := range types {
		t.Run(typ.String(), func(t *testing.T) {
			if typ.String() != expected[i] {
				t.Errorf("期望 %s，实际得到 %s", expected[i], typ.String())
			}
		})
	}
}

// TestTriggerTime 测试触发时机
func TestTriggerTime(t *testing.T) {
	times := []TriggerTime{
		WhenTrigger,
		WhenAfterTrigger,
		WhenRecover,
		WhenAfterRecover,
	}

	expected := []string{"on_trigger", "after_trigger", "on_recover", "after_recover"}

	for i, tt := range times {
		t.Run(tt.String(), func(t *testing.T) {
			if tt.String() != expected[i] {
				t.Errorf("期望 %s，实际得到 %s", expected[i], tt.String())
			}
		})
	}
}

// TestExecuteMode 测试执行模式
func TestExecuteMode(t *testing.T) {
	modes := []ExecuteMode{ModeOnce, ModeLoop}

	expected := []string{"once", "loop"}

	for i, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			if mode.String() != expected[i] {
				t.Errorf("期望 %s，实际得到 %s", expected[i], mode.String())
			}
		})
	}
}

// TestAlarmRuleGetVariableIDs 测试报警规则获取变量ID
func TestAlarmRuleGetVariableIDs(t *testing.T) {
	rule := &AlarmRule{
		ID:   "RULE_001",
		Name: "测试规则",
		Condition: &SingleCondition{
			VariableID: 100001,
			Operator:   OpGT,
			Value:      80.0,
		},
		EnableCond: &SingleCondition{
			VariableID: 100002,
			Operator:   OpEQ,
			Value:      0,
		},
		TriggerMsg: AlarmMessage{
			Type:      ContentDynamic,
			Content:   "设备{var:100003}温度过高: {var:100001}°C",
			Variables: []uint64{100001, 100003},
		},
		RecoverMsg: &AlarmMessage{
			Type:      ContentDynamic,
			Content:   "温度已恢复正常: {var:100001}°C",
			Variables: []uint64{100001},
		},
	}

	ids := rule.GetVariableIDs()
	// 应该包含：100001（主条件）, 100002（使能条件）, 100001, 100003（触发消息）, 100001（恢复消息）
	// 去重后：100001, 100002, 100003
	if len(ids) != 3 {
		t.Errorf("期望3个唯一变量ID，实际得到 %d", len(ids))
	}

	// 验证去重
	idMap := make(map[uint64]bool)
	for _, id := range ids {
		idMap[id] = true
	}
	if len(idMap) != 3 {
		t.Errorf("变量ID未正确去重")
	}
}

// TestNewConditionHelpers 测试条件构造辅助函数
func TestNewConditionHelpers(t *testing.T) {
	// 测试简单条件
	simple := NewSimpleCondition(100001, OpGT, 80.0)
	if simple.VariableID != 100001 {
		t.Errorf("NewSimpleCondition 失败")
	}
	if simple.Operator != OpGT {
		t.Errorf("NewSimpleCondition 操作符错误")
	}
	if simple.Value != 80.0 {
		t.Errorf("NewSimpleCondition 值错误")
	}

	// 测试AND条件
	and := NewAndCondition(
		NewSimpleCondition(100001, OpGT, 80.0),
		NewSimpleCondition(100002, OpLT, 100.0),
	)
	if and.Logic != LogicAND {
		t.Errorf("NewAndCondition 失败")
	}
	if len(and.Conditions) != 2 {
		t.Errorf("NewAndCondition 条件数量错误")
	}

	// 测试OR条件
	or := NewOrCondition(
		NewSimpleCondition(100001, OpGT, 80.0),
		NewSimpleCondition(100002, OpLT, 100.0),
	)
	if or.Logic != LogicOR {
		t.Errorf("NewOrCondition 失败")
	}
	if len(or.Conditions) != 2 {
		t.Errorf("NewOrCondition 条件数量错误")
	}
}

// TestParseCompareOp 测试操作符解析
func TestParseCompareOp(t *testing.T) {
	tests := []struct {
		input    string
		expected CompareOp
		wantErr  bool
	}{
		{">", OpGT, false},
		{"<", OpLT, false},
		{">=", OpGTE, false},
		{"<=", OpLTE, false},
		{"=", OpEQ, false},
		{"!=", OpNEQ, false},
		{"0→1", OpRise, false},
		{"1→0", OpFall, false},
		{"invalid", OpGT, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			op, err := ParseCompareOp(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCompareOp(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && op != tt.expected {
				t.Errorf("ParseCompareOp(%s) = %v, want %v", tt.input, op, tt.expected)
			}
		})
	}
}

// TestParseLogicOp 测试逻辑操作符解析
func TestParseLogicOp(t *testing.T) {
	tests := []struct {
		input    string
		expected LogicOp
		wantErr  bool
	}{
		{"AND", LogicAND, false},
		{"and", LogicAND, false},
		{"&", LogicAND, false},
		{"&&", LogicAND, false},
		{"OR", LogicOR, false},
		{"or", LogicOR, false},
		{"|", LogicOR, false},
		{"||", LogicOR, false},
		{"invalid", LogicAND, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			op, err := ParseLogicOp(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLogicOp(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && op != tt.expected {
				t.Errorf("ParseLogicOp(%s) = %v, want %v", tt.input, op, tt.expected)
			}
		})
	}
}

// 辅助函数
func uint64Ptr(v uint64) *uint64 {
	return &v
}
