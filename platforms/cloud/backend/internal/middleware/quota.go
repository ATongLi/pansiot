package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"pansiot-cloud/internal/models"
	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/response"
)

var (
	quotaDB       *gorm.DB
	quotaRedisCli *redis.Client
)

// InitQuotaMiddleware 初始化配额中间件
func InitQuotaMiddleware(db *gorm.DB, redisClient *redis.Client) {
	quotaDB = db
	quotaRedisCli = redisClient
}

// QuotaType 配额类型
type QuotaType string

const (
	QuotaSubTenants QuotaType = "sub_tenants" // 子租户数量
	QuotaUsers      QuotaType = "users"       // 用户数量
	QuotaDevices    QuotaType = "devices"     // 设备数量
	QuotaStorage    QuotaType = "storage_gb"  // 存储空间（GB）
)

// CheckQuota 检查配额中间件工厂函数
// 使用示例：router.POST("/users", middleware.CheckQuota(middleware.QuotaUsers, 1), handler.CreateUser)
func CheckQuota(quotaType QuotaType, required int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := GetTenantID(c)

		if tenantID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查配额
		hasEnough, err := checkQuota(c.Request.Context(), tenantID, quotaType, required)
		if err != nil {
			logger.Error(fmt.Sprintf("检查配额失败: tenant_id=%d, quota_type=%s, error=%v", tenantID, quotaType, err))
			response.Error(c, 500, "配额检查失败")
			c.Abort()
			return
		}

		if !hasEnough {
			response.Error(c, response.QuotaExceededCode, fmt.Sprintf("%s配额不足", quotaType))
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkQuota 检查租户配额是否足够
func checkQuota(ctx context.Context, tenantID int64, quotaType QuotaType, required int) (bool, error) {
	if quotaDB == nil {
		return false, fmt.Errorf("数据库未初始化")
	}

	// 1. 尝试从Redis缓存获取配额信息
	cacheKey := fmt.Sprintf("tenant:quota:%d:%s", tenantID, quotaType)
	if quotaRedisCli != nil {
		remainingStr, err := quotaRedisCli.Get(ctx, cacheKey).Result()
		if err == nil && remainingStr != "" {
			var remaining int
			if _, err := fmt.Sscanf(remainingStr, "%d", &remaining); err == nil {
				return remaining >= required, nil
			}
		}
	}

	// 2. 从数据库查询配额
	var quota models.TenantQuota
	err := quotaDB.WithContext(ctx).
		Where("tenant_id = ? AND quota_type = ?", tenantID, string(quotaType)).
		First(&quota).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 没有配额限制，返回true
			return true, nil
		}
		return false, err
	}

	// 3. 检查剩余配额是否足够
	hasEnough := quota.RemainingQuota >= required

	// 4. 更新Redis缓存（缓存1分钟）
	if quotaRedisCli != nil {
		quotaRedisCli.Set(ctx, cacheKey, quota.RemainingQuota, 60)
	}

	return hasEnough, nil
}

// UpdateQuotaUsage 更新配额使用量
// 创建资源时调用，增加used_quota
func UpdateQuotaUsage(ctx context.Context, tenantID int64, quotaType QuotaType, delta int) error {
	if quotaDB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 更新配额使用量
	err := quotaDB.WithContext(ctx).
		Model(&models.TenantQuota{}).
		Where("tenant_id = ? AND quota_type = ?", tenantID, string(quotaType)).
		UpdateColumn("used_quota", gorm.Expr("used_quota + ?", delta)).
		Error

	if err != nil {
		return err
	}

	// 清除缓存
	if quotaRedisCli != nil {
		cacheKey := fmt.Sprintf("tenant:quota:%d:%s", tenantID, quotaType)
		quotaRedisCli.Del(ctx, cacheKey)
	}

	logger.Info(fmt.Sprintf("更新配额使用量: tenant_id=%d, quota_type=%s, delta=%d", tenantID, quotaType, delta))
	return nil
}

// GetTenantQuota 获取租户配额信息
func GetTenantQuota(ctx context.Context, tenantID int64) (map[string]models.TenantQuota, error) {
	if quotaDB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	var quotas []models.TenantQuota
	err := quotaDB.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Find(&quotas).Error

	if err != nil {
		return nil, err
	}

	quotaMap := make(map[string]models.TenantQuota)
	for _, quota := range quotas {
		quotaMap[quota.QuotaType] = quota
	}

	return quotaMap, nil
}

// CheckAndConsumeQuota 检查并消费配额（原子操作）
// 用于创建资源时先检查配额，然后立即扣减
func CheckAndConsumeQuota(ctx context.Context, tenantID int64, quotaType QuotaType, required int) (bool, error) {
	if quotaDB == nil {
		return false, fmt.Errorf("数据库未初始化")
	}

	// 使用事务确保原子性
	tx := quotaDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询当前配额
	var quota models.TenantQuota
	err := tx.WithContext(ctx).
		Where("tenant_id = ? AND quota_type = ?", tenantID, string(quotaType)).
		First(&quota).Error

	if err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			// 没有配额限制
			return true, nil
		}
		return false, err
	}

	// 检查配额是否足够
	if quota.RemainingQuota < required {
		tx.Rollback()
		return false, nil
	}

	// 扣减配额
	err = tx.WithContext(ctx).
		Model(&models.TenantQuota{}).
		Where("tenant_id = ? AND quota_type = ?", tenantID, string(quotaType)).
		Update("used_quota", quota.UsedQuota+required).Error

	if err != nil {
		tx.Rollback()
		return false, err
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return false, err
	}

	// 清除缓存
	if quotaRedisCli != nil {
		cacheKey := fmt.Sprintf("tenant:quota:%d:%s", tenantID, quotaType)
		quotaRedisCli.Del(ctx, cacheKey)
	}

	return true, nil
}
