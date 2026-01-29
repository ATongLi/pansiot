-- PansIot Cloud Platform - 数据库回滚脚本
-- 版本: 1.0.0
-- 说明: 回滚初始数据库结构

-- 注意：此脚本将删除所有表和数据，请谨慎使用！

-- =====================================================
-- 1. 删除表（按依赖关系逆序）
-- =====================================================

DROP TABLE IF EXISTS audit_logs_2026_01 CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS devices CASCADE;
DROP TABLE IF EXISTS quota_allocations CASCADE;
DROP TABLE IF EXISTS tenant_quotas CASCADE;
DROP TABLE IF EXISTS tenant_features CASCADE;
DROP TABLE IF EXISTS feature_modules CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS tenants CASCADE;

-- =====================================================
-- 2. 删除序列和函数
-- =====================================================

DROP SEQUENCE IF EXISTS serial_number_seq CASCADE;
DROP FUNCTION IF EXISTS generate_serial_number() CASCADE;

-- =====================================================
-- 3. 删除迁移记录
-- =====================================================

DROP TABLE IF EXISTS schema_migrations CASCADE;
