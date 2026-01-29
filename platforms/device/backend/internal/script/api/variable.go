package api

import (
	"fmt"
	"pansiot-device/internal/core"

	"github.com/dop251/goja"
)

// VariableAPI 变量操作 API
type VariableAPI struct {
	storage core.Storage
}

// NewVariableAPI 创建变量 API
func NewVariableAPI(storage core.Storage) *VariableAPI {
	return &VariableAPI{
		storage: storage,
	}
}

// InjectToVM 注入 API 到 VM
func (api *VariableAPI) InjectToVM(vm *goja.Runtime) error {
	variableObj := map[string]interface{}{
		"read":       api.read,
		"readString": api.readString,
		"readBatch":  api.readBatch,
		"getMeta":    api.getMeta,
		"write":      api.write,
		"writeBatch": api.writeBatch,
		"onChanged":  api.onChanged,
	}

	return vm.Set("Variable", variableObj)
}

// read 读取单个变量（数字 ID）
func (api *VariableAPI) read(variableID uint64) (interface{}, error) {
	variable, err := api.storage.ReadVar(variableID)
	if err != nil {
		return nil, fmt.Errorf("读取变量 %d 失败: %w", variableID, err)
	}
	return variable.Value, nil
}

// readString 读取单个变量（字符串 ID）
func (api *VariableAPI) readString(stringID string) (interface{}, error) {
	// 先通过字符串 ID 查找数字 ID
	// TODO: 需要实现从 stringID 到 uint64 的映射
	// 暂时返回错误
	return nil, fmt.Errorf("暂不支持通过字符串 ID 读取变量")
}

// readBatch 批量读取变量
func (api *VariableAPI) readBatch(variableIDs []uint64) (map[uint64]interface{}, error) {
	variables, err := api.storage.ReadVars(variableIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[uint64]interface{})
	for _, v := range variables {
		result[v.ID] = v.Value
	}
	return result, nil
}

// getMeta 读取变量元数据
func (api *VariableAPI) getMeta(variableID uint64) (map[string]interface{}, error) {
	variable, err := api.storage.ReadVar(variableID)
	if err != nil {
		return nil, fmt.Errorf("读取变量 %d 失败: %w", variableID, err)
	}

	meta := map[string]interface{}{
		"id":          variable.ID,
		"stringID":    variable.StringID,
		"name":        variable.Name,
		"description": variable.Description,
		"dataType":    variable.DataType.String(),
		"unit":        variable.Unit,
		"deviceID":    variable.DeviceID,
	}

	return meta, nil
}

// write 写入单个变量
func (api *VariableAPI) write(variableID uint64, value interface{}) error {
	variable := &core.Variable{
		ID:    variableID,
		Value: value,
	}
	return api.storage.WriteVar(variable)
}

// writeBatch 批量写入变量
func (api *VariableAPI) writeBatch(values map[uint64]interface{}) error {
	variables := make([]*core.Variable, 0, len(values))
	for id, val := range values {
		variables = append(variables, &core.Variable{
			ID:    id,
			Value: val,
		})
	}
	return api.storage.WriteVars(variables)
}

// onChanged 订阅变量变化（在脚本中定义回调）
func (api *VariableAPI) onChanged(variableID uint64, callback func(interface{}, interface{}, interface{})) error {
	// TODO: 实现变量变化订阅
	// 这需要与 ScriptConsumer 集成
	return fmt.Errorf("暂未实现")
}
