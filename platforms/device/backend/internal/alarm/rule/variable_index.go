package rule

import (
	"sync"
)

// VariableRuleIndex 变量-规则索引
// 维护变量ID与规则ID的双向映射关系，支持快速查找
type VariableRuleIndex struct {
	mu sync.RWMutex

	// variableID -> ruleIDs 映射
	variableToRules map[uint64]map[string]bool

	// ruleID -> variableIDs 映射
	ruleToVariables map[string][]uint64

	// 按设备分组（可选优化）
	// deviceID -> variableIDs 映射
	deviceToVariables map[string][]uint64
}

// RuleReference 规则引用信息
type RuleReference struct {
	RuleID      string
	VariableIDs []uint64
	DeviceIDs   []string // 从变量ID提取的设备ID
}

// NewVariableRuleIndex 创建变量-规则索引
func NewVariableRuleIndex() *VariableRuleIndex {
	return &VariableRuleIndex{
		variableToRules:  make(map[uint64]map[string]bool),
		ruleToVariables:  make(map[string][]uint64),
		deviceToVariables: make(map[string][]uint64),
	}
}

// AddRule 添加规则到索引
func (idx *VariableRuleIndex) AddRule(rule *AlarmRule) {
	if rule == nil {
		return
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	// 提取规则涉及的所有变量ID
	variableIDs := rule.GetVariableIDs()

	// 更新 ruleID -> variableIDs 映射
	idx.ruleToVariables[rule.ID] = variableIDs

	// 更新 variableID -> ruleIDs 映射
	for _, variableID := range variableIDs {
		if _, exists := idx.variableToRules[variableID]; !exists {
			idx.variableToRules[variableID] = make(map[string]bool)
		}
		idx.variableToRules[variableID][rule.ID] = true
	}

	// 可选：按设备分组（后续优化）
	// idx.updateDeviceIndex(rule, variableIDs)
}

// RemoveRule 从索引移除规则
func (idx *VariableRuleIndex) RemoveRule(ruleID string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// 获取规则涉及的所有变量ID
	variableIDs, exists := idx.ruleToVariables[ruleID]
	if !exists {
		return
	}

	// 从 variableID -> ruleIDs 映射中移除
	for _, variableID := range variableIDs {
		if ruleMap, ok := idx.variableToRules[variableID]; ok {
			delete(ruleMap, ruleID)
			// 如果该变量没有其他规则引用了，删除整个条目
			if len(ruleMap) == 0 {
				delete(idx.variableToRules, variableID)
			}
		}
	}

	// 从 ruleID -> variableIDs 映射中移除
	delete(idx.ruleToVariables, ruleID)
}

// UpdateRule 更新规则索引（先删除再添加）
func (idx *VariableRuleIndex) UpdateRule(rule *AlarmRule) {
	if rule == nil {
		return
	}

	// 先删除旧的索引
	idx.RemoveRule(rule.ID)

	// 再添加新的索引
	idx.AddRule(rule)
}

// GetRulesByVariable 获取引用指定变量的所有规则ID
func (idx *VariableRuleIndex) GetRulesByVariable(variableID uint64) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	ruleMap, exists := idx.variableToRules[variableID]
	if !exists {
		return []string{}
	}

	// 转换为切片
	ruleIDs := make([]string, 0, len(ruleMap))
	for ruleID := range ruleMap {
		ruleIDs = append(ruleIDs, ruleID)
	}
	return ruleIDs
}

// GetVariablesByRule 获取规则涉及的所有变量ID
func (idx *VariableRuleIndex) GetVariablesByRule(ruleID string) []uint64 {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	variableIDs, exists := idx.ruleToVariables[ruleID]
	if !exists {
		return []uint64{}
	}

	// 返回副本
	result := make([]uint64, len(variableIDs))
	copy(result, variableIDs)
	return result
}

// GetAllMappings 获取所有变量-规则映射
func (idx *VariableRuleIndex) GetAllMappings() map[uint64][]string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	mappings := make(map[uint64][]string)

	for variableID, ruleMap := range idx.variableToRules {
		ruleIDs := make([]string, 0, len(ruleMap))
		for ruleID := range ruleMap {
			ruleIDs = append(ruleIDs, ruleID)
		}
		mappings[variableID] = ruleIDs
	}

	return mappings
}

// GetAllRules 获取所有已索引的规则ID
func (idx *VariableRuleIndex) GetAllRules() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	ruleIDs := make([]string, 0, len(idx.ruleToVariables))
	for ruleID := range idx.ruleToVariables {
		ruleIDs = append(ruleIDs, ruleID)
	}
	return ruleIDs
}

// GetAllVariables 获取所有被引用的变量ID
func (idx *VariableRuleIndex) GetAllVariables() []uint64 {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	variableIDs := make([]uint64, 0, len(idx.variableToRules))
	for variableID := range idx.variableToRules {
		variableIDs = append(variableIDs, variableID)
	}
	return variableIDs
}

// HasRule 检查规则是否已被索引
func (idx *VariableRuleIndex) HasRule(ruleID string) bool {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	_, exists := idx.ruleToVariables[ruleID]
	return exists
}

// HasVariable 检查变量是否被任何规则引用
func (idx *VariableRuleIndex) HasVariable(variableID uint64) bool {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	_, exists := idx.variableToRules[variableID]
	return exists
}

// GetRuleCount 获取引用指定变量的规则数量
func (idx *VariableRuleIndex) GetRuleCount(variableID uint64) int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	ruleMap, exists := idx.variableToRules[variableID]
	if !exists {
		return 0
	}
	return len(ruleMap)
}

// GetVariableCount 获取规则涉及的变量数量
func (idx *VariableRuleIndex) GetVariableCount(ruleID string) int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	variables, exists := idx.ruleToVariables[ruleID]
	if !exists {
		return 0
	}
	return len(variables)
}

// GetTotalRules 获取索引中的规则总数
func (idx *VariableRuleIndex) GetTotalRules() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.ruleToVariables)
}

// GetTotalVariables 获取索引中的变量总数
func (idx *VariableRuleIndex) GetTotalVariables() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.variableToRules)
}

// FindSharedVariables 查找两个规则共享的变量
func (idx *VariableRuleIndex) FindSharedVariables(ruleID1, ruleID2 string) []uint64 {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	vars1, exists1 := idx.ruleToVariables[ruleID1]
	vars2, exists2 := idx.ruleToVariables[ruleID2]

	if !exists1 || !exists2 {
		return []uint64{}
	}

	// 创建 vars1 的集合
	set1 := make(map[uint64]bool)
	for _, v := range vars1 {
		set1[v] = true
	}

	// 查找交集
	shared := make([]uint64, 0)
	for _, v := range vars2 {
		if set1[v] {
			shared = append(shared, v)
		}
	}

	return shared
}

// Clear 清空所有索引
func (idx *VariableRuleIndex) Clear() {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.variableToRules = make(map[uint64]map[string]bool)
	idx.ruleToVariables = make(map[string][]uint64)
	idx.deviceToVariables = make(map[string][]uint64)
}

// GetStats 获取索引统计信息
type IndexStats struct {
	TotalRules       int
	TotalVariables   int
	AverageVariables float64 // 每个规则平均涉及的变量数
	AverageRules     float64 // 每个变量平均被引用的规则数
	MaxRules         int     // 单个变量被最多规则引用的数量
}

func (idx *VariableRuleIndex) GetStats() IndexStats {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	stats := IndexStats{
		TotalRules:     len(idx.ruleToVariables),
		TotalVariables: len(idx.variableToRules),
	}

	// 计算平均变量数
	totalVars := 0
	maxRules := 0

	for _, ruleMap := range idx.variableToRules {
		ruleCount := len(ruleMap)
		totalVars += ruleCount
		if ruleCount > maxRules {
			maxRules = ruleCount
		}
	}

	if stats.TotalVariables > 0 {
		stats.AverageRules = float64(totalVars) / float64(stats.TotalVariables)
	}

	if stats.TotalRules > 0 {
		stats.AverageVariables = float64(totalVars) / float64(stats.TotalRules)
	}

	stats.MaxRules = maxRules

	return stats
}
