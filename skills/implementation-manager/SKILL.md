---
name: implementation-manager
description: |
  Implementation planning and execution management skill with precision modification support.

  **PREREQUISITE**: Can ONLY be used after solution-designer phase is complete.
  workflow-orchestrator will validate prerequisites before activation.

  **When to use**:
  - Creating detailed implementation plans (IMP-XXX)
  - Executing implementation tasks
  - Recording feature-to-code mappings
  - Updating traceability matrix
  - Managing implementation logs
  - Modifying existing functionality (precision modification)

  **Use cases**:
  - "制定用户认证的实现计划" → Creates IMP-001
  - "实现登录功能" → Executes implementation and records mapping
  - "修改 FE-001-02 的验证逻辑" → Uses mapping to locate affected code
  - "显示 FE-001 的代码映射" → Queries feature-to-code-map.md

  **Workflow stage**: Stage 3 (Planning) & Stage 4 (Implementation)
  **Input**: Technical solutions (SOL-{N}) from solution-designer
  **Outputs**:
  - Implementation plans (04-implementation/implementation-plans/IMP-{N}.md)
  - Feature-to-code mappings (04-implementation/code-mapping/feature-to-code-map.md)
  - Impact analysis reports (04-implementation/code-mapping/impact-analysis.md)
  - Implementation logs (04-implementation/implementation-logs/YYYY-MM-DD-IMP-{N}.md)
  - Updated traceability matrix

  **Key features**:
  - Detailed implementation planning with task breakdown
  - Real-time code mapping during implementation
  - Precision modification using feature-to-code map
  - Dependency ignore mechanism support
  - Complete implementation logging
  - Impact analysis for changes
---

# Implementation Manager

## Overview

This skill manages both implementation planning and actual code implementation, with a strong emphasis on maintaining complete traceability between features and code.

## Prerequisites

### Must Have (Validated by workflow-orchestrator)
- ✅ SOL-{N} solution document exists
- ✅ SOL-{N} has complete technical design
- ✅ SOL-{N} has API/interface specifications
- ✅ SOL-{N} has been reviewed and approved

### Input From solution-designer
- Technical solution document (SOL-{N})
- Architecture diagrams
- API specifications
- Data models
- Implementation plan outline

## Implementation Planning (Stage 3)

### Step 1: Create Implementation Plan (IMP-{N})

1. **Copy template**:
   ```
   .claude-workflow/templates/IMP-template.md
   ```

2. **Fill in sections**:
   - Overview
   - Technology stack
   - Task breakdown (phases)
   - Technical implementation details
   - Code mapping table (initial)
   - Dependencies
   - Test strategy

3. **Design code mapping table**:
   ```markdown
   | 功能ID | 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
   ```

4. **Save to**:
   ```
   .claude-workflow/04-implementation/implementation-plans/IMP-{N}-{title}.md
   ```

### Step 2: Break Down Tasks

**Task granularity**:
- Each task should be completable in 2-4 hours
- Tasks should be independent where possible
- Define clear acceptance criteria for each task

**Example**:
```markdown
### Phase 1: 基础结构
- [ ] Task 1.1: 创建项目目录结构 - 预估: 1h
- [ ] Task 1.2: 创建基础类和接口 - 预估: 2h
- [ ] Task 1.3: 配置依赖注入 - 预估: 1h

### Phase 2: 核心功能
- [ ] Task 2.1: 实现登录逻辑 - 预估: 4h
- [ ] Task 2.2: 实现验证逻辑 - 预估: 3h
- [ ] Task 2.3: 实现 token 生成 - 预估: 2h
```

### Step 3: Identify Dependencies

**Types of dependencies**:
- **Implementation dependencies**: Other IMPs that must complete first
- **Functional dependencies**: Other FE-{N} this feature depends on
- **External dependencies**: Libraries, services, APIs

**Document in IMP-{N}**:
```markdown
## 依赖项
### 前置依赖
- IMP-{N}: {dependency description}
- FE-{N}: {dependency description}

### 外部依赖
- {Library}: {version} - {purpose}
```

### Step 4: Handle Dependencies with Dependency Ignore

If dependency can be mocked:

1. **Create dependency ignore record** in code:
   ```typescript
   // TODO(依赖): FE-001 - 用户认证
   // 说明: 此处需要调用认证服务验证用户身份
   // 当前状态: 使用 mock 实现
   // 依赖模块: IMP-001
   // 补齐优先级: P0
   // 预计补齐日期: 2025-01-20
   // 负责人: @username
   ```

2. **Update dependency-backlog.md**:
   ```markdown
   | D-001 | IMP-002 | IMP-001 | 认证验证 | src/data/collector.ts:45 | P0 | ⏳ | 2025-01-20 |
   ```

3. **Implement with mock**:
   ```typescript
   // === 依赖忽略开始 ===
   // Mock 实现:
   const user = { id: userId, name: 'Mock User' };
   // === 依赖忽略结束 ===
   ```

4. **Continue with main logic**

## Code Implementation (Stage 4)

### Step 1: Implement According to Plan

Follow IMP-{N} task breakdown:
1. Read task description
2. Implement feature
3. Write/update code mapping
4. Run tests
5. Update progress

### Step 2: Record Code Mapping

**Real-time mapping update**:
As you implement code, immediately update the mapping table:

```markdown
## 代码实现映射表
| 功能ID | 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
|--------|--------|---------|---------|----------|---------|------|
| FE-001-01 | 用户登录 | src/auth/login.ts | LoginService | authenticate() | 45-89 | ✅ |
| FE-001-02 | 密码验证 | src/auth/validator.ts | PasswordValidator | validate() | 12-34 | ✅ |
| FE-001-03 | Token生成 | src/auth/token.ts | TokenService | generate() | 23-56 | ⏳ |
```

**Update feature-to-code-map.md**:
```markdown
.claude-workflow/04-implementation/code-mapping/feature-to-code-map.md
```

### Step 3: Implementation Logging

Create daily log:
```
.claude-workflow/04-implementation/implementation-logs/YYYY-MM-DD-IMP-{N}.md
```

**Log entries**:
- ✅ Completed tasks
- ⏳ In-progress tasks
- ❌ Blocked tasks
- 📝 Notes and observations

**Example**:
```markdown
## 2025-01-19

### 完成任务
- ✅ Task 1.1: 创建项目目录结构
- ✅ Task 1.2: 创建基础类和接口
- ✅ Task 2.1: 实现登录逻辑 (部分)

### 进行中
- ⏳ Task 2.1: 实现登录逻辑 (60%)
  - 已完成: 基础认证流程
  - 待完成: 错误处理

### 遇到问题
- ❌ Task 1.3: 配置依赖注入
  - 问题: DI库版本冲突
  - 解决: 升级到最新版本

### 下一步
- 完成 Task 2.1
- 开始 Task 2.2
```

## Precision Modification

### Querying Feature-to-Code Map

When user wants to modify existing functionality:

```
User: "修改 FE-001-02 的密码验证逻辑"

→ implementation-manager activates
→ Queries feature-to-code-map.md
→ Locates: FE-001-02 → src/auth/validator.ts:12-34
→ Analyzes dependencies
→ Displays impact analysis
```

### Impact Analysis Report

**Display**:
```markdown
## 修改影响分析

**功能**: FE-001-02 - 密码验证
**位置**: src/auth/validator.ts:12-34

### 直接影响
- FE-001-01: 用户登录 (调用方)
- 测试: TC-001-002 (密码验证测试)

### 需要修改的文件
1. src/auth/validator.ts (主要)
2. src/auth/login.ts (调用方，可能需要调整)
3. tests/auth/validator.test.ts (测试)

### 风险评估
- 风险等级: 低
- 影响用户: 是 (所有登录用户)
- 回滚难度: 低

### 建议
1. 创建 feature branch
2. 实现新验证逻辑
3. 更新测试
4. 运行完整测试套件
5. 代码审查
```

### Executing Modification

1. **Backup**: Create branch
2. **Modify**: Update code at specified location
3. **Test**: Run relevant tests
4. **Update**: Update code mapping
5. **Document**: Update rt-matrix.md

## Dependency Completion and补齐

### When Dependency Completes

1. **Detection**: workflow-coordinator detects IMP-001 completed
2. **Notification**: Notifies waiting implementations
3. **补齐**: For each dependency:

   ```markdown
   ## 补齐 D-001: IMP-002 依赖 IMP-001

   1. ✅ 确认 IMP-001 已完成
   2. ⏳ 移除 src/data/collector.ts:45 的 mock
   3. ⏳ 替换为真实调用: this.authService.validate()
   4. ⏳ 更新错误处理
   5. ⏳ 运行集成测试
   6. ⏳ 更新 feature-to-code-map.md
   7. ⏳ 更新 dependency-backlog.md (标记完成)
   ```

4. **Verify**: Run integration tests
5. **Update**: Mark as complete in dependency-backlog.md

## Validation Checklist

Before marking implementation complete:

**IMP-{N} validation**:
- [ ] All tasks completed
- [ ] Code implemented according to SOL-{N}
- [ ] Code mapping table complete and accurate
- [ ] All dependencies handled (or mocked)
- [ ] Unit tests written
- [ ] Integration tests defined
- [ ] Implementation log up to date
- [ ] No uncompleted dependency ignores (unless planned)

**Code quality**:
- [ ] Code follows project standards
- [ ] Sufficient comments
- [ ] Error handling
- [ ] Logging added
- [ ] Performance considered

**Traceability**:
- [ ] feature-to-code-map.md updated
- [ ] rt-matrix.md updated
- [ ] dependency-backlog.md current

## Workflow Integration

### Entry Point
Activated by `workflow-orchestrator` for Stage 3 & 4.

### Exit Criteria
Transition to verification when:
1. All implementation tasks complete
2. Code mapping complete
3. Unit tests pass
4. Traceability updated
5. No critical blockers

### Handoff to Verification Stage
Provide to `verification-manager`:
- IMP-{N} document
- Code locations
- Test strategy
- Dependencies

## Usage Examples

### Example 1: Implementation Planning

```
User: "制定用户认证的实现计划"

→ implementation-manager activates
→ Reads SOL-001: 用户认证技术方案
→ Creates IMP-001: 用户认证实现

Task breakdown:
- Phase 1: 基础结构 (4h)
- Phase 2: 认证逻辑 (9h)
- Phase 3: 测试 (4h)

Code mapping designed:
- FE-001-01 → src/auth/login.ts
- FE-001-02 → src/auth/validator.ts
- FE-001-03 → src/auth/token.ts

Dependencies: None

→ Saves IMP-001
→ Notifies workflow-orchestrator
```

### Example 2: Code Implementation

```
User: "实现 FE-001-01 用户登录功能"

→ implementation-manager activates
→ Reads IMP-001
→ Locates task: Phase 2, Task 2.1
→ Creates src/auth/login.ts
→ Implements LoginService.authenticate()
→ Records mapping: FE-001-01 → LoginService.authenticate() :45-89
→ Updates feature-to-code-map.md
→ Creates unit test
→ Runs test: ✅ Pass
→ Updates implementation log

→ Progress: 60% complete
```

### Example 3: Precision Modification

```
User: "修改 FE-001-02，增加密码复杂度检查"

→ implementation-manager activates
→ Queries feature-to-code-map.md
→ Locates: src/auth/validator.ts:12-34

Impact analysis:
主要文件: src/auth/validator.ts
影响功能: FE-001-01 (登录)
影响测试: TC-001-002
风险等级: 低

→ Executes modification:
1. Reads src/auth/validator.ts
2. Adds complexity check:
   ```typescript
   function validatePassword(password: string): boolean {
     if (password.length < 8) return false;
     if (!/[A-Z]/.test(password)) return false;
     if (!/[a-z]/.test(password)) return false;
     if (!/[0-9]/.test(password)) return false;
     return true;
   }
   ```
3. Updates FE-001-02 mapping (line numbers changed)
4. Updates test case
5. Runs tests: ✅ All pass
6. Updates feature-to-code-map.md
7. Updates rt-matrix.md

→ Modification complete
```

### Example 4: Dependency Ignore

```
User: "FE-002 依赖 FE-001，但想并行开发"

→ implementation-manager activates
→ Analyzes dependency: FE-002 → FE-001 (弱依赖)
→ Suggests: Use dependency ignore mechanism

→ Creates dependency ignore in code:
   // TODO(依赖): FE-001 - 用户认证
   // Mock 实现
   const user = mockUser;

→ Updates dependency-backlog.md:
   D-001 | IMP-002 | IMP-001 | 认证验证 | src/data/collector.ts:45 | P0

→ Implements main logic with mock
→ Updates active-tasks.md

→ Later, when FE-001 completes:
   workflow-coordinator: "IMP-001 completed!"
   → Replaces mock with real call
   → Runs integration tests
   → Updates dependency-backlog.md (complete)
```

## Code Mapping Best Practices

### Granularity

**Function-level mapping** (recommended):
```markdown
| FE-001-01 | 登录 | login.ts | LoginService | authenticate() | 45-89 | ✅ |
```

**When to use class-level**:
- Simple utility classes
- Single-purpose classes

**When to use file-level**:
- Configuration files
- Static data files

### Mapping Updates

**When to update**:
- ✅ After implementing each function
- ✅ After refactoring
- ✅ After modifying existing code
- ✅ After resolving dependencies

**How to update**:
1. Locate function in code
2. Identify line numbers
3. Update mapping table
4. Update status (⏳ → ✅)

## Implementation Metrics

Track in IMP-{N}:

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| Tasks completed | {total} | {current} | {percentage}% |
| Code coverage | 80% | {actual}% | ✅/❌ |
| Tests passing | 100% | {actual}% | ✅/❌ |
| Dependencies resolved | {total} | {current} | {percentage}% |

## Templates Location

```
.claude-workflow/templates/
└── IMP-template.md
```

## Output Files

```
.claude-workflow/04-implementation/
├── implementation-plans/
│   └── IMP-{N}-{title}.md
├── code-mapping/
│   ├── feature-to-code-map.md
│   └── impact-analysis.md
└── implementation-logs/
    └── YYYY-MM-DD-IMP-{N}.md
```
