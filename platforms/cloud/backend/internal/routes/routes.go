package routes

import (
	"github.com/gin-gonic/gin"

	"pansiot-cloud/internal/controller"
	"pansiot-cloud/internal/middleware"
)

// SetupRoutes 设置所有路由
func SetupRoutes(router *gin.Engine, authController *controller.AuthController, userController *controller.UserController, roleController *controller.RoleController, tenantController *controller.TenantController) {
	// 应用全局中间件
	router.Use(middleware.CORS())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestIDMiddleware())

	// 健康检查（不需要认证）
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "pansiot-cloud",
		})
	})

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 认证相关路由（不需要登录）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/send-code", authController.SendVerificationCode)
			auth.POST("/reset-password", authController.ResetPassword)

			// 需要认证的路由
			authAuthorized := auth.Group("")
			authAuthorized.Use(middleware.Auth())
			{
				authAuthorized.POST("/logout", authController.Logout)
				authAuthorized.GET("/me", authController.GetCurrentUser)
				authAuthorized.POST("/change-password", authController.ChangePassword)
			}
		}

		// 用户管理路由（需要认证）
		users := v1.Group("/users")
		users.Use(middleware.Auth())
		users.Use(middleware.TenantIsolation())
		{
			// 查询和查看（需要VIEW权限）
			users.GET("", middleware.RequirePermission("USER_MANAGEMENT", "VIEW"), userController.ListUsers)
			users.GET("/me/roles", userController.GetCurrentUserRoles)
			users.GET("/:id", middleware.RequirePermission("USER_MANAGEMENT", "VIEW"), userController.GetUser)
			users.GET("/:id/roles", middleware.RequirePermission("USER_MANAGEMENT", "VIEW"), userController.GetUserRoles)

			// 创建用户（需要CREATE权限）
			users.POST("", middleware.RequirePermission("USER_MANAGEMENT", "CREATE"), userController.CreateUser)
			users.POST("/batch-delete", middleware.RequirePermission("USER_MANAGEMENT", "DELETE"), userController.BatchDeleteUsers)
			users.POST("/batch-update-status", middleware.RequirePermission("USER_MANAGEMENT", "EDIT"), userController.UpdateUserStatus)

			// 更新用户（需要EDIT权限）
			users.PUT("/:id", middleware.RequirePermission("USER_MANAGEMENT", "EDIT"), userController.UpdateUser)
			users.PUT("/:id/roles", middleware.RequirePermission("USER_MANAGEMENT", "EDIT"), userController.AssignRoles)
			users.POST("/:id/reset-password", middleware.RequirePermission("USER_MANAGEMENT", "EDIT"), userController.ResetUserPassword)

			// 删除用户（需要DELETE权限）
			users.DELETE("/:id", middleware.RequirePermission("USER_MANAGEMENT", "DELETE"), userController.DeleteUser)
		}

		// 角色管理路由（需要认证和权限）
		roles := v1.Group("/roles")
		roles.Use(middleware.Auth())
		roles.Use(middleware.TenantIsolation())
		{
			// 查询和查看（需要VIEW权限）
			roles.GET("", middleware.RequirePermission("ROLE_MANAGEMENT", "VIEW"), roleController.ListRoles)
			roles.GET("/:id", middleware.RequirePermission("ROLE_MANAGEMENT", "VIEW"), roleController.GetRole)
			roles.GET("/:id/permissions", middleware.RequirePermission("ROLE_MANAGEMENT", "VIEW"), roleController.GetRolePermissions)
			roles.GET("/permissions/all", roleController.GetAllPermissions)

			// 创建角色（需要CREATE权限）
			roles.POST("", middleware.RequirePermission("ROLE_MANAGEMENT", "CREATE"), roleController.CreateRole)

			// 更新角色（需要EDIT权限）
			roles.PUT("/:id", middleware.RequirePermission("ROLE_MANAGEMENT", "EDIT"), roleController.UpdateRole)
			roles.PUT("/:id/permissions", middleware.RequirePermission("ROLE_MANAGEMENT", "EDIT"), roleController.AssignPermissionsToRole)

			// 删除角色（需要DELETE权限）
			roles.DELETE("/:id", middleware.RequirePermission("ROLE_MANAGEMENT", "DELETE"), roleController.DeleteRole)
		}

		// 租户管理路由（需要认证）
		tenants := v1.Group("/tenants")
		tenants.Use(middleware.Auth())
		tenants.Use(middleware.TenantIsolation())
		{
			// 当前租户信息（所有租户可访问）
			tenants.GET("/me", tenantController.GetCurrentTenant)
			tenants.PUT("/me", tenantController.UpdateCurrentTenant)
			tenants.GET("/stats", tenantController.GetTenantStats)

			// 子租户管理（仅集成商可访问）
			tenants.GET("/subs", tenantController.ListSubTenants)
			tenants.POST("/subs", tenantController.CreateSubTenant)
			tenants.GET("/subs/:id", tenantController.GetSubTenant)
		}

		// 设备管理路由（需要认证）
		devices := v1.Group("/devices")
		devices.Use(middleware.Auth())
		devices.Use(middleware.TenantIsolation())
		{
			// TODO: 添加设备管理路由
			// devices.GET("", deviceController.ListDevices)
			// devices.POST("", middleware.CheckQuota(middleware.QuotaDevices, 1), deviceController.CreateDevice)
			// devices.GET("/:id", deviceController.GetDevice)
			// devices.PUT("/:id", deviceController.UpdateDevice)
			// devices.DELETE("/:id", deviceController.DeleteDevice)
		}

		// 审计日志路由（需要认证和权限）
		audit := v1.Group("/audit-logs")
		audit.Use(middleware.Auth())
		audit.Use(middleware.TenantIsolation())
		{
			// TODO: 添加审计日志路由
			// audit.GET("", middleware.RequirePermission("AUDIT_LOG_VIEW", "VIEW"), auditController.ListAuditLogs)
			// audit.GET("/:id", middleware.RequirePermission("AUDIT_LOG_VIEW", "VIEW"), auditController.GetAuditLog)
		}
	}
}
