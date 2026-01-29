---
name: requirement-manager
description: |
  Comprehensive requirements engineering and management skill.

  **When to use**:
  - Collecting and documenting raw requirements (REQ-XXX)
  - Decomposing requirements into functional requirements (FE-XXX)
  - Creating user stories (US-XXX)
  - Allocating requirements to platforms (gateway/HMI/cloud/APP/edge AI/Scada/Web Editor)
  - Validating requirement completeness
  - Creating requirement traceability matrix

  **Use cases**:
  - "我们需要添加用户认证功能" → Creates REQ-001
  - "数据采集需要在网关和云端都实现" → Decomposes to platform-specific FEs
  - "检查需求是否完整" → Validates requirement phase completion
  - "拆分需求到各个平台" → Allocates requirements to platforms

  **Workflow stage**: Stage 1 - Requirements & Planning
  **Prerequisite**: Activated by workflow-orchestrator
  **Outputs**:
  - Raw requirement documents (01-requirements/raw-requirements/)
  - Functional requirement documents (01-requirements/functional-requirements/)
  - User stories (01-requirements/user-stories/)
  - Requirement mapping table (01-requirements/requirement-mapping.md)
  - Platform-specific requirement breakdown (07-platforms/*/requirements/)

  **Key responsibilities**:
  - Guide users through requirement elicitation
  - Ensure complete requirement documentation
  - Facilitate requirement decomposition
  - Manage requirement allocations across platforms
  - Validate requirement quality and completeness
---

# Requirement Manager

## Overview

This skill guides users through the complete requirements engineering process, from initial requirement collection to detailed functional requirements decomposition.

## Requirement Collection

### Creating Raw Requirements (REQ-{N})

When user expresses a need:

1. **Elicit information**:
   - What is the business need?
   - Who are the stakeholders?
   - What is the priority?
   - Which platforms are affected?

2. **Create REQ-{N}** using template:
   ```
   .claude-workflow/templates/REQ-template.md
   ```

3. **Document**:
   - Requirement description
   - Business value
   - Acceptance criteria
   - Related requirements
   - Risks and constraints

4. **Save to**:
   ```
   .claude-workflow/01-requirements/raw-requirements/REQ-{N}-{title}.md
   ```

### Requirement Template Fields

**Required**:
- Requirement ID: REQ-{N}
- Title: Clear, concise
- Description: Detailed business need
- Priority: P0 (critical) / P1 (high) / P2 (medium) / P3 (low)
- Platforms: Affected platforms
- Acceptance criteria: Measurable criteria

**Optional**:
- User stories
- Business value
- Dependencies
- Risks
- Constraints

## Requirement Decomposition

### Creating Functional Requirements (FE-{N})

After REQ-{N} is created:

1. **Analyze REQ-{N}**:
   - Break down into major features
   - Identify platform-specific needs
   - Define functional and non-functional requirements

2. **Create FE-{N}** for each major feature:
   ```
   .claude-workflow/templates/FE-template.md
   ```

3. **Allocate to platforms**:
   - Gateway: 网关端
   - HMI: HMI运行端
   - Configuration: 组态端
   - Cloud: 云平台端
   - APP: APP端
   - Edge AI: 边缘智能服务器
   - Scada: Scada软件
   - Web Editor: Web可视化编辑器

4. **Save to**:
   ```
   .claude-workflow/01-requirements/functional-requirements/FE-{N}-{title}.md
   ```

5. **Create platform-specific copies**:
   ```
   .claude-workflow/07-platforms/{platform}/requirements/FE-{N}-{title}.md
   ```

### Decomposition Guidelines

**Level of granularity**:
- One FE should implement one cohesive feature
- FE should be implementable in 1-2 weeks
- FE should have clear acceptance criteria
- FE should map to specific platforms

**Platform allocation**:
- Analyze which platforms need this feature
- Create platform-specific sub-requirements if needed
- Document platform-specific variations

## User Stories

### Creating User Stories (US-{N})

For each FE-{N}:

1. **Format**:
   ```
   As a {role}
   I want {feature}
   So that {value}
   ```

2. **Save to**:
   ```
   .claude-workflow/01-requirements/user-stories/US-{N}-{title}.md
   ```

## Requirement Mapping

### Creating Requirement Mapping Table

Track relationships between requirements:

```
.claude-workflow/01-requirements/requirement-mapping.md
```

**Structure**:
| REQ-{N} | FE-{N} | Platform | Status |
|---------|--------|----------|--------|
| REQ-001 | FE-001 | Gateway | ⏳ |
| REQ-001 | FE-002 | Cloud | ⏳ |
| REQ-002 | FE-003 | HMI | 📋 |

## Requirement Validation

### Completeness Checklist

Before allowing transition to design stage:

**REQ-{N} validation**:
- [ ] All required fields filled
- [ ] Description is clear and unambiguous
- [ ] Acceptance criteria are measurable
- [ ] Priority is assigned
- [ ] Platforms are identified
- [ ] Dependencies are documented

**FE-{N} validation**:
- [ ] Linked to parent REQ-{N}
- [ ] Functional description is complete
- [ ] Platform allocation is clear
- [ ] Non-functional requirements specified
- [ ] Acceptance criteria are testable
- [ ] User stories created

**Overall validation**:
- [ ] All REQs decomposed to FEs
- [ ] All FEs have platform assignments
- [ ] Requirement mapping table updated
- [ ] No orphan requirements
- [ ] Traceability established

## Multi-Platform Requirements

### Platform-Specific Considerations

**Gateway**:
- Focus on data collection, edge processing
- Resource constraints
- Real-time requirements

**HMI**:
- Focus on user interface, visualization
- Response time requirements
- Usability requirements

**Configuration**:
- Focus on configuration tools
- Ease of use
- Validation requirements

**Cloud**:
- Focus on data aggregation, analytics
- Scalability requirements
- Security requirements

**APP**:
- Focus on mobile experience
- Offline support
- Performance requirements

**Edge AI**:
- Focus on AI/ML capabilities
- Model performance
- Resource optimization

**Scada**:
- Focus on industrial control
- Reliability requirements
- Real-time constraints

**Web Editor**:
- Focus on web-based editing
- Browser compatibility
- Collaboration features

### Cross-Platform Features

For features spanning multiple platforms:

1. **Identify common functionality**
2. **Define platform-specific variations**
3. **Create separate FEs for each platform**:
   - FE-001-Gateway: 网关端数据采集
   - FE-001-Cloud: 云端数据聚合
4. **Document platform interfaces**
5. **Define data synchronization**

## Dependency Management

### Identifying Dependencies

Between requirements:

- **Functional dependency**: FE-002 requires FE-001
- **Data dependency**: FE-003 needs data from FE-001
- **Platform dependency**: Gateway FE requires Cloud FE

### Documenting Dependencies

In FE-{N} document:
```markdown
## Dependency Relations
- **Preceding**: FE-{N}
- **Following**: FE-{N}
- **Parallel**: FE-{N}
- **Cross-platform**: FE-{N}
```

Update task-dependencies.md:
```markdown
FE-{N} → depends on → FE-{M}
Type: Strong/Weak/None
```

## Quality Assurance

### Requirement Quality Metrics

**Good requirements**:
- ✅ Unambiguous: Clear, single interpretation
- ✅ Complete: All necessary information
- ✅ Consistent: No contradictions
- ✅ Testable: Can verify acceptance
- ✅ Traceable: Can track through lifecycle
- ✅ Feasible: Technically possible
- ✅ Prioritized: Clear priority

**Common issues to avoid**:
- ❌ Vague language ("should be fast")
- ❌ Subjective criteria ("user-friendly")
- ❌ Technical solutions in requirements
- ❌ Missing acceptance criteria
- ❌ Unprioritized requirements

## Workflow Integration

### Entry Point
Activated by `workflow-orchestrator` for Stage 1.

### Exit Criteria
Transition to design stage when:
1. All requirements documented
2. All requirements decomposed to FEs
3. All FEs allocated to platforms
4. Requirement mapping complete
5. Validation passed

### Handoff to Design Stage
Provide to `solution-designer`:
- List of FEs for design
- Platform allocations
- Dependencies between FEs
- Priority information
- Non-functional requirements

## Usage Examples

### Example 1: Simple Requirement

```
User: "我们需要添加用户登录功能"

→ requirement-manager activates
→ Asks clarifying questions:
  - Which platforms need login?
  - What authentication methods?
  - Any specific requirements?

→ Creates REQ-001: 用户认证
→ Decomposes to:
  - FE-001: 用户登录 (Gateway, Cloud)
  - FE-002: 权限管理 (All platforms)

→ Validates completeness
→ Notifies workflow-orchestrator
```

### Example 2: Multi-Platform Requirement

```
User: "数据采集需要在网关和云端都实现"

→ requirement-manager activates
→ Creates REQ-002: 数据采集
→ Decomposes to:
  - FE-003-Gateway: 网关数据采集
    - Focus: Edge collection, preprocessing
    - Non-functional: Real-time, low latency
  - FE-003-Cloud: 云端数据聚合
    - Focus: Data aggregation, storage
    - Non-functional: Scalable, durable

→ Creates platform-specific requirement documents
→ Updates requirement mapping
→ Validates completeness
```

### Example 3: Completeness Check

```
User: "检查 REQ-001 的需求是否完整"

→ requirement-manager activates
→ Reads REQ-001 document
→ Checks validation checklist:
  ✅ All required fields present
  ✅ Description clear
  ✅ Acceptance criteria measurable
  ✅ Priority assigned
  ✅ Platforms identified
  ✅ Decomposed to FE-001, FE-002
  ✅ FEs have platform assignments

→ Reports: "REQ-001 is complete and ready for design stage"
```

## Templates Location

All templates in:
```
.claude-workflow/templates/
├── REQ-template.md
├── FE-template.md
└── US-template.md
```

## Output Files

Created/Updated files:
```
.claude-workflow/
├── 01-requirements/
│   ├── raw-requirements/
│   │   └── REQ-{N}-{title}.md
│   ├── functional-requirements/
│   │   └── FE-{N}-{title}.md
│   ├── user-stories/
│   │   └── US-{N}-{title}.md
│   └── requirement-mapping.md
└── 07-platforms/
    ├── gateway/requirements/
    ├── hmi/requirements/
    ├── cloud/requirements/
    └── ...
```
