package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`              // 业务状态码
	Message string      `json:"message"`           // 消息
	Data    interface{} `json:"data,omitempty"`    // 数据
	Meta    interface{} `json:"meta,omitempty"`    // 元数据（分页等）
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Total       int64       `json:"total"`        // 总数
	Page        int         `json:"page"`         // 当前页
	PageSize    int         `json:"page_size"`    // 每页数量
	TotalPages  int         `json:"total_pages"`  // 总页数
	Data        interface{} `json:"data"`         // 数据
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPagination 分页成功响应
func SuccessWithPagination(c *gin.Context, total int64, page, pageSize int, data interface{}) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Meta: PaginationResponse{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 错误响应（带数据）
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalServerError 500错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, 500, message)
}

// 业务错误码定义
const (
	SuccessCode        = 200     // 成功
	ErrorCode          = 400     // 通用错误
	UnauthorizedCode   = 401     // 未授权
	ForbiddenCode      = 403     // 禁止访问
	NotFoundCode       = 404     // 未找到
	InternalErrorCode  = 500     // 内部错误
	InvalidParamsCode  = 10001   // 参数错误
	InvalidCaptchaCode = 10002   // 验证码错误
	UserExistsCode     = 10003   // 用户已存在
	UserNotFoundCode   = 10004   // 用户不存在
	PasswordErrorCode  = 10005   // 密码错误
	TokenExpiredCode   = 10006   // Token过期
	TokenInvalidCode   = 10007   // Token无效
	PermissionDenyCode = 10008   // 权限不足
	QuotaExceededCode  = 10009   // 配额超限
)
