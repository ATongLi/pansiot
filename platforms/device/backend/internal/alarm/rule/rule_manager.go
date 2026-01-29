package rule

import (
	"fmt"
	"sync"

	"pansiot-device/internal/core"
)

// RuleManager 报警规则管理器
// 管理规则的完整生命周期，包括增删改查、启用禁用等操作
type RuleManager struct {
	mu      sync.RWMutex
	storage core.Storage
	index   *VariableRuleIndex

	// 规则存储
	rules map[string]*AlarmRule

	// 按类别分组
	categories map[string][]string // categoryID -> ruleIDs

	// 统计信息
	stats RuleManagerStats
}

// RuleManagerStats 统计信息
type RuleManagerStats struct {
	TotalRules     int
	EnabledRules   int
	DisabledRules  int
	RulesByCategory map[string]int
	RulesByLevel   map[core.AlarmLevel]int
}

// NewRuleManager 创建规则管理器
func NewRuleManager(storage core.Storage) *RuleManager {
	return &RuleManager{
		storage:    storage,
		index:      NewVariableRuleIndex(),
		rules:      make(map[string]*AlarmRule),
		categories: make(map[string][]string),
		stats: RuleManagerStats{
			RulesByCategory: make(map[string]int),
			RulesByLevel:   make(map[core.AlarmLevel]int),
		},
	}
}

// AddRule 添加规则
func (rm *RuleManager) AddRule(rule *AlarmRule) error {
	if rule == nil {
		return fmt.Errorf("规则不能为nil")
	}

	if rule.ID == "" {
		return fmt.Errorf("规则ID不能为空")
	}

	// 验证规则
	if err := rm.Validate(rule); err != nil {
		return fmt.Errorf("规则验证失败: %w", err)
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	// 检查规则是否已存在
	if _, exists := rm.rules[rule.ID]; exists {
		return fmt.Errorf("规则已存在: %s", rule.ID)
	}

	// 添加到规则存储
	rm.rules[rule.ID] = rule

	// 更新变量索引
	rm.index.AddRule(rule)

	// 更新类别分组
	rm.categories[rule.Category] = append(rm.categories[rule.Category], rule.ID)

	// 更新统计
	rm.updateStats()

	return nil
}

// RemoveRule 删除规则
func (rm *RuleManager) RemoveRule(ruleID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rule, exists := rm.rules[ruleID]
	if !exists {
		return fmt.Errorf("规则不存在: %s", ruleID)
	}

	// 从变量索引移除
	rm.index.RemoveRule(ruleID)

	// 从类别分组移除
	if ruleIDs, ok := rm.categories[rule.Category]; ok {
		newRuleIDs := make([]string, 0, len(ruleIDs))
		for _, id := range ruleIDs {
			if id != ruleID {
				newRuleIDs = append(newRuleIDs, id)
			}
		}
		rm.categories[rule.Category] = newRuleIDs
	}

	// 从规则存储移除
	delete(rm.rules, ruleID)

	// 更新统计
	rm.updateStats()

	return nil
}

// GetRule 获取规则
func (rm *RuleManager) GetRule(ruleID string) (*AlarmRule, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	rule, exists := rm.rules[ruleID]
	if !exists {
		return nil, false
	}

	// 返回副本，避免外部修改
	return rule.DeepCopy(), true
}

// UpdateRule 更新规则
func (rm *RuleManager) UpdateRule(rule *AlarmRule) error {
	if rule == nil {
		return fmt.Errorf("规则不能为nil")
	}

	if rule.ID == "" {
		return fmt.Errorf("规则ID不能为空")
	}

	// 验证规则
	if err := rm.Validate(rule); err != nil {
		return fmt.Errorf("规则验证失败: %w", err)
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	// 检查规则是否存在
	if _, exists := rm.rules[rule.ID]; !exists {
		return fmt.Errorf("规则不存在: %s", rule.ID)
	}

	// 更新规则
	rm.rules[rule.ID] = rule

	// 更新变量索引
	rm.index.UpdateRule(rule)

	// 更新类别分组（如果类别改变了）
	rm.updateCategoryMapping(rule)

	// 更新统计
	rm.updateStats()

	return nil
}

// EnableRule 启用规则
func (rm *RuleManager) EnableRule(ruleID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rule, exists := rm.rules[ruleID]
	if !exists {
		return fmt.Errorf("规则不存在: %s", ruleID)
	}

	if rule.Enabled {
		return nil // 已经启用
	}

	// 创建副本并修改
	ruleCopy := rule.DeepCopy()
	ruleCopy.Enabled = true

	rm.rules[ruleID] = ruleCopy
	rm.updateStats()

	return nil
}

// DisableRule 禁用规则
func (rm *RuleManager) DisableRule(ruleID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rule, exists := rm.rules[ruleID]
	if !exists {
		return fmt.Errorf("规则不存在: %s", ruleID)
	}

	if !rule.Enabled {
		return nil // 已经禁用
	}

	// 创建副本并修改
	ruleCopy := rule.DeepCopy()
	ruleCopy.Enabled = false

	rm.rules[ruleID] = ruleCopy
	rm.updateStats()

	return nil
}

// ListRules 列出所有规则
func (rm *RuleManager) ListRules() []*AlarmRule {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	rules := make([]*AlarmRule, 0, len(rm.rules))
	for _, rule := range rm.rules {
		rules = append(rules, rule.DeepCopy())
	}
	return rules
}

// ListEnabledRules 列出所有启用的规则
func (rm *RuleManager) ListEnabledRules() []*AlarmRule {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	rules := make([]*AlarmRule, 0)
	for _, rule := range rm.rules {
		if rule.Enabled {
			rules = append(rules, rule.DeepCopy())
		}
	}
	return rules
}

// GetRulesByVariable 获取引用指定变量的规则
func (rm *RuleManager) GetRulesByVariable(variableID uint64) []*AlarmRule {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	ruleIDs := rm.index.GetRulesByVariable(variableID)
	rules := make([]*AlarmRule, 0, len(ruleIDs))

	for _, ruleID := range ruleIDs {
		if rule, exists := rm.rules[ruleID]; exists {
			rules = append(rules, rule.DeepCopy())
		}
	}

	return rules
}

// GetRulesByCategory 获取指定类别的规则
func (rm *RuleManager) GetRulesByCategory(category string) []*AlarmRule {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	ruleIDs, exists := rm.categories[category]
	if !exists {
		return []*AlarmRule{}
	}

	rules := make([]*AlarmRule, 0, len(ruleIDs))
	for _, ruleID := range ruleIDs {
		if rule, exists := rm.rules[ruleID]; exists {
			rules = append(rules, rule.DeepCopy())
		}
	}

	return rules
}

// GetStats 获取统计信息
func (rm *RuleManager) GetStats() RuleManagerStats {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// 返回统计信息的副本
	stats := RuleManagerStats{
		TotalRules:     rm.stats.TotalRules,
		EnabledRules:   rm.stats.EnabledRules,
		DisabledRules:  rm.stats.DisabledRules,
		RulesByCategory: make(map[string]int),
		RulesByLevel:   make(map[core.AlarmLevel]int),
	}

	// 复制 categories 统计
	for k, v := range rm.stats.RulesByCategory {
		stats.RulesByCategory[k] = v
	}

	// 复制 levels 统计
	for k, v := range rm.stats.RulesByLevel {
		stats.RulesByLevel[k] = v
	}

	return stats
}

// Validate 验证规则配置
func (rm *RuleManager) Validate(rule *AlarmRule) error {
	if rule == nil {
		return fmt.Errorf("规则不能为nil")
	}

	// 基础字段验证
	if rule.ID == "" {
		return fmt.Errorf("规则ID不能为空")
	}

	if rule.Name == "" {
		return fmt.Errorf("规则名称不能为空")
	}

	if rule.Category == "" {
		return fmt.Errorf("规则类别不能为空")
	}

	if rule.Level < 1 || rule.Level > 4 {
		return fmt.Errorf("报警等级必须在1-4之间")
	}

	// 条件验证
	if rule.Condition == nil {
		return fmt.Errorf("报警条件不能为空")
	}

	if err := rule.Condition.Validate(); err != nil {
		return fmt.Errorf("条件验证失败: %w", err)
	}

	// 使能条件验证（如果有）
	if rule.EnableCond != nil {
		if err := rule.EnableCond.Validate(); err != nil {
			return fmt.Errorf("使能条件验证失败: %w", err)
		}
	}

	// 报警内容验证
	if err := rm.validateAlarmMessage(&rule.TriggerMsg, true); err != nil {
		return fmt.Errorf("触发内容验证失败: %w", err)
	}

	if rule.RecoverMsg != nil {
		if err := rm.validateAlarmMessage(rule.RecoverMsg, false); err != nil {
			return fmt.Errorf("恢复内容验证失败: %w", err)
		}
	}

	return nil
}

// validateAlarmMessage 验证报警内容
func (rm *RuleManager) validateAlarmMessage(msg *AlarmMessage, isRequired bool) error {
	if msg == nil {
		if isRequired {
			return fmt.Errorf("报警内容不能为空")
		}
		return nil
	}

	// 检查内容类型
	switch msg.Type {
	case ContentStatic:
		if msg.Content == "" {
			return fmt.Errorf("静态内容不能为空")
		}
	case ContentDynamic:
		if msg.Content == "" {
			return fmt.Errorf("动态内容模板不能为空")
		}
		if len(msg.Variables) == 0 {
			return fmt.Errorf("动态内容必须指定变量列表")
		}
		// 验证变量存在性
		for _, varID := range msg.Variables {
			_, err := rm.storage.ReadVar(varID)
			if err != nil {
				return fmt.Errorf("变量不存在: %d", varID)
			}
		}
	case ContentLibrary:
		if msg.Content == "" {
			return fmt.Errorf("文本库内容不能为空")
		}
	}

	return nil
}

// GetVariableIndex 获取变量索引（供外部使用）
func (rm *RuleManager) GetVariableIndex() *VariableRuleIndex {
	return rm.index
}

// updateStats 更新统计信息
func (rm *RuleManager) updateStats() {
	// 重置统计
	rm.stats = RuleManagerStats{
		RulesByCategory: make(map[string]int),
		RulesByLevel:   make(map[core.AlarmLevel]int),
	}

	// 遍历所有规则
	for _, rule := range rm.rules {
		rm.stats.TotalRules++

		// 统计启用/禁用
		if rule.Enabled {
			rm.stats.EnabledRules++
		} else {
			rm.stats.DisabledRules++
		}

		// 按类别统计
		rm.stats.RulesByCategory[rule.Category]++

		// 按等级统计
		rm.stats.RulesByLevel[rule.Level]++
	}
}

// updateCategoryMapping 更新类别映射
func (rm *RuleManager) updateCategoryMapping(rule *AlarmRule) {
	// 查找旧规则
	oldRule, exists := rm.rules[rule.ID]
	if !exists {
		return
	}

	// 如果类别没变，不需要更新
	if oldRule.Category == rule.Category {
		return
	}

	// 从旧类别移除
	oldCategory := oldRule.Category
	if ruleIDs, ok := rm.categories[oldCategory]; ok {
		newRuleIDs := make([]string, 0, len(ruleIDs))
		for _, id := range ruleIDs {
			if id != rule.ID {
				newRuleIDs = append(newRuleIDs, id)
			}
		}
		rm.categories[oldCategory] = newRuleIDs
	}

	// 添加到新类别
	newCategory := rule.Category
	rm.categories[newCategory] = append(rm.categories[newCategory], rule.ID)
}

// GetCategories 获取所有类别
func (rm *RuleManager) GetCategories() []string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	categories := make([]string, 0, len(rm.categories))
	for category := range rm.categories {
		categories = append(categories, category)
	}
	return categories
}

// Clear 清空所有规则
func (rm *RuleManager) Clear() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.rules = make(map[string]*AlarmRule)
	rm.categories = make(map[string][]string)
	rm.index.Clear()
	rm.updateStats()
}

// GetIndex 获取索引（供内部使用）
func (rm *RuleManager) GetIndex() *VariableRuleIndex {
	return rm.index
}
