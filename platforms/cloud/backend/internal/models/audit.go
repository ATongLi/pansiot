package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// AuditLog 审计日志表
type AuditLog struct {
	ID           int64              `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantID     int64              `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_id" json:"tenant_id"`
	ModuleCode   string             `gorm:"column:module_code;type:varchar(50);index" json:"module_code"`   // 模块代码
	ActionType   string             `gorm:"column:action_type;type:varchar(20);index" json:"action_type"` // 操作类型：CREATE, UPDATE, DELETE
	EntityType   string             `gorm:"column:entity_type;type:varchar(50);index" json:"entity_type"` // 实体类型
	EntityID     *int64             `gorm:"column:entity_id;type:bigint" json:"entity_id,omitempty"`              // 实体ID
	ActionDetail ActionDetailJSON   `gorm:"column:action_detail;type:jsonb;not null" json:"action_detail"`         // 操作详情（JSON）
	OperatorID   int64              `gorm:"column:operator_id;type:bigint;not null;index:idx_operator_id" json:"operator_id"`
	OperatorName string             `gorm:"column:operator_name;type:varchar(100)" json:"operator_name"`
	IPAddress    string             `gorm:"column:ip_address;type:varchar(45)" json:"ip_address"`
	UserAgent    string             `gorm:"column:user_agent;type:varchar(500)" json:"user_agent"`
	Status       string             `gorm:"column:status;type:varchar(20);not null" json:"status"` // success, failed
	CreatedAt    time.Time          `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// ActionDetailJSON 操作详情JSON类型
type ActionDetailJSON struct {
	Before  interface{} `json:"before"`
	After   interface{} `json:"after"`
	Changes []Change     `json:"changes,omitempty"`
}

// Change 变更记录
type Change struct {
	Field string      `json:"field"`
	Old   interface{} `json:"old"`
	New   interface{} `json:"new"`
}

// Scan 实现sql.Scanner接口
func (j *ActionDetailJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j ActionDetailJSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// AuditLogActionType 操作类型常量
const (
	AuditActionCreate = "CREATE"
	AuditActionUpdate = "UPDATE"
	AuditActionDelete = "DELETE"
	AuditActionLogin  = "LOGIN"
	AuditActionLogout = "LOGOUT"
	AuditActionExport = "EXPORT"
	AuditActionImport = "IMPORT"
)

// AuditLogStatus 状态常量
const (
	AuditStatusSuccess = "success"
	AuditStatusFailed  = "failed"
)
