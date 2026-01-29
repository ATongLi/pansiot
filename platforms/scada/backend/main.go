/**
 * Scada Backend 主入口
 * 启动 HTTP API 服务器
 */

package main

import (
	"log"
	"pansiot-scada/internal/api"
	"pansiot-scada/internal/database"
	"pansiot-scada/internal/service"

	"github.com/gofiber/fiber/v2"
)

/**
 * 主函数
 */
func main() {
	// 初始化数据库
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// 自动迁移
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run auto migration: %v", err)
	}
	log.Println("Database migration completed")

	// 初始化服务
	projectService := service.NewProjectService(db)
	platformService := service.NewPlatformService(db)

	// 创建 Fiber 应用
	app := fiber.New(fiber.Config{
		AppName:      "PanTools Scada API",
		ServerHeader: "PanTools-Scada",
	})

	// 注册全局中间件和健康检查
	projectAPI := api.NewProjectAPI(projectService)
	projectAPI.RegisterGlobalMiddleware(app)

	// 创建平台API实例
	platformAPI := api.NewPlatformAPI(platformService)

	// 创建API路由组（所有API共享同一个 /api 前缀和中间件）
	apiGroup := app.Group("/api")

	// 注册工程路由
	projectAPI.RegisterRoutes(apiGroup)

	// 注册平台路由
	platformAPI.RegisterRoutes(apiGroup)

	// 启动服务器
	port := ":3000"
	log.Printf("Starting PanTools Scada API Server on %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
