package script

import (
	"pansiot-device/internal/core"
	"pansiot-device/internal/script/api"

	"github.com/dop251/goja"
)

// Sandbox 沙箱管理器
type Sandbox struct {
	allowedAPIs map[string]bool // API 白名单
	vmPool      *VMPool
}

// APICategory API 分类
type APICategory string

const (
	APICatVariable      APICategory = "variable"       // 变量操作
	APICatSystem        APICategory = "system"         // 系统指令
	APICatCommunication APICategory = "communication"  // 通讯功能
	APICatData          APICategory = "data"           // 数据处理
	APICatUI            APICategory = "ui"             // 界面控制
)

// NewSandbox 创建沙箱
func NewSandbox(vmPool *VMPool) *Sandbox {
	return &Sandbox{
		allowedAPIs: map[string]bool{
			"variable":      true,
			"system":        true,
			"communication": true,
			"data":          true,
			"ui":            true,
		},
		vmPool: vmPool,
	}
}

// SetupVM 为 VM 设置沙箱环境（注入 API）
func (s *Sandbox) SetupVM(vm *goja.Runtime, storage core.Storage) error {
	// 移除危险的全局对象
	if err := vm.Set("require", nil); err != nil {
		return err
	}
	if err := vm.Set("eval", nil); err != nil {
		return err
	}
	// 注意：不能完全移除 Function，因为它是 JavaScript 的核心部分
	// 但我们可以限制其使用

	// 注入允许的 API
	if s.allowedAPIs["variable"] {
		if err := s.InjectVariableAPI(vm, storage); err != nil {
			return err
		}
	}

	if s.allowedAPIs["system"] {
		if err := s.InjectSystemAPI(vm); err != nil {
			return err
		}
	}

	if s.allowedAPIs["communication"] {
		if err := s.InjectCommunicationAPI(vm); err != nil {
			return err
		}
	}

	return nil
}

// InjectVariableAPI 注入变量操作 API
func (s *Sandbox) InjectVariableAPI(vm *goja.Runtime, storage core.Storage) error {
	variableAPI := api.NewVariableAPI(storage)
	return variableAPI.InjectToVM(vm)
}

// InjectSystemAPI 注入系统指令 API
// Phase 3: 系统指令API - 文件、日志、字符串、JSON、时间日期
func (s *Sandbox) InjectSystemAPI(vm *goja.Runtime) error {
	// 创建系统API实例,基础路径为 data/scripts/
	// 日志默认启用
	systemAPI := api.NewSystemAPI("data/scripts/", true)
	return systemAPI.InjectToVM(vm)
}

// InjectCommunicationAPI 注入通讯功能 API
// Phase 4: 通讯功能API - HTTP、MQTT、Modbus
func (s *Sandbox) InjectCommunicationAPI(vm *goja.Runtime) error {
	commAPI := api.NewCommunicationAPI()
	return commAPI.InjectToVM(vm)
}

// SetAllowedAPIs 设置允许的 API
func (s *Sandbox) SetAllowedAPIs(apis map[string]bool) {
	s.allowedAPIs = apis
}

// IsAllowed 检查 API 是否允许
func (s *Sandbox) IsAllowed(api string) bool {
	return s.allowedAPIs[api]
}

// Validate 验证脚本是否只使用允许的 API
func (s *Sandbox) Validate(program *goja.Program) error {
	// TODO: 实现 AST 分析，检查脚本使用的 API
	// 这需要解析 JavaScript AST，比较复杂
	// 暂时跳过，依赖运行时检查

	return nil
}
