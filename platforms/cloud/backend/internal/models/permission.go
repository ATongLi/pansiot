package models

import (
	"time"
)

// Permission 权限表
type Permission struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantID    int64     `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_id" json:"tenant_id"` // 归属租户
	FeatureCode string    `gorm:"column:feature_code;type:varchar(50);not null;index:idx_feature_code" json:"feature_code"` // 功能代码
	ActionCode  string    `gorm:"column:action_code;type:varchar(20);not null;index:idx_action_code" json:"action_code"`     // 操作代码
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID       int64     `gorm:"column:role_id;type:bigint;not null;index:idx_role_id,index:idx_role_permission" json:"role_id"`
	PermissionID int64     `gorm:"column:permission_id;type:bigint;not null;index:idx_permission_id,index:idx_role_permission" json:"permission_id"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// 功能权限常量
const (
	FeatureSystemConfig       = "SYSTEM_CONFIG"
	FeatureOrganizationMgmt   = "ORGANIZATION_MANAGEMENT"
	FeatureUserMgmt           = "USER_MANAGEMENT"
	FeatureRoleMgmt           = "ROLE_MANAGEMENT"
	FeatureDeviceMgmt         = "DEVICE_MANAGEMENT"
	FeatureDataView           = "DATA_VIEW"
	FeatureAlertMgmt          = "ALERT_MANAGEMENT"
	FeatureQuotaMgmt          = "QUOTA_MANAGEMENT"
	FeatureAuditLogView       = "AUDIT_LOG_VIEW"
)

// 操作权限常量
const (
	ActionView   = "VIEW"
	ActionCreate = "CREATE"
	ActionEdit   = "EDIT"
	ActionDelete = "DELETE"
	ActionExport = "EXPORT"
	ActionImport = "IMPORT"
)
