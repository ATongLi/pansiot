# ADR-004: 采用租户级RBAC权限模型

## 元数据
- **决策ID**: ADR-004
- **决策状态**: 已接受
- **决策日期**: 2026-01-27
- **决策人**: Claude Code
- **评审人**: 待指定
- **相关功能**: FE-007-03（RBAC权限模型）

## 上下文（Context）

我们需要设计权限模型，支持：

1. **多租户隔离**: 每个租户独立管理自己的角色和权限
2. **三级权限**: 角色 → 功能权限 → 操作权限
3. **灵活配置**: 支持自定义角色和权限组合

### 候选方案

1. **全局RBAC**
   - 所有租户共享一套角色和权限
   - 平台超管统一管理

2. **租户级RBAC**（推荐）
   - 每个租户独立管理角色和权限
   - 租户管理员可以自定义角色

3. **混合模式**
   - 预定义角色全局共享
   - 自定义角色租户级管理

## 决策（Decision）

**采用租户级RBAC权限模型**

## 理由（Rationale）

### 权限模型设计

#### 三级权限架构

```
Level 1: 系统角色（Role）
  ↓ 包含多个
Level 2: 功能权限（Feature Permission）
  ↓ 包含多个
Level 3: 操作权限（Action Permission）
```

**示例**:
```
系统管理员（Role）
  ↓
  ├─ 用户管理（Feature Permission）
  │   ├─ 查看（Action Permission）
  │   ├─ 新增（Action Permission）
  │   ├─ 编辑（Action Permission）
  │   └─ 删除（Action Permission）
  │
  └─ 设备管理（Feature Permission）
      ├─ 查看
      ├─ 新增
      ├─ 编辑
      └─ 删除
```

### 方案对比

#### 方案1: 全局RBAC

**数据模型**:
```sql
-- 全局角色表
CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,  -- 全局唯一
  ...
);

-- 全局角色权限关联表
CREATE TABLE role_permissions (
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  ...
);
```

**问题**:
- ❌ 所有租户共享角色（无法定制）
- ❌ 租户无法创建自定义角色
- ❌ 权限不够灵活

#### 方案2: 租户级RBAC（推荐）

**数据模型**:
```sql
-- 租户级角色表
CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,      -- 租户隔离
  name VARCHAR(100) NOT NULL,      -- 租户内唯一
  is_system BOOLEAN DEFAULT FALSE, -- 是否系统预定义角色
  ...
);

-- 租户级角色权限关联表
CREATE TABLE role_feature_permissions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,        -- 租户隔离
  role_id BIGINT NOT NULL,
  feature_permission_id BIGINT NOT NULL,
  ...
);
```

**优势**:
- ✅ 每个租户独立管理角色
- ✅ 支持租户自定义角色
- ✅ 权限完全隔离

#### 方案3: 混合模式

**数据模型**:
```sql
-- 预定义角色（全局）
CREATE TABLE system_roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  ...
);

-- 自定义角色（租户级）
CREATE TABLE custom_roles (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  name VARCHAR(100) NOT NULL,
  ...
);
```

**问题**:
- ❌ 数据模型复杂
- ❌ 权限验证逻辑复杂
- ❌ 维护成本高

### 预定义角色设计

每个租户注册时，自动创建3个预定义角色：

#### 1. 系统管理员（SYSTEM_ADMIN）

**特征**:
- `is_system = true`（不可删除）
- 拥有所有功能权限
- 每个租户独立创建

**权限**:
```json
{
  "features": [
    {
      "code": "SYSTEM_CONFIG",
      "actions": ["VIEW", "EDIT"]
    },
    {
      "code": "ORGANIZATION_MANAGEMENT",
      "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
    },
    {
      "code": "USER_MANAGEMENT",
      "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
    },
    {
      "code": "ROLE_MANAGEMENT",
      "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
    },
    ...
  ]
}
```

#### 2. 组织管理员（ORGANIZATION_ADMIN）

**特征**:
- `is_system = false`（可编辑、可删除）
- 管理组织内部用户和角色

**权限**:
```json
{
  "features": [
    {
      "code": "USER_MANAGEMENT",
      "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
    },
    {
      "code": "ROLE_MANAGEMENT",
      "actions": ["VIEW", "CREATE", "EDIT"]
    },
    {
      "code": "DEVICE_MANAGEMENT",
      "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
    },
    ...
  ]
}
```

#### 3. 普通用户（NORMAL_USER）

**特征**:
- `is_system = false`（可编辑、可删除）
- 基础查看权限

**权限**:
```json
{
  "features": [
    {
      "code": "DEVICE_MANAGEMENT",
      "actions": ["VIEW"]
    },
    {
      "code": "DATA_VIEW",
      "actions": ["VIEW"]
    },
    {
      "code": "ALERT_MANAGEMENT",
      "actions": ["VIEW"]
    }
  ]
}
```

### 权限验证流程

```
┌──────────┐                  ┌──────────┐                  ┌──────────┐
│  Client  │                  │   API    │                  │  Redis   │
└────┬─────┘                  └────┬─────┘                  └────┬─────┘
     │                             │                             │
     │  1. 请求API (携带Token)     │                             │
     │────────────────────────────>│                             │
     │                             │  2. 提取user_id              │
     │                             │  3. 查询权限缓存             │
     │                             │  key: user_permissions:{user_id}│
     │                             │────────────────────────────>│
     │                             │                             │
     │                             │  4. 返回权限数据             │
     │                             │<────────────────────────────│
     │                             │                             │
     │                             │  5. 验证权限                 │
     │                             │  if (hasPermission(         │
     │                             │    feature: "USER_MANAGEMENT",│
     │                             │    action: "DELETE"          │
     │                             │  )) {                       │
     │                             │    next()                   │
     │                             │  } else {                   │
     │                             │    return 403               │
     │                             │  }                          │
     │                             │                             │
     │  6. 返回数据/403            │                             │
     │<────────────────────────────│                             │
```

**权限验证代码**:
```go
// PermissionMiddleware 权限验证中间件
func PermissionMiddleware(featureCode, actionCode string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取当前用户
        userID := c.GetString("user_id")

        // 2. 从Redis缓存获取权限
        permissions, err := getPermissionsFromCache(userID)
        if err != nil {
            // 缓存未命中，从数据库加载
            permissions, err = loadPermissionsFromDB(userID)
            if err != nil {
                c.JSON(500, gin.H{"error": "Failed to load permissions"})
                c.Abort()
                return
            }
            // 写入缓存
            setPermissionsToCache(userID, permissions, 1*time.Hour)
        }

        // 3. 验证权限
        if !hasPermission(permissions, featureCode, actionCode) {
            c.JSON(403, gin.H{"error": "Permission denied"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// hasPermission 检查权限
func hasPermission(permissions []Permission, featureCode, actionCode string) bool {
    for _, permission := range permissions {
        if permission.FeatureCode == featureCode {
            for _, action := range permission.Actions {
                if action == actionCode {
                    return true
                }
            }
        }
    }
    return false
}
```

### 性能优化

**Redis缓存设计**:
```redis
# 用户权限缓存
key: user_permissions:{user_id}
value: {
  "features": [
    {
      "code": "USER_MANAGEMENT",
      "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
    },
    ...
  ]
}
ttl: 3600 (1小时)

# 缓存更新策略:
# - 用户角色变更时删除缓存
# - 权限配置变更时删除相关用户缓存
```

**性能指标**:
- 权限验证响应时间: < 10ms（Redis缓存）
- 权限加载响应时间: < 100ms（数据库查询）
- 缓存命中率: > 95%

## 后果（Consequences）

### 正面影响

1. **灵活性**
   - 租户可以自定义角色
   - 支持复杂的权限组合
   - 权限完全隔离

2. **可扩展性**
   - 新增功能模块无需修改权限模型
   - 新增操作权限灵活

3. **安全性**
   - 租户间权限完全隔离
   - 细粒度权限控制

### 负面影响

1. **复杂度**
   - 数据模型稍复杂
   - 权限验证逻辑复杂

2. **性能**
   - 需要缓存优化（Redis）
   - 权限检查增加响应时间

### 缓解措施

1. **封装**
   - 封装权限检查中间件
   - 提供前端SDK（MobX Store）

2. **缓存**
   - Redis缓存用户权限
   - 缓存失效策略

3. **文档**
   - 编写权限管理文档
   - 提供代码示例

## 实施方案

### 数据库表设计

详见：`03-design/database-design/cloud-table-design-specs.md`

**核心表**:
1. `roles` - 角色表
2. `feature_permissions` - 功能权限表
3. `action_permissions` - 操作权限表
4. `role_feature_permissions` - 角色功能权限关联表
5. `feature_action_permissions` - 功能操作权限关联表
6. `user_roles` - 用户角色关联表

### API设计

**权限检查接口**:
```http
POST /api/v1/permissions/check
Content-Type: application/json

{
  "feature": "USER_MANAGEMENT",
  "action": "DELETE"
}

Response:
{
  "code": 200,
  "data": {
    "allowed": true
  }
}
```

### 前端集成

**MobX Store**:
```typescript
class PermissionStore {
  // Observable State
  featurePermissions: Map<string, FeaturePermission> = new Map()

  // Computed
  get hasPermission(): (feature: string, action: string) => boolean {
    return (feature: string, action: string) => {
      const fp = this.featurePermissions.get(feature)
      return fp?.actions.includes(action) || false
    }
  }

  // Action
  async loadPermissions(userId: string) {
    const permissions = await api.getUserPermissions(userId)
    this.featurePermissions = new Map(
      permissions.features.map(f => [f.code, f])
    )
  }
}
```

**React组件使用**:
```tsx
import { usePermissionStore } from '@/stores/permissionStore'

const UserListPage: React.FC = () => {
  const { hasPermission } = usePermissionStore()

  return (
    <div>
      {/* 只有DELETE权限才显示删除按钮 */}
      {hasPermission('USER_MANAGEMENT', 'DELETE') && (
        <Button danger onClick={handleDelete}>
          删除
        </Button>
      )}
    </div>
  )
}
```

## 测试策略

### 单元测试

```go
func TestPermissionCheck(t *testing.T) {
    // Setup
    role := createSystemRole()
    user := createUserWithRole(role)

    // Test 1: 系统管理员拥有所有权限
    permissions := loadUserPermissions(user.ID)
    assert.True(t, hasPermission(permissions, "USER_MANAGEMENT", "DELETE"))
    assert.True(t, hasPermission(permissions, "DEVICE_MANAGEMENT", "CREATE"))

    // Test 2: 普通用户只有查看权限
    normalUserRole := createNormalUserRole()
    user2 := createUserWithRole(normalUserRole)
    permissions2 := loadUserPermissions(user2.ID)
    assert.True(t, hasPermission(permissions2, "DEVICE_MANAGEMENT", "VIEW"))
    assert.False(t, hasPermission(permissions2, "DEVICE_MANAGEMENT", "DELETE"))
}
```

## 参考资料

1. **需求文档**: FE-007-03
2. **技术方案**: SOL-007
3. **RBAC最佳实践**: https://en.wikipedia.org/wiki/Role-based_access_control

## 变更历史

| 日期 | 版本 | 变更内容 | 变更人 |
|------|------|---------|--------|
| 2026-01-27 | 1.0 | 初始创建 | Claude Code |
