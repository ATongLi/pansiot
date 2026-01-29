package quota

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// QuotaService 配额服务
type QuotaService struct {
	db  *gorm.DB
	rdb *redis.Client
	mu  sync.RWMutex
}

// NewQuotaService 创建配额服务
func NewQuotaService(db *gorm.DB, rdb *redis.Client) *QuotaService {
	return &QuotaService{
		db:  db,
		rdb: rdb,
	}
}

// TenantQuotaUsage 租户配额使用情况
type TenantQuotaUsage struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	TenantID        uint   `gorm:"not null;index" json:"tenant_id"`
	ModuleCode      string `gorm:"size:50;not null;index" json:"module_code"`
	TotalQuota      int    `json:"total_quota"`
	UsedQuota       int    `json:"used_quota"`
	RemainingQuota  int    `json:"remaining_quota"`
	AllocatedQuota  *int   `json:"allocated_quota"` // 仅集成商使用
}

// TableName 指定表名
func (TenantQuotaUsage) TableName() string {
	return "tenant_quota_usage"
}

// FeatureModule 功能模块
type FeatureModule struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ModuleCode  string `gorm:"size:50;uniqueIndex;not null" json:"module_code"`
	ModuleName  string `gorm:"size:100;not null" json:"module_name"`
	ModuleType  string `gorm:"size:20;not null" json:"module_type"` // SYSTEM_DEFAULT, OPTIONAL
	Description string `gorm:"size:500" json:"description"`
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
}

// TableName 指定表名
func (FeatureModule) TableName() string {
	return "feature_modules"
}

// TenantFeature 租户功能开通情况
type TenantFeature struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	TenantID   uint   `gorm:"not null;index" json:"tenant_id"`
	ModuleCode string `gorm:"size:50;not null;index" json:"module_code"`
	Quota      int    `gorm:"not null" json:"quota"`
	IsEnabled  bool   `gorm:"default:true" json:"is_enabled"`
}

// TableName 指定表名
func (TenantFeature) TableName() string {
	return "tenant_features"
}

// AllocateQuota 分配配额（集成商给下游租户分配）
func (s *QuotaService) AllocateQuota(ctx context.Context, integratorTenantID uint, downstreamTenantID uint, moduleCode string, quota int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 检查集成商是否有足够的配额
	var integratorQuota TenantQuotaUsage
	if err := s.db.Where("tenant_id = ? AND module_code = ?", integratorTenantID, moduleCode).First(&integratorQuota).Error; err != nil {
		return fmt.Errorf("集成商配额不存在: %w", err)
	}

	availableQuota := integratorQuota.TotalQuota - integratorQuota.UsedQuota
	if integratorQuota.AllocatedQuota != nil {
		availableQuota -= *integratorQuota.AllocatedQuota
	}

	if availableQuota < quota {
		return fmt.Errorf("集成商配额不足（可用: %d，需要: %d）", availableQuota, quota)
	}

	// 2. 为下游租户开通功能
	var tenantFeature TenantFeature
	result := s.db.Where("tenant_id = ? AND module_code = ?", downstreamTenantID, moduleCode).First(&tenantFeature)

	if result.Error == gorm.ErrRecordNotFound {
		// 新开通
		tenantFeature = TenantFeature{
			TenantID:   downstreamTenantID,
			ModuleCode: moduleCode,
			Quota:      quota,
			IsEnabled:  true,
		}
		if err := s.db.Create(&tenantFeature).Error; err != nil {
			return fmt.Errorf("开通功能失败: %w", err)
		}
	} else if result.Error != nil {
		return fmt.Errorf("查询租户功能失败: %w", result.Error)
	} else {
		// 已开通，增加配额
		tenantFeature.Quota += quota
		if err := s.db.Save(&tenantFeature).Error; err != nil {
			return fmt.Errorf("更新配额失败: %w", err)
		}
	}

	// 3. 更新集成商的已分配配额
	if integratorQuota.AllocatedQuota == nil {
		allocated := quota
		integratorQuota.AllocatedQuota = &allocated
	} else {
		*integratorQuota.AllocatedQuota += quota
	}

	if err := s.db.Save(&integratorQuota).Error; err != nil {
		return fmt.Errorf("更新集成商配额失败: %w", err)
	}

	// 4. 清除Redis缓存
	s.clearQuotaCache(ctx, downstreamTenantID, moduleCode)
	s.clearQuotaCache(ctx, integratorTenantID, moduleCode)

	return nil
}

// CheckQuota 检查配额是否足够
func (s *QuotaService) CheckQuota(ctx context.Context, tenantID uint, moduleCode string, required int) (bool, error) {
	// 1. 尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("quota:%d:%s", tenantID, moduleCode)
	val, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var remaining int
		if _, err := fmt.Sscanf(val, "%d", &remaining); err == nil {
			return remaining >= required, nil
		}
	}

	// 2. 从数据库查询
	var quotaUsage TenantQuotaUsage
	if err := s.db.Where("tenant_id = ? AND module_code = ?", tenantID, moduleCode).First(&quotaUsage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("租户未开通该功能模块")
		}
		return false, fmt.Errorf("查询配额失败: %w", err)
	}

	// 3. 缓存到Redis（5分钟过期）
	s.rdb.Set(ctx, cacheKey, quotaUsage.RemainingQuota, 300)

	return quotaUsage.RemainingQuota >= required, nil
}

// UseQuota 使用配额
func (s *QuotaService) UseQuota(ctx context.Context, tenantID uint, moduleCode string, count int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var quotaUsage TenantQuotaUsage
	if err := s.db.Where("tenant_id = ? AND module_code = ?", tenantID, moduleCode).First(&quotaUsage).Error; err != nil {
		return fmt.Errorf("查询配额失败: %w", err)
	}

	// 检查配额是否足够
	if quotaUsage.RemainingQuota < count {
		return fmt.Errorf("配额不足（剩余: %d，需要: %d）", quotaUsage.RemainingQuota, count)
	}

	// 更新使用量
	quotaUsage.UsedQuota += count
	quotaUsage.RemainingQuota -= count

	if err := s.db.Save(&quotaUsage).Error; err != nil {
		return fmt.Errorf("更新配额失败: %w", err)
	}

	// 清除缓存
	s.clearQuotaCache(ctx, tenantID, moduleCode)

	return nil
}

// ReleaseQuota 释放配额
func (s *QuotaService) ReleaseQuota(ctx context.Context, tenantID uint, moduleCode string, count int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var quotaUsage TenantQuotaUsage
	if err := s.db.Where("tenant_id = ? AND module_code = ?", tenantID, moduleCode).First(&quotaUsage).Error; err != nil {
		return fmt.Errorf("查询配额失败: %w", err)
	}

	// 更新使用量
	quotaUsage.UsedQuota -= count
	if quotaUsage.UsedQuota < 0 {
		quotaUsage.UsedQuota = 0
	}
	quotaUsage.RemainingQuota = quotaUsage.TotalQuota - quotaUsage.UsedQuota

	if err := s.db.Save(&quotaUsage).Error; err != nil {
		return fmt.Errorf("更新配额失败: %w", err)
	}

	// 清除缓存
	s.clearQuotaCache(ctx, tenantID, moduleCode)

	return nil
}

// GetQuotaStatistics 获取配额统计信息
func (s *QuotaService) GetQuotaStatistics(ctx context.Context, tenantID uint) ([]TenantQuotaUsage, error) {
	var quotas []TenantQuotaUsage
	if err := s.db.Where("tenant_id = ?", tenantID).Find(&quotas).Error; err != nil {
		return nil, fmt.Errorf("查询配额统计失败: %w", err)
	}

	// 计算剩余配额
	for i := range quotas {
		quotas[i].RemainingQuota = quotas[i].TotalQuota - quotas[i].UsedQuota
	}

	return quotas, nil
}

// ActivateFeature 为租户开通功能模块
func (s *QuotaService) ActivateFeature(ctx context.Context, tenantID uint, moduleCode string, quota int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 检查功能模块是否存在
	var module FeatureModule
	if err := s.db.Where("module_code = ?", moduleCode).First(&module).Error; err != nil {
		return fmt.Errorf("功能模块不存在: %w", err)
	}

	// 2. 检查是否已开通
	var tenantFeature TenantFeature
	result := s.db.Where("tenant_id = ? AND module_code = ?", tenantID, moduleCode).First(&tenantFeature)

	if result.Error == gorm.ErrRecordNotFound {
		// 新开通
		tenantFeature = TenantFeature{
			TenantID:   tenantID,
			ModuleCode: moduleCode,
			Quota:      quota,
			IsEnabled:  true,
		}
		if err := s.db.Create(&tenantFeature).Error; err != nil {
			return fmt.Errorf("开通功能失败: %w", err)
		}

		// 创建配额使用记录
		quotaUsage := TenantQuotaUsage{
			TenantID:       tenantID,
			ModuleCode:     moduleCode,
			TotalQuota:     quota,
			UsedQuota:      0,
			RemainingQuota: quota,
		}
		if err := s.db.Create(&quotaUsage).Error; err != nil {
			return fmt.Errorf("创建配额记录失败: %w", err)
		}
	} else if result.Error != nil {
		return fmt.Errorf("查询租户功能失败: %w", result.Error)
	} else {
		// 已开通，更新配额
		tenantFeature.Quota = quota
		if err := s.db.Save(&tenantFeature).Error; err != nil {
			return fmt.Errorf("更新配额失败: %w", err)
		}
	}

	return nil
}

// DeactivateFeature 停用租户功能模块
func (s *QuotaService) DeactivateFeature(ctx context.Context, tenantID uint, moduleCode string) error {
	return s.db.Model(&TenantFeature{}).
		Where("tenant_id = ? AND module_code = ?", tenantID, moduleCode).
		Update("is_enabled", false).Error
}

// GetActivatedFeatures 获取租户已开通的功能列表
func (s *QuotaService) GetActivatedFeatures(ctx context.Context, tenantID uint) ([]TenantFeature, error) {
	var features []TenantFeature
	err := s.db.Where("tenant_id = ? AND is_enabled = ?", tenantID, true).
		Find(&features).Error

	if err != nil {
		return nil, fmt.Errorf("查询已开通功能失败: %w", err)
	}

	return features, nil
}

// clearQuotaCache 清除配额缓存
func (s *QuotaService) clearQuotaCache(ctx context.Context, tenantID uint, moduleCode string) {
	cacheKey := fmt.Sprintf("quota:%d:%s", tenantID, moduleCode)
	s.rdb.Del(ctx, cacheKey)
}
