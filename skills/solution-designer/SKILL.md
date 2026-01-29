---
name: solution-designer
description: |
  Technical solution and architecture design skill.

  **PREREQUISITE**: Can ONLY be used after requirement-manager phase is complete.
  workflow-orchestrator will validate prerequisites before activation.

  **When to use**:
  - Designing system architecture (03-design/architecture/)
  - Creating technical solutions (SOL-XXX)
  - Designing APIs and interfaces
  - Making architectural decisions (ADR)
  - Database schema design
  - Component design

  **Use cases**:
  - "设计用户认证的技术方案" → Creates SOL-001
  - "设计网关和云平台的通信接口" → Creates API specifications
  - "选择MQTT还是HTTP？" → Creates ADR documenting the decision
  - "设计数据库结构" → Creates schema design
  - "设计系统架构" → Creates architecture documents

  **Workflow stage**: Stage 2 - Solution Design
  **Input**: Functional requirements (FE-{N}) from requirement-manager
  **Outputs**:
  - Architecture documents (03-design/architecture/)
  - Technical solution documents (03-design/technical-solutions/SOL-{N}.md)
  - API specifications (03-design/api-design/)
  - ADR records (03-design/architecture/adr/)
  - Database schemas (03-design/database-design/)

  **Key responsibilities**:
  - Guide users through technical design process
  - Ensure complete solution documentation
  - Facilitate architecture decision making
  - Create API and interface specifications
  - Document design rationale and trade-offs
---

# Solution Designer

## Overview

This skill guides users through creating complete technical solutions for functional requirements, ensuring all design decisions are documented and justified.

## Prerequisites

### Must Have (Validated by workflow-orchestrator)
- ✅ FE-{N} requirement exists
- ✅ FE-{N} has complete functional description
- ✅ FE-{N} has platform allocation
- ✅ FE-{N} has non-functional requirements

### Input From requirement-manager
- Functional requirement document (FE-{N})
- Platform assignments
- Dependencies
- Non-functional requirements

## Solution Design Process

### Step 1: Analyze Functional Requirement

1. **Read FE-{N} document**
2. **Extract key information**:
   - Functional requirements
   - Platform constraints
   - Non-functional requirements (performance, security, reliability)
   - Dependencies on other features

### Step 2: Create Technical Solution (SOL-{N})

1. **Copy template**:
   ```
   .claude-workflow/templates/SOL-template.md
   ```

2. **Fill in sections**:
   - Background and objectives
   - Solution overview
   - Technology selection
   - Architecture design
   - Interface design
   - Data models
   - Implementation plan
   - Risk assessment
   - Impact analysis

3. **Save to**:
   ```
   .claude-workflow/03-design/technical-solutions/SOL-{N}-{title}.md
   ```

### Step 3: Design Architecture

#### System Architecture

Create architecture diagram:
```
.claude-workflow/03-design/architecture/system-architecture.md
```

**Include**:
- System components
- Component interactions
- Data flows
- External interfaces
- Technology stack

#### Component Design

Document each component:
```
.claude-workflow/03-design/architecture/component-design.md
```

**For each component**:
- Name
- Responsibility
- Interfaces
- Dependencies
- Technology stack

### Step 4: Design APIs

#### API Specifications

For each API endpoint:
```
.claude-workflow/03-design/api-design/api-specifications.md
```

**Include**:
- Endpoint URL
- HTTP method
- Request parameters
- Response format
- Error codes
- Authentication requirements
- Rate limits

**Example**:
```markdown
### POST /api/auth/login

**Description**: Authenticates user and returns access token

**Request**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Response**:
```json
{
  "code": 200,
  "data": {
    "token": "string",
    "expiresIn": 3600
  },
  "message": "success"
}
```

**Error Responses**:
- 400: Invalid credentials
- 401: Unauthorized
- 500: Server error
```

### Step 5: Design Data Models

#### Database Schema

```
.claude-workflow/03-design/database-design/schema.md
```

**For each table/collection**:
- Table name
- Columns/Fields
- Data types
- Constraints
- Indexes
- Relationships

**Example**:
```markdown
### Table: users

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK, NOT NULL | User ID |
| username | VARCHAR(50) | UNIQUE, NOT NULL | Username |
| password_hash | VARCHAR(255) | NOT NULL | Password hash |
| created_at | TIMESTAMP | NOT NULL | Creation time |
| updated_at | TIMESTAMP | NOT NULL | Last update |

**Indexes**:
- idx_username: username
- idx_created_at: created_at

**Relationships**:
- One-to-many with sessions
```

### Step 6: Create Architecture Decision Records (ADR)

For significant architectural decisions:

```
.claude-workflow/03-design/architecture/adr/adr-{N}-{decision}.md
```

**ADR Format**:
```markdown
# ADR-{N}: {Decision Title}

## Status
Proposed / Accepted / Rejected / Deprecated / Superseded

## Date
{YYYY-MM-DD}

## Context
{What is the problem we're trying to solve?}

## Decision
{What did we decide?}

## Consequences
{What are the results of this decision?}

**Positive**:
- {Benefit 1}
- {Benefit 2}

**Negative**:
- {Drawback 1}
- {Drawback 2}

## Alternatives Considered
- {Alternative 1}: {Reason for rejection}
- {Alternative 2}: {Reason for rejection}
```

**Common ADR topics**:
- Technology selection (MQTT vs HTTP)
- Architecture pattern (Microservices vs Monolith)
- Database choice (SQL vs NoSQL)
- Authentication method (JWT vs Session)
- Communication protocol (REST vs GraphQL vs gRPC)

### Step 7: Document Implementation Plan

In SOL-{N} document:

**Phases**:
1. **Phase 1**: Core functionality
2. **Phase 2**: Integration
3. **Phase 3**: Optimization

**Tasks**:
- Break down into implementable tasks
- Estimate effort
- Identify dependencies
- Sequence tasks

### Step 8: Risk Assessment

In SOL-{N} document:

**Types of risks**:
- Technical: Feasibility, complexity
- Performance: Scalability, latency
- Security: Data protection, authentication
- Operational: Deployment, monitoring

**For each risk**:
- Impact: High / Medium / Low
- Probability: High / Medium / Low
- Mitigation strategy

### Step 9: Impact Analysis

In SOL-{N} document:

**Analyze**:
- Affected systems
- Affected components
- Data migration needs
- Backward compatibility
- Performance impact
- Security implications

## Validation Checklist

Before allowing transition to implementation stage:

**SOL-{N} validation**:
- [ ] Background and objectives clear
- [ ] Solution overview complete
- [ ] Technology selection justified
- [ ] Architecture diagram included
- [ ] Component design documented
- [ ] API/interface specifications complete
- [ ] Data models defined
- [ ] Implementation plan detailed
- [ ] Risks identified and mitigated
- [ ] Impact analysis complete
- [ ] ADRs created for key decisions

**Cross-check with FE-{N}**:
- [ ] All functional requirements addressed
- [ ] All non-functional requirements addressed
- [ ] Platform constraints considered
- [ ] Dependencies handled

## Design Principles

### Quality Attributes

**Design for**:
- **Performance**: Response time, throughput
- **Scalability**: Handle increased load
- **Reliability**: Availability, fault tolerance
- **Security**: Authentication, authorization, data protection
- **Maintainability**: Code organization, documentation
- **Usability**: API design, error messages

### Best Practices

1. **Simplicity**: Simple solutions over complex ones
2. **Modularity**: Loosely coupled components
3. **Extensibility**: Easy to add features
4. **Testability**: Design for testing
5. **Observability**: Logging, metrics
6. **Documentation**: Document decisions and rationale

## Technology Selection Guidance

### Evaluation Criteria

When selecting technologies, consider:

| Criterion | Questions |
|-----------|-----------|
| **Fit for purpose** | Does it solve our problem? |
| **Maturity** | Is it production-ready? |
| **Community** | Active community and support? |
| **Performance** | Meets performance requirements? |
| **Scalability** | Can handle future growth? |
| **Security** | Security track record? |
| **Learning curve** | Team can learn it? |
| **License** | Compatible with project? |
| **Integration** | Works with existing stack? |

### Common Technology Choices

**Communication**:
- REST: Simple, widely adopted
- GraphQL: Flexible queries
- gRPC: High performance, type-safe
- WebSocket: Real-time communication

**Messaging**:
- MQTT: Lightweight, IoT
- Kafka: High throughput, stream processing
- RabbitMQ: Feature-rich message broker
- Redis Pub/Sub: Simple pub/sub

**Database**:
- PostgreSQL: Relational, feature-rich
- MongoDB: Document, flexible schema
- Redis: In-memory, fast
- InfluxDB: Time-series data

**Authentication**:
- JWT: Stateless, scalable
- OAuth 2.0: Third-party authentication
- Session: Simple, stateful

## Multi-Platform Design

### Platform-Specific Considerations

**Gateway**:
- Resource constraints (CPU, memory)
- Edge computing capabilities
- Offline operation
- Local data caching

**HMI**:
- Responsive UI
- Real-time updates
- Touch-friendly interface
- Data visualization

**Cloud**:
- Horizontal scalability
- Data aggregation
- Analytics capabilities
- High availability

**APP**:
- Mobile-optimized UI
- Offline support
- Push notifications
- Battery efficiency

**Edge AI**:
- ML model deployment
- Model optimization
- Inference acceleration
- Resource monitoring

**Scada**:
- Real-time control
- Deterministic response
- Industrial protocols
- High reliability

**Web Editor**:
- Rich editing experience
- Collaboration features
- Browser compatibility
- Performance optimization

### Cross-Platform Integration

**Design for**:
- Consistent APIs across platforms
- Data synchronization
- Conflict resolution
- Graceful degradation
- Platform-specific optimizations

## Workflow Integration

### Entry Point
Activated by `workflow-orchestrator` for Stage 2.

### Exit Criteria
Transition to implementation planning when:
1. SOL-{N} document complete
2. Architecture designed
3. APIs specified
4. Data models defined
5. Risks assessed
6. Validation passed

### Handoff to Implementation Stage
Provide to `implementation-manager`:
- SOL-{N} document
- Architecture diagrams
- API specifications
- Data models
- Implementation plan
- Technology stack

## Usage Examples

### Example 1: Simple Solution Design

```
User: "设计用户认证的技术方案"

→ solution-designer activates
→ Reads FE-001: 用户登录
→ Creates SOL-001: 用户认证技术方案

Design decisions:
- Authentication: JWT tokens
- Password hashing: bcrypt
- Database: PostgreSQL users table
- API: POST /api/auth/login

→ Creates architecture diagram
→ Defines API specifications
→ Designs database schema
→ Documents decision in ADR

→ Validates completeness
→ Notifies workflow-orchestrator
```

### Example 2: Architecture Decision

```
User: "选择MQTT还是HTTP进行设备通信？"

→ solution-designer activates
→ Evaluates options:

MQTT:
  Pros: Lightweight, pub/sub, QoS
  Cons: Less familiar, needs broker

HTTP:
  Pros: Simple, widely adopted
  Cons: Request/response, heavier

→ Creates ADR-001: 设备通信协议选择
→ Decision: MQTT for device communication
→ Rationale: Lightweight, pub/sub pattern fits IoT
→ Consequences documented

→ Updates architecture diagram
```

### Example 3: API Design

```
User: "设计网关和云平台的通信接口"

→ solution-designer activates
→ Analyzes requirements:
  - Real-time data upload
  - Command delivery
  - Configuration sync
  - Status reporting

→ Designs API specification:
  POST /api/gateway/data - Upload data
  POST /api/gateway/command - Deliver command
  GET /api/gateway/config - Get configuration
  PUT /api/gateway/config - Update configuration

→ Documents request/response formats
→ Defines error codes
→ Specifies authentication
→ Adds rate limiting

→ Saves to api-specifications.md
```

## Templates Location

All templates in:
```
.claude-workflow/templates/
└── SOL-template.md
```

## Output Files

Created/Updated files:
```
.claude-workflow/
├── 03-design/
│   ├── architecture/
│   │   ├── system-architecture.md
│   │   ├── component-design.md
│   │   └── adr/
│   │       └── adr-{N}-{decision}.md
│   ├── technical-solutions/
│   │   └── SOL-{N}-{title}.md
│   ├── api-design/
│   │   └── api-specifications.md
│   └── database-design/
│       └── schema.md
```

## Integration Points

### Upstream
- Receives FE-{N} from requirement-manager

### Downstream
- Provides SOL-{N} to implementation-manager
- Provides architecture to implementation-manager
- Provides API specs to implementation-manager
