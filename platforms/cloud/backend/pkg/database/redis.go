package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"pansiot-cloud/internal/config"
	"pansiot-cloud/pkg/logger"
)

var rdb *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	logger.Info("Redis connected successfully")
	return nil
}

// GetRedis 获取Redis客户端
func GetRedis() *redis.Client {
	return rdb
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if rdb == nil {
		return nil
	}

	return rdb.Close()
}
