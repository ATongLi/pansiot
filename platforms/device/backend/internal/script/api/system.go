package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
)

// SystemAPI 系统API(文件、日志、时间、字符串、JSON)
// 系统指令API - Phase 3
type SystemAPI struct {
	baseDir     string        // 安全基础路径
	maxFileSize int64         // 最大文件大小(字节)
	allowDelete bool          // 是否允许删除操作
	logEnabled  bool          // 日志开关
	logFile     *os.File      // 日志文件句柄
	logWriter   *bufio.Writer // 日志写入缓冲
	logMutex    sync.Mutex    // 日志写入锁
}

// NewSystemAPI 创建系统API实例
// baseDir: 安全基础路径,所有文件操作限制在此目录下
// logEnabled: 是否启用日志
func NewSystemAPI(baseDir string, logEnabled bool) *SystemAPI {
	// 确保基础目录存在
	os.MkdirAll(baseDir, 0755)

	api := &SystemAPI{
		baseDir:     baseDir,
		maxFileSize: 10 * 1024 * 1024, // 10MB
		allowDelete: false,             // 默认不允许删除
		logEnabled:  logEnabled,
	}

	// 初始化日志
	if logEnabled {
		if err := api.initLogger(); err != nil {
			// 日志初始化失败不应该阻止API创建
			fmt.Printf("[SystemAPI] 日志初始化失败: %v\n", err)
		}
	}

	return api
}

// InjectToVM 注入API到VM
// 将所有系统API注册到Goja运行时
func (api *SystemAPI) InjectToVM(vm *goja.Runtime) error {
	// 注入 System.File 对象
	fileAPI := api.createFileAPI(vm)
	systemObj := map[string]interface{}{
		"File": fileAPI,
	}
	vm.Set("System", systemObj)

	// 注入 Log API
	logAPI := api.createLogAPI(vm)
	vm.Set("Log", logAPI)

	// 注入 String API
	stringAPI := api.createStringAPI(vm)
	vm.Set("String", stringAPI)

	// 注入 JSON API
	jsonAPI := api.createJSONAPI(vm)
	vm.Set("JSON", jsonAPI)

	// 扩展Date对象
	if err := api.extendDateObject(vm); err != nil {
		return fmt.Errorf("扩展Date对象失败: %v", err)
	}

	return nil
}

// ==================== 文件操作 API ====================

// createFileAPI 创建文件API
func (api *SystemAPI) createFileAPI(vm *goja.Runtime) map[string]interface{} {
	return map[string]interface{}{
		"readText": func(call goja.FunctionCall) goja.Value {
			path := call.Argument(0).String()
			content, err := api.readText(path)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(content)
		},
		"writeText": func(call goja.FunctionCall) goja.Value {
			path := call.Argument(0).String()
			content := call.Argument(1).String()
			err := api.writeText(path, content)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return goja.Undefined()
		},
		"appendText": func(call goja.FunctionCall) goja.Value {
			path := call.Argument(0).String()
			content := call.Argument(1).String()
			err := api.appendText(path, content)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return goja.Undefined()
		},
		"exists": func(call goja.FunctionCall) goja.Value {
			path := call.Argument(0).String()
			exists := api.exists(path)
			return vm.ToValue(exists)
		},
		"list": func(call goja.FunctionCall) goja.Value {
			path := call.Argument(0).String()
			files, err := api.list(path)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(files)
		},
		"size": func(call goja.FunctionCall) goja.Value {
			path := call.Argument(0).String()
			size, err := api.size(path)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(size)
		},
		"delete": func(call goja.FunctionCall) goja.Value {
			if !api.allowDelete {
				panic(vm.NewGoError(fmt.Errorf("文件删除操作未启用")))
			}
			path := call.Argument(0).String()
			err := api.delete(path)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return goja.Undefined()
		},
	}
}

// readText 读取文本文件
func (api *SystemAPI) readText(path string) (string, error) {
	safePath, err := api.validatePath(path)
	if err != nil {
		return "", err
	}

	// 检查文件大小
	info, err := os.Stat(safePath)
	if err != nil {
		return "", err
	}
	if info.Size() > api.maxFileSize {
		return "", fmt.Errorf("文件过大: %d 字节 (最大 %d 字节)", info.Size(), api.maxFileSize)
	}

	content, err := os.ReadFile(safePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	return string(content), nil
}

// writeText 写入文本文件(覆盖模式)
func (api *SystemAPI) writeText(path, content string) error {
	safePath, err := api.validatePath(path)
	if err != nil {
		return err
	}

	// 检查内容大小
	if int64(len(content)) > api.maxFileSize {
		return fmt.Errorf("内容过大: %d 字节 (最大 %d 字节)", len(content), api.maxFileSize)
	}

	// 确保目录存在
	dir := filepath.Dir(safePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	if err := os.WriteFile(safePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// appendText 追加文本到文件
func (api *SystemAPI) appendText(path, content string) error {
	safePath, err := api.validatePath(path)
	if err != nil {
		return err
	}

	// 检查内容大小
	if int64(len(content)) > api.maxFileSize {
		return fmt.Errorf("内容过大: %d 字节 (最大 %d 字节)", len(content), api.maxFileSize)
	}

	// 确保目录存在
	dir := filepath.Dir(safePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	file, err := os.OpenFile(safePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("追加内容失败: %v", err)
	}

	return nil
}

// exists 检查文件是否存在
func (api *SystemAPI) exists(path string) bool {
	safePath, err := api.validatePath(path)
	if err != nil {
		return false
	}
	_, err = os.Stat(safePath)
	return err == nil
}

// list 列出目录内容
func (api *SystemAPI) list(path string) ([]string, error) {
	safePath, err := api.validatePath(path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(safePath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	return files, nil
}

// size 获取文件大小(字节)
func (api *SystemAPI) size(path string) (int64, error) {
	safePath, err := api.validatePath(path)
	if err != nil {
		return 0, err
	}

	info, err := os.Stat(safePath)
	if err != nil {
		return 0, fmt.Errorf("获取文件信息失败: %v", err)
	}

	return info.Size(), nil
}

// delete 删除文件
func (api *SystemAPI) delete(path string) error {
	safePath, err := api.validatePath(path)
	if err != nil {
		return err
	}

	if err := os.Remove(safePath); err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// validatePath 验证并规范化路径(安全检查)
// 防止路径遍历攻击
func (api *SystemAPI) validatePath(path string) (string, error) {
	// 1. 检查路径遍历
	if strings.Contains(path, "..") {
		return "", fmt.Errorf("禁止路径遍历: %s", path)
	}

	// 2. 转换为绝对路径
	absPath := filepath.Join(api.baseDir, path)

	// 3. 检查是否在基础路径内
	relPath, err := filepath.Rel(api.baseDir, absPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("路径超出安全范围: %s", path)
	}

	return absPath, nil
}

// ==================== 日志记录 API ====================

// createLogAPI 创建日志API
func (api *SystemAPI) createLogAPI(vm *goja.Runtime) map[string]interface{} {
	return map[string]interface{}{
		"info": func(call goja.FunctionCall) goja.Value {
			msg := api.formatMessage(call.Arguments)
			api.log("[INFO] " + msg)
			return goja.Undefined()
		},
		"warn": func(call goja.FunctionCall) goja.Value {
			msg := api.formatMessage(call.Arguments)
			api.log("[WARN] " + msg)
			return goja.Undefined()
		},
		"error": func(call goja.FunctionCall) goja.Value {
			msg := api.formatMessage(call.Arguments)
			api.log("[ERROR] " + msg)
			return goja.Undefined()
		},
		"debug": func(call goja.FunctionCall) goja.Value {
			msg := api.formatMessage(call.Arguments)
			api.log("[DEBUG] " + msg)
			return goja.Undefined()
		},
	}
}

// formatMessage 格式化日志消息
func (api *SystemAPI) formatMessage(args []goja.Value) string {
	if len(args) == 0 {
		return ""
	}

	// 简单实现:直接拼接参数
	var parts []string
	for _, arg := range args {
		parts = append(parts, arg.String())
	}
	return strings.Join(parts, " ")
}

// initLogger 初始化日志系统
func (api *SystemAPI) initLogger() error {
	logPath := filepath.Join(api.baseDir, "script.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	api.logFile = file
	api.logWriter = bufio.NewWriter(file)
	return nil
}

// log 写入日志
func (api *SystemAPI) log(msg string) {
	if !api.logEnabled {
		return
	}

	api.logMutex.Lock()
	defer api.logMutex.Unlock()

	// 格式: [时间] [级别] 消息
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, msg)

	// 写入文件
	if api.logWriter != nil {
		api.logWriter.WriteString(logLine)
		api.logWriter.Flush()
	}

	// 同时输出到控制台
	fmt.Print(logLine)
}

// Close 关闭日志文件
func (api *SystemAPI) Close() error {
	api.logMutex.Lock()
	defer api.logMutex.Unlock()

	if api.logWriter != nil {
		api.logWriter.Flush()
	}

	if api.logFile != nil {
		return api.logFile.Close()
	}

	return nil
}

// ==================== 字符串处理 API ====================

// createStringAPI 创建字符串API
func (api *SystemAPI) createStringAPI(vm *goja.Runtime) map[string]interface{} {
	return map[string]interface{}{
		"split": func(call goja.FunctionCall) goja.Value {
			str := call.Argument(0).String()
			sep := call.Argument(1).String()
			parts := strings.Split(str, sep)
			return vm.ToValue(parts)
		},
		"join": func(call goja.FunctionCall) goja.Value {
			arr := call.Argument(0).Export()
			sep := call.Argument(1).String()

			var strParts []string
			switch v := arr.(type) {
			case []interface{}:
				for _, item := range v {
					strParts = append(strParts, fmt.Sprintf("%v", item))
				}
			case []string:
				strParts = v
			default:
				return vm.ToValue([]string{})
			}

			result := strings.Join(strParts, sep)
			return vm.ToValue(result)
		},
		"replace": func(call goja.FunctionCall) goja.Value {
			str := call.Argument(0).String()
			old := call.Argument(1).String()
			newStr := call.Argument(2).String()
			result := strings.Replace(str, old, newStr, 1)
			return vm.ToValue(result)
		},
		"replaceAll": func(call goja.FunctionCall) goja.Value {
			str := call.Argument(0).String()
			old := call.Argument(1).String()
			newStr := call.Argument(2).String()
			result := strings.ReplaceAll(str, old, newStr)
			return vm.ToValue(result)
		},
		"trim": func(call goja.FunctionCall) goja.Value {
			str := call.Argument(0).String()
			result := strings.TrimSpace(str)
			return vm.ToValue(result)
		},
		"toUpper": func(call goja.FunctionCall) goja.Value {
			str := call.Argument(0).String()
			result := strings.ToUpper(str)
			return vm.ToValue(result)
		},
		"toLower": func(call goja.FunctionCall) goja.Value {
			str := call.Argument(0).String()
			result := strings.ToLower(str)
			return vm.ToValue(result)
		},
		"contains": func(call goja.FunctionCall) goja.Value {
			str := call.Argument(0).String()
			substr := call.Argument(1).String()
			result := strings.Contains(str, substr)
			return vm.ToValue(result)
		},
		"substring": func(call goja.FunctionCall) goja.Value {
			str := call.Argument(0).String()
			start := int(call.Argument(1).ToInteger())
			end := int(call.Argument(2).ToInteger())

			// 边界检查
			if start < 0 {
				start = 0
			}
			if end > len(str) {
				end = len(str)
			}
			if start > end {
				start, end = end, start
			}

			result := str[start:end]
			return vm.ToValue(result)
		},
	}
}

// ==================== JSON 处理 API ====================

// createJSONAPI 创建JSON API
func (api *SystemAPI) createJSONAPI(vm *goja.Runtime) map[string]interface{} {
	return map[string]interface{}{
		"parse": func(call goja.FunctionCall) goja.Value {
			jsonStr := call.Argument(0).String()
			var result interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				panic(vm.NewGoError(fmt.Errorf("JSON解析失败: %v", err)))
			}
			return vm.ToValue(result)
		},
		"stringify": func(call goja.FunctionCall) goja.Value {
			value := call.Argument(0).Export()
			data, err := json.Marshal(value)
			if err != nil {
				panic(vm.NewGoError(fmt.Errorf("JSON序列化失败: %v", err)))
			}
			return vm.ToValue(string(data))
		},
	}
}

// ==================== 时间日期 API ====================

// extendDateObject 扩展Date对象
func (api *SystemAPI) extendDateObject(vm *goja.Runtime) error {
	// 获取Date构造函数
	dateObj := vm.Get("Date")
	if dateObj == nil || goja.IsUndefined(dateObj) {
		return fmt.Errorf("Date对象不存在")
	}

	dateConstructor := dateObj.ToObject(vm)

	// 添加静态方法 Date.format()
	formatFunc := func(call goja.FunctionCall) goja.Value {
		timestamp := int64(call.Argument(0).ToInteger())
		format := call.Argument(1).String()

		// 转换时间戳
		t := time.UnixMilli(timestamp)

		// 格式化
		result := api.formatDate(t, format)
		return vm.ToValue(result)
	}

	if err := dateConstructor.Set("format", formatFunc); err != nil {
		return fmt.Errorf("设置Date.format失败: %v", err)
	}

	// 添加静态方法 Date.parseFormat()
	parseFormatFunc := func(call goja.FunctionCall) goja.Value {
		dateStr := call.Argument(0).String()
		format := call.Argument(1).String()

		// 解析时间字符串
		t, err := api.parseDate(dateStr, format)
		if err != nil {
			panic(vm.NewGoError(err))
		}

		return vm.ToValue(t.UnixMilli())
	}

	if err := dateConstructor.Set("parseFormat", parseFormatFunc); err != nil {
		return fmt.Errorf("设置Date.parseFormat失败: %v", err)
	}

	return nil
}

// formatDate 格式化时间
func (api *SystemAPI) formatDate(t time.Time, format string) string {
	result := format
	result = strings.ReplaceAll(result, "YYYY", t.Format("2006"))
	result = strings.ReplaceAll(result, "YY", t.Format("06"))
	result = strings.ReplaceAll(result, "MM", t.Format("01"))
	result = strings.ReplaceAll(result, "DD", t.Format("02"))
	result = strings.ReplaceAll(result, "HH", t.Format("15"))
	result = strings.ReplaceAll(result, "mm", t.Format("04"))
	result = strings.ReplaceAll(result, "ss", t.Format("05"))
	result = strings.ReplaceAll(result, "SSS", fmt.Sprintf("%03d", t.Nanosecond()/1e6))
	return result
}

// parseDate 解析时间字符串(简化版)
func (api *SystemAPI) parseDate(dateStr, format string) (time.Time, error) {
	// 支持的格式映射
	layouts := map[string]string{
		"YYYY-MM-DD HH:mm:ss": "2006-01-02 15:04:05",
		"YYYY-MM-DD":          "2006-01-02",
		"YYYY/MM/DD":          "2006/01/02",
		"HH:mm:ss":            "15:04:05",
		"HH:mm":               "15:04",
	}

	if layout, ok := layouts[format]; ok {
		return time.Parse(layout, dateStr)
	}

	return time.Time{}, fmt.Errorf("不支持的时间格式: %s", format)
}
