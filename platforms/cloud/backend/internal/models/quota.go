package models

import (
	"time"

	"gorm.io/gorm"
)

// FeatureModule 功能模块表
type FeatureModule struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ModuleCode  string    `gorm:"column:module_code;type:varchar(50);not null;uniqueIndex" json:"module_code"` // 功能代码
	ModuleName  string    `gorm:"column:module_name;type:varchar(100);not null" json:"module_name"`           // 功能名称
	Description string    `gorm:"column:description;type:varchar(500)" json:"description"`                   // 功能描述
	Enabled     bool      `gorm:"column:enabled;type:boolean;not null;default:true" json:"enabled"`               // 是否启用
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (FeatureModule) TableName() string {
	return "feature_modules"
}

// TenantFeature 租户功能开通表
type TenantFeature struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantID         int64      `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_id" json:"tenant_id"`
	ModuleCode       string     `gorm:"column:module_code;type:varchar(50);not null;index:idx_module_code" json:"module_code"`
	Enabled          bool       `gorm:"column:enabled;type:boolean;not null;default:true" json:"enabled"`
	ExpiresAt        *time.Time `gorm:"column:expires_at;type:timestamp" json:"expires_at,omitempty"`
	CreatedAt        time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (TenantFeature) TableName() string {
	return "tenant_features"
}

// TenantQuota 租户配额表
type TenantQuota struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantID      int64      `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_id" json:"tenant_id"`
	QuotaType     string     `gorm:"column:quota_type;type:varchar(50);not null;index:idx_quota_type" json:"quota_type"` // 配额类型：sub_tenants, users, devices, storage_gb
	TotalQuota    int        `gorm:"column:total_quota;type:int;not null;default:0" json:"total_quota"`                // 总配额
	UsedQuota     int        `gorm:"column:used_quota;type:int;not null;default:0" json:"used_quota"`                  // 已使用配额
	RemainingQuota int       `gorm:"column:remaining_quota;type:int;not null;generated;column:remaining_quota" json:"remaining_quota"` // 剩余配额（计算列）
	ExpiresAt     *time.Time `gorm:"column:expires_at;type:timestamp" json:"expires_at,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// BeforeUpdate GORM Hook - 更新前自动计算剩余配额
func (tq *TenantQuota) BeforeUpdate(tx *gorm.DB) error {
	tq.RemainingQuota = tq.TotalQuota - tq.UsedQuota
	return nil
}

// TableName 指定表名
func (TenantQuota) TableName() string {
	return "tenant_quotas"
}

// QuotaAllocation 配额分配表（集成商为下游分配配额）
type QuotaAllocation struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentTenantID int64      `gorm:"column:parent_tenant_id;type:bigint;not null;index:idx_parent_tenant" json:"parent_tenant_id"` // 集成商ID
	ChildTenantID  int64      `gorm:"column:child_tenant_id;type:bigint;not null;index:idx_child_tenant" json:"child_tenant_id"`   // 下游客户ID
	QuotaType      string     `gorm:"column:quota_type;type:varchar(50);not null" json:"quota_type"`
	AllocatedQuota int        `gorm:"column:allocated_quota;type:int;not null;default:0" json:"allocated_quota"` // 分配的配额
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (QuotaAllocation) TableName() string {
	return "quota_allocations"
}
