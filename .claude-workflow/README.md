# Claude Code 工作流管理

本目录使用 Claude Code 结构化研发工作流管理体系进行管理。

## 目录结构

```
.claude-workflow/
├── 01-requirements/        # 需求管理
├── 02-planning/           # 计划管理
├── 03-design/             # 设计管理
├── 04-implementation/      # 实现管理
├── 05-verification/       # 验证管理
├── 06-documentation/      # 文档管理
├── 07-platforms/          # 多平台管理
├── parallel-tasks/        # 并行任务管理
├── templates/             # 文档模板
├── state-history/         # 状态历史
├── config.yml             # 项目配置
├── current-phase.md       # 当前阶段
├── feature-to-code-map.md # 功能代码映射
├── dependency-backlog.md  # 依赖跟踪
└── rt-matrix.md           # 需求追溯矩阵
```

## 使用方式

### 开始新功能
```
"我们需要添加[新功能]"
→ workflow-orchestrator 自动激活
→ 跟随引导完成需求收集、方案设计、实现...
```

### 查看当前状态
```
"我上次做到哪了？"
→ workflow-coordinator 显示当前任务和进度
```

### 暂停和恢复
```
"暂停当前任务，处理其他任务"
→ 任务切换和上下文保存
"继续之前的任务"
→ 恢复到之前的任务
```

## 工作流程

1. **需求阶段** (requirement-manager)
   - 创建原始需求 (REQ-{N})
   - 分解功能需求 (FE-{N})

2. **设计阶段** (solution-designer)
   - 创建技术方案 (SOL-{N})
   - 设计架构和接口

3. **实现计划** (implementation-manager)
   - 创建实现计划 (IMP-{N})
   - 设计代码映射

4. **代码实现** (implementation-manager)
   - 按计划实现
   - 记录代码映射

5. **验证和文档** (verification-manager, documentation-generator)
   - 测试验证
   - 生成文档

## 更多信息

参考完整文档: `Claude-Code-Workflow-Design-v1.0.0.md`
