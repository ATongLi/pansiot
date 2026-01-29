package api_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"pansiot-device/internal/script/api"
)

// TestFileAPI 测试文件API
func TestFileAPI(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	systemAPI := api.NewSystemAPI(tmpDir, false)
	defer systemAPI.Close()

	vm := goja.New()
	err := systemAPI.InjectToVM(vm)
	if err != nil {
		t.Fatalf("注入API失败: %v", err)
	}

	t.Run("Write and Read", func(t *testing.T) {
		path := "test.txt"
		content := "Hello, World!"

		_, err := vm.RunString(`System.File.writeText("` + path + `", "` + content + `")`)
		if err != nil {
			t.Fatalf("写入失败: %v", err)
		}

		result, err := vm.RunString(`System.File.readText("` + path + `")`)
		if err != nil {
			t.Fatalf("读取失败: %v", err)
		}

		read := result.String()
		if read != content {
			t.Errorf("内容不匹配: 期望 %q, 实际 %q", content, read)
		}
	})

	t.Run("Append", func(t *testing.T) {
		path := "append.txt"
		vm.RunString(`System.File.writeText("` + path + `", "Line1\n")`)
		vm.RunString(`System.File.appendText("` + path + `", "Line2\n")`)

		result, err := vm.RunString(`System.File.readText("` + path + `")`)
		if err != nil {
			t.Fatalf("读取失败: %v", err)
		}

		content := result.String()
		expected := "Line1\nLine2\n"
		if content != expected {
			t.Errorf("追加内容不匹配: 期望 %q, 实际 %q", expected, content)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		vm.RunString(`System.File.writeText("exists.txt", "test")`)

		result1, _ := vm.RunString(`System.File.exists("exists.txt")`)
		if !result1.ToBoolean() {
			t.Error("Exists()应该返回true")
		}

		result2, _ := vm.RunString(`System.File.exists("nonexistent.txt")`)
		if result2.ToBoolean() {
			t.Error("Exists()应该返回false")
		}
	})

	t.Run("PathTraversal", func(t *testing.T) {
		_, err := vm.RunString(`System.File.readText("../../../etc/passwd")`)
		if err == nil {
			t.Error("应该拒绝路径遍历攻击")
		}
	})

	t.Run("FileSize", func(t *testing.T) {
		vm.RunString(`System.File.writeText("size.txt", "12345")`)
		result, err := vm.RunString(`System.File.size("size.txt")`)
		if err != nil {
			t.Fatalf("获取文件大小失败: %v", err)
		}
		size := result.ToInteger()
		if size != 5 {
			t.Errorf("文件大小错误: %d, 期望 5", size)
		}
	})

	t.Run("List", func(t *testing.T) {
		vm.RunString(`System.File.writeText("file1.txt", "content1")`)
		vm.RunString(`System.File.writeText("file2.txt", "content2")`)

		result, err := vm.RunString(`System.File.list("")`)
		if err != nil {
			t.Fatalf("列出目录失败: %v", err)
		}

		// result.Export() 可以返回 []string 或 []interface{}
		files := make([]string, 0)
		switch v := result.Export().(type) {
		case []string:
			files = v
		case []interface{}:
			for _, item := range v {
				files = append(files, fmt.Sprintf("%v", item))
			}
		}

		if len(files) < 2 {
			t.Errorf("文件数量不足: %d", len(files))
		}
	})
}

// TestStringAPI 测试字符串API
func TestStringAPI(t *testing.T) {
	vm := goja.New()
	systemAPI := api.NewSystemAPI("", false)
	defer systemAPI.Close()
	systemAPI.InjectToVM(vm)

	t.Run("Split", func(t *testing.T) {
		result, err := vm.RunString(`String.split("a,b,c", ",")`)
		if err != nil {
			t.Fatal(err)
		}

		// 处理 []string 或 []interface{}
		var arr []interface{}
		switch v := result.Export().(type) {
		case []string:
			arr = make([]interface{}, len(v))
			for i, s := range v {
				arr[i] = s
			}
		case []interface{}:
			arr = v
		}

		if len(arr) != 3 {
			t.Errorf("分割结果错误: 期望3个元素, 实际%d个", len(arr))
		}
		if fmt.Sprintf("%v", arr[0]) != "a" || fmt.Sprintf("%v", arr[1]) != "b" || fmt.Sprintf("%v", arr[2]) != "c" {
			t.Errorf("分割内容错误: %v", arr)
		}
	})

	t.Run("Join", func(t *testing.T) {
		result, err := vm.RunString(`String.join(["a", "b", "c"], "-")`)
		if err != nil {
			t.Fatal(err)
		}
		joined := result.String()
		if joined != "a-b-c" {
			t.Errorf("拼接失败: 期望 'a-b-c', 实际 '%s'", joined)
		}
	})

	t.Run("Replace", func(t *testing.T) {
		result, err := vm.RunString(`String.replace("Hello World", "World", "Go")`)
		if err != nil {
			t.Fatal(err)
		}
		if result.String() != "Hello Go" {
			t.Errorf("替换失败: 期望 'Hello Go', 实际 '%s'", result.String())
		}
	})

	t.Run("ReplaceAll", func(t *testing.T) {
		result, err := vm.RunString(`String.replaceAll("aaa", "a", "b")`)
		if err != nil {
			t.Fatal(err)
		}
		if result.String() != "bbb" {
			t.Errorf("全部替换失败: 期望 'bbb', 实际 '%s'", result.String())
		}
	})

	t.Run("Trim", func(t *testing.T) {
		result, err := vm.RunString(`String.trim("  hello  ")`)
		if err != nil {
			t.Fatal(err)
		}
		if result.String() != "hello" {
			t.Errorf("去空格失败: 期望 'hello', 实际 '%s'", result.String())
		}
	})

	t.Run("ToUpper", func(t *testing.T) {
		result, err := vm.RunString(`String.toUpper("hello")`)
		if err != nil {
			t.Fatal(err)
		}
		if result.String() != "HELLO" {
			t.Errorf("转大写失败: 期望 'HELLO', 实际 '%s'", result.String())
		}
	})

	t.Run("ToLower", func(t *testing.T) {
		result, err := vm.RunString(`String.toLower("HELLO")`)
		if err != nil {
			t.Fatal(err)
		}
		if result.String() != "hello" {
			t.Errorf("转小写失败: 期望 'hello', 实际 '%s'", result.String())
		}
	})

	t.Run("Contains", func(t *testing.T) {
		result, err := vm.RunString(`String.contains("Hello World", "World")`)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("Contains应该返回true")
		}
	})

	t.Run("Substring", func(t *testing.T) {
		result, err := vm.RunString(`String.substring("Hello World", 0, 5)`)
		if err != nil {
			t.Fatal(err)
		}
		if result.String() != "Hello" {
			t.Errorf("子字符串失败: 期望 'Hello', 实际 '%s'", result.String())
		}
	})
}

// TestJSONAPI 测试JSON API
func TestJSONAPI(t *testing.T) {
	vm := goja.New()
	systemAPI := api.NewSystemAPI("", false)
	defer systemAPI.Close()
	systemAPI.InjectToVM(vm)

	t.Run("Parse", func(t *testing.T) {
		result, err := vm.RunString(`JSON.parse('{"name":"test","value":123}')`)
		if err != nil {
			t.Fatal(err)
		}
		obj := result.Export().(map[string]interface{})
		if obj["name"] != "test" {
			t.Errorf("JSON解析失败: name字段错误")
		}
		if obj["value"].(float64) != 123 {
			t.Errorf("JSON解析失败: value字段错误")
		}
	})

	t.Run("Stringify", func(t *testing.T) {
		result, err := vm.RunString(`JSON.stringify({name:"test",value:123})`)
		if err != nil {
			t.Fatal(err)
		}
		jsonStr := result.String()
		if !strings.Contains(jsonStr, "test") {
			t.Error("JSON序列化失败")
		}
		if !strings.Contains(jsonStr, "123") {
			t.Error("JSON序列化失败")
		}
	})

	t.Run("ParseError", func(t *testing.T) {
		_, err := vm.RunString(`JSON.parse('invalid json')`)
		if err == nil {
			t.Error("应该返回JSON解析错误")
		}
	})
}

// TestDateAPI 测试日期API
func TestDateAPI(t *testing.T) {
	vm := goja.New()
	systemAPI := api.NewSystemAPI("", false)
	defer systemAPI.Close()
	systemAPI.InjectToVM(vm)

	t.Run("Format", func(t *testing.T) {
		// 2024-01-15 12:30:45.123
		timestamp := int64(1705317045123)
		result, err := vm.RunString(`Date.format(` + fmt.Sprintf("%d", timestamp) + `, "YYYY-MM-DD HH:mm:ss")`)
		if err != nil {
			t.Fatal(err)
		}
		formatted := result.String()
		t.Logf("格式化结果: %s", formatted)
		// 验证格式
		if !strings.Contains(formatted, "-") {
			t.Error("日期格式错误: 应该包含 '-'")
		}
		if !strings.Contains(formatted, ":") {
			t.Error("时间格式错误: 应该包含 ':'")
		}
	})

	t.Run("FormatShort", func(t *testing.T) {
		timestamp := int64(1705317045123)
		result, err := vm.RunString(`Date.format(` + fmt.Sprintf("%d", timestamp) + `, "YYYY-MM-DD")`)
		if err != nil {
			t.Fatal(err)
		}
		formatted := result.String()
		t.Logf("短格式化结果: %s", formatted)
		if !strings.Contains(formatted, "-") {
			t.Error("日期格式错误")
		}
	})

	t.Run("ParseFormat", func(t *testing.T) {
		result, err := vm.RunString(`Date.parseFormat("2024-01-15 12:30:45", "YYYY-MM-DD HH:mm:ss")`)
		if err != nil {
			t.Fatal(err)
		}
		timestamp := result.Export().(int64)
		t.Logf("解析时间戳: %d", timestamp)
		if timestamp == 0 {
			t.Error("时间戳解析失败")
		}
	})
}

// TestLogAPI 测试日志API
func TestLogAPI(t *testing.T) {
	tmpDir := t.TempDir()
	systemAPI := api.NewSystemAPI(tmpDir, true)
	defer systemAPI.Close()

	vm := goja.New()
	systemAPI.InjectToVM(vm)

	t.Run("Info", func(t *testing.T) {
		_, err := vm.RunString(`Log.info("测试信息日志")`)
		if err != nil {
			t.Errorf("Log.info失败: %v", err)
		}
	})

	t.Run("Warn", func(t *testing.T) {
		_, err := vm.RunString(`Log.warn("测试警告日志")`)
		if err != nil {
			t.Errorf("Log.warn失败: %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		_, err := vm.RunString(`Log.error("测试错误日志")`)
		if err != nil {
			t.Errorf("Log.error失败: %v", err)
		}
	})

	t.Run("Debug", func(t *testing.T) {
		_, err := vm.RunString(`Log.debug("测试调试日志")`)
		if err != nil {
			t.Errorf("Log.debug失败: %v", err)
		}
	})

	// 检查日志文件是否创建
	logPath := tmpDir + "/script.log"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("日志文件未创建")
	}
}

// TestSystemAPIIntegration 测试系统API集成
func TestSystemAPIIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	systemAPI := api.NewSystemAPI(tmpDir, true)
	defer systemAPI.Close()

	vm := goja.New()
	err := systemAPI.InjectToVM(vm)
	if err != nil {
		t.Fatalf("注入API失败: %v", err)
	}

	// 测试脚本:文件操作 + 字符串 + JSON + 日志
	script := `
		// 写入配置文件
		System.File.writeText("config.json", '{"key":"value"}');

		// 读取配置
		var content = System.File.readText("config.json");

		// 解析JSON
		var config = JSON.parse(content);

		// 字符串操作
		var upper = String.toUpper(config.key);

		// 记录日志
		Log.info("配置加载成功: " + upper);

		// 返回结果
		JSON.stringify({status:"ok", value:upper});
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("脚本执行失败: %v", err)
	}

	resultStr := result.String()
	if !strings.Contains(resultStr, "VALUE") {
		t.Errorf("集成测试失败: 结果不正确 %s", resultStr)
	}
}
