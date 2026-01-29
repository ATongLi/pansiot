# 云平台业务表设计规范

## 文档信息
- **文档ID**: DB-SPEC-001
- **版本**: v1.0
- **创建日期**: 2026-01-27
- **适用平台**: Cloud (云平台端)
- **来源需求**: REQ-007

## 概述

本文档定义云平台业务表的统一设计规范，确保多租户数据隔离、权限控制和审计追踪的一致性。

**核心设计原则**：
1. **数据隔离优先**：所有业务表必须包含租户隔离字段
2. **审计可追溯**：关键操作必须记录审计日志
3. **字段命名统一**：遵循统一的命名规范
4. **性能优化**：合理的索引设计
5. **扩展性**：预留扩展字段

---

## 1. 租户隔离字段规范

### 1.1 双字段设计（必选）

所有业务表（设备、用户、告警、日志等）**必须**包含以下两个核心字段：

#### `tenant_id`（归属租户ID）
- **类型**: BIGINT
- **约束**: NOT NULL
- **说明**: 标识数据的归属主体
- **用途**:
  - 数据创建时自动设置为当前用户所属租户ID
  - 用于租户数据隔离（下游客户只能看到自己的数据）
- **索引**: 必须创建索引 `idx_tenant_id`

#### `managed_tenant_id`（管理租户ID）
- **类型**: BIGINT
- **约束**: NULL
- **说明**: 标识数据的管控集成商
- **用途**:
  - 数据创建时自动设置为当前租户的managed_tenant_id
  - 集成商可以看到managed_tenant_id指向自己的所有下游数据
  - 集成商租户创建的数据，此字段为NULL
- **索引**: 必须创建索引 `idx_managed_tenant_id`

### 1.2 字段组合索引

建议创建复合索引以优化查询性能：
```sql
CREATE INDEX idx_tenant_managed ON {table_name}(tenant_id, managed_tenant_id);
```

### 1.3 字段填充规则

| 创建者 | tenant_id | managed_tenant_id |
|--------|-----------|-------------------|
| 集成商为自己创建 | 集成商自己的ID | NULL |
| 下游客户为自己创建 | 下游客户自己的ID | 下溯客户的managed_tenant_id（指向集成商） |
| 集成商为下游创建 | 下游租户的ID | 集成商自己的ID |

### 1.4 适用范围

**必须包含双字段的表**：
- `devices`（设备表）
- `users`（用户表）
- `alerts`（告警表）
- `logs`（日志表）
- `data_points`（数据点表）
- `dashboards`（仪表板表）
- `reports`（报表表）
- `projects`（工程表）
- `web_pages`（Web页面表）
- 任何其他业务对象表

**不包含双字段的表**：
- 系统配置表（全局共享）
- 权限定义表（全局共享）
- 字典表（全局共享）
- `tenants`（租户表本身，仅包含managed_tenant_id）

---

## 2. 标准字段规范

### 2.1 主键字段

所有表**必须**包含主键字段：

```sql
-- PostgreSQL语法（推荐）
id BIGSERIAL PRIMARY KEY

-- 或使用UUID
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
```

**推荐**: 对于高并发表，使用BIGINT自增ID；对于分布式系统，使用UUID。

### 2.2 时间戳字段（必选）

所有表**必须**包含以下时间戳字段：

#### `created_at`（创建时间）
```sql
-- PostgreSQL语法
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
```

#### `updated_at`（更新时间）
```sql
-- PostgreSQL语法（需要触发器或函数自动更新）
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
```

**说明**: PostgreSQL不支持`ON UPDATE CURRENT_TIMESTAMP`，需要通过触发器或GORM钩子实现自动更新。

#### `deleted_at`（软删除时间，可选）
```sql
deleted_at TIMESTAMP NULL
```

**说明**:
- 使用软删除机制，数据不物理删除
- 查询时添加 `WHERE deleted_at IS NULL` 条件
- 索引：`CREATE INDEX idx_deleted_at ON {table_name}(deleted_at)`

### 2.3 操作人字段（推荐）

关键业务表**建议**包含操作人字段：

#### `created_by`（创建人ID）
```sql
created_by BIGINT NOT NULL
-- 索引: idx_created_by
```

#### `updated_by`（更新人ID）
```sql
updated_by BIGINT
-- 索引: idx_updated_by
```

---

## 3. 审计日志字段规范

### 3.1 审计日志表设计

对于需要记录详细变更的表，建议使用审计日志表（而非在原表添加字段）：

```sql
-- PostgreSQL语法
CREATE TABLE audit_logs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  module_code VARCHAR(50),
  action_type VARCHAR(20) NOT NULL,
  entity_type VARCHAR(50) NOT NULL,
  entity_id BIGINT,
  action_detail JSONB,
  operator_id BIGINT NOT NULL,
  operator_name VARCHAR(100),
  ip_address VARCHAR(45),
  user_agent VARCHAR(500),
  status VARCHAR(20) NOT NULL,
  error_message TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 添加索引
CREATE INDEX idx_audit_tenant_module ON audit_logs(tenant_id, module_code);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_operator ON audit_logs(operator_id);
CREATE INDEX idx_audit_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_action ON audit_logs(action_type);

-- 添加注释
COMMENT ON COLUMN audit_logs.module_code IS '模块代码';
COMMENT ON COLUMN audit_logs.action_type IS '操作类型：CREATE/UPDATE/DELETE/LOGIN等';
COMMENT ON COLUMN audit_logs.entity_type IS '实体类型';
COMMENT ON COLUMN audit_logs.entity_id IS '实体ID';
COMMENT ON COLUMN audit_logs.action_detail IS '操作详情（before/after/changes）';
COMMENT ON COLUMN audit_logs.operator_id IS '操作人ID';
COMMENT ON COLUMN audit_logs.operator_name IS '操作人姓名（冗余字段）';
COMMENT ON COLUMN audit_logs.ip_address IS '操作IP';
COMMENT ON COLUMN audit_logs.user_agent IS '用户代理';
COMMENT ON COLUMN audit_logs.status IS '操作状态：SUCCESS/FAILED/PARTIAL';
COMMENT ON COLUMN audit_logs.error_message IS '错误信息';
```

### 3.2 action_detail 字段规范

使用JSON格式记录字段级变更：

**CREATE操作**：
```json
{
  "before": null,
  "after": {
    "device_id": 12345,
    "device_name": "温度传感器01",
    "tenant_id": 1001
  },
  "changes": null
}
```

**UPDATE操作**：
```json
{
  "before": {
    "device_name": "旧名称",
    "status": "offline"
  },
  "after": {
    "device_name": "新名称",
    "status": "online"
  },
  "changes": [
    {"field": "device_name", "old": "旧名称", "new": "新名称"},
    {"field": "status", "old": "offline", "new": "online"}
  ]
}
```

**DELETE操作**：
```json
{
  "before": {
    "device_id": 12345,
    "device_name": "温度传感器01"
  },
  "after": null,
  "changes": null
}
```

---

## 4. 字段命名规范

### 4.1 命名规则

- **使用小写字母**: 所有字段名使用小写
- **单词间用下划线分隔**: `created_at` 而不是 `createdAt`
- **使用英文单词**: 避免拼音或缩写
- **布尔字段**: 使用 `is_` 前缀，如 `is_active`, `is_deleted`
- **状态字段**: 使用 `_status` 后缀，如 `order_status`
- **类型字段**: 使用 `_type` 后缀，如 `user_type`
- **ID字段**: 使用 `_id` 后缀，如 `tenant_id`

### 4.2 数据类型规范

| 数据类型 | 用途 | 示例 |
|---------|------|------|
| BIGINT | 主键、外键、ID字段 | `id`, `tenant_id` |
| INT | 计数、数量 | `device_count`, `user_count` |
| VARCHAR(n) | 短字符串、名称 | `name`, `email` |
| TEXT | 长文本、描述 | `description`, `content` |
| JSON | 结构化数据 | `action_detail`, `config` |
| TIMESTAMP | 日期时间 | `created_at`, `updated_at` |
| DECIMAL(m, d) | 金额、精确小数 | `price`, `quota` |
| TINYINT/BOOLEAN | 布尔值、状态 | `is_active`, `is_deleted` |
| ENUM | 枚举类型 | `tenant_type`, `status` |

### 4.3 字段长度规范

| 字段类型 | 推荐长度 | 说明 |
|---------|---------|------|
| 企业名称 | VARCHAR(200) | |
| 用户名 | VARCHAR(50) | |
| 邮箱 | VARCHAR(100) | |
| 手机号 | VARCHAR(20) | 支持国际区号 |
| 密码哈希 | VARCHAR(255) | bcrypt等算法 |
| 序列号 | VARCHAR(8) | 企业序列号（8位）、设备序列号 |
| URL | VARCHAR(500) | |
| IP地址 | VARCHAR(45) | 支持IPv6 |
| 描述 | TEXT | 无长度限制 |

---

## 5. 索引设计规范

### 5.1 索引命名规范

- 主键索引：`PRIMARY KEY`
- 普通索引：`idx_字段名`，如 `idx_tenant_id`
- 唯一索引：`uk_字段名`，如 `uk_email`
- 全文索引：`ft_字段名`，如 `ft_content`
- 复合索引：`idx_字段1_字段2`，如 `idx_tenant_type`

### 5.2 索引设计原则

**必须创建索引的字段**：
- 主键字段
- 外键字段（tenant_id, managed_tenant_id等）
- 经常查询的字段（status, type等）
- 经常排序的字段（created_at等）
- 唯一性约束字段（email, username等）

**复合索引设计原则**：
- 最左前缀原则
- 区分度高的字段放前面
- 覆盖常用查询条件

### 5.3 索引示例

```sql
-- PostgreSQL语法 - 租户隔离表的标准索引
CREATE TABLE devices (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  managed_tenant_id BIGINT,
  device_name VARCHAR(100) NOT NULL,
  device_type VARCHAR(50) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'offline',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NOT NULL
);

-- 单列索引
CREATE INDEX idx_devices_tenant_id ON devices(tenant_id);
CREATE INDEX idx_devices_managed_tenant_id ON devices(managed_tenant_id);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_deleted_at ON devices(deleted_at);
CREATE INDEX idx_devices_created_by ON devices(created_by);
CREATE INDEX idx_devices_created_at ON devices(created_at);

-- 复合索引
CREATE INDEX idx_devices_tenant_status ON devices(tenant_id, status);
CREATE INDEX idx_devices_tenant_type ON devices(tenant_id, device_type);
CREATE INDEX idx_devices_tenant_managed ON devices(tenant_id, managed_tenant_id);

-- 唯一索引
CREATE UNIQUE INDEX uk_devices_device_name ON devices(tenant_id, device_name);
```

---

## 6. 表设计模板

### 6.1 标准业务表模板

```sql
-- PostgreSQL语法 - 标准业务表模板
CREATE TABLE {table_name} (
  -- 主键
  id BIGSERIAL PRIMARY KEY,

  -- 租户隔离字段（必选）
  tenant_id BIGINT NOT NULL,
  managed_tenant_id BIGINT,

  -- 业务字段
  name VARCHAR(200) NOT NULL,
  type VARCHAR(50),
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  description TEXT,

  -- 扩展字段（JSONB，可选）
  extra_config JSONB,

  -- 操作人字段（推荐）
  created_by BIGINT NOT NULL,
  updated_by BIGINT,

  -- 时间戳字段（必选）
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);

-- 索引
CREATE INDEX idx_{table}_tenant_id ON {table_name}(tenant_id);
CREATE INDEX idx_{table}_managed_tenant_id ON {table_name}(managed_tenant_id);
CREATE INDEX idx_{table}_status ON {table_name}(status);
CREATE INDEX idx_{table}_deleted_at ON {table_name}(deleted_at);
CREATE INDEX idx_{table}_created_at ON {table_name}(created_at);
CREATE INDEX idx_{table}_tenant_status ON {table_name}(tenant_id, status);
CREATE UNIQUE INDEX uk_{table}_name ON {table_name}(tenant_id, name);

-- 添加注释
COMMENT ON TABLE {table_name} IS '{table_comment}';
COMMENT ON COLUMN {table_name}.tenant_id IS '归属租户ID';
COMMENT ON COLUMN {table_name}.managed_tenant_id IS '管理租户ID';
COMMENT ON COLUMN {table_name}.name IS '名称';
COMMENT ON COLUMN {table_name}.status IS '状态';
COMMENT ON COLUMN {table_name}.extra_config IS '扩展配置';
```

---

## 7. 特殊表设计规范

### 7.1 配额表设计

配额表用于记录功能模块的资源使用情况：

```sql
-- PostgreSQL语法 - 配额表设计
CREATE TABLE tenant_quota_usage (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  module_code VARCHAR(50) NOT NULL,
  total_quota INT,
  used_quota INT NOT NULL DEFAULT 0,
  remaining_quota INT GENERATED ALWAYS AS (total_quota - used_quota) STORED,
  allocated_quota INT NOT NULL DEFAULT 0,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_quota_tenant_module ON tenant_quota_usage(tenant_id, module_code);
CREATE INDEX idx_quota_remaining ON tenant_quota_usage(remaining_quota);
CREATE UNIQUE INDEX uk_quota_tenant_module ON tenant_quota_usage(tenant_id, module_code);

-- 添加注释
COMMENT ON TABLE tenant_quota_usage IS '租户配额使用统计表';
COMMENT ON COLUMN tenant_quota_usage.module_code IS '功能模块代码';
COMMENT ON COLUMN tenant_quota_usage.total_quota IS '总配额';
COMMENT ON COLUMN tenant_quota_usage.used_quota IS '已使用配额';
COMMENT ON COLUMN tenant_quota_usage.remaining_quota IS '剩余配额';
COMMENT ON COLUMN tenant_quota_usage.allocated_quota IS '已分配给下游的配额（仅集成商）';
```

### 7.2 功能开通表设计

```sql
-- PostgreSQL语法 - 功能开通表设计
CREATE TABLE tenant_features (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  module_code VARCHAR(50) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  granted_by BIGINT,
  granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NULL
);

-- 索引
CREATE INDEX idx_features_tenant_module ON tenant_features(tenant_id, module_code);
CREATE INDEX idx_features_active ON tenant_features(is_active);
CREATE UNIQUE INDEX uk_features_tenant_module ON tenant_features(tenant_id, module_code);

-- 添加注释
COMMENT ON TABLE tenant_features IS '租户功能开通记录表';
COMMENT ON COLUMN tenant_features.module_code IS '功能模块代码';
COMMENT ON COLUMN tenant_features.is_active IS '是否开通';
COMMENT ON COLUMN tenant_features.granted_by IS '开通人ID（平台超管或集成商）';
COMMENT ON COLUMN tenant_features.granted_at IS '开通时间';
COMMENT ON COLUMN tenant_features.expires_at IS '过期时间';
```

### 7.3 组织表设计

```sql
-- PostgreSQL语法 - 租户类型枚举
CREATE TYPE tenant_type_enum AS ENUM ('TERMINAL', 'INTEGRATOR');

-- 组织表设计
CREATE TABLE tenants (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  tenant_type tenant_type_enum NOT NULL DEFAULT 'TERMINAL',
  managed_tenant_id BIGINT,
  parent_tenant_id BIGINT,
  serial_number VARCHAR(8) NOT NULL,
  industry VARCHAR(100),
  integrator_since TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);

-- 索引
CREATE INDEX idx_tenants_serial_number ON tenants(serial_number);
CREATE INDEX idx_tenants_managed_tenant_id ON tenants(managed_tenant_id);
CREATE INDEX idx_tenants_parent_tenant_id ON tenants(parent_tenant_id);
CREATE INDEX idx_tenants_tenant_type ON tenants(tenant_type);
CREATE UNIQUE INDEX uk_tenants_serial_number ON tenants(serial_number);

-- 添加注释
COMMENT ON TABLE tenants IS '租户（组织）表';
COMMENT ON COLUMN tenants.name IS '企业名称';
COMMENT ON COLUMN tenants.tenant_type IS '租户类型';
COMMENT ON COLUMN tenants.managed_tenant_id IS '管理租户ID（集成商）';
COMMENT ON COLUMN tenants.parent_tenant_id IS '父租户ID（直接上级）';
COMMENT ON COLUMN tenants.serial_number IS '企业序列号（8位）';
COMMENT ON COLUMN tenants.industry IS '所属行业';
COMMENT ON COLUMN tenants.integrator_since IS '开通集成商功能时间';
```

---

## 8. 数据库引擎和字符集规范

### 8.1 引擎选择

- **默认引擎**: InnoDB
- **理由**:
  - 支持事务（ACID）
  - 支持外键约束
  - 支持行级锁
  - 支持崩溃恢复

### 8.2 字符集规范

- **默认字符集**: `utf8mb4`
- **排序规则**: `utf8mb4_unicode_ci`
- **理由**:
  - 完整支持UTF-8（包括emoji）
  - 多语言支持
  - 避免字符截断问题

---

## 9. 注释规范

### 9.1 表注释

所有表**必须**添加表注释：

```sql
-- PostgreSQL语法
COMMENT ON TABLE users IS '用户表';
```

### 9.2 字段注释

所有字段**必须**添加字段注释：

```sql
-- PostgreSQL语法
COMMENT ON COLUMN users.tenant_id IS '归属租户ID';
COMMENT ON COLUMN users.managed_tenant_id IS '管理租户ID';
```

### 9.3 索引注释

PostgreSQL不支持直接在CREATE INDEX时添加注释，需要使用COMMENT ON INDEX：

```sql
-- PostgreSQL语法
CREATE INDEX idx_tenant_status ON users(tenant_id, status);
COMMENT ON INDEX idx_tenant_status IS '租户状态查询优化';
```

---

## 10. 性能优化建议

### 10.1 分区表

对于数据量大的表（如日志表），建议使用分区：

```sql
-- PostgreSQL语法 - 分区表设计
-- 1. 创建主表
CREATE TABLE audit_logs (
  id BIGSERIAL,
  tenant_id BIGINT NOT NULL,
  module_code VARCHAR(50),
  action_type VARCHAR(20) NOT NULL,
  entity_type VARCHAR(50) NOT NULL,
  entity_id BIGINT,
  action_detail JSONB,
  operator_id BIGINT NOT NULL,
  operator_name VARCHAR(100),
  ip_address VARCHAR(45),
  user_agent VARCHAR(500),
  status VARCHAR(20) NOT NULL,
  error_message TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- 2. 创建分区
CREATE TABLE audit_logs_2023 PARTITION OF audit_logs
  FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');

CREATE TABLE audit_logs_2024 PARTITION OF audit_logs
  FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE audit_logs_2025 PARTITION OF audit_logs
  FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

-- 3. 创建默认分区（可选）
CREATE TABLE audit_logs_default PARTITION OF audit_logs DEFAULT;
```

### 10.2 读写分离

- 读操作：从库（Slave）
- 写操作：主库（Master）
- 审计日志：可考虑使用NoSQL（如MongoDB）

### 10.3 缓存策略

- 热点数据缓存到Redis
- 配额数据缓存（带过期时间）
- 权限数据缓存

---

## 11. 数据完整性约束

### 11.1 外键约束

```sql
ALTER TABLE devices
ADD CONSTRAINT fk_device_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
```

**注意**: 外键约束可能影响性能，高并发场景可考虑应用层校验。

### 11.2 唯一性约束

```sql
ALTER TABLE users
ADD CONSTRAINT uk_username UNIQUE KEY (username);

ALTER TABLE users
ADD CONSTRAINT uk_email UNIQUE KEY (email);
```

### 11.3 检查约束（MySQL 8.0+）

```sql
ALTER TABLE devices
ADD CONSTRAINT chk_status
CHECK (status IN ('online', 'offline', 'error'));
```

---

## 12. 迁移和版本控制

### 12.1 数据库迁移脚本

使用版本化的迁移脚本：

```
migrations/
  ├── V1.0.0__init_schema.sql
  ├── V1.0.1__add_device_table.sql
  ├── V1.0.2__add_audit_log_table.sql
  └── V1.0.3__add_quota_tables.sql
```

### 12.2 回滚脚本

每个迁移脚本**必须**包含回滚脚本：

```
migrations/
  ├── V1.0.3__add_quota_tables.sql
  └── V1.0.3__rollback_quota_tables.sql
```

---

## 13. 安全规范

### 13.1 敏感数据加密

- **密码**: 使用bcrypt、Argon2等单向加密
- **手机号**: 考虑脱敏存储（如138****1234）
- **API密钥**: 使用AES加密存储

### 13.2 SQL注入防护

- 使用参数化查询（PreparedStatement）
- 避免字符串拼接SQL
- 使用ORM框架（如GORM、Ent）

### 13.3 数据备份

- 定期全量备份（每日）
- 增量备份（每小时）
- 异地备份
- 备份恢复测试

---

## 14. 监控和告警

### 14.1 慢查询监控

```sql
-- 启用慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1; -- 超过1秒的查询
```

### 14.2 表空间监控

- 监控表大小增长
- 监控索引使用率
- 监控碎片率

### 14.3 连接数监控

- 监控活跃连接数
- 监控连接池使用情况
- 设置最大连接数限制

---

## 15. 常见问题FAQ

### Q1: 何时使用软删除？

**A**: 关键业务数据（用户、订单、设备）使用软删除；日志类数据可物理删除。

### Q2: 如何选择BIGINT vs INT？

**A**: 主键、外键使用BIGINT；状态码、小计数使用INT。

### Q3: JSON字段性能如何？

**A**: MySQL 5.7+的JSON字段性能良好，但避免频繁更新大JSON。

### Q4: 索引越多越好吗？

**A**: 不是。索引会占用空间并影响写入性能，只对必要字段建索引。

### Q5: 如何处理分布式ID？

**A**: 使用UUID、Snowflake算法或数据库序列。

---

## 附录

### A. 快速检查清单

设计新表时，请检查以下项：

- [ ] 包含租户隔离字段（tenant_id, managed_tenant_id）
- [ ] 包含主键字段（id）
- [ ] 包含时间戳字段（created_at, updated_at）
- [ ] 包含软删除字段（deleted_at）
- [ ] 包含表注释和字段注释
- [ ] 创建必要的索引
- [ ] 设置合适的数据类型和长度
- [ ] 考虑外键约束
- [ ] 考虑唯一性约束
- [ ] 使用utf8mb4字符集

### B. 相关文档

- **需求文档**: REQ-007 (云平台账号系统)
- **功能需求**: FE-007-01 ~ FE-007-09
- **技术方案**: SOL-007 (待创建)
- **API文档**: /api-docs (待创建)

### C. 变更历史

| 日期 | 版本 | 变更内容 | 变更人 |
|------|------|---------|--------|
| 2026-01-27 | 1.0 | 初始创建 | Claude Code |

---

**文档结束**
