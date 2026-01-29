/**
 * Scada 硬件平台服务
 * 提供硬件平台配置的业务逻辑
 */

package service

import (
	"pansiot-scada/internal/model"

	"gorm.io/gorm"
)

/**
 * PlatformService 处理器
 */
type PlatformService struct {
	db *gorm.DB
}

/**
 * 创建新的 PlatformService 实例
 */
func NewPlatformService(db *gorm.DB) *PlatformService {
	return &PlatformService{
		db: db,
	}
}

/**
 * 获取所有启用的硬件平台
 */
func (s *PlatformService) GetAllEnabledPlatforms() ([]model.HardwarePlatformConfig, error) {
	// TODO(依赖): 后续从数据库读取平台配置
	// 当前状态: 返回Mock数据

	allPlatforms := model.GetAllHardwarePlatforms()

	// 过滤出启用的平台
	var enabledPlatforms []model.HardwarePlatformConfig
	for _, platform := range allPlatforms {
		if platform.Enabled {
			enabledPlatforms = append(enabledPlatforms, platform)
		}
	}

	return enabledPlatforms, nil
}

/**
 * 根据ID获取硬件平台
 */
func (s *PlatformService) GetPlatformByID(id string) (*model.HardwarePlatformConfig, error) {
	// TODO(依赖): 后续从数据库读取平台配置
	// 当前状态: 返回Mock数据

	allPlatforms := model.GetAllHardwarePlatforms()

	for _, platform := range allPlatforms {
		if platform.ID == id {
			return &platform, nil
		}
	}

	return nil, nil
}
