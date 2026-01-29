-- PansIot Cloud Platform - 数据库迁移脚本
-- 版本: 1.0.0
-- 日期: 2026-01-28
-- 说明: 云平台账号系统所有表结构

-- =====================================================
-- 1. 序列和函数
-- =====================================================

-- 企业序列号序列
CREATE SEQUENCE IF NOT EXISTS serial_number_seq START 1;

-- 企业序列号生成函数
CREATE OR REPLACE FUNCTION generate_serial_number()
RETURNS VARCHAR(8) AS $$
DECLARE
  prefix VARCHAR(4);
  suffix INT;
  serial_number VARCHAR(8);
  max_attempts INT := 10;
  attempts INT := 0;
BEGIN
  LOOP
    attempts := attempts + 1;
    IF attempts > max_attempts THEN
      RAISE EXCEPTION 'Failed to generate unique serial number after % attempts', max_attempts;
    END IF;

    -- 1. 生成4位随机字符（大小写字母+数字）
    prefix := upper(substring(encode(
      gen_random_bytes(3),
      'base64'
    ), 1, 4));

    -- 移除非字母数字字符
    prefix := regexp_replace(prefix, '[^A-Z0-9]', '', 'g');

    -- 确保长度为4
    WHILE length(prefix) < 4 LOOP
      prefix := prefix || substr('ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', ceil(random() * 36)::INT, 1);
    END LOOP;

    -- 2. 获取下一个自增ID
    suffix := nextval('serial_number_seq');

    -- 3. 拼接（4位随机+4位自增）
    serial_number := prefix || LPAD(suffix::TEXT, 4, '0');

    -- 4. 检查唯一性
    IF NOT EXISTS(SELECT 1 FROM tenants WHERE serial_number = serial_number) THEN
      RETURN serial_number;
    END IF;
  END LOOP;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- 2. 租户和用户相关表
-- =====================================================

-- 租户表
CREATE TABLE IF NOT EXISTS tenants (
  id BIGSERIAL PRIMARY KEY,
  serial_number VARCHAR(8) NOT NULL DEFAULT generate_serial_number() UNIQUE,
  name VARCHAR(200) NOT NULL,
  tenant_type VARCHAR(20) NOT NULL DEFAULT 'TERMINAL',  -- INTEGRATOR, TERMINAL
  industry VARCHAR(100) DEFAULT '其他',
  contact_person VARCHAR(100),
  contact_phone VARCHAR(20),
  contact_email VARCHAR(100),
  parent_tenant_id BIGINT,
  status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',  -- ACTIVE, SUSPENDED, DELETED
  expire_date TIMESTAMP,
  max_sub_tenants INT NOT NULL DEFAULT 0,
  max_users INT NOT NULL DEFAULT 0,
  max_devices INT NOT NULL DEFAULT 0,
  max_storage_gb INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- 租户表索引
CREATE INDEX idx_tenants_tenant_type ON tenants(tenant_type);
CREATE INDEX idx_tenants_parent_tenant_id ON tenants(parent_tenant_id);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_serial_number ON tenants(serial_number);

-- 用户表
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL,
  phone VARCHAR(20),
  phone_country_code VARCHAR(5) DEFAULT '+86',
  password_hash VARCHAR(255) NOT NULL,
  real_name VARCHAR(100),
  avatar VARCHAR(500),
  status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',  -- ACTIVE, SUSPENDED, DELETED
  last_login_at TIMESTAMP,
  last_login_ip VARCHAR(45),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  UNIQUE(tenant_id, username)
);

-- 用户表索引
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  role_code VARCHAR(50) NOT NULL,
  role_name VARCHAR(100) NOT NULL,
  description VARCHAR(500),
  is_system BOOLEAN NOT NULL DEFAULT false,
  is_deletable BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  UNIQUE(tenant_id, role_code)
);

-- 角色表索引
CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
CREATE INDEX idx_roles_role_code ON roles(role_code);
CREATE INDEX idx_roles_deleted_at ON roles(deleted_at);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, role_id)
);

-- 用户角色关联表索引
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- =====================================================
-- 3. 权限相关表
-- =====================================================

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  feature_code VARCHAR(50) NOT NULL,  -- SYSTEM_CONFIG, USER_MANAGEMENT, etc.
  action_code VARCHAR(20) NOT NULL,    -- VIEW, CREATE, EDIT, DELETE, etc.
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(tenant_id, feature_code, action_code)
);

-- 权限表索引
CREATE INDEX idx_permissions_tenant_id ON permissions(tenant_id);
CREATE INDEX idx_permissions_feature_code ON permissions(feature_code);
CREATE INDEX idx_permissions_action_code ON permissions(action_code);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
  id BIGSERIAL PRIMARY KEY,
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(role_id, permission_id)
);

-- 角色权限关联表索引
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- =====================================================
-- 4. 功能模块和配额相关表
-- =====================================================

-- 功能模块表
CREATE TABLE IF NOT EXISTS feature_modules (
  id BIGSERIAL PRIMARY KEY,
  module_code VARCHAR(50) NOT NULL UNIQUE,  -- DEVICE_MANAGEMENT, DATA_VIEW, etc.
  module_name VARCHAR(100) NOT NULL,
  description VARCHAR(500),
  enabled BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 租户功能开通表
CREATE TABLE IF NOT EXISTS tenant_features (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  module_code VARCHAR(50) NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT true,
  expires_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(tenant_id, module_code)
);

-- 租户功能开通表索引
CREATE INDEX idx_tenant_features_tenant_id ON tenant_features(tenant_id);
CREATE INDEX idx_tenant_features_module_code ON tenant_features(module_code);

-- 租户配额表
CREATE TABLE IF NOT EXISTS tenant_quotas (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  quota_type VARCHAR(50) NOT NULL,  -- sub_tenants, users, devices, storage_gb
  total_quota INT NOT NULL DEFAULT 0,
  used_quota INT NOT NULL DEFAULT 0,
  remaining_quota INT NOT NULL DEFAULT 0,
  expires_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(tenant_id, quota_type)
);

-- 租户配额表索引
CREATE INDEX idx_tenant_quotas_tenant_id ON tenant_quotas(tenant_id);
CREATE INDEX idx_tenant_quotas_quota_type ON tenant_quotas(quota_type);

-- 配额分配表（集成商为下游分配配额）
CREATE TABLE IF NOT EXISTS quota_allocations (
  id BIGSERIAL PRIMARY KEY,
  parent_tenant_id BIGINT NOT NULL,  -- 集成商ID
  child_tenant_id BIGINT NOT NULL,   -- 下游客户ID
  quota_type VARCHAR(50) NOT NULL,
  allocated_quota INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(parent_tenant_id, child_tenant_id, quota_type)
);

-- 配额分配表索引
CREATE INDEX idx_quota_allocations_parent_tenant_id ON quota_allocations(parent_tenant_id);
CREATE INDEX idx_quota_allocations_child_tenant_id ON quota_allocations(child_tenant_id);

-- =====================================================
-- 5. 审计日志表（分区表）
-- =====================================================

-- 审计日志主表
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  module_code VARCHAR(50),
  action_type VARCHAR(20),  -- CREATE, UPDATE, DELETE, LOGIN, LOGOUT
  entity_type VARCHAR(50),
  entity_id BIGINT,
  action_detail JSONB NOT NULL,  -- {before: {}, after: {}, changes: []}
  operator_id BIGINT NOT NULL,
  operator_name VARCHAR(100),
  ip_address VARCHAR(45),
  user_agent VARCHAR(500),
  status VARCHAR(20) NOT NULL,  -- success, failed
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- 审计日志索引
CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_module_code ON audit_logs(module_code);
CREATE INDEX idx_audit_logs_action_type ON audit_logs(action_type);
CREATE INDEX idx_audit_logs_operator_id ON audit_logs(operator_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_action_detail ON audit_logs USING GIN (action_detail);

-- 创建审计日志分区（按月）
-- 注意：需要为每个月创建分区，这里只创建当前月份的分区
-- 实际使用时需要定期创建新分区
CREATE TABLE IF NOT EXISTS audit_logs_2026_01 PARTITION OF audit_logs
  FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

-- =====================================================
-- 6. 设备相关表（后续功能使用）
-- =====================================================

-- 设备表
CREATE TABLE IF NOT EXISTS devices (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  managed_tenant_id BIGINT,  -- 管理租户（集成商）
  device_code VARCHAR(50) NOT NULL,
  device_name VARCHAR(100) NOT NULL,
  device_type VARCHAR(50) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'offline',  -- online, offline, error
  last_online_at TIMESTAMP,
  firmware_version VARCHAR(50),
  description VARCHAR(500),
  location VARCHAR(200),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  UNIQUE(tenant_id, device_code)
);

-- 设备表索引
CREATE INDEX idx_devices_tenant_id ON devices(tenant_id);
CREATE INDEX idx_devices_managed_tenant_id ON devices(managed_tenant_id);
CREATE INDEX idx_devices_device_code ON devices(device_code);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_deleted_at ON devices(deleted_at);

-- =====================================================
-- 7. 表注释
-- =====================================================

COMMENT ON TABLE tenants IS '租户表';
COMMENT ON COLUMN tenants.serial_number IS '企业序列号（8位：4位随机+4位自增）';
COMMENT ON COLUMN tenants.tenant_type IS '租户类型：INTEGRATOR(集成商)、TERMINAL(下游客户)';
COMMENT ON COLUMN tenants.parent_tenant_id IS '上级租户ID（用于组织树）';

COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.tenant_id IS '归属租户';
COMMENT ON COLUMN users.phone_country_code IS '手机号国家码（支持国际化）';

COMMENT ON TABLE roles IS '角色表';
COMMENT ON COLUMN roles.is_system IS '是否系统角色（系统角色不可删除）';

COMMENT ON TABLE permissions IS '权限表';
COMMENT ON COLUMN permissions.feature_code IS '功能代码';
COMMENT ON COLUMN permissions.action_code IS '操作代码';

COMMENT ON TABLE audit_logs IS '审计日志表';
COMMENT ON COLUMN audit_logs.action_detail IS '操作详情（JSON格式：before/after/changes）';

COMMENT ON TABLE tenant_features IS '租户功能开通表';
COMMENT ON TABLE tenant_quotas IS '租户配额表';
COMMENT ON TABLE quota_allocations IS '配额分配表（集成商为下游分配）';

COMMENT ON TABLE devices IS '设备表';
COMMENT ON COLUMN devices.managed_tenant_id IS '管理租户（集成商可以看到所有下游设备）';

-- =====================================================
-- 8. 初始化数据
-- =====================================================

-- 插入默认功能模块
INSERT INTO feature_modules (module_code, module_name, description, enabled) VALUES
  ('SYSTEM_CONFIG', '系统配置', '系统级配置管理', true),
  ('ORGANIZATION_MANAGEMENT', '组织管理', '子组织管理功能', true),
  ('USER_MANAGEMENT', '用户管理', '用户账号管理', true),
  ('ROLE_MANAGEMENT', '角色管理', '角色和权限管理', true),
  ('DEVICE_MANAGEMENT', '设备管理', '设备接入和管理', true),
  ('DATA_VIEW', '数据查看', '数据查看和导出', true),
  ('ALERT_MANAGEMENT', '告警管理', '告警规则和通知', true),
  ('QUOTA_MANAGEMENT', '配额管理', '功能配额分配', true),
  ('AUDIT_LOG_VIEW', '审计日志', '操作日志查看', true)
ON CONFLICT (module_code) DO NOTHING;

-- =====================================================
-- 9. 完成标记
-- =====================================================

-- 创建迁移记录表（用于记录迁移历史）
CREATE TABLE IF NOT EXISTS schema_migrations (
  id BIGSERIAL PRIMARY KEY,
  version VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(500),
  executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 记录本次迁移
INSERT INTO schema_migrations (version, description) VALUES
  ('1.0.0', '初始数据库结构：租户、用户、角色、权限、审计日志等表')
ON CONFLICT (version) DO NOTHING;
