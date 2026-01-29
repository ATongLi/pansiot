/**
 * Scada 工程管理 API
 * 提供工程CRUD操作的HTTP接口
 */

package api

import (
	"pansiot-scada/internal/model"
	"pansiot-scada/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

/**
 * ProjectAPI 处理器
 */
type ProjectAPI struct {
	projectService *service.ProjectService
}

/**
 * 创建新的 ProjectAPI 实例
 */
func NewProjectAPI(projectService *service.ProjectService) *ProjectAPI {
	return &ProjectAPI{
		projectService: projectService,
	}
}

/**
 * 注册路由
 * 注意：此方法接受已存在的 /api 路由组
 */
func (api *ProjectAPI) RegisterRoutes(apiGroup fiber.Router) {
	// 工程管理路由
	projects := apiGroup.Group("/projects")
	{
		// TODO(依赖): 后端服务实现 - 需要实现 ProjectService
		// 当前状态: Mock实现，返回固定数据

		// 创建工程
		projects.Post("/create", api.createProject)

		// 打开工程
		projects.Post("/open", api.openProject)

		// 保存工程
		projects.Post("/save", api.saveProject)

		// 验证密码
		projects.Post("/validate-password", api.validatePassword)

		// 最近工程列表
		projects.Get("/recent", api.getRecentProjects)

		// 添加或更新最近工程
		projects.Post("/recent", api.addOrUpdateRecentProject)

		// 删除最近工程
		projects.Delete("/recent/:projectId", api.removeRecentProject)
	}
}

/**
 * 注册全局中间件和健康检查
 */
func (api *ProjectAPI) RegisterGlobalMiddleware(app *fiber.App) {
	// 中间件
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	}))

	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Scada Backend API is running",
		})
	})
}

/**
 * 创建工程
 * POST /api/projects/create
 */
func (api *ProjectAPI) createProject(c *fiber.Ctx) error {
	var req model.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "INVALID_REQUEST",
			"message": "请求格式错误",
		})
	}

	// 调用服务创建工程
	project, err := api.projectService.CreateProject(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "CREATE_FAILED",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"projectId": project.ProjectID,
			"filePath":  req.SavePath,
		},
	})
}

/**
 * 打开工程
 * POST /api/projects/open
 */
func (api *ProjectAPI) openProject(c *fiber.Ctx) error {
	var req model.OpenProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "INVALID_REQUEST",
			"message": "请求格式错误",
		})
	}

	// 调用服务打开工程
	project, err := api.projectService.OpenProject(req.FilePath, req.Password)
	if err != nil {
		// 区分错误类型
		if err == model.ErrInvalidPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "INVALID_PASSWORD",
				"message": "密码错误",
			})
		}
		if err == model.ErrInvalidSignature {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "INVALID_SIGNATURE",
				"message": "工程文件签名验证失败，文件可能已被篡改",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "OPEN_FAILED",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"project": project,
		},
	})
}

/**
 * 保存工程
 * POST /api/projects/save
 */
func (api *ProjectAPI) saveProject(c *fiber.Ctx) error {
	var req model.SaveProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "INVALID_REQUEST",
			"message": "请求格式错误",
		})
	}

	// 调用服务保存工程
	filePath, err := api.projectService.SaveProject(req.Project)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "SAVE_FAILED",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"filePath": filePath,
		},
	})
}

/**
 * 验证密码
 * POST /api/projects/validate-password
 */
func (api *ProjectAPI) validatePassword(c *fiber.Ctx) error {
	var req model.ValidatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "INVALID_REQUEST",
			"message": "请求格式错误",
		})
	}

	// 调用服务验证密码
	valid, err := api.projectService.ValidatePassword(req.FilePath, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "VALIDATE_FAILED",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"valid": valid,
		},
	})
}

/**
 * 获取最近工程列表
 * GET /api/projects/recent
 */
func (api *ProjectAPI) getRecentProjects(c *fiber.Ctx) error {
	// 调用服务获取最近工程列表
	projects, err := api.projectService.GetRecentProjects()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "GET_RECENT_FAILED",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    projects,
	})
}

/**
 * 添加或更新最近工程
 * POST /api/projects/recent
 */
func (api *ProjectAPI) addOrUpdateRecentProject(c *fiber.Ctx) error {
	var project model.RecentProject
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "INVALID_REQUEST",
			"message": "请求格式错误",
		})
	}

	// 调用服务添加或更新最近工程
	if err := api.projectService.AddOrUpdateRecentProject(project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ADD_RECENT_FAILED",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"success": true,
		},
	})
}

/**
 * 删除最近工程
 * DELETE /api/projects/recent/:projectId
 */
func (api *ProjectAPI) removeRecentProject(c *fiber.Ctx) error {
	projectID := c.Params("projectId")

	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "INVALID_PROJECT_ID",
			"message": "工程ID不能为空",
		})
	}

	// 调用服务删除最近工程
	if err := api.projectService.RemoveRecentProject(projectID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "REMOVE_RECENT_FAILED",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"success": true,
		},
	})
}
