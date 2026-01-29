package models

import (
	"time"

	"gorm.io/gorm"
)

// Tenant 租户表
type Tenant struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SerialNumber     string     `gorm:"column:serial_number;type:varchar(8);not null;uniqueIndex" json:"serial_number"` // 企业序列号
	Name             string     `gorm:"column:name;type:varchar(200);not null" json:"name"`                          // 企业名称
	TenantType       string     `gorm:"column:tenant_type;type:varchar(20);not null;default:'TERMINAL'" json:"tenant_type"` // 租户类型：INTEGRATOR(集成商)、TERMINAL(下游客户)
	Industry         string     `gorm:"column:industry;type:varchar(100);default:'其他'" json:"industry"`                 // 所属行业
	ContactPerson    string     `gorm:"column:contact_person;type:varchar(100)" json:"contact_person"`                // 联系人
	ContactPhone     string     `gorm:"column:contact_phone;type:varchar(20)" json:"contact_phone"`                   // 联系电话
	ContactEmail     string     `gorm:"column:contact_email;type:varchar(100)" json:"contact_email"`                   // 联系邮箱
	ParentTenantID   *int64     `gorm:"column:parent_tenant_id;type:bigint" json:"parent_tenant_id,omitempty"`         // 上级租户ID
	Status           string     `gorm:"column:status;type:varchar(20);not null;default:'ACTIVE'" json:"status"`         // 状态：ACTIVE, SUSPENDED, DELETED
	ExpireDate       *time.Time `gorm:"column:expire_date;type:timestamp" json:"expire_date,omitempty"`                // 过期日期
	MaxSubTenants    int        `gorm:"column:max_sub_tenants;type:int;not null;default:0" json:"max_sub_tenants"`     // 最大子租户数（0=不限制）
	MaxUsers         int        `gorm:"column:max_users;type:int;not null;default:0" json:"max_users"`                 // 最大用户数（0=不限制）
	MaxDevices       int        `gorm:"column:max_devices;type:int;not null;default:0" json:"max_devices"`             // 最大设备数（0=不限制)
	MaxStorageGB     int        `gorm:"column:max_storage_gb;type:int;not null;default:0" json:"max_storage_gb"`       // 最大存储空间（GB，0=不限制）
	CreatedAt        time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at;type:timestamp;index" json:"deleted_at,omitempty"` // 软删除
}

// TableName 指定表名
func (Tenant) TableName() string {
	return "tenants"
}

// BeforeCreate GORM Hook - 创建前自动生成序列号
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	// TODO: 实现序列号生成逻辑
	// 4位随机字符 + 4位自增ID
	return nil
}
