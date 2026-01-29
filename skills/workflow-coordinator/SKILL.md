---
name: workflow-coordinator
description: |
  Workflow state management and context recovery skill.

  **When to use**:
  - Managing task stack (push/pop/pause/resume)
  - Validating phase transition prerequisites
  - Providing context recovery for interrupted work
  - Switching between tasks
  - Merging parallel work streams
  - Querying current state and progress

  **Use cases**:
  - "暂停当前任务，处理 IMP-002" → Pushes current task, switches
  - "继续执行之前的任务" → Pops from stack, restores context
  - "验证是否可以进入实现阶段" → Validates phase prerequisites
  - "我上次做到哪了？" → Reads current-phase.md and reports
  - "显示所有并行任务" → Shows active-tasks.md
  - "检查依赖状态" → Queries dependency-backlog.md

  **Key responsibilities**:
  - Maintain task stack for pause/resume
  - Track all parallel tasks
  - Validate phase transitions
  - Detect dependency completions
  - Provide context recovery
  - Manage state persistence

  **Files managed**:
  - .claude-workflow/current-phase.md (current task and stack)
  - .claude-workflow/parallel-tasks/active-tasks.md (parallel tasks)
  - .claude-workflow/parallel-tasks/task-dependencies.md (dependencies)
  - .claude-workflow/dependency-backlog.md (pending completions)

  **Integration**: Works with all other skills to maintain workflow state
---

# Workflow Coordinator

## Overview

This skill manages the overall workflow state, handles task switching, maintains context, and ensures smooth parallel task execution.

## Core Responsibilities

### 1. Task Stack Management

#### Task Stack Structure

```
[0] Current task (top of stack)
    ├─ Current step
    └─ Next steps

[1] Paused task
    └─ Suspended at step X

[2] Waiting task
    └─ Not started

[N] ...
```

#### Push Operation (Pause Current Task)

```
User: "暂停当前任务，处理 IMP-002"

→ workflow-coordinator activates
→ Reads current-phase.md
→ Pushes current task to stack
→ Updates task-stack.md
→ Switches to IMP-002
→ Updates current-phase.md with new context
```

**Stack update**:
```markdown
## 堆栈结构
[0] IMP-002 (当前任务)
    └─ 全部任务

[1] IMP-001 (已暂停)
    └─ 4.1.2 [暂停中] 实现登录逻辑
```

#### Pop Operation (Resume Previous Task)

```
User: "继续执行之前的任务"

→ workflow-coordinator activates
→ Reads task-stack.md
→ Pops from stack
→ Restores context for previous task
→ Updates current-phase.md
→ Loads relevant files
→ Reports: "恢复到 IMP-001，步骤 4.1.2"
```

### 2. Phase Transition Validation

#### Validation Process

Before allowing phase transition, validate:

**Feature-level validation** (per FE-{N}):
```markdown
### FE-{N} → Design Stage Validation
检查项:
✅ REQ-{N} 已创建
✅ FE-{N} 已创建并关联到 REQ-{N}
✅ FE-{N} 已分配到具体平台
✅ FE-{N} 的功能描述完整

结果: ✅ 允许进入设计阶段
```

#### Validation Commands

```
User: "验证 FE-001 是否可以进入实现阶段"

→ workflow-coordinator activates
→ Reads SOL-001 document
→ Checks validation checklist:
  ✅ FE-001 已完成需求分解
  ✅ SOL-001 技术方案文档已创建
  ✅ SOL-001 接口设计已明确
  ✅ SOL-001 方案已通过评审

→ Reports: "✅ FE-001 可以进入实现阶段"
```

### 3. Context Recovery

#### Recovery Process

```
User: "我上次做到哪了？"

→ workflow-coordinator activates
→ Reads current-phase.md
→ Reads task-stack.md
→ Reads active-tasks.md

→ Reports:
================================
📍 当前任务状态
================================
项目: 物联网平台
当前任务: IMP-001 - 用户认证实现
当前步骤: 4.1.2 - 实现登录逻辑
完成度: 60%
状态: ⏳ 进行中

📋 任务链
1. ✅ REQ-001 - 原始需求
2. ✅ FE-001 - 功能需求
3. ✅ SOL-001 - 技术方案
4. ⏳ IMP-001 - 实现计划
   - 4.1.1 ✅ 创建项目结构
   - 4.1.2 ⏳ 实现登录逻辑 (60%)
   - 4.1.3 📋 实现验证逻辑

🔄 并行任务
- FE-002: 方案设计中 (40%)
- FE-003: 需求分解中 (20%)

⚠️ 依赖忽略
- D-001: IMP-002 依赖 IMP-001 (等待中)

快速恢复: "继续执行 IMP-001 的步骤 4.1.2"
================================
```

#### Context Snapshot

When pausing task, save context:

```markdown
## 上下文快照
- **任务ID**: IMP-001
- **任务名称**: 用户认证实现
- **当前步骤**: 4.1.2 - 实现登录逻辑
- **完成度**: 60%
- **相关文件**: src/auth/login.ts, src/auth/validator.ts
- **依赖项**: 无
- **下一步**: 完成 authenticate() 方法
- **暂停时间**: 2025-01-19 14:30
- **暂停原因**: 切换到 IMP-002
```

### 4. Parallel Task Management

#### Active Tasks Tracking

Maintain `active-tasks.md`:
```markdown
## 正在进行的功能

### FE-001: 用户认证 (网关、云平台)
- **当前阶段**: IMP-001 实现中
- **完成度**: 60%
- **当前步骤**: 实现登录逻辑 (4.1.2)
- **状态**: 🟢 正常进行
- **依赖**: 无

### FE-002: 数据采集 (网关、云端)
- **当前阶段**: SOL-002 方案设计中
- **完成度**: 40%
- **当前步骤**: 设计通信接口
- **状态**: 🟢 正常进行
- **依赖**: FE-001 (弱依赖，已忽略)
```

#### Updates

**When to update active-tasks.md**:
- Starting new feature
- Completing a phase
- Updating progress
- Changing status
- Adding/removing dependencies

### 5. Dependency Management

#### Dependency Detection

```
User: "FE-002 依赖 FE-001，如何处理？"

→ workflow-coordinator activates
→ Reads task-dependencies.md
→ Analyzes dependency type

→ Reports:
================================
依赖分析: FE-002 → FE-001
类型: 弱依赖 (可以使用 mock)

建议策略: ✅ 并行开发 + 依赖忽略
1. FE-001 继续实现
2. FE-002 使用 mock 继续开发
3. FE-001 完成后，FE-002 补齐真实调用

预计影响: 低 (mock 足够用于开发)
补齐时间: 约 1 小时
================================
```

#### Dependency Completion Detection

When IMP-{N} completes:

1. **Scan dependency-backlog.md**
2. **Find all entries** depending on IMP-{N}
3. **Create notification**:
   ```
   ⚠️ 依赖完成通知
   依赖模块: IMP-001 (用户认证)
   完成时间: 2025-01-19 16:00

   影响范围:
   - D-001: IMP-002 需要补齐认证验证逻辑
   - D-002: IMP-003 需要补齐 Token 生成逻辑

   建议:
   1. 按优先级 P0 → P1 → P2 依次补齐
   2. 每个补齐后进行集成测试
   3. 更新代码映射表
   ```

4. **Update dependency-backlog.md** statuses

### 6. State Persistence

#### Files Managed

1. **current-phase.md**: Current task and stack
2. **active-tasks.md**: All parallel tasks
3. **task-dependencies.md**: Dependency graph
4. **dependency-backlog.md**: Pending completions

#### Update Frequency

**Update after**:
- Every task status change
- Phase transitions
- Dependency changes
- Progress updates
- Task switches

#### Backup Strategy

Keep history:
```
.claude-workflow/state-history/
├── 2025-01-19-current-phase.md.bak
├── 2025-01-18-current-phase.md.bak
└── ...
```

## Commands and Patterns

### Status Queries

```
"显示当前状态" → Shows current task, progress, parallel tasks
"显示所有任务" → Lists all active and pending tasks
"显示任务堆栈" → Shows task stack structure
"检查依赖" → Shows dependency status and backlog
"显示进度" → Shows overall project progress
```

### Task Management

```
"暂停当前任务" → Push to stack, save context
"继续上一个任务" → Pop from stack, restore context
"切换到 IMP-XXX" → Switch to specific task
"完成任务" → Mark complete, update status
```

### Validation

```
"验证 FE-XXX 状态" → Validate feature readiness
"检查是否可以进入设计阶段" → Validate design stage entry
"验证所有前置条件" → Validate all prerequisites
```

## State File Formats

### current-phase.md

```markdown
# 当前执行阶段

## 项目信息
- **项目名称**: {项目名称}
- **当前阶段**: 阶段{N} - {阶段名称}
- **开始日期**: {YYYY-MM-DD}
- **整体进度**: {percentage}%

## 当前任务
- **任务ID**: {ID}
- **任务名称**: {名称}
- **状态**: 进行中 / 暂停 / 等待
- **当前步骤**: {步骤描述}
- **完成度**: {percentage}%

## 任务堆栈
### 主任务链
1. ✅ {完成项}
2. ⏳ {进行项}
   - {子项1}
   - {子项2}

### 暂停的任务
- [暂停] {任务ID} - 在步骤{X}暂停

### 待执行任务
- [待执行] {任务ID}

## 上下文信息
- **相关平台**: {平台列表}
- **相关文件**: {文件列表}
- **依赖项**: {依赖}
- **阻塞项**: {阻塞}
```

### active-tasks.md

```markdown
# 活跃并行任务列表

更新时间: {YYYY-MM-DD HH:MM}

## 正在进行的功能
{任务列表}

## 并行统计
- **活跃功能数**: {数量}
- **分布阶段**: {阶段分布}
- **预计完成日期**: {YYYY-MM-DD}

## 阻塞警告
{阻塞信息}
```

### task-dependencies.md

```markdown
# 任务依赖关系图

## 依赖关系可视化
{依赖图}

## 依赖类型定义
{类型表}

## 依赖忽略记录
{忽略记录}
```

## Error Recovery

### State Corruption

If state files are corrupted:

1. **Detect inconsistency**
2. **Report specific issue**
3. **Offer recovery**:
   ```
   ⚠️ 检测到状态不一致

   问题: current-phase.md 和 active-tasks.md 不匹配
   建议:
   1. 使用备份恢复
   2. 重建状态文件
   3. 从日志重建
   ```
4. **Restore from backup** if available

### Lost Context Recovery

If context is lost:

1. **Check state files**
2. **Check implementation logs**
3. **Check git history**
4. **Reconstruct from available data**
5. **Ask user for confirmation**

## Integration with Other Skills

### workflow-orchestrator
- Validates phase transitions
- Activates on orchestrator request
- Reports state changes

### requirement-manager
- Updates active-tasks.md when FE created
- Updates task-dependencies.md

### solution-designer
- Updates active-tasks.md when SOL created
- Tracks design progress

### implementation-manager
- Updates current-phase.md during implementation
- Manages task stack for switches
- Tracks dependency completions

## Best Practices

1. **Always update state files**: After every action
2. **Save context before switch**: Complete context snapshots
3. **Validate before transition**: Check prerequisites
4. **Monitor parallel tasks**: Regular status checks
5. **Clean up completed tasks**: Archive finished tasks
6. **Backup regularly**: Keep state history

## Usage Examples

### Example 1: Task Switch

```
User: "暂停当前任务，处理 IMP-002"

→ workflow-coordinator activates
→ Reads current-phase.md: Current = IMP-001
→ Pushes IMP-001 to stack [1]
→ Updates task-stack.md
→ Switches to IMP-002
→ Updates current-phase.md:
  Current task: IMP-002
  Stack: [0]=IMP-002, [1]=IMP-001
→ Reports: "已切换到 IMP-002，IMP-001 已暂停"
```

### Example 2: Context Recovery

```
User: "继续之前的任务"

→ workflow-coordinator activates
→ Reads task-stack.md
→ Pops IMP-001 from [1]
→ Restores context:
  - Task: IMP-001
  - Step: 4.1.2
  - Progress: 60%
  - Files: src/auth/login.ts
→ Updates current-phase.md
→ Loads context files
→ Reports: "恢复到 IMP-001，步骤 4.1.2"
```

### Example 3: Progress Report

```
User: "显示所有并行任务"

→ workflow-coordinator activates
→ Reads active-tasks.md
→ Generates report:

================================
并行任务概览
================================
活跃任务数: 3

✅ FE-001: 用户认证 (60%)
   阶段: 实现
   平台: 网关, 云平台

⏳ FE-002: 数据采集 (40%)
   阶段: 方案设计
   平台: 网关, 云端
   依赖: FE-001 (弱依赖, 已忽略)

📋 FE-003: 实时监控 (20%)
   阶段: 需求分解
   平台: HMI, Web编辑器

整体进度: 40%
预计完成: 2025-01-25
================================
```

### Example 4: Dependency Notification

```
(IMP-001 completes)

→ workflow-coordinator detects completion
→ Scans dependency-backlog.md
→ Finds D-001, D-002 depending on IMP-001
→ Creates notification:
  "⚠️ IMP-001 已完成!
   需要补齐:
   - D-001: IMP-002 认证验证
   - D-002: IMP-003 Token 生成"

→ Updates active-tasks.md
→ Updates dependency-backlog.md
→ Notifies implementation-manager
```
