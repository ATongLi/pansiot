package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/response"
)

var (
	rateLimitRedis *redis.Client
)

// InitRateLimitMiddleware 初始化限流中间件
func InitRateLimitMiddleware(redisClient *redis.Client) {
	rateLimitRedis = redisClient
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerMinute int // 每分钟请求数
	RequestsPerHour   int // 每小时请求数
	BurstSize         int // 突发流量大小
}

// RateLimit 限流中间件（基于Redis + Token Bucket算法）
// 使用示例：router.GET("/api", middleware.RateLimit(100, 1000, 10), handler)
func RateLimit(requestsPerMinute, requestsPerHour, burstSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rateLimitRedis == nil {
			// Redis未初始化，跳过限流
			c.Next()
			return
		}

		ctx := c.Request.Context()
		key := getRateLimitKey(c)

		// 1. 检查分钟级限流
		allowed, err := checkTokenBucket(ctx, key+":minute", requestsPerMinute, time.Minute)
		if err != nil {
			logger.Error(fmt.Sprintf("检查分钟级限流失败: %v", err))
			// 限流检查失败，允许通过
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		// 2. 检查小时级限流
		allowed, err = checkTokenBucket(ctx, key+":hour", requestsPerHour, time.Hour)
		if err != nil {
			logger.Error(fmt.Sprintf("检查小时级限流失败: %v", err))
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, 429, "请求次数超过限制")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByIP 基于IP的限流
func RateLimitByIP(requestsPerMinute int) gin.HandlerFunc {
	return RateLimit(requestsPerMinute, requestsPerMinute*10, requestsPerMinute/2)
}

// RateLimitByUser 基于用户的限流
func RateLimitByUser(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rateLimitRedis == nil {
			c.Next()
			return
		}

		userID := GetUserID(c)
		if userID == 0 {
			// 未登录，使用IP限流
			RateLimitByIP(requestsPerMinute)(c)
			return
		}

		// 已登录，使用用户ID限流
		key := fmt.Sprintf("ratelimit:user:%d", userID)
		ctx := c.Request.Context()

		allowed, err := checkTokenBucket(ctx, key, requestsPerMinute, time.Minute)
		if err != nil {
			logger.Error(fmt.Sprintf("检查用户限流失败: user_id=%d, error=%v", userID, err))
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// getRateLimitKey 获取限流键
func getRateLimitKey(c *gin.Context) string {
	// 优先使用用户ID
	userID := GetUserID(c)
	if userID > 0 {
		return fmt.Sprintf("ratelimit:user:%d", userID)
	}

	// 其次使用IP地址
	ip := c.ClientIP()
	return fmt.Sprintf("ratelimit:ip:%s", ip)
}

// checkTokenBucket 检查Token Bucket算法限流
func checkTokenBucket(ctx context.Context, key string, maxTokens int, interval time.Duration) (bool, error) {
	// 使用Redis的Lua脚本确保原子性
	luaScript := `
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local interval = tonumber(ARGV[2])
		local max_tokens = tonumber(ARGV[3])
		local requested = tonumber(ARGV[4])

		-- 获取当前令牌数和最后更新时间
		local info = redis.call("HMGET", key, "tokens", "last_time")
		local tokens = tonumber(info[1]) or max_tokens
		local last_time = tonumber(info[2]) or now

		-- 计算需要添加的令牌数
		local elapsed = now - last_time
		local new_tokens = math.min(max_tokens, tokens + (elapsed * max_tokens / interval))

		-- 检查是否有足够的令牌
		if new_tokens >= requested then
			new_tokens = new_tokens - requested
			redis.call("HMSET", key, "tokens", new_tokens, "last_time", now)
			redis.call("EXPIRE", key, interval)
			return 1
		else
			-- 更新令牌数（不扣减）
			redis.call("HMSET", key, "tokens", new_tokens, "last_time", now)
			redis.call("EXPIRE", key, interval)
			return 0
		end
	`

	now := time.Now().Unix()
	result, err := rateLimitRedis.Eval(ctx, luaScript, []string{key}, now, int(interval.Seconds()), maxTokens, 1).Result()
	if err != nil {
		return false, err
	}

	allowed, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("invalid result type")
	}

	return allowed == 1, nil
}

// FixedWindowRateLimit 固定窗口限流（简单实现）
func FixedWindowRateLimit(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rateLimitRedis == nil {
			c.Next()
			return
		}

		key := getRateLimitKey(c)
		ctx := c.Request.Context()

		// 获取当前分钟的窗口
		now := time.Now()
		window := now.Format("2006-01-02 15:04")
		windowKey := fmt.Sprintf("%s:%s", key, window)

		// 增加计数器
		count, err := rateLimitRedis.Incr(ctx, windowKey).Result()
		if err != nil {
			logger.Error(fmt.Sprintf("增加限流计数器失败: %v", err))
			c.Next()
			return
		}

		// 第一次设置，添加过期时间
		if count == 1 {
			rateLimitRedis.Expire(ctx, windowKey, time.Minute+time.Second)
		}

		// 检查是否超过限制
		if count > int64(requestsPerMinute) {
			response.Error(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(int(requestsPerMinute)-int(count)))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(now.Add(time.Minute).Unix(), 10))

		c.Next()
	}
}

// GetRateLimitStatus 获取限流状态
func GetRateLimitStatus(ctx context.Context, c *gin.Context) (map[string]interface{}, error) {
	if rateLimitRedis == nil {
		return nil, fmt.Errorf("Redis未初始化")
	}

	key := getRateLimitKey(c)
	minuteKey := key + ":minute"
	hourKey := key + ":hour"

	// 获取分钟级限流状态
	minuteInfo, err := rateLimitRedis.HMGet(ctx, minuteKey, "tokens", "last_time").Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 获取小时级限流状态
	hourInfo, err := rateLimitRedis.HMGet(ctx, hourKey, "tokens", "last_time").Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	status := map[string]interface{}{
		"minute": map[string]interface{}{
			"tokens":    minuteInfo[0],
			"last_time": minuteInfo[1],
		},
		"hour": map[string]interface{}{
			"tokens":    hourInfo[0],
			"last_time": hourInfo[1],
		},
	}

	return status, nil
}

// ResetRateLimit 重置限流（管理员功能）
func ResetRateLimit(ctx context.Context, key string) error {
	if rateLimitRedis == nil {
		return fmt.Errorf("Redis未初始化")
	}

	// 删除所有限流key
	pattern := key + ":*"
	keys, err := rateLimitRedis.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		err = rateLimitRedis.Del(ctx, keys...).Err()
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("已重置限流: key=%s, affected_keys=%d", key, len(keys)))
	}

	return nil
}
