---
name: workflow-orchestrator
description: |
  Master orchestrator for comprehensive software development workflow management.
  Coordinates all phases from requirements through implementation to documentation.

  **MANDATORY USAGE**: This MUST be used as the entry point for ALL development work.
  No coding or design work should start without this skill being activated first.

  **When to use**:
  - Starting ANY new project or feature
  - Managing complex multi-platform projects (gateway, HMI, cloud, APP, edge AI, Scada, Web Editor)
  - Coordinating cross-team development
  - Ensuring proper workflow compliance
  - Recovering interrupted work sessions

  **Workflow phases enforced** (per feature):
  1. Requirements & Planning (REQ → FE)
  2. Solution Design (FE → SOL)
  3. Implementation Planning (SOL → IMP)
  4. Code Implementation (IMP → Code)
  5. Verification & Documentation

  **Key principles**:
  - Each feature progresses independently (not blocking other features)
  - Parallel development is encouraged
  - Dependency ignore mechanism for parallel work
  - Complete traceability from requirements to code

  **Key responsibilities**:
  - Validates phase transition prerequisites (per feature)
  - Activates appropriate specialist skills
  - Maintains task stack and state
  - Enforces documentation completeness
  - Provides context recovery mechanism

  **Use cases**:
  - "我们需要添加用户认证功能" → Orchestrates full workflow
  - "添加数据采集功能到网关和云端" → Coordinates multi-platform work
  - "同时开发 FE-001 和 FE-002" → Manages parallel tasks
  - "我上次做到哪了？" → Restores context and continues
---

# Workflow Orchestrator

## Overview

This skill manages the complete software development lifecycle, ensuring proper workflow adherence while enabling parallel development and complete traceability.

## Core Principles

### 1. Single-Feature Progression
Each feature progresses independently through stages:
```
REQ-{N} → FE-{N} → SOL-{N} → IMP-{N} → Code
```
Features do NOT block each other. FE-001 can enter design phase even if REQ-002 is not complete.

### 2. Parallel Development
Multiple features can be in different phases simultaneously:
- FE-001 in implementation
- FE-002 in design
- FE-003 in requirements

### 3. Dependency Management
- **Strong dependency**: Serial development
- **Weak dependency**: Parallel development with dependency ignore mechanism
- **No dependency**: Free parallel development

## Workflow Stages

### Stage 1: Requirements & Planning
**Prerequisites** (per feature):
- [ ] REQ-{N} created with complete metadata
- [ ] REQ-{N} decomposed into FE-{N}
- [ ] FE-{N} assigned to specific platforms
- [ ] FE-{N} has complete functional description

**Activates**: `requirement-manager`

**Validation**: Only checks the specific feature's requirements, not other features.

### Stage 2: Solution Design
**Prerequisites** (per feature):
- [ ] FE-{N} requirement decomposition complete
- [ ] SOL-{N} technical solution document created
- [ ] SOL-{N} interface design defined
- [ ] SOL-{N} design reviewed and approved

**Activates**: `solution-designer`

**Validation**: Only checks SOL-{N} completeness, not other solutions.

### Stage 3: Implementation Planning
**Prerequisites** (per feature):
- [ ] SOL-{N} design complete
- [ ] IMP-{N} implementation plan created
- [ ] IMP-{N} code mapping table designed
- [ ] IMP-{N} test strategy defined
- [ ] IMP-{N} implementation plan reviewed

**Activates**: `implementation-manager`

**Validation**: Only checks IMP-{N} readiness, not other implementations.

### Stage 4: Code Implementation
**Prerequisites** (per feature):
- [ ] IMP-{N} approved
- [ ] Development environment ready
- [ ] Dependencies available or mocked

**Activates**: `implementation-manager`

### Stage 5: Verification & Documentation
**Prerequisites** (per feature):
- [ ] Code implementation complete
- [ ] Feature-to-code mapping updated
- [ ] RT matrix updated

**Activates**: `verification-manager` then `documentation-generator`

## State Management

### Files Maintained

1. **`.claude-workflow/current-phase.md`**
   - Current active task
   - Task stack for pause/resume
   - Context information

2. **`.claude-workflow/parallel-tasks/active-tasks.md`**
   - All active parallel tasks
   - Task status and completion percentage
   - Dependency information

3. **`.claude-workflow/parallel-tasks/task-dependencies.md`**
   - Dependency graph
   - Dependency type (strong/weak/none)
   - Dependency ignore records

4. **`.claude-workflow/dependency-backlog.md`**
   - Pending dependency completions
   - Mock locations
   - Priorities

## Phase Transition Logic

### Entering Design Stage
```
IF FE-{N}.status == "已分解"
AND FE-{N}.platforms != empty
AND FE-{N}.description IS complete
THEN ALLOW SOL-{N} creation
```

### Entering Implementation Stage
```
IF SOL-{N}.status == "已批准"
AND SOL-{N}.interface_design IS complete
AND SOL-{N}.technical_design IS complete
THEN ALLOW IMP-{N} creation
```

### Entering Coding Stage
```
IF IMP-{N}.status == "已批准"
AND IMP-{N}.code_mapping EXISTS
AND IMP-{N}.test_strategy EXISTS
THEN ALLOW coding
```

## Usage Patterns

### Starting a New Feature

```
User: "我们需要添加用户认证功能"

→ workflow-orchestrator activates
→ Reads project config from .claude-workflow/config.yml
→ Activates requirement-manager
→ Guides user through REQ-001 → FE-001 creation
→ Validates FE-001 completeness
→ Activates solution-designer for SOL-001
→ Continues through all stages...
```

### Parallel Development

```
User: "同时开发 FE-001 和 FE-002"

→ workflow-orchestrator activates
→ Creates entries in active-tasks.md for both
→ Manages independent progress
→ Tracks dependencies between features
→ Enables dependency ignore if needed
```

### Context Recovery

```
User: "我上次做到哪了？"

→ workflow-orchestrator activates
→ Reads current-phase.md
→ Reports current task and progress
→ Offers to continue from last point
```

### Modifying Existing Feature

```
User: "修改 FE-001-02 的密码验证逻辑"

→ workflow-orchestrator activates
→ Activates implementation-manager
→ Queries feature-to-code-map.md
→ Displays impact analysis
→ Executes precise modification
→ Updates all tracking files
```

## Dependency Ignore Mechanism

### When to Use
- Feature B depends on Feature A
- Want to develop both in parallel
- Dependency is weak (can use mock)

### Process

1. **Detect dependency**: FE-002 depends on FE-001
2. **Suggest parallel + ignore**: "FE-002 can proceed with mock"
3. **Create dependency ignore record**:
   ```typescript
   // TODO(依赖): FE-001 - 用户认证
   // 当前状态: 使用 mock 实现
   // 依赖模块: IMP-001
   // 补齐优先级: P0
   ```
4. **Update dependency-backlog.md**
5. **Complete main logic**: Implement FE-002 with mock
6. **Wait for dependency**: FE-001 completes
7. **Quick completion**: Replace mock with real call

## Checklist Validation

### Feature-Level Validation

Before allowing phase transition, validate:

**Requirements Stage**:
- [ ] REQ-{N} has all required metadata
- [ ] REQ-{N} has clear acceptance criteria
- [ ] FE-{N} is decomposed from REQ-{N}
- [ ] FE-{N} has platform assignment
- [ ] FE-{N} has functional decomposition

**Design Stage**:
- [ ] SOL-{N} has complete technical design
- [ ] SOL-{N} has API/interface specifications
- [ ] SOL-{N} has architecture diagram
- [ ] SOL-{N} has risk assessment
- [ ] SOL-{N} is reviewed

**Implementation Stage**:
- [ ] IMP-{N} has task breakdown
- [ ] IMP-{N} has code mapping table
- [ ] IMP-{N} has test strategy
- [ ] IMP-{N} has dependency analysis
- [ ] IMP-{N} is reviewed

## Error Handling

### Validation Failure
When validation fails:
1. Report specific missing items
2. Show what's required for next stage
3. Offer to help complete missing items
4. Do NOT allow transition until fixed

### Dependency Conflict
When dependency conflict detected:
1. Analyze dependency type (strong/weak)
2. Suggest appropriate strategy:
   - Strong: Serial development
   - Weak: Parallel with dependency ignore
3. Document decision in task-dependencies.md

### State Corruption
If state files are corrupted:
1. Detect inconsistency
2. Report specific issue
3. Offer recovery options
4. Restore from backup if available

## Integration with Other Skills

### requirement-manager
- Activated for Stage 1
- Creates REQ-{N} and FE-{N}
- Validates requirement completeness

### solution-designer
- Activated for Stage 2
- Creates SOL-{N}
- Validates design completeness

### implementation-manager
- Activated for Stages 3 & 4
- Creates IMP-{N}
- Executes implementation
- Maintains code mapping

### verification-manager
- Activated for Stage 5
- Creates test plans
- Executes tests
- Generates verification reports

### documentation-generator
- Activated for Stage 5
- Generates all documentation
- Updates RT matrix

### workflow-coordinator
- Manages task stack
- Handles pause/resume
- Validates phase transitions
- Provides context recovery

## Best Practices

1. **Always validate before transition**: Never skip validation
2. **Maintain complete traceability**: Update all tracking files
3. **Encourage parallel work**: Use dependency ignore when appropriate
4. **Keep state files current**: Update after every action
5. **Document decisions**: Use ADR for important decisions
6. **Review regularly**: Periodic reviews of progress and blockers

## Troubleshooting

### Issue: Skill not activating
**Solution**:
- Check if SKILL.md format is correct
- Verify description field is clear
- Restart Claude Code

### Issue: Phase transition blocked
**Solution**:
- Check current-phase.md for errors
- Verify all prerequisites are met
- Review validation error messages

### Issue: Parallel tasks conflicting
**Solution**:
- Review task-dependencies.md
- Adjust task priorities
- Consider serial development for strong dependencies

### Issue: Lost context
**Solution**:
- Use "我上次做到哪了？" to recover
- Check current-phase.md
- Review implementation logs
