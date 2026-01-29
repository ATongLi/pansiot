/**
 * Scada 硬件平台配置 API
 * 提供硬件平台配置查询的HTTP接口
 */

package api

import (
	"pansiot-scada/internal/service"

	"github.com/gofiber/fiber/v2"
)

/**
 * PlatformAPI 处理器
 */
type PlatformAPI struct {
	platformService *service.PlatformService
}

/**
 * 创建新的 PlatformAPI 实例
 */
func NewPlatformAPI(platformService *service.PlatformService) *PlatformAPI {
	return &PlatformAPI{
		platformService: platformService,
	}
}

/**
 * 注册路由
 * 注意：此方法接受已存在的 /api 路由组，避免重复创建路由组
 */
func (api *PlatformAPI) RegisterRoutes(apiGroup fiber.Router) {
	// 硬件平台路由
	apiGroup.Get("/platforms", api.getAllPlatforms)
}

/**
 * 获取所有启用的硬件平台
 * GET /api/platforms
 */
func (api *PlatformAPI) getAllPlatforms(c *fiber.Ctx) error {
	// 调用服务获取平台列表
	platforms, err := api.platformService.GetAllEnabledPlatforms()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "GET_PLATFORMS_FAILED",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    platforms,
	})
}
