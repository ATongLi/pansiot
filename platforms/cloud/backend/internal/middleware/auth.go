package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"pansiot-cloud/internal/config"
	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/response"
)

// Claims JWT Claims
type Claims struct {
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
	TenantID   int64  `json:"tenant_id"`
	TenantType string `json:"tenant_type,omitempty"` // INTEGRATOR, TERMINAL
	TokenType  string `json:"token_type"`           // access, refresh
	jwt.RegisteredClaims
}

var (
	jwtSecret []byte
	redisClient *redis.Client
)

// InitJWT 初始化JWT
func InitJWT(cfg *config.Config) {
	jwtSecret = []byte(cfg.JWT.Secret)
}

// InitRedisClient 初始化Redis客户端（用于Token黑名单）
func InitRedisClient(client *redis.Client) {
	redisClient = client
}

// GenerateToken 生成JWT Token
func GenerateToken(userID int64, username string, tenantID int64, tenantType string, tokenType string, expireHours int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(expireHours) * time.Hour)

	// 生成JTI（JWT ID）用于黑名单
	jti := fmt.Sprintf("%d:%d:%s", userID, tenantID, tokenType)

	claims := Claims{
		UserID:     userID,
		Username:   username,
		TenantID:   tenantID,
		TenantType: tenantType,
		TokenType:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "pansiot-cloud",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateTokenPair 生成Access Token和Refresh Token
func GenerateTokenPair(userID int64, username string, tenantID int64, tenantType string) (accessToken, refreshToken string, err error) {
	// Access Token: 2小时有效期
	accessToken, err = GenerateToken(userID, username, tenantID, tenantType, "access", 2)
	if err != nil {
		return "", "", err
	}

	// Refresh Token: 7天有效期
	refreshToken, err = GenerateToken(userID, username, tenantID, tenantType, "refresh", 24*7)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 检查Token是否在黑名单中
		if IsTokenBlacklisted(claims.ID) {
			return nil, errors.New("token已失效")
		}
		return claims, nil
	}

	return nil, errors.New("token无效")
}

// IsTokenBlacklisted 检查Token是否在黑名单中
func IsTokenBlacklisted(jti string) bool {
	if redisClient == nil {
		return false
	}

	ctx := context.Background()
	key := fmt.Sprintf("token:blacklist:%s", jti)
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		logger.Error(fmt.Sprintf("检查Token黑名单失败: error=%v", err))
		return false
	}

	return exists > 0
}

// BlacklistToken 将Token加入黑名单
func BlacklistToken(tokenString string) error {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return err
	}

	if redisClient == nil {
		return errors.New("Redis客户端未初始化")
	}

	ctx := context.Background()
	key := fmt.Sprintf("token:blacklist:%s", claims.ID)

	// 计算Token剩余有效期
	expiration := time.Until(claims.ExpiresAt.Time)
	if expiration <= 0 {
		return nil // Token已过期，无需加入黑名单
	}

	// 将Token加入黑名单，过期时间与Token过期时间一致
	err = redisClient.Set(ctx, key, "1", expiration).Err()
	if err != nil {
		logger.Error(fmt.Sprintf("添加Token到黑名单失败: error=%v", err))
		return err
	}

	logger.Info(fmt.Sprintf("Token已加入黑名单: jti=%s, user_id=%d", claims.ID, claims.UserID))
	return nil
}

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// Bearer Token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		// 解析Token
		claims, err := ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Token无效或已过期")
			c.Abort()
			return
		}

		// 验证Token类型
		if claims.TokenType != "access" {
			response.Unauthorized(c, "Token类型错误")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("tenant_id", claims.TenantID)
		c.Set("tenant_type", claims.TenantType)
		c.Set("jti", claims.ID)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（允许未登录访问，但会解析Token）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 未提供Token，继续处理请求
			c.Next()
			return
		}

		// Bearer Token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		// 尝试解析Token
		claims, err := ParseToken(parts[1])
		if err != nil {
			// Token无效，但不阻止请求
			c.Next()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("tenant_id", claims.TenantID)
		c.Set("tenant_type", claims.TenantType)
		c.Set("jti", claims.ID)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) int64 {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(int64)
	}
	return 0
}

// GetTenantID 从上下文获取租户ID
func GetTenantID(c *gin.Context) int64 {
	if tenantID, exists := c.Get("tenant_id"); exists {
		return tenantID.(int64)
	}
	return 0
}

// GetTenantType 从上下文获取租户类型
func GetTenantType(c *gin.Context) string {
	if tenantType, exists := c.Get("tenant_type"); exists {
		return tenantType.(string)
	}
	return ""
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// GetJTI 从上下文获取JWT ID
func GetJTI(c *gin.Context) string {
	if jti, exists := c.Get("jti"); exists {
		return jti.(string)
	}
	return ""
}

// IsAuthenticated 检查用户是否已认证
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// IsIntegrator 检查当前用户是否为集成商
func IsIntegrator(c *gin.Context) bool {
	tenantType := GetTenantType(c)
	return tenantType == "INTEGRATOR"
}

// IsTerminal 检查当前用户是否为下游客户
func IsTerminal(c *gin.Context) bool {
	tenantType := GetTenantType(c)
	return tenantType == "TERMINAL"
}
