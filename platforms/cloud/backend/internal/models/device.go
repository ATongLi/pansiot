package models

import (
	"time"
)

// Device 设备表
type Device struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantID         int64      `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_id" json:"tenant_id"`         // 归属租户
	ManagedTenantID  *int64     `gorm:"column:managed_tenant_id;type:bigint;index:idx_managed_tenant_id" json:"managed_tenant_id,omitempty"` // 管理租户（集成商）
	DeviceCode       string     `gorm:"column:device_code;type:varchar(50);not null;index:idx_device_code" json:"device_code"` // 设备编码
	DeviceName       string     `gorm:"column:device_name;type:varchar(100);not null" json:"device_name"`       // 设备名称
	DeviceType       string     `gorm:"column:device_type;type:varchar(50);not null" json:"device_type"`       // 设备类型
	Status           string     `gorm:"column:status;type:varchar(20);not null;default:'offline'" json:"status"` // online, offline, error
	LastOnlineAt     *time.Time `gorm:"column:last_online_at;type:timestamp" json:"last_online_at,omitempty"`
	FirmwareVersion  string     `gorm:"column:firmware_version;type:varchar(50)" json:"firmware_version"`
	Description      string     `gorm:"column:description;type:varchar(500)" json:"description"`
	Location         string     `gorm:"column:location;type:varchar(200)" json:"location"`
	CreatedAt        time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at;type:timestamp;index" json:"deleted_at,omitempty"` // 软删除
}

// TableName 指定表名
func (Device) TableName() string {
	return "devices"
}

// DeviceStatus 设备状态常量
const (
	DeviceStatusOnline  = "online"
	DeviceStatusOffline = "offline"
	DeviceStatusError   = "error"
)
