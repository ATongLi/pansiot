package database

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"pansiot-cloud/pkg/logger"
)

// Tenant 租户信息（用于上下文）
type Tenant struct {
	ID        int64
	TenantID  int64
	IsAdmin   bool
	TenantType string // INTEGRATOR, TERMINAL
}

// TenantScope 租户隔离Scope（下游租户）
func TenantScope(tenantID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", tenantID)
	}
}

// ManagedTenantScope 管理租户Scope（集成商查看所有下游）
func ManagedTenantScope(integratorID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 集成商可以看到：
		// 1. 自己的数据（tenant_id = integratorID）
		// 2. 所有下游租户的数据（managed_tenant_id = integratorID）
		return db.Where("tenant_id = ? OR managed_tenant_id = ?", integratorID, integratorID)
	}
}

// TenantIsolationHook 租户隔离Hook（自动应用租户过滤）
func TenantIsolationHook(tenantID int64, tenantType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 根据租户类型应用不同的Scope
		if tenantType == "INTEGRATOR" {
			// 集成商：可以看到自己和所有下游的数据
			return db.Scopes(ManagedTenantScope(tenantID))
		} else {
			// 下游客户：只能看到自己的数据
			return db.Scopes(TenantScope(tenantID))
		}
	}
}

// SoftDeleteHook 软删除Hook
func SoftDeleteHook() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 自动添加deleted_at条件
		return db.Where("deleted_at IS NULL")
	}
}

// AutoFillTenantHook 自动填充租户ID Hook
func AutoFillTenantHook(tenantID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 在创建时自动填充tenant_id
		db.Statement.SetColumn("tenant_id", tenantID)
		return db
	}
}

// AutoFillManagedTenantHook 自动填充管理租户ID Hook（集成商创建数据）
func AutoFillManagedTenantHook(managedTenantID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 在创建时自动填充managed_tenant_id
		db.Statement.SetColumn("managed_tenant_id", managedTenantID)
		return db
	}
}

// UpdateTimestampHook 自动更新时间戳Hook
func UpdateTimestampHook() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 在更新时自动更新updated_at
		db.Statement.SetColumn("updated_at", time.Now())
		return db
	}
}

// RegisterGORMHooks 注册所有GORM Hooks
func RegisterGORMHooks(db *gorm.DB, tenant *Tenant) error {
	// 1. 回调注册：自动填充租户ID
	err := db.Callback().Create().Before("gorm:create").
		Register("auto_fill_tenant", func(db *gorm.DB) {
		tenantID := getTenantIDFromContext(db)
		if tenantID > 0 {
			// 填充tenant_id
			db.Statement.SetColumn("tenant_id", tenantID)

			// 如果是集成商，还需要填充managed_tenant_id
			if tenant != nil && tenant.TenantType == "INTEGRATOR" {
				// 集成商创建的数据，managed_tenant_id指向自己
				db.Statement.SetColumn("managed_tenant_id", tenantID)
			}
		}
	})
	if err != nil {
		return err
	}

	// 2. 回调注册：自动更新时间戳
	err = db.Callback().Update().Before("gorm:update").
		Register("auto_update_timestamp", func(db *gorm.DB) {
		db.Statement.SetColumn("updated_at", time.Now())
	})
	if err != nil {
		return err
	}

	// 3. 回调注册：软删除
	err = db.Callback().Query().Before("gorm:query").
		Register("soft_delete_filter", func(db *gorm.DB) {
		// 检查模型是否实现了软删除
		if db.Statement.Schema != nil {
			hasDeletedAt := false
			for _, f := range db.Statement.Schema.Fields {
				if f.Name == "deleted_at" {
					hasDeletedAt = true
					break
				}
			}
			if hasDeletedAt {
				// 自动添加deleted_at IS NULL条件
				db.Statement.AddClause(clause.Where{
					Exprs: []clause.Expression{clause.Expr{SQL: "deleted_at IS NULL"}},
				})
			}
		}
	})
	if err != nil {
		return err
	}

	logger.Info("GORM hooks registered successfully")
	return nil
}

// getTenantIDFromContext 从上下文获取租户ID
func getTenantIDFromContext(db *gorm.DB) int64 {
	// 从GORM的上下文中获取租户ID
	if tenantID, ok := db.InstanceGet("tenant_id"); ok {
		if id, ok := tenantID.(int64); ok {
			return id
		}
	}
	return 0
}

// SetTenantToContext 设置租户信息到GORM上下文
func SetTenantToContext(db *gorm.DB, tenant *Tenant) {
	db.InstanceSet("tenant_id", tenant.TenantID)
	db.InstanceSet("tenant_type", tenant.TenantType)
	db.InstanceSet("is_admin", tenant.IsAdmin)
}

// WithTenantScope 使用租户Scope执行查询（用于手动指定租户）
func WithTenantScope(db *gorm.DB, tenant *Tenant) *gorm.DB {
	if tenant == nil {
		return db
	}

	var scopeFunc func(*gorm.DB) *gorm.DB
	if tenant.TenantType == "INTEGRATOR" {
		scopeFunc = ManagedTenantScope(tenant.TenantID)
	} else {
		scopeFunc = TenantScope(tenant.TenantID)
	}

	return db.Scopes(scopeFunc, SoftDeleteHook())
}
