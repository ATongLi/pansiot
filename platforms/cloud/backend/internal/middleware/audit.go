package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"pansiot-cloud/internal/models"
	"pansiot-cloud/pkg/logger"
)

var (
	auditDB       *gorm.DB
	auditRedisCli *redis.Client
	auditChan     chan *models.AuditLog
)

const (
	auditChannelSize = 1000 // 审计日志缓冲区大小
	auditBatchSize   = 100  // 批量插入大小
	auditFlushTime   = 5 * time.Second
)

// InitAuditMiddleware 初始化审计日志中间件
func InitAuditMiddleware(db *gorm.DB, redisClient *redis.Client) {
	auditDB = db
	auditRedisCli = redisClient
	auditChan = make(chan *models.AuditLog, auditChannelSize)

	// 启动后台协程处理审计日志
	go processAuditLogs()
}

// AuditLog 审计日志中间件
// 记录所有请求的操作日志
func AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只记录需要审计的操作
		if !shouldAudit(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		startTime := time.Now()

		// 读取请求体（用于记录）
		var bodyBytes []byte
		if c.Request.Body != nil && c.Request.Body != http.NoBody {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// 恢复请求体供后续处理器使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 使用响应writer包装器捕获响应
		writer := &responseWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = writer

		// 执行请求
		c.Next()

		// 请求完成后记录审计日志
		duration := time.Since(startTime)
		go func() {
			auditLog := buildAuditLog(c, bodyBytes, writer.body.Bytes(), duration)
			if auditLog != nil {
				select {
				case auditChan <- auditLog:
					// 成功发送到通道
				default:
					// 通道已满，记录警告
					logger.Error("审计日志通道已满，丢弃日志")
				}
			}
		}()
	}
}

// shouldAudit 判断是否需要审计
func shouldAudit(method, path string) bool {
	// 只审计修改数据的操作
	if method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
		return false
	}

	// 排除健康检查等接口
	if strings.Contains(path, "/health") || strings.Contains(path, "/metrics") {
		return false
	}

	// 排除登录接口
	if strings.Contains(path, "/login") || strings.Contains(path, "/register") {
		return false
	}

	return true
}

// responseWriter 响应writer包装器（用于捕获响应体）
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// buildAuditLog 构建审计日志
func buildAuditLog(c *gin.Context, requestBody, responseBody []byte, duration time.Duration) *models.AuditLog {
	userID := GetUserID(c)
	tenantID := GetTenantID(c)
	username := GetUsername(c)

	if userID == 0 || tenantID == 0 {
		return nil
	}

	// 解析请求路径和操作类型
	moduleCode, actionType := parseAuditInfo(c.Request.Method, c.Request.URL.Path)

	// 构建action_detail（JSONB格式）
	actionDetail := models.ActionDetailJSON{
		Before:  nil, // CREATE操作没有before
		After:   parseJSONSafe(requestBody),
		Changes: nil,
	}

	// 对于UPDATE和DELETE操作，如果有entity_id，记录before和changes
	// 这里简化处理，实际应该在handler层面记录

	auditLog := &models.AuditLog{
		TenantID:     tenantID,
		ModuleCode:   moduleCode,
		ActionType:   actionType,
		EntityType:   parseEntityType(c.Request.URL.Path),
		EntityID:     parseEntityID(c.Request.URL.Path),
		ActionDetail: actionDetail,
		OperatorID:   userID,
		OperatorName: username,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		Status:       "success",
		CreatedAt:    time.Now(),
	}

	// 如果请求失败，记录失败状态
	if len(c.Errors) > 0 || c.Writer.Status() >= 400 {
		auditLog.Status = "failed"
	}

	return auditLog
}

// parseAuditInfo 解析审计信息（模块代码和操作类型）
func parseAuditInfo(method, path string) (moduleCode, actionType string) {
	// 操作类型
	switch method {
	case "POST":
		actionType = models.AuditActionCreate
	case "PUT", "PATCH":
		actionType = models.AuditActionUpdate
	case "DELETE":
		actionType = models.AuditActionDelete
	default:
		actionType = "VIEW"
	}

	// 模块代码（从路径解析）
	path = strings.TrimPrefix(path, "/api/v1/")
	path = strings.TrimPrefix(path, "/api/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		resource := parts[0]
		moduleCode = resourceToFeatureCode(resource)
	}

	return moduleCode, actionType
}

// parseEntityType 解析实体类型
func parseEntityType(path string) string {
	path = strings.TrimPrefix(path, "/api/v1/")
	path = strings.TrimPrefix(path, "/api/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return strings.ToUpper(parts[0])
	}
	return "UNKNOWN"
}

// parseEntityID 从路径解析实体ID
// 例如：/api/v1/users/123 -> 123
func parseEntityID(path string) *int64 {
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		var id int64
		if _, err := fmt.Sscanf(lastPart, "%d", &id); err == nil {
			return &id
		}
	}
	return nil
}

// parseJSONSafe 安全解析JSON，失败则返回原始字符串
func parseJSONSafe(data []byte) interface{} {
	if len(data) == 0 {
		return nil
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return string(data)
	}

	return result
}

// processAuditLogs 后台协程处理审计日志
func processAuditLogs() {
	ticker := time.NewTicker(auditFlushTime)
	defer ticker.Stop()

	batch := make([]*models.AuditLog, 0, auditBatchSize)

	for {
		select {
		case log, ok := <-auditChan:
			if !ok {
				// 通道关闭，保存剩余日志
				if len(batch) > 0 {
					saveAuditLogs(batch)
				}
				return
			}

			batch = append(batch, log)

			// 达到批量大小，保存
			if len(batch) >= auditBatchSize {
				saveAuditLogs(batch)
				batch = make([]*models.AuditLog, 0, auditBatchSize)
			}

		case <-ticker.C:
			// 定时刷新
			if len(batch) > 0 {
				saveAuditLogs(batch)
				batch = make([]*models.AuditLog, 0, auditBatchSize)
			}
		}
	}
}

// saveAuditLogs 批量保存审计日志
func saveAuditLogs(logs []*models.AuditLog) {
	if len(logs) == 0 || auditDB == nil {
		return
	}

	ctx := context.Background()
	err := auditDB.WithContext(ctx).Create(&logs).Error
	if err != nil {
		logger.Error(fmt.Sprintf("批量保存审计日志失败: %v", err))
		// 可以在这里实现重试逻辑或fallback到Redis
	} else {
		logger.Info(fmt.Sprintf("成功保存 %d 条审计日志", len(logs)))
	}
}

// AuditCustomOperation 记录自定义操作（用于手动记录）
// 使用场景：在业务逻辑中记录特定的操作
func AuditCustomOperation(ctx context.Context, tenantID, operatorID int64, operatorName, moduleCode, actionType, entityType string, entityID *int64, before, after interface{}, changes []models.Change) error {
	if auditDB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	actionDetail := models.ActionDetailJSON{
		Before:  before,
		After:   after,
		Changes: changes,
	}

	auditLog := &models.AuditLog{
		TenantID:     tenantID,
		ModuleCode:   moduleCode,
		ActionType:   actionType,
		EntityType:   entityType,
		EntityID:     entityID,
		ActionDetail: actionDetail,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Status:       "success",
		CreatedAt:    time.Now(),
	}

	// 直接保存（不通过通道，用于重要操作）
	err := auditDB.WithContext(ctx).Create(auditLog).Error
	if err != nil {
		logger.Error(fmt.Sprintf("保存自定义审计日志失败: %v", err))
		return err
	}

	return nil
}

// AuditLogin 记录登录操作
func AuditLogin(ctx context.Context, userID int64, username string, tenantID int64, ipAddress, userAgent string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}

	actionDetail := models.ActionDetailJSON{
		Before:  nil,
		After:   map[string]interface{}{"username": username, "tenant_id": tenantID},
		Changes: nil,
	}

	auditLog := &models.AuditLog{
		TenantID:     tenantID,
		ModuleCode:   "SYSTEM",
		ActionType:   models.AuditActionLogin,
		EntityType:   "USER",
		EntityID:     &userID,
		ActionDetail: actionDetail,
		OperatorID:   userID,
		OperatorName: username,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Status:       status,
		CreatedAt:    time.Now(),
	}

	// 发送到通道
	select {
	case auditChan <- auditLog:
	default:
		logger.Error("审计日志通道已满，丢弃登录日志")
	}
}

// AuditLogout 记录登出操作
func AuditLogout(ctx context.Context, userID int64, username string, tenantID int64, ipAddress, userAgent string) {
	actionDetail := models.ActionDetailJSON{
		Before:  map[string]interface{}{"username": username, "tenant_id": tenantID},
		After:   nil,
		Changes: nil,
	}

	auditLog := &models.AuditLog{
		TenantID:     tenantID,
		ModuleCode:   "SYSTEM",
		ActionType:   models.AuditActionLogout,
		EntityType:   "USER",
		EntityID:     &userID,
		ActionDetail: actionDetail,
		OperatorID:   userID,
		OperatorName: username,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Status:       "success",
		CreatedAt:    time.Now(),
	}

	// 发送到通道
	select {
	case auditChan <- auditLog:
	default:
		logger.Error("审计日志通道已满，丢弃登出日志")
	}
}

// CloseAuditMiddleware 关闭审计日志中间件
func CloseAuditMiddleware() {
	if auditChan != nil {
		close(auditChan)
	}
}
