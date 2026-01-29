package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"pansiot-cloud/internal/config"
	"pansiot-cloud/internal/controller"
	"pansiot-cloud/internal/middleware"
	"pansiot-cloud/internal/routes"
	"pansiot-cloud/internal/service"
	"pansiot-cloud/pkg/database"
	"pansiot-cloud/pkg/logger"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志
	if err := logger.InitLogger(cfg); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting PansIot Cloud Platform...")

	// 3. 初始化数据库
	if err := database.InitPostgres(cfg); err != nil {
		logger.Error(fmt.Sprintf("Failed to init database: %v", err))
		log.Fatalf("Failed to init database: %v", err)
	}
	defer database.CloseDB()

	// 4. 初始化Redis
	if err := database.InitRedis(cfg); err != nil {
		logger.Error(fmt.Sprintf("Failed to init redis: %v", err))
		log.Fatalf("Failed to init redis: %v", err)
	}
	defer database.CloseRedis()

	// 5. 获取DB和Redis实例
	db := database.GetDB()
	rdb := database.GetRedis()

	// 6. 初始化中间件
	middleware.InitJWT(cfg)
	middleware.InitRedisClient(rdb)
	middleware.InitPermissionMiddleware(db, rdb)
	middleware.InitTenantMiddleware(db)
	middleware.InitAuditMiddleware(db, rdb)
	middleware.InitQuotaMiddleware(db, rdb)
	middleware.InitRateLimitMiddleware(rdb)

	// 7. 注册GORM Hooks
	// TODO: 创建租户对象并传入
	// err = database.RegisterGORMHooks(db, tenant)
	// if err != nil {
	// 	logger.Error(fmt.Sprintf("Failed to register GORM hooks: %v", err))
	// }

	// 8. 创建服务层
	authService := service.NewAuthService(db, rdb)
	userService := service.NewUserService(db)
	roleService := service.NewRoleService(db)
	tenantService := service.NewTenantService(db)

	// 9. 创建控制器
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	roleController := controller.NewRoleController(roleService)
	tenantController := controller.NewTenantController(tenantService)

	// 10. 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 11. 创建路由
	router := gin.New()

	// 12. 注册路由
	routes.SetupRoutes(router, authController, userController, roleController, tenantController)

	// 13. 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info(fmt.Sprintf("Server starting on %s", addr))

	// 优雅关闭
	go func() {
		if err := router.Run(addr); err != nil {
			logger.Error(fmt.Sprintf("Failed to start server: %v", err))
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 关闭审计日志中间件
	middleware.CloseAuditMiddleware()

	// 优雅关闭，等待5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建HTTP服务器以支持优雅关闭
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	logger.Info("Server exited")
}
