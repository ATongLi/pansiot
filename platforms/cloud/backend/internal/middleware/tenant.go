package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"pansiot-cloud/pkg/database"
	"pansiot-cloud/pkg/response"
)

var (
	tenantDB *gorm.DB
)

// InitTenantMiddleware 初始化租户中间件
func InitTenantMiddleware(db *gorm.DB) {
	tenantDB = db
}

// TenantScope 租户隔离Scope（用于手动应用租户过滤）
func TenantScope(c *gin.Context) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tenantID := GetTenantID(c)
		tenantType := GetTenantType(c)

		if tenantID == 0 {
			// 未登录，返回空结果
			return db.Where("1 = 0")
		}

		// 根据租户类型应用不同的Scope
		if tenantType == "INTEGRATOR" {
			// 集成商：可以看到自己和所有下游的数据
			return db.Scopes(database.ManagedTenantScope(tenantID))
		} else {
			// 下游客户：只能看到自己的数据
			return db.Scopes(database.TenantScope(tenantID))
		}
	}
}

// TenantIsolation 租户隔离中间件
// 自动将租户信息设置到GORM上下文，确保所有查询都应用租户过滤
func TenantIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tenantDB == nil {
			c.Next()
			return
		}

		tenantID := GetTenantID(c)
		tenantType := GetTenantType(c)

		if tenantID == 0 {
			// 未登录，不设置租户上下文
			c.Next()
			return
		}

		// 创建租户信息对象
		tenant := &database.Tenant{
			TenantID:  tenantID,
			TenantType: tenantType,
		}

		// 将租户信息设置到GORM实例上下文
		// 这样在后续的数据库操作中可以自动获取租户信息
		db := tenantDB.WithContext(c.Request.Context())
		database.SetTenantToContext(db, tenant)

		// 将带有租户上下文的DB实例存入Gin上下文
		c.Set("db", db)

		c.Next()
	}
}

// GetDB 获取带有租户上下文的DB实例
func GetDB(c *gin.Context) *gorm.DB {
	if db, exists := c.Get("db"); exists {
		return db.(*gorm.DB)
	}
	return tenantDB
}

// RequireTenantType 要求特定租户类型
func RequireTenantType(tenantType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentType := GetTenantType(c)
		if currentType != tenantType {
			response.Error(c, response.PermissionDenyCode, "此功能仅对"+tenantType+"类型租户开放")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireIntegrator 要求集成商角色
func RequireIntegrator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsIntegrator(c) {
			response.Error(c, response.PermissionDenyCode, "此功能仅对集成商开放")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireTerminal 要求下游客户角色
func RequireTerminal() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsTerminal(c) {
			response.Error(c, response.PermissionDenyCode, "此功能仅对下游客户开放")
			c.Abort()
			return
		}
		c.Next()
	}
}
