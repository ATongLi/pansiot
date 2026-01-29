package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"pansiot-cloud/internal/models"
	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/response"
)

var (
	permDB       *gorm.DB
	permRedisCli *redis.Client
)

// InitPermissionMiddleware 初始化权限中间件
func InitPermissionMiddleware(db *gorm.DB, redisClient *redis.Client) {
	permDB = db
	permRedisCli = redisClient
}

// PermissionRequirement 权限要求
type PermissionRequirement struct {
	FeatureCode string // 功能代码：USER_MANAGEMENT, DEVICE_MANAGEMENT等
	ActionCode  string // 操作代码：VIEW, CREATE, EDIT, DELETE等
}

// RequirePermission 权限验证中间件工厂函数
// 使用示例：router.POST("/users", middleware.RequirePermission("USER_MANAGEMENT", "CREATE"), handler.CreateUser)
func RequirePermission(featureCode, actionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		tenantID := GetTenantID(c)

		if userID == 0 || tenantID == 0 {
			response.Unauthorized(c, "未登录或登录已过期")
			c.Abort()
			return
		}

		// 检查权限
		hasPermission, err := CheckPermission(c.Request.Context(), userID, tenantID, featureCode, actionCode)
		if err != nil {
			logger.Error(fmt.Sprintf("检查权限失败: user_id=%d, feature=%s, action=%s, error=%v", userID, featureCode, actionCode, err))
			response.Error(c, 500, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			logger.Warn(fmt.Sprintf("权限不足: user_id=%d, tenant_id=%d, feature=%s, action=%s", userID, tenantID, featureCode, actionCode))
			response.Error(c, response.PermissionDenyCode, fmt.Sprintf("没有权限执行此操作：需要 %s.%s 权限", featureCode, actionCode))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole 角色验证中间件工厂函数
// 使用示例：router.POST("/admin", middleware.RequireRole("SYSTEM_ADMIN"), handler.AdminHandler)
func RequireRole(roleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		tenantID := GetTenantID(c)

		if userID == 0 || tenantID == 0 {
			response.Unauthorized(c, "未登录或登录已过期")
			c.Abort()
			return
		}

		// 检查角色
		hasRole, err := CheckRole(c.Request.Context(), userID, tenantID, roleCode)
		if err != nil {
			logger.Error(fmt.Sprintf("检查角色失败: user_id=%d, role=%s, error=%v", userID, roleCode, err))
			response.Error(c, 500, "角色检查失败")
			c.Abort()
			return
		}

		if !hasRole {
			logger.Warn(fmt.Sprintf("角色不足: user_id=%d, tenant_id=%d, required_role=%s", userID, tenantID, roleCode))
			response.Error(c, response.PermissionDenyCode, fmt.Sprintf("需要 %s 角色才能访问", roleCode))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 要求系统管理员角色
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("SYSTEM_ADMIN")
}

// CheckPermission 检查用户是否有指定权限
func CheckPermission(ctx context.Context, userID, tenantID int64, featureCode, actionCode string) (bool, error) {
	// 1. 尝试从Redis缓存获取权限
	cacheKey := fmt.Sprintf("user:permissions:%d:%d", tenantID, userID)
	if permRedisCli != nil {
		permissionKey := fmt.Sprintf("%s:%s", featureCode, actionCode)
		exists, err := permRedisCli.SIsMember(ctx, cacheKey, permissionKey).Result()
		if err == nil && exists {
			return true, nil
		}
	}

	// 2. 从数据库查询权限
	if permDB == nil {
		return false, fmt.Errorf("数据库未初始化")
	}

	// 查询用户的角色
	var userRoles []models.UserRole
	err := permDB.WithContext(ctx).
		Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.tenant_id = ? AND roles.deleted_at IS NULL", userID, tenantID).
		Find(&userRoles).Error
	if err != nil {
		return false, err
	}

	if len(userRoles) == 0 {
		return false, nil
	}

	// 提取角色ID列表
	roleIDs := make([]int64, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	// 查询角色权限
	var rolePermissions []models.RolePermission
	err = permDB.WithContext(ctx).
		Table("role_permissions").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN ? AND permissions.tenant_id = ? AND permissions.feature_code = ? AND permissions.action_code = ?",
			roleIDs, tenantID, featureCode, actionCode).
		Find(&rolePermissions).Error
	if err != nil {
		return false, err
	}

	hasPermission := len(rolePermissions) > 0

	// 3. 如果有权限，更新Redis缓存
	if hasPermission && permRedisCli != nil {
		permissionKey := fmt.Sprintf("%s:%s", featureCode, actionCode)
		permRedisCli.SAdd(ctx, cacheKey, permissionKey)
		// 缓存1小时
		permRedisCli.Expire(ctx, cacheKey, time.Hour)
	}

	return hasPermission, nil
}

// CheckRole 检查用户是否有指定角色
func CheckRole(ctx context.Context, userID, tenantID int64, roleCode string) (bool, error) {
	if permDB == nil {
		return false, fmt.Errorf("数据库未初始化")
	}

	// 查询用户角色
	var userRole models.UserRole
	err := permDB.WithContext(ctx).
		Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.tenant_id = ? AND roles.role_code = ? AND roles.deleted_at IS NULL",
			userID, tenantID, roleCode).
		First(&userRole).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetUserPermissions 获取用户所有权限（带缓存）
func GetUserPermissions(ctx context.Context, userID, tenantID int64) (map[string][]string, error) {
	// 尝试从Redis获取
	cacheKey := fmt.Sprintf("user:permissions:all:%d:%d", tenantID, userID)
	if permRedisCli != nil {
		data, err := permRedisCli.Get(ctx, cacheKey).Result()
		if err == nil && data != "" {
			// TODO: 解析缓存数据
			// 简化处理：这里返回空map，实际应该解析JSON
			return make(map[string][]string), nil
		}
	}

	if permDB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 查询用户所有角色
	var userRoles []models.UserRole
	err := permDB.WithContext(ctx).
		Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.tenant_id = ? AND roles.deleted_at IS NULL", userID, tenantID).
		Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return make(map[string][]string), nil
	}

	roleIDs := make([]int64, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	// 查询所有权限
	var permissions []models.Permission
	err = permDB.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN ? AND permissions.tenant_id = ?", roleIDs, tenantID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	// 组织权限数据：feature_code -> [action_code1, action_code2, ...]
	permMap := make(map[string][]string)
	for _, perm := range permissions {
		permMap[perm.FeatureCode] = append(permMap[perm.FeatureCode], perm.ActionCode)
	}

	return permMap, nil
}

// InvalidateUserPermissionsCache 清除用户权限缓存
// 当用户角色或权限变更时调用
func InvalidateUserPermissionsCache(ctx context.Context, userID, tenantID int64) {
	if permRedisCli == nil {
		return
	}

	cacheKey := fmt.Sprintf("user:permissions:%d:%d", tenantID, userID)
	err := permRedisCli.Del(ctx, cacheKey).Err()
	if err != nil {
		logger.Error(fmt.Sprintf("清除权限缓存失败: user_id=%d, tenant_id=%d, error=%v", userID, tenantID, err))
	} else {
		logger.Info(fmt.Sprintf("已清除用户权限缓存: user_id=%d, tenant_id=%d", userID, tenantID))
	}
}

// InvalidateRolePermissionsCache 清除角色所有用户的权限缓存
// 当角色权限变更时调用
func InvalidateRolePermissionsCache(ctx context.Context, roleID int64) {
	if permDB == nil || permRedisCli == nil {
		return
	}

	// 查询该角色下的所有用户
	// 查询使用该角色的所有用户及其租户ID
	type UserRoleWithTenant struct {
		UserID   int64
		TenantID int64
	}
	var userRoles []UserRoleWithTenant
	err := permDB.WithContext(ctx).
		Table("user_roles").
		Select("user_roles.user_id, users.tenant_id").
		Joins("JOIN users ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", roleID).
		Scan(&userRoles).Error
	if err != nil {
		logger.Error(fmt.Sprintf("查询角色用户失败: role_id=%d, error=%v", roleID, err))
		return
	}

	// 清除所有用户的权限缓存
	for _, ur := range userRoles {
		cacheKey := fmt.Sprintf("user:permissions:%d:%d", ur.TenantID, ur.UserID)
		permRedisCli.Del(ctx, cacheKey)
	}

	logger.Info(fmt.Sprintf("已清除角色权限缓存: role_id=%d, affected_users=%d", roleID, len(userRoles)))
}

// ParsePermissionFromRoute 从路由路径解析权限要求
// 例如：/api/v1/users -> USER_MANAGEMENT + VIEW
//        POST /api/v1/users -> USER_MANAGEMENT + CREATE
func ParsePermissionFromRoute(method, path string) (featureCode, actionCode string) {
	// 移除 /api/v1 前缀
	path = strings.TrimPrefix(path, "/api/v1/")
	path = strings.TrimPrefix(path, "/api/")

	// 分割路径
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return "UNKNOWN", "VIEW"
	}

	// 资源类型（复数转单数）
	resource := parts[0]
	featureCode = resourceToFeatureCode(resource)

	// 根据HTTP方法确定操作
	switch method {
	case "GET":
		actionCode = "VIEW"
	case "POST":
		actionCode = "CREATE"
	case "PUT", "PATCH":
		actionCode = "EDIT"
	case "DELETE":
		actionCode = "DELETE"
	default:
		actionCode = "VIEW"
	}

	return featureCode, actionCode
}

// resourceToFeatureCode 将资源路径转换为功能代码
func resourceToFeatureCode(resource string) string {
	// 移除复数s
	resource = strings.TrimSuffix(resource, "s")

	// 转换为大写
	resource = strings.ToUpper(resource)

	// 特殊映射
	mappings := map[string]string{
		"USER":       "USER_MANAGEMENT",
		"TENANT":     "ORGANIZATION_MANAGEMENT",
		"ROLE":       "ROLE_MANAGEMENT",
		"DEVICE":     "DEVICE_MANAGEMENT",
		"QUOTA":      "QUOTA_MANAGEMENT",
		"AUDIT":      "AUDIT_LOG_VIEW",
		"PERMISSION": "ROLE_MANAGEMENT",
	}

	if code, ok := mappings[resource]; ok {
		return code
	}

	return resource + "_MANAGEMENT"
}
