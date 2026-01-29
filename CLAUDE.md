# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the **Claude Code Workflow Template** - a meta-project providing a structured development workflow management system for Claude Code. It enables:

- Single-feature independent progression through REQ → FE → SOL → IMP → Code
- Parallel development of multiple features with dependency management
- Complete traceability from requirements to code implementation
- Multi-platform support (Gateway, HMI, Cloud, APP, Edge AI, Scada, Web Editor)

## Architecture

### Core Workflow System

The workflow is managed through specialized Skills located in `skills/`:

1. **workflow-orchestrator** (Master entry point)
   - MANDATORY for ALL development work
   - Coordinates all phases from requirements through implementation
   - Validates phase transitions before proceeding
   - Manages task stack and state

2. **Specialist Skills** (activated by orchestrator):
   - **requirement-manager**: Creates REQ-{N}, decomposes to FE-{N}, allocates to platforms
   - **solution-designer**: Creates SOL-{N} with architecture, APIs, and technical design
   - **implementation-manager**: Creates IMP-{N}, executes implementation, maintains code mapping
   - **workflow-coordinator**: Manages task stack, pause/resume, parallel tasks

### Workflow Stages (Per Feature)

```
Stage 1: Requirements & Planning
  ↓ REQ-{N} created, decomposed to FE-{N}, allocated to platforms

Stage 2: Solution Design
  ↓ SOL-{N} with architecture, APIs, technical specifications

Stage 3: Implementation Planning
  ↓ IMP-{N} with task breakdown and code mapping table

Stage 4: Code Implementation
  ↓ Actual code written, feature-to-code mapping recorded

Stage 5: Verification & Documentation
  ↓ Tests, verification, documentation generation
```

**Key Principle**: Each feature progresses independently. FE-001 can be in implementation while FE-002 is in design and FE-003 is in requirements.

### State Management Files

Located in `.claude-workflow/`:

- **config.yml**: Project configuration and enabled platforms
- **current-phase.md**: Current active task, task stack for pause/resume
- **parallel-tasks/active-tasks.md**: All parallel tasks with status
- **parallel-tasks/task-dependencies.md**: Dependency graph (strong/weak/none)
- **dependency-backlog.md**: Pending dependency completions and mock locations
- **feature-to-code-map.md**: Bidirectional mapping between FE-ID and code locations
- **rt-matrix.md**: Requirements traceability matrix

### Directory Structure

```
.claude-workflow/
├── 01-requirements/          # REQ-{N}, FE-{N}, US-{N}
├── 02-planning/             # Project planning documents
├── 03-design/               # SOL-{N}, architecture, APIs, ADRs
├── 04-implementation/        # IMP-{N}, code mapping, implementation logs
├── 05-verification/         # Test plans, test cases, verification reports
├── 06-documentation/        # Technical docs, user guides, API docs
├── 07-platforms/            # Platform-specific breakdown (gateway, hmi, cloud, etc.)
└── templates/               # Document templates for REQ, FE, SOL, IMP

skills/                       # Claude Code Skills definitions
├── workflow-orchestrator/   # Master orchestrator skill
├── requirement-manager/     # Requirements engineering
├── solution-designer/       # Technical solution design
├── implementation-manager/  # Implementation planning & execution
└── workflow-coordinator/    # State management & context recovery
```

## Development Workflow

### Starting New Features

**CRITICAL**: ALWAYS use workflow-orchestrator as entry point. Never skip directly to coding.

```
User: "我们需要添加用户认证功能"
→ workflow-orchestrator activates
→ requirement-manager guides REQ-001 → FE-001 creation
→ solution-designer creates SOL-001
→ implementation-manager creates IMP-001
→ Code implementation executed
→ Verification and documentation
```

### Modifying Existing Features

```
User: "修改 FE-001-02 的验证逻辑"
→ implementation-manager queries feature-to-code-map.md
→ Displays impact analysis
→ Executes precise modification
→ Updates all tracking files
```

### Parallel Development

Multiple features can progress simultaneously:

```
FE-001 (Implementation) ━━━━━━━━→ Complete
FE-002 (Design) ━━━━━━━━━━━━━━━━━→ SOL-002
FE-003 (Requirements) ━━━━━━━━━━→ FE-003
```

### Dependency Ignore Mechanism

For weak dependencies, use mock implementations:

```typescript
// TODO(依赖): FE-001 - 用户认证
// 说明: 此处需要调用认证服务验证用户身份
// 当前状态: 使用 mock 实现
// 依赖模块: IMP-001
// 补齐优先级: P0

export class DataCollector {
  async collectData(userId: string) {
    // === 依赖忽略开始 ===
    const user = { id: userId, name: 'Mock User' };
    // === 依赖忽略结束 ===

    // Main business logic continues
  }
}
```

Record in `dependency-backlog.md` and complete when dependency is ready.

## Platform Support

The workflow supports 8 platforms (configure in `.claude-workflow/config.yml`):

- **gateway**: 网关端
- **hmi**: HMI运行端
- **configuration**: 组态端
- **cloud**: 云平台端
- **app**: APP端
- **edge-ai**: 边缘智能服务器
- **scada**: Scada软件
- **web-editor**: Web可视化编辑器

Requirements (FE-{N}) are allocated to specific platforms during Stage 1.

## Key Files Reference

### When User Asks...

- **"我上次做到哪了？"** → Read `.claude-workflow/current-phase.md`
- **"显示所有并行任务"** → Read `.claude-workflow/parallel-tasks/active-tasks.md`
- **"检查依赖状态"** → Read `.claude-workflow/dependency-backlog.md`
- **"修改 FE-XXX-YY"** → Read `.claude-workflow/feature-to-code-map.md` for location
- **"验证是否可以进入实现"** → Validate prerequisites in workflow-orchestrator

### Workflow State Queries

The workflow-coordinator skill provides:
- Task stack operations (push/pop/pause/resume)
- Phase transition validation
- Context recovery for interrupted sessions
- Parallel task status queries

## Best Practices

### DO:
- Always activate workflow-orchestrator for new features
- Maintain complete traceability (update all tracking files)
- Use dependency ignore mechanism for parallel development
- Record feature-to-code mapping during implementation
- Validate phase prerequisites before transitioning
- Keep state files current after every action

### DON'T:
- Skip workflow stages (REQ → FE → SOL → IMP → Code)
- Manually modify state files (let Claude Code manage them)
- Start coding without implementation plan (IMP-{N})
- Forget to update feature-to-code-map.md
- Leave dependency ignore markers unresolved
- Block features when parallel development is possible

## Initialization Script

For new projects using this template:

```bash
python .claude-workflow/scripts/init-workflow.py \
    --project-name "我的项目" \
    --project-type "iot-platform" \
    --platforms "gateway,cloud,hmi" \
    --template-source .
```

This creates the full directory structure, configuration, and initial state files.

## Traceability

The system maintains complete bidirectional traceability:

- **Forward**: REQ → FE → SOL → IMP → Code → Test
- **Backward**: Code → IMP → SOL → FE → REQ
- **Impact Analysis**: Change FE-001 → identify affected code/tests/platforms
- **Coverage**: Verify all requirements have corresponding implementation and tests

All mappings are maintained in:
- `.claude-workflow/feature-to-code-map.md`
- `.claude-workflow/rt-matrix.md`
