package script

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dop251/goja"
)

// GojaEngine Goja 脚本引擎
type GojaEngine struct {
	vmPool  *VMPool
	sandbox *Sandbox
	mu      sync.RWMutex
	programs map[string]*goja.Program // 缓存编译后的脚本
}

// ExecutionContext 执行上下文
type ExecutionContext struct {
	VM       *goja.Runtime
	ScriptID string
	Input    map[string]interface{}
	Output   map[string]interface{}
	Error    error
	Start    time.Time
	End      time.Time
}

// NewGojaEngine 创建脚本引擎
func NewGojaEngine(vmPool *VMPool, sandbox *Sandbox) *GojaEngine {
	return &GojaEngine{
		vmPool:   vmPool,
		sandbox:  sandbox,
		programs: make(map[string]*goja.Program),
	}
}

// Compile 编译脚本
func (e *GojaEngine) Compile(scriptID, content string) (*goja.Program, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 检查是否已编译
	if program, exists := e.programs[scriptID]; exists {
		return program, nil
	}

	// 包装脚本以支持顶层 return 语句
	// 将脚本包装在立即执行函数中
	wrappedScript := fmt.Sprintf(`
		(function() {
			try {
				%s
			} catch(e) {
				throw e;
			}
		})()
	`, content)

	// 编译脚本
	program, err := goja.Compile("", wrappedScript, true)
	if err != nil {
		return nil, fmt.Errorf("脚本编译失败: %w", err)
	}

	// 缓存编译结果
	e.programs[scriptID] = program

	return program, nil
}

// Execute 执行脚本
func (e *GojaEngine) Execute(scriptID string, program *goja.Program, input map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	// 创建执行上下文
	ctx := &ExecutionContext{
		ScriptID: scriptID,
		Input:    input,
		Start:    time.Now(),
	}

	// 执行脚本
	if err := e.ExecuteWithContext(ctx, program, timeout); err != nil {
		return nil, err
	}

	return ctx.Output, nil
}

// ExecuteWithContext 执行脚本（带完整上下文）
func (e *GojaEngine) ExecuteWithContext(ctx *ExecutionContext, program *goja.Program, timeout time.Duration) error {
	// 从 VM 池获取 VM
	pooledVM := e.vmPool.Get()
	defer func() {
		// 归还 VM 到池中
		if err := recover(); err != nil {
			log.Printf("[脚本引擎] 脚本执行 panic: %v", err)
			ctx.Error = fmt.Errorf("脚本执行 panic: %v", err)
		}
		e.vmPool.Put(pooledVM)
	}()

	vm := pooledVM.VM
	ctx.VM = vm

	// 设置输入参数
	if ctx.Input != nil {
		for k, v := range ctx.Input {
			if err := vm.Set(k, v); err != nil {
				return fmt.Errorf("设置输入参数 %s 失败: %w", k, err)
			}
		}
	}

	// 执行脚本（带超时控制）
	resultChan := make(chan interface{}, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errorChan <- fmt.Errorf("脚本执行 panic: %v", r)
			}
		}()

		// 执行编译后的程序（已包装为立即执行函数）
		result, err := vm.RunProgram(program)
		if err != nil {
			errorChan <- fmt.Errorf("脚本执行失败: %w", err)
			return
		}

		// 导出结果
		resultChan <- result
	}()

	// 等待执行完成或超时
	select {
	case result := <-resultChan:
		// 将结果转换为 map[string]interface{}
		if result != nil {
			// 如果结果是 goja value，先转换为 Go 类型
			if exportable, ok := result.(goja.Value); ok {
				result = exportable.Export()
			}

			// 如果结果是 map，直接使用
			if resultMap, ok := result.(map[string]interface{}); ok {
				ctx.Output = resultMap
			} else {
				// 否则包装成 map
				ctx.Output = map[string]interface{}{
					"result": result,
				}
			}
		} else {
			ctx.Output = make(map[string]interface{})
		}
	case err := <-errorChan:
		ctx.Error = err
	case <-time.After(timeout):
		ctx.Error = fmt.Errorf("脚本执行超时（%v）", timeout)
	}

	ctx.End = time.Now()

	return ctx.Error
}

// GetProgram 获取已编译的脚本
func (e *GojaEngine) GetProgram(scriptID string) (*goja.Program, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	program, exists := e.programs[scriptID]
	return program, exists
}

// RemoveProgram 移除已编译的脚本
func (e *GojaEngine) RemoveProgram(scriptID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.programs, scriptID)
}

// Clear 清空所有缓存的脚本
func (e *GojaEngine) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.programs = make(map[string]*goja.Program)
}
