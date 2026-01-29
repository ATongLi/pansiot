package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"pansiot-cloud/internal/middleware"
	"pansiot-cloud/internal/service"
	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/response"
)

type AuthController struct {
	authService *service.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户，支持创建新企业或加入已有企业
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=service.RegisterResponse}
// @Router /api/v1/auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误: "+err.Error())
		return
	}

	// 调用服务层注册
	result, err := ctrl.authService.Register(c.Request.Context(), &req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		logger.Error(fmt.Sprintf("用户注册失败: %v", err))
		if err == service.ErrUserExists {
			response.Error(c, response.UserExistsCode, err.Error())
			response.Error(c, 400, "企业序列码无效")
		}
		return
	}

	// 生成Token（待完善）
	// accessToken, refreshToken, err := middleware.GenerateTokenPair(result.UserID, result.Username, result.TenantID, "TERMINAL")
	// if err != nil {
	// 	logger.Error("生成Token失败", "error", err)
	// 	response.Error(c, 500, "登录失败")
	// 	return
	// }
	// result.Token = accessToken
	// result.RefreshToken = refreshToken

	// 记录审计日志（待完善）
	// middleware.AuditCustomOperation(c.Request.Context(), result.TenantID, result.UserID, result.Username, "USER_MANAGEMENT", "CREATE", "USER", &result.UserID, nil, result, nil)

	response.Success(c, result)
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名/邮箱/手机号登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=service.LoginResponse}
// @Router /api/v1/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误: "+err.Error())
		return
	}

	// 调用服务层登录
	result, err := ctrl.authService.Login(c.Request.Context(), &req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		logger.Error(fmt.Sprintf("用户登录失败: account=%s, ip=%s, error=%v", req.Username, c.ClientIP(), err))
		if err != nil {
			response.Error(c, 401, "账号或密码错误")
		}
		return
	}	// Token已在service层生成，直接返回
	// 记录登录审计日志
	middleware.AuditLogin(c.Request.Context(), result.UserID, result.Username, result.TenantID, c.ClientIP(), c.Request.UserAgent(), true)

	response.Success(c, result)
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Description 使用RefreshToken获取新的AccessToken
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "RefreshToken"
// @Success 200 {object} response.Response{data=object{token=string,refresh_token=string}}
// @Router /api/v1/auth/refresh [post]
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	// 解析RefreshToken
	claims, err := middleware.ParseToken(req.RefreshToken)
	if err != nil {
		response.Error(c, 401, "RefreshToken无效或已过期")
		return
	}

	// 验证Token类型
	if claims.TokenType != "refresh" {
		response.Error(c, 401, "Token类型错误")
		return
	}

	// 检查用户是否存在
	user, err := ctrl.authService.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		response.Error(c, 401, "用户不存在")
		return
	}

	if user.Status != "ACTIVE" {
		response.Error(c, 401, "用户已被禁用")
		return
	}

	// 生成新的Token对
	accessToken, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Username, user.TenantID, claims.TenantType)
	if err != nil {
		logger.Error(fmt.Sprintf("生成Token失败: user_id=%d, error=%v", user.ID, err))
		response.Error(c, 500, "Token刷新失败")
		return
	}

	// 将旧的RefreshToken加入黑名单
	middleware.BlacklistToken(req.RefreshToken)

	response.Success(c, gin.H{
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 登出并失效Token
// @Tags 认证
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
	// 获取Token（如果提供）
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// 解析Bearer Token
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token := authHeader[7:]
			// 将Token加入黑名单
			middleware.BlacklistToken(token)
		}
	}

	// 获取用户信息
	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	tenantID := middleware.GetTenantID(c)

	// 记录登出审计日志
	if userID > 0 {
		middleware.AuditLogout(c.Request.Context(), userID, username, tenantID, c.ClientIP(), c.Request.UserAgent())
	}

	response.Success(c, gin.H{
		"message": "登出成功",
	})
}

// GetCurrentUser 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取登录用户的详细信息
// @Tags 认证
// @Security Bearer
// @Success 200 {object} response.Response{data=object{user_id=int64,username=string,email=string,tenant=object}}
// @Router /api/v1/auth/me [get]
func (ctrl *AuthController) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取用户信息
	user, err := ctrl.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		logger.Error(fmt.Sprintf("获取用户信息失败: user_id=%d, error=%v", userID, err))
		response.Error(c, 500, "获取用户信息失败")
		return
	}

	// 获取租户信息（待实现）
	// tenant, err := ctrl.tenantService.GetTenantByID(c.Request.Context(), user.TenantID)

	// 获取用户角色（待实现）
	// roles, err := ctrl.roleService.GetUserRoles(c.Request.Context(), userID)

	response.Success(c, gin.H{
		"user_id":    user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"phone":      user.Phone,
		"real_name":  user.RealName,
		"avatar":     user.Avatar,
		"tenant_id":  user.TenantID,
		"status":     user.Status,
		"created_at": user.CreatedAt,
		"last_login_at": user.LastLoginAt,
	})
}

// SendVerificationCode 发送验证码
// @Summary 发送验证码
// @Description 发送手机或邮箱验证码
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{account=string,type=string} true "账号信息"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/send-code [post]
func (ctrl *AuthController) SendVerificationCode(c *gin.Context) {
	var req struct {
		Account string `json:"account" binding:"required"` // 手机号或邮箱
		Type    string `json:"type" binding:"required,oneof=phone email"` // phone或email
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	// TODO: 实现验证码发送逻辑
	// 1. 验证账号格式
	// 2. 生成6位随机验证码
	// 3. 调用短信/邮件服务发送
	// 4. 将验证码存储到Redis（5分钟有效期）

	logger.Info(fmt.Sprintf("发送验证码: account=%s, type=%s, ip=%s", req.Account, req.Type, c.ClientIP()))

	// Mock实现：返回成功
	response.Success(c, gin.H{
		"message":       "验证码已发送",
		"expire_seconds": 300,
	})
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 通过验证码重置密码
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{account=string,code=string,new_password=string} true "重置信息"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/reset-password [post]
func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req struct {
		Account     string `json:"account" binding:"required"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	// TODO: 实现密码重置逻辑
	// 1. 从Redis获取验证码并验证
	// 2. 查找用户
	// 3. 更新密码（使用bcrypt加密）
	// 4. 清除验证码

	logger.Info(fmt.Sprintf("重置密码: account=%s, ip=%s", req.Account, c.ClientIP()))

	// Mock实现：返回成功
	response.Success(c, gin.H{
		"message": "密码重置成功",
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 已登录用户修改密码
// @Tags 认证
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body object{old_password=string,new_password=string} true "密码信息"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/change-password [post]
func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	// TODO: 实现密码修改逻辑
	// 1. 获取用户信息
	// 2. 验证旧密码
	// 3. 更新新密码（使用bcrypt加密）
	// 4. 返回成功

	logger.Info(fmt.Sprintf("修改密码: user_id=%d, ip=%s", userID, c.ClientIP()))

	response.Success(c, gin.H{
		"message": "密码修改成功",
	})
}
