# Claude Code Workflow 使用指南

本指南将帮助你快速上手 Claude Code 结构化研发工作流管理体系。

## 🎯 核心理念

传统开发流程的问题：
- ❌ 需求、设计、实现脱节
- ❌ 无法追溯代码到需求的映射
- ❌ 并行开发时依赖阻塞
- ❌ 修改代码时影响范围不明确

本工作流的解决方案：
- ✅ 每个功能独立推进，不阻塞其他功能
- ✅ 完整的需求→代码映射和追溯
- ✅ 依赖忽略机制支持并行开发
- ✅ 精准的影响分析和代码定位

## 🚀 5分钟快速上手

### Step 1: 初始化项目

```bash
# 克隆模板到新项目
cp -r claude-workflow-template my-project
cd my-project

# 运行初始化脚本
python .claude-workflow/scripts/init-workflow.py \
    --project-name "物联网平台" \
    --platforms "gateway,cloud,hmi"
```

### Step 2: 配置项目

编辑 `.claude-workflow/config.yml`，启用你需要的平台。

### Step 3: 安装 Skills

```bash
# 复制 skills 到 Claude Code 目录
cp -r skills/* ~/.config/claude/skills/
```

### Step 4: 启动 Claude Code

```bash
claude-code .
```

### Step 5: 开始第一个需求

```
"我们需要添加用户登录功能"
```

Claude Code 会自动激活 `workflow-orchestrator`，引导你完成整个流程。

## 📖 工作流程详解

### 阶段 1: 需求收集

```
User: "我们需要添加用户认证功能"
```

**Claude Code 会做什么**:

1. **激活 requirement-manager**
2. **询问问题**:
   - 这个功能需要在哪些平台实现？
   - 有哪些具体的认证要求？
   - 优先级是什么？
3. **创建文档**:
   - `REQ-001-用户认证.md`
   - `FE-001-用户登录.md` (Gateway, Cloud)
   - `FE-002-权限管理.md` (All platforms)
4. **更新**:
   - `requirement-mapping.md`
   - `active-tasks.md`

**输出**: 需求文档完整，可以进入下一阶段

### 阶段 2: 方案设计

```
User: "设计用户认证的技术方案"
```

**Claude Code 会做什么**:

1. **激活 solution-designer**
2. **引导设计**:
   - 选择技术方案 (JWT vs Session?)
   - 设计 API 接口
   - 设计数据库表结构
   - 架构设计
3. **创建文档**:
   - `SOL-001-用户认证技术方案.md`
   - `system-architecture.md` (更新)
   - `api-specifications.md` (更新)
   - `schema.md` (数据库设计)
4. **创建 ADR**:
   - `adr-001-选择JWT作为认证方式.md`

**输出**: 技术方案完整，可以进入下一阶段

### 阶段 3: 实现计划

```
User: "制定用户认证的实现计划"
```

**Claude Code 会做什么**:

1. **激活 implementation-manager**
2. **制定计划**:
   - 任务分解 (Phase 1, 2, 3)
   - 时间估算
   - 依赖分析
3. **创建文档**:
   - `IMP-001-用户认证实现.md`
   - 代码映射表设计
4. **更新**:
   - `current-phase.md`

**输出**: 实现计划完整，可以开始编码

### 阶段 4: 代码实现

```
User: "实现 FE-001-01 用户登录功能"
```

**Claude Code 会做什么**:

1. **激活 implementation-manager**
2. **按计划实现**:
   - 创建 `src/auth/login.ts`
   - 实现 `LoginService.authenticate()`
3. **实时记录**:
   - 更新 `feature-to-code-map.md`
   - 记录: FE-001-01 → login.ts:45-89
4. **编写测试**
5. **更新日志**:
   - `implementation-logs/2025-01-19-IMP-001.md`

**输出**: 功能实现完成，映射记录完整

### 阶段 5: 验证和文档

```
User: "验证并生成文档"
```

**Claude Code 会做什么**:

1. **激活 verification-manager**
   - 运行测试
   - 生成测试报告
2. **激活 documentation-generator**
   - 生成技术文档
   - 生成 API 文档
   - 生成用户手册
3. **更新**:
   - `rt-matrix.md`

**输出**: 功能验证通过，文档完整

## 🔄 并行开发示例

### 场景: 同时开发两个功能

```
User: "同时开发 FE-001 (用户认证) 和 FE-002 (数据采集)"
```

**执行流程**:

1. **workflow-coordinator** 激活
2. **创建并行任务**:
   ```
   FE-001: 用户认证
   - 阶段: 实现
   - 完成度: 60%

   FE-002: 数据采集
   - 阶段: 设计
   - 完成度: 40%
   - 依赖: FE-001 (弱依赖)
   ```
3. **建议策略**:
   ```
   FE-001 继续实现
   FE-002 使用依赖忽略继续设计
   ```
4. **更新**:
   - `active-tasks.md`
   - `task-dependencies.md`
   - `dependency-backlog.md`

**结果**: 两个功能独立推进，互不阻塞

## 🎯 精准修改示例

### 场景: 修改现有功能

```
User: "修改 FE-001-02 的密码验证逻辑，增加复杂度检查"
```

**执行流程**:

1. **implementation-manager** 激活
2. **查询映射表**:
   ```
   FE-001-02 → src/auth/validator.ts:12-34
   ```
3. **影响分析**:
   ```
   主要文件: src/auth/validator.ts
   影响功能: FE-001-01 (登录)
   影响测试: TC-001-002
   风险等级: 低
   ```
4. **执行修改**:
   - 读取 `src/auth/validator.ts`
   - 修改 `validate()` 函数
   - 更新测试
5. **更新**:
   - `feature-to-code-map.md` (行号变化)
   - `rt-matrix.md`

**结果**: 精准定位，影响明确，安全修改

## 📊 状态查询

### 随时了解进度

```
User: "我上次做到哪了？"
```

**Claude Code 会显示**:

```
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

🔄 并行任务
- FE-002: 方案设计中 (40%)

⚠️ 依赖忽略
- D-001: IMP-002 依赖 IMP-001 (等待中)

快速恢复: "继续执行 IMP-001 的步骤 4.1.2"
================================
```

## 🔧 依赖管理

### 弱依赖并行开发

```
User: "FE-002 依赖 FE-001，但想并行开发"
```

**Claude Code 会建议**:

```
================================
依赖分析: FE-002 → FE-001
类型: 弱依赖 (可以使用 mock)

建议策略: ✅ 并行开发 + 依赖忽略
1. FE-001 继续实现
2. FE-002 使用 mock 继续开发
3. FE-001 完成后，FE-002 补齐真实调用

代码模板:
// TODO(依赖): FE-001 - 用户认证
// Mock 实现:
const user = mockUser;

预计影响: 低
补齐时间: 约 1 小时
================================
```

### 依赖补齐

当 FE-001 完成时:

```
⚠️ 依赖完成通知
依赖模块: IMP-001 (用户认证)

影响范围:
- D-001: IMP-002 需要补齐认证验证逻辑
  位置: src/data/collector.ts:45
  优先级: P0

补齐步骤:
1. ✅ 确认 IMP-001 已完成
2. ⏳ 移除 mock 代码
3. ⏳ 替换为真实调用
4. ⏳ 运行集成测试
5. ⏳ 更新映射表
```

## 📁 文件组织

### 需求文档

```
01-requirements/
├── raw-requirements/
│   └── REQ-001-用户认证.md
├── functional-requirements/
│   ├── FE-001-用户登录.md
│   └── FE-002-权限管理.md
└── requirement-mapping.md
```

### 设计文档

```
03-design/
├── technical-solutions/
│   └── SOL-001-用户认证技术方案.md
├── architecture/
│   ├── system-architecture.md
│   └── adr/
│       └── adr-001-选择JWT.md
├── api-design/
│   └── api-specifications.md
└── database-design/
    └── schema.md
```

### 实现文档

```
04-implementation/
├── implementation-plans/
│   └── IMP-001-用户认证实现.md
├── code-mapping/
│   └── feature-to-code-map.md
└── implementation-logs/
    └── 2025-01-19-IMP-001.md
```

## 💡 提示和技巧

### 命名规范

- **REQ-XXX**: 原始需求
- **FE-XXX**: 功能需求
- **SOL-XXX**: 技术方案
- **IMP-XXX**: 实现计划
- **TC-XXX**: 测试用例
- **US-XXX**: 用户故事
- **D-XXX**: 依赖项

### 快速导航

```
"显示 FE-001 的所有文档"
→ 列出所有关联的 REQ, SOL, IMP

"显示 FE-001 的代码映射"
→ 查询 feature-to-code-map.md

"显示所有待补齐的依赖"
→ 查询 dependency-backlog.md
```

### 状态恢复

```
"暂停当前任务"
→ 保存上下文到堆栈

"继续上一个任务"
→ 从堆栈恢复上下文

"切换到 IMP-002"
→ 切换任务，保存当前任务
```

## 🎓 学习路径

1. **新手**: 先完成一个完整功能的流程 (REQ → IMP → Code)
2. **进阶**: 尝试并行开发两个功能
3. **高级**: 使用依赖忽略机制处理复杂依赖
4. **专家**: 精准修改和影响分析

## ❓ 常见问题

**Q: 必须严格按照流程吗？**
A: 是的，但每个功能独立。FE-001 可以在设计阶段，同时 FE-002 在需求阶段。

**Q: 可以跳过设计直接实现吗？**
A: 不建议。工作流会验证前置条件，确保设计完整。

**Q: 依赖忽略安全吗？**
A: 安全。mock 是临时的，依赖完成后必须补齐，并有完整跟踪。

**Q: 修改代码时找不到映射怎么办？**
A: 检查 `feature-to-code-map.md`，确认功能已实现并记录映射。

**Q: 可以同时开发多个功能吗？**
A: 强烈推荐。工作流设计用于并行开发，提高效率。

## 📞 获取帮助

遇到问题时：
- 查看完整设计文档: `Claude-Code-Workflow-Design-v1.0.0.md`
- 查看模板文件: `.claude-workflow/templates/`
- 查看状态文件: `.claude-workflow/*.md`

---

**祝你使用愉快！🎉**
