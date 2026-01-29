# 脚本模块 (Script Module)

## 概述

脚本模块为 pans-runtime 提供了 JavaScript 脚本执行能力，基于 **Goja** 引擎实现。

### 核心特性

- ✅ **高性能**：VM 池化机制，复用 Goja VM 实例，避免频繁创建/销毁
- ✅ **简化架构**：继承 BaseConsumer，复用现有模式
- ✅ **安全可靠**：沙箱隔离、API 白名单、执行超时、资源限制
- ✅ **异步执行**：脚本不阻塞主流程，支持并发执行
- ✅ **易于使用**：JavaScript 语言、丰富的 API

### 性能指标

| 指标 | 目标值 |
|------|--------|
| 脚本执行耗时（VM 池命中） | < 10ms |
| VM 池复用率 | > 80% |
| 最大并发执行数 | 100 |
| 脚本执行超时 | 5s（可配置） |

---

## 快速开始

### 1. 配置 Go 代理（可选）

如果网络无法访问 `proxy.golang.org`，建议配置国内镜像：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### 2. 更新依赖

```bash
go mod tidy
```

### 3. 运行示例程序

```bash
go run examples/script/quick_start.go
```

### 4. 运行测试

```bash
go test -v ./internal/script/...
```

---

## 架构设计

### 核心组件

```
internal/script/
├── consumer.go          # ScriptConsumer（继承 BaseConsumer）
├── engine.go            # GojaEngine（脚本引擎核心）
├── vm_pool.go           # VM 池管理
├── sandbox.go           # 沙箱和 API 白名单
├── types.go             # 核心类型定义
├── config.go            # 配置结构（暂未实现）
├── trigger.go           # 触发器管理（Phase 2）
├── scheduler.go         # 周期执行调度器（Phase 2）
└── api/                 # JavaScript API 实现
    ├── variable.go      # 变量读写 API
    ├── system.go        # 系统指令 API（Phase 3）
    ├── communication.go # 通讯功能 API（Phase 4）
    ├── data.go          # 数据处理 API（Phase 5）
    └── ui.go            # 界面控制 API（Phase 6）
```

### 数据流程

```
脚本加载
  ↓
编译脚本（Goja）
  ↓
缓存编译结果
  ↓
执行脚本
  ↓
从 VM 池获取 VM
  ↓
设置沙箱环境（注入 API）
  ↓
执行脚本代码
  ↓
返回结果
  ↓
归还 VM 到池中
```

---

## 使用示例

### 基本用法

```go
package main

import (
    "context"
    "log"

    "pans-runtime/internal/script"
    "pans-runtime/internal/storage"
)

func main() {
    // 1. 创建存储层
    storage := memory.NewMemoryStorage()

    // 2. 创建脚本消费者
    config := script.DefaultScriptConfig()
    consumer := script.NewScriptConsumer("my-script", storage, config)

    // 3. 启动消费者
    ctx := context.Background()
    if err := consumer.Start(ctx); err != nil {
        log.Fatal(err)
    }
    defer consumer.Stop()

    // 4. 加载脚本
    script := &script.Script{
        ID:      "MY_SCRIPT",
        Name:    "我的脚本",
        Content: `
            var a = 10;
            var b = 20;
            return a + b;
        `,
        Enabled: true,
    }

    if err := consumer.LoadScript(script); err != nil {
        log.Fatal(err)
    }

    // 5. 执行脚本（同步）
    result, err := consumer.ExecuteScript("MY_SCRIPT", nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("执行结果: %v", result)
}
```

### JavaScript API

#### 变量操作 API

```javascript
// 读取单个变量
var value = Variable.read(100001);

// 批量读取
var values = Variable.readBatch([100001, 100002, 100003]);

// 写入变量
Variable.write(100001, 26.5);

// 批量写入
Variable.writeBatch({
    100001: 26.5,
    100002: true
});
```

#### 系统指令 API

```javascript
// 日志记录
Log.info("信息日志");
Log.warn("警告日志");
Log.error("错误日志");

// 文件操作
var content = System.File.readText("/data/config.txt");
System.File.writeText("/data/output.txt", "Hello World");

// JSON 处理
var obj = JSON.parse('{"name": "test"}');
var str = JSON.stringify({name: "test"});
```

---

## 配置说明

### ScriptConfig

```go
type ScriptConfig struct {
    // VM 池配置
    VMPoolSize    int           // VM 池大小（默认 10）
    VMMaxIdle     time.Duration // VM 最大空闲时间（默认 5m）
    VMMaxLifetime time.Duration // VM 最大生命周期（默认 30m）

    // 执行配置
    DefaultTimeout time.Duration // 默认执行超时（默认 5s）
    MaxConcurrent  int           // 最大并发执行数（默认 100）
    QueueSize      int           // 执行队列大小（默认 1000）

    // 安全配置
    EnableSandbox bool  // 是否启用沙箱（默认 true）
    MemoryLimit   int64 // 内存限制（字节，默认 10MB）
    MaxExecutions int   // 单脚本最大执行次数/分钟（默认 60）

    // 日志配置
    EnableLog bool   // 是否启用脚本日志（默认 false）
    LogPath   string // 日志路径（默认 "./data/logs/script/"）
}
```

### 默认配置

```go
config := script.DefaultScriptConfig()
// 可以覆盖默认配置
config.EnableLog = true
config.VMPoolSize = 20
```

---

## 测试

### 运行所有测试

```bash
go test -v ./internal/script/...
```

### 运行特定测试

```bash
# 测试脚本加载和执行
go test -v ./internal/script/ -run TestScriptLoadAndExecute

# 测试 VM 池
go test -v ./internal/script/ -run TestVMPool

# 测试异步执行
go test -v ./internal/script/ -run TestScriptAsyncExecution
```

### 基准测试

```bash
go test -bench=. -benchmem ./internal/script/
```

---

## Phase 规划

### ✅ Phase 1: 基础脚本引擎（已完成）

- 集成 Goja 引擎
- 实现 ScriptConsumer
- 实现 VM 池管理
- 实现基础沙箱
- 变量读写 API

### ⏳ Phase 2: 触发器系统（待实施）

- 变量变化触发
- 周期执行触发
- 系统触发（开机/关机）
- 报警触发

### ⏳ Phase 3: 系统指令 API（待实施）

- 文件操作 API
- 日志记录 API
- 时间日期 API
- 字符串处理 API
- JSON 处理 API

### ⏳ Phase 4: 通讯功能 API（待实施）

- HTTP 请求 API
- MQTT 发布/订阅 API
- Modbus 读写 API（预留接口）

### ⏳ Phase 5: 数据处理 API（待实施）

- 配方管理 API
- 数据记录 API
- 数据统计 API

### ⏳ Phase 6: 界面控制 API（待实施）

- 前端交互接口设计
- WebSocket 通信机制
- 页面跳转 API
- 元件属性修改 API

---

## 安全机制

### 1. 沙箱隔离

- API 白名单：只暴露安全的 API
- 移除危险对象：`eval`, `require`, `Function`
- 全局变量隔离：每个脚本独立的 VM 实例

### 2. 资源限制

- 执行超时：默认 5s，可配置
- 内存限制：每个 VM 限制 10MB
- 并发限制：最大 100 个并发执行
- 频率限制：单脚本最多 60 次/分钟

### 3. 错误处理

- 脚本异常不影响系统：捕获所有 panic
- 详细错误日志：记录脚本 ID、错误信息、堆栈

---

## 性能优化

### VM 池化

- 复用 Goja VM 实例，避免频繁创建/销毁
- VM 创建耗时：~5ms
- VM 复用耗时：~0.01ms
- 目标复用率：> 80%

### 异步执行

- 脚本执行不阻塞主流程
- 支持并发执行（最大 100 个并发）
- 执行队列缓冲（1000 个任务）

### 编译缓存

- 脚本编译结果缓存
- 避免重复编译
- 提升执行性能

---

## 故障排查

### 问题1：编译失败

**错误**：`missing go.sum entry for package github.com/dop251/goja`

**解决**：
```bash
# 配置 Go 代理
go env -w GOPROXY=https://goproxy.cn,direct

# 更新依赖
go mod tidy
```

### 问题2：脚本执行超时

**错误**：`脚本执行超时`

**解决**：
- 检查脚本是否有死循环
- 增加超时时间配置：`config.DefaultTimeout = 10 * time.Second`

### 问题3：VM 池复用率低

**原因**：并发执行过多，VM 池大小不足

**解决**：
- 增加 VM 池大小：`config.VMPoolSize = 20`
- 减少并发执行数：`config.MaxConcurrent = 50`

---

## 后续优化

### 短期优化

1. **调试工具**：
   - Web 调试界面
   - 脚本断点设置
   - 变量监视窗口

2. **脚本编辑器**：
   - 语法高亮
   - 自动补全
   - 错误提示

3. **性能监控**：
   - 脚本执行统计
   - VM 池使用率监控
   - 资源占用报告

### 中期优化

1. **脚本市场**：
   - 预置常用脚本库
   - 脚本分享机制
   - 脚本版本管理

2. **分布式执行**：
   - 脚本远程执行
   - 集群调度
   - 负载均衡

---

## 相关文档

- [脚本模块完整开发方案](../../.claude/plans/radiant-tickling-emson.md)
- [项目整体架构](../../CLAUDE.md)
- [实现计划](../../IMPLEMENTATION_PLAN.md)

---

## 许可证

MIT License
