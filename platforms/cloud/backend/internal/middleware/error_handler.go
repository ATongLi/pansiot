package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/response"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 捕获panic
				handlePanic(c, err)
			}
		}()

		c.Next()

		// 处理请求过程中的错误
		if len(c.Errors) > 0 {
			handleErrors(c)
		}
	}
}

// handlePanic 处理panic
func handlePanic(c *gin.Context, err interface{}) {
	// 记录panic堆栈信息
	stack := debug.Stack()
	logger.Error(fmt.Sprintf("PANIC RECOVERED: %v\n%s", err, string(stack)))

	// 返回友好错误信息
	response.Error(c, 500, "服务器内部错误，请稍后重试")
	c.Abort()
}

// handleErrors 处理错误
func handleErrors(c *gin.Context) {
	// 获取第一个错误
	err := c.Errors.Last().Err

	// 根据错误类型返回不同的响应
	var apiError *APIError
	if errors.As(err, &apiError) {
		// 自定义API错误
		response.Error(c, apiError.Code, apiError.Message)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 记录不存在
		response.Error(c, 404, "记录不存在")
	} else if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// 参数验证错误
		handleValidationError(c, validationErrors)
	} else {
		// 其他错误
		logger.Error(fmt.Sprintf("Unhandled error: %v", err))
		response.Error(c, 500, "服务器内部错误")
	}

	c.Abort()
}

// handleValidationError 处理验证错误
func handleValidationError(c *gin.Context, errors validator.ValidationErrors) {
	messages := make([]string, 0, len(errors))
	for _, e := range errors {
		messages = append(messages, formatValidationError(e))
	}

	response.Error(c, response.InvalidParamsCode, strings.Join(messages, "; "))
}

// formatValidationError 格式化验证错误
func formatValidationError(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s不能为空", field)
	case "email":
		return fmt.Sprintf("%s格式不正确", field)
	case "min":
		return fmt.Sprintf("%s不能小于%s", field, param)
	case "max":
		return fmt.Sprintf("%s不能大于%s", field, param)
	case "len":
		return fmt.Sprintf("%s长度必须为%s", field, param)
	case "gte":
		return fmt.Sprintf("%s必须大于或等于%s", field, param)
	case "lte":
		return fmt.Sprintf("%s必须小于或等于%s", field, param)
	default:
		return fmt.Sprintf("%s验证失败: %s", field, tag)
	}
}

// APIError 自定义API错误
type APIError struct {
	Code    int
	Message string
	Err     error
}

// Error 实现error接口
func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 实现错误包装
func (e *APIError) Unwrap() error {
	return e.Err
}

// NewAPIError 创建API错误
func NewAPIError(code int, message string, err error) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 常用错误构造函数
func BadRequest(message string, err error) *APIError {
	return NewAPIError(400, message, err)
}

func Unauthorized(message string) *APIError {
	return NewAPIError(401, message, nil)
}

func Forbidden(message string) *APIError {
	return NewAPIError(403, message, nil)
}

func NotFound(message string) *APIError {
	return NewAPIError(404, message, nil)
}

func Conflict(message string, err error) *APIError {
	return NewAPIError(409, message, err)
}

func InternalServer(message string, err error) *APIError {
	return NewAPIError(500, message, err)
}

// BusinessError 业务错误
type BusinessError struct {
	Code    int
	Message string
	Details map[string]interface{}
}

// Error 实现error接口
func (e *BusinessError) Error() string {
	return e.Message
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message string, details map[string]interface{}) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// 常用业务错误
var (
	ErrUserExists      = NewBusinessError(response.UserExistsCode, "用户已存在", nil)
	ErrUserNotFound    = NewBusinessError(response.UserNotFoundCode, "用户不存在", nil)
	ErrPermissionDeny  = NewBusinessError(response.PermissionDenyCode, "权限不足", nil)
	ErrQuotaExceeded   = NewBusinessError(response.QuotaExceededCode, "配额已用完", nil)
	ErrInvalidParams   = NewBusinessError(response.InvalidParamsCode, "参数错误", nil)
	ErrInvalidPassword = NewBusinessError(10002, "密码错误", nil)
	ErrTokenExpired    = NewBusinessError(10003, "Token已过期", nil)
	ErrTokenInvalid    = NewBusinessError(10004, "Token无效", nil)
)

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取或生成请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 设置请求ID到上下文和响应头
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简单实现：使用时间戳和随机数
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GetRequestID 获取请求ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// RecoveryMiddleware 自定义Recovery中间件（替代Gin默认的）
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)
				stack := debug.Stack()

				logger.Error(fmt.Sprintf("RequestID: %s, PANIC: %v\n%s", requestID, err, string(stack)))

				// 检查是否是已断开的连接
				if isBrokenConnection(c) {
					c.Abort()
					return
				}

				response.Error(c, 500, "服务器内部错误")
				c.Abort()
			}
		}()

		c.Next()
	}
}

// isBrokenConnection 检查是否是断开的连接
func isBrokenConnection(c *gin.Context) bool {
	// 检查响应是否已经写入
	if !c.Writer.Written() {
		return false
	}

	// 检查连接状态
	if c.Writer.Status() == http.StatusEarlyHints {
		return false
	}

	return true
}
