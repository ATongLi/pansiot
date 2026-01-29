package models

import (
	"time"
)

// User 用户表
type User struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantID         int64      `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_id" json:"tenant_id"` // 归属租户
	Username         string     `gorm:"column:username;type:varchar(50);not null;uniqueIndex" json:"username"`       // 用户名
	Email            string     `gorm:"column:email;type:varchar(100);not null;index" json:"email"`                  // 邮箱
	Phone            string     `gorm:"column:phone;type:varchar(20);index" json:"phone"`                            // 手机号
	PhoneCountryCode string     `gorm:"column:phone_country_code;type:varchar(5);default:'+86'" json:"phone_country_code"` // 手机号国家码
	PasswordHash     string     `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`                     // 密码哈希（不返回给前端）
	RealName         string     `gorm:"column:real_name;type:varchar(100)" json:"real_name"`                         // 真实姓名
	Avatar           string     `gorm:"column:avatar;type:varchar(500)" json:"avatar"`                               // 头像URL
	Status           string     `gorm:"column:status;type:varchar(20);not null;default:'ACTIVE'" json:"status"`       // 状态：ACTIVE, SUSPENDED, DELETED
	LastLoginAt      *time.Time `gorm:"column:last_login_at;type:timestamp" json:"last_login_at,omitempty"`           // 最后登录时间
	LastLoginIP      string     `gorm:"column:last_login_ip;type:varchar(45)" json:"last_login_ip"`                  // 最后登录IP
	CreatedAt        time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at;type:timestamp;index" json:"deleted_at,omitempty"` // 软删除
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Role 角色表
type Role struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantID    int64      `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_id" json:"tenant_id"` // 归属租户
	RoleCode    string     `gorm:"column:role_code;type:varchar(50);not null;index" json:"role_code"`          // 角色代码
	RoleName    string     `gorm:"column:role_name;type:varchar(100);not null" json:"role_name"`              // 角色名称
	Description string     `gorm:"column:description;type:varchar(500)" json:"description"`                    // 角色描述
	IsSystem    bool       `gorm:"column:is_system;type:boolean;not null;default:false" json:"is_system"`       // 是否系统角色
	IsDeletable bool       `gorm:"column:is_deletable;type:boolean;not null;default:true" json:"is_deletable"` // 是否可删除
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:timestamp;index" json:"deleted_at,omitempty"` // 软删除
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;type:bigint;not null;index:idx_user_id,index:idx_user_role" json:"user_id"`
	RoleID    int64     `gorm:"column:role_id;type:bigint;not null;index:idx_role_id,index:idx_user_role" json:"role_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
