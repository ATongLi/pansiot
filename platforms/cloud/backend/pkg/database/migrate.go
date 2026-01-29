package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"pansiot-cloud/internal/config"
	"pansiot-cloud/pkg/logger"
)

// RunMigrations 运行数据库迁移
func RunMigrations(cfg *config.Config) error {
	dsn := cfg.Database.GetDSN()

	// 添加sslmode到DSN
	dsn += " sslmode=disable"

	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// 执行迁移
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// RollbackMigrations 回滚数据库迁移
func RollbackMigrations(cfg *config.Config, steps int) error {
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// 回滚指定步数
	if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	logger.Info(fmt.Sprintf("Rolled back %d migration step(s)", steps))
	return nil
}

// RunRawSQL 直接执行SQL文件（用于开发调试）
func RunRawSQL(cfg *config.Config, sqlFile string) error {
	// 读取SQL文件
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	// 连接数据库
	dsn := cfg.Database.GetDSN() + " sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	defer db.Close()

	// 执行SQL
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	logger.Info(fmt.Sprintf("Executed SQL file: %s", filepath.Base(sqlFile)))
	return nil
}

// CreateDatabase 创建数据库（如果不存在）
func CreateDatabase(cfg *config.Config) error {
	// 连接到默认postgres数据库
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	// 检查数据库是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pg_database WHERE datname = $1", cfg.Database.DBName).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// 如果不存在则创建
	if count == 0 {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.DBName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		logger.Info(fmt.Sprintf("Created database: %s", cfg.Database.DBName))
	} else {
		logger.Info(fmt.Sprintf("Database already exists: %s", cfg.Database.DBName))
	}

	return nil
}

// ListMigrationFiles 列出所有迁移文件
func ListMigrationFiles() ([]string, error) {
	migrationsDir := "migrations"

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob migration files: %w", err)
	}

	sort.Strings(files)
	return files, nil
}
