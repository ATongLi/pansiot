-- 初始化数据库脚本
-- 创建序列号序列
CREATE SEQUENCE IF NOT EXISTS serial_number_seq START 1;

-- 创建企业序列号生成函数
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

-- 创建审核日志分区表（按月）
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  module_code VARCHAR(50),
  action_type VARCHAR(20),
  entity_type VARCHAR(50),
  entity_id BIGINT,
  action_detail JSONB,
  operator_id BIGINT NOT NULL,
  operator_name VARCHAR(100),
  ip_address VARCHAR(45),
  user_agent VARCHAR(500),
  status VARCHAR(20),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- 创建审核日志索引
CREATE INDEX idx_audit_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_module ON audit_logs(module_code);
CREATE INDEX idx_audit_action_type ON audit_logs(action_type);
CREATE INDEX idx_audit_operator_id ON audit_logs(operator_id);
CREATE INDEX idx_audit_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_action_detail ON audit_logs USING GIN (action_detail);

-- 注释
COMMENT ON TABLE audit_logs IS '系统审计日志表';
COMMENT ON COLUMN audit_logs.action_detail IS '操作详情（JSON格式，包含before/after/changes）';
