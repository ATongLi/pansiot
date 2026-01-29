/**
 * Scada 数据库初始化和迁移
 * 管理SQLite数据库连接和表结构
 */

package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"pansiot-scada/internal/model"

	_ "modernc.org/sqlite" // 纯 Go SQLite 驱动
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/**
 * PanTool 数据库目录
 */
const PanToolDir = ".pansiot"

/**
 * 数据库文件名
 */
const DBFileName = "pantool.db"

/**
 * InitDB 初始化数据库连接
 * 创建数据库文件并建立连接
 */
func InitDB() (*gorm.DB, error) {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// 创建 .pansiot 目录
	pantoolDir := filepath.Join(homeDir, PanToolDir)
	if err := os.MkdirAll(pantoolDir, 0755); err != nil {
		return nil, err
	}

	// 构建数据库路径
	dbPath := filepath.Join(pantoolDir, DBFileName)

	// 配置 GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	// 先用database/sql打开数据库
	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// 然后用GORM包装
	var db *gorm.DB
	db, err = gorm.Open(sqlite.Dialector{Conn: sqlDB}, config)
	if err != nil {
		return nil, err
	}

	// 获取底层数据库连接并配置
	dbSQL, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	dbSQL.SetMaxIdleConns(10)
	dbSQL.SetMaxOpenConns(100)

	return db, nil
}

/**
 * AutoMigrate 自动迁移数据库表结构
 * 创建所有必需的表和索引
 */
func AutoMigrate(db *gorm.DB) error {
	// 迁移所有表结构
	err := db.AutoMigrate(
		// 1. 最近工程表
		&model.RecentProject{},
		// 2. 自定义分类表
		&model.CustomCategory{},
		// 3. 应用配置表
		&model.AppConfig{},
		// 4. 用户偏好表
		&model.UserPreference{},
		// 5. 审计日志表
		&model.AuditLog{},
		// 6. KEK版本管理表
		&model.KEKVersion{},
	)

	if err != nil {
		return err
	}

	// 初始化默认配置
	if err := initDefaultConfig(db); err != nil {
		return err
	}

	return nil
}

/**
 * initDefaultConfig 初始化默认配置
 */
func initDefaultConfig(db *gorm.DB) error {
	// 检查是否已初始化
	var count int64
	db.Model(&model.AppConfig{}).Where("`key` = ?", "version").Count(&count)
	if count > 0 {
		return nil // 已初始化，跳过
	}

	// 插入默认配置
	defaultConfigs := []model.AppConfig{
		{Key: "version", Value: "1.0.0"},
		{Key: "theme", Value: "light"},
		{Key: "language", Value: "zh-CN"},
	}

	for _, config := range defaultConfigs {
		if err := db.Create(&config).Error; err != nil {
			return err
		}
	}

	return nil
}

/**
 * GetDBPath 获取数据库文件完整路径
 */
func GetDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, PanToolDir, DBFileName), nil
}

/**
 * CloseDB 关闭数据库连接
 */
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

/**
 * DropAllTables 删除所有表（仅用于测试）
 * 警告：此操作会删除所有数据！
 */
func DropAllTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&model.RecentProject{},
		&model.CustomCategory{},
		&model.AppConfig{},
		&model.UserPreference{},
		&model.AuditLog{},
		&model.KEKVersion{},
	)
}
