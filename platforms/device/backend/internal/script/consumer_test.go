package script

import (
	"context"
	"testing"
	"time"

	"pansiot-device/internal/core"
	"pansiot-device/internal/storage"
)

// TestScriptLoadAndExecute 测试脚本加载和执行
func TestScriptLoadAndExecute(t *testing.T) {
	// 1. 创建存储层
	strg := storage.NewMemoryStorage()

	// 2. 创建测试变量
	testVar := &core.Variable{
		ID:       100001,
		StringID: "DV-TEST-TEMP",
		Name:     "测试温度",
		Value:    75.0,
	}
	if err := strg.CreateVariable(testVar); err != nil {
		t.Fatalf("创建变量失败: %v", err)
	}

	// 3. 创建脚本消费者
	config := DefaultScriptConfig()
	consumer := NewScriptConsumer("test-script", strg, config)

	// 4. 启动消费者
	ctx := context.Background()
	if err := consumer.Start(ctx); err != nil {
		t.Fatalf("启动失败: %v", err)
	}
	defer consumer.Stop()

	// 5. 加载脚本
	script := &Script{
		ID:      "TEST_SCRIPT",
		Name:    "测试脚本",
		Content: `
			// 简单的数学运算
			var result = 10 + 20;
			return { status: "ok", value: result };
		`,
		Enabled: true,
	}

	if err := consumer.LoadScript(script); err != nil {
		t.Fatalf("加载脚本失败: %v", err)
	}

	// 6. 执行脚本
	result, err := consumer.ExecuteScript("TEST_SCRIPT", nil)
	if err != nil {
		t.Fatalf("执行脚本失败: %v", err)
	}

	// 7. 验证结果
	if result["status"] != "ok" {
		t.Errorf("期望状态 'ok'，得到 '%v'", result["status"])
	}

	// 检查值是否为 30（支持 int 和 float64 类型）
	value := result["value"]
	var valueNum float64
	switch v := value.(type) {
	case int:
		valueNum = float64(v)
	case int64:
		valueNum = float64(v)
	case float64:
		valueNum = v
	case float32:
		valueNum = float64(v)
	default:
		t.Errorf("期望值类型为数字，得到 %T", value)
	}

	if valueNum != 30 {
		t.Errorf("期望值 30，得到 %v", value)
	}

	t.Logf("脚本执行成功，结果: %v", result)
}

// TestVMPool 测试 VM 池
func TestVMPool(t *testing.T) {
	sandbox := NewSandbox(nil)
	pool := NewVMPool(5, 5*time.Minute, 30*time.Minute, sandbox)

	// 获取 VM
	vm1 := pool.Get()
	if vm1 == nil {
		t.Fatal("获取 VM 失败")
	}

	if vm1.VM == nil {
		t.Fatal("VM 实例为空")
	}

	// 归还 VM
	if err := pool.Put(vm1); err != nil {
		t.Fatalf("归还 VM 失败: %v", err)
	}

	// 再次获取，应该复用
	vm2 := pool.Get()
	if vm2 == nil {
		t.Fatal("获取 VM 失败")
	}

	// 验证复用
	if vm1.VM != vm2.VM {
		t.Error("VM 应该被复用")
	}

	stats := pool.GetStats()
	if stats.TotalReused != 1 {
		t.Errorf("期望复用 1 次，得到 %d", stats.TotalReused)
	}

	if stats.TotalCreated != 1 {
		t.Errorf("期望创建 1 个，得到 %d", stats.TotalCreated)
	}

	// 归还 VM
	pool.Put(vm2)

	t.Logf("VM 池测试成功，统计: %+v", stats)
}

// TestScriptAsyncExecution 测试异步脚本执行
func TestScriptAsyncExecution(t *testing.T) {
	// 1. 创建存储层
	strg := storage.NewMemoryStorage()

	// 2. 创建脚本消费者
	config := DefaultScriptConfig()
	consumer := NewScriptConsumer("test-async", strg, config)

	// 3. 启动消费者
	ctx := context.Background()
	if err := consumer.Start(ctx); err != nil {
		t.Fatalf("启动失败: %v", err)
	}
	defer consumer.Stop()

	// 4. 加载脚本
	script := &Script{
		ID:      "ASYNC_SCRIPT",
		Name:    "异步测试脚本",
		Content: `
			return { async: true, value: 42 };
		`,
		Enabled: true,
	}

	if err := consumer.LoadScript(script); err != nil {
		t.Fatalf("加载脚本失败: %v", err)
	}

	// 5. 异步执行脚本
	if err := consumer.ExecuteScriptAsync("ASYNC_SCRIPT", nil); err != nil {
		t.Fatalf("异步执行脚本失败: %v", err)
	}

	// 等待执行完成
	time.Sleep(100 * time.Millisecond)

	// 6. 检查脚本状态
	status := consumer.GetScriptStatus("ASYNC_SCRIPT")
	if status == nil {
		t.Fatal("未找到脚本状态")
	}

	if status.ExecCount != 1 {
		t.Errorf("期望执行 1 次，得到 %d", status.ExecCount)
	}

	t.Logf("异步执行测试成功，状态: %+v", status)
}

// TestScriptStatus 测试脚本状态
func TestScriptStatus(t *testing.T) {
	// 1. 创建存储层
	strg := storage.NewMemoryStorage()

	// 2. 创建脚本消费者
	config := DefaultScriptConfig()
	consumer := NewScriptConsumer("test-status", strg, config)

	// 3. 启动消费者
	ctx := context.Background()
	if err := consumer.Start(ctx); err != nil {
		t.Fatalf("启动失败: %v", err)
	}
	defer consumer.Stop()

	// 4. 加载脚本
	script := &Script{
		ID:      "STATUS_SCRIPT",
		Name:    "状态测试脚本",
		Content: `return { status: "test" };`,
		Enabled: true,
	}

	if err := consumer.LoadScript(script); err != nil {
		t.Fatalf("加载脚本失败: %v", err)
	}

	// 5. 检查状态
	status := consumer.GetScriptStatus("STATUS_SCRIPT")
	if status == nil {
		t.Fatal("未找到脚本状态")
	}

	if !status.Loaded {
		t.Error("脚本应该标记为已加载")
	}

	if !status.Enabled {
		t.Error("脚本应该标记为已启用")
	}

	if status.State != ScriptStateIdle {
		t.Errorf("期望状态 %d，得到 %d", ScriptStateIdle, status.State)
	}

	t.Logf("脚本状态测试成功，状态: %+v", status)
}

// TestListScripts 测试列出所有脚本
func TestListScripts(t *testing.T) {
	// 1. 创建存储层
	strg := storage.NewMemoryStorage()

	// 2. 创建脚本消费者
	config := DefaultScriptConfig()
	consumer := NewScriptConsumer("test-list", strg, config)

	// 3. 加载多个脚本
	scripts := []*Script{
		{
			ID:      "SCRIPT_1",
			Name:    "脚本1",
			Content: `return 1;`,
			Enabled: true,
		},
		{
			ID:      "SCRIPT_2",
			Name:    "脚本2",
			Content: `return 2;`,
			Enabled: true,
		},
		{
			ID:      "SCRIPT_3",
			Name:    "脚本3",
			Content: `return 3;`,
			Enabled: false,
		},
	}

	for _, script := range scripts {
		if err := consumer.LoadScript(script); err != nil {
			t.Fatalf("加载脚本失败: %v", err)
		}
	}

	// 4. 列出所有脚本
	list := consumer.ListScripts()
	if len(list) != 3 {
		t.Errorf("期望 3 个脚本，得到 %d", len(list))
	}

	// 5. 验证脚本内容
	scriptMap := make(map[string]*Script)
	for _, s := range list {
		scriptMap[s.ID] = s
	}

	if scriptMap["SCRIPT_1"].Name != "脚本1" {
		t.Error("脚本1 名称不匹配")
	}

	if scriptMap["SCRIPT_2"].Name != "脚本2" {
		t.Error("脚本2 名称不匹配")
	}

	if scriptMap["SCRIPT_3"].Enabled {
		t.Error("脚本3 应该是禁用状态")
	}

	t.Logf("列出脚本测试成功，共 %d 个脚本", len(list))
}

// TestUnloadScript 测试卸载脚本
func TestUnloadScript(t *testing.T) {
	// 1. 创建存储层
	strg := storage.NewMemoryStorage()

	// 2. 创建脚本消费者
	config := DefaultScriptConfig()
	consumer := NewScriptConsumer("test-unload", strg, config)

	// 3. 加载脚本
	script := &Script{
		ID:      "UNLOAD_SCRIPT",
		Name:    "卸载测试脚本",
		Content: `return "test";`,
		Enabled: true,
	}

	if err := consumer.LoadScript(script); err != nil {
		t.Fatalf("加载脚本失败: %v", err)
	}

	// 4. 验证脚本已加载
	if _, exists := consumer.GetEngine().GetProgram("UNLOAD_SCRIPT"); !exists {
		t.Error("脚本应该已编译")
	}

	// 5. 卸载脚本
	if err := consumer.UnloadScript("UNLOAD_SCRIPT"); err != nil {
		t.Fatalf("卸载脚本失败: %v", err)
	}

	// 6. 验证脚本已删除
	status := consumer.GetScriptStatus("UNLOAD_SCRIPT")
	if status != nil {
		t.Error("脚本状态应该已删除")
	}

	if _, exists := consumer.GetEngine().GetProgram("UNLOAD_SCRIPT"); exists {
		t.Error("脚本编译缓存应该已删除")
	}

	// 7. 尝试执行已卸载的脚本
	_, err := consumer.ExecuteScript("UNLOAD_SCRIPT", nil)
	if err == nil {
		t.Error("执行已卸载的脚本应该返回错误")
	}

	t.Log("卸载脚本测试成功")
}
